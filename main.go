package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jakub-szewczyk/career-compass-gin/api/routes"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DATABASE_URL")

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	r := routes.Setup(ctx, queries)

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
