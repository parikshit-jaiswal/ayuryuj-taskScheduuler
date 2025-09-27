# Task Scheduler Development Guide

## 🏗️ Project Structure

```
ayuryuj-task/
├── cmd/api/                    # Application entry point
│   └── main.go                 # Main application file
├── internal/                   # Private application code
│   ├── database/              # Database connection & migrations
│   │   ├── database.go        # Database service interface
│   │   └── database_test.go   # Database tests
│   ├── handlers/              # HTTP request handlers
│   │   ├── task_handler.go    # Task-related endpoints
│   │   ├── task_result_handler.go # Task result endpoints
│   │   ├── health_handler.go  # Health check endpoints
│   │   └── *_test.go         # Handler tests
│   ├── middleware/            # HTTP middleware
│   │   └── middleware.go      # Logging, CORS, JSON middleware
│   ├── models/                # Data models
│   │   ├── task.go           # Task model & DTOs
│   │   └── task_result.go    # Task result model & DTOs
│   ├── repository/            # Data access layer
│   │   ├── task_repository.go # Task database operations
│   │   └── task_result_repository.go # Task result DB operations
│   ├── scheduler/             # Task scheduling engine
│   │   └── scheduler.go       # Main scheduler implementation
│   ├── server/                # HTTP server setup
│   │   ├── server.go         # Server initialization
│   │   └── routes.go         # Route definitions
│   └── services/              # Business logic layer
│       ├── task_service.go    # Task business logic
│       ├── task_result_service.go # Task result business logic
│       └── http_executor_service.go # HTTP request executor
├── migrations/                 # Database migrations
│   ├── 001_initial_schema.up.sql   # Create tables
│   └── 001_initial_schema.down.sql # Drop tables
├── api/                       # API documentation
│   └── openapi.yaml          # OpenAPI/Swagger specification
├── docs/                      # Additional documentation
│   └── TaskScheduler.postman_collection.json # Postman collection
├── docker-compose.yml         # Docker orchestration
├── Dockerfile                # Docker image definition
├── Makefile                  # Build automation
├── go.mod & go.sum          # Go module dependencies
├── .env & .env.example      # Environment configuration
└── README.md                # Main documentation
```

## 🔧 Development Setup

### Prerequisites
- Go 1.21 or higher
- Docker & Docker Compose
- Git
- PostgreSQL (optional, for local development)

### Initial Setup
1. **Clone the repository**
2. **Install dependencies**: `make deps`
3. **Setup environment**: `make dev-setup`
4. **Update .env file** with your configuration
5. **Start development**: `make docker-run`

## 🏃 Running the Application

### Using Docker (Recommended)
```bash
# Start everything (app + database)
make docker-run

# View logs
make docker-logs

# Stop everything
make docker-down
```

### Local Development
```bash
# Run with live reload
make watch

# Run normally
make run

# Run tests
make test
```

## 🧪 Testing

### Unit Tests
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make itest
```

### Manual Testing
Use the provided Postman collection in `docs/TaskScheduler.postman_collection.json`

## 🛠️ Development Workflow

### 1. Code Structure

**Controllers (Handlers)**
- Handle HTTP requests/responses
- Input validation
- Call appropriate services
- Return JSON responses

**Services**
- Business logic
- Validation
- Coordinate between repositories
- Handle complex operations

**Repositories**
- Data access layer
- SQL queries
- Database operations
- No business logic

**Models**
- Data structures
- Validation tags
- JSON serialization
- Database mapping

### 2. Adding New Features

#### Adding a New Endpoint
1. **Define the model** in `internal/models/`
2. **Add repository methods** in `internal/repository/`
3. **Implement service logic** in `internal/services/`
4. **Create handler** in `internal/handlers/`
5. **Register route** in `internal/server/routes.go`
6. **Write tests**
7. **Update API documentation**

#### Example: Adding Task Tags
```go
// 1. Update model
type Task struct {
    // ... existing fields
    Tags []string `json:"tags" db:"tags"`
}

// 2. Update repository
func (r *TaskRepository) GetByTag(tag string) ([]Task, error) {
    // Implementation
}

