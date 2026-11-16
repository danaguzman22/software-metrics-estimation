-- +goose Up
CREATE TABLE esfuerzos (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id text NOT NULL REFERENCES proyectos(id),
    item_numero integer NOT NULL,
    integrante  text NOT NULL,
    fecha       date NOT NULL,
    actividad   text NOT NULL,
    horas       numeric(4,2) NOT NULL CHECK (horas > 0 AND horas <= 24),
    FOREIGN KEY (proyecto_id, item_numero)
        REFERENCES items_backlog (proyecto_id, numero)
);

-- +goose Down
DROP TABLE esfuerzos;
