# Development Guide

## Overview

This guide covers the development workflow, best practices, and architectural patterns for contributing to UserCenter. The project follows clean architecture principles and event-driven design patterns.

## Development Philosophy

### Clean Architecture
- **Separation of Concerns**: Clear boundaries between business logic and infrastructure
- **Dependency Inversion**: Business layer defines interfaces, infrastructure implements them
- **Testability**: Easy to test through dependency injection and interface abstraction

### Event-Driven Design
- **Asynchronous Processing**: Non-blocking event processing for better performance
- **Loose Coupling**: Services communicate through events, not direct calls
- **Scalability**: Horizontal scaling through event-driven patterns

## Development Workflow

### 1. Environment Setup

Ensure your development environment is properly configured:
- Go 1.23.1+ installed
- Docker and Docker Compose running
- Dependencies started: `docker-compose up -d`
- Environment variables configured in `.env`

### 2. Code Generation

UserCenter uses code generation for dependency injection:

```bash
# Generate Wire dependency injection code
go generate ./cmd/usercenter

# Generate API documentation
make swagger

# Generate mocks for testing
make mockgen
```

### 3. Development Cycle

1. **Write Code**: Implement features following architectural patterns
2. **Run Tests**: Ensure all tests pass: `make test`
3. **Check Quality**: Run linting: `make lint`
4. **Test Integration**: Verify with integration tests
5. **Document Changes**: Update relevant documentation

## Project Structure

### Directory Organization

```
internal/
├── infrastructure/     # External dependencies
│   ├── cache/         # Cache implementations
│   ├── database/      # Database implementations
│   ├── messaging/     # Message queue implementations
│   └── manager.go     # Infrastructure manager
├── events/            # Event-driven architecture
│   ├── types/         # Event definitions
│   ├── publisher/     # Event publishing
│   └── handlers/      # Event handling
├── service/           # Business logic
├── repository/        # Data access layer
├── handler/           # HTTP handlers
├── middleware/        # HTTP middleware
├── config/            # Configuration
├── model/             # Domain models
└── dto/               # Data transfer objects
```

### Key Principles

- **Infrastructure Layer**: Contains all external dependency implementations
- **Business Layer**: Contains core business logic and domain models
- **Interface Layer**: HTTP handlers and middleware
- **Events Layer**: Event definitions, publishing, and handling

## Architectural Patterns

### 1. Dependency Injection

UserCenter uses Google Wire for compile-time dependency injection:

- **Provider Functions**: Define how dependencies are created
- **Interface-Based Design**: Depend on interfaces, not concrete implementations
- **Testability**: Easy to mock dependencies for testing

### 2. Repository Pattern

Data access is abstracted through repository interfaces:

- **Interface Definition**: Business layer defines repository interfaces
- **Implementation**: Infrastructure layer provides concrete implementations
- **Testing**: Easy to mock repositories for unit testing

### 3. Event-Driven Architecture

Asynchronous event processing for better scalability:

- **Event Publishing**: Business services publish events through publishers
- **Event Handling**: Separate handlers process events asynchronously
- **Infrastructure Decoupling**: Events use generic interfaces for infrastructure independence

## Development Best Practices

### 1. Code Organization

- **Single Responsibility**: Each function and struct has a single, well-defined purpose
- **Interface Segregation**: Keep interfaces small and focused
- **Dependency Direction**: Dependencies point inward (infrastructure → business)

### 2. Error Handling

- **Structured Errors**: Use wrapped errors with context
- **Graceful Degradation**: Don't let infrastructure failures affect business logic
- **Logging**: Log errors with appropriate levels and context

### 3. Testing Strategy

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **Mock Testing**: Use interfaces for easy mocking
- **Test Coverage**: Aim for 80%+ test coverage

### 4. Configuration Management

- **Environment Variables**: Use environment variables for configuration
- **Validation**: Validate configuration at startup
- **Defaults**: Provide sensible defaults for all configuration options

## Adding New Features

### 1. Feature Planning

1. **Define Requirements**: Clear understanding of what needs to be built
2. **Architecture Review**: Ensure the feature fits the architectural patterns
3. **Interface Design**: Define interfaces before implementation

### 2. Implementation Steps

1. **Define Interfaces**: Business layer defines required interfaces
2. **Implement Infrastructure**: Infrastructure layer implements interfaces
3. **Implement Business Logic**: Business layer implements core logic
4. **Add Tests**: Comprehensive test coverage
5. **Update Documentation**: Update relevant documentation

### 3. Event Integration

When adding new features that require event processing:

1. **Define Event Types**: Add new event types in `internal/events/types/`
2. **Implement Event Interface**: Ensure events implement the `Event` interface
3. **Add Event Handlers**: Create handlers in `internal/events/handlers/`
4. **Update Publishers**: Add publishing methods to event publishers
5. **Update Services**: Integrate event publishing in business services

