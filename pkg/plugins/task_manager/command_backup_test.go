package task_manager_test

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"
)

// TestBackupCommand tests basic backup creation
func TestBackupCommand(t *testing.T) {
	// CRITICAL: Use temp directory for HOME to avoid writing to real user directory
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Initialize roadmap to create database
	ctx := context.Background()
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

	// Create backup
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	stdout := &bytes.Buffer{}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     stdout,
		logger:     &stubLogger{},
	}

	if err := backupCmd.Execute(ctx, cmdCtx, []string{}); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	// Verify output contains success message
	output := stdout.String()
	if !strings.Contains(output, "Backup created successfully") {
		t.Errorf("expected success message in output, got: %s", output)
	}

	// Verify backup file was created in temp directory
	backupDir := filepath.Join(tmpHome, ".darwinflow", "task-manager", "backups", "default")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("failed to read backup directory: %v", err)
	}

	backupFound := false
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "task-manager-") && strings.HasSuffix(entry.Name(), ".db") {
			backupFound = true
			break
		}
	}

	if !backupFound {
		t.Error("backup file was not created")
	}
}

// TestBackupCommand_NoDatabase tests backup when database doesn't exist
func TestBackupCommand_NoDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	logger := &stubLogger{}
	plugin, err := task_manager.NewTaskManagerPlugin(logger, tmpDir, nil)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Create default project but don't initialize roadmap
	projectCmd := &task_manager.ProjectCreateCommand{Plugin: plugin}
	projectCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     logger,
	}
	if err := projectCmd.Execute(context.Background(), projectCtx, []string{"default"}); err != nil {
		t.Fatalf("failed to create default project: %v", err)
	}

	// Try to backup (should fail - no roadmap initialized)
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     logger,
	}

	// Remove the database to simulate no initialization
	dbPath := filepath.Join(tmpDir, ".darwinflow", "projects", "default", "roadmap.db")
	os.Remove(dbPath)

	err = backupCmd.Execute(context.Background(), cmdCtx, []string{})
	if err == nil {
		t.Fatal("expected error when database doesn't exist")
	}

	if !strings.Contains(err.Error(), "database not found") {
		t.Errorf("expected 'database not found' error, got: %v", err)
	}
}

// TestBackupCommand_AutoCleanup tests that old backups are cleaned up
func TestBackupCommand_AutoCleanup(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Initialize roadmap
	ctx := context.Background()
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

	// Create 12 backups (should keep only 10)
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	for i := 0; i < 12; i++ {
		stdout := &bytes.Buffer{}
		cmdCtx := &mockCommandContext{
			workingDir: tmpDir,
			stdout:     stdout,
			logger:     &stubLogger{},
		}

		if err := backupCmd.Execute(ctx, cmdCtx, []string{}); err != nil {
			t.Fatalf("backup %d failed: %v", i, err)
		}

		// Sleep 1 second to ensure different timestamps (backup uses second precision)
		time.Sleep(1 * time.Second)
	}

	// Verify only 10 backups remain in temp directory
	backupDir := filepath.Join(tmpHome, ".darwinflow", "task-manager", "backups", "default")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("failed to read backup directory: %v", err)
	}

	backupCount := 0
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "task-manager-") && strings.HasSuffix(entry.Name(), ".db") {
			backupCount++
		}
	}

	if backupCount != 10 {
		t.Errorf("expected 10 backups after cleanup, got %d", backupCount)
	}
}

// TestBackupCommand_TimestampFormat tests backup filename format
func TestBackupCommand_TimestampFormat(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Initialize roadmap
	ctx := context.Background()
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

	// Create backup
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	stdout := &bytes.Buffer{}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     stdout,
		logger:     &stubLogger{},
	}

	if err := backupCmd.Execute(ctx, cmdCtx, []string{}); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	// Verify filename format (task-manager-YYYY-MM-DD-HHMMSS.db) in temp directory
	backupDir := filepath.Join(tmpHome, ".darwinflow", "task-manager", "backups", "default")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("failed to read backup directory: %v", err)
	}

	foundValidFilename := false
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "task-manager-") && strings.HasSuffix(name, ".db") {
			// Check format: task-manager-2025-11-03-143022.db
			// Should be: task-manager- (13) + YYYY-MM-DD (10) + - (1) + HHMMSS (6) + .db (3) = 33 characters
			if len(name) == 33 {
				foundValidFilename = true
			}
		}
	}

	if !foundValidFilename {
		t.Errorf("backup filename does not match expected format (task-manager-YYYY-MM-DD-HHMMSS.db), found: %v", entries)
	}
}

