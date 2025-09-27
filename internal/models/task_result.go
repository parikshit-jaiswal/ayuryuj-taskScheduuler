package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TaskResult represents the result of a task execution
type TaskResult struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	TaskID          uuid.UUID       `json:"task_id" db:"task_id"`
	RunAt           time.Time       `json:"run_at" db:"run_at"`
	StatusCode      int             `json:"status_code" db:"status_code"`
	Success         bool            `json:"success" db:"success"`
	ResponseHeaders json.RawMessage `json:"response_headers" db:"response_headers"`
	ResponseBody    string          `json:"response_body" db:"response_body"`
	ErrorMessage    *string         `json:"error_message,omitempty" db:"error_message"`
	DurationMs      int64           `json:"duration_ms" db:"duration_ms"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}

// TaskResultListResponse represents the response for listing task results
type TaskResultListResponse struct {
	Results    []TaskResult `json:"results"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

// TaskResultFilter represents filters for task result listing
type TaskResultFilter struct {
	TaskID   *uuid.UUID `json:"task_id,omitempty"`
	Success  *bool      `json:"success,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// HTTPHeaders is a custom type for HTTP headers
type HTTPHeaders map[string]string

// Value implements driver.Valuer interface for database storage
func (h HTTPHeaders) Value() (driver.Value, error) {
	if h == nil {
		return nil, nil
	}
	return json.Marshal(h)
}

// Scan implements sql.Scanner interface for database retrieval
func (h *HTTPHeaders) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), h)
	}
	return json.Unmarshal(bytes, h)
}
