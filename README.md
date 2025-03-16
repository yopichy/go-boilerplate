# Go Microservices Boilerplate

A production-ready microservices boilerplate to accelerate your Go application development. This project provides a solid foundation with two main services that you can use as a starting point for building robust, scalable applications.

## Project Overview

This boilerplate consists of two main services:

### 1. Identity Server (Authentication Service)
A complete authentication and authorization service that provides:
- User authentication and management
- OAuth2 authorization server
- JWT-based authentication
- Client application registration
- Multiple OAuth2 grant types support
- Secure user management

### 2. Web API (Resource Service)
A feature-rich API service template that includes:
- Ready-to-use API structure
- Multi-tenancy support
- Localization middleware
- OAuth integration
- Health monitoring
- Example weather endpoints
- Swagger documentation

## Key Features

- **Production-Ready Architecture**: Follows best practices and modern Go development patterns
- **Microservices Design**: Clear separation of concerns between authentication and resource servers
- **Comprehensive Security**: Built-in JWT authentication, OAuth2 support, and secure user management
- **Developer Experience**:
  - Detailed Swagger documentation
  - Easy-to-use configuration system
  - Clear project structure
  - Ready-to-use middleware components
  - Database integration with GORM
  - Logging and monitoring included

## Technology Stack

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Documentation**: Swagger/OpenAPI
- **Configuration**: Viper
- **Authentication**: JWT, OAuth2
- **Logging**: Structured logging with support for multiple formats
- **Development Tools**: Built-in scripts for documentation generation

## Quick Start

1. Clone the repository
2. Each service can be run independently:

   Identity Server (Default port: 8000):
   ```bash
   cd identity-server
   go mod download
   go run main.go
   ```

   Web API (Default port: 8001):
   ```bash
   cd webapi
   go mod download
   go run main.go
   ```

3. Configure your services by editing their respective `config.yaml` files

## Project Structure

```
├── identity-server/         # Authentication & Authorization Service
│   ├── config/             # Configuration management
│   ├── database/           # Database initialization
│   ├── internal/           # Internal packages
│   │   ├── handlers/       # HTTP handlers
│   │   ├── models/         # Data models
│   │   └── services/       # Business logic
│   └── scripts/            # Utility scripts
│
├── webapi/                 # Resource Service
    ├── config/             # Configuration management
    ├── database/           # Database initialization
    ├── docs/               # API documentation
    ├── internal/           # Internal packages
    │   ├── api/           # API definitions
    │   ├── handlers/      # HTTP handlers
    │   ├── middleware/    # Custom middleware
    │   └── models/        # Data models
    └── scripts/           # Utility scripts
```

## Getting Started with Development

1. Review the documentation in each service's directory
2. Configure your database settings in the respective `config.yaml` files
3. Generate API documentation using the provided scripts
4. Start with the identity-server for authentication
5. Configure the webapi to work with the identity-server
6. Build your additional services following the established patterns

## Contributing

Feel free to:
- Submit bug reports
- Propose new features
- Create pull requests

## Roadmap

The following features and improvements are planned for future releases:

### Observability & Monitoring
- [ ] Elasticsearch integration for telemetry and logging
- [ ] Distributed tracing with OpenTelemetry
- [ ] Prometheus metrics integration
- [ ] Grafana dashboards templates
- [ ] Enhanced logging with structured logging patterns

### Database & Storage
- [ ] Database provider abstraction layer
  - [ ] Support for MySQL
  - [ ] Support for MongoDB
  - [ ] Support for SQLite
  - [ ] Easy database switch through configuration
- [ ] Database migration tools integration
- [ ] Connection pooling optimization
- [ ] Read/Write splitting support

### Security & Authentication
- [ ] Enhanced rate limiting
- [ ] API key management system
- [ ] Support for additional OAuth2 providers
- [ ] Two-factor authentication (2FA)
- [ ] Role-based access control (RBAC) improvements

### Development Experience
- [ ] Docker compose setup for local development
- [ ] Kubernetes deployment templates
- [ ] CI/CD pipeline templates
- [ ] Development environment automation
- [ ] Enhanced API testing tools
- [ ] Performance testing suite

### Scalability & Performance
- [ ] Caching layer with Redis
- [ ] Message queue integration (RabbitMQ/Kafka)
- [ ] Horizontal scaling support
- [ ] Load balancing configurations
- [ ] Circuit breaker implementation

### Additional Features
- [ ] WebSocket support
- [ ] GraphQL API support
- [ ] gRPC service templates
- [ ] File upload handling
- [ ] Email service integration
- [ ] Notification service
- [ ] Background job processing

## License

MIT License - feel free to use this boilerplate for your projects.