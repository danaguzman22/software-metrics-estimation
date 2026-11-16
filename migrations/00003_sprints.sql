-- +goose Up
CREATE TABLE sprints (
    proyecto_id  text NOT NULL REFERENCES proyectos(id),
    numero       integer NOT NULL CHECK (numero > 0),
    objetivo     text NOT NULL,
    fecha_inicio date NOT NULL,
    fecha_fin    date NOT NULL,
    estado       text NOT NULL CHECK (estado IN ('Activo')),
    PRIMARY KEY (proyecto_id, numero),
    CHECK (fecha_fin >= fecha_inicio)
);

CREATE UNIQUE INDEX sprints_proyecto_activo_idx
    ON sprints (proyecto_id)
    WHERE estado = 'Activo';

-- +goose Down
DROP TABLE sprints;
