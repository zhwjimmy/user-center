# Architecture Design

## Overview

UserCenter follows clean architecture principles with clear separation between infrastructure and business layers. The system is designed to be maintainable, testable, and scalable through proper abstraction and dependency management.

## 🏗️ Clean Architecture

### Architecture Layers

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

1. **Dependency Inversion**: Business layer defines interfaces, infrastructure implements them
2. **Separation of Concerns**: Clear boundaries between layers
3. **Single Responsibility**: Each component has a single, well-defined purpose
4. **Open/Closed Principle**: Open for extension, closed for modification
5. **Interface Segregation**: Small, focused interfaces

## 📁 Project Structure

```
user-center/
├── cmd/usercenter/          # Application entry point
│   ├── main.go             # Main application
│   └── wire.go             # Wire dependency injection
├── internal/               # Private application code
│   ├── infrastructure/     # External dependencies
│   │   ├── cache/         # Cache implementations
│   │   ├── database/      # Database implementations
│   │   ├── messaging/     # Message queue implementations
│   │   ├── manager.go     # Infrastructure manager
│   │   └── health.go      # Health checks
│   ├── events/            # Event-driven architecture
│   │   ├── types/         # Event definitions
│   │   ├── publisher/     # Event publishing
│   │   └── handlers/      # Event handling
│   ├── service/           # Business logic
│   ├── repository/        # Data access layer
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   ├── config/            # Configuration
│   ├── model/             # Domain models
│   └── dto/               # Data transfer objects
├── pkg/                   # Shared packages
│   └── jwt/               # JWT utilities
├── docs/                  # Documentation
├── configs/               # Configuration files
└── migrations/            # Database migrations
```

## 🔄 Dependency Flow

### Dependency Direction
- **Business Layer** → **Infrastructure Layer** (through interfaces)
- **Infrastructure Layer** never imports business types
- **Business Layer** defines interfaces, infrastructure implements them

### Interface Definitions

#### Infrastructure Interfaces
```go
// Cache interface in infrastructure layer
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    Health(ctx context.Context) error
}

// Event interface for messaging
type Event interface {
    GetTopic() string
    GetEventType() string
    GetUserID() string
    GetRequestID() string
    GetTimestamp() string
}
```

#### Business Layer Implementation
```go
// Business events implement infrastructure interfaces
func (e *UserRegisteredEvent) GetTopic() string {
    return "user_registered"
}

func (e *UserRegisteredEvent) GetEventType() string {
    return string(e.Type)
}
```

## 🚀 Event-Driven Architecture

### Event Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Service   │───▶│  Publisher  │───▶│   Kafka     │
└─────────────┘    └─────────────┘    └─────────────┘
                                              │
                                              ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Handlers   │◀───│  Consumer   │◀───│   Topic     │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Event Types

- **User Registration**: `user.registered`
- **User Login**: `user.logged_in`
- **Password Change**: `user.password_changed`
- **Status Change**: `user.status_changed`
- **User Deletion**: `user.deleted`
- **User Update**: `user.updated`

### Event Structure

```go
type BaseEvent struct {
    ID        string                 `json:"id"`
    Type      EventType              `json:"type"`
    Source    string                 `json:"source"`
    Timestamp time.Time              `json:"timestamp"`
    Version   string                 `json:"version"`
    RequestID string                 `json:"request_id,omitempty"`
    UserID    string                 `json:"user_id,omitempty"`
    Data      map[string]interface{} `json:"data"`
}
```

## 🏢 Infrastructure Management

### Infrastructure Manager

The `InfrastructureManager` centralizes all external dependency management:

```go
type Manager struct {
    postgreSQL *postgreSQL
    mongoDB    *mongoDB
    redis      *redisImpl
    kafka      *kafkaService
    logger     *zap.Logger
}

func (m *Manager) Start() error {
    // Start all infrastructure components
}

func (m *Manager) Stop() error {
    // Stop all infrastructure components
}
```

### Health Checks

Unified health checking for all infrastructure components:

```go
func (m *Manager) Health(ctx context.Context) error {
    // Check all infrastructure components
}

func (m *Manager) Ready(ctx context.Context) error {
    // Check if all components are ready
}
```

## 🔧 Dependency Injection

### Wire Configuration

Using Google Wire for compile-time dependency injection:

