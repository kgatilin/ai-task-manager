package repositories

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// TrackRepository defines the contract for persistent storage of track entities.
type TrackRepository interface {
	// SaveTrack persists a new track to storage.
	// Returns ErrAlreadyExists if a track with the same ID already exists.
	SaveTrack(ctx context.Context, track *entities.TrackEntity) error

	// GetTrack retrieves a track by its ID.
	// Returns ErrNotFound if the track doesn't exist.
	GetTrack(ctx context.Context, id string) (*entities.TrackEntity, error)

	// ListTracks returns all tracks for a roadmap, optionally filtered.
	// Returns empty slice if no tracks match the filters.
	ListTracks(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error)

	// UpdateTrack updates an existing track.
	// Returns ErrNotFound if the track doesn't exist.
	UpdateTrack(ctx context.Context, track *entities.TrackEntity) error

	// DeleteTrack removes a track from storage.
	// Returns ErrNotFound if the track doesn't exist.
	DeleteTrack(ctx context.Context, id string) error

	// AddTrackDependency adds a dependency from trackID to dependsOnID.
	// Returns ErrNotFound if either track doesn't exist.
	// Returns ErrInvalidArgument if it would create a self-dependency.
	// Returns ErrAlreadyExists if the dependency already exists.
	AddTrackDependency(ctx context.Context, trackID, dependsOnID string) error

	// RemoveTrackDependency removes a dependency from trackID to dependsOnID.
	// Returns ErrNotFound if the dependency doesn't exist.
	RemoveTrackDependency(ctx context.Context, trackID, dependsOnID string) error

	// GetTrackDependencies returns the IDs of all tracks that trackID depends on.
	// Returns empty slice if there are no dependencies.
	GetTrackDependencies(ctx context.Context, trackID string) ([]string, error)

	// ValidateNoCycles checks if adding/updating the track would create a circular dependency.
	// Returns ErrInvalidArgument if a cycle is detected.
	ValidateNoCycles(ctx context.Context, trackID string) error

	// GetTrackWithTasks retrieves a track with all its tasks.
	// The track is returned with Dependencies populated from the database.
	// Returns ErrNotFound if the track doesn't exist.
	GetTrackWithTasks(ctx context.Context, trackID string) (*entities.TrackEntity, error)
}
