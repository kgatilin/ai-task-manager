package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/persistence"
)

// ============================================================================
// Repository Composite Tests
// ============================================================================

func TestNewSQLiteRepositoryComposite(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	logger := createTestLogger()
	composite := persistence.NewSQLiteRepositoryComposite(db, logger)

	if composite == nil {
		t.Fatal("expected non-nil composite repository")
	}
	if composite.DB != db {
		t.Error("expected DB to be set")
	}
	if composite.Roadmap == nil {
		t.Error("expected Roadmap repository to be set")
	}
	if composite.Track == nil {
		t.Error("expected Track repository to be set")
	}
	if composite.Task == nil {
		t.Error("expected Task repository to be set")
	}
	if composite.Iteration == nil {
		t.Error("expected Iteration repository to be set")
	}
	if composite.ADR == nil {
		t.Error("expected ADR repository to be set")
	}
	if composite.AC == nil {
		t.Error("expected AC repository to be set")
	}
	if composite.Aggregate == nil {
		t.Error("expected Aggregate repository to be set")
	}
}

func TestGetDB(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	if composite.GetDB() != db {
		t.Error("expected GetDB to return the database connection")
	}
}

// ============================================================================
// Roadmap Operation Delegation Tests
// ============================================================================

func TestComposite_SaveRoadmap_DelegatesToRoadmapRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())

	err := composite.SaveRoadmap(ctx, roadmap)
	if err != nil {
		t.Fatalf("SaveRoadmap failed: %v", err)
	}

	// Verify it was saved
	retrieved, err := composite.GetRoadmap(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("GetRoadmap failed: %v", err)
	}
	if retrieved.ID != "roadmap-1" {
		t.Errorf("expected roadmap-1, got %s", retrieved.ID)
	}
}

func TestComposite_GetActiveRoadmap_DelegatesToRoadmapRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	active, err := composite.GetActiveRoadmap(ctx)
	if err != nil {
		t.Fatalf("GetActiveRoadmap failed: %v", err)
	}
	if active.ID != "roadmap-1" {
		t.Errorf("expected roadmap-1, got %s", active.ID)
	}
}

func TestComposite_UpdateRoadmap_DelegatesToRoadmapRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	roadmap.Vision = "updated vision"
	err := composite.UpdateRoadmap(ctx, roadmap)
	if err != nil {
		t.Fatalf("UpdateRoadmap failed: %v", err)
	}

	retrieved, _ := composite.GetRoadmap(ctx, "roadmap-1")
	if retrieved.Vision != "updated vision" {
		t.Errorf("expected updated vision, got %s", retrieved.Vision)
	}
}

// ============================================================================
// Track Operation Delegation Tests
// ============================================================================

func TestComposite_SaveTrack_DelegatesToTrackRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	// Setup roadmap
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())

	err := composite.SaveTrack(ctx, track)
	if err != nil {
		t.Fatalf("SaveTrack failed: %v", err)
	}

	retrieved, err := composite.GetTrack(ctx, "track-1")
	if err != nil {
		t.Fatalf("GetTrack failed: %v", err)
	}
	if retrieved.ID != "track-1" {
		t.Errorf("expected track-1, got %s", retrieved.ID)
	}
}

func TestComposite_ListTracks_DelegatesToTrackRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track1)

	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "desc", "in-progress", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track2)

	tracks, err := composite.ListTracks(ctx, "roadmap-1", entities.TrackFilters{})
	if err != nil {
		t.Fatalf("ListTracks failed: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("expected 2 tracks, got %d", len(tracks))
	}
}

func TestComposite_UpdateTrack_DelegatesToTrackRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	track.Status = "in-progress"
	err := composite.UpdateTrack(ctx, track)
	if err != nil {
		t.Fatalf("UpdateTrack failed: %v", err)
	}

	retrieved, _ := composite.GetTrack(ctx, "track-1")
	if retrieved.Status != "in-progress" {
		t.Errorf("expected in-progress, got %s", retrieved.Status)
	}
}

