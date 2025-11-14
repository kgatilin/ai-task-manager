package mocks

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// MockAggregateRepository is a mock implementation of repositories.AggregateRepository for testing.
type MockAggregateRepository struct {
	// GetRoadmapWithTracksFunc is called by GetRoadmapWithTracks. If nil, returns nil, nil.
	GetRoadmapWithTracksFunc func(ctx context.Context, roadmapID string) (*entities.RoadmapEntity, error)

	// GetProjectMetadataFunc is called by GetProjectMetadata. If nil, returns "", nil.
	GetProjectMetadataFunc func(ctx context.Context, key string) (string, error)

	// SetProjectMetadataFunc is called by SetProjectMetadata. If nil, returns nil.
	SetProjectMetadataFunc func(ctx context.Context, key, value string) error

	// GetProjectCodeFunc is called by GetProjectCode. If nil, returns "TM".
	GetProjectCodeFunc func(ctx context.Context) string

	// GetNextSequenceNumberFunc is called by GetNextSequenceNumber. If nil, returns 1, nil.
	GetNextSequenceNumberFunc func(ctx context.Context, entityType string) (int, error)
}

// GetRoadmapWithTracks implements repositories.AggregateRepository.
func (m *MockAggregateRepository) GetRoadmapWithTracks(ctx context.Context, roadmapID string) (*entities.RoadmapEntity, error) {
	if m.GetRoadmapWithTracksFunc != nil {
		return m.GetRoadmapWithTracksFunc(ctx, roadmapID)
	}
	return nil, nil
}

// GetProjectMetadata implements repositories.AggregateRepository.
func (m *MockAggregateRepository) GetProjectMetadata(ctx context.Context, key string) (string, error) {
	if m.GetProjectMetadataFunc != nil {
		return m.GetProjectMetadataFunc(ctx, key)
	}
	return "", nil
}

// SetProjectMetadata implements repositories.AggregateRepository.
func (m *MockAggregateRepository) SetProjectMetadata(ctx context.Context, key, value string) error {
	if m.SetProjectMetadataFunc != nil {
		return m.SetProjectMetadataFunc(ctx, key, value)
	}
	return nil
}

// GetProjectCode implements repositories.AggregateRepository.
func (m *MockAggregateRepository) GetProjectCode(ctx context.Context) string {
	if m.GetProjectCodeFunc != nil {
		return m.GetProjectCodeFunc(ctx)
	}
	return "TM"
}

// GetNextSequenceNumber implements repositories.AggregateRepository.
func (m *MockAggregateRepository) GetNextSequenceNumber(ctx context.Context, entityType string) (int, error) {
	if m.GetNextSequenceNumberFunc != nil {
		return m.GetNextSequenceNumberFunc(ctx, entityType)
	}
	return 1, nil
}

// Reset clears all configured behavior.
func (m *MockAggregateRepository) Reset() {
	m.GetRoadmapWithTracksFunc = nil
	m.GetProjectMetadataFunc = nil
	m.SetProjectMetadataFunc = nil
	m.GetProjectCodeFunc = nil
	m.GetNextSequenceNumberFunc = nil
}

// WithError configures the mock to return the specified error for methods that can fail.
func (m *MockAggregateRepository) WithError(err error) *MockAggregateRepository {
	m.GetRoadmapWithTracksFunc = func(ctx context.Context, roadmapID string) (*entities.RoadmapEntity, error) { return nil, err }
	m.GetProjectMetadataFunc = func(ctx context.Context, key string) (string, error) { return "", err }
	m.SetProjectMetadataFunc = func(ctx context.Context, key, value string) error { return err }
	m.GetNextSequenceNumberFunc = func(ctx context.Context, entityType string) (int, error) { return 0, err }
	return m
}
