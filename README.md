# UserCenter - User Management Service

[![Go Version](https://img.shields.io/badge/Go-1.23.1-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI](https://github.com/zhwjimmy/user-center/workflows/CI/badge.svg)](https://github.com/zhwjimmy/user-center/actions/workflows/ci.yml)
[![Release](https://github.com/zhwjimmy/user-center/workflows/Release/badge.svg)](https://github.com/zhwjimmy/user-center/actions/workflows/release.yml)
[![Security Scan](https://github.com/zhwjimmy/user-center/workflows/Security%20Scan/badge.svg)](https://github.com/zhwjimmy/user-center/actions/workflows/security.yml)

**English** | [中文](README.zh-CN.md)

## 🎯 Overview

UserCenter is a production-ready user management service built with Go, featuring clean architecture design and event-driven patterns. It provides comprehensive user management capabilities including registration, authentication, querying, and listing with high concurrency, availability, and scalability.

### ✨ Key Features

- 🔐 **JWT-based Authentication** with secure password hashing
- 👥 **User Management** with UUID-based identification and soft delete
- 🚀 **Event-Driven Architecture** using Apache Kafka for async processing
- 🏗️ **Clean Architecture** with clear separation of concerns
- 🛡️ **Security Features** including rate limiting, CORS, and input validation
- 📊 **Observability** with health checks, metrics, and structured logging
- 🌍 **Internationalization** support (Chinese/English)
- 🔄 **Graceful Shutdown** and dependency management

## 🚀 Quick Start

### Prerequisites
- Go 1.23.1+
- PostgreSQL 13+
- Redis 6.0+
- Apache Kafka 2.8+

### 1. Clone and Setup
```bash
git clone <repository-url>
cd user-center
go mod download
```

### 2. Start Dependencies
```bash
docker-compose up -d
```

### 3. Configure Environment
```bash
# Create .env file
cat > .env << EOF
USERCENTER_DATABASE_POSTGRES_HOST=localhost
USERCENTER_DATABASE_POSTGRES_PORT=5432
USERCENTER_DATABASE_POSTGRES_USER=postgres
USERCENTER_DATABASE_POSTGRES_PASSWORD=password
USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter
USERCENTER_DATABASE_POSTGRES_SSLMODE=disable
EOF
```

### 4. Run the Service
```bash
# Generate dependency injection code
go generate ./cmd/usercenter

# Run in development mode
make run-dev
```

### 5. Verify Installation
```bash
# Health check
curl http://localhost:8080/health

# API documentation
open http://localhost:8080/swagger/index.html
```

## 📚 Documentation

📖 **Complete Documentation**: [docs/README.md](docs/README.md)

### Quick Links
- 🏗️ [Architecture Design](docs/architecture.md)
- 🚀 [Getting Started Guide](docs/getting-started.md)
- 📖 [API Reference](docs/api-reference.md)
- 🛠️ [Development Guide](docs/development.md)
- 🚀 [Deployment Guide](docs/deployment.md)
- 🔧 [Configuration Guide](docs/configuration.md)
- 🐛 [Troubleshooting](docs/troubleshooting.md)
- 📊 [Kafka Integration](docs/kafka-integration.md)

## 🏗️ Architecture

UserCenter follows clean architecture principles with clear separation between infrastructure and business layers:

```
┌─────────────────────────────────────────────────────────────┐
│                    Business Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Services  │  │  Handlers   │  │   Events    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Database   │  │    Cache    │  │  Messaging  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Principles
- **Dependency Inversion**: Business layer defines interfaces, infrastructure implements them
- **Event-Driven**: Asynchronous event processing with Kafka
- **Separation of Concerns**: Clear boundaries between layers
- **Testability**: Dependency injection enables easy testing

## 🛠️ Development

### Project Structure
```
user-center/
├── cmd/usercenter/          # Application entry point
├── internal/               # Private application code
│   ├── infrastructure/     # External dependencies (DB, Cache, MQ)
│   ├── events/            # Event-driven architecture
│   ├── service/           # Business logic
│   ├── handler/           # HTTP handlers
│   └── middleware/        # HTTP middleware
├── docs/                  # Documentation
├── configs/               # Configuration files
└── migrations/            # Database migrations
```

### Available Commands
```bash
# Development
make run-dev              # Run with hot reload
make build                # Build binary
make wire                 # Generate dependency injection

# Testing
make test                 # Run all tests
make test-coverage        # Run with coverage

# Database
make migrate-up           # Run migrations
make migrate-down         # Rollback migrations

# Documentation
make swagger              # Generate API docs
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/contributing.md) for details.

### Quick Contribution Steps
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -am 'Add amazing feature'`
4. Push branch: `git push origin feature/amazing-feature`
5. Create Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🔗 Links

- [中文文档](README.zh-CN.md)
- [Project Homepage](https://github.com/zhwjimmy/user-center)
- [Issues](https://github.com/zhwjimmy/user-center/issues)
- [Discussions](https://github.com/zhwjimmy/user-center/discussions) 