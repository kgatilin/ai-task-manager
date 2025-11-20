package repositories

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// TaskRepository defines the contract for persistent storage of task entities.
type TaskRepository interface {
	// SaveTask persists a new task to storage.
	// Returns ErrAlreadyExists if a task with the same ID already exists.
	// Returns ErrNotFound if the track doesn't exist.
	SaveTask(ctx context.Context, task *entities.TaskEntity) error

	// GetTask retrieves a task by its ID.
	// Returns ErrNotFound if the task doesn't exist.
	GetTask(ctx context.Context, id string) (*entities.TaskEntity, error)

	// ListTasks returns all tasks matching the filters.
	// Returns empty slice if no tasks match the filters.
	ListTasks(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error)

	// UpdateTask updates an existing task.
	// Returns ErrNotFound if the task doesn't exist.
	UpdateTask(ctx context.Context, task *entities.TaskEntity) error

	// DeleteTask removes a task from storage.
	// Returns ErrNotFound if the task doesn't exist.
	DeleteTask(ctx context.Context, id string) error

	// MoveTaskToTrack moves a task from its current track to a new track.
	// Returns ErrNotFound if the task or new track doesn't exist.
	MoveTaskToTrack(ctx context.Context, taskID, newTrackID string) error

	// GetBacklogTasks returns all tasks that are not in any iteration and not done.
	// Returns empty slice if there are no backlog tasks.
	// Ordered by created_at ascending.
	GetBacklogTasks(ctx context.Context) ([]*entities.TaskEntity, error)

	// GetIterationsForTask returns all iterations that contain a specific task.
	// Returns empty slice if the task is not in any iterations.
	// Ordered by iteration number ascending.
	GetIterationsForTask(ctx context.Context, taskID string) ([]*entities.IterationEntity, error)
}
