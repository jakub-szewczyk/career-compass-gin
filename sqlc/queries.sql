-- name: Purge :exec
TRUNCATE TABLE users;

-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password)
VALUES (
  UPPER(LEFT(sqlc.arg(first_name)::text, 1)) || LOWER(SUBSTRING(sqlc.arg(first_name)::text FROM 2)),
  UPPER(LEFT(sqlc.arg(last_name)::text, 1)) || LOWER(SUBSTRING(sqlc.arg(last_name)::text FROM 2)),
  sqlc.arg(email)::text,
  sqlc.arg(password)::text
) RETURNING id, first_name, last_name, email;

-- name: GetUserOnSignIn :one
SELECT id, first_name, last_name, email, password FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT id, first_name, last_name, email FROM users where id = $1;
