package task_manager_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	task_manager "github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"
)

// ============================================================================
// MigrateIDsCommand Tests
// ============================================================================

func TestMigrateIDsCommand_NoIDsToMigrate(t *testing.T) {
	plugin, tmpDir, _ := setupTestWithRoadmapAndTrack(t)
	ctx := context.Background()

	// The track created by helper uses new ID format already
	// Execute migrate command - should report nothing to migrate
	migrateCmd := &task_manager.MigrateIDsCommand{Plugin: plugin}
	migrateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	if err := migrateCmd.Execute(ctx, migrateCtx, []string{}); err != nil {
		t.Fatalf("MigrateIDsCommand failed: %v", err)
	}

	output := migrateCtx.stdout.String()

	// Should indicate no migration needed
	if !strings.Contains(output, "No IDs need migration") || !strings.Contains(output, "already in new format") {
		t.Errorf("Expected 'No IDs need migration' message, got: %s", output)
	}
}

func TestMigrateIDsCommand_DryRun(t *testing.T) {
	plugin, tmpDir, _ := setupTestWithRoadmapAndTrack(t)
	ctx := context.Background()

	// Execute migrate with dry-run (even though nothing to migrate)
	migrateCmd := &task_manager.MigrateIDsCommand{Plugin: plugin}
	migrateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	if err := migrateCmd.Execute(ctx, migrateCtx, []string{"--dry-run"}); err != nil {
		t.Fatalf("MigrateIDsCommand with --dry-run failed: %v", err)
	}

	output := migrateCtx.stdout.String()

	// Dry run should complete successfully
	// Output will say "No IDs need migration" since we're using new format
	if !strings.Contains(output, "No IDs need migration") && !strings.Contains(output, "Dry run") {
		t.Errorf("Expected migration analysis output, got: %s", output)
	}
}

func TestMigrateIDsCommand_Idempotent(t *testing.T) {
	plugin, tmpDir, _ := setupTestWithRoadmapAndTrack(t)
	ctx := context.Background()

	// Run migration twice - should be idempotent
	migrateCmd1 := &task_manager.MigrateIDsCommand{Plugin: plugin}
	migrateCtx1 := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	if err := migrateCmd1.Execute(ctx, migrateCtx1, []string{}); err != nil {
		t.Fatalf("First migration failed: %v", err)
	}

	// Run again
	migrateCmd2 := &task_manager.MigrateIDsCommand{Plugin: plugin}
	migrateCtx2 := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	if err := migrateCmd2.Execute(ctx, migrateCtx2, []string{}); err != nil {
		t.Fatalf("Second migration failed: %v", err)
	}

	output2 := migrateCtx2.stdout.String()

	// Second run should also say no migration needed
	if !strings.Contains(output2, "No IDs need migration") {
		t.Errorf("Expected idempotent behavior, got: %s", output2)
	}
}

func TestMigrateIDsCommand_WithProject(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)
	ctx := context.Background()

	// Create a second project
	projectCmd := &task_manager.ProjectCreateCommand{Plugin: plugin}
	projectCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	if err := projectCmd.Execute(ctx, projectCtx, []string{"test-project"}); err != nil {
		t.Fatalf("failed to create test project: %v", err)
	}

	// Switch to test project
	switchCmd := &task_manager.ProjectSwitchCommand{Plugin: plugin}
	switchCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	if err := switchCmd.Execute(ctx, switchCtx, []string{"test-project"}); err != nil {
		t.Fatalf("failed to switch project: %v", err)
	}

	// Create roadmap in test project
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	if err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Test vision",
		"--success-criteria", "Test criteria",
	}); err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	// Run migration on test-project explicitly
	migrateCmd := &task_manager.MigrateIDsCommand{Plugin: plugin}
	migrateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	if err := migrateCmd.Execute(ctx, migrateCtx, []string{"--project", "test-project"}); err != nil {
		t.Fatalf("MigrateIDsCommand with --project failed: %v", err)
	}

	output := migrateCtx.stdout.String()

	// Should complete successfully
	if !strings.Contains(output, "test-project") {
		t.Errorf("Expected project name in output, got: %s", output)
	}
}
