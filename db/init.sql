CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS pessoas (
    id UUID PRIMARY KEY,
    apelido VARCHAR(32) NOT NULL UNIQUE,
    nome VARCHAR(100) NOT NULL,
    nascimento DATE NOT NULL,
    stack JSONB,
    busca TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pessoas_busca_trgm
ON pessoas USING GIN (busca gin_trgm_ops);
