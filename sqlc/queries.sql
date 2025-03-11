-- name: Purge :exec
TRUNCATE TABLE users, verification_tokens, password_reset_tokens;

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
SELECT id, first_name, last_name, email, password, is_email_verified FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT id, first_name, last_name, email, is_email_verified FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT u.id, u.first_name, u.last_name, u.email, u.is_email_verified, v.token as verification_token
FROM users AS u
JOIN verification_tokens as v ON u.id = v.user_id
WHERE u.email = $1;

-- name: GetVerificationToken :one
SELECT token, expires_at FROM verification_tokens WHERE user_id = $1;

-- name: VerifyEmail :one
UPDATE users SET is_email_verified = true WHERE id = $1 RETURNING id, first_name, last_name, email, is_email_verified;

-- name: UpdateVerificationToken :one
UPDATE verification_tokens SET
  token = encode(gen_random_bytes(32), 'hex'),
  expires_at = NOW() + INTERVAL '1 day'
WHERE user_id = $1
RETURNING token, expires_at;

-- name: ExpireVerificationToken :exec
UPDATE verification_tokens SET expires_at = NOW() - INTERVAL '1 day' WHERE user_id = $1;

-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id)
VALUES ($1)
ON CONFLICT (user_id)
DO UPDATE SET token = encode(gen_random_bytes(32), 'hex'), expires_at = NOW() + INTERVAL '15 minutes'
RETURNING token;
