package tests

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jakub-szewczyk/career-compass-gin/api/handlers"
	"github.com/jakub-szewczyk/career-compass-gin/api/routes"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var R *gin.Engine

var Ctx context.Context
var Queries *db.Queries

func TestMain(m *testing.M) {
	Ctx = context.Background()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to obtain the current working directory: %s", err)
	}

	postgresContainer, err := postgres.Run(Ctx, "postgres", postgres.WithInitScripts(filepath.Join(dir, "..", "..", "sqlc", "schema.sql")), postgres.WithDatabase("career-compass-gin-test"), postgres.WithUsername("career-compass-gin-test"), postgres.WithPassword("career-compass-gin-test"), testcontainers.WithWaitStrategy(
		wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	databaseURL, err := postgresContainer.ConnectionString(Ctx)
	if err != nil {
		log.Fatalf("failed to obtain container connection string: %s", err)
	}

	conn, err := pgx.Connect(Ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer conn.Close(Ctx)

	Queries = db.New(conn)

	R = routes.Setup(Ctx, handlers.NewEnv("", databaseURL, ""), Queries)

	code := m.Run()

	os.Exit(code)
}
