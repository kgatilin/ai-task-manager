package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// defaultDir is the default task manager directory name
	defaultDir = ".tm"
)

// ResolveWorkingDirectory determines the working directory for the task manager.
// Returns the .tm directory path.
// Priority:
//  1. Environment variable TM_WORKING_DIR (used as-is if set)
//  2. Current working directory + "/.tm"
//  3. Fallback to ".tm" if current directory cannot be determined
func ResolveWorkingDirectory() string {
	if wd := os.Getenv("TM_WORKING_DIR"); wd != "" {
		return wd
	}

	wd, err := os.Getwd()
	if err != nil {
		return defaultDir
	}
	return filepath.Join(wd, defaultDir)
}

// GetActiveProject reads the active project name from the working directory.
// Reads from <workingDir>/active-project.txt.
// Returns "default" if the file doesn't exist or is empty.
func GetActiveProject(workingDir string) (string, error) {
	activeProjectFile := filepath.Join(workingDir, "active-project.txt")

	data, err := os.ReadFile(activeProjectFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "default", nil
		}
		return "", fmt.Errorf("failed to read active project file: %w", err)
	}

	projectName := strings.TrimSpace(string(data))
	if projectName == "" {
		return "default", nil
	}
	return projectName, nil
}

// GetProjectDatabasePath returns the database file path for a project.
// Does not open the database or verify its existence.
func GetProjectDatabasePath(workingDir, projectName string) string {
	return filepath.Join(workingDir, "projects", projectName, "roadmap.db")
}

// OpenProjectDatabase opens or creates a project database and runs migrations.
// Creates the project directory structure if it doesn't exist.
// Returns the database path and an open database connection.
func OpenProjectDatabase(workingDir, projectName string) (string, *sql.DB, error) {
	projectDir := filepath.Join(workingDir, "projects", projectName)

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return "", nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	dbPath := GetProjectDatabasePath(workingDir, projectName)

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return "", nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize schema (run migrations)
	if err := InitSchema(db); err != nil {
		db.Close()
		return "", nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return dbPath, db, nil
}
