package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/infrastructure/persistence"
)

// ============================================================================
// Track Tests
// ============================================================================

func TestSaveAndGetTrack(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Create roadmap first
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	// Create track
	track, err := entities.NewTrackEntity(
		"track-core",
		"roadmap-1",
		"Core Features",
		"Essential features",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track entity: %v", err)
	}

	// Save track
	if err := trackRepo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	// Get track
	retrieved, err := trackRepo.GetTrack(ctx, "track-core")
	if err != nil {
		t.Fatalf("failed to get track: %v", err)
	}

	if retrieved.ID != track.ID {
		t.Errorf("expected track ID %s, got %s", track.ID, retrieved.ID)
	}
}

func TestListTracks(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Create roadmap
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	// Create tracks
	for i := 1; i <= 3; i++ {
		id := "track-" + string(rune(48+i))
		track, _ := entities.NewTrackEntity(id, "roadmap-1", "Track "+string(rune(48+i)), "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
		trackRepo.SaveTrack(ctx, track)
	}

	// List all tracks
	tracks, err := trackRepo.ListTracks(ctx, "roadmap-1", entities.TrackFilters{})
	if err != nil {
		t.Fatalf("failed to list tracks: %v", err)
	}

	if len(tracks) != 3 {
		t.Errorf("expected 3 tracks, got %d", len(tracks))
	}
}

func TestListTracksWithFilters(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Create roadmap
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	// Create tracks with different statuses
	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "", "in-progress", 200, []string{}, time.Now().UTC(), time.Now().UTC())

	trackRepo.SaveTrack(ctx, track1)
	trackRepo.SaveTrack(ctx, track2)

	// Filter by status
	tracks, err := trackRepo.ListTracks(ctx, "roadmap-1", entities.TrackFilters{Status: []string{"in-progress"}})
	if err != nil {
		t.Fatalf("failed to list tracks: %v", err)
	}

	if len(tracks) != 1 || tracks[0].Status != "in-progress" {
		t.Errorf("expected 1 in-progress track, got %d", len(tracks))
	}
}

func TestTrackDependencies(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())

	trackRepo.SaveTrack(ctx, track1)
	trackRepo.SaveTrack(ctx, track2)

	// Add dependency
	if err := trackRepo.AddTrackDependency(ctx, "track-2", "track-1"); err != nil {
		t.Fatalf("failed to add dependency: %v", err)
	}

	// Get dependencies
	deps, err := trackRepo.GetTrackDependencies(ctx, "track-2")
	if err != nil {
		t.Fatalf("failed to get dependencies: %v", err)
	}

	if len(deps) != 1 || deps[0] != "track-1" {
		t.Errorf("expected track-1 dependency, got %v", deps)
	}

	// Remove dependency
	if err := trackRepo.RemoveTrackDependency(ctx, "track-2", "track-1"); err != nil {
		t.Fatalf("failed to remove dependency: %v", err)
	}

	deps, _ = trackRepo.GetTrackDependencies(ctx, "track-2")
	if len(deps) != 0 {
		t.Errorf("expected no dependencies, got %v", deps)
	}
}

func TestValidateNoCycles(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	track3, _ := entities.NewTrackEntity("track-3", "roadmap-1", "Track 3", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())

	trackRepo.SaveTrack(ctx, track1)
	trackRepo.SaveTrack(ctx, track2)
	trackRepo.SaveTrack(ctx, track3)

	// Create a cycle: 1 -> 2 -> 3 -> 1
	trackRepo.AddTrackDependency(ctx, "track-2", "track-1")
	trackRepo.AddTrackDependency(ctx, "track-3", "track-2")
	trackRepo.AddTrackDependency(ctx, "track-1", "track-3")

	// Validate should detect cycle
	err := trackRepo.ValidateNoCycles(ctx, "track-1")
	if err == nil {
		t.Error("expected error for cycle detection")
	} else if !errors.Is(err, pluginsdk.ErrInvalidArgument) {
		t.Errorf("expected ErrInvalidArgument, got: %v", err)
	}
}

func TestAddDependencyToNonexistentTrack(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	err := trackRepo.AddTrackDependency(ctx, "nonexistent", "also-nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent track")
	} else if !errors.Is(err, pluginsdk.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestSelfDependency(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Try self dependency
	err := trackRepo.AddTrackDependency(ctx, "track-1", "track-1")
	if err == nil {
		t.Error("expected error for self dependency")
	} else if !errors.Is(err, pluginsdk.ErrInvalidArgument) {
		t.Errorf("expected ErrInvalidArgument, got: %v", err)
	}
}

func TestUpdateTrack(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Update track
	track.Title = "Updated Title"
	track.Status = "in-progress"
	track.UpdatedAt = time.Now().UTC()

	if err := trackRepo.UpdateTrack(ctx, track); err != nil {
		t.Fatalf("failed to update track: %v", err)
	}

	// Verify update
	retrieved, _ := trackRepo.GetTrack(ctx, "track-1")
	if retrieved.Title != "Updated Title" {
		t.Errorf("expected title to be updated, got %s", retrieved.Title)
	}
	if retrieved.Status != "in-progress" {
		t.Errorf("expected status in-progress, got %s", retrieved.Status)
	}
}

func TestDeleteTrack(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Delete track
	if err := trackRepo.DeleteTrack(ctx, "track-1"); err != nil {
		t.Fatalf("failed to delete track: %v", err)
	}

	// Verify deletion
	_, err := trackRepo.GetTrack(ctx, "track-1")
	if !errors.Is(err, pluginsdk.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
