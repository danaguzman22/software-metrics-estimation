-- +goose Up
ALTER TABLE items_backlog
    ADD COLUMN horas_estimadas numeric
    CHECK (horas_estimadas >= 0 AND horas_estimadas = round(horas_estimadas, 2));

ALTER TABLE esfuerzos
    ADD COLUMN sprint_numero integer;

ALTER TABLE esfuerzos
    ADD CONSTRAINT esfuerzos_sprint_fk
    FOREIGN KEY (proyecto_id, sprint_numero)
    REFERENCES sprints (proyecto_id, numero);

CREATE TABLE sprint_esfuerzo_snapshots (
    proyecto_id   text NOT NULL,
    sprint_numero integer NOT NULL,
    PRIMARY KEY (proyecto_id, sprint_numero),
    FOREIGN KEY (proyecto_id, sprint_numero)
        REFERENCES sprints (proyecto_id, numero)
);

CREATE TABLE sprint_esfuerzo_items (
    proyecto_id      text NOT NULL,
    sprint_numero    integer NOT NULL,
    item_numero      integer NOT NULL,
    horas_estimadas  numeric,
    PRIMARY KEY (proyecto_id, sprint_numero, item_numero),
    FOREIGN KEY (proyecto_id, sprint_numero)
        REFERENCES sprint_esfuerzo_snapshots (proyecto_id, sprint_numero),
    FOREIGN KEY (proyecto_id, item_numero)
        REFERENCES items_backlog (proyecto_id, numero)
);

-- +goose Down
DROP TABLE sprint_esfuerzo_items;
DROP TABLE sprint_esfuerzo_snapshots;
ALTER TABLE esfuerzos DROP CONSTRAINT esfuerzos_sprint_fk;
ALTER TABLE esfuerzos DROP COLUMN sprint_numero;
ALTER TABLE items_backlog DROP COLUMN horas_estimadas;
