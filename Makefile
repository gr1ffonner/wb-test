POSTGRES_DSN = "postgres://admin:admin@localhost:5432/wb_orders?sslmode=disable"
MIGRATION_DIR_PG = ./migrations/postgres/
APP_DIR= ./cmd/app
PRODUCER_DIR= ./cmd/producer

.PHONY: run up up-dev down migrate-up-pg migrate-down-pg migrate-status-pg migrate-create-pg test test-jwt test-verbose

run:
	@export $$(grep -v '^#' ./.env | xargs) >/dev/null 2>&1; \
	go run $(APP_DIR)/main.go

run-producer:
	@export $$(grep -v '^#' ./.env | xargs) >/dev/null 2>&1; \
	go run $(PRODUCER_DIR)/main.go

up: 
	COMPOSE_PROJECT_NAME=garage docker compose -f docker-compose.yml --profile=test up -d --build 

# Start the infrastructure with Docker Compose
up-dev: 
	COMPOSE_PROJECT_NAME=wb-test docker compose -f docker-compose.yml --env-file=.env up -d --build

# Stop the services with Docker Compose
down: 
	COMPOSE_PROJECT_NAME=wb-test docker compose -f docker-compose.yml --profile=test down

# Run all tests
test:
	@echo "Running all tests..."
	go test ./... -v

# Run JWT tests specifically
test-jwt:
	@echo "Running JWT tests..."
	go test ./tests -v -run "TestJWT"

# Run tests with verbose output and coverage
test-verbose:
	@echo "Running tests with verbose output and coverage..."
	go test ./... -v -cover

.PHONY: migrate-up
migrate-up:
	@echo "Applying PostgreSQL migrations..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back last PostgreSQL migration..."
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) down

.PHONY: migrate-status
migrate-status:
	@echo "PostgreSQL migration status:"
	goose -dir $(MIGRATION_DIR_PG) postgres $(POSTGRES_DSN) status

.PHONY: migrate-create
migrate-create:
	@read -p "Enter PostgreSQL migration name: " NAME; \
	goose -dir $(MIGRATION_DIR_PG) create $$NAME sql
