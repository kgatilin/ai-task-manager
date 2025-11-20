package mocks

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// MockTrackRepository is a mock implementation of repositories.TrackRepository for testing.
type MockTrackRepository struct {
	// In-memory storage for testing
	tracks map[string]*entities.TrackEntity

	// SaveTrackFunc is called by SaveTrack. If nil, uses default implementation.
	SaveTrackFunc func(ctx context.Context, track *entities.TrackEntity) error

	// GetTrackFunc is called by GetTrack. If nil, returns nil, nil.
	GetTrackFunc func(ctx context.Context, id string) (*entities.TrackEntity, error)

	// ListTracksFunc is called by ListTracks. If nil, returns empty slice, nil.
	ListTracksFunc func(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error)

	// UpdateTrackFunc is called by UpdateTrack. If nil, returns nil.
	UpdateTrackFunc func(ctx context.Context, track *entities.TrackEntity) error

	// DeleteTrackFunc is called by DeleteTrack. If nil, returns nil.
	DeleteTrackFunc func(ctx context.Context, id string) error

	// AddTrackDependencyFunc is called by AddTrackDependency. If nil, returns nil.
	AddTrackDependencyFunc func(ctx context.Context, trackID, dependsOnID string) error

	// RemoveTrackDependencyFunc is called by RemoveTrackDependency. If nil, returns nil.
	RemoveTrackDependencyFunc func(ctx context.Context, trackID, dependsOnID string) error

	// GetTrackDependenciesFunc is called by GetTrackDependencies. If nil, returns empty slice, nil.
	GetTrackDependenciesFunc func(ctx context.Context, trackID string) ([]string, error)

	// ValidateNoCyclesFunc is called by ValidateNoCycles. If nil, returns nil.
	ValidateNoCyclesFunc func(ctx context.Context, trackID string) error

	// GetTrackWithTasksFunc is called by GetTrackWithTasks. If nil, returns nil, nil.
	GetTrackWithTasksFunc func(ctx context.Context, trackID string) (*entities.TrackEntity, error)
}

// NewMockTrackRepository creates a new mock track repository with in-memory storage
func NewMockTrackRepository() *MockTrackRepository {
	return &MockTrackRepository{
		tracks: make(map[string]*entities.TrackEntity),
	}
}

// SaveTrack implements repositories.TrackRepository.
func (m *MockTrackRepository) SaveTrack(ctx context.Context, track *entities.TrackEntity) error {
	if m.SaveTrackFunc != nil {
		return m.SaveTrackFunc(ctx, track)
	}
	// Default implementation: store in memory
	m.tracks[track.ID] = track
	return nil
}

// GetTrack implements repositories.TrackRepository.
func (m *MockTrackRepository) GetTrack(ctx context.Context, id string) (*entities.TrackEntity, error) {
	if m.GetTrackFunc != nil {
		return m.GetTrackFunc(ctx, id)
	}
	return nil, nil
}

// ListTracks implements repositories.TrackRepository.
func (m *MockTrackRepository) ListTracks(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
	if m.ListTracksFunc != nil {
		return m.ListTracksFunc(ctx, roadmapID, filters)
	}
	// Default implementation: return all tracks for the roadmap
	var result []*entities.TrackEntity
	for _, track := range m.tracks {
		if track.RoadmapID == roadmapID {
			result = append(result, track)
		}
	}
	return result, nil
}

// UpdateTrack implements repositories.TrackRepository.
func (m *MockTrackRepository) UpdateTrack(ctx context.Context, track *entities.TrackEntity) error {
	if m.UpdateTrackFunc != nil {
		return m.UpdateTrackFunc(ctx, track)
	}
	return nil
}

// DeleteTrack implements repositories.TrackRepository.
func (m *MockTrackRepository) DeleteTrack(ctx context.Context, id string) error {
	if m.DeleteTrackFunc != nil {
		return m.DeleteTrackFunc(ctx, id)
	}
	return nil
}

// AddTrackDependency implements repositories.TrackRepository.
func (m *MockTrackRepository) AddTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	if m.AddTrackDependencyFunc != nil {
		return m.AddTrackDependencyFunc(ctx, trackID, dependsOnID)
	}
	return nil
}

// RemoveTrackDependency implements repositories.TrackRepository.
func (m *MockTrackRepository) RemoveTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	if m.RemoveTrackDependencyFunc != nil {
		return m.RemoveTrackDependencyFunc(ctx, trackID, dependsOnID)
	}
	return nil
}

// GetTrackDependencies implements repositories.TrackRepository.
func (m *MockTrackRepository) GetTrackDependencies(ctx context.Context, trackID string) ([]string, error) {
	if m.GetTrackDependenciesFunc != nil {
		return m.GetTrackDependenciesFunc(ctx, trackID)
	}
	return []string{}, nil
}

// ValidateNoCycles implements repositories.TrackRepository.
func (m *MockTrackRepository) ValidateNoCycles(ctx context.Context, trackID string) error {
	if m.ValidateNoCyclesFunc != nil {
		return m.ValidateNoCyclesFunc(ctx, trackID)
	}
	return nil
}

// GetTrackWithTasks implements repositories.TrackRepository.
func (m *MockTrackRepository) GetTrackWithTasks(ctx context.Context, trackID string) (*entities.TrackEntity, error) {
	if m.GetTrackWithTasksFunc != nil {
		return m.GetTrackWithTasksFunc(ctx, trackID)
	}
	return nil, nil
}

// Reset clears all configured behavior.
func (m *MockTrackRepository) Reset() {
	m.SaveTrackFunc = nil
	m.GetTrackFunc = nil
	m.ListTracksFunc = nil
	m.UpdateTrackFunc = nil
	m.DeleteTrackFunc = nil
	m.AddTrackDependencyFunc = nil
	m.RemoveTrackDependencyFunc = nil
	m.GetTrackDependenciesFunc = nil
	m.ValidateNoCyclesFunc = nil
	m.GetTrackWithTasksFunc = nil
}

// WithError configures the mock to return the specified error for all methods.
func (m *MockTrackRepository) WithError(err error) *MockTrackRepository {
	m.SaveTrackFunc = func(ctx context.Context, track *entities.TrackEntity) error { return err }
	m.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) { return nil, err }
	m.ListTracksFunc = func(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
		return nil, err
	}
	m.UpdateTrackFunc = func(ctx context.Context, track *entities.TrackEntity) error { return err }
	m.DeleteTrackFunc = func(ctx context.Context, id string) error { return err }
	m.AddTrackDependencyFunc = func(ctx context.Context, trackID, dependsOnID string) error { return err }
	m.RemoveTrackDependencyFunc = func(ctx context.Context, trackID, dependsOnID string) error { return err }
	m.GetTrackDependenciesFunc = func(ctx context.Context, trackID string) ([]string, error) { return nil, err }
	m.ValidateNoCyclesFunc = func(ctx context.Context, trackID string) error { return err }
	m.GetTrackWithTasksFunc = func(ctx context.Context, trackID string) (*entities.TrackEntity, error) { return nil, err }
	return m
}
