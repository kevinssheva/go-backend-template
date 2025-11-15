.PHONY: help run test lint build migrate-up migrate-down docker-up docker-down clean tidy

# Variables
BINARY_NAME=palmvue-server
MAIN_PATH=./cmd/server
MIGRATION_PATH=./db/migrations
DB_URL=postgres://postgres:postgres@localhost:5432/palmvue?sslmode=disable

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

run: ## Run the application
	@echo "Running application..."
	go run $(MAIN_PATH)/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Coverage report:"
	go tool cover -func=coverage.out

lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run ./...

build: ## Build the application
	@echo "Building application..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)/main.go
	@echo "Binary created at bin/$(BINARY_NAME)"

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: migration name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)..."
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

docker-up: ## Start Docker services
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 3

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

clean: ## Clean build artifacts and test cache
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out
	go clean -testcache

tidy: ## Tidy and verify go.mod
	@echo "Tidying go.mod..."
	go mod tidy
	go mod verify

install-tools: ## Install required development tools
	@echo "Installing development tools..."
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	@echo "Installing golang-migrate..."
	@which migrate > /dev/null || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed successfully!"

dev: docker-up ## Start development environment (Docker + Run app)
	@echo "Starting development environment..."
	@sleep 5
	$(MAKE) run

all: tidy lint test build ## Run tidy, lint, test, and build
