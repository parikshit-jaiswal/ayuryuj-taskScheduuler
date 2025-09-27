package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"ayuryuj-task/internal/models"
	"ayuryuj-task/internal/services"

	"github.com/google/uuid"
)

type TaskResultHandler struct {
	taskResultService *services.TaskResultService
}

func NewTaskResultHandler(taskResultService *services.TaskResultService) *TaskResultHandler {
	return &TaskResultHandler{
		taskResultService: taskResultService,
	}
}

func (h *TaskResultHandler) ListResults(w http.ResponseWriter, r *http.Request) {
	var filter models.TaskResultFilter

	if taskID := r.URL.Query().Get("task_id"); taskID != "" {
		if id, err := uuid.Parse(taskID); err == nil {
			filter.TaskID = &id
		}
	}

	if success := r.URL.Query().Get("success"); success != "" {
		if s, err := strconv.ParseBool(success); err == nil {
			filter.Success = &s
		}
	}

	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			filter.DateFrom = &t
		}
	}

	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
			filter.DateTo = &t
		}
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filter.Page = p
		}
	}

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			filter.PageSize = ps
		}
	}

	response, err := h.taskResultService.ListTaskResults(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
