# Web API

A Go-based Web API service with features for weather data, authentication, and multi-tenancy support.

## Features

- Weather information endpoints
- Health check monitoring
- JWT-based authentication
- Multi-tenancy support
- Localization middleware
- OAuth integration
- Swagger API documentation

## Prerequisites

- Go 1.24.1 or higher
- PostgreSQL database
- Swagger for API documentation

## Getting Started

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Configure the application by editing `config.yaml`:
```yaml
server:
  port: 8001

database:
  driver: postgres
  host: localhost
  port: 5432
  database: webapi
  username: your_username
  password: your_password

logging:
  level: info
  format: json
```

4. Run the application:
```bash
go run main.go
```

## API Documentation

The API is documented using Swagger/OpenAPI. To generate the documentation:

Windows:
```bash
./scripts/generate-swagger.bat
```

Linux/Mac:
```bash
./scripts/generate-swagger.sh
```

After starting the application, access the Swagger UI at:
```
http://localhost:8001/swagger/index.html
```

## Available Endpoints

### Health Check
- `GET /health` - Service health status

### Weather
- Weather-related endpoints (check Swagger documentation for details)

### Authentication
- Authentication endpoints integrated with the Identity Server

## Development

The application uses:
- Gin web framework
- GORM for database operations
- Logrus for logging
- Viper for configuration management
- JWT for authentication
- Swagger for API documentation