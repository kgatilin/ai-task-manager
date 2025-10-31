package task_manager_test

import (
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	tm "github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"
)

// MockRepository is a mock implementation of RoadmapRepository for testing
type MockRepository struct {
	activeRoadmap *tm.RoadmapEntity
	tracks        []*tm.TrackEntity
	tasks         []*tm.TaskEntity
	shouldError   bool
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		tracks: []*tm.TrackEntity{},
		tasks:  []*tm.TaskEntity{},
	}
}

func (m *MockRepository) SaveRoadmap(ctx context.Context, roadmap *tm.RoadmapEntity) error {
	if m.shouldError {
		return pluginsdk.ErrInternal
	}
	m.activeRoadmap = roadmap
	return nil
}

func (m *MockRepository) GetRoadmap(ctx context.Context, id string) (*tm.RoadmapEntity, error) {
	if m.shouldError {
		return nil, pluginsdk.ErrInternal
	}
	if m.activeRoadmap != nil && m.activeRoadmap.ID == id {
		return m.activeRoadmap, nil
	}
	return nil, pluginsdk.ErrNotFound
}

func (m *MockRepository) GetActiveRoadmap(ctx context.Context) (*tm.RoadmapEntity, error) {
	if m.shouldError {
		return nil, pluginsdk.ErrInternal
	}
	if m.activeRoadmap == nil {
		return nil, pluginsdk.ErrNotFound
	}
	return m.activeRoadmap, nil
}

func (m *MockRepository) UpdateRoadmap(ctx context.Context, roadmap *tm.RoadmapEntity) error {
	return nil
}

func (m *MockRepository) SaveTrack(ctx context.Context, track *tm.TrackEntity) error {
	return nil
}

func (m *MockRepository) GetTrack(ctx context.Context, id string) (*tm.TrackEntity, error) {
	if m.shouldError {
		return nil, pluginsdk.ErrInternal
	}
	for _, track := range m.tracks {
		if track.ID == id {
			return track, nil
		}
	}
	return nil, pluginsdk.ErrNotFound
}

func (m *MockRepository) ListTracks(ctx context.Context, roadmapID string, filters tm.TrackFilters) ([]*tm.TrackEntity, error) {
	if m.shouldError {
		return nil, pluginsdk.ErrInternal
	}
	return m.tracks, nil
}

func (m *MockRepository) UpdateTrack(ctx context.Context, track *tm.TrackEntity) error {
	return nil
}

func (m *MockRepository) DeleteTrack(ctx context.Context, id string) error {
	return nil
}

func (m *MockRepository) AddTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	return nil
}

func (m *MockRepository) RemoveTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	return nil
}

func (m *MockRepository) GetTrackDependencies(ctx context.Context, trackID string) ([]string, error) {
	return []string{}, nil
}

func (m *MockRepository) ValidateNoCycles(ctx context.Context, trackID string) error {
	return nil
}

func (m *MockRepository) SaveTask(ctx context.Context, task *tm.TaskEntity) error {
	return nil
}

func (m *MockRepository) GetTask(ctx context.Context, id string) (*tm.TaskEntity, error) {
	return nil, pluginsdk.ErrNotFound
}

func (m *MockRepository) ListTasks(ctx context.Context, filters tm.TaskFilters) ([]*tm.TaskEntity, error) {
	if m.shouldError {
		return nil, pluginsdk.ErrInternal
	}
	return m.tasks, nil
}

func (m *MockRepository) UpdateTask(ctx context.Context, task *tm.TaskEntity) error {
	return nil
}

func (m *MockRepository) DeleteTask(ctx context.Context, id string) error {
	return nil
}

func (m *MockRepository) MoveTaskToTrack(ctx context.Context, taskID, newTrackID string) error {
	return nil
}

func (m *MockRepository) SaveIteration(ctx context.Context, iteration *tm.IterationEntity) error {
	return nil
}

func (m *MockRepository) GetIteration(ctx context.Context, number int) (*tm.IterationEntity, error) {
	return nil, pluginsdk.ErrNotFound
}

func (m *MockRepository) GetCurrentIteration(ctx context.Context) (*tm.IterationEntity, error) {
	return nil, pluginsdk.ErrNotFound
}

func (m *MockRepository) ListIterations(ctx context.Context) ([]*tm.IterationEntity, error) {
	return []*tm.IterationEntity{}, nil
}

func (m *MockRepository) UpdateIteration(ctx context.Context, iteration *tm.IterationEntity) error {
	return nil
}

func (m *MockRepository) DeleteIteration(ctx context.Context, number int) error {
	return nil
}

func (m *MockRepository) AddTaskToIteration(ctx context.Context, iterationNum int, taskID string) error {
	return nil
}

func (m *MockRepository) RemoveTaskFromIteration(ctx context.Context, iterationNum int, taskID string) error {
	return nil
}

func (m *MockRepository) GetIterationTasks(ctx context.Context, iterationNum int) ([]*tm.TaskEntity, error) {
	return []*tm.TaskEntity{}, nil
}

func (m *MockRepository) StartIteration(ctx context.Context, iterationNum int) error {
	return nil
}

