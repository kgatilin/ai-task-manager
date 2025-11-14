package repositories

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// AcceptanceCriteriaRepository defines the contract for persistent storage of acceptance criteria entities.
type AcceptanceCriteriaRepository interface {
	// SaveAC persists a new acceptance criterion to storage.
	// Returns ErrAlreadyExists if an AC with the same ID already exists.
	// Returns ErrNotFound if the task doesn't exist.
	SaveAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error

	// GetAC retrieves an acceptance criterion by its ID.
	// Returns ErrNotFound if the AC doesn't exist.
	GetAC(ctx context.Context, id string) (*entities.AcceptanceCriteriaEntity, error)

	// ListAC returns all acceptance criteria for a task.
	// Returns empty slice if the task has no ACs.
	ListAC(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error)

	// UpdateAC updates an existing acceptance criterion.
	// Returns ErrNotFound if the AC doesn't exist.
	UpdateAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error

	// DeleteAC removes an acceptance criterion from storage.
	// Returns ErrNotFound if the AC doesn't exist.
	DeleteAC(ctx context.Context, id string) error

	// ListACByTask is an alias for ListAC for consistency with other repositories.
	ListACByTask(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error)

	// ListACByIteration returns all acceptance criteria for all tasks in an iteration.
	// Returns empty slice if the iteration has no ACs.
	ListACByIteration(ctx context.Context, iterationNum int) ([]*entities.AcceptanceCriteriaEntity, error)

	// ListFailedAC returns all acceptance criteria with status "failed".
	// Supports optional filtering by iteration, track, or task.
	// Returns empty slice if no failed ACs match the filters.
	ListFailedAC(ctx context.Context, filters entities.ACFilters) ([]*entities.AcceptanceCriteriaEntity, error)
}
