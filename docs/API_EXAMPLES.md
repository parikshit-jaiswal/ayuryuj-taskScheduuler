# Task Scheduler API Examples

This document provides practical examples of how to use the Task Scheduler API.

## 🚀 Getting Started

### 1. Start the Application

```bash
# Using Docker (Recommended)
make docker-run

# Or locally (requires PostgreSQL)
make run
```

The API will be available at `http://localhost:8080`

### 2. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "up",
  "message": "It's healthy"
}
```

## 📝 Creating Tasks

### Example 1: One-time HTTP GET Request

Create a task that will make a GET request to httpbin.org once:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Check API Status",
    "trigger": {
      "type": "one-off",
      "datetime": "2025-09-28T10:00:00Z"
    },
    "action": {
      "method": "GET",
      "url": "https://httpbin.org/json",
      "headers": {
        "User-Agent": "TaskScheduler/1.0"
      }
    }
  }'
```

### Example 2: Recurring Daily Task

Create a task that runs every weekday at 9 AM:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Daily Health Check",
    "trigger": {
      "type": "cron",
      "cron": "0 9 * * 1-5"
    },
    "action": {
      "method": "GET",
      "url": "https://httpbin.org/status/200",
      "headers": {
        "Authorization": "Bearer your-api-key"
      }
    }
  }'
```

### Example 3: POST Request with JSON Payload

Create a task that sends data to an API:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Send Daily Report",
    "trigger": {
      "type": "cron",
      "cron": "0 18 * * 1-5"
    },
    "action": {
      "method": "POST",
      "url": "https://httpbin.org/post",
      "headers": {
        "Content-Type": "application/json",
        "X-API-Key": "your-secret-key"
      },
      "payload": {
        "report_type": "daily",
        "timestamp": "{{current_time}}",
        "data": {
          "metrics": {
            "tasks_completed": 42,
            "success_rate": 98.5
          }
        }
      }
    }
  }'
```

## 📊 Managing Tasks

### List All Tasks

```bash
curl http://localhost:8080/tasks
```

### List Tasks with Pagination

```bash
curl "http://localhost:8080/tasks?page=1&page_size=10"
```

### Filter Tasks by Status

```bash
# Get only scheduled tasks
curl "http://localhost:8080/tasks?status=scheduled"

# Get cancelled tasks
curl "http://localhost:8080/tasks?status=cancelled"
```

### Get Specific Task

```bash
# Replace {task-id} with actual UUID
curl http://localhost:8080/tasks/{task-id}
```

### Update a Task

```bash
curl -X PUT http://localhost:8080/tasks/{task-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Task Name",
    "trigger": {
      "type": "cron",
      "cron": "0 10 * * 1-5"
    }
  }'
```

### Cancel a Task

```bash
curl -X DELETE http://localhost:8080/tasks/{task-id}
```

## 📈 Viewing Results

### Get Results for a Specific Task

```bash
curl http://localhost:8080/tasks/{task-id}/results
```

### Get Results with Filters

```bash
# Get only successful results
curl "http://localhost:8080/tasks/{task-id}/results?success=true"

# Get results from a date range
curl "http://localhost:8080/tasks/{task-id}/results?date_from=2025-09-01T00:00:00Z&date_to=2025-09-30T23:59:59Z"
```

### List All Results

```bash
curl http://localhost:8080/results
```

### Filter All Results

```bash
# Get all failed executions
curl "http://localhost:8080/results?success=false"

# Get results for specific task
curl "http://localhost:8080/results?task_id={task-id}"
```

## 🕐 Cron Expression Examples

| Expression | Description |
|------------|-------------|
| `0 9 * * 1-5` | Every weekday at 9:00 AM |
| `0 */6 * * *` | Every 6 hours |
| `30 2 * * 0` | Every Sunday at 2:30 AM |
| `0 0 1 * *` | First day of every month at midnight |
| `*/15 * * * *` | Every 15 minutes |
| `0 12 * * 1` | Every Monday at noon |
| `0 0 * * 0` | Every Sunday at midnight |

## 🧪 Testing Scenarios

### Immediate Execution Test

Create a task that executes immediately (or very soon):

```bash
# Set datetime to 1 minute from now
FUTURE_TIME=$(date -u -d '+1 minute' '+%Y-%m-%dT%H:%M:%SZ')

curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Immediate Test\",
    \"trigger\": {
      \"type\": \"one-off\",
      \"datetime\": \"$FUTURE_TIME\"
    },
    \"action\": {
      \"method\": \"GET\",
      \"url\": \"https://httpbin.org/uuid\"
    }
  }"
```

