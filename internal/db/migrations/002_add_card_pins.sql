-- +goose Up
ALTER TABLE cards ADD COLUMN pinned BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE cards DROP COLUMN pinned;
