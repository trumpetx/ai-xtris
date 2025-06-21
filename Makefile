# Xtris Clone Makefile
# Common development tasks for the Xtris game

.PHONY: help test test-coverage test-race build build-all clean lint fmt coverage-report coverage-html

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Testing
test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...

test-race: ## Run tests with race detection
	go test -race -v ./...

# Coverage analysis
coverage-report: test-coverage ## Generate coverage report
	@echo "Test Coverage Report:"
	@go tool cover -func=coverage.out
	@echo "\nCoverage Summary:"
	@go tool cover -func=coverage.out | tail -1

coverage-html: test-coverage ## Generate HTML coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

# Code quality
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

# Building
build: ## Build the game for current platform
	go build -o xtris main.go

build-race: ## Build with race detection
	go build -race -o xtris main.go

build-all: ## Build for all platforms (Windows, macOS, Linux)
	@echo "Building for Windows (AMD64)..."
	GOOS=windows GOARCH=amd64 go build -v -o xtris.exe main.go
	@echo "Building for macOS (AMD64)..."
	GOOS=darwin GOARCH=amd64 go build -v -o xtris-darwin-amd64 main.go
	@echo "Building for macOS (ARM64)..."
	GOOS=darwin GOARCH=arm64 go build -v -o xtris-darwin-arm64 main.go
	@echo "Building for Linux (AMD64)..."
	GOOS=linux GOARCH=amd64 go build -v -o xtris-linux-amd64 main.go
	@echo "Building for Linux (ARM64)..."
	GOOS=linux GOARCH=arm64 go build -v -o xtris-linux-arm64 main.go
	@echo "✅ All builds completed!"

build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 go build -v -o xtris.exe main.go

build-macos: ## Build for macOS (both AMD64 and ARM64)
	GOOS=darwin GOARCH=amd64 go build -v -o xtris-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -v -o xtris-darwin-arm64 main.go

build-linux: ## Build for Linux (both AMD64 and ARM64)
	GOOS=linux GOARCH=amd64 go build -v -o xtris-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -v -o xtris-linux-arm64 main.go

# Cleanup
clean: ## Clean build artifacts
	rm -f xtris xtris.exe xtris-*-amd64 xtris-*-arm64 coverage.out coverage.html

# Development workflow
dev: fmt lint test ## Run full development workflow (fmt, lint, test)

# Coverage threshold check (fails if coverage < 80%)
coverage-check: test-coverage ## Check if coverage meets minimum threshold (80%)
	@coverage=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$coverage >= 80" | bc -l) -eq 1 ]; then \
		echo "✅ Coverage: $$coverage% (meets 80% threshold)"; \
	else \
		echo "❌ Coverage: $$coverage% (below 80% threshold)"; \
		exit 1; \
	fi

# Install development tools
install-tools: ## Install development tools (golangci-lint, etc.)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Run the game
run: ## Run the game
	go run main.go

# Release preparation
release: build-all ## Build all platforms for release
	@echo "Creating release archive..."
	@mkdir -p release
	@cp xtris.exe release/ 2>/dev/null || true
	@cp xtris-*-amd64 release/ 2>/dev/null || true
	@cp xtris-*-arm64 release/ 2>/dev/null || true
	@echo "✅ Release builds ready in ./release/ directory" 