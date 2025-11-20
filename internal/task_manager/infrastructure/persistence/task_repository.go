package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
)

// Compile-time check that SQLiteTaskRepository implements repositories.TaskRepository
var _ repositories.TaskRepository = (*SQLiteTaskRepository)(nil)

// SQLiteTaskRepository implements repositories.TaskRepository using SQLite as the backend.
type SQLiteTaskRepository struct {
	DB     *sql.DB
	logger logger.Logger
}

// NewSQLiteTaskRepository creates a new SQLite-backed repository.
func NewSQLiteTaskRepository(db *sql.DB, logger logger.Logger) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{
		DB:     db,
		logger: logger,
	}
}

// ============================================================================
// Task Operations
// ============================================================================

// SaveTask persists a new task to storage.
func (r *SQLiteTaskRepository) SaveTask(ctx context.Context, task *entities.TaskEntity) error {
	// Check if task already exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tasks WHERE id = ?", task.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check task existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w: task %s already exists", tmerrors.ErrAlreadyExists, task.ID)
	}

	// Check if track exists
	var trackExists int
	err = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tracks WHERE id = ?", task.TrackID).Scan(&trackExists)
	if err != nil {
		return fmt.Errorf("failed to check track existence: %w", err)
	}
	if trackExists == 0 {
		return fmt.Errorf("%w: track %s not found", tmerrors.ErrNotFound, task.TrackID)
	}

	_, err = r.DB.ExecContext(
		ctx,
		"INSERT INTO tasks (id, track_id, title, description, status, rank, branch, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		task.ID, task.TrackID, task.Title, task.Description, task.Status, task.Rank, task.Branch, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

// GetTask retrieves a task by its ID.
func (r *SQLiteTaskRepository) GetTask(ctx context.Context, id string) (*entities.TaskEntity, error) {
	var task entities.TaskEntity
	var branch sql.NullString

	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id, track_id, title, description, status, rank, branch, created_at, updated_at FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.TrackID, &task.Title, &task.Description, &task.Status, &task.Rank, &branch, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: task %s not found", tmerrors.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to query task: %w", err)
	}

	if branch.Valid {
		task.Branch = branch.String
	}

	return &task, nil
}