func TestComposite_DeleteTrack_DelegatesToTrackRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	err := composite.DeleteTrack(ctx, "track-1")
	if err != nil {
		t.Fatalf("DeleteTrack failed: %v", err)
	}

	// Verify it was deleted
	_, err = composite.GetTrack(ctx, "track-1")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestComposite_TrackDependencies_DelegatesToTrackRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track1)

	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "desc", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track2)

	// Add dependency
	err := composite.AddTrackDependency(ctx, "track-2", "track-1")
	if err != nil {
		t.Fatalf("AddTrackDependency failed: %v", err)
	}

	// Get dependencies
	deps, err := composite.GetTrackDependencies(ctx, "track-2")
	if err != nil {
		t.Fatalf("GetTrackDependencies failed: %v", err)
	}
	if len(deps) != 1 || deps[0] != "track-1" {
		t.Errorf("expected [track-1], got %v", deps)
	}

	// Validate no cycles
	err = composite.ValidateNoCycles(ctx, "track-2")
	if err != nil {
		t.Fatalf("ValidateNoCycles failed: %v", err)
	}

	// Remove dependency
	err = composite.RemoveTrackDependency(ctx, "track-2", "track-1")
	if err != nil {
		t.Fatalf("RemoveTrackDependency failed: %v", err)
	}

	deps, _ = composite.GetTrackDependencies(ctx, "track-2")
	if len(deps) != 0 {
		t.Errorf("expected no dependencies, got %v", deps)
	}
}

func TestComposite_GetTrackWithTasks_DelegatesToTrackRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	retrieved, err := composite.GetTrackWithTasks(ctx, "track-1")
	if err != nil {
		t.Fatalf("GetTrackWithTasks failed: %v", err)
	}
	if retrieved.ID != "track-1" {
		t.Errorf("expected track-1, got %s", retrieved.ID)
	}
}

// ============================================================================
// Task Operation Delegation Tests
// ============================================================================

func TestComposite_SaveTask_DelegatesToTaskRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "desc", "todo", 100, "", time.Now().UTC(), time.Now().UTC())

	err := composite.SaveTask(ctx, task)
	if err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	retrieved, err := composite.GetTask(ctx, "task-1")
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if retrieved.ID != "task-1" {
		t.Errorf("expected task-1, got %s", retrieved.ID)
	}
}

func TestComposite_ListTasks_DelegatesToTaskRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task1)

	task2, _ := entities.NewTaskEntity("task-2", "track-1", "Task 2", "", "in-progress", 200, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task2)

	tasks, err := composite.ListTasks(ctx, entities.TaskFilters{TrackID: "track-1"})
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestComposite_UpdateTask_DelegatesToTaskRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	task.Status = "in-progress"
	err := composite.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	retrieved, _ := composite.GetTask(ctx, "task-1")
	if retrieved.Status != "in-progress" {
		t.Errorf("expected in-progress, got %s", retrieved.Status)
	}
}

func TestComposite_DeleteTask_DelegatesToTaskRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	err := composite.DeleteTask(ctx, "task-1")
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	_, err = composite.GetTask(ctx, "task-1")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestComposite_MoveTaskToTrack_DelegatesToTaskRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track1)

	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "desc", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track2)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	err := composite.MoveTaskToTrack(ctx, "task-1", "track-2")
	if err != nil {
		t.Fatalf("MoveTaskToTrack failed: %v", err)
	}

	retrieved, _ := composite.GetTask(ctx, "task-1")
	if retrieved.TrackID != "track-2" {
		t.Errorf("expected track-2, got %s", retrieved.TrackID)
	}
}

func TestComposite_GetBacklogTasks_DelegatesToTaskRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	backlog, err := composite.GetBacklogTasks(ctx)
	if err != nil {
		t.Fatalf("GetBacklogTasks failed: %v", err)
	}
	if len(backlog) != 1 {
		t.Errorf("expected 1 backlog task, got %d", len(backlog))
	}
}

