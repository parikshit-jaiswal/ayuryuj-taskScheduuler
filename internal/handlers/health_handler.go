package handlers

import (
	"encoding/json"
	"net/http"

	"ayuryuj-task/internal/database"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db database.Service
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db database.Service) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Health handles GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	health := h.db.Health()

	statusCode := http.StatusOK
	if health["status"] != "up" {
		statusCode = http.StatusServiceUnavailable
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}

// Metrics handles GET /metrics (basic implementation)
func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	// Basic metrics response
	metrics := map[string]interface{}{
		"uptime":  "up",
		"service": "task-scheduler",
	}

	json.NewEncoder(w).Encode(metrics)
}
