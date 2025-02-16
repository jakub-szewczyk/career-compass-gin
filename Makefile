run:
	go run main.go

air:
	air

test:
	go test ./api/tests

test-v:
	go test ./api/tests -v

slqc:
	sqlc -f sqlc/sqlc.yaml generate

migrate-dev:
	goose create $(MIGRATION) sql --env .env.dev

migrate-prod:
	goose create $(MIGRATION) sql --env .env.prod

up-dev:
	goose up -env .env.dev

up-prod:
	goose up -env .env.prod

psql-dev:
	docker exec -it postgres-dev psql -U career-compass-gin-dev

psql-prod:
	docker exec -it postgres-prod psql -U career-compass-gin-prod

swag-fmt:
	swag fmt

swag-init: swag-fmt
	swag init -g api/routes/routes.go

help:
	@echo "Usage: make [target]\n"
	@echo "Targets:"
	@echo "  run                    - Starts local dev server"
	@echo "  air                    - Starts local dev server with auto-reload"
	@echo "  test                   - Runs tests"
	@echo "  test-v                 - Runs tests in verbose mode"
	@echo "  sqlc                   - Generate source code from SQL"
	@echo "  migrate-dev            - Create dev DB migration file"
	@echo "  migrate-prod           - Create prod DB migration file"
	@echo "  up-dev                 - Migrate dev DB to the most recent version available"
	@echo "  up-prod                - Migrate prod DB to the most recent version available"
	@echo "  psqlc-dev              - Start dev DB interactive shell"
	@echo "  psqlc-prod             - Start prod DB interactive shell"
	@echo "  swag-fmt               - Format swag comments"
	@echo "  swag-init              - Generate docs.go"
