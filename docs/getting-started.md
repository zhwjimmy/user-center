# Getting Started Guide

## Overview

This guide will help you get UserCenter up and running in your development environment. UserCenter is a user management service built with Go, featuring clean architecture and event-driven patterns.

## Prerequisites

### Required Software
- **Go 1.23.1+**: Latest stable version of Go
- **Docker & Docker Compose**: For running dependencies
- **Git**: For version control

### Required Services
- **PostgreSQL 13+**: Primary database
- **Redis 6.0+**: Caching layer
- **Apache Kafka 2.8+**: Message queue for events

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd user-center
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start Infrastructure Services

UserCenter uses Docker Compose to manage external dependencies:

```bash
docker-compose up -d
```

This starts:
- PostgreSQL database
- Redis cache
- Apache Kafka with Zookeeper
- Prometheus for metrics
- Jaeger for tracing

### 4. Configure Environment

Create a `.env` file for local development:

```bash
# Database configuration
USERCENTER_DATABASE_POSTGRES_HOST=localhost
USERCENTER_DATABASE_POSTGRES_PORT=5432
USERCENTER_DATABASE_POSTGRES_USER=postgres
USERCENTER_DATABASE_POSTGRES_PASSWORD=password
USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter
USERCENTER_DATABASE_POSTGRES_SSLMODE=disable

# Redis configuration
USERCENTER_CACHE_REDIS_ADDR=localhost:6379
USERCENTER_CACHE_REDIS_PASSWORD=
USERCENTER_CACHE_REDIS_DB=0

# Kafka configuration
USERCENTER_KAFKA_BROKERS=localhost:9092
USERCENTER_KAFKA_GROUP_ID=usercenter
```

### 5. Generate Application Code

UserCenter uses Google Wire for dependency injection:

```bash
go generate ./cmd/usercenter
```

### 6. Run Database Migrations

Initialize the database schema:

```bash
make migrate-up
```

### 7. Start the Application

Run in development mode with hot reload:

```bash
make run-dev
```

### 8. Verify Installation

Check that the service is running:

```bash
# Health check
curl http://localhost:8080/health

# API documentation
open http://localhost:8080/swagger/index.html
```

## Development Workflow

### Local Development

UserCenter follows a hybrid development approach:
- **Application Service**: Runs locally for fast development
- **Dependency Services**: Managed via Docker Compose
- **Hot Reload**: Code changes trigger automatic rebuild

### Available Commands

```bash
# Development
make run-dev              # Run with hot reload
make build                # Build binary
make clean                # Clean build artifacts

# Testing
make test                 # Run all tests
make test-coverage        # Run with coverage

# Database
make migrate-up           # Run migrations
make migrate-down         # Rollback migrations

# Documentation
make swagger              # Generate API docs
```

## Service Access

Once running, you can access:

- **API Service**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Metrics**: http://localhost:8080/metrics
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090

## Architecture Overview

UserCenter follows clean architecture principles:

### Layers
- **Business Layer**: Core business logic and domain models
- **Infrastructure Layer**: External dependencies (database, cache, messaging)
- **Interface Layer**: HTTP handlers and middleware

### Key Components
- **Infrastructure Manager**: Centralized dependency management
- **Event-Driven Architecture**: Asynchronous event processing with Kafka
- **Dependency Injection**: Google Wire for compile-time DI

## Next Steps

### For Developers
1. Read the [Architecture Guide](architecture.md) to understand the system design
2. Review the [Development Guide](development.md) for detailed development workflow
3. Check the [API Reference](api-reference.md) for available endpoints

### For Operations
1. Review the [Deployment Guide](deployment.md) for production deployment
2. Check the [Configuration Guide](configuration.md) for environment setup
3. Read the [Monitoring Guide](monitoring.md) for observability

### For Integration
1. Review the [Kafka Integration](kafka-integration.md) for event processing
2. Check the [API Reference](api-reference.md) for integration endpoints

## Troubleshooting

### Common Issues

**Service won't start**
- Check if all dependencies are running: `docker-compose ps`
- Verify environment variables are set correctly
- Check logs: `docker-compose logs`

**Database connection failed**
- Ensure PostgreSQL is running: `docker-compose ps postgres`
- Verify database credentials in `.env` file
- Check if migrations ran successfully

**Kafka connection failed**
- Ensure Kafka is running: `docker-compose ps kafka`
- Check Kafka broker configuration
- Verify topic creation

### Getting Help

- **Documentation**: Check the [Troubleshooting Guide](troubleshooting.md)
- **Issues**: Create an issue on GitHub
- **Discussions**: Join GitHub Discussions for community help

## Environment Variables

### Required Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `USERCENTER_DATABASE_POSTGRES_HOST` | PostgreSQL host | localhost |
| `USERCENTER_DATABASE_POSTGRES_PORT` | PostgreSQL port | 5432 |
| `USERCENTER_DATABASE_POSTGRES_USER` | Database user | postgres |
| `USERCENTER_DATABASE_POSTGRES_PASSWORD` | Database password | password |
| `USERCENTER_DATABASE_POSTGRES_DBNAME` | Database name | usercenter |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `USERCENTER_SERVER_PORT` | HTTP server port | 8080 |
| `USERCENTER_LOG_LEVEL` | Logging level | info |
| `USERCENTER_ENV` | Environment | development |

## Development Tips

### Hot Reload
The development server automatically rebuilds and restarts when you make changes to Go files.

### Environment Management
Use `.env` files for local development. The application automatically loads environment variables from `.env` files.

### Database Migrations
Always run migrations after pulling new changes: `make migrate-up`

### Testing
Run tests frequently: `make test`. The project aims for 80%+ test coverage.

---

## 🔗 Related Documentation

- [Architecture Guide](architecture.md)
- [Development Guide](development.md)
- [API Reference](api-reference.md)
- [Configuration Guide](configuration.md)
- [Troubleshooting Guide](troubleshooting.md) 