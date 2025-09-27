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
		checkInterval:     10 * time.Second,
	}
}

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

func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

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

func (s *Scheduler) processPendingTasks() {
	tasks, err := s.taskService.GetScheduledTasks()
	if err != nil {
		log.Printf("Error getting scheduled tasks: %v", err)
		return
	}

	for _, task := range tasks {
		go s.executeTask(task)
	}
}

func (s *Scheduler) executeTask(task models.Task) {
	log.Printf("Executing task: %s (ID: %s)", task.Name, task.ID)

	result := s.httpExecutor.ExecuteTask(task)

	if err := s.taskResultService.CreateTaskResult(result); err != nil {
		log.Printf("Error saving task result for task %s: %v", task.ID, err)
	}

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

func (s *Scheduler) updateTaskAfterExecution(task models.Task) error {
	switch task.Trigger.Type {
	case models.TriggerTypeOneOff:
		return s.taskService.MarkTaskCompleted(task.ID)
	case models.TriggerTypeCron:
		return s.updateCronTaskNextRun(task)
	}
	return nil
}

func (s *Scheduler) updateCronTaskNextRun(task models.Task) error {
	if task.Trigger.Type != models.TriggerTypeCron {
		return nil
	}

	schedule, err := cron.ParseStandard(task.Trigger.Cron)
	if err != nil {
		return fmt.Errorf("failed to parse cron expression: %w", err)
	}

	nextRun := schedule.Next(time.Now())
	return s.taskService.UpdateTaskNextRun(task.ID, &nextRun)
}
