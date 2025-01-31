package handlers

import (
	"context"

	"github.com/jakub-szewczyk/career-compass-gin/db"
)

type Handler struct {
	ctx     context.Context
	queries *db.Queries
}

func NewHandler(ctx context.Context, queries *db.Queries) *Handler {
	return &Handler{
		ctx:     ctx,
		queries: queries,
	}
}
