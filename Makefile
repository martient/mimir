# Mimir Makefile

.PHONY: help build run test lint fmt tidy clean deps install

help: ## Display this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -v 'awk'

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the mimir binary
	@echo "Building mimir..."
	@go build -o mimir -ldflags="-s -w" ./cmd
	@echo "Build complete: ./mimir"

run: ## Run the mimir server
	go run ./cmd/main.go

test: ## Run tests
	go test -v -race -cover ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	@echo "Running golangci-lint..."
	golangci-lint run --config=.golangci.yml

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Done"

tidy: ## Tidy go.mod
	go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f mimir
	@rm -f coverage.out coverage.html
	@rm -rf mimir/
	@echo "Clean complete"

install: ## Install mimir globally
	@echo "Installing mimir..."
	@go install -ldflags="-s -w" ./cmd
	@echo "Installed to $(go env GOPATH)/bin/mimir"

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run ./scripts/migrate.go

generate: ## Generate code (wire, mockgen, etc.)
	@echo "Generating code..."
	@if command -v wire >/dev/null 2>&1; then \
		wire ./...; \
	else \
		echo "wire not installed. Install with: go install github.com/google/wire/cmd/wire@latest"; \
	fi
