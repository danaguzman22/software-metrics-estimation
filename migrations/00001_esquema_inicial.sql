-- +goose Up
CREATE TABLE proyectos (
    id           text PRIMARY KEY DEFAULT gen_random_uuid()::text,
    nombre       text NOT NULL,
    fecha_inicio date NOT NULL,
    fecha_fin    date NOT NULL,
    estado       text NOT NULL CHECK (estado IN ('Planificado', 'En curso', 'Finalizado'))
);

CREATE TABLE integrantes (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id text NOT NULL REFERENCES proyectos(id),
    nombre      text NOT NULL,
    rol         text NOT NULL
);

CREATE UNIQUE INDEX integrantes_proyecto_nombre_idx
    ON integrantes (proyecto_id, lower(nombre));

CREATE TABLE items_backlog (
    proyecto_id  text NOT NULL REFERENCES proyectos(id),
    numero       integer NOT NULL,
    titulo       text NOT NULL,
    descripcion  text NOT NULL,
    prioridad    text NOT NULL CHECK (prioridad IN ('Alta', 'Media', 'Baja')),
    estado       text NOT NULL CHECK (estado IN ('Pendiente')),
    story_points integer NOT NULL CHECK (story_points IN (0, 1, 2, 3, 5, 8, 13, 21)),
    criterios_aceptacion text[] NOT NULL DEFAULT '{}',
    PRIMARY KEY (proyecto_id, numero)
);

-- +goose Down
DROP TABLE items_backlog;
DROP TABLE integrantes;
DROP TABLE proyectos;
