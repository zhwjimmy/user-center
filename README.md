# UserCenter - User Management Service

[![Go Version](https://img.shields.io/badge/Go-1.23.1-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI](https://github.com/zhwjimmy/user-center/workflows/CI/badge.svg)](https://github.com/zhwjimmy/user-center/actions/workflows/ci.yml)
[![Release](https://github.com/zhwjimmy/user-center/workflows/Release/badge.svg)](https://github.com/zhwjimmy/user-center/actions/workflows/release.yml)
[![Security Scan](https://github.com/zhwjimmy/user-center/workflows/Security%20Scan/badge.svg)](https://github.com/zhwjimmy/user-center/actions/workflows/security.yml)

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
- [CI/CD](#cicd)
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
- User profile management with UUID-based identification
- Account status management (active, inactive, suspended)
- Soft delete support
- Bulk user operations
- UUID-based user identification for enhanced security

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

### Event-Driven Architecture
- **Asynchronous Event Processing**: Kafka-based event-driven architecture for user lifecycle events
- **Event Types**: User registration, login, password change, status change, deletion, and update events
- **Reliable Delivery**: Idempotent producers with retry mechanisms and message acknowledgment
- **Scalable Processing**: Consumer groups with load balancing and horizontal scaling
- **Observability**: Comprehensive logging, metrics, and health checks for Kafka operations
- **Graceful Degradation**: Event publishing failures don't affect main business flows

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

## 📚 Kafka Event Processing

### Event-Driven Architecture

The application implements a robust event-driven architecture using Apache Kafka for asynchronous processing of user lifecycle events.

#### Supported Event Types
- **User Registration**: `user.registered` - Triggered when a new user registers
- **User Login**: `user.logged_in` - Triggered when a user successfully logs in
- **Password Change**: `user.password_changed` - Triggered when a user changes their password
- **Status Change**: `user.status_changed` - Triggered when user status is modified
- **User Deletion**: `user.deleted` - Triggered when a user account is deleted
- **User Update**: `user.updated` - Triggered when user profile is updated

#### Event Processing Features
- **Reliable Delivery**: Idempotent producers with retry mechanisms
- **Message Acknowledgment**: Consumer group with automatic offset management
- **Scalable Processing**: Horizontal scaling with consumer groups
- **Graceful Degradation**: Event publishing failures don't affect main business flows
- **Comprehensive Logging**: Structured logging with request ID tracking
- **Health Monitoring**: Kafka connectivity and consumer group health checks

### Kafka Configuration

```yaml
kafka:
  brokers: ["localhost:9092"]
  topics:
    user_events: "user.events"
    user_notifications: "user.notifications"
    user_analytics: "user.analytics"
  group_id: "usercenter"
```

### Testing Kafka Integration

```bash
# Run Kafka integration tests
make test-kafka

# Manual Kafka test script
./scripts/test-kafka.sh

# Test user registration (triggers Kafka event)
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# Monitor Kafka messages
kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic user.events \
  --from-beginning
```

### Kafka Documentation

- **Integration Guide**: `docs/kafka-integration.md` - Detailed setup and usage guide
- **Implementation Summary**: `KAFKA_IMPLEMENTATION_SUMMARY.md` - High-level architecture overview
- **Test Script**: `scripts/test-kafka.sh` - Automated Kafka functionality testing

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
make test-coverage-xml      # Run tests with XML coverage report
make test-short             # Run only short tests
make test-race              # Run tests with race detection
make mockgen                # Generate mocks for testing

# Kafka Testing
make test-kafka             # Run Kafka integration tests
./scripts/test-kafka.sh     # Manual Kafka test script

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

# Run tests with XML coverage report (for CI)
make test-coverage-xml

# Run only unit tests (skip integration)
make test-short

# Run tests with race detection
make test-race

# Generate mocks for testing
make mockgen

# Run specific test
go test -run TestUserService_CreateUser ./...
```

### Test Coverage
The project aims for 80%+ test coverage. Coverage reports are generated in:
- `coverage.out` - Raw coverage data
- `coverage.html` - HTML coverage report
- `coverage.xml` - XML coverage report (for CI integration)

### Test Structure
- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test database operations and API endpoints
- **Mock Tests**: Use gomock for dependency mocking
- **Mock Generation**: Automatically generate mocks using `mockgen`

## 🔄 CI/CD

This project uses GitHub Actions for continuous integration and deployment. The CI/CD pipeline includes code quality checks, testing, building, and automated deployment.

### Workflows

#### 1. CI Workflow (`ci.yml`)
- **Triggers**: Push to `main`/`develop` branches, Pull Requests
- **Features**:
  - Unit and integration tests with coverage
  - Mock generation and dependency injection code generation
  - XML coverage reporting for CI integration
  - Optimized for fast execution (parallel testing, caching)

#### 2. Release Workflow (`release.yml`)
- **Triggers**: Version tag pushes (e.g., `v1.0.0`)
- **Features**:
  - Builds and publishes Docker images to GitHub Container Registry
  - Creates GitHub Releases with assets
  - Multi-architecture support (linux/amd64, linux/arm64)

#### 3. Deploy Workflow (`deploy.yml`)
- **Triggers**: After successful CI completion on `main` branch
- **Features**:
  - Automated deployment to staging environment
  - Automated deployment to production environment
  - Deployment notifications

#### 4. Security Workflow (`security.yml`)
- **Triggers**: Weekly schedule, manual dispatch, dependency changes
- **Features**:
  - Code security scanning (gosec)
  - Dependency vulnerability checks (govulncheck)
  - Docker image security scanning (Trivy)
  - File system security scanning

### Setup

1. **Enable GitHub Actions** in your repository settings
2. **Configure Secrets** for database connections and deployment
3. **Set up Environments** for staging and production
4. **Configure Dependabot** for automated dependency updates
5. **Ensure `go.sum` is committed** (not ignored) for reproducible builds

### Usage

```bash
# Create a new release
git tag v1.0.0
git push origin v1.0.0

# Check workflow status
# Visit: https://github.com/username/user-center/actions

# View security scan results
# Visit: https://github.com/username/user-center/security
```

For detailed configuration and troubleshooting, see [GitHub Actions Documentation](docs/github-actions.md).

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