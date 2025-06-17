CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    nome TEXT,
    limite INT NOT NULL,
    saldo INT NOT NULL DEFAULT 0
);

CREATE TABLE transacoes (
    id SERIAL PRIMARY KEY,
    cliente_id INT REFERENCES clientes(id),
    valor INT NOT NULL,
    tipo CHAR(1) NOT NULL CHECK (tipo IN ('c','d')),
    descricao TEXT NOT NULL,
    realizada_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO clientes (id, nome, limite, saldo) VALUES
(1, 'o barato sai caro', 100000, 0),
(2, 'zan corp ltda', 80000, 0),
(3, 'les cruders', 1000000, 0),
(4, 'padaria joia de cocaia', 10000000, 0),
(5, 'kid mais', 500000, 0);
