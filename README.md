# Task Scheduler Backend

A production-ready, RESTful Task Scheduler Service built with Go that manages HTTP tasks which can be scheduled (one-off or recurring) and reliably executed. All tasks and their results are persistently stored in PostgreSQL.

## 🚀 Features

- **Task Management**: Create, read, update, and cancel HTTP tasks
- **Flexible Scheduling**: Support for both one-off (datetime) and recurring (cron) tasks
- **Persistent Storage**: PostgreSQL database with proper migrations
- **Task Execution**: Reliable HTTP request execution with detailed logging
- **Result Tracking**: Complete execution history with response details
- **RESTful API**: Full REST API with OpenAPI/Swagger documentation
- **Production Ready**: Docker containerization, health checks, and graceful shutdown
- **Error Handling**: Comprehensive error handling and logging
- **Pagination**: Efficient pagination for large datasets

## 🏗️ Architecture

```
ayuryuj-task/
├── cmd/api/                    # Application entry point
├── internal/                   # Private application packages
│   ├── database/              # Database service and connection management
│   ├── handlers/              # HTTP request handlers (controllers)
│   ├── middleware/            # HTTP middleware components
│   ├── models/                # Domain models and data structures
│   ├── repository/            # Data access layer (repository pattern)
│   ├── scheduler/             # Task scheduling engine
│   ├── server/                # HTTP server configuration
│   └── services/              # Business logic layer
├── migrations/                 # Database migration files (up/down)
├── api/                       # OpenAPI/Swagger specifications
├── docs/                      # Project documentation
├── postman/                   # Postman collections for API testing
├── bin/                       # Compiled binaries (gitignored)
├── logs/                      # Application logs (gitignored)
├── docker-compose.yml         # Docker orchestration
├── Dockerfile                 # Container build instructions
└── Makefile                   # Development and build commands
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally without Docker)

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone and navigate to the project**:
   ```bash
   cd ayuryuj-task
   ```

2. **Copy environment variables**:
   ```bash
   cp .env.example .env
   ```

3. **Start the application**:
   ```bash
   make docker-run
   ```

   Or manually with docker-compose:
   ```bash
   docker-compose up --build
   ```

4. **The API will be available at**: `http://localhost:8080`

### Local Development

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL database locally and update `.env` file**

3. **Run database migrations**:
   ```bash
   # Migrations run automatically when the application starts
   ```

4. **Start the application**:
   ```bash
   make run
   ```

## 📝 API Documentation

### OpenAPI/Swagger
- **Specification**: [`api/openapi.yaml`](api/openapi.yaml)
- **Interactive Docs**: Visit `http://localhost:8080/docs` (when implemented with Swagger UI)

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/tasks` | Create a new task |
| `GET` | `/tasks` | List all tasks (with pagination/filtering) |
| `GET` | `/tasks/{id}` | Get a specific task |
| `PUT` | `/tasks/{id}` | Update a task |
| `DELETE` | `/tasks/{id}` | Cancel a task |
| `GET` | `/tasks/{id}/results` | Get execution results for a task |
| `GET` | `/results` | List all execution results |
| `GET` | `/health` | Health check endpoint |
| `GET` | `/metrics` | Service metrics |

### Task Entity

```json
{
  "id": "uuid",
  "name": "string",
  "trigger": {
    "type": "one-off|cron",
    "datetime": "2025-09-28T10:00:00Z", // for one-off
    "cron": "0 9 * * 1-5"              // for cron
  },
  "action": {
    "method": "GET|POST|PUT|DELETE|...",
    "url": "https://api.example.com/endpoint",
    "headers": {"key": "value"},
    "payload": {"data": "value"}
  },
  "status": "scheduled|cancelled|completed",
  "created_at": "2025-09-27T10:00:00Z",
  "updated_at": "2025-09-27T10:00:00Z",
  "next_run": "2025-09-28T10:00:00Z"
}
```

### Example API Calls

#### Create a One-off Task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "One-time API call",
    "trigger": {
      "type": "one-off",
      "datetime": "2025-09-28T10:00:00Z"
    },
    "action": {
      "method": "GET",
      "url": "https://httpbin.org/json",
      "headers": {"User-Agent": "TaskScheduler/1.0"}
    }
  }'
```

