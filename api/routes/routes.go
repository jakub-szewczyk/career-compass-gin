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

	// NOTE: Root routes
	api := r.Group("/api")

	api.GET("/health-check", h.HealthCheck)

	// NOTE: Auth routes
	api.POST("/sign-up", h.SignUp)
	api.POST("/sign-in", h.SignIn)

	return r
}