## Code Quality

### 1. Linting and Formatting

```bash
# Format code
make fmt

# Run linter
make lint

# Run vet
make vet
```

### 2. Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific tests
go test ./internal/service -v
```

### 3. Performance

- **Benchmarking**: Use Go's benchmarking tools for performance-critical code
- **Profiling**: Use pprof for performance analysis
- **Monitoring**: Monitor performance metrics in production

## Database Development

### 1. Migrations

- **Version Control**: All schema changes go through migrations
- **Rollback Support**: Migrations should be reversible
- **Testing**: Test migrations in development before production

### 2. Schema Design

- **Normalization**: Follow database normalization principles
- **Indexing**: Proper indexing for query performance
- **Constraints**: Use database constraints for data integrity

### 3. Query Optimization

- **Connection Pooling**: Efficient database connection management
- **Query Analysis**: Monitor and optimize slow queries
- **Caching**: Use Redis for frequently accessed data

## Event-Driven Development

### 1. Event Design

- **Event Naming**: Use descriptive, past-tense event names
- **Event Structure**: Include all necessary context in events
- **Event Versioning**: Plan for event schema evolution

### 2. Event Publishing

- **Asynchronous**: Don't block business logic for event publishing
- **Error Handling**: Handle publishing failures gracefully
- **Monitoring**: Monitor event publishing success rates

### 3. Event Handling

- **Idempotency**: Ensure handlers are idempotent
- **Error Recovery**: Implement proper error recovery mechanisms
- **Monitoring**: Monitor event processing performance

## Security Considerations

### 1. Input Validation

- **Request Validation**: Validate all incoming requests
- **SQL Injection**: Use parameterized queries
- **XSS Prevention**: Sanitize user inputs

### 2. Authentication & Authorization

- **JWT Security**: Secure JWT token handling
- **Role-Based Access**: Implement proper RBAC
- **Rate Limiting**: Protect against abuse

### 3. Data Protection

- **Encryption**: Encrypt sensitive data at rest and in transit
- **Access Control**: Implement proper access controls
- **Audit Logging**: Log security-relevant events

## Performance Optimization

### 1. Caching Strategy

- **Multi-Level Caching**: Use Redis for distributed caching
- **Cache Invalidation**: Proper cache invalidation strategies
- **Cache Warming**: Pre-populate frequently accessed data

### 2. Database Optimization

- **Query Optimization**: Optimize database queries
- **Connection Pooling**: Efficient connection management
- **Read Replicas**: Use read replicas for read-heavy workloads

### 3. Application Optimization

- **Goroutine Management**: Efficient goroutine usage
- **Memory Management**: Monitor and optimize memory usage
- **Concurrency**: Use appropriate concurrency patterns

## Monitoring and Observability

### 1. Logging

- **Structured Logging**: Use structured logging with Zap
- **Log Levels**: Appropriate log levels for different types of information
- **Context**: Include relevant context in log entries

### 2. Metrics

- **Business Metrics**: Track business-relevant metrics
- **Performance Metrics**: Monitor performance indicators
- **Error Rates**: Track error rates and types

### 3. Tracing

- **Request Tracing**: End-to-end request tracing
- **Performance Analysis**: Identify performance bottlenecks
- **Dependency Mapping**: Understand service dependencies

## Deployment Considerations

### 1. Environment Management

- **Configuration**: Environment-specific configuration
- **Secrets Management**: Secure handling of secrets
- **Feature Flags**: Gradual feature rollouts

### 2. Health Checks

- **Readiness Probes**: Ensure service is ready to handle requests
- **Liveness Probes**: Ensure service is healthy
- **Dependency Checks**: Check external dependency health

### 3. Graceful Shutdown

- **Signal Handling**: Proper signal handling for shutdown
- **Resource Cleanup**: Clean up resources on shutdown
- **In-Flight Requests**: Handle in-flight requests gracefully

## Contributing Guidelines

### 1. Code Review Process

- **Pull Requests**: All changes go through pull requests
- **Code Review**: Thorough code review by maintainers
- **Testing**: Ensure all tests pass before merging

### 2. Documentation

- **Code Comments**: Clear and helpful code comments
- **API Documentation**: Keep API documentation up-to-date
- **Architecture Documentation**: Update architecture documentation for significant changes

### 3. Commit Messages

- **Conventional Commits**: Use conventional commit message format
- **Descriptive**: Clear and descriptive commit messages
- **Atomic**: Each commit should be atomic and focused

---

## 🔗 Related Documentation

- [Getting Started Guide](getting-started.md)
- [Architecture Guide](architecture.md)
- [API Reference](api-reference.md)
- [Testing Guide](testing.md)
- [Contributing Guide](contributing.md) 