CREATE TABLE IF NOT EXISTS materiais (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nome TEXT NOT NULL,
    unidade TEXT NOT NULL,
    estoque_atual REAL NOT NULL,
    estoque_minimo REAL NOT NULL,
    estoque_emergencial REAL NOT NULL,
    ponto_recompra REAL NOT NULL,
    lead_time_dias INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS movimentacoes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    material_id INTEGER NOT NULL,
    tipo TEXT CHECK(tipo IN ('entrada', 'saida')),
    quantidade REAL NOT NULL,
    data_movimentacao DATETIME DEFAULT CURRENT_TIMESTAMP,
    observacao TEXT,
    FOREIGN KEY(material_id) REFERENCES materiais(id)
);
