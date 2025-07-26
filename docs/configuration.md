# Configuration Guide

## Overview

UserCenter supports multiple configuration methods with a hierarchical approach. Configuration can be provided through environment variables, configuration files, or default values.

## Configuration Hierarchy

Configuration is loaded in the following order (highest priority first):

1. **Environment Variables** - Runtime configuration
2. **Configuration Files** - YAML configuration files
3. **Default Values** - Hardcoded defaults in the application

## Configuration Methods

### 1. Environment Variables

Environment variables provide the most flexible configuration method, especially for containerized deployments.

**Naming Convention**: All environment variables use the `USERCENTER_` prefix followed by the configuration path.

**Examples**:
```bash
USERCENTER_SERVER_PORT=8080
USERCENTER_DATABASE_POSTGRES_HOST=localhost
USERCENTER_KAFKA_BROKERS=localhost:9092
```

### 2. Configuration Files

YAML configuration files provide a structured way to manage configuration.

**Default Location**: `configs/config.yaml`

**File Structure**:
```yaml
server:
  port: 8080
  mode: debug

database:
  postgres:
    host: localhost
    port: 5432
    # ... other database settings

kafka:
  brokers: ["localhost:9092"]
  # ... other kafka settings
```

### 3. Default Values

The application provides sensible defaults for all configuration options, ensuring it can run with minimal configuration.

## Configuration Categories

### Server Configuration

Controls the HTTP server behavior and settings.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_SERVER_PORT` | HTTP server port | 8080 | No |
| `USERCENTER_SERVER_HOST` | HTTP server host | 0.0.0.0 | No |
| `USERCENTER_SERVER_MODE` | Server mode (debug/release) | debug | No |
| `USERCENTER_SERVER_TIMEOUT` | Request timeout | 30s | No |

### Database Configuration

PostgreSQL database connection and settings.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_DATABASE_POSTGRES_HOST` | PostgreSQL host | localhost | Yes |
| `USERCENTER_DATABASE_POSTGRES_PORT` | PostgreSQL port | 5432 | No |
| `USERCENTER_DATABASE_POSTGRES_USER` | Database user | postgres | Yes |
| `USERCENTER_DATABASE_POSTGRES_PASSWORD` | Database password | - | Yes |
| `USERCENTER_DATABASE_POSTGRES_DBNAME` | Database name | usercenter | Yes |
| `USERCENTER_DATABASE_POSTGRES_SSLMODE` | SSL mode | disable | No |
| `USERCENTER_DATABASE_POSTGRES_MAX_OPEN_CONNS` | Max open connections | 25 | No |
| `USERCENTER_DATABASE_POSTGRES_MAX_IDLE_CONNS` | Max idle connections | 5 | No |

### Cache Configuration

Redis cache connection and settings.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_CACHE_REDIS_ADDR` | Redis address | localhost:6379 | No |
| `USERCENTER_CACHE_REDIS_PASSWORD` | Redis password | - | No |
| `USERCENTER_CACHE_REDIS_DB` | Redis database | 0 | No |
| `USERCENTER_CACHE_REDIS_POOL_SIZE` | Connection pool size | 10 | No |
| `USERCENTER_CACHE_REDIS_MIN_IDLE_CONNS` | Min idle connections | 5 | No |

### Kafka Configuration

Apache Kafka connection and messaging settings.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_KAFKA_BROKERS` | Kafka broker addresses | localhost:9092 | No |
| `USERCENTER_KAFKA_GROUP_ID` | Consumer group ID | usercenter | No |
| `USERCENTER_KAFKA_TOPICS_USER_EVENTS` | User events topic | user.events | No |
| `USERCENTER_KAFKA_TOPICS_USER_NOTIFICATIONS` | Notifications topic | user.notifications | No |
| `USERCENTER_KAFKA_TOPICS_USER_ANALYTICS` | Analytics topic | user.analytics | No |

### JWT Configuration

JWT token settings for authentication.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_JWT_SECRET` | JWT signing secret | - | Yes |
| `USERCENTER_JWT_ISSUER` | JWT issuer | usercenter | No |
| `USERCENTER_JWT_EXPIRY` | Token expiry time | 24h | No |
| `USERCENTER_JWT_REFRESH_EXPIRY` | Refresh token expiry | 168h | No |

### Logging Configuration

Logging behavior and output settings.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_LOG_LEVEL` | Log level (debug/info/warn/error) | info | No |
| `USERCENTER_LOG_FORMAT` | Log format (json/console) | json | No |
| `USERCENTER_LOG_OUTPUT_PATH` | Log output path | - | No |
| `USERCENTER_LOG_MAX_SIZE` | Max log file size | 100MB | No |
| `USERCENTER_LOG_MAX_AGE` | Max log file age | 30 days | No |

### Monitoring Configuration

