package mocks

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// MockRoadmapRepository is a mock implementation of repositories.RoadmapRepository for testing.
type MockRoadmapRepository struct {
	// In-memory storage for testing
	roadmaps map[string]*entities.RoadmapEntity

	// SaveRoadmapFunc is called by SaveRoadmap. If nil, uses default implementation.
	SaveRoadmapFunc func(ctx context.Context, roadmap *entities.RoadmapEntity) error

	// GetRoadmapFunc is called by GetRoadmap. If nil, uses default implementation.
	GetRoadmapFunc func(ctx context.Context, id string) (*entities.RoadmapEntity, error)

	// GetActiveRoadmapFunc is called by GetActiveRoadmap. If nil, uses default implementation.
	GetActiveRoadmapFunc func(ctx context.Context) (*entities.RoadmapEntity, error)

	// UpdateRoadmapFunc is called by UpdateRoadmap. If nil, uses default implementation.
	UpdateRoadmapFunc func(ctx context.Context, roadmap *entities.RoadmapEntity) error
}

// NewMockRoadmapRepository creates a new mock roadmap repository with in-memory storage
func NewMockRoadmapRepository() *MockRoadmapRepository {
	return &MockRoadmapRepository{
		roadmaps: make(map[string]*entities.RoadmapEntity),
	}
}

// SaveRoadmap implements repositories.RoadmapRepository.
func (m *MockRoadmapRepository) SaveRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error {
	if m.SaveRoadmapFunc != nil {
		return m.SaveRoadmapFunc(ctx, roadmap)
	}
	// Default implementation: store in memory
	if _, exists := m.roadmaps[roadmap.ID]; exists {
		return pluginsdk.ErrAlreadyExists
	}
	m.roadmaps[roadmap.ID] = roadmap
	return nil
}

// GetRoadmap implements repositories.RoadmapRepository.
func (m *MockRoadmapRepository) GetRoadmap(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
	if m.GetRoadmapFunc != nil {
		return m.GetRoadmapFunc(ctx, id)
	}
	// Default implementation: get from memory
	roadmap, exists := m.roadmaps[id]
	if !exists {
		return nil, pluginsdk.ErrNotFound
	}
	return roadmap, nil
}

// GetActiveRoadmap implements repositories.RoadmapRepository.
func (m *MockRoadmapRepository) GetActiveRoadmap(ctx context.Context) (*entities.RoadmapEntity, error) {
	if m.GetActiveRoadmapFunc != nil {
		return m.GetActiveRoadmapFunc(ctx)
	}
	// Default implementation: return the first roadmap (assuming only one)
	for _, roadmap := range m.roadmaps {
		return roadmap, nil
	}
	return nil, pluginsdk.ErrNotFound
}

// UpdateRoadmap implements repositories.RoadmapRepository.
func (m *MockRoadmapRepository) UpdateRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error {
	if m.UpdateRoadmapFunc != nil {
		return m.UpdateRoadmapFunc(ctx, roadmap)
	}
	// Default implementation: update in memory
	if _, exists := m.roadmaps[roadmap.ID]; !exists {
		return pluginsdk.ErrNotFound
	}
	m.roadmaps[roadmap.ID] = roadmap
	return nil
}

// Reset clears all configured behavior.
func (m *MockRoadmapRepository) Reset() {
	m.SaveRoadmapFunc = nil
	m.GetRoadmapFunc = nil
	m.GetActiveRoadmapFunc = nil
	m.UpdateRoadmapFunc = nil
}

// WithError configures the mock to return the specified error for all methods.
func (m *MockRoadmapRepository) WithError(err error) *MockRoadmapRepository {
	m.SaveRoadmapFunc = func(ctx context.Context, roadmap *entities.RoadmapEntity) error { return err }
	m.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) { return nil, err }
	m.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) { return nil, err }
	m.UpdateRoadmapFunc = func(ctx context.Context, roadmap *entities.RoadmapEntity) error { return err }
	return m
}
