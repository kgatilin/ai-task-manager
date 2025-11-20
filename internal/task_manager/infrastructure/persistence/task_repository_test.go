package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/persistence"
)

// ============================================================================
// Task Tests
// ============================================================================

func TestSaveAndGetTask(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create and save task
	task, _ := entities.NewTaskEntity("task-1", "track-1", "Implement feature", "Do something", "todo", 200, "feat/impl", time.Now().UTC(), time.Now().UTC())

	if err := taskRepo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Get task
	retrieved, err := taskRepo.GetTask(ctx, "task-1")
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}

	if retrieved.ID != task.ID || retrieved.Title != task.Title {
		t.Errorf("task mismatch")
	}
	if retrieved.Branch != "feat/impl" {
		t.Errorf("expected branch feat/impl, got %s", retrieved.Branch)
	}
}

func TestListTasks(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create multiple tasks
	for i := 1; i <= 3; i++ {
		id := "task-" + string(rune(48+i))
		task, _ := entities.NewTaskEntity(id, "track-1", "Task "+string(rune(48+i)), "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
		taskRepo.SaveTask(ctx, task)
	}

	// List tasks
	tasks, err := taskRepo.ListTasks(ctx, entities.TaskFilters{})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestListTasksWithFilters(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create tasks with different statuses
	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	task2, _ := entities.NewTaskEntity("task-2", "track-1", "Task 2", "", "done", 200, "", time.Now().UTC(), time.Now().UTC())

	taskRepo.SaveTask(ctx, task1)
	taskRepo.SaveTask(ctx, task2)

	// Filter by status
	tasks, err := taskRepo.ListTasks(ctx, entities.TaskFilters{Status: []string{"done"}})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}

	if len(tasks) != 1 || tasks[0].Status != "done" {
		t.Errorf("expected 1 done task, got %d", len(tasks))
	}
}

func TestMoveTaskToTrack(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())

	trackRepo.SaveTrack(ctx, track1)
	trackRepo.SaveTrack(ctx, track2)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	// Move task to track-2
	if err := taskRepo.MoveTaskToTrack(ctx, "task-1", "track-2"); err != nil {
		t.Fatalf("failed to move task: %v", err)
	}

	// Verify move
	updated, _ := taskRepo.GetTask(ctx, "task-1")
	if updated.TrackID != "track-2" {
		t.Errorf("expected track-2, got %s", updated.TrackID)
	}
}

func TestGetBacklogTasks(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	iterationRepo := persistence.NewSQLiteIterationRepository(db, createTestLogger(), persistence.NewSQLiteAcceptanceCriteriaRepository(db, createTestLogger()))
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create tasks:
	// - task-1: in backlog (todo, not in any iteration)
	// - task-2: in backlog (in-progress, not in any iteration)
	// - task-3: NOT in backlog (done, not in any iteration)
	// - task-4: NOT in backlog (todo, but in an iteration)
	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	task2, _ := entities.NewTaskEntity("task-2", "track-1", "Task 2", "", "in-progress", 200, "", time.Now().UTC(), time.Now().UTC())
	task3, _ := entities.NewTaskEntity("task-3", "track-1", "Task 3", "", "done", 200, "", time.Now().UTC(), time.Now().UTC())
	task4, _ := entities.NewTaskEntity("task-4", "track-1", "Task 4", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())

	taskRepo.SaveTask(ctx, task1)
	taskRepo.SaveTask(ctx, task2)
	taskRepo.SaveTask(ctx, task3)
	taskRepo.SaveTask(ctx, task4)

	// Add task-4 to an iteration
	iter1, _ := entities.NewIterationEntity(1, "Sprint 1", "Goal", "", []string{}, "planned", 500, time.Time{}, time.Time{}, time.Now().UTC(), time.Now().UTC())
	iterationRepo.SaveIteration(ctx, iter1)
	iterationRepo.AddTaskToIteration(ctx, 1, "task-4")

	// Get backlog tasks
	backlog, err := taskRepo.GetBacklogTasks(ctx)
	if err != nil {
		t.Fatalf("failed to get backlog tasks: %v", err)
	}

	// Should return task-1 and task-2 only
	if len(backlog) != 2 {
		t.Errorf("expected 2 backlog tasks, got %d", len(backlog))
	}

	// Check that we got the right tasks (task-1 and task-2)
	foundTask1 := false
	foundTask2 := false
	for _, task := range backlog {
		if task.ID == "task-1" {
			foundTask1 = true
		}
		if task.ID == "task-2" {
			foundTask2 = true
		}
	}

	if !foundTask1 || !foundTask2 {
		t.Errorf("expected task-1 and task-2 in backlog, got %v", backlog)
	}
}

func TestGetBacklogTasksEmpty(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create only done tasks
	task1, _ := entities.NewTaskEntity("task-1", "track-1", "Task 1", "", "done", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task1)

	// Get backlog tasks
	backlog, err := taskRepo.GetBacklogTasks(ctx)
	if err != nil {
		t.Fatalf("failed to get backlog tasks: %v", err)
	}

	// Should return empty slice
	if len(backlog) != 0 {
		t.Errorf("expected 0 backlog tasks, got %d", len(backlog))
	}
}

func TestUpdateTask(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	// Update task
	task.Title = "Updated Task"
	task.Status = "done"
	task.UpdatedAt = time.Now().UTC()

	if err := taskRepo.UpdateTask(ctx, task); err != nil {
		t.Fatalf("failed to update task: %v", err)
	}

	// Verify update
	retrieved, _ := taskRepo.GetTask(ctx, "task-1")
	if retrieved.Title != "Updated Task" {
		t.Errorf("expected title to be updated, got %s", retrieved.Title)
	}
	if retrieved.Status != "done" {
		t.Errorf("expected status done, got %s", retrieved.Status)
	}
}

func TestDeleteTask(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	taskRepo := persistence.NewSQLiteTaskRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 200, "", time.Now().UTC(), time.Now().UTC())
	taskRepo.SaveTask(ctx, task)

	// Delete task
	if err := taskRepo.DeleteTask(ctx, "task-1"); err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}

	// Verify deletion
	_, err := taskRepo.GetTask(ctx, "task-1")
	if !errors.Is(err, tmerrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