// ============================================================================
// Iteration Operation Delegation Tests
// ============================================================================

func TestComposite_SaveIteration_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)

	err := composite.SaveIteration(ctx, iteration)
	if err != nil {
		t.Fatalf("SaveIteration failed: %v", err)
	}

	retrieved, err := composite.GetIteration(ctx, 1)
	if err != nil {
		t.Fatalf("GetIteration failed: %v", err)
	}
	if retrieved.Number != 1 {
		t.Errorf("expected 1, got %d", retrieved.Number)
	}
}

func TestComposite_ListIterations_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iter1, _ := entities.NewIterationEntity(1, "Iter 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iter1)

	iter2, _ := entities.NewIterationEntity(2, "Iter 2", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iter2)

	iterations, err := composite.ListIterations(ctx)
	if err != nil {
		t.Fatalf("ListIterations failed: %v", err)
	}
	if len(iterations) != 2 {
		t.Errorf("expected 2 iterations, got %d", len(iterations))
	}
}

func TestComposite_UpdateIteration_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)

	iteration.Name = "Updated Iteration"
	err := composite.UpdateIteration(ctx, iteration)
	if err != nil {
		t.Fatalf("UpdateIteration failed: %v", err)
	}

	retrieved, _ := composite.GetIteration(ctx, 1)
	if retrieved.Name != "Updated Iteration" {
		t.Errorf("expected Updated Iteration, got %s", retrieved.Name)
	}
}

func TestComposite_DeleteIteration_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)

	err := composite.DeleteIteration(ctx, 1)
	if err != nil {
		t.Fatalf("DeleteIteration failed: %v", err)
	}

	_, err = composite.GetIteration(ctx, 1)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestComposite_IterationTaskManagement_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)

	// Add task to iteration
	err := composite.AddTaskToIteration(ctx, 1, "task-1")
	if err != nil {
		t.Fatalf("AddTaskToIteration failed: %v", err)
	}

	// Get iteration tasks
	tasks, err := composite.GetIterationTasks(ctx, 1)
	if err != nil {
		t.Fatalf("GetIterationTasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}

	// Remove task from iteration
	err = composite.RemoveTaskFromIteration(ctx, 1, "task-1")
	if err != nil {
		t.Fatalf("RemoveTaskFromIteration failed: %v", err)
	}

	tasks, _ = composite.GetIterationTasks(ctx, 1)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestComposite_GetIterationTasksWithWarnings_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)

	tasks, warnings, err := composite.GetIterationTasksWithWarnings(ctx, 1)
	if err != nil {
		t.Fatalf("GetIterationTasksWithWarnings failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
	if len(warnings) != 0 {
		t.Errorf("expected 0 warnings, got %d", len(warnings))
	}
}

func TestComposite_StartAndCompleteIteration_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)

	// Start iteration
	err := composite.StartIteration(ctx, 1)
	if err != nil {
		t.Fatalf("StartIteration failed: %v", err)
	}

	retrieved, _ := composite.GetIteration(ctx, 1)
	if retrieved.Status != "current" {
		t.Errorf("expected current, got %s", retrieved.Status)
	}

	// Complete iteration
	err = composite.CompleteIteration(ctx, 1)
	if err != nil {
		t.Fatalf("CompleteIteration failed: %v", err)
	}

	retrieved, _ = composite.GetIteration(ctx, 1)
	if retrieved.Status != "complete" {
		t.Errorf("expected complete, got %s", retrieved.Status)
	}
}

func TestComposite_GetCurrentIteration_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)
	composite.StartIteration(ctx, 1)

	current, err := composite.GetCurrentIteration(ctx)
	if err != nil {
		t.Fatalf("GetCurrentIteration failed: %v", err)
	}
	if current.Number != 1 {
		t.Errorf("expected 1, got %d", current.Number)
	}
}

func TestComposite_GetIterationByNumber_DelegatesToIterationRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)

	retrieved, err := composite.GetIterationByNumber(ctx, 1)
	if err != nil {
		t.Fatalf("GetIterationByNumber failed: %v", err)
	}
	if retrieved.Number != 1 {
		t.Errorf("expected 1, got %d", retrieved.Number)
	}
}

