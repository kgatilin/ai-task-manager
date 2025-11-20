package repositories

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// AggregateRepository defines the contract for aggregate queries and project metadata operations.
type AggregateRepository interface {
	// GetRoadmapWithTracks retrieves a roadmap with all its tracks.
	// The roadmap is returned with Dependencies populated from the database.
	// Returns ErrNotFound if the roadmap doesn't exist.
	GetRoadmapWithTracks(ctx context.Context, roadmapID string) (*entities.RoadmapEntity, error)

	// GetProjectMetadata retrieves a metadata value by key.
	// Returns ErrNotFound if the key doesn't exist.
	GetProjectMetadata(ctx context.Context, key string) (string, error)

	// SetProjectMetadata sets a metadata value by key.
	// Creates or updates the key-value pair.
	SetProjectMetadata(ctx context.Context, key, value string) error

	// GetProjectCode retrieves the project code (e.g., "DW" for darwinflow).
	// Returns "DW" as default if not set.
	GetProjectCode(ctx context.Context) string

	// GetNextSequenceNumber retrieves the next sequence number for an entity type.
	// Entity types: "task", "track", "iter", "adr", "ac"
	GetNextSequenceNumber(ctx context.Context, entityType string) (int, error)
}