### Error Handling Test

Test how the system handles failed requests:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Error Test",
    "trigger": {
      "type": "one-off",
      "datetime": "2025-09-28T10:05:00Z"
    },
    "action": {
      "method": "GET",
      "url": "https://httpbin.org/status/500"
    }
  }'
```

### Timeout Test

Test request timeouts:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Timeout Test",
    "trigger": {
      "type": "one-off",
      "datetime": "2025-09-28T10:10:00Z"
    },
    "action": {
      "method": "GET",
      "url": "https://httpbin.org/delay/45"
    }
  }'
```

## 📋 Sample Responses

### Task Creation Response

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Daily Health Check",
  "trigger": {
    "type": "cron",
    "cron": "0 9 * * 1-5"
  },
  "action": {
    "method": "GET",
    "url": "https://httpbin.org/status/200",
    "headers": {
      "Authorization": "Bearer your-api-key"
    }
  },
  "status": "scheduled",
  "created_at": "2025-09-27T10:00:00Z",
  "updated_at": "2025-09-27T10:00:00Z",
  "next_run": "2025-09-28T09:00:00Z"
}
```

### Task Execution Result

```json
{
  "id": "c27493e8-1234-5678-9abc-def012345678",
  "task_id": "123e4567-e89b-12d3-a456-426614174000",
  "run_at": "2025-09-28T09:00:00Z",
  "status_code": 200,
  "success": true,
  "response_headers": {
    "Content-Type": "application/json",
    "Server": "gunicorn/19.9.0"
  },
  "response_body": "{\"status\": \"ok\"}",
  "error_message": null,
  "duration_ms": 245,
  "created_at": "2025-09-28T09:00:00Z"
}
```

### Error Response

```json
{
  "id": "d38594f9-2345-6789-abcd-ef0123456789",
  "task_id": "123e4567-e89b-12d3-a456-426614174000",
  "run_at": "2025-09-28T09:00:00Z",
  "status_code": 0,
  "success": false,
  "response_headers": null,
  "response_body": "",
  "error_message": "Get \"https://invalid-url.com\": dial tcp: lookup invalid-url.com: no such host",
  "duration_ms": 5000,
  "created_at": "2025-09-28T09:00:00Z"
}
```

## 🔧 PowerShell Examples (Windows)

For Windows users using PowerShell:

### Create a Task

```powershell
$headers = @{
    'Content-Type' = 'application/json'
}

$body = @{
    name = "PowerShell Test Task"
    trigger = @{
        type = "one-off"
        datetime = "2025-09-28T10:00:00Z"
    }
    action = @{
        method = "GET"
        url = "https://httpbin.org/json"
        headers = @{
            'User-Agent' = 'PowerShell/7.0'
        }
    }
} | ConvertTo-Json -Depth 3

Invoke-RestMethod -Uri "http://localhost:8080/tasks" -Method Post -Headers $headers -Body $body
```

### Get Task Results

```powershell
$taskId = "your-task-id-here"
Invoke-RestMethod -Uri "http://localhost:8080/tasks/$taskId/results" -Method Get
```

## 🎯 Real-World Use Cases

### 1. Website Health Monitoring

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Website Health Check",
    "trigger": {
      "type": "cron",
      "cron": "*/5 * * * *"
    },
    "action": {
      "method": "GET",
      "url": "https://your-website.com/health",
      "headers": {
        "User-Agent": "HealthMonitor/1.0"
      }
    }
  }'
```

### 2. API Data Sync

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sync User Data",
    "trigger": {
      "type": "cron",
      "cron": "0 2 * * *"
    },
    "action": {
      "method": "POST",
      "url": "https://api.your-service.com/sync",
      "headers": {
        "Authorization": "Bearer your-token",
        "Content-Type": "application/json"
      },
      "payload": {
        "sync_type": "users",
        "full_sync": true
      }
    }
  }'
```

### 3. Webhook Trigger

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Daily Report Webhook",
    "trigger": {
      "type": "cron",
      "cron": "0 18 * * 1-5"
    },
    "action": {
      "method": "POST",
      "url": "https://hooks.slack.com/services/your/webhook/url",
      "headers": {
        "Content-Type": "application/json"
      },
      "payload": {
        "text": "Daily report: All systems operational"
      }
    }
  }'
```

This comprehensive guide should help you get started with the Task Scheduler API!