-- +goose Up
CREATE TABLE rondas_poker (
    proyecto_id  text NOT NULL REFERENCES proyectos(id),
    item_numero  integer NOT NULL,
    numero       integer NOT NULL CHECK (numero > 0),
    estado       text NOT NULL CHECK (estado IN ('Abierta', 'Revelada')),
    PRIMARY KEY (proyecto_id, item_numero, numero),
    CONSTRAINT rondas_poker_item_fk
        FOREIGN KEY (proyecto_id, item_numero)
        REFERENCES items_backlog (proyecto_id, numero)
);

CREATE UNIQUE INDEX rondas_poker_historia_abierta_idx
    ON rondas_poker (proyecto_id, item_numero)
    WHERE estado = 'Abierta';

CREATE TABLE votos_poker (
    proyecto_id       text NOT NULL,
    item_numero       integer NOT NULL,
    ronda_numero      integer NOT NULL,
    integrante_nombre text NOT NULL,
    valor             integer NOT NULL CHECK (valor IN (1, 2, 3, 5, 8, 13, 21)),
    PRIMARY KEY (proyecto_id, item_numero, ronda_numero, integrante_nombre),
    CONSTRAINT votos_poker_ronda_fk
        FOREIGN KEY (proyecto_id, item_numero, ronda_numero)
        REFERENCES rondas_poker (proyecto_id, item_numero, numero)
);

-- +goose Down
DROP TABLE votos_poker;
DROP TABLE rondas_poker;