// ListTasks returns all tasks matching the filters.
func (r *SQLiteTaskRepository) ListTasks(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error) {
	query := "SELECT id, track_id, title, description, status, rank, branch, created_at, updated_at FROM tasks WHERE 1=1"
	args := []interface{}{}

	// Add track filter if provided
	if filters.TrackID != "" {
		query += " AND track_id = ?"
		args = append(args, filters.TrackID)
	}

	// Add status filter if provided
	if len(filters.Status) > 0 {
		placeholders := ""
		for i := range filters.Status {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, filters.Status[i])
		}
		query += " AND status IN (" + placeholders + ")"
	}

	// Add priority filter if provided
	if len(filters.Priority) > 0 {
		placeholders := ""
		for i := range filters.Priority {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, filters.Priority[i])
		}
		query += " AND rank IN (" + placeholders + ")"
	}

	query += " ORDER BY id"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*entities.TaskEntity
	for rows.Next() {
		var task entities.TaskEntity
		var branch sql.NullString

		err := rows.Scan(&task.ID, &task.TrackID, &task.Title, &task.Description, &task.Status, &task.Rank, &branch, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if branch.Valid {
			task.Branch = branch.String
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// UpdateTask updates an existing task.
func (r *SQLiteTaskRepository) UpdateTask(ctx context.Context, task *entities.TaskEntity) error {
	result, err := r.DB.ExecContext(
		ctx,
		"UPDATE tasks SET track_id = ?, title = ?, description = ?, status = ?, rank = ?, branch = ?, updated_at = ? WHERE id = ?",
		task.TrackID, task.Title, task.Description, task.Status, task.Rank, task.Branch, task.UpdatedAt, task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: task %s not found", tmerrors.ErrNotFound, task.ID)
	}

	return nil
}

// DeleteTask removes a task from storage.
func (r *SQLiteTaskRepository) DeleteTask(ctx context.Context, id string) error {
	result, err := r.DB.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: task %s not found", tmerrors.ErrNotFound, id)
	}

	return nil
}

// MoveTaskToTrack moves a task from its current track to a new track.
func (r *SQLiteTaskRepository) MoveTaskToTrack(ctx context.Context, taskID, newTrackID string) error {
	// Check if task exists
	task, err := r.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	// Check if new track exists
	var trackExists int
	err = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tracks WHERE id = ?", newTrackID).Scan(&trackExists)
	if err != nil {
		return fmt.Errorf("failed to check track existence: %w", err)
	}
	if trackExists == 0 {
		return fmt.Errorf("%w: track %s not found", tmerrors.ErrNotFound, newTrackID)
	}

	// Update task's track
	task.TrackID = newTrackID
	task.UpdatedAt = time.Now().UTC()
	return r.UpdateTask(ctx, task)
}

// GetBacklogTasks returns all tasks that are not in any iteration and not done.
func (r *SQLiteTaskRepository) GetBacklogTasks(ctx context.Context) ([]*entities.TaskEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT t.id, t.track_id, t.title, t.description, t.status, t.rank, t.branch, t.created_at, t.updated_at
		 FROM tasks t
		 LEFT JOIN iteration_tasks it ON t.id = it.task_id
		 WHERE it.task_id IS NULL AND t.status != 'done'
		 ORDER BY t.created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query backlog tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*entities.TaskEntity
	for rows.Next() {
		var task entities.TaskEntity
		var branch sql.NullString

		err := rows.Scan(&task.ID, &task.TrackID, &task.Title, &task.Description, &task.Status, &task.Rank, &branch, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if branch.Valid {
			task.Branch = branch.String
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// GetIterationsForTask returns all iterations that contain a specific task.
func (r *SQLiteTaskRepository) GetIterationsForTask(ctx context.Context, taskID string) ([]*entities.IterationEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT i.number, i.name, i.goal, i.status, i.rank, i.deliverable, i.started_at, i.completed_at, i.created_at, i.updated_at
		 FROM iterations i
		 JOIN iteration_tasks it ON i.number = it.iteration_number
		 WHERE it.task_id = ?
		 ORDER BY i.number ASC`,
		taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query iterations for task: %w", err)
	}
	defer rows.Close()

	var iterations []*entities.IterationEntity
	for rows.Next() {
		var iteration entities.IterationEntity
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(&iteration.Number, &iteration.Name, &iteration.Goal, &iteration.Status, &iteration.Rank, &iteration.Deliverable, &startedAt, &completedAt, &iteration.CreatedAt, &iteration.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan iteration: %w", err)
		}

		if startedAt.Valid {
			iteration.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			iteration.CompletedAt = &completedAt.Time
		}

		// Load task IDs for each iteration
		taskIDs, err := r.getIterationTaskIDs(ctx, iteration.Number)
		if err != nil {
			return nil, fmt.Errorf("failed to load iteration tasks: %w", err)
		}
		iteration.TaskIDs = taskIDs

		iterations = append(iterations, &iteration)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating iterations: %w", err)
	}

	return iterations, nil
}

// ============================================================================
// Helper Methods
// ============================================================================

// getIterationTaskIDs retrieves all task IDs for an iteration.
func (r *SQLiteTaskRepository) getIterationTaskIDs(ctx context.Context, iterationNum int) ([]string, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT task_id FROM iteration_tasks WHERE iteration_number = ? ORDER BY task_id",
		iterationNum,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query iteration tasks: %w", err)
	}
	defer rows.Close()

	var taskIDs []string
	for rows.Next() {
		var taskID string
		if err := rows.Scan(&taskID); err != nil {
			return nil, fmt.Errorf("failed to scan task ID: %w", err)
		}
		taskIDs = append(taskIDs, taskID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task IDs: %w", err)
	}

	return taskIDs, nil
}
