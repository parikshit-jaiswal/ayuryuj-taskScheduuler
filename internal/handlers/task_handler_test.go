package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ayuryuj-task/internal/handlers"
	"ayuryuj-task/internal/models"
	"ayuryuj-task/internal/services"

	"github.com/google/uuid"
)

// Mock repository for testing
type mockTaskRepository struct {
	tasks map[uuid.UUID]*models.Task
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks: make(map[uuid.UUID]*models.Task),
	}
}

func (m *mockTaskRepository) Create(task *models.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) GetByID(id uuid.UUID) (*models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, nil
	}
	return task, nil
}

func (m *mockTaskRepository) List(filter models.TaskFilter) ([]models.Task, int, error) {
	tasks := make([]models.Task, 0)
	for _, task := range m.tasks {
		if filter.Status == "" || task.Status == filter.Status {
			tasks = append(tasks, *task)
		}
	}
	return tasks, len(tasks), nil
}

func (m *mockTaskRepository) Update(task *models.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) Delete(id uuid.UUID) error {
	if task, exists := m.tasks[id]; exists {
		task.Status = models.TaskStatusCancelled
		task.UpdatedAt = time.Now()
		m.tasks[id] = task
	}
	return nil
}

func (m *mockTaskRepository) GetScheduledTasks() ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	now := time.Now()
	for _, task := range m.tasks {
		if task.Status == models.TaskStatusScheduled &&
			(task.NextRun == nil || task.NextRun.Before(now)) {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func TestCreateTask(t *testing.T) {
	// Setup
	mockRepo := newMockTaskRepository()
	taskService := services.NewTaskService(mockRepo)
	mockResultRepo := &mockTaskResultRepository{results: make(map[uuid.UUID]*models.TaskResult)}
	resultService := services.NewTaskResultService(mockResultRepo)
	handler := handlers.NewTaskHandler(taskService, resultService)

	// Test data
	futureTime := time.Now().Add(1 * time.Hour)
	reqBody := models.CreateTaskRequest{
		Name: "Test Task",
		Trigger: models.Trigger{
			Type:     models.TriggerTypeOneOff,
			DateTime: &futureTime,
		},
		Action: models.Action{
			Method: "GET",
			URL:    "https://httpbin.org/get",
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Execute
	handler.CreateTask(rr, req)

	// Verify
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var response models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Name != reqBody.Name {
		t.Errorf("Expected name %s, got %s", reqBody.Name, response.Name)
	}
}

// Mock task result repository
type mockTaskResultRepository struct {
	results map[uuid.UUID]*models.TaskResult
}

func (m *mockTaskResultRepository) Create(result *models.TaskResult) error {
	m.results[result.ID] = result
	return nil
}

func (m *mockTaskResultRepository) GetByID(id uuid.UUID) (*models.TaskResult, error) {
	result, exists := m.results[id]
	if !exists {
		return nil, nil
	}
	return result, nil
}

func (m *mockTaskResultRepository) ListByTaskID(taskID uuid.UUID, filter models.TaskResultFilter) ([]models.TaskResult, int, error) {
	results := make([]models.TaskResult, 0)
	for _, result := range m.results {
		if result.TaskID == taskID {
			results = append(results, *result)
		}
	}
	return results, len(results), nil
}

func (m *mockTaskResultRepository) List(filter models.TaskResultFilter) ([]models.TaskResult, int, error) {
	results := make([]models.TaskResult, 0)
	for _, result := range m.results {
		if filter.TaskID == nil || result.TaskID == *filter.TaskID {
			results = append(results, *result)
		}
	}
	return results, len(results), nil
}
