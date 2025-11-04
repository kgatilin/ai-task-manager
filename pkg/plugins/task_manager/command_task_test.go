package task_manager_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"
)

// ============================================================================
// TaskCreateCommand Tests
// ============================================================================

func TestTaskCreateCommand_Success(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)
	ctx := context.Background()

	// Setup: Create roadmap and track using commands
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Test vision",
		"--success-criteria", "Test criteria",
	})
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	trackCmd := &task_manager.TrackCreateCommand{Plugin: plugin}
	trackCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = trackCmd.Execute(ctx, trackCtx, []string{
		"--title", "Test Track",
		"--description", "Test description",
		"--rank", "200",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Extract track ID from output (format: "ID:          <ID>")
	trackOutput := trackCtx.stdout.String()
	trackIDPrefix := "ID:"
	trackIDStart := strings.Index(trackOutput, trackIDPrefix)
	if trackIDStart == -1 {
		t.Fatalf("failed to find track ID in output: %s", trackOutput)
	}
	trackIDStart += len(trackIDPrefix)
	trackIDEnd := strings.Index(trackOutput[trackIDStart:], "\n")
	if trackIDEnd == -1 {
		trackIDEnd = len(trackOutput)
	} else {
		trackIDEnd += trackIDStart
	}
	trackID := strings.TrimSpace(trackOutput[trackIDStart:trackIDEnd])

	// Execute command
	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", trackID,
		"--title", "Implement feature",
		"--description", "Add new feature",
		"--rank", "200",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify output
	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task created successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
	if !strings.Contains(output, "Implement feature") {
		t.Errorf("expected title in output, got: %s", output)
	}

	// Verify task was saved
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})

	tasks, err := repo.ListTasks(ctx, task_manager.TaskFilters{TrackID: trackID})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
	if len(tasks) > 0 {
		if tasks[0].Title != "Implement feature" {
			t.Errorf("expected title 'Implement feature', got '%s'", tasks[0].Title)
		}
		if tasks[0].Rank != 200 {
			t.Errorf("expected rank 200, got %d", tasks[0].Rank)
		}
		if tasks[0].Status != "todo" {
			t.Errorf("expected status 'todo', got '%s'", tasks[0].Status)
		}
	}
}

func TestTaskCreateCommand_MissingTrack(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--title", "Implement feature",
		"--description", "Add new feature",
	})
	if err == nil {
		t.Errorf("expected error for missing track flag")
	}
	if !strings.Contains(err.Error(), "--track") {
		t.Errorf("expected error about --track, got: %v", err)
	}
}

func TestTaskCreateCommand_MissingTitle(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", "track-test",
		"--description", "Add new feature",
	})
	if err == nil {
		t.Errorf("expected error for missing title flag")
	}
	if !strings.Contains(err.Error(), "--title") {
		t.Errorf("expected error about --title, got: %v", err)
	}
}

func TestTaskCreateCommand_TrackNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", "nonexistent-track",
		"--title", "Implement feature",
	})
	if err == nil {
		t.Errorf("expected error for nonexistent track")
	}
	if !strings.Contains(err.Error(), "track not found") {
		t.Errorf("expected 'track not found' error, got: %v", err)
	}
}

// ============================================================================
// TaskListCommand Tests
// ============================================================================

func TestTaskListCommand_NoTasks(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskListCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "No tasks found") {
		t.Errorf("expected 'No tasks found', got: %s", output)
	}
}

func TestTaskListCommand_ListAllTasks(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup: Create roadmap and track using direct database
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	track, err := task_manager.NewTrackEntity(
		"DW-track-1",
		"roadmap-test",
		"Test Track",
		"Test description",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	// Create multiple tasks
	for i := 0; i < 3; i++ {
		taskID := fmt.Sprintf("DW-task-%d", i+1)
		task := task_manager.NewTaskEntity(
			taskID,
			"DW-track-1",
			"Task "+string(rune(i+49)),
			"",
			"todo",
			300,
			"",
			time.Now().UTC(),
			time.Now().UTC(),
		)
		if err := repo.SaveTask(ctx, task); err != nil {
			t.Fatalf("failed to save task: %v", err)
		}
	}

	// Execute command
	cmd := &task_manager.TaskListCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Total: 3 task(s)") {
		t.Errorf("expected '3 task(s)', got: %s", output)
	}
	if !strings.Contains(output, "Task") {
		t.Errorf("expected task titles in output, got: %s", output)
	}
}

