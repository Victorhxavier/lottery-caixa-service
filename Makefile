.PHONY: help build run test clean docker docker-push fmt lint

# Variáveis
BINARY_NAME=lottery-service
DOCKER_IMAGE=lottery-caixa-service
DOCKER_TAG=latest
GO_VERSION=1.21

help:
	@echo "Lottery Caixa Service - Available targets:"
	@echo "  make build          - Build the binary"
	@echo "  make run            - Run the service"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make docker         - Build Docker image"
	@echo "  make docker-run     - Run Docker image"
	@echo "  make dev            - Start development environment"
	@echo "  make deps           - Download dependencies"

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 go build -o $(BINARY_NAME) ./cmd/main.go
	@echo "Build complete: $(BINARY_NAME)"

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

run-dev:
	@echo "Running in development mode..."
	go run ./cmd/main.go

test:
	@echo "Running tests..."
	go test -v -race -timeout=30s ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean -cache -testcache
	@echo "Clean complete"

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

lint:
	@echo "Running linter..."
	golangci-lint run ./...

docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run:
	@echo "Running Docker container..."
	docker run -it --rm \
		-p 8080:8080 \
		-e DOWNSTREAM_SERVICE_URL=http://localhost:8081 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-push:
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

dev: deps
	@echo "Starting development environment..."
	docker-compose up -d
	@echo "Services running. Check docker-compose.yml for details"

dev-logs:
	docker-compose logs -f lottery-caixa-service

dev-down:
	@echo "Stopping development environment..."
	docker-compose down

check-go:
	@echo "Go version check..."
	go version

check-tools:
	@echo "Checking required tools..."
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

all: deps clean fmt lint test build
	@echo "All tasks completed successfully!"
