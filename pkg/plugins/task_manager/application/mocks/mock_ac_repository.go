package mocks

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// MockAcceptanceCriteriaRepository is a mock implementation of repositories.AcceptanceCriteriaRepository for testing.
type MockAcceptanceCriteriaRepository struct {
	// SaveACFunc is called by SaveAC. If nil, returns nil.
	SaveACFunc func(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error

	// GetACFunc is called by GetAC. If nil, returns nil, nil.
	GetACFunc func(ctx context.Context, id string) (*entities.AcceptanceCriteriaEntity, error)

	// ListACFunc is called by ListAC. If nil, returns empty slice, nil.
	ListACFunc func(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error)

	// UpdateACFunc is called by UpdateAC. If nil, returns nil.
	UpdateACFunc func(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error

	// DeleteACFunc is called by DeleteAC. If nil, returns nil.
	DeleteACFunc func(ctx context.Context, id string) error

	// ListACByTaskFunc is called by ListACByTask. If nil, returns empty slice, nil.
	ListACByTaskFunc func(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error)

	// ListACByIterationFunc is called by ListACByIteration. If nil, returns empty slice, nil.
	ListACByIterationFunc func(ctx context.Context, iterationNum int) ([]*entities.AcceptanceCriteriaEntity, error)

	// ListFailedACFunc is called by ListFailedAC. If nil, returns empty slice, nil.
	ListFailedACFunc func(ctx context.Context, filters entities.ACFilters) ([]*entities.AcceptanceCriteriaEntity, error)
}

// SaveAC implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) SaveAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error {
	if m.SaveACFunc != nil {
		return m.SaveACFunc(ctx, ac)
	}
	return nil
}

// GetAC implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) GetAC(ctx context.Context, id string) (*entities.AcceptanceCriteriaEntity, error) {
	if m.GetACFunc != nil {
		return m.GetACFunc(ctx, id)
	}
	return nil, nil
}

// ListAC implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) ListAC(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
	if m.ListACFunc != nil {
		return m.ListACFunc(ctx, taskID)
	}
	return []*entities.AcceptanceCriteriaEntity{}, nil
}

// UpdateAC implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) UpdateAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error {
	if m.UpdateACFunc != nil {
		return m.UpdateACFunc(ctx, ac)
	}
	return nil
}

// DeleteAC implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) DeleteAC(ctx context.Context, id string) error {
	if m.DeleteACFunc != nil {
		return m.DeleteACFunc(ctx, id)
	}
	return nil
}

// ListACByTask implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) ListACByTask(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
	if m.ListACByTaskFunc != nil {
		return m.ListACByTaskFunc(ctx, taskID)
	}
	return []*entities.AcceptanceCriteriaEntity{}, nil
}

// ListACByIteration implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) ListACByIteration(ctx context.Context, iterationNum int) ([]*entities.AcceptanceCriteriaEntity, error) {
	if m.ListACByIterationFunc != nil {
		return m.ListACByIterationFunc(ctx, iterationNum)
	}
	return []*entities.AcceptanceCriteriaEntity{}, nil
}

// ListFailedAC implements repositories.AcceptanceCriteriaRepository.
func (m *MockAcceptanceCriteriaRepository) ListFailedAC(ctx context.Context, filters entities.ACFilters) ([]*entities.AcceptanceCriteriaEntity, error) {
	if m.ListFailedACFunc != nil {
		return m.ListFailedACFunc(ctx, filters)
	}
	return []*entities.AcceptanceCriteriaEntity{}, nil
}

// Reset clears all configured behavior.
func (m *MockAcceptanceCriteriaRepository) Reset() {
	m.SaveACFunc = nil
	m.GetACFunc = nil
	m.ListACFunc = nil
	m.UpdateACFunc = nil
	m.DeleteACFunc = nil
	m.ListACByTaskFunc = nil
	m.ListACByIterationFunc = nil
	m.ListFailedACFunc = nil
}

// WithError configures the mock to return the specified error for all methods.
func (m *MockAcceptanceCriteriaRepository) WithError(err error) *MockAcceptanceCriteriaRepository {
	m.SaveACFunc = func(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error { return err }
	m.GetACFunc = func(ctx context.Context, id string) (*entities.AcceptanceCriteriaEntity, error) {
		return nil, err
	}
	m.ListACFunc = func(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
		return nil, err
	}
	m.UpdateACFunc = func(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error { return err }
	m.DeleteACFunc = func(ctx context.Context, id string) error { return err }
	m.ListACByTaskFunc = func(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
		return nil, err
	}
	m.ListACByIterationFunc = func(ctx context.Context, iterationNum int) ([]*entities.AcceptanceCriteriaEntity, error) {
		return nil, err
	}
	m.ListFailedACFunc = func(ctx context.Context, filters entities.ACFilters) ([]*entities.AcceptanceCriteriaEntity, error) {
		return nil, err
	}
	return m
}
