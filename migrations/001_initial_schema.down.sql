-- Drop triggers
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_task_results_created_at;
DROP INDEX IF EXISTS idx_task_results_success;
DROP INDEX IF EXISTS idx_task_results_run_at;
DROP INDEX IF EXISTS idx_task_results_task_id;

DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_next_run;
DROP INDEX IF EXISTS idx_tasks_status;

-- Drop tables
DROP TABLE IF EXISTS task_results;
DROP TABLE IF EXISTS tasks;