func TestComposite_GetIterationsForTask_DelegatesToTaskRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)
	composite.AddTaskToIteration(ctx, 1, "task-1")

	iterations, err := composite.GetIterationsForTask(ctx, "task-1")
	if err != nil {
		t.Fatalf("GetIterationsForTask failed: %v", err)
	}
	if len(iterations) != 1 {
		t.Errorf("expected 1 iteration, got %d", len(iterations))
	}
}

// ============================================================================
// ADR Operation Delegation Tests
// ============================================================================

func TestComposite_ADROperations_DelegateToADRRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	now := time.Now().UTC()
	adr, _ := entities.NewADREntity("adr-1", "track-1", "ADR", "proposed", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)

	err := composite.SaveADR(ctx, adr)
	if err != nil {
		t.Fatalf("SaveADR failed: %v", err)
	}

	retrieved, err := composite.GetADR(ctx, "adr-1")
	if err != nil {
		t.Fatalf("GetADR failed: %v", err)
	}
	if retrieved.ID != "adr-1" {
		t.Errorf("expected adr-1, got %s", retrieved.ID)
	}

	adrs, err := composite.ListADRs(ctx, nil)
	if err != nil {
		t.Fatalf("ListADRs failed: %v", err)
	}
	if len(adrs) != 1 {
		t.Errorf("expected 1 ADR, got %d", len(adrs))
	}

	adr.Title = "Updated ADR"
	err = composite.UpdateADR(ctx, adr)
	if err != nil {
		t.Fatalf("UpdateADR failed: %v", err)
	}

	trackADRs, err := composite.GetADRsByTrack(ctx, "track-1")
	if err != nil {
		t.Fatalf("GetADRsByTrack failed: %v", err)
	}
	if len(trackADRs) != 1 {
		t.Errorf("expected 1 ADR for track, got %d", len(trackADRs))
	}
}

func TestComposite_SupersedeAndDeprecateADR_DelegateToADRRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("adr-1", "track-1", "ADR 1", "accepted", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	composite.SaveADR(ctx, adr1)

	adr2, _ := entities.NewADREntity("adr-2", "track-1", "ADR 2", "proposed", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	composite.SaveADR(ctx, adr2)

	err := composite.SupersedeADR(ctx, "adr-1", "adr-2")
	if err != nil {
		t.Fatalf("SupersedeADR failed: %v", err)
	}

	err = composite.DeprecateADR(ctx, "adr-2")
	if err != nil {
		t.Fatalf("DeprecateADR failed: %v", err)
	}
}

// ============================================================================
// AC Operation Delegation Tests
// ============================================================================

func TestComposite_ACOperations_DelegateToACRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "instructions", now, now)

	err := composite.SaveAC(ctx, ac)
	if err != nil {
		t.Fatalf("SaveAC failed: %v", err)
	}

	retrieved, err := composite.GetAC(ctx, "ac-1")
	if err != nil {
		t.Fatalf("GetAC failed: %v", err)
	}
	if retrieved.ID != "ac-1" {
		t.Errorf("expected ac-1, got %s", retrieved.ID)
	}

	acs, err := composite.ListAC(ctx, "task-1")
	if err != nil {
		t.Fatalf("ListAC failed: %v", err)
	}
	if len(acs) != 1 {
		t.Errorf("expected 1 AC, got %d", len(acs))
	}

	ac.Description = "Updated AC"
	err = composite.UpdateAC(ctx, ac)
	if err != nil {
		t.Fatalf("UpdateAC failed: %v", err)
	}

	err = composite.DeleteAC(ctx, "ac-1")
	if err != nil {
		t.Fatalf("DeleteAC failed: %v", err)
	}
}

