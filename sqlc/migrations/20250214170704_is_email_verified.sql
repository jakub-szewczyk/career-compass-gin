-- +goose Up
ALTER TABLE users ADD COLUMN is_email_verified BOOLEAN DEFAULT false;

-- +goose Down
ALTER TABLE users DROP COLUMN is_email_verified;
