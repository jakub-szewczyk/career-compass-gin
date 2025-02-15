-- name: Purge :exec
TRUNCATE TABLE users, verification_tokens;

-- name: CreateUser :one
WITH new_user AS (
  INSERT INTO users (first_name, last_name, email, password)
  VALUES (
    UPPER(LEFT(sqlc.arg(first_name)::text, 1)) || LOWER(SUBSTRING(sqlc.arg(first_name)::text FROM 2)),
    UPPER(LEFT(sqlc.arg(last_name)::text, 1)) || LOWER(SUBSTRING(sqlc.arg(last_name)::text FROM 2)),
    sqlc.arg(email)::text,
    sqlc.arg(password)::text
  )
  RETURNING id, first_name, last_name, email, is_email_verified
),
new_token AS (
  INSERT INTO verification_tokens (user_id) SELECT id FROM new_user RETURNING token
)
SELECT 
  new_user.id,
  new_user.first_name,
  new_user.last_name,
  new_user.email,
  new_user.is_email_verified,
  new_token.token as verification_token
FROM new_user, new_token;

-- name: GetUserOnSignIn :one
SELECT u.id, u.first_name, u.last_name, u.email, u.password, u.is_email_verified, v.token AS verification_token
FROM users AS u
JOIN verification_tokens AS v ON u.id = v.user_id
WHERE email = $1;

-- name: GetUserById :one
SELECT u.id, u.first_name, u.last_name, u.email, u.is_email_verified, v.token AS verification_token
FROM users AS u
JOIN verification_tokens AS v ON u.id = v.user_id
WHERE u.id = $1;