func TestTaskListCommand_FilterByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	track, err := task_manager.NewTrackEntity(
		"track-test",
		"roadmap-test",
		"Test Track",
		"Test description",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	// Create tasks with different statuses
	statuses := []string{"todo", "in-progress", "done"}
	for i, status := range statuses {
		taskID := fmt.Sprintf("DEF-task-%d", i+1)
		task := task_manager.NewTaskEntity(
			taskID,
			"track-test",
			"Task "+string(rune(i+49)),
			"",
			status,
			300,
			"",
			time.Now().UTC(),
			time.Now().UTC(),
		)
		if err := repo.SaveTask(ctx, task); err != nil {
			t.Fatalf("failed to save task: %v", err)
		}
	}

	// Execute command with status filter
	cmd := &task_manager.TaskListCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"--status", "done"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Total: 1 task(s)") {
		t.Errorf("expected '1 task(s)', got: %s", output)
	}
	if !strings.Contains(output, "done") {
		t.Errorf("expected 'done' status in output, got: %s", output)
	}
}

func TestTaskListCommand_FilterByTrack(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	// Create two tracks
	for i := 0; i < 2; i++ {
		track, err := task_manager.NewTrackEntity(
			"track-test-"+string(rune(i+49)),
			"roadmap-test",
			"Track "+string(rune(i+49)),
			"",
			"not-started",
			200,
			[]string{},
			time.Now().UTC(),
			time.Now().UTC(),
		)
		if err != nil {
			t.Fatalf("failed to create track: %v", err)
		}
		if err := repo.SaveTrack(ctx, track); err != nil {
			t.Fatalf("failed to save track: %v", err)
		}
	}

	// Create tasks in different tracks
	for i := 0; i < 2; i++ {
		trackID := "track-test-" + string(rune(i+49))
		taskID := fmt.Sprintf("DEF-task-%d", i+1)
		task := task_manager.NewTaskEntity(
			taskID,
			trackID,
			"Task "+string(rune(i+49)),
			"",
			"todo",
			300,
			"",
			time.Now().UTC(),
			time.Now().UTC(),
		)
		if err := repo.SaveTask(ctx, task); err != nil {
			t.Fatalf("failed to save task: %v", err)
		}
	}

	// Execute command with track filter
	cmd := &task_manager.TaskListCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"--track", "track-test-1"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Total: 1 task(s)") {
		t.Errorf("expected '1 task(s)', got: %s", output)
	}
}

// ============================================================================
// TaskShowCommand Tests
// ============================================================================

func TestTaskShowCommand_Success(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	track, err := task_manager.NewTrackEntity(
		"track-test",
		"roadmap-test",
		"Test Track",
		"Test description",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	task := task_manager.NewTaskEntity(
		"DEF-task-1",
		"track-test",
		"Test Task",
		"Test description",
		"in-progress",
		200,
		"feat/test",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Execute command
	cmd := &task_manager.TaskShowCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"DEF-task-1"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Test Task") {
		t.Errorf("expected task title, got: %s", output)
	}
	if !strings.Contains(output, "in-progress") {
		t.Errorf("expected status, got: %s", output)
	}
	if !strings.Contains(output, "feat/test") {
		t.Errorf("expected branch, got: %s", output)
	}
}

func TestTaskShowCommand_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskShowCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"nonexistent-task"})
	if err == nil {
		t.Errorf("expected error for nonexistent task")
	}
	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("expected 'task not found' error, got: %v", err)
	}
}

