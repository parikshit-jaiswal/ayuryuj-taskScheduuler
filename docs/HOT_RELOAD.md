# 🔥 Hot Reload Development with Air

This project is configured with [Air](https://github.com/air-verse/air) for automatic server restart during development.

## 🚀 Quick Start

### Option 1: Using Makefile (Recommended)
```bash
make watch
```

### Option 2: Direct Air Command
```bash
air
```

### Option 3: For Windows PowerShell
```powershell
# If make is not available
air
```

## ⚙️ Configuration

The `.air.toml` file contains the hot reload configuration:

- **Watch Directories**: `cmd/`, `internal/`
- **Watch Extensions**: `.go`, `.yml`, `.yaml`, `.env`
- **Exclude**: Test files, build artifacts, vendor, node_modules
- **Build Command**: `go build -o ./tmp/main cmd/api/main.go`
- **Auto Restart**: On any file change

## 📁 What Happens

1. **File Watch**: Air monitors files in `cmd/` and `internal/` directories
2. **Auto Build**: When you save changes, Air automatically builds the project
3. **Auto Restart**: The server restarts with your changes
4. **Clean Logs**: Clear output on each rebuild for better readability

## 🎯 Development Workflow

```bash
# 1. Start hot reload
make watch

# 2. Edit any .go file in cmd/ or internal/
# 3. Save the file
# 4. Air automatically:
#    - Detects the change
#    - Rebuilds the application
#    - Restarts the server
#    - Shows you the logs
```

## 📋 Air Features Enabled

- ✅ **Fast Rebuild**: Only rebuilds when files change
- ✅ **Clean Logs**: Clear screen on each rebuild
- ✅ **Colored Output**: Colored logs for better readability
- ✅ **Error Handling**: Shows build errors clearly
- ✅ **Graceful Restart**: Properly stops and starts the server
- ✅ **Smart Watching**: Ignores test files and build artifacts

## 🔧 Configuration Details

```toml
[build]
  cmd = "go build -o ./tmp/main cmd/api/main.go"
  bin = "./tmp/main"
  include_dir = ["cmd", "internal"]
  include_ext = ["go", "yml", "yaml", "env"]
  exclude_dir = ["tmp", "vendor", "testdata", "node_modules", "bin", "logs"]
  exclude_regex = ["_test.go"]
```

## 🚨 Troubleshooting

### Air not found
```bash
# Install Air
go install github.com/air-verse/air@latest

# Or use make command (auto-installs)
make watch
```

### Port already in use
```bash
# Stop any running instances
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Or change port in .env file
PORT=8081
```

### Permission denied
```bash
# On Windows, run PowerShell as Administrator
# Or use Windows Defender exclusions for the project folder
```

## 📊 Example Output

```
$ make watch
Starting Task Scheduler with hot reload...
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.63.0, built with Go 1.21+

watching .
!exclude tmp
!exclude vendor
!exclude testdata
building...
running...

2025/09/27 22:30:15 Server starting on :8080
2025/09/27 22:30:15 🚀 Task Scheduler starting...
2025/09/27 22:30:15 📊 Database connection established
2025/09/27 22:30:15 ⚙️  Task scheduler started
2025/09/27 22:30:15 🌐 Server listening on http://localhost:8080

# When you save a file:
main.go has changed
building...
running...
2025/09/27 22:30:45 Server starting on :8080
# ... continues with new logs
```

## 💡 Pro Tips

1. **Keep Air Running**: Leave `make watch` running in a dedicated terminal
2. **Multiple Changes**: Air batches rapid changes to avoid excessive restarts  
3. **Build Errors**: Fix build errors quickly - Air shows them clearly
4. **Environment**: Changes to `.env` also trigger restart
5. **Clean Start**: Use Ctrl+C to stop, then `make watch` for clean restart

## 🎉 Benefits

- **⚡ Faster Development**: No manual restart needed
- **🔍 Immediate Feedback**: See changes instantly
- **🚫 No Stale State**: Fresh server state on every change
- **📈 Productivity**: Focus on coding, not restarting
- **🔧 Easy Debugging**: Clear logs and error messages

Happy coding with hot reload! 🔥