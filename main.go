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

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	r := routes.Setup(ctx, handlers.NewEnv(port, databaseURL, jwtSecret), queries)

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