func TestTaskShowCommand_WithIterations(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	// Create roadmap
	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	// Create track
	track, err := task_manager.NewTrackEntity(
		"track-test",
		"roadmap-test",
		"Test Track",
		"Test description",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	// Create task
	task := task_manager.NewTaskEntity(
		"DEF-task-1",
		"track-test",
		"Test Task",
		"Test description",
		"in-progress",
		200,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Create iterations
	now := time.Now().UTC()
	var zeroTime time.Time
	iter1, err := task_manager.NewIterationEntity(
		1,
		"Sprint 1",
		"First sprint",
		"Deliverable 1",
		[]string{"DEF-task-1"},
		"current",
		100,
		now,
		zeroTime,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("failed to create iteration 1: %v", err)
	}
	if err := repo.SaveIteration(ctx, iter1); err != nil {
		t.Fatalf("failed to save iteration 1: %v", err)
	}

	iter2, err := task_manager.NewIterationEntity(
		2,
		"Sprint 2",
		"Second sprint",
		"Deliverable 2",
		[]string{"DEF-task-1"},
		"planned",
		200,
		zeroTime,
		zeroTime,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("failed to create iteration 2: %v", err)
	}
	if err := repo.SaveIteration(ctx, iter2); err != nil {
		t.Fatalf("failed to save iteration 2: %v", err)
	}

	// Execute command
	cmd := &task_manager.TaskShowCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"DEF-task-1"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify output contains iterations section
	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Iterations:") {
		t.Errorf("expected 'Iterations:' section, got: %s", output)
	}
	if !strings.Contains(output, "Sprint 1") {
		t.Errorf("expected 'Sprint 1' in output, got: %s", output)
	}
	if !strings.Contains(output, "Sprint 2") {
		t.Errorf("expected 'Sprint 2' in output, got: %s", output)
	}
	if !strings.Contains(output, "current") {
		t.Errorf("expected 'current' status in output, got: %s", output)
	}
	if !strings.Contains(output, "planned") {
		t.Errorf("expected 'planned' status in output, got: %s", output)
	}
}

func TestTaskShowCommand_NoIterations(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	// Create roadmap
	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	// Create track
	track, err := task_manager.NewTrackEntity(
		"track-test",
		"roadmap-test",
		"Test Track",
		"Test description",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	// Create task (not assigned to any iteration)
	task := task_manager.NewTaskEntity(
		"DEF-task-1",
		"track-test",
		"Test Task",
		"Test description",
		"todo",
		200,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Execute command
	cmd := &task_manager.TaskShowCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"DEF-task-1"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify output contains "not assigned to any iteration" message
	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Iterations:") {
		t.Errorf("expected 'Iterations:' section, got: %s", output)
	}
	if !strings.Contains(output, "Not assigned to any iteration") {
		t.Errorf("expected 'Not assigned to any iteration' message, got: %s", output)
	}
}

// ============================================================================
// TaskUpdateCommand Tests
// ============================================================================

func TestTaskUpdateCommand_UpdateStatus(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	track, err := task_manager.NewTrackEntity(
		"track-test",
		"roadmap-test",
		"Test Track",
		"",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	task := task_manager.NewTaskEntity(
		"DEF-task-1",
		"track-test",
		"Test Task",
		"",
		"todo",
		300,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Execute command
	cmd := &task_manager.TaskUpdateCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"DEF-task-1", "--status", "in-progress"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify update
	updated, err := repo.GetTask(ctx, "DEF-task-1")
	if err != nil {
		t.Fatalf("failed to get updated task: %v", err)
	}
	if updated.Status != "in-progress" {
		t.Errorf("expected status 'in-progress', got '%s'", updated.Status)
	}
}

func TestTaskUpdateCommand_NoUpdates(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskUpdateCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"task-123"})
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
	if !strings.Contains(err.Error(), "at least one flag") {
		t.Errorf("expected error about flags, got: %v", err)
	}
}

func TestTaskUpdateCommand_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskUpdateCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"nonexistent-task", "--status", "done"})
	if err == nil {
		t.Errorf("expected error for nonexistent task")
	}
	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("expected 'task not found' error, got: %v", err)
	}
}

