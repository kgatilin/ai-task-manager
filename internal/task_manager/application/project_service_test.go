package application_test

import (
	"fmt"
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/mocks"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
)

// setupProjectTestService creates a test service with mock repository
func setupProjectTestService(t *testing.T) (*application.ProjectApplicationService, *mocks.MockProjectManagementRepository) {
	mockRepo := mocks.NewMockProjectManagementRepository()
	validationService := services.NewValidationService()
	service := application.NewProjectService(mockRepo, validationService)
	return service, mockRepo
}

// TestNewProjectService tests the constructor
func TestNewProjectService(t *testing.T) {
	mockRepo := mocks.NewMockProjectManagementRepository()
	validationService := services.NewValidationService()

	service := application.NewProjectService(mockRepo, validationService)

	if service == nil {
		t.Fatal("NewProjectService returned nil")
	}
}

// TestCreateProject_Success tests successful project creation
func TestCreateProject_Success(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Test with explicit project code
	_, err := service.CreateProject("my-project", "MP")
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	// Verify project was created
	exists, _ := mockRepo.ProjectExists("my-project")
	if !exists {
		t.Error("Project was not created")
	}
}

// TestCreateProject_WithDefaultCode tests project creation with auto-generated code
func TestCreateProject_WithDefaultCode(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Test with empty project code (should auto-generate)
	_, err := service.CreateProject("test-project", "")
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	// Verify project was created
	exists, _ := mockRepo.ProjectExists("test-project")
	if !exists {
		t.Error("Project was not created")
	}
}

