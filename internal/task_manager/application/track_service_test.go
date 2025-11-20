package application_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/mocks"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
)

// setupTrackTestService creates a test service with mock repositories
func setupTrackTestService(t *testing.T) (*application.TrackApplicationService, context.Context, *mocks.MockTrackRepository, *mocks.MockRoadmapRepository, *mocks.MockAggregateRepository) {
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockRoadmapRepo := &mocks.MockRoadmapRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}
	validationService := services.NewValidationService()

	service := application.NewTrackApplicationService(mockTrackRepo, mockRoadmapRepo, mockAggregateRepo, validationService)
	ctx := context.Background()

	return service, ctx, mockTrackRepo, mockRoadmapRepo, mockAggregateRepo
}

// createTestRoadmap creates a test roadmap entity for mock configuration
func createTestRoadmap(t *testing.T, roadmapID string) *entities.RoadmapEntity {
	now := time.Now().UTC()
	roadmap, err := entities.NewRoadmapEntity(roadmapID, "Test Roadmap", "Success criteria", now, now)
	if err != nil {
		t.Fatalf("failed to create test roadmap: %v", err)
	}
	return roadmap
}

// TestTrackService_CreateTrack_Success tests successful track creation
func TestTrackService_CreateTrack_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	// Configure mocks
	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		if id == roadmap.ID {
			return roadmap, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	mockTrackRepo.SaveTrackFunc = func(ctx context.Context, track *entities.TrackEntity) error {
		return nil // Success
	}

	input := dto.CreateTrackDTO{
		RoadmapID:   roadmap.ID,
		Title:       "Test Track",
		Description: "Test description",
		Status:      "not-started",
		Rank:        100,
	}

	track, err := service.CreateTrack(ctx, input)
	if err != nil {
		t.Fatalf("CreateTrack() failed: %v", err)
	}

	if track.ID == "" {
		t.Error("track.ID should not be empty (auto-generated)")
	}
	if track.Title != input.Title {
		t.Errorf("track.Title = %q, want %q", track.Title, input.Title)
	}
	if track.Status != input.Status {
		t.Errorf("track.Status = %q, want %q", track.Status, input.Status)
	}
	if track.Rank != input.Rank {
		t.Errorf("track.Rank = %d, want %d", track.Rank, input.Rank)
	}
}

// TestTrackService_CreateTrack_InvalidID tests track creation with invalid ID
// NOTE: This test is now obsolete because CreateTrackDTO no longer has an ID field.
// The service auto-generates IDs internally, so there's no "invalid ID" scenario for create operations.
// Keeping this test as a stub for documentation purposes.
func TestTrackService_CreateTrack_InvalidID(t *testing.T) {
	t.Skip("Test obsolete: CreateTrackDTO no longer accepts ID field (service auto-generates)")
}

// TestTrackService_CreateTrack_EmptyTitle tests track creation with empty title
func TestTrackService_CreateTrack_EmptyTitle(t *testing.T) {
	service, ctx, _, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	input := dto.CreateTrackDTO{
		RoadmapID:   roadmap.ID,
		Title:       "",
		Description: "Test description",
		Status:      "not-started",
		Rank:        100,
	}

	_, err := service.CreateTrack(ctx, input)
	if err == nil {
		t.Fatal("CreateTrack() should fail with empty title")
	}
}

// TestTrackService_CreateTrack_RoadmapNotFound tests track creation with non-existent roadmap
func TestTrackService_CreateTrack_RoadmapNotFound(t *testing.T) {
	service, ctx, _, mockRoadmapRepo, _ := setupTrackTestService(t)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return nil, tmerrors.ErrNotFound
	}

	input := dto.CreateTrackDTO{
		RoadmapID:   "nonexistent",
		Title:       "Test Track",
		Description: "Test description",
		Status:      "not-started",
		Rank:        100,
	}

	_, err := service.CreateTrack(ctx, input)
	if err == nil {
		t.Fatal("CreateTrack() should fail with non-existent roadmap")
	}
}

