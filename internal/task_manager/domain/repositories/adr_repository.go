package repositories

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// ADRRepository defines the contract for persistent storage of ADR (Architecture Decision Record) entities.
type ADRRepository interface {
	// SaveADR persists a new ADR to storage.
	// Returns ErrAlreadyExists if an ADR with the same ID already exists.
	// Returns ErrNotFound if the track doesn't exist.
	SaveADR(ctx context.Context, adr *entities.ADREntity) error

	// GetADR retrieves an ADR by its ID.
	// Returns ErrNotFound if the ADR doesn't exist.
	GetADR(ctx context.Context, id string) (*entities.ADREntity, error)

	// ListADRs returns all ADRs, optionally filtered by track.
	// Returns empty slice if no ADRs match the filters.
	ListADRs(ctx context.Context, trackID *string) ([]*entities.ADREntity, error)

	// UpdateADR updates an existing ADR.
	// Returns ErrNotFound if the ADR doesn't exist.
	UpdateADR(ctx context.Context, adr *entities.ADREntity) error

	// SupersedeADR marks an ADR as superseded by another ADR.
	// Returns ErrNotFound if either ADR doesn't exist.
	SupersedeADR(ctx context.Context, adrID, supersededByID string) error

	// DeprecateADR marks an ADR as deprecated.
	// Returns ErrNotFound if the ADR doesn't exist.
	DeprecateADR(ctx context.Context, adrID string) error

	// GetADRsByTrack returns all ADRs for a specific track.
	// Returns empty slice if the track has no ADRs.
	GetADRsByTrack(ctx context.Context, trackID string) ([]*entities.ADREntity, error)
}
