package application

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
)

// ProjectApplicationService handles project management operations
type ProjectApplicationService struct {
	repo              repositories.ProjectManagementRepository
	validationService *services.ValidationService
}

// NewProjectService creates a new project application service
func NewProjectService(
	repo repositories.ProjectManagementRepository,
	validationService *services.ValidationService,
) *ProjectApplicationService {
	return &ProjectApplicationService{
		repo:              repo,
		validationService: validationService,
	}
}

// CreateProject creates a new project with validation
// Returns the generated or provided project code
func (s *ProjectApplicationService) CreateProject(projectName, projectCode string) (string, error) {
	// Generate default project code if not provided
	if projectCode == "" {
		projectCode = generateDefaultProjectCode(projectName)
	}

	// Validate project name
	projectNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !projectNameRegex.MatchString(projectName) {
		return "", fmt.Errorf("invalid project name: must be alphanumeric with hyphens or underscores only")
	}

	// Validate project code
	if !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(projectCode) {
		return "", fmt.Errorf("invalid project code: must be alphanumeric uppercase (e.g., DW, PROD, TEST)")
	}

	// Check if project already exists
	exists, err := s.repo.ProjectExists(projectName)
	if err != nil {
		return "", fmt.Errorf("failed to check project existence: %w", err)
	}
	if exists {
		return "", fmt.Errorf("project '%s' already exists", projectName)
	}

	// Create project via repository
	if err := s.repo.CreateProject(projectName, projectCode); err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}

	return projectCode, nil
}

// ListProjects lists all projects
func (s *ProjectApplicationService) ListProjects() ([]string, error) {
	projects, err := s.repo.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	return projects, nil
}

// GetActiveProject returns the currently active project
func (s *ProjectApplicationService) GetActiveProject() (string, error) {
	activeProject, err := s.repo.GetActiveProject()
	if err != nil {
		return "", fmt.Errorf("failed to get active project: %w", err)
	}
	return activeProject, nil
}

// SwitchProject switches to a different project with validation
func (s *ProjectApplicationService) SwitchProject(projectName string) error {
	// Verify project exists
	exists, err := s.repo.ProjectExists(projectName)
	if err != nil {
		return fmt.Errorf("failed to check project existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("project '%s' does not exist", projectName)
	}

	// Set active project via repository
	if err := s.repo.SetActiveProject(projectName); err != nil {
		return fmt.Errorf("failed to switch project: %w", err)
	}

	return nil
}

// DeleteProject deletes a project
func (s *ProjectApplicationService) DeleteProject(projectName string) error {
	// Verify project exists
	exists, err := s.repo.ProjectExists(projectName)
	if err != nil {
		return fmt.Errorf("failed to check project existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("project '%s' does not exist", projectName)
	}

	// Check if project is active
	activeProject, err := s.repo.GetActiveProject()
	if err == nil && activeProject == projectName {
		return fmt.Errorf("cannot delete active project '%s': switch to another project first", projectName)
	}

	// Delete project via repository
	if err := s.repo.DeleteProject(projectName); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// GetProjectInfo returns information about a project
func (s *ProjectApplicationService) GetProjectInfo(projectName string) (map[string]string, error) {
	// Verify project exists
	exists, err := s.repo.ProjectExists(projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to check project existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("project '%s' does not exist", projectName)
	}

	// Get project info via repository
	info, err := s.repo.GetProjectInfo(projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get project info: %w", err)
	}

	return info, nil
}

// GenerateDefaultProjectCode generates a default project code from the project name
func generateDefaultProjectCode(projectName string) string {
	// Extract uppercase letters from project name
	var code string
	parts := strings.FieldsFunc(projectName, func(r rune) bool {
		return r == '-' || r == '_'
	})

	for _, part := range parts {
		if len(part) > 0 {
			code += strings.ToUpper(string(part[0]))
		}
	}

	// If not enough letters, use first 2-3 uppercase letters
	if len(code) < 2 && len(projectName) > 0 {
		code = strings.ToUpper(projectName[:1])
		if len(projectName) > 1 {
			code += strings.ToUpper(projectName[1:2])
		}
	}

	// Default to "TM" if still empty
	if code == "" {
		code = "TM"
	}

	return code
}
