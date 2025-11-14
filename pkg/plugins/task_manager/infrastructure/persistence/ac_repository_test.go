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
// Acceptance Criteria Tests
// ============================================================================

func TestSaveAndGetAC(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	// Create AC
	ac := entities.NewAcceptanceCriteriaEntity(
		"ac-1",
		"task-1",
		"Feature works correctly",
		entities.VerificationTypeManual,
		"Run tests and verify",
		time.Now().UTC(),
		time.Now().UTC(),
	)

	// Save AC
	if err := acRepo.SaveAC(ctx, ac); err != nil {
		t.Fatalf("failed to save AC: %v", err)
	}

	// Get AC
	retrieved, err := acRepo.GetAC(ctx, "ac-1")
	if err != nil {
		t.Fatalf("failed to get AC: %v", err)
	}

	if retrieved.ID != ac.ID || retrieved.Description != ac.Description {
		t.Errorf("AC mismatch")
	}
	if retrieved.Status != entities.ACStatusNotStarted {
		t.Errorf("expected status not-started, got %s", retrieved.Status)
	}
}

func TestListAC(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	// Create multiple ACs
	for i := 1; i <= 3; i++ {
		id := "ac-" + string(rune(48+i))
		ac := entities.NewAcceptanceCriteriaEntity(id, "task-1", "AC "+string(rune(48+i)), entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
		acRepo.SaveAC(ctx, ac)
	}

	// List ACs
	acs, err := acRepo.ListAC(ctx, "task-1")
	if err != nil {
		t.Fatalf("failed to list ACs: %v", err)
	}

	if len(acs) != 3 {
		t.Errorf("expected 3 ACs, got %d", len(acs))
	}
}

func TestUpdateAC(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "Old Description", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	acRepo.SaveAC(ctx, ac)

	// Update AC
	ac.Description = "New Description"
	ac.Status = entities.ACStatusVerified
	ac.UpdatedAt = time.Now().UTC()

	if err := acRepo.UpdateAC(ctx, ac); err != nil {
		t.Fatalf("failed to update AC: %v", err)
	}

	// Verify update
	retrieved, _ := acRepo.GetAC(ctx, "ac-1")
	if retrieved.Description != "New Description" {
		t.Errorf("expected description to be updated, got %s", retrieved.Description)
	}
	if retrieved.Status != entities.ACStatusVerified {
		t.Errorf("expected status verified, got %s", retrieved.Status)
	}
}

func TestListFailedAC(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	// Create ACs with different statuses
	ac1 := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC 1", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	ac2 := entities.NewAcceptanceCriteriaEntity("ac-2", "task-1", "AC 2", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	ac3 := entities.NewAcceptanceCriteriaEntity("ac-3", "task-1", "AC 3", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())

	acRepo.SaveAC(ctx, ac1)
	acRepo.SaveAC(ctx, ac2)
	acRepo.SaveAC(ctx, ac3)

	// Mark ac-1 and ac-2 as failed, ac-3 as verified
	ac1.Status = entities.ACStatusFailed
	ac1.Notes = "Failed reason 1"
	acRepo.UpdateAC(ctx, ac1)

	ac2.Status = entities.ACStatusFailed
	ac2.Notes = "Failed reason 2"
	acRepo.UpdateAC(ctx, ac2)

	ac3.Status = entities.ACStatusVerified
	acRepo.UpdateAC(ctx, ac3)

	// List failed ACs (no filter)
	failedACs, err := acRepo.ListFailedAC(ctx, entities.ACFilters{})
	if err != nil {
		t.Fatalf("failed to list failed ACs: %v", err)
	}

	// Should return ac-1 and ac-2 only
	if len(failedACs) != 2 {
		t.Errorf("expected 2 failed ACs, got %d", len(failedACs))
	}

	// Verify statuses
	for _, ac := range failedACs {
		if ac.Status != entities.ACStatusFailed {
			t.Errorf("expected failed status, got %s", ac.Status)
		}
	}
}

func TestListFailedACWithTaskFilter(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	task2, _ := entities.NewTaskEntity("task-2", "track-1", "Task 2", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task1)
	taskRepo.SaveTask(ctx, task2)

	// Create failed ACs for both tasks
	ac1 := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC 1", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	ac2 := entities.NewAcceptanceCriteriaEntity("ac-2", "task-2", "AC 2", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())

	acRepo.SaveAC(ctx, ac1)
	acRepo.SaveAC(ctx, ac2)

	ac1.Status = entities.ACStatusFailed
	acRepo.UpdateAC(ctx, ac1)

	ac2.Status = entities.ACStatusFailed
	acRepo.UpdateAC(ctx, ac2)

	// Filter by task-1
	failedACs, err := acRepo.ListFailedAC(ctx, entities.ACFilters{TaskID: "task-1"})
	if err != nil {
		t.Fatalf("failed to list failed ACs: %v", err)
	}

	// Should return only ac-1
	if len(failedACs) != 1 {
		t.Errorf("expected 1 failed AC, got %d", len(failedACs))
	}
	if failedACs[0].ID != "ac-1" {
		t.Errorf("expected ac-1, got %s", failedACs[0].ID)
	}
}

func TestListFailedACWithIterationFilter(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	iterationRepo := persistence.NewSQLiteIterationRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	task2, _ := entities.NewTaskEntity("task-2", "track-1", "Task 2", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task1)
	taskRepo.SaveTask(ctx, task2)

	// Create iteration and add only task-1
	iter1, _ := entities.NewIterationEntity(1, "Sprint 1", "Goal", "", []string{}, "planned", 500, time.Time{}, time.Time{}, time.Now().UTC(), time.Now().UTC())
	iterationRepo.SaveIteration(ctx, iter1)
	iterationRepo.AddTaskToIteration(ctx, 1, "task-1")

	// Create failed ACs for both tasks
	ac1 := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC 1", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	ac2 := entities.NewAcceptanceCriteriaEntity("ac-2", "task-2", "AC 2", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())

	acRepo.SaveAC(ctx, ac1)
	acRepo.SaveAC(ctx, ac2)

	ac1.Status = entities.ACStatusFailed
	acRepo.UpdateAC(ctx, ac1)

	ac2.Status = entities.ACStatusFailed
	acRepo.UpdateAC(ctx, ac2)

	// Filter by iteration 1
	iterNum := 1
	failedACs, err := acRepo.ListFailedAC(ctx, entities.ACFilters{IterationNum: &iterNum})
	if err != nil {
		t.Fatalf("failed to list failed ACs: %v", err)
	}

	// Should return only ac-1 (task-1 is in iteration 1)
	if len(failedACs) != 1 {
		t.Errorf("expected 1 failed AC, got %d", len(failedACs))
	}
	if failedACs[0].ID != "ac-1" {
		t.Errorf("expected ac-1, got %s", failedACs[0].ID)
	}
}

func TestListFailedACWithTrackFilter(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track1)
	trackRepo.SaveTrack(ctx, track2)

	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	task2, _ := entities.NewTaskEntity("task-2", "track-2", "Task 2", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task1)
	taskRepo.SaveTask(ctx, task2)

	// Create failed ACs for both tasks
	ac1 := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC 1", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	ac2 := entities.NewAcceptanceCriteriaEntity("ac-2", "task-2", "AC 2", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())

	acRepo.SaveAC(ctx, ac1)
	acRepo.SaveAC(ctx, ac2)

	ac1.Status = entities.ACStatusFailed
	acRepo.UpdateAC(ctx, ac1)

	ac2.Status = entities.ACStatusFailed
	acRepo.UpdateAC(ctx, ac2)

	// Filter by track-1
	failedACs, err := acRepo.ListFailedAC(ctx, entities.ACFilters{TrackID: "track-1"})
	if err != nil {
		t.Fatalf("failed to list failed ACs: %v", err)
	}

	// Should return only ac-1 (task-1 is in track-1)
	if len(failedACs) != 1 {
		t.Errorf("expected 1 failed AC, got %d", len(failedACs))
	}
	if failedACs[0].ID != "ac-1" {
		t.Errorf("expected ac-1, got %s", failedACs[0].ID)
	}
}

