# UserCenter - User Management Service

[![Go Version](https://img.shields.io/badge/Go-1.23.1-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/Coverage-80%25-brightgreen.svg)](./coverage.html)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

**English** | [中文](README.zh-CN.md)

## 📖 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

---

## 🎯 Overview

UserCenter is a production-ready user management service built with Go, providing comprehensive user management capabilities including registration, authentication, querying, and listing. The project follows standard Go project layout and adopts modern technology stack to support high concurrency, high availability, and scalability.

### Core Features

- 🔐 **User Authentication**: JWT-based user registration and login
- 🔍 **User Query**: Conditional filtering for user information queries
- 📋 **User Listing**: Paginated and sortable user lists
- 🏥 **Health Checks**: Service status monitoring endpoints
- 🛡️ **Security Features**: Input validation, rate limiting, CORS support
- 🌍 **Internationalization**: Multi-language support (Chinese/English)
- 🔄 **Graceful Shutdown**: Safe service termination mechanism
- 📊 **Observability**: Complete monitoring, logging, and distributed tracing

## 🚀 Features

### Authentication & Authorization
- JWT-based stateless authentication
- Password hashing with bcrypt (cost 12)
- Role-based access control
- Token refresh mechanism
- Secure session management

### User Management
- User registration with email verification
- User profile management
- Account status management (active, inactive, suspended)
- Soft delete support
- Bulk user operations

### API Features
- RESTful API design
- Comprehensive input validation
- Rate limiting (general, login-specific, registration-specific)
- Request ID tracking
- CORS configuration
- Swagger/OpenAPI documentation

### Monitoring & Observability
- Health check endpoints for all dependencies
- Prometheus metrics collection
- Structured logging with Zap
- Distributed tracing with OpenTelemetry
- Performance monitoring

## 🛠️ Technology Stack

### Core Framework
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- **Dependency Injection**: [Wire](https://github.com/google/wire) - Compile-time dependency injection
- **API Documentation**: [Swagger](https://github.com/swaggo/gin-swagger) - Auto-generated OpenAPI 3.0 documentation

### Data Storage
- **Primary Database**: [PostgreSQL](https://www.postgresql.org/) + [GORM](https://gorm.io/) - User core data
- **Auxiliary Database**: [MongoDB](https://www.mongodb.com/) - Logs and session data
- **Cache**: [Redis](https://redis.io/) - High-performance caching
- **Database Migration**: [Goose](https://github.com/pressly/goose) - Database version control

### Message & Task Processing
- **Message Queue**: [Kafka](https://kafka.apache.org/) - Event consumption
- **Async Tasks**: [Asynq](https://github.com/hibiken/asynq) - Background task processing

### Monitoring & Logging
- **Logging**: [Zap](https://github.com/uber-go/zap) - High-performance structured logging
- **Monitoring**: [Prometheus](https://prometheus.io/) - Metrics collection
- **Distributed Tracing**: [OpenTelemetry](https://opentelemetry.io/) - Distributed tracing

### Security & Utilities
- **Authentication**: [JWT](https://github.com/golang-jwt/jwt) - Stateless authentication
- **Internationalization**: [go-i18n](https://github.com/nicksnyder/go-i18n) - Multi-language support
- **Configuration**: YAML configuration files
- **Code Quality**: [golangci-lint](https://golangci-lint.run/) - Code quality checks

## 📋 Prerequisites

### System Requirements
- Go 1.23.1 or higher
- PostgreSQL 13+
- MongoDB 5.0+
- Redis 6.0+
- Apache Kafka 2.8+

### Development Tools
```bash
# Install Go
# Reference: https://golang.org/doc/install

# Install Wire
go install github.com/google/wire/cmd/wire@latest

# Install Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# Install Swagger generator
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🚀 Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd user-center
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Configure Environment
```bash
# Copy configuration file
cp configs/config.example.yaml configs/config.yaml

# Edit configuration
vim configs/config.yaml
```

### 4. Initialize Database
```bash
# Run database migrations
make migrate-up

# Or manually
goose -dir migrations postgres "user=username password=password dbname=usercenter sslmode=disable" up
```

### 5. Generate Wire Dependency Injection Code
```bash
make wire
```

### 6. Generate Swagger Documentation
```bash
make swagger
```

### 7. Run the Service
```bash
# Development environment
make run

# Or run directly
go run cmd/usercenter/main.go

# Production environment
make build
./bin/usercenter
```

## ⚙️ Configuration

The project supports multiple configuration methods with priority from high to low:

1. **Environment Variables**: `USERCENTER_` prefix
2. **Configuration File**: `configs/config.yaml`
3. **Default Values**: Default configuration in code

### Main Configuration Items

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # debug, release, test

database:
  postgres:
    host: "localhost"
    port: 5432
    user: "username"
    password: "password"
    dbname: "usercenter"
    sslmode: "disable"
  
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "usercenter_logs"
  
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0

kafka:
  brokers: ["localhost:9092"]
  topics:
    user_events: "user.events"

jwt:
  secret: "your-secret-key"
  expiry: "24h"

logging:
  level: "info"
  format: "json"

monitoring:
  prometheus:
    enabled: true
    port: 9090
  
  tracing:
    enabled: true
    endpoint: "http://localhost:14268/api/traces"
```

## 📚 API Documentation

### Swagger Documentation Access

After starting the service, you can access the API documentation at:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI JSON**: http://localhost:8080/swagger/doc.json

### API Endpoints

#### 1. Health Check
```bash
# Basic health check
GET /health

# Detailed health check
GET /health/detailed

# Metrics endpoint
GET /metrics
```

#### 2. User Management
```bash
# User registration
POST /api/v1/users/register
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password"
}

# User login
POST /api/v1/users/login
{
  "email": "john@example.com",
  "password": "secure_password"
}

# Get user profile
GET /api/v1/users/profile
Authorization: Bearer <jwt_token>

# Update user profile
PUT /api/v1/users/profile
Authorization: Bearer <jwt_token>
{
  "username": "john_doe_updated",
  "email": "john.updated@example.com"
}

# Get user list (with pagination and filtering)
GET /api/v1/users?page=1&limit=20&status=active&search=john
Authorization: Bearer <jwt_token>

# Get specific user
GET /api/v1/users/{id}
Authorization: Bearer <jwt_token>

# Delete user
DELETE /api/v1/users/{id}
Authorization: Bearer <jwt_token>
```

## 🛠️ Development

### Project Structure
```
user-center/
├── cmd/usercenter/          # Application entry point
│   ├── main.go             # Main application
│   └── wire.go             # Wire dependency injection
├── internal/               # Private application code
│   ├── config/             # Configuration management
│   ├── model/              # Domain entities (GORM models)
│   ├── dto/                # Data transfer objects
│   ├── service/            # Business logic layer
│   ├── repository/         # Data access layer
│   ├── handler/            # HTTP handlers (controllers)
│   ├── middleware/         # HTTP middleware
│   ├── server/             # Server setup and routing
│   └── database/           # Database connections
├── pkg/                    # Shared packages
│   ├── logger/             # Logging utilities
│   └── jwt/                # JWT utilities
├── configs/                # Configuration files
├── migrations/             # Database migrations
├── docs/                   # Generated documentation
├── Makefile                # Build and development tasks
├── Dockerfile              # Container configuration
└── README.md               # This file
```

### Available Make Commands
```bash
# Development
make run                    # Run in development mode
make build                  # Build binary
make clean                  # Clean build artifacts
make wire                   # Generate Wire dependency injection
make swagger                # Generate Swagger documentation

# Testing
make test                   # Run all tests
make test-coverage          # Run tests with coverage
make test-short             # Run only short tests
make test-race              # Run tests with race detection

# Database
make migrate-up             # Run database migrations
make migrate-down           # Rollback database migrations
make migrate-status         # Check migration status

# Code Quality
make lint                   # Run golangci-lint
make fmt                    # Format code
make vet                    # Run go vet

# Docker
make docker-build           # Build Docker image
make docker-run             # Run Docker container
make docker-clean           # Clean Docker artifacts

# Utilities
make help                   # Show all available commands
```

## 🧪 Testing

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run only unit tests (skip integration)
make test-short

# Run tests with race detection
make test-race

# Run specific test
go test -run TestUserService_CreateUser ./...
```

### Test Coverage
The project aims for 80%+ test coverage. Coverage reports are generated in:
- `coverage.out` - Raw coverage data
- `coverage.html` - HTML coverage report

### Test Structure
- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test database operations and API endpoints
- **Mock Tests**: Use gomock for dependency mocking

## 🚀 Deployment

### Docker Deployment
```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Or use docker-compose
docker-compose up -d
```

### Production Deployment
```bash
# Build for production
make build

# Set environment variables
export USERCENTER_ENV=production
export USERCENTER_DB_HOST=your-db-host
export USERCENTER_DB_PASSWORD=your-db-password

# Run the service
./bin/usercenter
```

### Kubernetes Deployment
```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -l app=usercenter
```

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push branch: `git push origin feature/new-feature`
5. Create Pull Request

### Development Guidelines
- Follow Go coding standards
- Write comprehensive tests
- Update documentation
- Use conventional commit messages
- Ensure all tests pass before submitting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🔗 Related Links

- [中文文档](README.zh-CN.md)
- [Project Homepage](https://github.com/zhwjimmy/user-center)
- [Issues](https://github.com/zhwjimmy/user-center/issues)
- [Discussions](https://github.com/zhwjimmy/user-center/discussions) 