// 3. Update service
func (s *TaskService) GetTasksByTag(tag string) ([]Task, error) {
    return s.taskRepo.GetByTag(tag)
}

// 4. Add handler
func (h *TaskHandler) GetTasksByTag(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// 5. Register route
mux.HandleFunc("GET /tasks/tag/{tag}", taskHandler.GetTasksByTag)
```

### 3. Database Migrations

#### Creating a Migration
```sql
-- migrations/002_add_tags.up.sql
ALTER TABLE tasks ADD COLUMN tags TEXT[];
CREATE INDEX idx_tasks_tags ON tasks USING GIN(tags);

-- migrations/002_add_tags.down.sql
DROP INDEX IF EXISTS idx_tasks_tags;
ALTER TABLE tasks DROP COLUMN tags;
```

#### Running Migrations
Migrations run automatically on application startup, or manually:
```bash
make db-migrate    # Run pending migrations
make db-rollback   # Rollback last migration
```

## 📊 Architecture Patterns

### Dependency Injection
```go
// Server creates and wires dependencies
type Server struct {
    db           database.Service
    taskService  *services.TaskService
    // ...
}

func NewServer() *http.Server {
    db := database.New()
    taskRepo := repository.NewTaskRepository(db.GetDB())
    taskService := services.NewTaskService(taskRepo)
    // ...
}
```

### Error Handling
```go
// Service layer
func (s *TaskService) GetTask(id uuid.UUID) (*models.Task, error) {
    task, err := s.taskRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get task: %w", err)
    }
    if task == nil {
        return nil, fmt.Errorf("task not found")
    }
    return task, nil
}

// Handler layer
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    // ... parse ID
    
    task, err := h.taskService.GetTask(id)
    if err != nil {
        if err.Error() == "task not found" {
            http.Error(w, "Task not found", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(task)
}
```

## 🔒 Security Considerations

### Input Validation
- Validate all user inputs
- Use struct tags for validation
- Sanitize SQL queries (use parameterized queries)
- Validate URLs before making HTTP requests

### Authentication (Future Enhancement)
```go
// middleware/auth.go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## 📈 Performance Optimization

### Database
- Use indexes on frequently queried columns
- Implement connection pooling
- Use prepared statements
- Consider read replicas for heavy read workloads

### Application
- Use goroutines for concurrent task execution
- Implement caching where appropriate
- Use context for request timeouts
- Monitor memory usage

### Monitoring
```go
// Add metrics collection
var (
    tasksCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "tasks_created_total",
            Help: "Total number of tasks created",
        },
        []string{"status"},
    )
)
```

## 🚀 Deployment

### Production Checklist
- [ ] Environment variables configured
- [ ] Database migrations tested
- [ ] HTTPS enabled
- [ ] Logging configured
- [ ] Health checks working
- [ ] Monitoring setup
- [ ] Backup strategy in place

### Docker Production
```dockerfile
# Multi-stage build for smaller image
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o main cmd/api/main.go

FROM scratch
COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
CMD ["./main"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-scheduler
spec:
  replicas: 3
  selector:
    matchLabels:
      app: task-scheduler
  template:
    metadata:
      labels:
        app: task-scheduler
    spec:
      containers:
      - name: task-scheduler
        image: task-scheduler:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
```

## 🐛 Debugging

### Common Issues
1. **Database connection fails**
   - Check environment variables
   - Verify database is running
   - Check firewall/network settings

2. **Tasks not executing**
   - Check scheduler logs
   - Verify task trigger configuration
   - Check database for scheduled tasks

3. **HTTP requests timing out**
   - Check target URL accessibility
   - Verify timeout settings
   - Check firewall rules

### Logging
```go
log.Printf("Executing task: %s (ID: %s)", task.Name, task.ID)
log.Printf("Task result: %+v", result)
```

## 🤝 Contributing

### Code Standards
- Follow Go conventions
- Write tests for new features
- Update documentation
- Use meaningful commit messages
- Run linter before committing

### Pull Request Process
1. Fork the repository
2. Create feature branch
3. Write tests
4. Update documentation
5. Submit pull request

## 📚 Additional Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Best Practices](https://docs.docker.com/develop/best-practices/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Cron Expression Format](https://pkg.go.dev/github.com/robfig/cron)