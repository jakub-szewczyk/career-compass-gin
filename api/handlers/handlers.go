package handlers

import (
	"context"

	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

type Env struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func NewEnv(port, databaseURL, jwtSecret string) Env {
	return Env{
		Port:        port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
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
