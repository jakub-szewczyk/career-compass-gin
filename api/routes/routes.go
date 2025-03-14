package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/handlers"
	"github.com/jakub-szewczyk/career-compass-gin/docs"
	_ "github.com/jakub-szewczyk/career-compass-gin/docs"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title						Career Compass REST API
// @BasePath					/api
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func Setup(ctx context.Context, env handlers.Env, queries *db.Queries) *gin.Engine {
	// TODO: Read from env vars
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + env.Port

	r := gin.Default()

	h := handlers.NewHandler(ctx, env, queries)

	api := r.Group("/api")

	// NOTE: Public routes
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.GET("/health-check", h.HealthCheck)

	api.POST("/sign-up", h.SignUp)
	api.POST("/sign-in", h.SignIn)

	api.POST("/password/reset", h.InitPasswordReset)
	api.PUT("/password/reset", h.ResetPassword)

	// NOTE: Private routes
	api.Use(h.Auth())

	api.GET("/profile", h.Profile)

	api.GET("/profile/verify-email", h.SendVerificationEmail)
	api.PATCH("/profile/verify-email", h.VerifyEmail)

	api.GET("/job-applications", h.JobApplications)
	api.GET("/job-applications/:jobApplicationId", h.JobApplication)
	api.POST("/job-applications", h.CreateJobApplication)
	api.PUT("/job-applications/:jobApplicationId", h.UpdateJobApplication)
	api.DELETE("/job-applications/:jobApplicationId", h.DeleteJobApplication)

	return r
}