// TestRestoreCommand tests basic restore functionality
func TestRestoreCommand(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Initialize roadmap
	ctx := context.Background()
	roadmapCmd := &task_manager.RoadmapInitCommand{Plugin: plugin}
	roadmapCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	if err := roadmapCmd.Execute(ctx, roadmapCtx, []string{
		"--vision", "Original vision",
		"--success-criteria", "Original criteria",
	}); err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	// Create backup
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	stdout := &bytes.Buffer{}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     stdout,
		logger:     &stubLogger{},
	}

	if err := backupCmd.Execute(ctx, cmdCtx, []string{}); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	// Get backup file path from temp directory
	backupDir := filepath.Join(tmpHome, ".darwinflow", "task-manager", "backups", "default")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("failed to read backup directory: %v", err)
	}

	var backupPath string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "task-manager-") && strings.HasSuffix(entry.Name(), ".db") {
			backupPath = filepath.Join(backupDir, entry.Name())
			break
		}
	}

	if backupPath == "" {
		t.Fatal("backup file not found")
	}

	// Modify roadmap
	updateCmd := &task_manager.RoadmapUpdateCommand{Plugin: plugin}
	updateCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	if err := updateCmd.Execute(ctx, updateCtx, []string{
		"--vision", "Modified vision",
	}); err != nil {
		t.Fatalf("failed to update roadmap: %v", err)
	}

	// Restore from backup
	restoreCmd := &task_manager.RestoreCommand{Plugin: plugin}
	restoreStdout := &bytes.Buffer{}
	restoreCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     restoreStdout,
		logger:     &stubLogger{},
	}

	if err := restoreCmd.Execute(ctx, restoreCtx, []string{backupPath}); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	// Verify output contains success message
	output := restoreStdout.String()
	if !strings.Contains(output, "Database restored successfully") {
		t.Errorf("expected success message in output, got: %s", output)
	}

	// Verify safety backup was created
	if !strings.Contains(output, "Safety backup created") {
		t.Error("expected safety backup message in output")
	}

	// Verify roadmap was restored
	showCmd := &task_manager.RoadmapShowCommand{Plugin: plugin}
	showStdout := &bytes.Buffer{}
	showCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     showStdout,
		logger:     &stubLogger{},
	}
	if err := showCmd.Execute(ctx, showCtx, []string{}); err != nil {
		t.Fatalf("failed to show roadmap: %v", err)
	}

	roadmapOutput := showStdout.String()
	if !strings.Contains(roadmapOutput, "Original vision") {
		t.Error("roadmap was not restored correctly")
	}
}

// TestRestoreCommand_InvalidBackup tests restore with non-existent backup
func TestRestoreCommand_InvalidBackup(t *testing.T) {
	plugin, tmpDir := setupTestPlugin(t)

	restoreCmd := &task_manager.RestoreCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err := restoreCmd.Execute(ctx, cmdCtx, []string{"/nonexistent/backup.db"})
	if err == nil {
		t.Fatal("expected error for non-existent backup file")
	}

	if !strings.Contains(err.Error(), "backup file not found") {
		t.Errorf("expected 'backup file not found' error, got: %v", err)
	}
}

// TestRestoreCommand_CorruptedBackup tests restore with corrupted backup
func TestRestoreCommand_CorruptedBackup(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Create a corrupted backup file
	corruptedBackup := filepath.Join(tmpDir, "corrupted.db")
	if err := os.WriteFile(corruptedBackup, []byte("not a valid sqlite database"), 0644); err != nil {
		t.Fatalf("failed to create corrupted backup: %v", err)
	}

	restoreCmd := &task_manager.RestoreCommand{Plugin: plugin}
	ctx := context.Background()
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err := restoreCmd.Execute(ctx, cmdCtx, []string{corruptedBackup})
	if err == nil {
		t.Fatal("expected error for corrupted backup file")
	}

	if !strings.Contains(err.Error(), "integrity check failed") {
		t.Errorf("expected integrity check error, got: %v", err)
	}
}

