package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"ayuryuj-task/internal/database"
	"ayuryuj-task/internal/repository"
	"ayuryuj-task/internal/scheduler"
	"ayuryuj-task/internal/services"
)

type Server struct {
	port              int
	db                database.Service
	taskService       *services.TaskService
	taskResultService *services.TaskResultService
	httpExecutor      *services.HTTPExecutorService
	scheduler         *scheduler.Scheduler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db := database.New()

	if err := db.RunMigrations(); err != nil {
		log.Printf("Failed to run migrations: %v", err)
	}

	taskRepo := repository.NewTaskRepository(db.GetDB())
	taskResultRepo := repository.NewTaskResultRepository(db.GetDB())

	taskService := services.NewTaskService(taskRepo)
	taskResultService := services.NewTaskResultService(taskResultRepo)
	httpExecutor := services.NewHTTPExecutorService()

	taskScheduler := scheduler.NewScheduler(taskService, taskResultService, httpExecutor)

	newServer := &Server{
		port:              port,
		db:                db,
		taskService:       taskService,
		taskResultService: taskResultService,
		httpExecutor:      httpExecutor,
		scheduler:         taskScheduler,
	}

	taskScheduler.Start()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) Shutdown() {
	if s.scheduler != nil {
		s.scheduler.Stop()
	}
	if s.db != nil {
		s.db.Close()
	}
}