// TestTrackService_CreateTrack_InvalidRank tests track creation with invalid rank
func TestTrackService_CreateTrack_InvalidRank(t *testing.T) {
	service, ctx, _, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	input := dto.CreateTrackDTO{
		RoadmapID:   roadmap.ID,
		Title:       "Test Track",
		Description: "Test description",
		Status:      "not-started",
		Rank:        9999, // Invalid: must be 1-1000
	}

	_, err := service.CreateTrack(ctx, input)
	if err == nil {
		t.Fatal("CreateTrack() should fail with invalid rank")
	}
}

// TestTrackService_CreateTrack_DefaultStatus tests track creation with default status
func TestTrackService_CreateTrack_DefaultStatus(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.SaveTrackFunc = func(ctx context.Context, track *entities.TrackEntity) error {
		return nil
	}

	input := dto.CreateTrackDTO{
		RoadmapID:   roadmap.ID,
		Title:       "Test Track",
		Description: "Test description",
		Status:      "", // Empty status should default to not-started
		Rank:        100,
	}

	track, err := service.CreateTrack(ctx, input)
	if err != nil {
		t.Fatalf("CreateTrack() failed: %v", err)
	}

	if track.Status != "not-started" {
		t.Errorf("track.Status = %q, want %q", track.Status, "not-started")
	}
}

// TestTrackService_UpdateTrack_Success tests successful track update
func TestTrackService_UpdateTrack_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	existingTrack, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Original Title", "Original description", "not-started", 100, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == existingTrack.ID {
			return existingTrack, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	mockTrackRepo.SaveTrackFunc = func(ctx context.Context, track *entities.TrackEntity) error {
		return nil
	}

	mockTrackRepo.UpdateTrackFunc = func(ctx context.Context, track *entities.TrackEntity) error {
		return nil
	}

	// Update track
	newTitle := "Updated Title"
	newStatus := "in-progress"
	newRank := 200
	updateInput := dto.UpdateTrackDTO{
		ID:     existingTrack.ID, // MUST set ID for update operations
		Title:  &newTitle,
		Status: &newStatus,
		Rank:   &newRank,
	}

	track, err := service.UpdateTrack(ctx, updateInput)
	if err != nil {
		t.Fatalf("UpdateTrack() failed: %v", err)
	}

	if track.Title != newTitle {
		t.Errorf("track.Title = %q, want %q", track.Title, newTitle)
	}
	if track.Status != newStatus {
		t.Errorf("track.Status = %q, want %q", track.Status, newStatus)
	}
	if track.Rank != newRank {
		t.Errorf("track.Rank = %d, want %d", track.Rank, newRank)
	}
}

// TestTrackService_UpdateTrack_NotFound tests updating non-existent track
func TestTrackService_UpdateTrack_NotFound(t *testing.T) {
	service, ctx, mockTrackRepo, _, _ := setupTrackTestService(t)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return nil, tmerrors.ErrNotFound
	}

	newTitle := "Updated Title"
	updateInput := dto.UpdateTrackDTO{
		ID:    "nonexistent",
		Title: &newTitle,
	}

	_, err := service.UpdateTrack(ctx, updateInput)
	if err == nil {
		t.Fatal("UpdateTrack() should fail for non-existent track")
	}
}

// TestTrackService_UpdateTrack_PartialUpdate tests partial track update
func TestTrackService_UpdateTrack_PartialUpdate(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	existingTrack, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Original Title", "Original description", "not-started", 100, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == existingTrack.ID {
			return existingTrack, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	mockTrackRepo.UpdateTrackFunc = func(ctx context.Context, track *entities.TrackEntity) error {
		return nil
	}

	// Update only title
	newTitle := "Updated Title"
	updateInput := dto.UpdateTrackDTO{
		ID:    existingTrack.ID, // MUST set ID for update operations
		Title: &newTitle,
	}

	track, err := service.UpdateTrack(ctx, updateInput)
	if err != nil {
		t.Fatalf("UpdateTrack() failed: %v", err)
	}

	if track.Title != newTitle {
		t.Errorf("track.Title = %q, want %q", track.Title, newTitle)
	}
	// Other fields should remain unchanged
	if track.Description != existingTrack.Description {
		t.Errorf("track.Description changed: got %q, want %q", track.Description, existingTrack.Description)
	}
	if track.Status != existingTrack.Status {
		t.Errorf("track.Status changed: got %q, want %q", track.Status, existingTrack.Status)
	}
}

