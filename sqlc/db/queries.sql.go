// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createJobApplication = `-- name: CreateJobApplication :one
INSERT INTO job_applications (user_id, company_name, job_title, date_applied, status, min_salary, max_salary, job_posting_url, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, notes
`

type CreateJobApplicationParams struct {
	UserID        pgtype.UUID        `json:"userId"`
	CompanyName   string             `json:"companyName"`
	JobTitle      string             `json:"jobTitle"`
	DateApplied   pgtype.Timestamptz `json:"dateApplied"`
	Status        Status             `json:"status"`
	MinSalary     pgtype.Float8      `json:"minSalary"`
	MaxSalary     pgtype.Float8      `json:"maxSalary"`
	JobPostingUrl pgtype.Text        `json:"jobPostingUrl"`
	Notes         pgtype.Text        `json:"notes"`
}

type CreateJobApplicationRow struct {
	ID            pgtype.UUID        `json:"id"`
	CompanyName   string             `json:"companyName"`
	JobTitle      string             `json:"jobTitle"`
	DateApplied   pgtype.Timestamptz `json:"dateApplied"`
	Status        Status             `json:"status"`
	IsReplied     bool               `json:"isReplied"`
	MinSalary     pgtype.Float8      `json:"minSalary"`
	MaxSalary     pgtype.Float8      `json:"maxSalary"`
	JobPostingUrl pgtype.Text        `json:"jobPostingUrl"`
	Notes         pgtype.Text        `json:"notes"`
}

func (q *Queries) CreateJobApplication(ctx context.Context, arg CreateJobApplicationParams) (CreateJobApplicationRow, error) {
	row := q.db.QueryRow(ctx, createJobApplication,
		arg.UserID,
		arg.CompanyName,
		arg.JobTitle,
		arg.DateApplied,
		arg.Status,
		arg.MinSalary,
		arg.MaxSalary,
		arg.JobPostingUrl,
		arg.Notes,
	)
	var i CreateJobApplicationRow
	err := row.Scan(
		&i.ID,
		&i.CompanyName,
		&i.JobTitle,
		&i.DateApplied,
		&i.Status,
		&i.IsReplied,
		&i.MinSalary,
		&i.MaxSalary,
		&i.JobPostingUrl,
		&i.Notes,
	)
	return i, err
}

const createPasswordResetToken = `-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id)
VALUES ($1)
ON CONFLICT (user_id)
DO UPDATE SET token = encode(gen_random_bytes(32), 'hex'), expires_at = NOW() + INTERVAL '15 minutes'
RETURNING token
`

func (q *Queries) CreatePasswordResetToken(ctx context.Context, userID pgtype.UUID) (string, error) {
	row := q.db.QueryRow(ctx, createPasswordResetToken, userID)
	var token string
	err := row.Scan(&token)
	return token, err
}

const createUser = `-- name: CreateUser :one
WITH new_user AS (
  INSERT INTO users (first_name, last_name, email, password)
  VALUES (
    UPPER(LEFT($1::text, 1)) || LOWER(SUBSTRING($1::text FROM 2)),
    UPPER(LEFT($2::text, 1)) || LOWER(SUBSTRING($2::text FROM 2)),
    $3::text,
    $4::text
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
FROM new_user, new_token
`

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type CreateUserRow struct {
	ID                pgtype.UUID `json:"id"`
	FirstName         string      `json:"firstName"`
	LastName          string      `json:"lastName"`
	Email             string      `json:"email"`
	IsEmailVerified   pgtype.Bool `json:"isEmailVerified"`
	VerificationToken string      `json:"verificationToken"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.IsEmailVerified,
		&i.VerificationToken,
	)
	return i, err
}

const deletePasswordResetToken = `-- name: DeletePasswordResetToken :exec
DELETE FROM password_reset_tokens WHERE token = $1
`

func (q *Queries) DeletePasswordResetToken(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, deletePasswordResetToken, token)
	return err
}

const expireVerificationToken = `-- name: ExpireVerificationToken :exec
UPDATE verification_tokens SET expires_at = NOW() - INTERVAL '1 day' WHERE user_id = $1
`

func (q *Queries) ExpireVerificationToken(ctx context.Context, userID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, expireVerificationToken, userID)
	return err
}

const getJobApplication = `-- name: GetJobApplication :one
SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, notes FROM job_applications WHERE id = $1 AND user_id = $2
`

type GetJobApplicationParams struct {
	ID     pgtype.UUID `json:"id"`
	UserID pgtype.UUID `json:"userId"`
}

