package services

import (
	"fmt"
	"time"

	"ayuryuj-task/internal/models"
	"ayuryuj-task/internal/repository"

	"github.com/google/uuid"
)

// TaskResultService handles business logic for task results
type TaskResultService struct {
	taskResultRepo *repository.TaskResultRepository
}

// NewTaskResultService creates a new task result service
func NewTaskResultService(taskResultRepo *repository.TaskResultRepository) *TaskResultService {
	return &TaskResultService{
		taskResultRepo: taskResultRepo,
	}
}

// CreateTaskResult creates a new task result
func (s *TaskResultService) CreateTaskResult(result *models.TaskResult) error {
	result.ID = uuid.New()
	result.CreatedAt = time.Now()

	if err := s.taskResultRepo.Create(result); err != nil {
		return fmt.Errorf("failed to create task result: %w", err)
	}

	return nil
}

// GetTaskResult retrieves a task result by ID
func (s *TaskResultService) GetTaskResult(id uuid.UUID) (*models.TaskResult, error) {
	result, err := s.taskResultRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task result: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("task result not found")
	}
	return result, nil
}

// ListTaskResults lists all task results with filtering and pagination
func (s *TaskResultService) ListTaskResults(filter models.TaskResultFilter) (*models.TaskResultListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	results, total, err := s.taskResultRepo.List(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list task results: %w", err)
	}

	totalPages := (total + filter.PageSize - 1) / filter.PageSize

	return &models.TaskResultListResponse{
		Results:    results,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ListTaskResultsByTaskID lists task results for a specific task
func (s *TaskResultService) ListTaskResultsByTaskID(taskID uuid.UUID, filter models.TaskResultFilter) (*models.TaskResultListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	results, total, err := s.taskResultRepo.ListByTaskID(taskID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list task results: %w", err)
	}

	totalPages := (total + filter.PageSize - 1) / filter.PageSize

	return &models.TaskResultListResponse{
		Results:    results,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}
