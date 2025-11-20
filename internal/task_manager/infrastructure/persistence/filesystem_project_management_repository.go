package persistence

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
)

// Compile-time interface check
var _ repositories.ProjectManagementRepository = (*FileSystemProjectManagementRepository)(nil)

// FileSystemProjectManagementRepository implements ProjectManagementRepository using filesystem operations
type FileSystemProjectManagementRepository struct {
	workingDir string
}

// NewFileSystemProjectManagementRepository creates a new filesystem-based project management repository
func NewFileSystemProjectManagementRepository(workingDir string) *FileSystemProjectManagementRepository {
	return &FileSystemProjectManagementRepository{
		workingDir: workingDir,
	}
}

// CreateProject creates a new project with its own database
func (r *FileSystemProjectManagementRepository) CreateProject(projectName, projectCode string) error {
	// Check if project already exists
	projectsDir := filepath.Join(r.workingDir, "projects")
	projectDir := filepath.Join(projectsDir, projectName)

	if _, err := os.Stat(projectDir); err == nil {
		return fmt.Errorf("project '%s' already exists", projectName)
	}

	// Create project directory
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create and initialize database for the project
	_, db, err := OpenProjectDatabase(r.workingDir, projectName)
	if err != nil {
		return fmt.Errorf("failed to create project database: %w", err)
	}
	defer db.Close()

	// Store project code in a config file
	projectConfig := filepath.Join(projectDir, "config.txt")
	configContent := fmt.Sprintf("code=%s\ncreated=%s\n", projectCode, time.Now().Format(time.RFC3339))
	if err := os.WriteFile(projectConfig, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	return nil
}

// ListProjects returns all project names in the workspace
func (r *FileSystemProjectManagementRepository) ListProjects() ([]string, error) {
	projectsDir := filepath.Join(r.workingDir, "projects")

	// Check if projects directory exists
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	var projects []string
	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}

	sort.Strings(projects)
	return projects, nil
}

// GetActiveProject returns the currently active project name
func (r *FileSystemProjectManagementRepository) GetActiveProject() (string, error) {
	activeProjectFile := filepath.Join(r.workingDir, "active-project.txt")

	content, err := os.ReadFile(activeProjectFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "default", nil
		}
		return "", fmt.Errorf("failed to read active project: %w", err)
	}

	projectName := strings.TrimSpace(string(content))
	if projectName == "" {
		return "default", nil
	}

	return projectName, nil
}

// SetActiveProject sets the active project
func (r *FileSystemProjectManagementRepository) SetActiveProject(projectName string) error {
	// Ensure working directory exists
	if err := os.MkdirAll(r.workingDir, 0755); err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}

	activeProjectFile := filepath.Join(r.workingDir, "active-project.txt")
	if err := os.WriteFile(activeProjectFile, []byte(projectName), 0644); err != nil {
		return fmt.Errorf("failed to write active project: %w", err)
	}

	return nil
}

// DeleteProject deletes a project and its database
func (r *FileSystemProjectManagementRepository) DeleteProject(projectName string) error {
	projectsDir := filepath.Join(r.workingDir, "projects")
	projectDir := filepath.Join(projectsDir, projectName)

	// Check if project exists
	if _, err := os.Stat(projectDir); err != nil {
		return fmt.Errorf("project '%s' does not exist", projectName)
	}

	// Remove project directory
	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	// If this was the active project, reset to default
	activeProject, _ := r.GetActiveProject()
	if activeProject == projectName {
		_ = r.SetActiveProject("default")
	}

	return nil
}

// GetProjectInfo returns metadata about a project
func (r *FileSystemProjectManagementRepository) GetProjectInfo(projectName string) (map[string]string, error) {
	projectsDir := filepath.Join(r.workingDir, "projects")
	projectDir := filepath.Join(projectsDir, projectName)

	// Check if project exists
	if _, err := os.Stat(projectDir); err != nil {
		return nil, fmt.Errorf("project '%s' does not exist", projectName)
	}

	info := make(map[string]string)
	info["name"] = projectName
	info["location"] = projectDir

	// Try to read project config
	projectConfig := filepath.Join(projectDir, "config.txt")
	if content, err := os.ReadFile(projectConfig); err == nil {
		// Parse config (key=value format)
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				info[key] = value
			}
		}
	}

	// Get database file info
	dbFile := filepath.Join(projectDir, "roadmap.db")
	if stat, err := os.Stat(dbFile); err == nil {
		info["database_size"] = fmt.Sprintf("%d bytes", stat.Size())
		info["modified"] = stat.ModTime().Format("2006-01-02 15:04:05")
	}

	return info, nil
}

// ProjectExists checks if a project exists
func (r *FileSystemProjectManagementRepository) ProjectExists(projectName string) (bool, error) {
	projectsDir := filepath.Join(r.workingDir, "projects")
	projectDir := filepath.Join(projectsDir, projectName)

	_, err := os.Stat(projectDir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
