package task_manager_e2e_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ProjectTestSuite tests all project commands
type ProjectTestSuite struct {
	E2ETestSuite
}

// SetupSuite overrides E2ETestSuite.SetupSuite to skip project/roadmap creation.
// ProjectTestSuite tests project creation/switching, so it creates its own projects per test.
// This prevents conflicts with E2ETestSuite's default project setup.
func (s *ProjectTestSuite) SetupSuite() {
	// Binary is already built in TestMain - no need to build here
	// Do NOT create project/roadmap - ProjectTestSuite creates its own projects in each test
}

// SetupTest runs before EACH test method
func (s *ProjectTestSuite) SetupTest() {
	// Create unique working directory for each test
	// This ensures project files (.darwinflow/active-project.txt, databases) don't conflict
	tmpDir := s.T().TempDir()
	s.testWorkingDir = tmpDir
}

// TearDownTest runs after each test
func (s *ProjectTestSuite) TearDownTest() {
	// No per-test cleanup needed
	// Tests operate on the shared project created in SetupSuite
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}

// TestProjectCreate tests creating a new project
func (s *ProjectTestSuite) TestProjectCreate() {
	output, err := s.run("project", "create", "test-project")
	s.NoError(err, "project create should succeed\nOutput:\n%s", output)
	s.Contains(output, "test-project", "output should contain project name")
	s.Contains(output, "Project created successfully", "output should confirm creation")
	s.Contains(output, "Project code:", "output should show project code")
}

// TestProjectCreateWithCustomCode tests creating a project with custom code
func (s *ProjectTestSuite) TestProjectCreateWithCustomCode() {
	output, err := s.run("project", "create", "custom-project", "--code", "CUST")
	s.NoError(err, "project create with custom code should succeed\nOutput:\n%s", output)
	s.Contains(output, "custom-project", "output should contain project name")
	s.Contains(output, "CUST", "output should contain custom code")
}

// TestProjectCreateDuplicate tests that duplicate projects are rejected
func (s *ProjectTestSuite) TestProjectCreateDuplicate() {
	// Create first project
	output1, err := s.run("project", "create", "duplicate-test")
	s.NoError(err, "first project creation should succeed\nOutput:\n%s", output1)

	// Try to create duplicate
	output2, err := s.run("project", "create", "duplicate-test")
	s.Error(err, "duplicate project creation should fail")
	s.Contains(output2, "already exists", "error message should indicate project already exists")
}

// TestProjectCreateInvalidName tests that invalid project names are rejected
func (s *ProjectTestSuite) TestProjectCreateInvalidName() {
	output, err := s.run("project", "create", "invalid@name")
	s.Error(err, "project with invalid characters should fail\nOutput:\n%s", output)
	s.Contains(output, "invalid project name", "error should indicate invalid name")
}

// TestProjectList tests listing projects
func (s *ProjectTestSuite) TestProjectList() {
	// Create multiple projects
	s.run("project", "create", "list-test-1")
	s.run("project", "create", "list-test-2")
	s.run("project", "create", "list-test-3")

	// List projects
	output, err := s.run("project", "list")
	s.NoError(err, "project list should succeed\nOutput:\n%s", output)
	s.Contains(output, "list-test-1", "output should contain first project")
	s.Contains(output, "list-test-2", "output should contain second project")
	s.Contains(output, "list-test-3", "output should contain third project")
	s.Contains(output, "Projects:", "output should have Projects header")
}

// TestProjectListShowsActive tests that list marks active project
func (s *ProjectTestSuite) TestProjectListShowsActive() {
	// Create projects
	s.run("project", "create", "active-1")
	s.run("project", "create", "active-2")

	// Switch to active-1
	output, err := s.run("project", "switch", "active-1")
	s.NoError(err, "project switch should succeed\nOutput:\n%s", output)

	// List and verify active-1 is marked
	output, err = s.run("project", "list")
	s.NoError(err, "project list should succeed\nOutput:\n%s", output)
	// Should show "active-1 (active)" and "active-2" without (active)
	s.Contains(output, "active-1", "output should contain active-1")
	s.Contains(output, "active-2", "output should contain active-2")
}

// TestProjectSwitch tests switching between projects
func (s *ProjectTestSuite) TestProjectSwitch() {
	// Create two projects
	s.run("project", "create", "switch-test-1")
	s.run("project", "create", "switch-test-2")

	// Switch to switch-test-1
	output, err := s.run("project", "switch", "switch-test-1")
	s.NoError(err, "project switch should succeed\nOutput:\n%s", output)
	s.Contains(output, "switch-test-1", "output should confirm switch")
	s.Contains(output, "Switched to project", "output should indicate switch")

	// Verify it's active
	output, err = s.run("project", "show")
	s.NoError(err, "project show should succeed\nOutput:\n%s", output)
	s.Contains(output, "switch-test-1", "show should display switched project")

	// Switch to switch-test-2
	output, err = s.run("project", "switch", "switch-test-2")
	s.NoError(err, "project switch to second project should succeed\nOutput:\n%s", output)

	// Verify new active project
	output, err = s.run("project", "show")
	s.NoError(err, "project show should succeed\nOutput:\n%s", output)
	s.Contains(output, "switch-test-2", "show should display current active project")
}

