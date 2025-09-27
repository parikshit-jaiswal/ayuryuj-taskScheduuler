package handlers

import (
	"encoding/json"
	"net/http"

	"ayuryuj-task/internal/database"
)

type HealthHandler struct {
	db database.Service
}

func NewHealthHandler(db database.Service) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	health := h.db.Health()

	statusCode := http.StatusOK
	if health["status"] != "up" {
		statusCode = http.StatusServiceUnavailable
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}

func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"uptime":  "up",
		"service": "task-scheduler",
	}

	json.NewEncoder(w).Encode(metrics)
}