```go
func InitializeApp() (*server.Server, error) {
    wire.Build(
        // Configuration
        config.Load,
        
        // Infrastructure
        provideInfrastructureManager,
        provideGormDB,
        provideCache,
        provideKafkaService,
        
        // Business Layer
        provideEventPublisher,
        service.NewEventService,
        service.NewAuthService,
        
        // Handlers
        handler.NewUserHandler,
        
        // Server
        provideServer,
    )
    return &server.Server{}, nil
}
```

## 🛡️ Security Architecture

### Authentication & Authorization

- **JWT-based Authentication**: Stateless token-based authentication
- **Role-based Access Control**: Granular permission system
- **Rate Limiting**: Multiple rate limiting strategies
- **Input Validation**: Comprehensive request validation

### Security Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   CORS      │  │ Rate Limit  │  │   Auth      │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Business Layer                           │
└─────────────────────────────────────────────────────────────┘
```

## 📊 Observability

### Logging

- **Structured Logging**: Using Zap for high-performance logging
- **Request Tracing**: Request ID propagation across services
- **Contextual Information**: Rich context in log entries

### Metrics

- **Prometheus Integration**: Custom business metrics
- **Health Checks**: Comprehensive health monitoring
- **Performance Metrics**: Response times, throughput, error rates

### Tracing

- **Distributed Tracing**: OpenTelemetry integration
- **Request Flow**: End-to-end request tracking
- **Performance Analysis**: Bottleneck identification

## 🔄 Data Flow

### Request Flow

```
Client Request
    │
    ▼
┌─────────────┐
│ Middleware  │ ← CORS, Auth, Rate Limit
└─────────────┘
    │
    ▼
┌─────────────┐
│   Handler   │ ← Request validation
└─────────────┘
    │
    ▼
┌─────────────┐
│   Service   │ ← Business logic
└─────────────┘
    │
    ▼
┌─────────────┐
│ Repository  │ ← Data access
└─────────────┘
    │
    ▼
┌─────────────┐
│Infrastructure│ ← Database, Cache, MQ
└─────────────┘
```

### Event Flow

```
Business Logic
    │
    ▼
┌─────────────┐
│ Event Service│ ← Event creation
└─────────────┘
    │
    ▼
┌─────────────┐
│  Publisher  │ ← Event publishing
└─────────────┘
    │
    ▼
┌─────────────┐
│   Kafka     │ ← Message queue
└─────────────┘
    │
    ▼
┌─────────────┐
│  Consumer   │ ← Event consumption
└─────────────┘
    │
    ▼
┌─────────────┐
│   Handler   │ ← Event processing
└─────────────┘
```

## 🧪 Testing Architecture

### Testing Strategy

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete user workflows
- **Mock Testing**: Use interfaces for easy mocking

### Test Structure

```
tests/
├── unit/              # Unit tests
├── integration/       # Integration tests
├── e2e/              # End-to-end tests
└── mocks/            # Mock implementations
```

## 🚀 Scalability Considerations

### Horizontal Scaling

- **Stateless Design**: No server-side session storage
- **Database Sharding**: Support for database partitioning
- **Load Balancing**: Multiple service instances
- **Caching Strategy**: Multi-level caching

### Performance Optimization

- **Connection Pooling**: Database and cache connection pools
- **Batch Processing**: Efficient bulk operations
- **Async Processing**: Non-blocking event processing
- **Compression**: Data compression for network efficiency

## 🔄 Migration Strategy

### Database Migrations

- **Version Control**: Goose for database versioning
- **Rollback Support**: Safe migration rollbacks
- **Zero Downtime**: Online schema changes
- **Data Integrity**: Referential integrity maintenance

### Application Updates

- **Blue-Green Deployment**: Zero-downtime deployments
- **Feature Flags**: Gradual feature rollouts
- **Backward Compatibility**: API versioning
- **Rollback Capability**: Quick service rollbacks

## 📈 Future Considerations

### Planned Enhancements

- **Microservices**: Service decomposition
- **Event Sourcing**: Complete event history
- **CQRS**: Command Query Responsibility Segregation
- **GraphQL**: Flexible API querying
- **gRPC**: High-performance RPC

### Technology Evolution

- **Cloud Native**: Kubernetes-native deployment
- **Serverless**: Function-as-a-Service integration
- **Edge Computing**: Distributed processing
- **AI/ML Integration**: Intelligent user management

---

## 🔗 Related Documentation

- [Getting Started Guide](getting-started.md)
- [Development Guide](development.md)
- [Kafka Integration](kafka-integration.md)
- [API Reference](api-reference.md)
- [Deployment Guide](deployment.md) 