package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

type Material struct {
    ID                 int     `json:"id"`
    Nome               string  `json:"nome"`
    Unidade            string  `json:"unidade"`
    EstoqueAtual       float64 `json:"estoque_atual"`
    EstoqueMinimo      float64 `json:"estoque_minimo"`
    EstoqueEmergencial float64 `json:"estoque_emergencial"`
    PontoRecompra      float64 `json:"ponto_recompra"`   // ⬅️ novo campo
    LeadTimeDias       int     `json:"lead_time_dias"`
}


var db *sql.DB

func main() {
    os.MkdirAll("data", os.ModePerm)
    dbConn, err := sql.Open("sqlite3", "./data/estoque.db")
    if err != nil {
        log.Fatal(err)
    }
    db = dbConn
    initDatabase()

    http.HandleFunc("/api/materiais", listarMateriais)
    http.HandleFunc("/api/materiais/criar", criarMaterial)
    http.Handle("/", http.FileServer(http.Dir("public")))

    fmt.Println("Servidor rodando em http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func initDatabase() {
    sqlBytes, err := os.ReadFile("db/init.sql")
    if err != nil {
        log.Fatal("Erro ao ler init.sql:", err)
    }
    _, err = db.Exec(string(sqlBytes))
    if err != nil {
        log.Fatal("Erro ao criar tabelas:", err)
    }
}

func listarMateriais(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, nome, unidade, estoque_atual, estoque_minimo, estoque_emergencial, ponto_recompra, lead_time_dias FROM materiais")

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var materiais []Material
    for rows.Next() {
        var m Material
        err := rows.Scan(&m.ID, &m.Nome, &m.Unidade, &m.EstoqueAtual, &m.EstoqueMinimo, &m.EstoqueEmergencial, &m.PontoRecompra, &m.LeadTimeDias)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        materiais = append(materiais, m)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(materiais)
}

func criarMaterial(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    var m Material
    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    pontoRecompra := m.EstoqueMinimo + m.EstoqueEmergencial
    m.PontoRecompra = pontoRecompra

    result, err := db.Exec(`INSERT INTO materiais 
    (nome, unidade, estoque_atual, estoque_minimo, estoque_emergencial, ponto_recompra, lead_time_dias)
    VALUES (?, ?, ?, ?, ?, ?, ?)`,
    m.Nome, m.Unidade, m.EstoqueAtual, m.EstoqueMinimo, m.EstoqueEmergencial, m.PontoRecompra, m.LeadTimeDias)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    id, err := result.LastInsertId()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    m.ID = int(id)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(m)
}

