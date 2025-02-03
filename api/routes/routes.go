package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/handlers"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

func Setup(ctx context.Context, queries *db.Queries) *gin.Engine {
	r := gin.Default()

	h := handlers.NewHandler(ctx, queries)

	// NOTE: Root routes
	api := r.Group("/api")

	api.GET("/health-check", h.HealthCheck)

	// NOTE: Auth routes
	auth := api.Group("/auth")

	auth.POST("/sign-up", h.SignUp)
	// auth.POST("/sign-in", h.SignIn)

	return r
}
