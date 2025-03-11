package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jakub-szewczyk/career-compass-gin/api/handlers"
	"github.com/jakub-szewczyk/career-compass-gin/api/routes"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.dev")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing env var: PORT")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("missing env var: DATABASE_URL")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing env var: JWT_SECRET")
	}

	smtpIdentity := os.Getenv("SMTP_IDENTITY")
	if smtpIdentity == "" {
		log.Fatal("missing env var: SMTP_IDENTITY")
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	if smtpUsername == "" {
		log.Fatal("missing env var: SMTP_USERNAME")
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		log.Fatal("missing env var: SMTP_PASSWORD")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Fatal("missing env var: SMTP_HOST")
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		log.Fatal("missing env var: SMTP_PORT")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatal("missing env var: FRONTEND_URL")
	}

	emailVerificationURL := os.Getenv("EMAIL_VERIFICATION_URL")
	if emailVerificationURL == "" {
		log.Fatal("missing env var: EMAIL_VERIFICATION_URL")
	}

	resetPasswordURL := os.Getenv("RESET_PASSWORD_URL")
	if resetPasswordURL == "" {
		log.Fatal("missing env var: RESET_PASSWORD_URL")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	r := routes.Setup(ctx, handlers.NewEnv(port, databaseURL, jwtSecret, smtpIdentity, smtpUsername, smtpPassword, smtpHost, smtpPort, frontendURL, emailVerificationURL, resetPasswordURL), queries)

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
