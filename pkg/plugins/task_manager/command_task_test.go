package task_manager_test

import (
	"bytes"
	"context"
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
		"--priority", "high",
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
		"--priority", "high",
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
		if tasks[0].Priority != "high" {
			t.Errorf("expected priority 'high', got '%s'", tasks[0].Priority)
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
		"--priority", "high",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Extract track ID from output
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

	// Create multiple tasks
	for i := 0; i < 3; i++ {
		taskCmd := &task_manager.TaskCreateCommand{Plugin: plugin}
		taskCtx := &mockCommandContext{
			workingDir: tmpDir,
			stdout:     &bytes.Buffer{},
			logger:     &stubLogger{},
		}
		err = taskCmd.Execute(ctx, taskCtx, []string{
			"--track", trackID,
			"--title", "Task " + string(rune(i+49)),
			"--priority", "medium",
		})
		if err != nil {
			t.Fatalf("failed to create task: %v", err)
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

	// Setup
	repo := plugin.GetRepository()
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
		"high",
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
		task := task_manager.NewTaskEntity(
			"task-"+string(rune(i+49)),
			"track-test",
			"Task "+string(rune(i+49)),
			"",
			status,
			"medium",
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

	// Setup
	repo := plugin.GetRepository()
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
			"high",
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
		task := task_manager.NewTaskEntity(
			"task-"+string(rune(i+49)),
			trackID,
			"Task "+string(rune(i+49)),
			"",
			"todo",
			"medium",
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

	// Setup
	repo := plugin.GetRepository()
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
		"high",
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
		"task-123",
		"track-test",
		"Test Task",
		"Test description",
		"in-progress",
		"high",
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

	err = cmd.Execute(ctx, cmdCtx, []string{"task-123"})
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

// ============================================================================
// TaskUpdateCommand Tests
// ============================================================================

func TestTaskUpdateCommand_UpdateStatus(t *testing.T) {
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

	// Setup
	repo := plugin.GetRepository()
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
		"high",
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
		"task-123",
		"track-test",
		"Test Task",
		"",
		"todo",
		"medium",
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

	err = cmd.Execute(ctx, cmdCtx, []string{"task-123", "--status", "done"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task updated successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify update
	updated, err := repo.GetTask(ctx, "task-123")
	if err != nil {
		t.Fatalf("failed to get updated task: %v", err)
	}
	if updated.Status != "done" {
		t.Errorf("expected status 'done', got '%s'", updated.Status)
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

	// Setup
	repo := plugin.GetRepository()
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
		"high",
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
		"task-123",
		"track-test",
		"Test Task",
		"",
		"todo",
		"medium",
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

	err = cmd.Execute(ctx, cmdCtx, []string{"task-123", "--force"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task deleted successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify task was deleted
	_, err = repo.GetTask(ctx, "task-123")
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

	// Setup
	repo := plugin.GetRepository()
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
			"high",
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
		"task-123",
		"track-test-1",
		"Test Task",
		"",
		"todo",
		"medium",
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

	err = cmd.Execute(ctx, cmdCtx, []string{"task-123", "--track", "track-test-2"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := cmdCtx.stdout.String()
	if !strings.Contains(output, "Task moved successfully") {
		t.Errorf("expected success message, got: %s", output)
	}

	// Verify move
	moved, err := repo.GetTask(ctx, "task-123")
	if err != nil {
		t.Fatalf("failed to get moved task: %v", err)
	}
	if moved.TrackID != "track-test-2" {
		t.Errorf("expected task to be in 'track-test-2', got '%s'", moved.TrackID)
	}
}

func TestTaskMoveCommand_NewTrackNotFound(t *testing.T) {
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

	// Setup
	repo := plugin.GetRepository()
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
		"high",
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
		"task-123",
		"track-test",
		"Test Task",
		"",
		"todo",
		"medium",
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

	err = cmd.Execute(ctx, cmdCtx, []string{"task-123", "--track", "nonexistent-track"})
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
