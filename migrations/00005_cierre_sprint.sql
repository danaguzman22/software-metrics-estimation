-- +goose Up
ALTER TABLE items_backlog
    DROP CONSTRAINT items_backlog_estado_check;
ALTER TABLE items_backlog
    ADD CONSTRAINT items_backlog_estado_check
    CHECK (estado IN ('Pendiente', 'Completada'));

ALTER TABLE sprints
    DROP CONSTRAINT sprints_estado_check;
ALTER TABLE sprints
    ADD CONSTRAINT sprints_estado_check
    CHECK (estado IN ('Activo', 'Cerrado'));

ALTER TABLE sprints
    ADD COLUMN sp_planificados integer;

-- +goose Down
ALTER TABLE items_backlog
    DROP CONSTRAINT items_backlog_estado_check;
ALTER TABLE items_backlog
    ADD CONSTRAINT items_backlog_estado_check
    CHECK (estado IN ('Pendiente'));

ALTER TABLE sprints
    DROP CONSTRAINT sprints_estado_check;
ALTER TABLE sprints
    ADD CONSTRAINT sprints_estado_check
    CHECK (estado IN ('Activo'));

ALTER TABLE sprints
    DROP COLUMN sp_planificados;
