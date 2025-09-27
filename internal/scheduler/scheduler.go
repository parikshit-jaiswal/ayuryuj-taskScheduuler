package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"ayuryuj-task/internal/models"
	"ayuryuj-task/internal/services"

	"github.com/robfig/cron/v3"
)

// Scheduler handles task scheduling and execution
type Scheduler struct {
	taskService       *services.TaskService
	taskResultService *services.TaskResultService
	httpExecutor      *services.HTTPExecutorService
	ticker            *time.Ticker
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	running           bool
	mu                sync.RWMutex
	checkInterval     time.Duration
}

// NewScheduler creates a new scheduler
func NewScheduler(
	taskService *services.TaskService,
	taskResultService *services.TaskResultService,
	httpExecutor *services.HTTPExecutorService,
) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		taskService:       taskService,
		taskResultService: taskResultService,
		httpExecutor:      httpExecutor,
		ctx:               ctx,
		cancel:            cancel,
		checkInterval:     10 * time.Second, // Check for scheduled tasks every 10 seconds
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}

	s.running = true
	s.ticker = time.NewTicker(s.checkInterval)

	s.wg.Add(1)
	go s.run()

	log.Println("Task scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	s.cancel()

	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.wg.Wait()
	log.Println("Task scheduler stopped")
}

// IsRunning returns whether the scheduler is running
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// run is the main scheduler loop
func (s *Scheduler) run() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ticker.C:
			s.processPendingTasks()
		}
	}
}

// processPendingTasks checks for and executes scheduled tasks
func (s *Scheduler) processPendingTasks() {
	tasks, err := s.taskService.GetScheduledTasks()
	if err != nil {
		log.Printf("Error getting scheduled tasks: %v", err)
		return
	}

	for _, task := range tasks {
		// Execute task in a separate goroutine to avoid blocking
		go s.executeTask(task)
	}
}

// executeTask executes a single task
func (s *Scheduler) executeTask(task models.Task) {
	log.Printf("Executing task: %s (ID: %s)", task.Name, task.ID)

	// Execute the HTTP request
	result := s.httpExecutor.ExecuteTask(task)

	// Save the result
	if err := s.taskResultService.CreateTaskResult(result); err != nil {
		log.Printf("Error saving task result for task %s: %v", task.ID, err)
	}

	// Update task based on trigger type
	if err := s.updateTaskAfterExecution(task); err != nil {
		log.Printf("Error updating task after execution %s: %v", task.ID, err)
	}

	if result.Success {
		log.Printf("Task %s executed successfully (Status: %d, Duration: %dms)",
			task.ID, result.StatusCode, result.DurationMs)
	} else {
		log.Printf("Task %s execution failed (Status: %d, Duration: %dms, Error: %s)",
			task.ID, result.StatusCode, result.DurationMs,
			func() string {
				if result.ErrorMessage != nil {
					return *result.ErrorMessage
				}
				return "No error message"
			}())
	}
}

// updateTaskAfterExecution updates the task after execution based on its trigger type
func (s *Scheduler) updateTaskAfterExecution(task models.Task) error {
	switch task.Trigger.Type {
	case models.TriggerTypeOneOff:
		// Mark one-off tasks as completed
		return s.taskService.MarkTaskCompleted(task.ID)
	case models.TriggerTypeCron:
		// Calculate next run time for cron tasks
		return s.updateCronTaskNextRun(task)
	}
	return nil
}

// updateCronTaskNextRun calculates and updates the next run time for cron tasks
func (s *Scheduler) updateCronTaskNextRun(task models.Task) error {
	if task.Trigger.Type != models.TriggerTypeCron {
		return nil
	}

	// Parse the cron expression and calculate next run
	schedule, err := cron.ParseStandard(task.Trigger.Cron)
	if err != nil {
		return fmt.Errorf("failed to parse cron expression: %w", err)
	}

	nextRun := schedule.Next(time.Now())
	return s.taskService.UpdateTaskNextRun(task.ID, &nextRun)
}
