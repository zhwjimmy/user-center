# UserCenter Documentation

Welcome to the UserCenter documentation! This guide will help you navigate through all available documentation based on your needs.

## 📚 Documentation Overview

UserCenter is a production-ready user management service built with Go, featuring clean architecture design and event-driven patterns. Our documentation is organized to help different types of users find the information they need quickly.

## 🎯 For New Users

### Getting Started
- **[Quick Start Guide](getting-started.md)** - Get up and running in 10 minutes
- **[Architecture Overview](architecture.md)** - Understand the system design
- **[Configuration Guide](configuration.md)** - Learn about configuration options

### First Steps
1. Read the [Quick Start Guide](getting-started.md) to set up your development environment
2. Review the [Architecture Overview](architecture.md) to understand the system design
3. Check the [API Reference](api-reference.md) to see available endpoints

## 🛠️ For Developers

### Development Guides
- **[Development Guide](development.md)** - Complete development workflow
- **[API Reference](api-reference.md)** - Detailed API documentation
- **[Testing Guide](testing.md)** - Testing strategies and best practices
- **[Contributing Guide](contributing.md)** - How to contribute to the project

### Architecture & Design
- **[Architecture Design](architecture.md)** - Clean architecture principles and implementation
- **[Event-Driven Architecture](kafka-integration.md)** - Kafka integration and event processing
- **[Database Design](database-design.md)** - Database schema and migrations

### Advanced Topics
- **[Performance Optimization](performance.md)** - Performance tuning and optimization
- **[Security Best Practices](security.md)** - Security considerations and best practices
- **[Monitoring & Observability](monitoring.md)** - Logging, metrics, and tracing

## 🚀 For DevOps & Operations

### Deployment
- **[Deployment Guide](deployment.md)** - Production deployment instructions
- **[Docker Setup](docker-setup.md)** - Container deployment with Docker
- **[Kubernetes Deployment](kubernetes.md)** - Kubernetes deployment manifests

### Operations
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
- **[Monitoring Guide](monitoring.md)** - Health checks, metrics, and alerting
- **[Backup & Recovery](backup-recovery.md)** - Data backup and disaster recovery

### CI/CD
- **[GitHub Actions](github-actions.md)** - CI/CD pipeline configuration
- **[Release Process](release-process.md)** - Versioning and release management

## 📊 For System Administrators

### Infrastructure
- **[Infrastructure Setup](infrastructure.md)** - Server and service setup
- **[Database Management](database-management.md)** - PostgreSQL and MongoDB administration
- **[Kafka Management](kafka-management.md)** - Apache Kafka setup and maintenance

### Security
- **[Security Hardening](security-hardening.md)** - Security configuration and hardening
- **[Access Control](access-control.md)** - User authentication and authorization
- **[Audit Logging](audit-logging.md)** - Security audit and compliance

## 🔍 Documentation by Category

### Architecture & Design
- [Architecture Overview](architecture.md) - System architecture and design principles
- [Clean Architecture](clean-architecture.md) - Clean architecture implementation details
- [Event-Driven Design](event-driven-design.md) - Event-driven architecture patterns
- [Database Design](database-design.md) - Database schema and relationships

### Development
- [Getting Started](getting-started.md) - Quick start for new developers
- [Development Guide](development.md) - Complete development workflow
- [API Reference](api-reference.md) - RESTful API documentation
- [Testing Guide](testing.md) - Testing strategies and examples
- [Code Standards](code-standards.md) - Coding conventions and best practices

### Configuration & Setup
- [Configuration Guide](configuration.md) - Configuration options and environment variables
- [Environment Setup](environment-setup.md) - Development and production environment setup
- [Docker Setup](docker-setup.md) - Container-based deployment
- [Local Development](local-development.md) - Local development environment

### Deployment & Operations
- [Deployment Guide](deployment.md) - Production deployment instructions
- [Kubernetes Deployment](kubernetes.md) - Kubernetes manifests and configuration
- [Monitoring & Observability](monitoring.md) - Logging, metrics, and health checks
- [Troubleshooting](troubleshooting.md) - Common issues and solutions

### Integration & APIs
- [Kafka Integration](kafka-integration.md) - Apache Kafka setup and event processing
- [External APIs](external-apis.md) - Integration with external services
- [Webhook Configuration](webhooks.md) - Webhook setup and configuration

### Security & Compliance
- [Security Guide](security.md) - Security best practices and considerations
- [Authentication & Authorization](auth.md) - JWT authentication and role-based access
- [Data Protection](data-protection.md) - Data privacy and protection measures

## 📖 Documentation Structure

```
docs/
├── README.md                    # This file - Documentation overview
├── architecture.md              # System architecture
├── getting-started.md           # Quick start guide
├── development.md               # Development guide
├── api-reference.md             # API documentation
├── configuration.md             # Configuration guide
├── deployment.md                # Deployment instructions
├── troubleshooting.md           # Troubleshooting guide
├── kafka-integration.md         # Kafka integration
├── contributing.md              # Contributing guidelines
├── testing.md                   # Testing guide
├── monitoring.md                # Monitoring and observability
├── security.md                  # Security guide
├── performance.md               # Performance optimization
├── docker-setup.md              # Docker configuration
├── kubernetes.md                # Kubernetes deployment
├── github-actions.md            # CI/CD pipeline
├── iterations/                  # Versioned feature documentation
│   └── README.md
└── assets/                      # Documentation assets
    ├── images/
    └── diagrams/
```

## 🔄 Documentation Maintenance

### Version Control
- All documentation is version-controlled with the codebase
- Documentation changes are reviewed as part of the PR process
- Breaking changes require documentation updates

### Contributing to Documentation
- Follow the [Contributing Guide](contributing.md) for documentation updates
- Use clear, concise language
- Include code examples where appropriate
- Keep documentation up-to-date with code changes

### Documentation Standards
- Use Markdown format for all documentation
- Include table of contents for long documents
- Use consistent formatting and structure
- Provide both English and Chinese versions where appropriate

## 🆘 Getting Help

### Support Channels
- **GitHub Issues**: [Create an issue](https://github.com/zhwjimmy/user-center/issues) for bugs or feature requests
- **GitHub Discussions**: [Join discussions](https://github.com/zhwjimmy/user-center/discussions) for questions and ideas
- **Documentation Issues**: Report documentation problems via GitHub Issues

### Community Resources
- **Code Examples**: Check the [examples](../examples/) directory
- **Sample Applications**: Review sample integrations
- **Blog Posts**: Read our technical blog for deep dives

## 📝 Documentation Feedback

We value your feedback on our documentation! If you find any issues or have suggestions for improvement:

1. **Create an issue** with the `documentation` label
2. **Submit a PR** with your proposed changes
3. **Join discussions** to share your thoughts

Your feedback helps us make our documentation better for everyone!

---

## 🔗 Quick Links

- [Main README](../README.md) - Project overview
- [API Documentation](api-reference.md) - RESTful API reference
- [Architecture](architecture.md) - System design
- [Contributing](contributing.md) - How to contribute
- [GitHub Repository](https://github.com/zhwjimmy/user-center) - Source code 