type GetJobApplicationRow struct {
	ID            pgtype.UUID        `json:"id"`
	CompanyName   string             `json:"companyName"`
	JobTitle      string             `json:"jobTitle"`
	DateApplied   pgtype.Timestamptz `json:"dateApplied"`
	Status        Status             `json:"status"`
	IsReplied     bool               `json:"isReplied"`
	MinSalary     pgtype.Float8      `json:"minSalary"`
	MaxSalary     pgtype.Float8      `json:"maxSalary"`
	JobPostingUrl pgtype.Text        `json:"jobPostingUrl"`
	Notes         pgtype.Text        `json:"notes"`
}

func (q *Queries) GetJobApplication(ctx context.Context, arg GetJobApplicationParams) (GetJobApplicationRow, error) {
	row := q.db.QueryRow(ctx, getJobApplication, arg.ID, arg.UserID)
	var i GetJobApplicationRow
	err := row.Scan(
		&i.ID,
		&i.CompanyName,
		&i.JobTitle,
		&i.DateApplied,
		&i.Status,
		&i.IsReplied,
		&i.MinSalary,
		&i.MaxSalary,
		&i.JobPostingUrl,
		&i.Notes,
	)
	return i, err
}

const getJobApplications = `-- name: GetJobApplications :many
WITH user_job_applications AS (
  SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url
  FROM job_applications
  WHERE user_id = $3
), 
filtered_job_applications AS (
  SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, COUNT(*) OVER() AS total
  FROM user_job_applications
  WHERE 
    (company_name ILIKE '%' || $4::text || '%' OR job_title ILIKE '%' || $4::text || '%' OR $4::text IS NULL)
    AND (((date_applied AT TIME ZONE 'Europe/Warsaw')::date = (($5) AT TIME ZONE 'Europe/Warsaw')::date) OR $5 IS NULL)
    AND (status = nullif($6, '')::status OR nullif($6, '') IS NULL)
  ORDER BY
    CASE WHEN $7::bool THEN company_name END ASC,
    CASE WHEN $8::bool THEN company_name END DESC,
    CASE WHEN $9::bool THEN job_title END ASC,
    CASE WHEN $10::bool THEN job_title END DESC,
    CASE WHEN $11::bool THEN date_applied END ASC,
    CASE WHEN $12::bool THEN date_applied END DESC,
    CASE WHEN $13::bool THEN status END ASC,
    CASE WHEN $14::bool THEN status END DESC,
    CASE WHEN $15::bool THEN greatest(min_salary, max_salary) END ASC,
    CASE WHEN $16::bool THEN greatest(min_salary, max_salary) END DESC,
    CASE WHEN $17::bool THEN is_replied END ASC,
    CASE WHEN $18::bool THEN is_replied END DESC
)
SELECT id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, total
FROM filtered_job_applications
LIMIT $1 OFFSET $2
`

type GetJobApplicationsParams struct {
	Limit                 int32       `json:"limit"`
	Offset                int32       `json:"offset"`
	UserID                pgtype.UUID `json:"userId"`
	CompanyNameOrJobTitle string      `json:"companyNameOrJobTitle"`
	DateApplied           interface{} `json:"dateApplied"`
	Status                interface{} `json:"status"`
	CompanyNameAsc        bool        `json:"companyNameAsc"`
	CompanyNameDesc       bool        `json:"companyNameDesc"`
	JobTitleAsc           bool        `json:"jobTitleAsc"`
	JobTitleDesc          bool        `json:"jobTitleDesc"`
	DateAppliedAsc        bool        `json:"dateAppliedAsc"`
	DateAppliedDesc       bool        `json:"dateAppliedDesc"`
	StatusAsc             bool        `json:"statusAsc"`
	StatusDesc            bool        `json:"statusDesc"`
	SalaryAsc             bool        `json:"salaryAsc"`
	SalaryDesc            bool        `json:"salaryDesc"`
	IsRepliedAsc          bool        `json:"isRepliedAsc"`
	IsRepliedDesc         bool        `json:"isRepliedDesc"`
}

type GetJobApplicationsRow struct {
	ID            pgtype.UUID        `json:"id"`
	CompanyName   string             `json:"companyName"`
	JobTitle      string             `json:"jobTitle"`
	DateApplied   pgtype.Timestamptz `json:"dateApplied"`
	Status        Status             `json:"status"`
	IsReplied     bool               `json:"isReplied"`
	MinSalary     pgtype.Float8      `json:"minSalary"`
	MaxSalary     pgtype.Float8      `json:"maxSalary"`
	JobPostingUrl pgtype.Text        `json:"jobPostingUrl"`
	Total         int64              `json:"total"`
}

