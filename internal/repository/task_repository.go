package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"ayuryuj-task/internal/models"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (id, name, trigger, action, status, created_at, updated_at, next_run)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		task.ID,
		task.Name,
		task.Trigger,
		task.Action,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
		task.NextRun,
	)

	return err
}

func (r *TaskRepository) GetByID(id uuid.UUID) (*models.Task, error) {
	query := `
		SELECT id, name, trigger, action, status, created_at, updated_at, next_run
		FROM tasks
		WHERE id = $1
	`

	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Name,
		&task.Trigger,
		&task.Action,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.NextRun,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return task, err
}

func (r *TaskRepository) List(filter models.TaskFilter) ([]models.Task, int, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
		SELECT id, name, trigger, action, status, created_at, updated_at, next_run
		FROM tasks %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Trigger,
			&task.Action,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.NextRun,
		)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, rows.Err()
}

func (r *TaskRepository) Update(task *models.Task) error {
	query := `
		UPDATE tasks
		SET name = $1, trigger = $2, action = $3, status = $4, updated_at = $5, next_run = $6
		WHERE id = $7
	`

	result, err := r.db.Exec(query,
		task.Name,
		task.Trigger,
		task.Action,
		task.Status,
		task.UpdatedAt,
		task.NextRun,
		task.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE tasks
		SET status = 'cancelled', updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) GetScheduledTasks() ([]models.Task, error) {
	query := `
		SELECT id, name, trigger, action, status, created_at, updated_at, next_run
		FROM tasks
		WHERE status = 'scheduled' 
		AND (next_run IS NULL OR next_run <= $1)
		ORDER BY next_run ASC NULLS FIRST
	`

	rows, err := r.db.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Trigger,
			&task.Action,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.NextRun,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}