#### Create a Recurring Task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Daily health check",
    "trigger": {
      "type": "cron",
      "cron": "0 9 * * 1-5"
    },
    "action": {
      "method": "GET",
      "url": "https://httpbin.org/status/200"
    }
  }'
```

#### List Tasks with Filtering
```bash
# Get all scheduled tasks
curl "http://localhost:8080/tasks?status=scheduled&page=1&page_size=10"

# Get all tasks
curl "http://localhost:8080/tasks"
```

## 🗄️ Database Schema

### Tasks Table
- `id` (UUID, PK)
- `name` (VARCHAR)
- `trigger` (JSONB)
- `action` (JSONB)
- `status` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `next_run` (TIMESTAMP)

### Task Results Table
- `id` (UUID, PK)
- `task_id` (UUID, FK)
- `run_at` (TIMESTAMP)
- `status_code` (INTEGER)
- `success` (BOOLEAN)
- `response_headers` (JSONB)
- `response_body` (TEXT)
- `error_message` (TEXT)
- `duration_ms` (BIGINT)
- `created_at` (TIMESTAMP)

## 🔧 Configuration

Configuration is handled through environment variables:

```env
# Application
APP_ENV=development
PORT=8080

# Database
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=taskscheduler
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=password123
BLUEPRINT_DB_SCHEMA=public
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run integration tests
make itest

# Run tests with coverage
go test ./... -v -cover
```

## 📊 Development Commands

The project includes a comprehensive Makefile:

```bash
make build        # Build the application
make run          # Run the application
make test         # Run tests
make itest        # Run integration tests
make docker-run   # Start with Docker Compose
make docker-down  # Stop Docker containers
make watch        # Live reload during development
make clean        # Clean build artifacts
```

## 🔍 Monitoring and Health

### Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "up",
  "message": "It's healthy",
  "open_connections": "1",
  "in_use": "0",
  "idle": "1"
}
```

### Metrics
```bash
curl http://localhost:8080/metrics
```

## 🔄 Task Scheduling

The application includes a built-in scheduler that:

1. **Polls for Ready Tasks**: Checks every 10 seconds for tasks ready to execute
2. **Executes HTTP Requests**: Performs the configured HTTP action
3. **Records Results**: Stores detailed execution results
4. **Updates Task State**: 
   - One-off tasks → marked as `completed`
   - Cron tasks → next run time calculated
5. **Error Handling**: Captures and logs execution errors

### Cron Expression Examples

- `0 9 * * 1-5` - Every weekday at 9:00 AM
- `0 */6 * * *` - Every 6 hours
- `30 2 * * 0` - Every Sunday at 2:30 AM
- `0 0 1 * *` - First day of every month at midnight

## 🚀 Production Deployment

### Docker Production Build

```dockerfile
# Multi-stage build for optimized production image
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment Variables for Production

```env
APP_ENV=production
PORT=8080
BLUEPRINT_DB_HOST=your-postgres-host
BLUEPRINT_DB_DATABASE=taskscheduler_prod
# ... other production configs
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For support, please check:
1. This README
2. OpenAPI documentation in `api/openapi.yaml`
3. Create an issue in the repository

## 📋 Task Result Example

When a task executes, a result record like this is created:

```json
{
  "id": "c27493e8-1234-5678-9abc-def012345678",
  "task_id": "a2349de4-1234-5678-9abc-def012345678",
  "run_at": "2025-07-01T10:00:00Z",
  "status_code": 200,
  "success": true,
  "response_headers": {
    "Content-Type": "application/json"
  },
  "response_body": "{\"message\":\"ok\"}",
  "error_message": null,
  "duration_ms": 233,
  "created_at": "2025-07-01T10:00:00Z"
}
```
```bash
make clean
```