func (q *Queries) GetJobApplications(ctx context.Context, arg GetJobApplicationsParams) ([]GetJobApplicationsRow, error) {
	rows, err := q.db.Query(ctx, getJobApplications,
		arg.Limit,
		arg.Offset,
		arg.UserID,
		arg.CompanyNameOrJobTitle,
		arg.DateApplied,
		arg.Status,
		arg.CompanyNameAsc,
		arg.CompanyNameDesc,
		arg.JobTitleAsc,
		arg.JobTitleDesc,
		arg.DateAppliedAsc,
		arg.DateAppliedDesc,
		arg.StatusAsc,
		arg.StatusDesc,
		arg.SalaryAsc,
		arg.SalaryDesc,
		arg.IsRepliedAsc,
		arg.IsRepliedDesc,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetJobApplicationsRow
	for rows.Next() {
		var i GetJobApplicationsRow
		if err := rows.Scan(
			&i.ID,
			&i.CompanyName,
			&i.JobTitle,
			&i.DateApplied,
			&i.Status,
			&i.IsReplied,
			&i.MinSalary,
			&i.MaxSalary,
			&i.JobPostingUrl,
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPasswordResetToken = `-- name: GetPasswordResetToken :one
SELECT token, expires_at, user_id FROM password_reset_tokens WHERE token = $1
`

type GetPasswordResetTokenRow struct {
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
	UserID    pgtype.UUID        `json:"userId"`
}

func (q *Queries) GetPasswordResetToken(ctx context.Context, token string) (GetPasswordResetTokenRow, error) {
	row := q.db.QueryRow(ctx, getPasswordResetToken, token)
	var i GetPasswordResetTokenRow
	err := row.Scan(&i.Token, &i.ExpiresAt, &i.UserID)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT u.id, u.first_name, u.last_name, u.email, u.is_email_verified, v.token as verification_token
FROM users AS u
JOIN verification_tokens as v ON u.id = v.user_id
WHERE u.email = $1
`

type GetUserByEmailRow struct {
	ID                pgtype.UUID `json:"id"`
	FirstName         string      `json:"firstName"`
	LastName          string      `json:"lastName"`
	Email             string      `json:"email"`
	IsEmailVerified   pgtype.Bool `json:"isEmailVerified"`
	VerificationToken string      `json:"verificationToken"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.IsEmailVerified,
		&i.VerificationToken,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, first_name, last_name, email, is_email_verified FROM users WHERE id = $1
`

type GetUserByIdRow struct {
	ID              pgtype.UUID `json:"id"`
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Email           string      `json:"email"`
	IsEmailVerified pgtype.Bool `json:"isEmailVerified"`
}

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.IsEmailVerified,
	)
	return i, err
}

const getUserOnSignIn = `-- name: GetUserOnSignIn :one
SELECT id, first_name, last_name, email, password, is_email_verified FROM users WHERE email = $1
`

type GetUserOnSignInRow struct {
	ID              pgtype.UUID `json:"id"`
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Email           string      `json:"email"`
	Password        string      `json:"password"`
	IsEmailVerified pgtype.Bool `json:"isEmailVerified"`
}

func (q *Queries) GetUserOnSignIn(ctx context.Context, email string) (GetUserOnSignInRow, error) {
	row := q.db.QueryRow(ctx, getUserOnSignIn, email)
	var i GetUserOnSignInRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.IsEmailVerified,
	)
	return i, err
}

const getVerificationToken = `-- name: GetVerificationToken :one
SELECT token, expires_at FROM verification_tokens WHERE user_id = $1
`

type GetVerificationTokenRow struct {
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
}

func (q *Queries) GetVerificationToken(ctx context.Context, userID pgtype.UUID) (GetVerificationTokenRow, error) {
	row := q.db.QueryRow(ctx, getVerificationToken, userID)
	var i GetVerificationTokenRow
	err := row.Scan(&i.Token, &i.ExpiresAt)
	return i, err
}

const purge = `-- name: Purge :exec
TRUNCATE TABLE users, verification_tokens, password_reset_tokens, job_applications
`

func (q *Queries) Purge(ctx context.Context) error {
	_, err := q.db.Exec(ctx, purge)
	return err
}

const updateJobApplication = `-- name: UpdateJobApplication :one
UPDATE job_applications
SET 
  company_name = COALESCE($3, company_name),
  job_title = COALESCE($4, job_title),
  date_applied = COALESCE($5, date_applied),
  status = COALESCE($6, status),
  is_replied = COALESCE($7, is_replied),
  min_salary = COALESCE($8, min_salary),
  max_salary = COALESCE($9, max_salary),
  job_posting_url = COALESCE($10, job_posting_url),
  notes = COALESCE($11, notes)
WHERE id = $1 AND user_id = $2
RETURNING id, company_name, job_title, date_applied, status, is_replied, min_salary, max_salary, job_posting_url, notes
`

type UpdateJobApplicationParams struct {
	ID            pgtype.UUID        `json:"id"`
	UserID        pgtype.UUID        `json:"userId"`
	CompanyName   pgtype.Text        `json:"companyName"`
	JobTitle      pgtype.Text        `json:"jobTitle"`
	DateApplied   pgtype.Timestamptz `json:"dateApplied"`
	Status        NullStatus         `json:"status"`
	IsReplied     pgtype.Bool        `json:"isReplied"`
	MinSalary     pgtype.Float8      `json:"minSalary"`
	MaxSalary     pgtype.Float8      `json:"maxSalary"`
	JobPostingUrl pgtype.Text        `json:"jobPostingUrl"`
	Notes         pgtype.Text        `json:"notes"`
}

type UpdateJobApplicationRow struct {
	ID            pgtype.UUID        `json:"id"`
	CompanyName   string             `json:"companyName"`
	JobTitle      string             `json:"jobTitle"`
	DateApplied   pgtype.Timestamptz `json:"dateApplied"`
	Status        Status             `json:"status"`
	IsReplied     bool               `json:"isReplied"`
	MinSalary     pgtype.Float8      `json:"minSalary"`
	MaxSalary     pgtype.Float8      `json:"maxSalary"`
	JobPostingUrl pgtype.Text        `json:"jobPostingUrl"`
	Notes         pgtype.Text        `json:"notes"`
}

func (q *Queries) UpdateJobApplication(ctx context.Context, arg UpdateJobApplicationParams) (UpdateJobApplicationRow, error) {
	row := q.db.QueryRow(ctx, updateJobApplication,
		arg.ID,
		arg.UserID,
		arg.CompanyName,
		arg.JobTitle,
		arg.DateApplied,
		arg.Status,
		arg.IsReplied,
		arg.MinSalary,
		arg.MaxSalary,
		arg.JobPostingUrl,
		arg.Notes,
	)
	var i UpdateJobApplicationRow
	err := row.Scan(
		&i.ID,
		&i.CompanyName,
		&i.JobTitle,
		&i.DateApplied,
		&i.Status,
		&i.IsReplied,
		&i.MinSalary,
		&i.MaxSalary,
		&i.JobPostingUrl,
		&i.Notes,
	)
	return i, err
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users SET password = $2 WHERE id = $1
`

type UpdatePasswordParams struct {
	ID       pgtype.UUID `json:"id"`
	Password string      `json:"password"`
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error {
	_, err := q.db.Exec(ctx, updatePassword, arg.ID, arg.Password)
	return err
}

const updateVerificationToken = `-- name: UpdateVerificationToken :one
UPDATE verification_tokens SET
  token = encode(gen_random_bytes(32), 'hex'),
  expires_at = NOW() + INTERVAL '1 day'
WHERE user_id = $1
RETURNING token, expires_at
`

type UpdateVerificationTokenRow struct {
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
}

func (q *Queries) UpdateVerificationToken(ctx context.Context, userID pgtype.UUID) (UpdateVerificationTokenRow, error) {
	row := q.db.QueryRow(ctx, updateVerificationToken, userID)
	var i UpdateVerificationTokenRow
	err := row.Scan(&i.Token, &i.ExpiresAt)
	return i, err
}

const verifyEmail = `-- name: VerifyEmail :one
UPDATE users SET is_email_verified = true WHERE id = $1 RETURNING id, first_name, last_name, email, is_email_verified
`

type VerifyEmailRow struct {
	ID              pgtype.UUID `json:"id"`
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Email           string      `json:"email"`
	IsEmailVerified pgtype.Bool `json:"isEmailVerified"`
}

func (q *Queries) VerifyEmail(ctx context.Context, id pgtype.UUID) (VerifyEmailRow, error) {
	row := q.db.QueryRow(ctx, verifyEmail, id)
	var i VerifyEmailRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.IsEmailVerified,
	)
	return i, err
}
