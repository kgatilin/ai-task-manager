package mocks

// MockProjectManagementRepository is a mock implementation of repositories.ProjectManagementRepository for testing.
type MockProjectManagementRepository struct {
	// In-memory storage for testing
	projects      map[string]map[string]string // project name → project info
	activeProject string

	// CreateProjectFunc is called by CreateProject. If nil, uses default implementation.
	CreateProjectFunc func(projectName, projectCode string) error

	// ListProjectsFunc is called by ListProjects. If nil, returns list of projects.
	ListProjectsFunc func() ([]string, error)

	// GetActiveProjectFunc is called by GetActiveProject. If nil, returns activeProject.
	GetActiveProjectFunc func() (string, error)

	// SetActiveProjectFunc is called by SetActiveProject. If nil, uses default implementation.
	SetActiveProjectFunc func(projectName string) error

	// DeleteProjectFunc is called by DeleteProject. If nil, uses default implementation.
	DeleteProjectFunc func(projectName string) error

	// GetProjectInfoFunc is called by GetProjectInfo. If nil, returns project info from storage.
	GetProjectInfoFunc func(projectName string) (map[string]string, error)

	// ProjectExistsFunc is called by ProjectExists. If nil, checks projects map.
	ProjectExistsFunc func(projectName string) (bool, error)
}

// NewMockProjectManagementRepository creates a new mock project management repository
func NewMockProjectManagementRepository() *MockProjectManagementRepository {
	return &MockProjectManagementRepository{
		projects:      make(map[string]map[string]string),
		activeProject: "",
	}
}

// CreateProject implements repositories.ProjectManagementRepository.
func (m *MockProjectManagementRepository) CreateProject(projectName, projectCode string) error {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(projectName, projectCode)
	}
	// Default implementation: store in memory
	m.projects[projectName] = map[string]string{
		"name": projectName,
		"code": projectCode,
	}
	// Set as active if it's the first project
	if m.activeProject == "" {
		m.activeProject = projectName
	}
	return nil
}

// ListProjects implements repositories.ProjectManagementRepository.
func (m *MockProjectManagementRepository) ListProjects() ([]string, error) {
	if m.ListProjectsFunc != nil {
		return m.ListProjectsFunc()
	}
	// Default implementation: return all project names
	var names []string
	for name := range m.projects {
		names = append(names, name)
	}
	return names, nil
}

// GetActiveProject implements repositories.ProjectManagementRepository.
func (m *MockProjectManagementRepository) GetActiveProject() (string, error) {
	if m.GetActiveProjectFunc != nil {
		return m.GetActiveProjectFunc()
	}
	// Default implementation: return activeProject
	return m.activeProject, nil
}

// SetActiveProject implements repositories.ProjectManagementRepository.
func (m *MockProjectManagementRepository) SetActiveProject(projectName string) error {
	if m.SetActiveProjectFunc != nil {
		return m.SetActiveProjectFunc(projectName)
	}
	// Default implementation: set activeProject
	m.activeProject = projectName
	return nil
}

// DeleteProject implements repositories.ProjectManagementRepository.
func (m *MockProjectManagementRepository) DeleteProject(projectName string) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(projectName)
	}
	// Default implementation: delete from memory
	delete(m.projects, projectName)
	// Clear active project if it was deleted
	if m.activeProject == projectName {
		m.activeProject = ""
	}
	return nil
}

// GetProjectInfo implements repositories.ProjectManagementRepository.
func (m *MockProjectManagementRepository) GetProjectInfo(projectName string) (map[string]string, error) {
	if m.GetProjectInfoFunc != nil {
		return m.GetProjectInfoFunc(projectName)
	}
	// Default implementation: return from storage
	info, exists := m.projects[projectName]
	if !exists {
		return nil, nil
	}
	return info, nil
}

// ProjectExists implements repositories.ProjectManagementRepository.
func (m *MockProjectManagementRepository) ProjectExists(projectName string) (bool, error) {
	if m.ProjectExistsFunc != nil {
		return m.ProjectExistsFunc(projectName)
	}
	// Default implementation: check projects map
	_, exists := m.projects[projectName]
	return exists, nil
}