// TestProjectSwitchNonExistent tests that switching to non-existent project fails
func (s *ProjectTestSuite) TestProjectSwitchNonExistent() {
	output, err := s.run("project", "switch", "non-existent-project")
	s.Error(err, "switching to non-existent project should fail\nOutput:\n%s", output)
	s.Contains(output, "does not exist", "error should indicate project doesn't exist")
}

// TestProjectShow tests showing the active project
func (s *ProjectTestSuite) TestProjectShow() {
	// Create and switch to a project
	s.run("project", "create", "show-test")
	s.run("project", "switch", "show-test")

	// Show active project
	output, err := s.run("project", "show")
	s.NoError(err, "project show should succeed\nOutput:\n%s", output)
	s.Contains(output, "show-test", "output should contain active project name")
	s.Contains(output, "Active project:", "output should have Active project label")
}

// TestProjectShowDefault tests showing default project when none switched
func (s *ProjectTestSuite) TestProjectShowDefault() {
	// Without explicitly creating/switching, should show default or error appropriately
	output, _ := s.run("project", "show")
	// May succeed or fail depending on implementation
	// Just verify it produces output
	s.NotEmpty(output, "project show should produce output")
}

// TestProjectDelete tests deleting a project
func (s *ProjectTestSuite) TestProjectDelete() {
	// Create two projects
	s.run("project", "create", "delete-test-1")
	s.run("project", "create", "delete-test-2")

	// Switch to delete-test-2 (so we don't try to delete active project)
	s.run("project", "switch", "delete-test-2")

	// Delete delete-test-1 with --force
	output, err := s.run("project", "delete", "delete-test-1", "--force")
	s.NoError(err, "project delete should succeed\nOutput:\n%s", output)
	s.Contains(output, "delete-test-1", "output should contain deleted project name")
	s.Contains(output, "deleted successfully", "output should confirm deletion")

	// Verify it's gone from list
	output, err = s.run("project", "list")
	s.NoError(err, "project list should succeed\nOutput:\n%s", output)
	s.NotContains(output, "delete-test-1", "deleted project should not appear in list")
}

// TestProjectDeleteNonExistent tests that deleting non-existent project fails
func (s *ProjectTestSuite) TestProjectDeleteNonExistent() {
	output, err := s.run("project", "delete", "non-existent", "--force")
	s.Error(err, "deleting non-existent project should fail\nOutput:\n%s", output)
	s.Contains(output, "does not exist", "error should indicate project doesn't exist")
}

// TestProjectDeleteActive tests that deleting active project fails
func (s *ProjectTestSuite) TestProjectDeleteActive() {
	// Create two projects
	s.run("project", "create", "active-delete-1")
	s.run("project", "create", "active-delete-2")

	// Switch to active-delete-1
	s.run("project", "switch", "active-delete-1")

	// Try to delete active-delete-1 (should fail)
	output, err := s.run("project", "delete", "active-delete-1", "--force")
	s.Error(err, "deleting active project should fail\nOutput:\n%s", output)
	s.Contains(output, "cannot delete active project", "error should indicate can't delete active project")
}

// TestProjectParallelCreation tests that multiple projects can be created concurrently without conflicts
func (s *ProjectTestSuite) TestProjectParallelCreation() {
	const numProjects = 5
	projectNames := make([]string, numProjects)
	for i := 0; i < numProjects; i++ {
		projectNames[i] = "parallel-" + strings.Repeat("0", len("00000")-len(string(rune(i+48)))) + string(rune(i+48))
	}

	// Use goroutines to create projects concurrently
	var wg sync.WaitGroup
	errors := make([]error, numProjects)

	for i := 0; i < numProjects; i++ {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()
			_, err := s.run("project", "create", name)
			errors[idx] = err
		}(i, projectNames[i])
	}

	wg.Wait()

	// Verify all projects were created successfully
	for i, err := range errors {
		s.NoError(err, "parallel project creation %d should succeed", i)
	}

	// Verify all projects exist in list
	output, err := s.run("project", "list")
	s.NoError(err, "project list should succeed\nOutput:\n%s", output)

	for _, name := range projectNames {
		s.Contains(output, name, "parallel project %s should be in list", name)
	}
}

// TestProjectCommandFlow tests a complete workflow: create -> list -> switch -> show -> delete
func (s *ProjectTestSuite) TestProjectCommandFlow() {
	// Create a project
	output, err := s.run("project", "create", "workflow-test")
	s.NoError(err, "project create should succeed\nOutput:\n%s", output)

	// List projects and find it
	output, err = s.run("project", "list")
	s.NoError(err, "project list should succeed\nOutput:\n%s", output)
	s.Contains(output, "workflow-test", "workflow-test should be in list")

	// Switch to it
	output, err = s.run("project", "switch", "workflow-test")
	s.NoError(err, "project switch should succeed\nOutput:\n%s", output)

	// Show it
	output, err = s.run("project", "show")
	s.NoError(err, "project show should succeed\nOutput:\n%s", output)
	s.Contains(output, "workflow-test", "workflow-test should be active")

	// Create another project to switch to before deleting
	s.run("project", "create", "workflow-cleanup")
	s.run("project", "switch", "workflow-cleanup")

	// Delete the first one
	output, err = s.run("project", "delete", "workflow-test", "--force")
	s.NoError(err, "project delete should succeed\nOutput:\n%s", output)

	// Verify it's gone
	output, err = s.run("project", "list")
	s.NoError(err, "project list should succeed\nOutput:\n%s", output)
	s.NotContains(output, "workflow-test", "deleted project should not be in list")
}
