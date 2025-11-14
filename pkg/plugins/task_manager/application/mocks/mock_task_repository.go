package mocks

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// MockTaskRepository is a mock implementation of repositories.TaskRepository for testing.
type MockTaskRepository struct {
	// In-memory storage for testing
	tasks map[string]*entities.TaskEntity

	// SaveTaskFunc is called by SaveTask. If nil, uses default implementation.
	SaveTaskFunc func(ctx context.Context, task *entities.TaskEntity) error

	// GetTaskFunc is called by GetTask. If nil, uses default implementation.
	GetTaskFunc func(ctx context.Context, id string) (*entities.TaskEntity, error)

	// ListTasksFunc is called by ListTasks. If nil, returns empty slice, nil.
	ListTasksFunc func(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error)

	// UpdateTaskFunc is called by UpdateTask. If nil, returns nil.
	UpdateTaskFunc func(ctx context.Context, task *entities.TaskEntity) error

	// DeleteTaskFunc is called by DeleteTask. If nil, returns nil.
	DeleteTaskFunc func(ctx context.Context, id string) error

	// MoveTaskToTrackFunc is called by MoveTaskToTrack. If nil, returns nil.
	MoveTaskToTrackFunc func(ctx context.Context, taskID, newTrackID string) error

	// GetBacklogTasksFunc is called by GetBacklogTasks. If nil, returns empty slice, nil.
	GetBacklogTasksFunc func(ctx context.Context) ([]*entities.TaskEntity, error)

	// GetIterationsForTaskFunc is called by GetIterationsForTask. If nil, returns empty slice, nil.
	GetIterationsForTaskFunc func(ctx context.Context, taskID string) ([]*entities.IterationEntity, error)
}

// NewMockTaskRepository creates a new mock task repository with in-memory storage
func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks: make(map[string]*entities.TaskEntity),
	}
}

// SaveTask implements repositories.TaskRepository.
func (m *MockTaskRepository) SaveTask(ctx context.Context, task *entities.TaskEntity) error {
	if m.SaveTaskFunc != nil {
		return m.SaveTaskFunc(ctx, task)
	}
	// Default implementation: store in memory
	m.tasks[task.ID] = task
	return nil
}

// GetTask implements repositories.TaskRepository.
func (m *MockTaskRepository) GetTask(ctx context.Context, id string) (*entities.TaskEntity, error) {
	if m.GetTaskFunc != nil {
		return m.GetTaskFunc(ctx, id)
	}
	// Default implementation: get from memory
	task, exists := m.tasks[id]
	if !exists {
		return nil, nil
	}
	return task, nil
}

// ListTasks implements repositories.TaskRepository.
func (m *MockTaskRepository) ListTasks(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error) {
	if m.ListTasksFunc != nil {
		return m.ListTasksFunc(ctx, filters)
	}
	// Default implementation: return all tasks
	var result []*entities.TaskEntity
	for _, task := range m.tasks {
		result = append(result, task)
	}
	return result, nil
}

// UpdateTask implements repositories.TaskRepository.
func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *entities.TaskEntity) error {
	if m.UpdateTaskFunc != nil {
		return m.UpdateTaskFunc(ctx, task)
	}
	return nil
}

// DeleteTask implements repositories.TaskRepository.
func (m *MockTaskRepository) DeleteTask(ctx context.Context, id string) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(ctx, id)
	}
	return nil
}

// MoveTaskToTrack implements repositories.TaskRepository.
func (m *MockTaskRepository) MoveTaskToTrack(ctx context.Context, taskID, newTrackID string) error {
	if m.MoveTaskToTrackFunc != nil {
		return m.MoveTaskToTrackFunc(ctx, taskID, newTrackID)
	}
	return nil
}

// GetBacklogTasks implements repositories.TaskRepository.
func (m *MockTaskRepository) GetBacklogTasks(ctx context.Context) ([]*entities.TaskEntity, error) {
	if m.GetBacklogTasksFunc != nil {
		return m.GetBacklogTasksFunc(ctx)
	}
	return []*entities.TaskEntity{}, nil
}

// GetIterationsForTask implements repositories.TaskRepository.
func (m *MockTaskRepository) GetIterationsForTask(ctx context.Context, taskID string) ([]*entities.IterationEntity, error) {
	if m.GetIterationsForTaskFunc != nil {
		return m.GetIterationsForTaskFunc(ctx, taskID)
	}
	return []*entities.IterationEntity{}, nil
}

// Reset clears all configured behavior.
func (m *MockTaskRepository) Reset() {
	m.SaveTaskFunc = nil
	m.GetTaskFunc = nil
	m.ListTasksFunc = nil
	m.UpdateTaskFunc = nil
	m.DeleteTaskFunc = nil
	m.MoveTaskToTrackFunc = nil
	m.GetBacklogTasksFunc = nil
	m.GetIterationsForTaskFunc = nil
}

// WithError configures the mock to return the specified error for all methods.
func (m *MockTaskRepository) WithError(err error) *MockTaskRepository {
	m.SaveTaskFunc = func(ctx context.Context, task *entities.TaskEntity) error { return err }
	m.GetTaskFunc = func(ctx context.Context, id string) (*entities.TaskEntity, error) { return nil, err }
	m.ListTasksFunc = func(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error) {
		return nil, err
	}
	m.UpdateTaskFunc = func(ctx context.Context, task *entities.TaskEntity) error { return err }
	m.DeleteTaskFunc = func(ctx context.Context, id string) error { return err }
	m.MoveTaskToTrackFunc = func(ctx context.Context, taskID, newTrackID string) error { return err }
	m.GetBacklogTasksFunc = func(ctx context.Context) ([]*entities.TaskEntity, error) { return nil, err }
	m.GetIterationsForTaskFunc = func(ctx context.Context, taskID string) ([]*entities.IterationEntity, error) {
		return nil, err
	}
	return m
}
