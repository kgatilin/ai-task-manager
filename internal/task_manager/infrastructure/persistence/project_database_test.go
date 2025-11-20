package persistence_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/persistence"
)

// ResolveWorkingDirectory tests

func TestResolveWorkingDirectory_FromEnv(t *testing.T) {
	// Setup: Set TM_WORKING_DIR env var
	tempDir := t.TempDir()
	t.Setenv("TM_WORKING_DIR", tempDir)

	// Test: Resolve working directory
	got := persistence.ResolveWorkingDirectory()

	// Verify: Should return env var value
	if got != tempDir {
		t.Errorf("ResolveWorkingDirectory() = %v, want %v", got, tempDir)
	}
}

func TestResolveWorkingDirectory_FromCwd(t *testing.T) {
	// Setup: Ensure TM_WORKING_DIR is not set
	t.Setenv("TM_WORKING_DIR", "")

	// Test: Resolve working directory
	got := persistence.ResolveWorkingDirectory()

	// Verify: Should return current working directory + "/.tm"
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() failed: %v", err)
	}
	expected := filepath.Join(cwd, ".tm")
	if got != expected {
		t.Errorf("ResolveWorkingDirectory() = %v, want %v", got, expected)
	}
}

// GetActiveProject tests

func TestGetActiveProject_Default(t *testing.T) {
	// Setup: Create temporary directory without active-project.txt
	tempDir := t.TempDir()

	// Test: Get active project
	got, err := persistence.GetActiveProject(tempDir)

	// Verify: Should return "default" when file doesn't exist
	if err != nil {
		t.Fatalf("GetActiveProject() error = %v", err)
	}
	if got != "default" {
		t.Errorf("GetActiveProject() = %v, want default", got)
	}
}

func TestGetActiveProject_Explicit(t *testing.T) {
	// Setup: Create active-project.txt with explicit project name
	tempDir := t.TempDir()
	activeProjectFile := filepath.Join(tempDir, "active-project.txt")
	if err := os.WriteFile(activeProjectFile, []byte("my-project\n"), 0644); err != nil {
		t.Fatalf("failed to write active-project.txt: %v", err)
	}

	// Test: Get active project
	got, err := persistence.GetActiveProject(tempDir)

	// Verify: Should return explicit project name
	if err != nil {
		t.Fatalf("GetActiveProject() error = %v", err)
	}
	if got != "my-project" {
		t.Errorf("GetActiveProject() = %v, want my-project", got)
	}
}

func TestGetActiveProject_Empty(t *testing.T) {
	// Setup: Create active-project.txt with whitespace only
	tempDir := t.TempDir()
	activeProjectFile := filepath.Join(tempDir, "active-project.txt")
	if err := os.WriteFile(activeProjectFile, []byte("  \n"), 0644); err != nil {
		t.Fatalf("failed to write active-project.txt: %v", err)
	}

	// Test: Get active project
	got, err := persistence.GetActiveProject(tempDir)

	// Verify: Should return "default" when file is empty
	if err != nil {
		t.Fatalf("GetActiveProject() error = %v", err)
	}
	if got != "default" {
		t.Errorf("GetActiveProject() = %v, want default", got)
	}
}

// GetProjectDatabasePath tests

func TestGetProjectDatabasePath(t *testing.T) {
	// Test: Get database path for project
	workingDir := "/home/user/projects/myapp/.tm"
	projectName := "test-project"

	got := persistence.GetProjectDatabasePath(workingDir, projectName)

	// Verify: Path follows expected structure
	expected := filepath.Join(workingDir, "projects", projectName, "roadmap.db")
	if got != expected {
		t.Errorf("GetProjectDatabasePath() = %v, want %v", got, expected)
	}
}

// OpenProjectDatabase tests

func TestOpenProjectDatabase_Success(t *testing.T) {
	// Setup: Create temporary working directory
	tempDir := t.TempDir()

	// Test: Open database
	dbPath, db, err := persistence.OpenProjectDatabase(tempDir, "test-project")
	if err != nil {
		t.Fatalf("OpenProjectDatabase() error = %v", err)
	}
	defer db.Close()

	// Verify: Database path is correct
	expectedPath := filepath.Join(tempDir, "projects", "test-project", "roadmap.db")
	if dbPath != expectedPath {
		t.Errorf("dbPath = %v, want %v", dbPath, expectedPath)
	}

	// Verify: Database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file does not exist at %v", dbPath)
	}

	// Verify: Database is functional
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Errorf("Database query failed: %v", err)
	}
	if result != 1 {
		t.Errorf("Database query result = %v, want 1", result)
	}

	// Verify: Schema was initialized (check for roadmaps table)
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='roadmaps'").Scan(&tableExists)
	if err != nil {
		t.Errorf("Failed to check for roadmaps table: %v", err)
	}
	if tableExists != 1 {
		t.Errorf("roadmaps table does not exist, schema initialization failed")
	}
}

func TestOpenProjectDatabase_MultipleProjects(t *testing.T) {
	// Setup: Create temporary working directory
	tempDir := t.TempDir()

	// Test: Open first project
	dbPath1, db1, err := persistence.OpenProjectDatabase(tempDir, "project1")
	if err != nil {
		t.Fatalf("OpenProjectDatabase(project1) error = %v", err)
	}
	defer db1.Close()

	// Test: Open second project
	dbPath2, db2, err := persistence.OpenProjectDatabase(tempDir, "project2")
	if err != nil {
		t.Fatalf("OpenProjectDatabase(project2) error = %v", err)
	}
	defer db2.Close()

	// Verify: Different database paths
	if dbPath1 == dbPath2 {
		t.Errorf("Different projects should have different database paths")
	}

	// Verify: Both database files exist
	if _, err := os.Stat(dbPath1); os.IsNotExist(err) {
		t.Errorf("Database file for project1 does not exist")
	}
	if _, err := os.Stat(dbPath2); os.IsNotExist(err) {
		t.Errorf("Database file for project2 does not exist")
	}
}

func TestOpenProjectDatabase_ExistingDatabase(t *testing.T) {
	// Setup: Create temporary working directory
	tempDir := t.TempDir()

	// Test: Create database first time
	dbPath1, db1, err := persistence.OpenProjectDatabase(tempDir, "existing-project")
	if err != nil {
		t.Fatalf("First OpenProjectDatabase() error = %v", err)
	}
	db1.Close()

	// Test: Open same database again
	dbPath2, db2, err := persistence.OpenProjectDatabase(tempDir, "existing-project")
	if err != nil {
		t.Fatalf("Second OpenProjectDatabase() error = %v", err)
	}
	defer db2.Close()

	// Verify: Same database path
	if dbPath1 != dbPath2 {
		t.Errorf("Same project should return same database path")
	}

	// Verify: Database is still functional
	var result int
	err = db2.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Errorf("Database query failed: %v", err)
	}
}
