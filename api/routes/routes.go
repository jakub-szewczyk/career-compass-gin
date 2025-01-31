package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/handlers"
	"github.com/jakub-szewczyk/career-compass-gin/db"
)

func Setup(ctx context.Context, queries *db.Queries) *gin.Engine {
	r := gin.Default()

	h := handlers.NewHandler(ctx, queries)

	r.GET("/health-check", h.HealthCheck)

	r.GET("/users", h.GetAllUsers)

	return r
}
