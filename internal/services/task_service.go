package services

import (
	"fmt"
	"time"

	"ayuryuj-task/internal/models"
	"ayuryuj-task/internal/repository"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) CreateTask(req models.CreateTaskRequest) (*models.Task, error) {
	if err := s.validateTrigger(req.Trigger); err != nil {
		return nil, fmt.Errorf("invalid trigger: %w", err)
	}

	task := &models.Task{
		ID:        uuid.New(),
		Name:      req.Name,
		Trigger:   req.Trigger,
		Action:    req.Action,
		Status:    models.TaskStatusScheduled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	nextRun, err := s.calculateNextRun(req.Trigger)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate next run: %w", err)
	}
	task.NextRun = nextRun

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (s *TaskService) GetTask(id uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (s *TaskService) ListTasks(filter models.TaskFilter) (*models.TaskListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	tasks, total, err := s.taskRepo.List(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	totalPages := (total + filter.PageSize - 1) / filter.PageSize

	return &models.TaskListResponse{
		Tasks:      tasks,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *TaskService) UpdateTask(id uuid.UUID, req models.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	if req.Name != nil {
		task.Name = *req.Name
	}
	if req.Trigger != nil {
		if err := s.validateTrigger(*req.Trigger); err != nil {
			return nil, fmt.Errorf("invalid trigger: %w", err)
		}
		task.Trigger = *req.Trigger

		nextRun, err := s.calculateNextRun(*req.Trigger)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate next run: %w", err)
		}
		task.NextRun = nextRun
	}
	if req.Action != nil {
		task.Action = *req.Action
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (s *TaskService) DeleteTask(id uuid.UUID) error {
	if err := s.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func (s *TaskService) GetScheduledTasks() ([]models.Task, error) {
	tasks, err := s.taskRepo.GetScheduledTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled tasks: %w", err)
	}
	return tasks, nil
}

func (s *TaskService) UpdateTaskNextRun(taskID uuid.UUID, nextRun *time.Time) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	task.NextRun = nextRun
	task.UpdatedAt = time.Now()

	return s.taskRepo.Update(task)
}

func (s *TaskService) MarkTaskCompleted(taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	if task.Trigger.Type == models.TriggerTypeOneOff {
		task.Status = models.TaskStatusCompleted
		task.NextRun = nil
		task.UpdatedAt = time.Now()
		return s.taskRepo.Update(task)
	}

	return nil
}

func (s *TaskService) validateTrigger(trigger models.Trigger) error {
	switch trigger.Type {
	case models.TriggerTypeOneOff:
		if trigger.DateTime == nil {
			return fmt.Errorf("datetime is required for one-off triggers")
		}
		if trigger.DateTime.Before(time.Now()) {
			return fmt.Errorf("datetime cannot be in the past")
		}
	case models.TriggerTypeCron:
		if trigger.Cron == "" {
			return fmt.Errorf("cron expression is required for cron triggers")
		}
		_, err := cron.ParseStandard(trigger.Cron)
		if err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}
	default:
		return fmt.Errorf("invalid trigger type: %s", trigger.Type)
	}
	return nil
}

func (s *TaskService) calculateNextRun(trigger models.Trigger) (*time.Time, error) {
	switch trigger.Type {
	case models.TriggerTypeOneOff:
		if trigger.DateTime == nil {
			return nil, fmt.Errorf("datetime is required for one-off triggers")
		}
		return trigger.DateTime, nil
	case models.TriggerTypeCron:
		if trigger.Cron == "" {
			return nil, fmt.Errorf("cron expression is required for cron triggers")
		}
		schedule, err := cron.ParseStandard(trigger.Cron)
		if err != nil {
			return nil, fmt.Errorf("invalid cron expression: %w", err)
		}
		nextTime := schedule.Next(time.Now())
		return &nextTime, nil
	default:
		return nil, fmt.Errorf("invalid trigger type: %s", trigger.Type)
	}
}