// TestCreateProject_InvalidProjectName tests validation of project name
func TestCreateProject_InvalidProjectName(t *testing.T) {
	service, _ := setupProjectTestService(t)

	tests := []struct {
		name        string
		projectName string
		wantErr     bool
	}{
		{
			name:        "empty name",
			projectName: "",
			wantErr:     true,
		},
		{
			name:        "spaces in name",
			projectName: "my project",
			wantErr:     true,
		},
		{
			name:        "special characters",
			projectName: "my@project",
			wantErr:     true,
		},
		{
			name:        "valid name with hyphens",
			projectName: "my-project",
			wantErr:     false,
		},
		{
			name:        "valid name with underscores",
			projectName: "my_project",
			wantErr:     false,
		},
		{
			name:        "valid alphanumeric",
			projectName: "project123",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateProject(tt.projectName, "TEST")
			if tt.wantErr && err == nil {
				t.Error("Expected error for invalid project name")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestCreateProject_InvalidProjectCode tests validation of project code
func TestCreateProject_InvalidProjectCode(t *testing.T) {
	service, _ := setupProjectTestService(t)

	tests := []struct {
		name        string
		projectCode string
		wantErr     bool
	}{
		{
			name:        "lowercase code",
			projectCode: "mp",
			wantErr:     true,
		},
		{
			name:        "code with special characters",
			projectCode: "M-P",
			wantErr:     true,
		},
		{
			name:        "code with spaces",
			projectCode: "M P",
			wantErr:     true,
		},
		{
			name:        "valid uppercase code",
			projectCode: "MP",
			wantErr:     false,
		},
		{
			name:        "valid code with numbers",
			projectCode: "MP123",
			wantErr:     false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use unique valid project name for each test (testing project code, not name)
			projectName := fmt.Sprintf("test-project-%d", i)
			_, err := service.CreateProject(projectName, tt.projectCode)
			if tt.wantErr && err == nil {
				t.Error("Expected error for invalid project code")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestCreateProject_AlreadyExists tests duplicate project name
func TestCreateProject_AlreadyExists(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create first project
	_, err := service.CreateProject("my-project", "MP")
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	// Try to create duplicate
	_, err = service.CreateProject("my-project", "MP2")
	if err == nil {
		t.Error("Expected error for duplicate project name")
	}

	// Verify only one project exists
	projects, _ := mockRepo.ListProjects()
	if len(projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects))
	}
}

// TestCreateProject_RepositoryError tests error handling from repository
func TestCreateProject_RepositoryError(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Configure mock to return error on CreateProject
	mockRepo.CreateProjectFunc = func(projectName, projectCode string) error {
		return fmt.Errorf("database error")
	}

	_, err := service.CreateProject("my-project", "MP")
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// TestListProjects tests listing all projects
func TestListProjects(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create multiple projects
	_ = mockRepo.CreateProject("project-1", "P1")
	_ = mockRepo.CreateProject("project-2", "P2")
	_ = mockRepo.CreateProject("project-3", "P3")

	projects, err := service.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	if len(projects) != 3 {
		t.Errorf("Expected 3 projects, got %d", len(projects))
	}
}

// TestListProjects_Empty tests listing with no projects
func TestListProjects_Empty(t *testing.T) {
	service, _ := setupProjectTestService(t)

	projects, err := service.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(projects))
	}
}

// TestListProjects_RepositoryError tests error handling
func TestListProjects_RepositoryError(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	mockRepo.ListProjectsFunc = func() ([]string, error) {
		return nil, fmt.Errorf("database error")
	}

	_, err := service.ListProjects()
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// TestGetActiveProject tests getting the active project
func TestGetActiveProject(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create and set active project
	_ = mockRepo.CreateProject("my-project", "MP")
	_ = mockRepo.SetActiveProject("my-project")

	activeProject, err := service.GetActiveProject()
	if err != nil {
		t.Fatalf("GetActiveProject failed: %v", err)
	}

	if activeProject != "my-project" {
		t.Errorf("Expected active project 'my-project', got '%s'", activeProject)
	}
}

// TestGetActiveProject_RepositoryError tests error handling
func TestGetActiveProject_RepositoryError(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	mockRepo.GetActiveProjectFunc = func() (string, error) {
		return "", fmt.Errorf("database error")
	}

	_, err := service.GetActiveProject()
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// TestSwitchProject_Success tests successful project switch
func TestSwitchProject_Success(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create two projects
	_ = mockRepo.CreateProject("project-1", "P1")
	_ = mockRepo.CreateProject("project-2", "P2")
	_ = mockRepo.SetActiveProject("project-1")

	// Switch to project-2
	err := service.SwitchProject("project-2")
	if err != nil {
		t.Fatalf("SwitchProject failed: %v", err)
	}

	// Verify active project changed
	activeProject, _ := mockRepo.GetActiveProject()
	if activeProject != "project-2" {
		t.Errorf("Expected active project 'project-2', got '%s'", activeProject)
	}
}

// TestSwitchProject_NotFound tests switching to non-existent project
func TestSwitchProject_NotFound(t *testing.T) {
	service, _ := setupProjectTestService(t)

	err := service.SwitchProject("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

// TestSwitchProject_RepositoryError tests error handling
func TestSwitchProject_RepositoryError(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create project
	_ = mockRepo.CreateProject("my-project", "MP")

	// Configure mock to return error on SetActiveProject
	mockRepo.SetActiveProjectFunc = func(projectName string) error {
		return fmt.Errorf("database error")
	}

	err := service.SwitchProject("my-project")
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// TestDeleteProject_Success tests successful project deletion
func TestDeleteProject_Success(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create two projects
	_ = mockRepo.CreateProject("my-project", "MP")
	_ = mockRepo.CreateProject("other-project", "OP")

	// Set other-project as active (so my-project can be deleted)
	_ = mockRepo.SetActiveProject("other-project")

	err := service.DeleteProject("my-project")
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}

	// Verify project was deleted
	exists, _ := mockRepo.ProjectExists("my-project")
	if exists {
		t.Error("Project should have been deleted")
	}
}

// TestDeleteProject_NotFound tests deleting non-existent project
func TestDeleteProject_NotFound(t *testing.T) {
	service, _ := setupProjectTestService(t)

	err := service.DeleteProject("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

// TestDeleteProject_RepositoryError tests error handling
func TestDeleteProject_RepositoryError(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create project
	_ = mockRepo.CreateProject("my-project", "MP")

	// Configure mock to return error on DeleteProject
	mockRepo.DeleteProjectFunc = func(projectName string) error {
		return fmt.Errorf("database error")
	}

	err := service.DeleteProject("my-project")
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// TestGetProjectInfo_Success tests getting project info
func TestGetProjectInfo_Success(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create project
	_ = mockRepo.CreateProject("my-project", "MP")

	info, err := service.GetProjectInfo("my-project")
	if err != nil {
		t.Fatalf("GetProjectInfo failed: %v", err)
	}

	if info == nil {
		t.Fatal("Expected project info, got nil")
	}

	if info["name"] != "my-project" {
		t.Errorf("Expected project name 'my-project', got '%s'", info["name"])
	}

	if info["code"] != "MP" {
		t.Errorf("Expected project code 'MP', got '%s'", info["code"])
	}
}

// TestGetProjectInfo_NotFound tests getting info for non-existent project
func TestGetProjectInfo_NotFound(t *testing.T) {
	service, _ := setupProjectTestService(t)

	_, err := service.GetProjectInfo("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}

// TestGetProjectInfo_RepositoryError tests error handling
func TestGetProjectInfo_RepositoryError(t *testing.T) {
	service, mockRepo := setupProjectTestService(t)

	// Create project
	_ = mockRepo.CreateProject("my-project", "MP")

	// Configure mock to return error on GetProjectInfo
	mockRepo.GetProjectInfoFunc = func(projectName string) (map[string]string, error) {
		return nil, fmt.Errorf("database error")
	}

	_, err := service.GetProjectInfo("my-project")
	if err == nil {
		t.Error("Expected error from repository")
	}
}