// TestTrackService_DeleteTrack_Success tests successful track deletion
func TestTrackService_DeleteTrack_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	existingTrack, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Test Track", "Test description", "not-started", 100, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	deleted := false
	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if deleted && id == existingTrack.ID {
			return nil, tmerrors.ErrNotFound
		}
		if id == existingTrack.ID {
			return existingTrack, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	mockTrackRepo.DeleteTrackFunc = func(ctx context.Context, id string) error {
		if id == existingTrack.ID {
			deleted = true
			return nil
		}
		return tmerrors.ErrNotFound
	}

	// Delete track
	err := service.DeleteTrack(ctx, "TM-track-1")
	if err != nil {
		t.Fatalf("DeleteTrack() failed: %v", err)
	}

	// Verify track is deleted
	_, err = service.GetTrack(ctx, "TM-track-1")
	if err == nil {
		t.Fatal("GetTrack() should fail after deletion")
	}
}

// TestTrackService_DeleteTrack_NotFound tests deleting non-existent track
func TestTrackService_DeleteTrack_NotFound(t *testing.T) {
	service, ctx, mockTrackRepo, _, _ := setupTrackTestService(t)

	mockTrackRepo.DeleteTrackFunc = func(ctx context.Context, id string) error {
		return tmerrors.ErrNotFound
	}

	err := service.DeleteTrack(ctx, "nonexistent")
	if err == nil {
		t.Fatal("DeleteTrack() should fail for non-existent track")
	}
}

// TestTrackService_GetTrack_Success tests successful track retrieval
func TestTrackService_GetTrack_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	existingTrack, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Test Track", "Test description", "not-started", 100, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == existingTrack.ID {
			return existingTrack, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	// Get track
	track, err := service.GetTrack(ctx, "TM-track-1")
	if err != nil {
		t.Fatalf("GetTrack() failed: %v", err)
	}

	if track.ID != existingTrack.ID {
		t.Errorf("track.ID = %q, want %q", track.ID, existingTrack.ID)
	}
	if track.Title != existingTrack.Title {
		t.Errorf("track.Title = %q, want %q", track.Title, existingTrack.Title)
	}
}

// TestTrackService_GetTrack_NotFound tests retrieving non-existent track
func TestTrackService_GetTrack_NotFound(t *testing.T) {
	service, ctx, mockTrackRepo, _, _ := setupTrackTestService(t)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return nil, tmerrors.ErrNotFound
	}

	_, err := service.GetTrack(ctx, "nonexistent")
	if err == nil {
		t.Fatal("GetTrack() should fail for non-existent track")
	}
}

