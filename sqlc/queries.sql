-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, first_name, last_name, email;

-- name: GetUserOnSignIn :one
SELECT id, first_name, last_name, email, password FROM users WHERE email = $1;
