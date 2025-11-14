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
// Aggregate Repository Tests
// ============================================================================

func TestNewSQLiteAggregateRepository(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	logger := createTestLogger()
	repo := persistence.NewSQLiteAggregateRepository(db, logger)

	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
	if repo.DB != db {
		t.Error("expected DB to be set")
	}
}

// ============================================================================
// GetRoadmapWithTracks Tests
// ============================================================================

func TestGetRoadmapWithTracks_Success(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())
	roadmapRepo := persistence.NewSQLiteRoadmapOnlyRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())

	// Create roadmap
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	// Create tracks
	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track1)

	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "desc", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track2)

	// Get roadmap with tracks
	result, err := repo.GetRoadmapWithTracks(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("GetRoadmapWithTracks failed: %v", err)
	}

	if result.ID != "roadmap-1" {
		t.Errorf("expected roadmap-1, got %s", result.ID)
	}
	if result.Vision != "vision" {
		t.Errorf("expected vision, got %s", result.Vision)
	}
}

func TestGetRoadmapWithTracks_NotFound(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	_, err := repo.GetRoadmapWithTracks(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent roadmap")
	}
	if !errors.Is(err, pluginsdk.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestGetRoadmapWithTracks_EmptyTracks(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())
	roadmapRepo := persistence.NewSQLiteRoadmapOnlyRepository(db, createTestLogger())

	// Create roadmap without tracks
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	// Get roadmap with no tracks
	result, err := repo.GetRoadmapWithTracks(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("GetRoadmapWithTracks failed: %v", err)
	}

	if result.ID != "roadmap-1" {
		t.Errorf("expected roadmap-1, got %s", result.ID)
	}
}

// ============================================================================
// Project Metadata Tests
// ============================================================================

func TestSetAndGetProjectMetadata(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Set metadata
	err := repo.SetProjectMetadata(ctx, "test_key", "test_value")
	if err != nil {
		t.Fatalf("SetProjectMetadata failed: %v", err)
	}

	// Get metadata
	value, err := repo.GetProjectMetadata(ctx, "test_key")
	if err != nil {
		t.Fatalf("GetProjectMetadata failed: %v", err)
	}

	if value != "test_value" {
		t.Errorf("expected test_value, got %s", value)
	}
}

func TestGetProjectMetadata_NotFound(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	_, err := repo.GetProjectMetadata(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent metadata")
	}
	if !errors.Is(err, pluginsdk.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestSetProjectMetadata_Replace(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Set initial value
	repo.SetProjectMetadata(ctx, "test_key", "value1")

	// Replace with new value
	err := repo.SetProjectMetadata(ctx, "test_key", "value2")
	if err != nil {
		t.Fatalf("SetProjectMetadata failed: %v", err)
	}

	// Verify new value
	value, _ := repo.GetProjectMetadata(ctx, "test_key")
	if value != "value2" {
		t.Errorf("expected value2, got %s", value)
	}
}

func TestGetProjectCode_Default(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Get project code without setting it (should return default "DW")
	code := repo.GetProjectCode(ctx)
	if code != "DW" {
		t.Errorf("expected default DW, got %s", code)
	}
}

func TestGetProjectCode_Custom(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Set custom project code
	repo.SetProjectMetadata(ctx, "project_code", "CUSTOM")

	// Get project code
	code := repo.GetProjectCode(ctx)
	if code != "CUSTOM" {
		t.Errorf("expected CUSTOM, got %s", code)
	}
}

// ============================================================================
// Sequence Number Tests
// ============================================================================

func TestGetNextSequenceNumber_Task_Empty(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Get next sequence for empty database
	seq, err := repo.GetNextSequenceNumber(ctx, "task")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 1 {
		t.Errorf("expected 1, got %d", seq)
	}
}

func TestGetNextSequenceNumber_Task_WithExisting(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())
	roadmapRepo := persistence.NewSQLiteRoadmapOnlyRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())

	// Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create tasks with IDs like "DW-task-5", "DW-task-12"
	task1, _ := entities.NewTaskEntity("DW-task-5", "track-1", "Task 1", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task1)

	task2, _ := entities.NewTaskEntity("DW-task-12", "track-1", "Task 2", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task2)

	// Get next sequence (should be 13)
	seq, err := repo.GetNextSequenceNumber(ctx, "task")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 13 {
		t.Errorf("expected 13, got %d", seq)
	}
}

func TestGetNextSequenceNumber_Track_Empty(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Get next sequence for empty database
	seq, err := repo.GetNextSequenceNumber(ctx, "track")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 1 {
		t.Errorf("expected 1, got %d", seq)
	}
}

func TestGetNextSequenceNumber_Track_WithExisting(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())
	roadmapRepo := persistence.NewSQLiteRoadmapOnlyRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())

	// Create roadmap
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	// Create tracks with IDs like "DW-track-3", "DW-track-7"
	track1, _ := entities.NewTrackEntity("DW-track-3", "roadmap-1", "Track 1", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track1)

	track2, _ := entities.NewTrackEntity("DW-track-7", "roadmap-1", "Track 2", "desc", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track2)

	// Get next sequence (should be 8)
	seq, err := repo.GetNextSequenceNumber(ctx, "track")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 8 {
		t.Errorf("expected 8, got %d", seq)
	}
}

