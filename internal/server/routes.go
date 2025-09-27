package server

import (
	"net/http"

	"ayuryuj-task/internal/handlers"
	"ayuryuj-task/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	// Create handlers
	taskHandler := handlers.NewTaskHandler(s.taskService, s.taskResultService)
	taskResultHandler := handlers.NewTaskResultHandler(s.taskResultService)
	healthHandler := handlers.NewHealthHandler(s.db)

	// Create router
	mux := http.NewServeMux()

	// Health endpoints
	mux.HandleFunc("GET /health", healthHandler.Health)
	mux.HandleFunc("GET /metrics", healthHandler.Metrics)

	// Task endpoints
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.ListTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)
	mux.HandleFunc("GET /tasks/{id}/results", taskHandler.GetTaskResults)

	// Task results endpoints
	mux.HandleFunc("GET /results", taskResultHandler.ListResults)

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.JSONMiddleware(handler)

	return handler
}
