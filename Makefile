# UserCenter Makefile

# Variables
APP_NAME := usercenter
BUILD_DIR := bin
GO_FILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet

# Binary names
BINARY_NAME := $(APP_NAME)
BINARY_UNIX := $(BINARY_NAME)_unix

# Default target
.PHONY: all
all: clean deps lint test build

# Help
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: deps
deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: wire
wire: ## Generate Wire dependency injection code
	@echo "Generating Wire code..."
	wire ./cmd/usercenter

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	swag init -g cmd/usercenter/main.go -o docs

.PHONY: mock
mock: ## Generate mock files
	@echo "Generating mocks..."
	mockgen -source=internal/service/user.go -destination=internal/mock/user_service_mock.go
	mockgen -source=internal/repository/user.go -destination=internal/mock/user_repository_mock.go

##@ Building

.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/usercenter

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) ./cmd/usercenter

.PHONY: build-all
build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/usercenter
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/usercenter
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/usercenter

##@ Running

.PHONY: run
run: ## Run the application
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run ./cmd/usercenter

.PHONY: run-dev
run-dev: wire swagger ## Run in development mode with code generation
	@echo "Running in development mode..."
	$(GOCMD) run ./cmd/usercenter

##@ Testing

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GOTEST) -v -race ./...

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -v -bench=. -benchmem ./...

##@ Code Quality

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w $(GO_FILES)

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run --config .golangci.yml

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	gosec ./...

##@ Database

.PHONY: migrate-create
migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	@echo "Creating migration: $(name)"
	goose -dir migrations create $(name) sql

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@echo "Applying migrations..."
	goose -dir migrations postgres "$(shell grep -A 5 'postgres:' configs/config.yaml | grep -E '(host|port|user|password|dbname)' | awk '{print $$2}' | tr -d '"' | paste -sd ' ' | awk '{print "host=" $$1 " port=" $$2 " user=" $$3 " password=" $$4 " dbname=" $$5 " sslmode=disable"}')" up

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@echo "Rolling back migration..."
	goose -dir migrations postgres "$(shell grep -A 5 'postgres:' configs/config.yaml | grep -E '(host|port|user|password|dbname)' | awk '{print $$2}' | tr -d '"' | paste -sd ' ' | awk '{print "host=" $$1 " port=" $$2 " user=" $$3 " password=" $$4 " dbname=" $$5 " sslmode=disable"}')" down

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "Migration status..."
	goose -dir migrations postgres "$(shell grep -A 5 'postgres:' configs/config.yaml | grep -E '(host|port|user|password|dbname)' | awk '{print $$2}' | tr -d '"' | paste -sd ' ' | awk '{print "host=" $$1 " port=" $$2 " user=" $$3 " password=" $$4 " dbname=" $$5 " sslmode=disable"}')" status

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(APP_NAME):latest

.PHONY: docker-compose-up
docker-compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services with docker-compose..."
	docker-compose down

##@ Cleanup

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: clean-cache
clean-cache: ## Clean Go module cache
	@echo "Cleaning module cache..."
	$(GOCMD) clean -modcache

##@ Tools Installation

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/golang/mock/mockgen@latest

##@ Information

.PHONY: version
version: ## Show version
	@echo "Version: $(VERSION)"

.PHONY: info
info: ## Show project information
	@echo "Application: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build date: $(shell date -u +%Y-%m-%dT%H:%M:%SZ)" 