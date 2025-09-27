package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusScheduled TaskStatus = "scheduled"
	TaskStatusCancelled TaskStatus = "cancelled"
	TaskStatusCompleted TaskStatus = "completed"
)

// TriggerType represents the type of trigger
type TriggerType string

const (
	TriggerTypeOneOff TriggerType = "one-off"
	TriggerTypeCron   TriggerType = "cron"
)

// Trigger represents the scheduling configuration for a task
type Trigger struct {
	Type     TriggerType `json:"type" validate:"required,oneof=one-off cron"`
	DateTime *time.Time  `json:"datetime,omitempty"` // For one-off tasks
	Cron     string      `json:"cron,omitempty"`     // For cron tasks
}

// Value implements driver.Valuer interface for database storage
func (t Trigger) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan implements sql.Scanner interface for database retrieval
func (t *Trigger) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), t)
	}
	return json.Unmarshal(bytes, t)
}

// Action represents the HTTP action to be performed
type Action struct {
	Method  string            `json:"method" validate:"required"`
	URL     string            `json:"url" validate:"required,url"`
	Headers map[string]string `json:"headers,omitempty"`
	Payload json.RawMessage   `json:"payload,omitempty"`
}

// Value implements driver.Valuer interface for database storage
func (a Action) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements sql.Scanner interface for database retrieval
func (a *Action) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), a)
	}
	return json.Unmarshal(bytes, a)
}

// Task represents a scheduled task
type Task struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name" validate:"required"`
	Trigger   Trigger    `json:"trigger" db:"trigger" validate:"required"`
	Action    Action     `json:"action" db:"action" validate:"required"`
	Status    TaskStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	NextRun   *time.Time `json:"next_run,omitempty" db:"next_run"`
}

// CreateTaskRequest represents the request payload for creating a task
type CreateTaskRequest struct {
	Name    string  `json:"name" validate:"required"`
	Trigger Trigger `json:"trigger" validate:"required"`
	Action  Action  `json:"action" validate:"required"`
}

// UpdateTaskRequest represents the request payload for updating a task
type UpdateTaskRequest struct {
	Name    *string  `json:"name,omitempty"`
	Trigger *Trigger `json:"trigger,omitempty"`
	Action  *Action  `json:"action,omitempty"`
}

// TaskListResponse represents the response for listing tasks
type TaskListResponse struct {
	Tasks      []Task `json:"tasks"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// TaskFilter represents filters for task listing
type TaskFilter struct {
	Status   TaskStatus `json:"status,omitempty"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}
