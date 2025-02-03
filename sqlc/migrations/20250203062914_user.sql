-- +goose Up
ALTER TABLE users ADD COLUMN first_name TEXT NOT NULL;

ALTER TABLE users ADD COLUMN last_name TEXT NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN first_name;

ALTER TABLE users DROP COLUMN last_name;