Observability and monitoring settings.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_MONITORING_PROMETHEUS_ENABLED` | Enable Prometheus metrics | true | No |
| `USERCENTER_MONITORING_PROMETHEUS_PORT` | Prometheus metrics port | 9090 | No |
| `USERCENTER_MONITORING_TRACING_ENABLED` | Enable distributed tracing | true | No |
| `USERCENTER_MONITORING_TRACING_ENDPOINT` | Tracing endpoint | - | No |

### Security Configuration

Security-related settings and policies.

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USERCENTER_SECURITY_CORS_ALLOWED_ORIGINS` | CORS allowed origins | * | No |
| `USERCENTER_SECURITY_RATE_LIMIT_ENABLED` | Enable rate limiting | true | No |
| `USERCENTER_SECURITY_RATE_LIMIT_REQUESTS` | Rate limit requests per minute | 100 | No |
| `USERCENTER_SECURITY_PASSWORD_MIN_LENGTH` | Minimum password length | 8 | No |

## Environment-Specific Configuration

### Development Environment

For local development, use the following configuration:

```bash
# Development environment variables
USERCENTER_ENV=development
USERCENTER_SERVER_MODE=debug
USERCENTER_LOG_LEVEL=debug
USERCENTER_LOG_FORMAT=console
USERCENTER_DATABASE_POSTGRES_HOST=localhost
USERCENTER_DATABASE_POSTGRES_PASSWORD=password
USERCENTER_JWT_SECRET=dev-secret-key
```

### Production Environment

For production deployments, use secure configuration:

```bash
# Production environment variables
USERCENTER_ENV=production
USERCENTER_SERVER_MODE=release
USERCENTER_LOG_LEVEL=info
USERCENTER_LOG_FORMAT=json
USERCENTER_DATABASE_POSTGRES_HOST=your-db-host
USERCENTER_DATABASE_POSTGRES_PASSWORD=secure-password
USERCENTER_JWT_SECRET=your-secure-jwt-secret
USERCENTER_SECURITY_CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

### Testing Environment

For automated testing:

```bash
# Testing environment variables
USERCENTER_ENV=test
USERCENTER_SERVER_MODE=test
USERCENTER_LOG_LEVEL=error
USERCENTER_DATABASE_POSTGRES_DBNAME=usercenter_test
USERCENTER_JWT_SECRET=test-secret-key
```

## Configuration Validation

### Startup Validation

The application validates configuration at startup:

- **Required Fields**: Ensures all required configuration is provided
- **Format Validation**: Validates configuration format and types
- **Connection Testing**: Tests database and external service connections
- **Security Checks**: Validates security-related configuration

### Validation Errors

If configuration validation fails, the application will:

1. Log detailed error messages
2. Exit with a non-zero status code
3. Provide guidance on how to fix the configuration

## Configuration Management

### Environment Files

Use `.env` files for local development:

```bash
# .env file example
USERCENTER_DATABASE_POSTGRES_HOST=localhost
USERCENTER_DATABASE_POSTGRES_PASSWORD=password
USERCENTER_JWT_SECRET=your-secret-key
```

### Configuration Templates

Provide configuration templates for different environments:

- `configs/config.example.yaml` - Example configuration file
- `configs/config.development.yaml` - Development configuration
- `configs/config.production.yaml` - Production configuration

### Secrets Management

For production deployments, use proper secrets management:

- **Environment Variables**: Pass secrets as environment variables
- **Secret Managers**: Use cloud provider secret managers
- **Configuration Management**: Use tools like Helm or Kustomize

## Configuration Best Practices

### 1. Security

- **Never commit secrets**: Keep secrets out of version control
- **Use strong secrets**: Use cryptographically strong secrets
- **Rotate secrets**: Regularly rotate sensitive configuration
- **Limit access**: Restrict access to production configuration

### 2. Environment Management

- **Environment separation**: Use different configuration for different environments
- **Configuration drift**: Prevent configuration drift between environments
- **Documentation**: Document all configuration options
- **Validation**: Validate configuration in all environments

### 3. Deployment

- **Immutable configuration**: Use immutable configuration in production
- **Configuration as code**: Version control configuration changes
- **Rollback capability**: Ensure configuration can be rolled back
- **Monitoring**: Monitor configuration changes and their impact

### 4. Development

- **Local overrides**: Allow local configuration overrides
- **Default values**: Provide sensible defaults for all options
- **Documentation**: Document configuration options and their impact
- **Validation**: Validate configuration during development

## Troubleshooting Configuration

### Common Issues

**Configuration not loaded**
- Check environment variable names and values
- Verify configuration file location and format
- Check application logs for configuration errors

**Database connection failed**
- Verify database host, port, and credentials
- Check network connectivity to database
- Ensure database is running and accessible

**Kafka connection failed**
- Verify Kafka broker addresses
- Check network connectivity to Kafka
- Ensure Kafka topics are created

**JWT authentication failed**
- Verify JWT secret is set correctly
- Check JWT token format and expiration
- Ensure JWT issuer and audience are correct

### Debugging Configuration

Enable debug logging to troubleshoot configuration issues:

```bash
USERCENTER_LOG_LEVEL=debug
USERCENTER_LOG_FORMAT=console
```

### Configuration Testing

Test configuration before deployment:

```bash
# Validate configuration
./usercenter --validate-config

# Test connections
./usercenter --test-connections
```

---

## 🔗 Related Documentation

- [Getting Started Guide](getting-started.md)
- [Deployment Guide](deployment.md)
- [Troubleshooting Guide](troubleshooting.md)
- [Architecture Guide](architecture.md) 