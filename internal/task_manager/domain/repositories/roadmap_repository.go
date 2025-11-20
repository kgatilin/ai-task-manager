package repositories

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// RoadmapRepository defines the contract for persistent storage of roadmap entities.
type RoadmapRepository interface {
	// SaveRoadmap persists a new roadmap to storage.
	// Returns ErrAlreadyExists if a roadmap with the same ID already exists.
	SaveRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error

	// GetRoadmap retrieves a roadmap by its ID.
	// Returns ErrNotFound if the roadmap doesn't exist.
	GetRoadmap(ctx context.Context, id string) (*entities.RoadmapEntity, error)

	// GetActiveRoadmap retrieves the most recently created roadmap.
	// Returns ErrNotFound if no roadmaps exist.
	GetActiveRoadmap(ctx context.Context) (*entities.RoadmapEntity, error)

	// UpdateRoadmap updates an existing roadmap.
	// Returns ErrNotFound if the roadmap doesn't exist.
	UpdateRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error
}
