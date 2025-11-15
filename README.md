# Palmvue Backend Revamp

A production-ready Go backend application built with clean architecture principles.

## Features

- **Clean Architecture**: Separation of concerns (handler → service → repository)
- **Standardized API Responses**: Consistent JSON response format
- **Error Handling**: Centralized error management with custom error codes
- **Request Validation**: go-playground/validator integration
- **PostgreSQL**: Database with pgx driver
- **Structured Logging**: Zap logger
- **Database Migrations**: golang-migrate
- **Testing**: Comprehensive unit tests with mockery
- **Docker Support**: Docker Compose for local development

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models (business layer)
│   ├── entity/          # Database entities
│   ├── errs/            # Error definitions and handling
│   ├── handler/         # HTTP handlers
│   │   ├── request/     # Request DTOs
│   │   └── response/    # Response DTOs
│   ├── repository/      # Data access layer
│   └── service/         # Business logic
│       └── mocks/       # Generated mocks (mockery)
├── pkg/                 # Public shared packages
│   ├── database/        # Database utilities
│   ├── logger/          # Logger setup
│   └── validation/      # Request validation
├── db/
│   └── migrations/      # Database migrations
├── .mockery.yml         # Mockery configuration
└── docker-compose.yml   # Docker services
```

## Prerequisites

- Go 1.25.3 or higher
- Docker and Docker Compose
- PostgreSQL 15
- golang-migrate (optional)
- mockery v3.6.0 (for generating mocks)

## Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd palmvue-revamp-backend
```

### 2. Copy environment configuration

```bash
cp .env.example .env
```

Edit `.env` if needed to match your environment.

### 3. Install development tools (optional)

```bash
make install-tools
```

This installs:

- golangci-lint
- golang-migrate

### 4. Start the database

```bash
make docker-up
```

This starts a PostgreSQL container using Docker Compose.

### 5. Run database migrations

```bash
make migrate-up
```

### 6. Run the application

```bash
make run
```

The server will start on `http://localhost:8080`

### 7. Test the API

```bash
curl -X POST http://localhost:8080/ping \
  -H "Content-Type: application/json" \
  -d '{"include_db": false}'
```

Expected response:

```json
{
  "success": true,
  "message": "pong",
  "data": {
    "message": "pong"
  }
}
```

## Available Make Commands

| Command                           | Description                  |
| --------------------------------- | ---------------------------- |
| `make run`                        | Run the application          |
| `make test`                       | Run tests with coverage      |
| `make build`                      | Build the binary             |
| `make migrate-up`                 | Run database migrations      |
| `make migrate-down`               | Rollback database migrations |
| `make migrate-create name=<name>` | Create a new migration       |
| `make docker-up`                  | Start Docker services        |
| `make docker-down`                | Stop Docker services         |
| `make clean`                      | Clean build artifacts        |
| `make tidy`                       | Tidy go.mod                  |

## API Endpoints

### Health Check

**POST** `/ping`

Check API and database health.

**Request Body:**

```json
{
  "include_db": false
}
```

**Response:**

```json
{
  "success": true,
  "message": "pong",
  "data": {
    "message": "pong"
  }
}
```

**Error Response:**

```json
{
  "success": false,
  "error": {
    "code": "database_unavailable",
    "message": "Database is unavailable",
    "details": "connection refused"
  }
}
```

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options:

- `SERVER_PORT`: HTTP server port (default: 8080)
- `SERVER_HOST`: Server host (default: 0.0.0.0)
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: postgres)
- `DB_NAME`: Database name (default: palmvue)
- `DB_SSLMODE`: Database SSL mode (default: disable)
- `APP_ENV`: Application environment (default: development)
- `LOG_LEVEL`: Logging level (default: debug)

## Testing

Run all tests:

```bash
make test
```

Run tests with verbose output:

```bash
go test -v ./...
```

### Generate Mocks

Generate mocks for interfaces using mockery:

```bash
mockery
```

Configuration is in `.mockery.yml`

## Database Migrations

### Create a new migration

```bash
make migrate-create name=add_users_table
```

This creates two files:

- `db/migrations/000002_add_users_table.up.sql`
- `db/migrations/000002_add_users_table.down.sql`

### Apply migrations

```bash
make migrate-up
```

### Rollback migrations

```bash
make migrate-down
```

## Development Workflow

1. Start the database:

   ```bash
   make docker-up
   ```

2. Run migrations:

   ```bash
   make migrate-up
   ```

3. Start the server:

   ```bash
   make run
   ```

4. Make your changes

5. Generate mocks (if interfaces changed):

   ```bash
   mockery
   ```

6. Run tests:
   ```bash
   make test
   ```

## Production Build

```bash
make build
```

The binary will be created at `bin/palmvue-server`

Run the binary:

```bash
./bin/palmvue-server
```

## Docker

### Start services

```bash
make docker-up
```

### Stop services

```bash
make docker-down
```
