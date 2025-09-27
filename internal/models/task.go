package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusScheduled TaskStatus = "scheduled"
	TaskStatusCancelled TaskStatus = "cancelled"
	TaskStatusCompleted TaskStatus = "completed"
)

type TriggerType string

const (
	TriggerTypeOneOff TriggerType = "one-off"
	TriggerTypeCron   TriggerType = "cron"
)

type Trigger struct {
	Type     TriggerType `json:"type" validate:"required,oneof=one-off cron"`
	DateTime *time.Time  `json:"datetime,omitempty"`
	Cron     string      `json:"cron,omitempty"`
}

func (t Trigger) Value() (driver.Value, error) {
	return json.Marshal(t)
}

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

type Action struct {
	Method  string            `json:"method" validate:"required"`
	URL     string            `json:"url" validate:"required,url"`
	Headers map[string]string `json:"headers,omitempty"`
	Payload json.RawMessage   `json:"payload,omitempty"`
}

func (a Action) Value() (driver.Value, error) {
	return json.Marshal(a)
}

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

type CreateTaskRequest struct {
	Name    string  `json:"name" validate:"required"`
	Trigger Trigger `json:"trigger" validate:"required"`
	Action  Action  `json:"action" validate:"required"`
}

type UpdateTaskRequest struct {
	Name    *string  `json:"name,omitempty"`
	Trigger *Trigger `json:"trigger,omitempty"`
	Action  *Action  `json:"action,omitempty"`
}

type TaskListResponse struct {
	Tasks      []Task `json:"tasks"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

type TaskFilter struct {
	Status   TaskStatus `json:"status,omitempty"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}