func TestComposite_ListACByTask_DelegateToACRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "instructions", now, now)
	composite.SaveAC(ctx, ac)

	acs, err := composite.ListACByTask(ctx, "task-1")
	if err != nil {
		t.Fatalf("ListACByTask failed: %v", err)
	}
	if len(acs) != 1 {
		t.Errorf("expected 1 AC, got %d", len(acs))
	}
}

func TestComposite_ListACByTrack_CrossEntityQuery(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "instructions", now, now)
	composite.SaveAC(ctx, ac)

	acs, err := composite.ListACByTrack(ctx, "track-1")
	if err != nil {
		t.Fatalf("ListACByTrack failed: %v", err)
	}
	if len(acs) != 1 {
		t.Errorf("expected 1 AC, got %d", len(acs))
	}
}

func TestComposite_ListACByIteration_DelegateToACRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iteration 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	composite.SaveIteration(ctx, iteration)
	composite.AddTaskToIteration(ctx, 1, "task-1")

	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "instructions", now, now)
	composite.SaveAC(ctx, ac)

	acs, err := composite.ListACByIteration(ctx, 1)
	if err != nil {
		t.Fatalf("ListACByIteration failed: %v", err)
	}
	if len(acs) != 1 {
		t.Errorf("expected 1 AC, got %d", len(acs))
	}
}

func TestComposite_ListFailedAC_DelegateToACRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	composite.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "instructions", now, now)
	ac.Status = entities.ACStatusFailed
	composite.SaveAC(ctx, ac)

	failedACs, err := composite.ListFailedAC(ctx, entities.ACFilters{})
	if err != nil {
		t.Fatalf("ListFailedAC failed: %v", err)
	}
	if len(failedACs) != 1 {
		t.Errorf("expected 1 failed AC, got %d", len(failedACs))
	}
}

// ============================================================================
// Aggregate Query Delegation Tests
// ============================================================================

func TestComposite_GetRoadmapWithTracks_DelegatesToAggregateRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	composite.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "desc", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	composite.SaveTrack(ctx, track)

	result, err := composite.GetRoadmapWithTracks(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("GetRoadmapWithTracks failed: %v", err)
	}
	if result.ID != "roadmap-1" {
		t.Errorf("expected roadmap-1, got %s", result.ID)
	}
}

// ============================================================================
// Project Metadata Delegation Tests
// ============================================================================

func TestComposite_ProjectMetadata_DelegatesToAggregateRepo(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	composite := persistence.NewSQLiteRepositoryComposite(db, createTestLogger())

	err := composite.SetProjectMetadata(ctx, "test_key", "test_value")
	if err != nil {
		t.Fatalf("SetProjectMetadata failed: %v", err)
	}

	value, err := composite.GetProjectMetadata(ctx, "test_key")
	if err != nil {
		t.Fatalf("GetProjectMetadata failed: %v", err)
	}
	if value != "test_value" {
		t.Errorf("expected test_value, got %s", value)
	}

	code := composite.GetProjectCode(ctx)
	if code != "DW" {
		t.Errorf("expected DW, got %s", code)
	}

	seq, err := composite.GetNextSequenceNumber(ctx, "task")
	if err != nil {
		t.Fatalf("GetNextSequenceNumber failed: %v", err)
	}
	if seq != 1 {
		t.Errorf("expected 1, got %d", seq)
	}
}

// ============================================================================
// Backward Compatibility Tests
// ============================================================================

func TestNewSQLiteRoadmapRepository_ReturnsComposite(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	logger := createTestLogger()
	repo := persistence.NewSQLiteRoadmapRepository(db, logger)

	if repo == nil {
		t.Fatal("expected non-nil repository")
	}

	// Verify it's a composite by testing it has all methods
	ctx := context.Background()
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())

	err := repo.SaveRoadmap(ctx, roadmap)
	if err != nil {
		t.Fatalf("SaveRoadmap failed: %v", err)
	}

	retrieved, err := repo.GetRoadmap(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("GetRoadmap failed: %v", err)
	}
	if retrieved.ID != "roadmap-1" {
		t.Errorf("expected roadmap-1, got %s", retrieved.ID)
	}
}
