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

var r *gin.Engine
var ctx context.Context
var token string
var queries *db.Queries

func setupUser(ctx context.Context) {
	queries.CreateUser(ctx, db.CreateUserParams{FirstName: "Jakub", LastName: "Szewczyk", Email: "jakub.szewczyk@test.com", Password: "qwerty!123456789"})
}

func TestMain(m *testing.M) {
	ctx = context.Background()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to obtain the current working directory: %s", err)
	}

	postgresContainer, err := postgres.Run(ctx, "postgres", postgres.WithInitScripts(filepath.Join(dir, "..", "..", "sqlc", "schema.sql")), postgres.WithDatabase("career-compass-gin-test"), postgres.WithUsername("career-compass-gin-test"), postgres.WithPassword("career-compass-gin-test"), testcontainers.WithWaitStrategy(
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

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to obtained mapped port: %s", err)
	}

	databaseURL, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to obtain container connection string: %s", err)
	}

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer conn.Close(ctx)

	queries = db.New(conn)

	r = routes.Setup(ctx, handlers.NewEnv(port.Port(), databaseURL, "testing", "", "", "", "", "", "", "", ""), queries)

	code := m.Run()

	os.Exit(code)
}
