-- +goose Up
CREATE TABLE users (
  id       BIGSERIAL PRIMARY KEY,
  email    TEXT NOT NULL,
  password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
