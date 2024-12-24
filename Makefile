# Variables
APP_NAME=my-go-app
PORT=8081
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_CONTAINER=$(APP_NAME)-container

# Default target
.PHONY: all
all: build

# Build the Go application
.PHONY: build
build:
	@echo "Building the Go application..."
	go build -o main .

# Run the application locally
.PHONY: run
run:
	@echo "Running the Go application locally..."
	./main up

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -f main

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -d --name $(DOCKER_CONTAINER) -p $(PORT):8080 $(DOCKER_IMAGE)

# Stop Docker container
.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

# Push Docker image to registry
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image to registry..."
	docker push $(DOCKER_IMAGE)

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Lint the code
.PHONY: lint
lint:
	@echo "Linting the code..."
	golangci-lint run

# Format the code
.PHONY: format
format:
	@echo "Formatting the code..."
	go fmt ./...

# Full clean up (build artifacts and Docker)
.PHONY: full-clean
full-clean: clean docker-stop
	@echo "Full clean completed."

# Help command to list all available targets
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          Build the Go application"
	@echo "  run            Run the Go application locally"
	@echo "  clean          Clean up build artifacts"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container"
	@echo "  docker-stop    Stop Docker container"
	@echo "  docker-push    Push Docker image to registry"
	@echo "  test           Run tests"
	@echo "  lint           Lint the code"
	@echo "  format         Format the code"
	@echo "  full-clean     Perform a full cleanup"
	@echo "  help           Show this help message"

# Migrations commands
.PHONY: migration
migration:
	@migrate create -ext sql -dir db/migrate/migrations $(filter-out $!,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Running database migrations (up)..."
	@go run db/migrate/migration.go up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database migrations (down)..."
	@go run db/migrate/migration.go down

.PHONY: migration-force
migration-force:
	@migrate -path cmd/migrate/migrations -database "postgres://postgres:password@localhost:5432/your_database_name?sslmode=disable" force $(filter-out $!,$(MAKECMDGOALS))