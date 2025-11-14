package mocks

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// MockADRRepository is a mock implementation of repositories.ADRRepository for testing.
type MockADRRepository struct {
	// SaveADRFunc is called by SaveADR. If nil, returns nil.
	SaveADRFunc func(ctx context.Context, adr *entities.ADREntity) error

	// GetADRFunc is called by GetADR. If nil, returns nil, nil.
	GetADRFunc func(ctx context.Context, id string) (*entities.ADREntity, error)

	// ListADRsFunc is called by ListADRs. If nil, returns empty slice, nil.
	ListADRsFunc func(ctx context.Context, trackID *string) ([]*entities.ADREntity, error)

	// UpdateADRFunc is called by UpdateADR. If nil, returns nil.
	UpdateADRFunc func(ctx context.Context, adr *entities.ADREntity) error

	// SupersedeADRFunc is called by SupersedeADR. If nil, returns nil.
	SupersedeADRFunc func(ctx context.Context, adrID, supersededByID string) error

	// DeprecateADRFunc is called by DeprecateADR. If nil, returns nil.
	DeprecateADRFunc func(ctx context.Context, adrID string) error

	// GetADRsByTrackFunc is called by GetADRsByTrack. If nil, returns empty slice, nil.
	GetADRsByTrackFunc func(ctx context.Context, trackID string) ([]*entities.ADREntity, error)
}

// SaveADR implements repositories.ADRRepository.
func (m *MockADRRepository) SaveADR(ctx context.Context, adr *entities.ADREntity) error {
	if m.SaveADRFunc != nil {
		return m.SaveADRFunc(ctx, adr)
	}
	return nil
}

// GetADR implements repositories.ADRRepository.
func (m *MockADRRepository) GetADR(ctx context.Context, id string) (*entities.ADREntity, error) {
	if m.GetADRFunc != nil {
		return m.GetADRFunc(ctx, id)
	}
	return nil, nil
}

// ListADRs implements repositories.ADRRepository.
func (m *MockADRRepository) ListADRs(ctx context.Context, trackID *string) ([]*entities.ADREntity, error) {
	if m.ListADRsFunc != nil {
		return m.ListADRsFunc(ctx, trackID)
	}
	return []*entities.ADREntity{}, nil
}

// UpdateADR implements repositories.ADRRepository.
func (m *MockADRRepository) UpdateADR(ctx context.Context, adr *entities.ADREntity) error {
	if m.UpdateADRFunc != nil {
		return m.UpdateADRFunc(ctx, adr)
	}
	return nil
}

// SupersedeADR implements repositories.ADRRepository.
func (m *MockADRRepository) SupersedeADR(ctx context.Context, adrID, supersededByID string) error {
	if m.SupersedeADRFunc != nil {
		return m.SupersedeADRFunc(ctx, adrID, supersededByID)
	}
	return nil
}

// DeprecateADR implements repositories.ADRRepository.
func (m *MockADRRepository) DeprecateADR(ctx context.Context, adrID string) error {
	if m.DeprecateADRFunc != nil {
		return m.DeprecateADRFunc(ctx, adrID)
	}
	return nil
}

// GetADRsByTrack implements repositories.ADRRepository.
func (m *MockADRRepository) GetADRsByTrack(ctx context.Context, trackID string) ([]*entities.ADREntity, error) {
	if m.GetADRsByTrackFunc != nil {
		return m.GetADRsByTrackFunc(ctx, trackID)
	}
	return []*entities.ADREntity{}, nil
}

// Reset clears all configured behavior.
func (m *MockADRRepository) Reset() {
	m.SaveADRFunc = nil
	m.GetADRFunc = nil
	m.ListADRsFunc = nil
	m.UpdateADRFunc = nil
	m.SupersedeADRFunc = nil
	m.DeprecateADRFunc = nil
	m.GetADRsByTrackFunc = nil
}

// WithError configures the mock to return the specified error for all methods.
func (m *MockADRRepository) WithError(err error) *MockADRRepository {
	m.SaveADRFunc = func(ctx context.Context, adr *entities.ADREntity) error { return err }
	m.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) { return nil, err }
	m.ListADRsFunc = func(ctx context.Context, trackID *string) ([]*entities.ADREntity, error) { return nil, err }
	m.UpdateADRFunc = func(ctx context.Context, adr *entities.ADREntity) error { return err }
	m.SupersedeADRFunc = func(ctx context.Context, adrID, supersededByID string) error { return err }
	m.DeprecateADRFunc = func(ctx context.Context, adrID string) error { return err }
	m.GetADRsByTrackFunc = func(ctx context.Context, trackID string) ([]*entities.ADREntity, error) {
		return nil, err
	}
	return m
}
