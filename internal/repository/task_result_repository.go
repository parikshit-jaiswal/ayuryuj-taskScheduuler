package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"ayuryuj-task/internal/models"

	"github.com/google/uuid"
)

// TaskResultRepository handles database operations for task results
type TaskResultRepository struct {
	db *sql.DB
}

// NewTaskResultRepository creates a new task result repository
func NewTaskResultRepository(db *sql.DB) *TaskResultRepository {
	return &TaskResultRepository{db: db}
}

// Create creates a new task result
func (r *TaskResultRepository) Create(result *models.TaskResult) error {
	query := `
		INSERT INTO task_results (id, task_id, run_at, status_code, success, response_headers, response_body, error_message, duration_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		result.ID,
		result.TaskID,
		result.RunAt,
		result.StatusCode,
		result.Success,
		result.ResponseHeaders,
		result.ResponseBody,
		result.ErrorMessage,
		result.DurationMs,
		result.CreatedAt,
	)

	return err
}

// GetByID retrieves a task result by its ID
func (r *TaskResultRepository) GetByID(id uuid.UUID) (*models.TaskResult, error) {
	query := `
		SELECT id, task_id, run_at, status_code, success, response_headers, response_body, error_message, duration_ms, created_at
		FROM task_results
		WHERE id = $1
	`

	result := &models.TaskResult{}
	err := r.db.QueryRow(query, id).Scan(
		&result.ID,
		&result.TaskID,
		&result.RunAt,
		&result.StatusCode,
		&result.Success,
		&result.ResponseHeaders,
		&result.ResponseBody,
		&result.ErrorMessage,
		&result.DurationMs,
		&result.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

// ListByTaskID retrieves task results for a specific task with pagination
func (r *TaskResultRepository) ListByTaskID(taskID uuid.UUID, filter models.TaskResultFilter) ([]models.TaskResult, int, error) {
	// Build the WHERE clause
	whereClauses := []string{"task_id = $1"}
	args := []interface{}{taskID}
	argIndex := 2

	if filter.Success != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("success = $%d", argIndex))
		args = append(args, *filter.Success)
		argIndex++
	}

	if filter.DateFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("run_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("run_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}

	whereClause := "WHERE " + strings.Join(whereClauses, " AND ")

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM task_results %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (filter.Page - 1) * filter.PageSize

	// Get results with pagination
	query := fmt.Sprintf(`
		SELECT id, task_id, run_at, status_code, success, response_headers, response_body, error_message, duration_ms, created_at
		FROM task_results %s
		ORDER BY run_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := make([]models.TaskResult, 0)
	for rows.Next() {
		result := models.TaskResult{}
		err := rows.Scan(
			&result.ID,
			&result.TaskID,
			&result.RunAt,
			&result.StatusCode,
			&result.Success,
			&result.ResponseHeaders,
			&result.ResponseBody,
			&result.ErrorMessage,
			&result.DurationMs,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}

	return results, total, rows.Err()
}

// List retrieves all task results with filtering and pagination
func (r *TaskResultRepository) List(filter models.TaskResultFilter) ([]models.TaskResult, int, error) {
	// Build the WHERE clause
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.TaskID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("task_id = $%d", argIndex))
		args = append(args, *filter.TaskID)
		argIndex++
	}

	if filter.Success != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("success = $%d", argIndex))
		args = append(args, *filter.Success)
		argIndex++
	}

	if filter.DateFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("run_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("run_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM task_results %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (filter.Page - 1) * filter.PageSize

	// Get results with pagination
	query := fmt.Sprintf(`
		SELECT id, task_id, run_at, status_code, success, response_headers, response_body, error_message, duration_ms, created_at
		FROM task_results %s
		ORDER BY run_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := make([]models.TaskResult, 0)
	for rows.Next() {
		result := models.TaskResult{}
		err := rows.Scan(
			&result.ID,
			&result.TaskID,
			&result.RunAt,
			&result.StatusCode,
			&result.Success,
			&result.ResponseHeaders,
			&result.ResponseBody,
			&result.ErrorMessage,
			&result.DurationMs,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}

	return results, total, rows.Err()
}