func (m *MockRepository) CompleteIteration(ctx context.Context, iterationNum int) error {
	return nil
}

func (m *MockRepository) GetRoadmapWithTracks(ctx context.Context, roadmapID string) (*tm.RoadmapEntity, error) {
	return nil, pluginsdk.ErrNotFound
}

func (m *MockRepository) GetTrackWithTasks(ctx context.Context, trackID string) (*tm.TrackEntity, error) {
	return nil, pluginsdk.ErrNotFound
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Tests

func TestNewAppModel(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepository()
	logger := NewMockLogger()

	model := tm.NewAppModel(ctx, repo, logger)

	if model == nil {
		t.Fatal("NewAppModel returned nil")
	}
}

func TestAppModelInit(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepository()
	logger := NewMockLogger()
	repo.activeRoadmap = &tm.RoadmapEntity{
		ID:               "roadmap-1",
		Vision:           "Test vision",
		SuccessCriteria:  "Test criteria",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	model := tm.NewAppModel(ctx, repo, logger)
	cmd := model.Init()

	if cmd == nil {
		t.Fatal("Init returned nil command")
	}
}

func TestAppModelRoadmapListView(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepository()
	logger := NewMockLogger()

	model := tm.NewAppModel(ctx, repo, logger)
	model.SetRoadmap(&tm.RoadmapEntity{
		ID:               "roadmap-1",
		Vision:           "Test vision",
		SuccessCriteria:  "Test criteria",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	})
	model.SetTracks([]*tm.TrackEntity{
		{
			ID:          "track-1",
			RoadmapID:   "roadmap-1",
			Title:       "Track 1",
			Description: "Description 1",
			Status:      "in-progress",
			Priority:    "high",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	})
	model.SetCurrentView(tm.ViewRoadmapList)
	model.SetDimensions(80, 24)

	view := model.View()

	if view == "" {
		t.Fatal("View returned empty string")
	}
	if !contains(view, "Track 1") {
		t.Fatalf("View should contain track title, got: %s", view)
	}
}

func TestAppModelTrackDetailView(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepository()
	logger := NewMockLogger()

	model := tm.NewAppModel(ctx, repo, logger)
	model.SetCurrentTrack(&tm.TrackEntity{
		ID:          "track-1",
		RoadmapID:   "roadmap-1",
		Title:       "Track 1",
		Description: "Description 1",
		Status:      "in-progress",
		Priority:    "high",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	model.SetTasks([]*tm.TaskEntity{
		{
			ID:          "task-1",
			TrackID:     "track-1",
			Title:       "Task 1",
			Description: "Description 1",
			Status:      "todo",
			Priority:    "medium",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	})
	model.SetCurrentView(tm.ViewTrackDetail)
	model.SetDimensions(80, 24)

	view := model.View()

	if view == "" {
		t.Fatal("View returned empty string")
	}
	if !contains(view, "Track 1") {
		t.Fatalf("View should contain track title, got: %s", view)
	}
	if !contains(view, "Task 1") {
		t.Fatalf("View should contain task title, got: %s", view)
	}
}

func TestAppModelKeyNavigation(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepository()
	logger := NewMockLogger()

	model := tm.NewAppModel(ctx, repo, logger)
	model.SetRoadmap(&tm.RoadmapEntity{
		ID:               "roadmap-1",
		Vision:           "Test vision",
		SuccessCriteria:  "Test criteria",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	})
	model.SetTracks([]*tm.TrackEntity{
		{
			ID:          "track-1",
			RoadmapID:   "roadmap-1",
			Title:       "Track 1",
			Description: "Description 1",
			Status:      "in-progress",
			Priority:    "high",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "track-2",
			RoadmapID:   "roadmap-1",
			Title:       "Track 2",
			Description: "Description 2",
			Status:      "todo",
			Priority:    "low",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	})
	model.SetCurrentView(tm.ViewRoadmapList)

	// Test navigation down
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if cmd != nil {
		t.Fatal("Unexpected command returned")
	}

	// Test navigation up
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if cmd != nil {
		t.Fatal("Unexpected command returned")
	}

	// Test quit
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("Expected quit command")
	}
}

func TestAppModelErrorView(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepository()
	logger := NewMockLogger()

	model := tm.NewAppModel(ctx, repo, logger)
	model.SetCurrentView(tm.ViewError)
	model.SetError(pluginsdk.ErrNotFound)

	view := model.View()

	if view == "" {
		t.Fatal("View returned empty string")
	}
	if !contains(view, "Error") {
		t.Fatalf("Error view should contain 'Error', got: %s", view)
	}
}

func TestAppModelLoadingView(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepository()
	logger := NewMockLogger()

	model := tm.NewAppModel(ctx, repo, logger)
	model.SetCurrentView(tm.ViewLoading)

	view := model.View()

	if view == "" {
		t.Fatal("View returned empty string")
	}
	if !contains(view, "Loading") {
		t.Fatalf("Loading view should contain 'Loading', got: %s", view)
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if len(s[i:]) < len(substr) {
			return false
		}
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
