package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/handlers"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

func Setup(ctx context.Context, env handlers.Env, queries *db.Queries) *gin.Engine {
	r := gin.Default()

	h := handlers.NewHandler(ctx, env, queries)

	api := r.Group("/api")

	// NOTE: Public routes
	api.GET("/health-check", h.HealthCheck)

	api.POST("/sign-up", h.SignUp)
	api.POST("/sign-in", h.SignIn)

	// NOTE: Private routes
	api.Use(h.Auth())

	api.GET("/profile", h.Profile)

	return r
}
