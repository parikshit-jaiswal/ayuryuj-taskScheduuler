# 🚀 How to Access Your Task Scheduler Backend

## 🔗 **Current Status**
- ✅ Database is running on port 5432
- ⚠️  Application container needs port adjustment due to conflicts

## 🎯 **3 Ways to Access the Backend**

### **Option 1: Run Locally (Recommended)**
```powershell
# Set environment variables and run locally
cd "c:\Users\parik\Desktop\GoLang Task\ayuryuj-task"

# Set database connection to use the running Docker database
$env:BLUEPRINT_DB_HOST="localhost"
$env:BLUEPRINT_DB_PORT="5432" 
$env:BLUEPRINT_DB_DATABASE="taskscheduler"
$env:BLUEPRINT_DB_USERNAME="postgres"
$env:BLUEPRINT_DB_PASSWORD="password123"
$env:PORT="3000"

# Run the application
go run cmd/api/main.go
```

**Access at: http://localhost:3000**

### **Option 2: Fix Docker Port**
```powershell
# Stop current containers
docker-compose down

# Edit .env file - change PORT from 8081 to 3000
# Then restart
docker-compose up --build -d
```

**Access at: http://localhost:3000**

### **Option 3: Use Different Port**
```powershell
# Find available port
netstat -ano | findstr :9000
# If port 9000 is free, use it

$env:PORT="9000"
go run cmd/api/main.go
```

**Access at: http://localhost:9000**

---

## 🧪 **Test the API**

Once the server is running, test these endpoints:

### **Health Check**
```powershell
Invoke-RestMethod -Uri "http://localhost:3000/health" -Method Get
```

### **Create a Task**
```powershell
$task = @{
    name = "Test Task"
    trigger = @{
        type = "one-off"
        datetime = "2025-09-28T10:00:00Z"
    }
    action = @{
        method = "GET"
        url = "https://httpbin.org/json"
        headers = @{
            "User-Agent" = "TaskScheduler/1.0"
        }
    }
} | ConvertTo-Json -Depth 3

Invoke-RestMethod -Uri "http://localhost:3000/tasks" -Method Post -Body $task -ContentType "application/json"
```

### **List Tasks**
```powershell
Invoke-RestMethod -Uri "http://localhost:3000/tasks" -Method Get
```

---

## 📊 **API Endpoints Available**

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/tasks` | Create new task |
| GET | `/tasks` | List all tasks |
| GET | `/tasks/{id}` | Get specific task |
| PUT | `/tasks/{id}` | Update task |
| DELETE | `/tasks/{id}` | Delete task |
| GET | `/tasks/{id}/results` | Get task results |
| GET | `/results` | List all results |

---

## 🎉 **Your Backend is Ready!**

The Task Scheduler backend includes:
- ✅ Complete REST API
- ✅ PostgreSQL database (running in Docker)
- ✅ Task scheduling engine
- ✅ HTTP request execution
- ✅ Result tracking

**Recommended**: Use **Option 1** (run locally) as it's the quickest way to get started!