// TestTrackService_ListTracks_Success tests successful track listing
func TestTrackService_ListTracks_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track1, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Track 1", "", "not-started", 100, []string{}, now, now)
	track2, _ := entities.NewTrackEntity("TM-track-2", roadmap.ID, "Track 2", "", "in-progress", 200, []string{}, now, now)
	track3, _ := entities.NewTrackEntity("TM-track-3", roadmap.ID, "Track 3", "", "complete", 300, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.ListTracksFunc = func(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
		return []*entities.TrackEntity{track1, track2, track3}, nil
	}

	// List all tracks
	filters := entities.TrackFilters{}
	results, err := service.ListTracks(ctx, roadmap.ID, filters)
	if err != nil {
		t.Fatalf("ListTracks() failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("ListTracks() returned %d tracks, want 3", len(results))
	}
}

// TestTrackService_ListTracks_WithFilters tests track listing with status filter
func TestTrackService_ListTracks_WithFilters(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track2, _ := entities.NewTrackEntity("TM-track-2", roadmap.ID, "Track 2", "", "in-progress", 200, []string{}, now, now)
	track3, _ := entities.NewTrackEntity("TM-track-3", roadmap.ID, "Track 3", "", "in-progress", 300, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.ListTracksFunc = func(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
		if len(filters.Status) > 0 && filters.Status[0] == "in-progress" {
			return []*entities.TrackEntity{track2, track3}, nil
		}
		return []*entities.TrackEntity{}, nil
	}

	// List tracks with status filter
	filters := entities.TrackFilters{Status: []string{"in-progress"}}
	results, err := service.ListTracks(ctx, roadmap.ID, filters)
	if err != nil {
		t.Fatalf("ListTracks() failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("ListTracks() returned %d tracks, want 2", len(results))
	}

	for _, track := range results {
		if track.Status != "in-progress" {
			t.Errorf("track.Status = %q, want %q", track.Status, "in-progress")
		}
	}
}

// TestTrackService_ListTracks_Empty tests listing tracks from empty roadmap
func TestTrackService_ListTracks_Empty(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.ListTracksFunc = func(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
		return []*entities.TrackEntity{}, nil
	}

	// List tracks from empty roadmap
	filters := entities.TrackFilters{}
	results, err := service.ListTracks(ctx, roadmap.ID, filters)
	if err != nil {
		t.Fatalf("ListTracks() failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("ListTracks() returned %d tracks, want 0", len(results))
	}
}

// TestTrackService_GetTrackWithTasks_Success tests retrieving track with tasks
func TestTrackService_GetTrackWithTasks_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	existingTrack, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Test Track", "Test description", "not-started", 100, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackWithTasksFunc = func(ctx context.Context, trackID string) (*entities.TrackEntity, error) {
		if trackID == existingTrack.ID {
			return existingTrack, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	// Get track with tasks
	track, err := service.GetTrackWithTasks(ctx, "TM-track-1")
	if err != nil {
		t.Fatalf("GetTrackWithTasks() failed: %v", err)
	}

	if track.ID != existingTrack.ID {
		t.Errorf("track.ID = %q, want %q", track.ID, existingTrack.ID)
	}
}

// TestTrackService_GetTrackWithTasks_NotFound tests retrieving non-existent track
func TestTrackService_GetTrackWithTasks_NotFound(t *testing.T) {
	service, ctx, mockTrackRepo, _, _ := setupTrackTestService(t)

	mockTrackRepo.GetTrackWithTasksFunc = func(ctx context.Context, trackID string) (*entities.TrackEntity, error) {
		return nil, tmerrors.ErrNotFound
	}

	_, err := service.GetTrackWithTasks(ctx, "nonexistent")
	if err == nil {
		t.Fatal("GetTrackWithTasks() should fail for non-existent track")
	}
}

// TestTrackService_AddDependency_Success tests adding a track dependency
func TestTrackService_AddDependency_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track1, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Track 1", "", "not-started", 100, []string{}, now, now)
	track2, _ := entities.NewTrackEntity("TM-track-2", roadmap.ID, "Track 2", "", "not-started", 200, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track1.ID {
			return track1, nil
		}
		if id == track2.ID {
			return track2, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	mockTrackRepo.ValidateNoCyclesFunc = func(ctx context.Context, trackID string) error {
		return nil // No cycles
	}

	mockTrackRepo.AddTrackDependencyFunc = func(ctx context.Context, trackID, dependsOnID string) error {
		return nil
	}

	dependencies := []string{}
	mockTrackRepo.GetTrackDependenciesFunc = func(ctx context.Context, trackID string) ([]string, error) {
		if trackID == track2.ID {
			return dependencies, nil
		}
		return []string{}, nil
	}

	// Add dependency: track-2 depends on track-1
	err := service.AddDependency(ctx, "TM-track-2", "TM-track-1")
	if err != nil {
		t.Fatalf("AddDependency() failed: %v", err)
	}

	dependencies = []string{"TM-track-1"}

	// Verify dependency was added
	deps, err := service.GetDependencies(ctx, "TM-track-2")
	if err != nil {
		t.Fatalf("GetDependencies() failed: %v", err)
	}

	if len(deps) != 1 || deps[0] != "TM-track-1" {
		t.Errorf("GetDependencies() = %v, want [TM-track-1]", deps)
	}
}

// TestTrackService_AddDependency_CircularDetection tests circular dependency detection
func TestTrackService_AddDependency_CircularDetection(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track1, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Track 1", "", "not-started", 100, []string{}, now, now)
	track2, _ := entities.NewTrackEntity("TM-track-2", roadmap.ID, "Track 2", "", "not-started", 200, []string{}, now, now)
	track3, _ := entities.NewTrackEntity("TM-track-3", roadmap.ID, "Track 3", "", "not-started", 300, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track1.ID {
			return track1, nil
		}
		if id == track2.ID {
			return track2, nil
		}
		if id == track3.ID {
			return track3, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	callCount := 0
	mockTrackRepo.ValidateNoCyclesFunc = func(ctx context.Context, trackID string) error {
		callCount++
		if callCount > 2 {
			// Third call would create a cycle
			return tmerrors.ErrInvalidArgument
		}
		return nil
	}

	mockTrackRepo.AddTrackDependencyFunc = func(ctx context.Context, trackID, dependsOnID string) error {
		return nil
	}

	// Create chain: track-2 -> track-1
	err := service.AddDependency(ctx, "TM-track-2", "TM-track-1")
	if err != nil {
		t.Fatalf("AddDependency() failed: %v", err)
	}

	// Create chain: track-3 -> track-2
	err = service.AddDependency(ctx, "TM-track-3", "TM-track-2")
	if err != nil {
		t.Fatalf("AddDependency() failed: %v", err)
	}

	// Try to create cycle: track-1 -> track-3 (should fail)
	err = service.AddDependency(ctx, "TM-track-1", "TM-track-3")
	if err == nil {
		t.Fatal("AddDependency() should fail with circular dependency")
	}
}

// TestTrackService_AddDependency_SelfDependency tests self-dependency prevention
func TestTrackService_AddDependency_SelfDependency(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track1, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Test Track", "Test description", "not-started", 100, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track1.ID {
			return track1, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	// Try to add self-dependency
	err := service.AddDependency(ctx, "TM-track-1", "TM-track-1")
	if err == nil {
		t.Fatal("AddDependency() should fail with self-dependency")
	}
}

// TestTrackService_AddDependency_TrackNotFound tests adding dependency with non-existent track
func TestTrackService_AddDependency_TrackNotFound(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track1, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Test Track", "Test description", "not-started", 100, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track1.ID {
			return track1, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	// Try to add dependency to non-existent track
	err := service.AddDependency(ctx, "TM-track-1", "nonexistent")
	if err == nil {
		t.Fatal("AddDependency() should fail with non-existent dependency track")
	}
}

// TestTrackService_RemoveDependency_Success tests removing a track dependency
func TestTrackService_RemoveDependency_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track1, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Track 1", "", "not-started", 100, []string{}, now, now)
	track2, _ := entities.NewTrackEntity("TM-track-2", roadmap.ID, "Track 2", "", "not-started", 200, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track1.ID {
			return track1, nil
		}
		if id == track2.ID {
			return track2, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	mockTrackRepo.ValidateNoCyclesFunc = func(ctx context.Context, trackID string) error {
		return nil
	}

	mockTrackRepo.AddTrackDependencyFunc = func(ctx context.Context, trackID, dependsOnID string) error {
		return nil
	}

	removed := false
	mockTrackRepo.RemoveTrackDependencyFunc = func(ctx context.Context, trackID, dependsOnID string) error {
		if trackID == track2.ID && dependsOnID == track1.ID {
			removed = true
			return nil
		}
		return tmerrors.ErrNotFound
	}

	mockTrackRepo.GetTrackDependenciesFunc = func(ctx context.Context, trackID string) ([]string, error) {
		if trackID == track2.ID && !removed {
			return []string{"TM-track-1"}, nil
		}
		return []string{}, nil
	}

	// Add dependency first
	err := service.AddDependency(ctx, "TM-track-2", "TM-track-1")
	if err != nil {
		t.Fatalf("AddDependency() failed: %v", err)
	}

	// Remove dependency
	err = service.RemoveDependency(ctx, "TM-track-2", "TM-track-1")
	if err != nil {
		t.Fatalf("RemoveDependency() failed: %v", err)
	}

	// Verify dependency was removed
	deps, err := service.GetDependencies(ctx, "TM-track-2")
	if err != nil {
		t.Fatalf("GetDependencies() failed: %v", err)
	}

	if len(deps) != 0 {
		t.Errorf("GetDependencies() = %v, want []", deps)
	}
}

// TestTrackService_GetDependencies_Success tests retrieving track dependencies
func TestTrackService_GetDependencies_Success(t *testing.T) {
	service, ctx, mockTrackRepo, mockRoadmapRepo, _ := setupTrackTestService(t)
	roadmap := createTestRoadmap(t, "roadmap-1")

	now := time.Now().UTC()
	track1, _ := entities.NewTrackEntity("TM-track-1", roadmap.ID, "Track 1", "", "not-started", 100, []string{}, now, now)
	track2, _ := entities.NewTrackEntity("TM-track-2", roadmap.ID, "Track 2", "", "not-started", 200, []string{}, now, now)
	track3, _ := entities.NewTrackEntity("TM-track-3", roadmap.ID, "Track 3", "", "not-started", 300, []string{}, now, now)

	mockRoadmapRepo.GetRoadmapFunc = func(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track1.ID {
			return track1, nil
		}
		if id == track2.ID {
			return track2, nil
		}
		if id == track3.ID {
			return track3, nil
		}
		return nil, tmerrors.ErrNotFound
	}

	mockTrackRepo.ValidateNoCyclesFunc = func(ctx context.Context, trackID string) error {
		return nil
	}

	mockTrackRepo.AddTrackDependencyFunc = func(ctx context.Context, trackID, dependsOnID string) error {
		return nil
	}

	dependencies := []string{}
	mockTrackRepo.GetTrackDependenciesFunc = func(ctx context.Context, trackID string) ([]string, error) {
		if trackID == track3.ID {
			return dependencies, nil
		}
		return []string{}, nil
	}

	// Add multiple dependencies
	err := service.AddDependency(ctx, "TM-track-3", "TM-track-1")
	if err != nil {
		t.Fatalf("AddDependency() failed: %v", err)
	}
	dependencies = append(dependencies, "TM-track-1")

	err = service.AddDependency(ctx, "TM-track-3", "TM-track-2")
	if err != nil {
		t.Fatalf("AddDependency() failed: %v", err)
	}
	dependencies = append(dependencies, "TM-track-2")

	// Get dependencies
	deps, err := service.GetDependencies(ctx, "TM-track-3")
	if err != nil {
		t.Fatalf("GetDependencies() failed: %v", err)
	}

	if len(deps) != 2 {
		t.Fatalf("GetDependencies() returned %d dependencies, want 2", len(deps))
	}
}

// TestTrackService_GetActiveRoadmap_Success tests getting the active roadmap
func TestTrackService_GetActiveRoadmap_Success(t *testing.T) {
	service, ctx, _, mockRoadmapRepo, _ := setupTrackTestService(t)

	expectedRoadmap := createTestRoadmap(t, "TM-roadmap-1")

	mockRoadmapRepo.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) {
		return expectedRoadmap, nil
	}

	roadmap, err := service.GetActiveRoadmap(ctx)
	if err != nil {
		t.Fatalf("GetActiveRoadmap() failed: %v", err)
	}

	if roadmap == nil {
		t.Fatal("Expected roadmap, got nil")
	}

	if roadmap.ID != expectedRoadmap.ID {
		t.Errorf("Roadmap ID = %s, want %s", roadmap.ID, expectedRoadmap.ID)
	}
}

// TestTrackService_GetActiveRoadmap_NotFound tests getting active roadmap when none exists
func TestTrackService_GetActiveRoadmap_NotFound(t *testing.T) {
	service, ctx, _, mockRoadmapRepo, _ := setupTrackTestService(t)

	mockRoadmapRepo.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) {
		return nil, tmerrors.ErrNotFound
	}

	_, err := service.GetActiveRoadmap(ctx)
	if err == nil {
		t.Error("Expected error when no active roadmap exists")
	}
}

// TestTrackService_GetActiveRoadmap_RepositoryError tests error handling
func TestTrackService_GetActiveRoadmap_RepositoryError(t *testing.T) {
	service, ctx, _, mockRoadmapRepo, _ := setupTrackTestService(t)

	mockRoadmapRepo.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) {
		return nil, fmt.Errorf("database error")
	}

	_, err := service.GetActiveRoadmap(ctx)
	if err == nil {
		t.Error("Expected error from repository")
	}
}
