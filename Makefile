POSTGRES_DSN = "postgres://admin:admin@localhost:5432/garage?sslmode=disable"
MIGRATION_DIR_PG = ./migrations/postgres/
APP_DIR= ./cmd/app

.PHONY: run up up-dev down migrate-up-pg migrate-down-pg migrate-status-pg migrate-create-pg 

run:
	@export $$(grep -v '^#' ./.env | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

up: 
	COMPOSE_PROJECT_NAME=garage docker compose -f docker-compose.yml --profile=test up -d --build 

# Start the infrastructure with Docker Compose
up-dev: 
	COMPOSE_PROJECT_NAME=wb-test docker compose -f docker-compose.yml --env-file=.env up -d --build

# Stop the services with Docker Compose
down: 
	COMPOSE_PROJECT_NAME=wb-test docker compose -f docker-compose.yml --profile=test down


.PHONY: migrate-up-pg
migrate-up-pg:
	@echo "Applying PostgreSQL migrations..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) up

.PHONY: migrate-down-pg
migrate-down-pg:
	@echo "Rolling back last PostgreSQL migration..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) down

.PHONY: migrate-status-pg
migrate-status-pg:
	@echo "PostgreSQL migration status:"
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) status

.PHONY: migrate-create-pg
migrate-create-pg:
	@read -p "Enter PostgreSQL migration name: " NAME; \
	goose -dir $(MIGRATION_DIR_PG) create $$NAME sql
