-- +goose Up
CREATE TABLE sprint_items (
    proyecto_id   text NOT NULL,
    sprint_numero integer NOT NULL,
    item_numero   integer NOT NULL,
    PRIMARY KEY (proyecto_id, item_numero),
    CONSTRAINT sprint_items_sprint_fk
        FOREIGN KEY (proyecto_id, sprint_numero)
        REFERENCES sprints (proyecto_id, numero),
    CONSTRAINT sprint_items_item_fk
        FOREIGN KEY (proyecto_id, item_numero)
        REFERENCES items_backlog (proyecto_id, numero)
);

CREATE INDEX sprint_items_proyecto_sprint_idx
    ON sprint_items (proyecto_id, sprint_numero);

-- +goose Down
DROP TABLE sprint_items;
