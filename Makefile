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

# Test coverage settings
COVERAGE_DIR := coverage
COVERAGE_OUT := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html
COVERAGE_XML := $(COVERAGE_DIR)/coverage.xml
COVERAGE_FUNC := $(COVERAGE_DIR)/coverage.txt

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
	mockgen -source=internal/service/user.go -destination=internal/mock/user_service_mock.go -package=mock
	mockgen -source=internal/repository/user.go -destination=internal/mock/user_repository_mock.go -package=mock
	mockgen -source=internal/service/auth.go -destination=internal/mock/auth_service_mock.go -package=mock

##@ Building

.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/usercenter

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

.PHONY: test-short
test-short: ## Run only short tests (unit tests)
	@echo "Running short tests..."
	$(GOTEST) -v -short ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	$(GOCMD) tool cover -func=$(COVERAGE_OUT) -o $(COVERAGE_FUNC)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@echo "Coverage summary: $(COVERAGE_FUNC)"

.PHONY: test-coverage-verbose
test-coverage-verbose: ## Run tests with detailed coverage
	@echo "Running tests with detailed coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_OUT) -covermode=atomic -coverpkg=./... ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	$(GOCMD) tool cover -func=$(COVERAGE_OUT) -o $(COVERAGE_FUNC)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@echo "Coverage summary: $(COVERAGE_FUNC)"
	@echo "Coverage percentage:"
	@tail -1 $(COVERAGE_FUNC)

.PHONY: test-coverage-xml
test-coverage-xml: ## Generate XML coverage report for CI/CD
	@echo "Generating XML coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...
	gocov convert $(COVERAGE_OUT) | gocov-xml > $(COVERAGE_XML)
	@echo "XML coverage report generated: $(COVERAGE_XML)"

.PHONY: test-race
test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GOTEST) -v -race ./...

.PHONY: test-benchmark
test-benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -v -bench=. -benchmem ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./...

.PHONY: test-all
test-all: test-short test-integration test-coverage ## Run all types of tests
	@echo "All tests completed"

.PHONY: test-watch
test-watch: ## Run tests in watch mode (requires air)
	@echo "Running tests in watch mode..."
	air -c .air.toml

.PHONY: test-kafka
test-kafka: ## Run Kafka integration tests
	@echo "Running Kafka integration tests..."
	$(GOTEST) -v -tags=kafka ./internal/service/...

.PHONY: test-kafka-manual
test-kafka-manual: ## Run manual Kafka test script
	@echo "Running manual Kafka test script..."
	@if [ -f "./scripts/test-kafka.sh" ]; then \
		chmod +x ./scripts/test-kafka.sh; \
		./scripts/test-kafka.sh; \
	else \
		echo "Kafka test script not found: ./scripts/test-kafka.sh"; \
		exit 1; \
	fi

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

.PHONY: docker-compose-logs
docker-compose-logs: ## Show docker-compose logs
	@echo "Showing docker-compose logs..."
	docker-compose logs -f

.PHONY: docker-compose-ps
docker-compose-ps: ## Show docker-compose service status
	@echo "Showing docker-compose service status..."
	docker-compose ps

.PHONY: setup-env
setup-env: ## Setup local environment variables
	@echo "Setting up local environment variables..."
	@if [ ! -f .env ]; then \
		echo "Creating .env file..."; \
		echo "USERCENTER_DATABASE_POSTGRES_HOST=localhost" > .env; \
		echo "USERCENTER_DATABASE_POSTGRES_PORT=5432" >> .env; \
		echo "USERCENTER_DATABASE_POSTGRES_USER=postgres" >> .env; \
		echo "USERCENTER_DATABASE_POSTGRES_PASSWORD=password" >> .env; \
		echo "USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter" >> .env; \
		echo "USERCENTER_DATABASE_POSTGRES_SSLMODE=disable" >> .env; \
		echo ".env file created successfully!"; \
	else \
		echo ".env file already exists"; \
	fi

.PHONY: dev-start
dev-start: setup-env docker-compose-up ## Start development environment
	@echo "Development environment started!"
	@echo "Next steps:"
	@echo "1. Run 'make run-dev' to start the application"
	@echo "2. Access the service at http://localhost:8080"
	@echo "3. View logs with 'make docker-compose-logs'"

.PHONY: dev-stop
dev-stop: docker-compose-down ## Stop development environment
	@echo "Development environment stopped!"

##@ Cleanup

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	rm -f coverage.out coverage.html coverage.xml

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
	go install github.com/axw/gocov/gocov@latest
	go install github.com/AlekSi/gocov-xml@latest
	go install github.com/air-verse/air@latest

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