// ============================================================================
// TaskDeleteCommand Tests
// ============================================================================

func TestTaskDeleteCommand_WithForce(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	track, err := task_manager.NewTrackEntity(
		"track-test",
		"roadmap-test",
		"Test Track",
		"",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	task := task_manager.NewTaskEntity(
		"DEF-task-1",
		"track-test",
		"Test Task",
		"",
		"todo",
		300,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Execute command with --force
	cmd := &task_manager.TaskDeleteCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"DEF-task-1", "--force"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task deleted successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify task was deleted
	_, err = repo.GetTask(ctx, "DEF-task-1")
	if err == nil {
		t.Errorf("expected task to be deleted")
	}
}

func TestTaskDeleteCommand_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskDeleteCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"nonexistent-task", "--force"})
	if err == nil {
		t.Errorf("expected error for nonexistent task")
	}
	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("expected 'task not found' error, got: %v", err)
	}
}

// ============================================================================
// TaskMoveCommand Tests
// ============================================================================

func TestTaskMoveCommand_Success(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	// Create two tracks
	for i := 0; i < 2; i++ {
		track, err := task_manager.NewTrackEntity(
			"track-test-"+string(rune(i+49)),
			"roadmap-test",
			"Track "+string(rune(i+49)),
			"",
			"not-started",
			200,
			[]string{},
			time.Now().UTC(),
			time.Now().UTC(),
		)
		if err != nil {
			t.Fatalf("failed to create track: %v", err)
		}
		if err := repo.SaveTrack(ctx, track); err != nil {
			t.Fatalf("failed to save track: %v", err)
		}
	}

	// Create task in track 1
	task := task_manager.NewTaskEntity(
		"DEF-task-1",
		"track-test-1",
		"Test Task",
		"",
		"todo",
		300,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Execute command
	cmd := &task_manager.TaskMoveCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"DEF-task-1", "--track", "track-test-2"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task moved successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify move
	moved, err := repo.GetTask(ctx, "DEF-task-1")
	if err != nil {
		t.Fatalf("failed to get moved task: %v", err)
	}
	if moved.TrackID != "track-test-2" {
		t.Errorf("expected task to be in 'track-test-2', got '%s'", moved.TrackID)
	}
}

func TestTaskMoveCommand_NewTrackNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project (getProjectDB created the directory, but we need to set it as active)
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	track, err := task_manager.NewTrackEntity(
		"track-test",
		"roadmap-test",
		"Test Track",
		"",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	task := task_manager.NewTaskEntity(
		"DEF-task-1",
		"track-test",
		"Test Task",
		"",
		"todo",
		300,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task); err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Execute command with nonexistent target track
	cmd := &task_manager.TaskMoveCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"DEF-task-1", "--track", "nonexistent-track"})
	if err == nil {
		t.Errorf("expected error for nonexistent track")
	}
	if !strings.Contains(err.Error(), "track not found") {
		t.Errorf("expected 'track not found' error, got: %v", err)
	}
}

func TestTaskMoveCommand_TaskNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := createRoadmapTestDB(t)
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPluginWithDatabase(
		&stubLogger{},
		tmpDir,
		db,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	cmd := &task_manager.TaskMoveCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{"nonexistent-task", "--track", "some-track"})
	if err == nil {
		t.Errorf("expected error for nonexistent task")
	}
	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("expected 'task not found' error, got: %v", err)
	}
}

// ============================================================================
// TaskBacklogCommand Tests
// ============================================================================