// TestRestoreCommand_SafetyBackup tests that safety backup is created
func TestRestoreCommand_SafetyBackup(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Initialize roadmap
	ctx := context.Background()
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

	// Create backup
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	if err := backupCmd.Execute(ctx, cmdCtx, []string{}); err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	// Get backup file path from temp directory
	backupDir := filepath.Join(tmpHome, ".darwinflow", "task-manager", "backups", "default")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("failed to read backup directory: %v", err)
	}

	var backupPath string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "task-manager-") && strings.HasSuffix(entry.Name(), ".db") {
			backupPath = filepath.Join(backupDir, entry.Name())
			break
		}
	}

	// Restore
	restoreCmd := &task_manager.RestoreCommand{Plugin: plugin}
	restoreCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}
	if err := restoreCmd.Execute(ctx, restoreCtx, []string{backupPath}); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	// Verify safety backup exists in project folder
	projectDir := filepath.Join(tmpDir, ".darwinflow", "projects", "default")
	safetyEntries, err := os.ReadDir(projectDir)
	if err != nil {
		t.Fatalf("failed to read project directory: %v", err)
	}

	safetyBackupFound := false
	for _, entry := range safetyEntries {
		if strings.HasPrefix(entry.Name(), "roadmap-safety-backup-") && strings.HasSuffix(entry.Name(), ".db") {
			safetyBackupFound = true
			break
		}
	}

	if !safetyBackupFound {
		t.Error("safety backup was not created in project folder")
	}
}

// TestBackupListCommand tests listing backups
func TestBackupListCommand(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Initialize roadmap
	ctx := context.Background()
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

	// Create 3 backups
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	for i := 0; i < 3; i++ {
		cmdCtx := &mockCommandContext{
			workingDir: tmpDir,
			stdout:     &bytes.Buffer{},
			logger:     &stubLogger{},
		}
		if err := backupCmd.Execute(ctx, cmdCtx, []string{}); err != nil {
			t.Fatalf("backup %d failed: %v", i, err)
		}
		// Sleep 1 second to ensure different timestamps (backup uses second precision)
		time.Sleep(1 * time.Second)
	}

	// List backups
	listCmd := &task_manager.BackupListCommand{Plugin: plugin}
	stdout := &bytes.Buffer{}
	listCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     stdout,
		logger:     &stubLogger{},
	}

	if err := listCmd.Execute(ctx, listCtx, []string{}); err != nil {
		t.Fatalf("backup list failed: %v", err)
	}

	// Verify output
	output := stdout.String()
	if !strings.Contains(output, "Backups for project: default") {
		t.Error("expected project name in output")
	}
	if !strings.Contains(output, "Total: 3 backups") {
		t.Errorf("expected 3 backups in output, got: %s", output)
	}

	// Verify backup files are listed
	backupCount := strings.Count(output, "task-manager-")
	if backupCount != 3 {
		t.Errorf("expected 3 backup files listed, got %d", backupCount)
	}
}

// TestBackupListCommand_NoBackups tests listing when no backups exist
func TestBackupListCommand_NoBackups(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	listCmd := &task_manager.BackupListCommand{Plugin: plugin}
	ctx := context.Background()
	stdout := &bytes.Buffer{}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     stdout,
		logger:     &stubLogger{},
	}

	if err := listCmd.Execute(ctx, cmdCtx, []string{}); err != nil {
		t.Fatalf("backup list failed: %v", err)
	}

	// Verify output indicates no backups
	output := stdout.String()
	if !strings.Contains(output, "No backups found") {
		t.Errorf("expected 'No backups found' in output, got: %s", output)
	}
}

// TestBackupIntegrity tests that backup validates database integrity
func TestBackupIntegrity(t *testing.T) {
	// CRITICAL: Use temp directory for HOME
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	plugin, tmpDir := setupTestPlugin(t)

	// Initialize roadmap
	ctx := context.Background()
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

	// Corrupt the database
	dbPath := filepath.Join(tmpDir, ".darwinflow", "projects", "default", "roadmap.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	db.Close()

	// Overwrite with corrupted data
	if err := os.WriteFile(dbPath, []byte("corrupted data"), 0644); err != nil {
		t.Fatalf("failed to corrupt database: %v", err)
	}

	// Try to backup (should fail)
	backupCmd := &task_manager.BackupCommand{Plugin: plugin}
	cmdCtx := &mockCommandContext{
		workingDir: tmpDir,
		stdout:     &bytes.Buffer{},
		logger:     &stubLogger{},
	}

	err = backupCmd.Execute(ctx, cmdCtx, []string{})
	if err == nil {
		t.Fatal("expected error when backing up corrupted database")
	}

	if !strings.Contains(err.Error(), "integrity check failed") {
		t.Errorf("expected integrity check error, got: %v", err)
	}
}