func TestGetNextSequenceNumber_Iteration_Empty(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Get next sequence for empty database
	seq, err := repo.GetNextSequenceNumber(ctx, "iter")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 1 {
		t.Errorf("expected 1, got %d", seq)
	}
}

func TestGetNextSequenceNumber_Iteration_WithExisting(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())
	iterRepo := persistence.NewSQLiteIterationRepository(db, createTestLogger())

	// Create iterations
	now := time.Now().UTC()
	iter1, _ := entities.NewIterationEntity(1, "Iter 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	iterRepo.SaveIteration(ctx, iter1)

	iter2, _ := entities.NewIterationEntity(5, "Iter 5", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	iterRepo.SaveIteration(ctx, iter2)

	// Get next sequence (should be 6)
	seq, err := repo.GetNextSequenceNumber(ctx, "iter")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 6 {
		t.Errorf("expected 6, got %d", seq)
	}
}

func TestGetNextSequenceNumber_AC_Empty(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Get next sequence for empty database
	seq, err := repo.GetNextSequenceNumber(ctx, "ac")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 1 {
		t.Errorf("expected 1, got %d", seq)
	}
}

func TestGetNextSequenceNumber_AC_WithExisting(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())
	roadmapRepo := persistence.NewSQLiteRoadmapOnlyRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())

	// Create roadmap, track, and task
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	// Create ACs with IDs like "DW-ac-2", "DW-ac-9"
	now := time.Now().UTC()
	ac1 := entities.NewAcceptanceCriteriaEntity("DW-ac-2", "task-1", "AC 1", entities.VerificationTypeManual, "instructions", now, now)
	acRepo.SaveAC(ctx, ac1)

	ac2 := entities.NewAcceptanceCriteriaEntity("DW-ac-9", "task-1", "AC 2", entities.VerificationTypeManual, "instructions", now, now)
	acRepo.SaveAC(ctx, ac2)

	// Get next sequence (should be 10)
	seq, err := repo.GetNextSequenceNumber(ctx, "ac")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 10 {
		t.Errorf("expected 10, got %d", seq)
	}
}

func TestGetNextSequenceNumber_ADR_Empty(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	// Get next sequence for empty database
	seq, err := repo.GetNextSequenceNumber(ctx, "adr")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 1 {
		t.Errorf("expected 1, got %d", seq)
	}
}

func TestGetNextSequenceNumber_ADR_WithExisting(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())
	roadmapRepo := persistence.NewSQLiteRoadmapOnlyRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	adrRepo := persistence.NewSQLiteADRRepository(db, createTestLogger())

	// Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create ADRs with IDs like "DW-adr-4", "DW-adr-11"
	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("DW-adr-4", "track-1", "ADR 1", "proposed", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	adrRepo.SaveADR(ctx, adr1)

	adr2, _ := entities.NewADREntity("DW-adr-11", "track-1", "ADR 2", "proposed", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	adrRepo.SaveADR(ctx, adr2)

	// Get next sequence (should be 12)
	seq, err := repo.GetNextSequenceNumber(ctx, "adr")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}

	if seq != 12 {
		t.Errorf("expected 12, got %d", seq)
	}
}

func TestGetNextSequenceNumber_InvalidEntityType(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := persistence.NewSQLiteAggregateRepository(db, createTestLogger())

	_, err := repo.GetNextSequenceNumber(ctx, "invalid_type")
	if err == nil {
		t.Error("expected error for invalid entity type")
	}
	if !errors.Is(err, pluginsdk.ErrInvalidArgument) {
		t.Errorf("expected ErrInvalidArgument, got: %v", err)
	}
}
