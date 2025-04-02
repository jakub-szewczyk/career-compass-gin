-- name: Purge :exec
TRUNCATE TABLE users, verification_tokens, password_reset_tokens, job_applications;

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

-- name: GetPasswordResetToken :one
SELECT token, expires_at, user_id FROM password_reset_tokens WHERE token = $1;

-- name: UpdatePassword :exec
UPDATE users SET password = $2 WHERE id = $1;

-- name: DeletePasswordResetToken :exec
DELETE FROM password_reset_tokens WHERE token = $1;

-- name: GetJobApplications :many
WITH user_job_applications AS (
  SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url
  FROM job_applications
  WHERE user_id = $3
), 
filtered_job_applications AS (
  SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, COUNT(*) OVER() AS total
  FROM user_job_applications
  WHERE 
    (company_name ILIKE '%' || @company_name_or_job_title::text || '%' OR job_title ILIKE '%' || @company_name_or_job_title::text || '%' OR @company_name_or_job_title::text IS NULL)
    AND (((date_applied AT TIME ZONE 'Europe/Warsaw')::date = ((@date_applied) AT TIME ZONE 'Europe/Warsaw')::date) OR @date_applied IS NULL)
    AND (status = nullif(@status, '')::status OR nullif(@status, '') IS NULL)
  ORDER BY
    CASE WHEN @company_name_asc::bool THEN company_name END ASC,
    CASE WHEN @company_name_desc::bool THEN company_name END DESC,
    CASE WHEN @job_title_asc::bool THEN job_title END ASC,
    CASE WHEN @job_title_desc::bool THEN job_title END DESC,
    CASE WHEN @date_applied_asc::bool THEN date_applied END ASC,
    CASE WHEN @date_applied_desc::bool THEN date_applied END DESC,
    CASE WHEN @status_asc::bool THEN status END ASC,
    CASE WHEN @status_desc::bool THEN status END DESC,
    CASE WHEN @salary_asc::bool THEN greatest(min_salary, max_salary) END ASC,
    CASE WHEN @salary_desc::bool THEN greatest(min_salary, max_salary) END DESC,
    CASE WHEN @is_replied_asc::bool THEN is_replied END ASC,
    CASE WHEN @is_replied_desc::bool THEN is_replied END DESC
)
SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, total
FROM filtered_job_applications
LIMIT $1 OFFSET $2;

-- name: GetJobApplication :one
SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, notes FROM job_applications WHERE id = $1 AND user_id = $2;

-- name: CreateJobApplication :one
INSERT INTO job_applications (user_id, company_name, job_title, date_applied, status, min_salary, max_salary, job_posting_url, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, notes;

-- name: UpdateJobApplication :one
UPDATE job_applications
SET 
  company_name = coalesce(nullif(sqlc.narg('company_name'), ''), company_name),
  job_title = coalesce(nullif(sqlc.narg('job_title'), ''), job_title),
  date_applied = coalesce(sqlc.narg('date_applied'), date_applied),
  status = coalesce(sqlc.narg('status'), status),
  is_replied = coalesce(sqlc.narg('is_replied'), is_replied),
  min_salary = coalesce(sqlc.narg('min_salary'), min_salary),
  max_salary = coalesce(sqlc.narg('max_salary'), max_salary),
  job_posting_url = coalesce(nullif(sqlc.narg('job_posting_url'), ''), job_posting_url),
  notes = coalesce(nullif(sqlc.narg('notes'), ''), notes)
WHERE id = $1 AND user_id = $2
RETURNING id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, notes;

-- name: DeleteJobApplication :one
DELETE FROM job_applications WHERE id = $1 AND user_id = $2
RETURNING id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, notes;