func TestTaskBacklogCommand_ShowUnassignedTasks(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup: Create roadmap, track, tasks, and iteration
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	track, err := task_manager.NewTrackEntity(
		"DW-track-1",
		"roadmap-test",
		"Test Track",
		"Test description",
		"not-started",
		200,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
	if err := repo.SaveTrack(ctx, track); err != nil {
		t.Fatalf("failed to save track: %v", err)
	}

	// Create 3 tasks: 1 in iteration, 1 done (not in iteration), 1 backlog
	task1 := task_manager.NewTaskEntity(
		"DW-task-1",
		"DW-track-1",
		"Task in iteration",
		"",
		"in-progress",
		300,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task1); err != nil {
		t.Fatalf("failed to save task1: %v", err)
	}

	task2 := task_manager.NewTaskEntity(
		"DW-task-2",
		"DW-track-1",
		"Task done",
		"",
		"done",
		300,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task2); err != nil {
		t.Fatalf("failed to save task2: %v", err)
	}

	task3 := task_manager.NewTaskEntity(
		"DW-task-3",
		"DW-track-1",
		"Task in backlog",
		"",
		"todo",
		300,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err := repo.SaveTask(ctx, task3); err != nil {
		t.Fatalf("failed to save task3: %v", err)
	}

	// Create iteration and add task1 to it
	iteration, err := task_manager.NewIterationEntity(
		1,
		"Sprint 1",
		"Test goal",
		"Test deliverable",
		[]string{},
		"planned",
		500,
		time.Time{},
		time.Time{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}
	if err := repo.SaveIteration(ctx, iteration); err != nil {
		t.Fatalf("failed to save iteration: %v", err)
	}
	if err := repo.AddTaskToIteration(ctx, 1, "DW-task-1"); err != nil {
		t.Fatalf("failed to add task to iteration: %v", err)
	}

	// Execute backlog command
	cmd := &task_manager.TaskBacklogCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()

	// Should show only task3 (not in iteration, not done)
	if !strings.Contains(output, "DW-task-3") {
		t.Errorf("expected backlog to contain DW-task-3, got: %s", output)
	}
	if !strings.Contains(output, "Task in backlog") {
		t.Errorf("expected backlog to contain 'Task in backlog', got: %s", output)
	}

	// Should NOT show task1 (in iteration) or task2 (done)
	if strings.Contains(output, "DW-task-1") {
		t.Errorf("expected backlog to NOT contain DW-task-1 (in iteration), got: %s", output)
	}
	if strings.Contains(output, "DW-task-2") {
		t.Errorf("expected backlog to NOT contain DW-task-2 (done), got: %s", output)
	}

	// Should show total count
	if !strings.Contains(output, "Total: 1 backlog task(s)") {
		t.Errorf("expected total count of 1, got: %s", output)
	}
}

func TestTaskBacklogCommand_EmptyBacklog(t *testing.T) {
	tmpDir := t.TempDir()
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()

	plugin, err := task_manager.NewTaskManagerPlugin(
		&stubLogger{},
		tmpDir,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Set active project
	if err := os.WriteFile(filepath.Join(tmpDir, ".darwinflow", "active-project.txt"), []byte("default"), 0644); err != nil {
		t.Fatalf("failed to set active project: %v", err)
	}

	// Setup: Create roadmap only (no tasks)
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})
	ctx := context.Background()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-test",
		"Test vision",
		"Test criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	// Execute backlog command
	cmd := &task_manager.TaskBacklogCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "No backlog tasks found") {
		t.Errorf("expected 'No backlog tasks found', got: %s", output)
	}
}

// ============================================================================
// TaskCreateCommand - Track Status Validation Tests
// ============================================================================

func TestTaskCreateCommand_TrackCompleted(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)
	ctx := context.Background()

	// Setup: Create roadmap and completed track
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Test vision",
		"--success-criteria", "Test criteria",
	})
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	trackCmd := &task_manager.TrackCreateCommand{Plugin: plugin}
	trackCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = trackCmd.Execute(ctx, trackCtx, []string{
		"--title", "Completed Track",
		"--description", "A track that is complete",
		"--rank", "200",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Extract track ID
	trackOutput := trackCtx.stdout.String()
	trackIDPrefix := "ID:"
	trackIDStart := strings.Index(trackOutput, trackIDPrefix)
	if trackIDStart == -1 {
		t.Fatalf("failed to find track ID in output: %s", trackOutput)
	}
	trackIDStart += len(trackIDPrefix)
	trackIDEnd := strings.Index(trackOutput[trackIDStart:], "\n")
	if trackIDEnd == -1 {
		trackIDEnd = len(trackOutput)
	} else {
		trackIDEnd += trackIDStart
	}
	trackID := strings.TrimSpace(trackOutput[trackIDStart:trackIDEnd])

	// Update track to completed status
	updateCmd := &task_manager.TrackUpdateCommand{Plugin: plugin}
	updateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = updateCmd.Execute(ctx, updateCtx, []string{
		trackID,
		"--status", "complete",
	})
	if err != nil {
		t.Fatalf("failed to update track status: %v", err)
	}

	// Now try to create a task in the completed track - should fail
	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", trackID,
		"--title", "This should fail",
	})

	// Verify the error occurred (nil because error was printed to stdout to avoid help text)
	if err != nil {
		t.Errorf("expected nil error (validation error printed to stdout), got: %v", err)
	}

	// Verify error message is helpful
	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Cannot create task in completed track") {
		t.Errorf("expected error message about completed track, got: %s", output)
	}
	if !strings.Contains(output, "Reopen the track") {
		t.Errorf("expected suggestion to reopen track, got: %s", output)
	}
	if !strings.Contains(output, "different track") {
		t.Errorf("expected suggestion for different track, got: %s", output)
	}

	// Verify task was not created
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})

	tasks, err := repo.ListTasks(ctx, task_manager.TaskFilters{TrackID: trackID})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestTaskCreateCommand_TrackInProgress(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)
	ctx := context.Background()

	// Setup: Create roadmap and in-progress track
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Test vision",
		"--success-criteria", "Test criteria",
	})
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	trackCmd := &task_manager.TrackCreateCommand{Plugin: plugin}
	trackCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = trackCmd.Execute(ctx, trackCtx, []string{
		"--title", "In Progress Track",
		"--description", "A track in progress",
		"--rank", "200",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Extract track ID
	trackOutput := trackCtx.stdout.String()
	trackIDPrefix := "ID:"
	trackIDStart := strings.Index(trackOutput, trackIDPrefix)
	if trackIDStart == -1 {
		t.Fatalf("failed to find track ID in output: %s", trackOutput)
	}
	trackIDStart += len(trackIDPrefix)
	trackIDEnd := strings.Index(trackOutput[trackIDStart:], "\n")
	if trackIDEnd == -1 {
		trackIDEnd = len(trackOutput)
	} else {
		trackIDEnd += trackIDStart
	}
	trackID := strings.TrimSpace(trackOutput[trackIDStart:trackIDEnd])

	// Update track to in-progress status
	updateCmd := &task_manager.TrackUpdateCommand{Plugin: plugin}
	updateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = updateCmd.Execute(ctx, updateCtx, []string{
		trackID,
		"--status", "in-progress",
	})
	if err != nil {
		t.Fatalf("failed to update track status: %v", err)
	}

	// Create task in in-progress track - should succeed
	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", trackID,
		"--title", "This should succeed",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify success message
	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task created successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify task was created
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})

	tasks, err := repo.ListTasks(ctx, task_manager.TaskFilters{TrackID: trackID})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestTaskCreateCommand_TrackNotStarted(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)
	ctx := context.Background()

	// Setup: Create roadmap and not-started track (default status)
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Test vision",
		"--success-criteria", "Test criteria",
	})
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	trackCmd := &task_manager.TrackCreateCommand{Plugin: plugin}
	trackCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = trackCmd.Execute(ctx, trackCtx, []string{
		"--title", "Not Started Track",
		"--description", "A new track",
		"--rank", "200",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Extract track ID
	trackOutput := trackCtx.stdout.String()
	trackIDPrefix := "ID:"
	trackIDStart := strings.Index(trackOutput, trackIDPrefix)
	if trackIDStart == -1 {
		t.Fatalf("failed to find track ID in output: %s", trackOutput)
	}
	trackIDStart += len(trackIDPrefix)
	trackIDEnd := strings.Index(trackOutput[trackIDStart:], "\n")
	if trackIDEnd == -1 {
		trackIDEnd = len(trackOutput)
	} else {
		trackIDEnd += trackIDStart
	}
	trackID := strings.TrimSpace(trackOutput[trackIDStart:trackIDEnd])

	// Create task in not-started track - should succeed
	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", trackID,
		"--title", "This should succeed",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify success message
	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task created successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify task was created
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})

	tasks, err := repo.ListTasks(ctx, task_manager.TaskFilters{TrackID: trackID})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestTaskCreateCommand_TrackBlocked(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)
	ctx := context.Background()

	// Setup: Create roadmap and blocked track
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Test vision",
		"--success-criteria", "Test criteria",
	})
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	trackCmd := &task_manager.TrackCreateCommand{Plugin: plugin}
	trackCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = trackCmd.Execute(ctx, trackCtx, []string{
		"--title", "Blocked Track",
		"--description", "A blocked track",
		"--rank", "200",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Extract track ID
	trackOutput := trackCtx.stdout.String()
	trackIDPrefix := "ID:"
	trackIDStart := strings.Index(trackOutput, trackIDPrefix)
	if trackIDStart == -1 {
		t.Fatalf("failed to find track ID in output: %s", trackOutput)
	}
	trackIDStart += len(trackIDPrefix)
	trackIDEnd := strings.Index(trackOutput[trackIDStart:], "\n")
	if trackIDEnd == -1 {
		trackIDEnd = len(trackOutput)
	} else {
		trackIDEnd += trackIDStart
	}
	trackID := strings.TrimSpace(trackOutput[trackIDStart:trackIDEnd])

	// Update track to blocked status
	updateCmd := &task_manager.TrackUpdateCommand{Plugin: plugin}
	updateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = updateCmd.Execute(ctx, updateCtx, []string{
		trackID,
		"--status", "blocked",
	})
	if err != nil {
		t.Fatalf("failed to update track status: %v", err)
	}

	// Create task in blocked track - should succeed
	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", trackID,
		"--title", "This should succeed",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify task was created
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})

	tasks, err := repo.ListTasks(ctx, task_manager.TaskFilters{TrackID: trackID})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestTaskCreateCommand_TrackWaiting(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)
	ctx := context.Background()

	// Setup: Create roadmap and waiting track
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Test vision",
		"--success-criteria", "Test criteria",
	})
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	trackCmd := &task_manager.TrackCreateCommand{Plugin: plugin}
	trackCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = trackCmd.Execute(ctx, trackCtx, []string{
		"--title", "Waiting Track",
		"--description", "A waiting track",
		"--rank", "200",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Extract track ID
	trackOutput := trackCtx.stdout.String()
	trackIDPrefix := "ID:"
	trackIDStart := strings.Index(trackOutput, trackIDPrefix)
	if trackIDStart == -1 {
		t.Fatalf("failed to find track ID in output: %s", trackOutput)
	}
	trackIDStart += len(trackIDPrefix)
	trackIDEnd := strings.Index(trackOutput[trackIDStart:], "\n")
	if trackIDEnd == -1 {
		trackIDEnd = len(trackOutput)
	} else {
		trackIDEnd += trackIDStart
	}
	trackID := strings.TrimSpace(trackOutput[trackIDStart:trackIDEnd])

	// Update track to waiting status
	updateCmd := &task_manager.TrackUpdateCommand{Plugin: plugin}
	updateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	err = updateCmd.Execute(ctx, updateCtx, []string{
		trackID,
		"--status", "waiting",
	})
	if err != nil {
		t.Fatalf("failed to update track status: %v", err)
	}

	// Create task in waiting track - should succeed
	cmd := &task_manager.TaskCreateCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = cmd.Execute(ctx, cmdCtx, []string{
		"--track", trackID,
		"--title", "This should succeed",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify task was created
	db := getProjectDB(t, tmpDir, "default")
	defer db.Close()
	repo := task_manager.NewSQLiteRoadmapRepository(db, &stubLogger{})

	tasks, err := repo.ListTasks(ctx, task_manager.TaskFilters{TrackID: trackID})
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}
