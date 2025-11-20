package repositories

// ProjectManagementRepository handles CRUD operations for project management
type ProjectManagementRepository interface {
	// CreateProject creates a new project with its own database
	CreateProject(projectName, projectCode string) error

	// ListProjects returns all project names in the workspace
	ListProjects() ([]string, error)

	// GetActiveProject returns the currently active project name
	GetActiveProject() (string, error)

	// SetActiveProject sets the active project
	SetActiveProject(projectName string) error

	// DeleteProject deletes a project and its database
	DeleteProject(projectName string) error

	// GetProjectInfo returns metadata about a project
	GetProjectInfo(projectName string) (map[string]string, error)

	// ProjectExists checks if a project exists
	ProjectExists(projectName string) (bool, error)
}
