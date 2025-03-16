# Identity Server

A Go-based OAuth2 and Authentication Server that provides secure user authentication and authorization services.

## Features

- User authentication and management
- OAuth2 authorization server
- JWT-based authentication
- Client application registration
- Supports multiple OAuth2 grant types:
  - Password grant
  - Authorization code grant
  - Client credentials grant

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
  port: 8000

database:
  driver: postgres
  host: localhost
  port: 5432
  database: identity_server
  username: your_username
  password: your_password

auth:
  jwtSecret: "your-secure-jwt-secret"

logging:
  filePath: ./logs/identity-server.log
```

4. Run the application:
```bash
go run main.go
```

## API Endpoints

### Public Routes
- `POST /oauth/clients` - Register new OAuth client
- `GET /oauth/authorize` - OAuth authorization endpoint
- `POST /users` - Create new user
- `POST /login` - User login

### Protected Routes
- `POST /oauth/token` - Get OAuth access token

## Generate Swagger Documentation

Windows:
```bash
./scripts/generate-swagger.bat
```

Linux/Mac:
```bash
./scripts/generate-swagger.sh
```

## Development

The application uses:
- Gin web framework
- GORM for database operations
- Viper for configuration management
- JWT for authentication
- Swagger for API documentation