package handlers

import (
	"context"

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

	FrontendURL string
}

func NewEnv(port, databaseURL, jwtSecret, smtpIdentity, smtpUsername, smtpPassword, smtpHost, smtpPort, frontendURL string) Env {
	return Env{
		Port:        port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,

		SMTPIdentity: smtpIdentity,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		SMTPHost:     smtpHost,
		SMTPPORT:     smtpPort,

		FrontendURL: frontendURL,
	}
}

type Handler struct {
	ctx     context.Context
	env     Env
	queries *db.Queries
}

func NewHandler(ctx context.Context, env Env, queries *db.Queries) *Handler {
	return &Handler{
		ctx:     ctx,
		env:     env,
		queries: queries,
	}
}