func TestListACForIteration(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	iterationRepo := persistence.NewSQLiteIterationRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	task2, _ := entities.NewTaskEntity("task-2", "track-1", "Task 2", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task1)
	taskRepo.SaveTask(ctx, task2)

	// Create iteration and add task-1
	iter1, _ := entities.NewIterationEntity(1, "Sprint 1", "Goal", "", []string{}, "planned", 500, time.Time{}, time.Time{}, time.Now().UTC(), time.Now().UTC())
	iterationRepo.SaveIteration(ctx, iter1)
	iterationRepo.AddTaskToIteration(ctx, 1, "task-1")

	// Create ACs for both tasks
	ac1 := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC 1", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	ac2 := entities.NewAcceptanceCriteriaEntity("ac-2", "task-1", "AC 2", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	ac3 := entities.NewAcceptanceCriteriaEntity("ac-3", "task-2", "AC 3", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())

	acRepo.SaveAC(ctx, ac1)
	acRepo.SaveAC(ctx, ac2)
	acRepo.SaveAC(ctx, ac3)

	// List ACs for iteration 1
	acs, err := acRepo.ListACByIteration(ctx, 1)
	if err != nil {
		t.Fatalf("failed to list ACs for iteration: %v", err)
	}

	// Should return ac-1 and ac-2 (task-1 is in iteration 1)
	if len(acs) != 2 {
		t.Errorf("expected 2 ACs for iteration, got %d", len(acs))
	}

	// Verify correct ACs returned
	foundAC1 := false
	foundAC2 := false
	for _, ac := range acs {
		if ac.ID == "ac-1" {
			foundAC1 = true
		}
		if ac.ID == "ac-2" {
			foundAC2 = true
		}
	}

	if !foundAC1 || !foundAC2 {
		t.Errorf("expected ac-1 and ac-2, got %v", acs)
	}
}

func TestDeleteAC(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	acRepo := persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC 1", entities.VerificationTypeManual, "", time.Now().UTC(), time.Now().UTC())
	acRepo.SaveAC(ctx, ac)

	// Delete AC
	if err := acRepo.DeleteAC(ctx, "ac-1"); err != nil {
		t.Fatalf("failed to delete AC: %v", err)
	}

	// Verify deletion
	_, err := acRepo.GetAC(ctx, "ac-1")
	if !errors.Is(err, pluginsdk.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
