package handlers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

type Env struct {
	Port        string
	DatabaseURL string
	JWTSecret   string

	SMTPIdentity string
	SMTPUsername string
	SMTPPassword string
	SMTPHost     string
	SMTPPORT     string

	FrontendURL          string
	EmailVerificationURL string
	ResetPasswordURL     string
}

func NewEnv(port, databaseURL, jwtSecret, smtpIdentity, smtpUsername, smtpPassword, smtpHost, smtpPort, frontendURL, emailVerificationURL, resetPasswordURL string) Env {
	return Env{
		Port:        port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,

		SMTPIdentity: smtpIdentity,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		SMTPHost:     smtpHost,
		SMTPPORT:     smtpPort,

		FrontendURL:          frontendURL,
		EmailVerificationURL: emailVerificationURL,
		ResetPasswordURL:     resetPasswordURL,
	}
}

type Handler struct {
	ctx     context.Context
	env     Env
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewHandler(ctx context.Context, env Env, queries *db.Queries, pool *pgxpool.Pool) *Handler {
	return &Handler{
		ctx:     ctx,
		env:     env,
		queries: queries,
		pool:    pool,
	}
}
