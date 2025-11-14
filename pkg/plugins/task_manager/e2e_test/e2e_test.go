package task_manager_e2e_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	// Binary name for E2E tests - built once in TestMain
	dwBinaryName = "dw-e2e-test"
)

var (
	// Path to built binary - set once in TestMain, immutable thereafter
	dwBinaryPath string
)

// TestMain runs once before all tests in the package
// It builds the dw binary from source and cleans it up afterward
func TestMain(m *testing.M) {
	// Build the dw binary from source before running tests
	// This ensures we're testing the current codebase, not a system-installed version

	// Get absolute path to project root (4 levels up from e2e_test directory)
	// pkg/plugins/task_manager/e2e_test -> ../../../../
	relativeRoot := "../../../../"
	projectRoot, err := filepath.Abs(relativeRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve project root path: %v\n", err)
		os.Exit(1)
	}

	dwBinaryPath = filepath.Join(projectRoot, dwBinaryName)

	buildCmd := exec.Command("go", "build", "-o", dwBinaryPath, "./cmd/dw")
	buildCmd.Dir = projectRoot
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build dw binary for E2E tests\nOutput: %s\n", string(output))
		os.Exit(1)
	}

	// Verify binary was created
	_, err = os.Stat(dwBinaryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dw binary not found at %s after build: %v\n", dwBinaryPath, err)
		os.Exit(1)
	}

	// Run all tests
	exitCode := m.Run()

	// Cleanup: remove the built binary
	os.Remove(dwBinaryPath)

	os.Exit(exitCode)
}

// E2ETestSuite is the base test suite for all E2E tests
type E2ETestSuite struct {
	suite.Suite
	projectName    string // Unique project name for this test suite
	testWorkingDir string // Consistent working directory for all commands
}

// SetupSuite runs ONCE per test suite before any tests
func (s *E2ETestSuite) SetupSuite() {
	// Binary is already built in TestMain - no need to build here

	// Generate unique project name for this suite to avoid conflicts
	s.projectName = "e2e-test-" + strconv.FormatInt(time.Now().UnixNano(), 36)

	// Create dedicated working directory for this test suite
	// This ensures all CLI commands run from the same working directory,
	// which is critical for .darwinflow/active-project.txt persistence
	tmpDir := os.TempDir()
	s.testWorkingDir = filepath.Join(tmpDir, "dw-e2e-wd-"+s.projectName)
	err := os.MkdirAll(s.testWorkingDir, 0755)
	s.Require().NoError(err, "failed to create test working directory")

	// Create project once for entire suite
	cmdOutput, err := s.run("project", "create", s.projectName)
	s.Require().NoError(err, "failed to create project for suite\nOutput: %s", cmdOutput)

	// Switch to the project
	cmdOutput, err = s.run("project", "switch", s.projectName)
	s.Require().NoError(err, "failed to switch to project\nOutput: %s", cmdOutput)

	// Initialize roadmap once for entire suite
	cmdOutput, err = s.run("roadmap", "init", "--vision", "E2E Test Suite Vision", "--success-criteria", "All E2E tests pass")
	s.Require().NoError(err, "failed to initialize roadmap\nOutput: %s", cmdOutput)
}

// SetupTest runs before EACH test method
func (s *E2ETestSuite) SetupTest() {
	// No per-test setup needed
	// Project and roadmap are shared across all tests in the suite (created in SetupSuite)
	// Test isolation is provided by the unique project name per suite
}

// TearDownTest runs after each test
func (s *E2ETestSuite) TearDownTest() {
	// No per-test cleanup needed
	// Tests operate on the shared project created in SetupSuite
}

// run executes a dw task-manager command and returns stdout/stderr combined
func (s *E2ETestSuite) run(args ...string) (string, error) {
	// Prepend "task-manager" to the arguments
	fullArgs := append([]string{"task-manager"}, args...)

	// Use the locally-built binary, not system PATH
	cmd := exec.Command(dwBinaryPath, fullArgs...)

	// CRITICAL: Set DARWINFLOW_WORKING_DIR environment variable
	// This ensures the binary uses a consistent working directory for all operations
	// (.darwinflow/active-project.txt, project databases, etc.)
	// Without this, each command invocation would use its own os.Getwd() which may vary
	cmd.Env = append(os.Environ(), "DARWINFLOW_WORKING_DIR="+s.testWorkingDir)

	// Execute the command and capture output
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// requireSuccess asserts that a command executed successfully
func (s *E2ETestSuite) requireSuccess(output string, err error, msg string, args ...interface{}) {
	s.Require().NoError(err, append([]interface{}{msg, "\nOutput:\n", output}, args...)...)
}

// requireError asserts that a command failed with an error
func (s *E2ETestSuite) requireError(err error, msg string, args ...interface{}) {
	s.Require().Error(err, append([]interface{}{msg}, args...)...)
}

// parseID extracts an entity ID from command output
// For example, extracts "ABC-track-1" from "ID: ABC-track-1" or "Created track: ABC-track-1"
// The prefix parameter can be "-track-" or "track" (both formats are accepted)
func (s *E2ETestSuite) parseID(output string, prefix string) string {
	// Normalize prefix to "-entity-" format if it's not already
	entityPattern := prefix
	if !strings.Contains(entityPattern, "-") || !strings.HasPrefix(entityPattern, "-") {
		entityPattern = "-" + entityPattern + "-"
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Look for lines containing the entity pattern (e.g., "-track-", "-task-", "-ac-")
		if strings.Contains(trimmedLine, entityPattern) {
			// Extract the ID - look for pattern like "XXX-track-1"
			// Try colon-separated format first (ID: XXX-track-1)
			if strings.Contains(trimmedLine, ":") {
				parts := strings.Split(trimmedLine, ":")
				if len(parts) > 1 {
					idPart := strings.TrimSpace(parts[len(parts)-1])
					if strings.Contains(idPart, entityPattern) {
						return idPart
					}
				}
			}

			// Try space-separated format
			parts := strings.Fields(trimmedLine)
			for _, part := range parts {
				// Remove trailing punctuation (commas, colons, etc.)
				cleanPart := strings.Trim(part, ",;:")
				if strings.Contains(cleanPart, entityPattern) {
					return cleanPart
				}
			}
		}
	}
	return ""
}

// parseIterationNumber extracts iteration number from command output
// For example, extracts "3" from "Number: 3"
func (s *E2ETestSuite) parseIterationNumber(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Look for "Number: X" pattern
		if strings.Contains(trimmedLine, "Number:") {
			parts := strings.Split(trimmedLine, ":")
			if len(parts) > 1 {
				number := strings.TrimSpace(parts[1])
				return number
			}
		}
	}
	return ""
}
