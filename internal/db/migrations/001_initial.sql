-- +goose Up
CREATE TABLE cards
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        STRING NOT NULL,
    description STRING NOT NULL,
    section     STRING NOT NULL,
    UNIQUE (section, name)
);

-- +goose Down
DROP TABLE cards;