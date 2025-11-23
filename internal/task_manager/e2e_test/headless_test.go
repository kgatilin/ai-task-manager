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

// HeadlessTestSuite tests headless build functionality
// This suite builds the headless binary and verifies it works correctly
type HeadlessTestSuite struct {
	E2ETestSuite
	headlessBinaryPath string // Path to built headless binary
}

// SetupSuite runs once per test suite
func (s *HeadlessTestSuite) SetupSuite() {
	// Base setup (project creation, roadmap init, working directory)
	s.E2ETestSuite.SetupSuite()

	// Build headless binary with build tag
	// Get absolute path to project root (4 levels up from e2e_test directory)
	// internal/task_manager/e2e_test -> ../../../../
	relativeRoot := "../../../"
	projectRoot, err := filepath.Abs(relativeRoot)
	s.Require().NoError(err, "failed to resolve project root")

	s.headlessBinaryPath = filepath.Join(projectRoot, "tm-headless-e2e-test")

	buildCmd := exec.Command("go", "build", "-tags", "headless", "-o", s.headlessBinaryPath, "./cmd/tm")
	buildCmd.Dir = projectRoot
	output, err := buildCmd.CombinedOutput()
	s.Require().NoError(err, "failed to build headless binary\nOutput: %s", string(output))

	// Verify headless binary exists
	_, err = os.Stat(s.headlessBinaryPath)
	s.Require().NoError(err, "headless binary not found at %s", s.headlessBinaryPath)
}

// TearDownSuite runs once after all tests in the suite
func (s *HeadlessTestSuite) TearDownSuite() {
	// Clean up headless binary
	if s.headlessBinaryPath != "" {
		os.Remove(s.headlessBinaryPath)
	}
}

// runHeadless executes a command with the headless binary
func (s *HeadlessTestSuite) runHeadless(args ...string) (string, error) {
	cmd := exec.Command(s.headlessBinaryPath, args...)
	cmd.Env = append(os.Environ(), "TM_WORKING_DIR="+s.testWorkingDir)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// TestHeadlessBinaryExists verifies the headless binary was built successfully
func (s *HeadlessTestSuite) TestHeadlessBinaryExists() {
	_, err := os.Stat(s.headlessBinaryPath)
	s.NoError(err, "headless binary should exist")
}

// TestHeadlessBinarySizeReduction verifies headless binary is smaller than full binary
func (s *HeadlessTestSuite) TestHeadlessBinarySizeReduction() {
	// Build full binary if it doesn't exist
	relativeRoot := "../../../"
	projectRoot, err := filepath.Abs(relativeRoot)
	s.Require().NoError(err, "failed to resolve project root")

	fullBinaryPath := filepath.Join(projectRoot, "tm-full-e2e-test")
	defer os.Remove(fullBinaryPath)

	buildCmd := exec.Command("go", "build", "-o", fullBinaryPath, "./cmd/tm")
	buildCmd.Dir = projectRoot
	output, err := buildCmd.CombinedOutput()
	s.Require().NoError(err, "failed to build full binary\nOutput: %s", string(output))

	// Get file sizes
	fullInfo, err := os.Stat(fullBinaryPath)
	s.Require().NoError(err, "failed to stat full binary")

	headlessInfo, err := os.Stat(s.headlessBinaryPath)
	s.Require().NoError(err, "failed to stat headless binary")

	fullSize := fullInfo.Size()
	headlessSize := headlessInfo.Size()

	// Calculate reduction percentage
	reduction := float64(fullSize-headlessSize) / float64(fullSize) * 100

	s.T().Logf("Full binary size: %d bytes", fullSize)
	s.T().Logf("Headless binary size: %d bytes", headlessSize)
	s.T().Logf("Size reduction: %.1f%%", reduction)

	// Verify at least 30% reduction (AC-670)
	s.GreaterOrEqual(reduction, 30.0, "headless binary should be at least 30%% smaller than full binary")
	s.Less(headlessSize, fullSize, "headless binary should be smaller than full binary")
}

// TestHeadlessProjectCommand verifies project commands work in headless build
func (s *HeadlessTestSuite) TestHeadlessProjectCommand() {
	// Create a new project using headless binary
	projectName := "headless-test-project-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	output, err := s.runHeadless("project", "create", projectName)
	s.NoError(err, "project create should work in headless binary\nOutput:\n%s", output)
	s.Contains(output, projectName, "output should contain project name")
}

// TestHeadlessTrackCommand verifies track commands work in headless build
func (s *HeadlessTestSuite) TestHeadlessTrackCommand() {
	// Create a track using headless binary
	output, err := s.runHeadless("track", "create", "--title", "Headless Test Track", "--description", "Test")
	s.NoError(err, "track create should work in headless binary\nOutput:\n%s", output)

	// Verify track was created
	s.Contains(output, "Headless Test Track", "output should contain track title")

	// List tracks
	output, err = s.runHeadless("track", "list")
	s.NoError(err, "track list should work in headless binary\nOutput:\n%s", output)
	s.Contains(output, "Headless Test Track", "track list should show created track")
}

// TestHeadlessTaskCommand verifies task commands work in headless build
func (s *HeadlessTestSuite) TestHeadlessTaskCommand() {
	// Create a track first
	output, err := s.runHeadless("track", "create", "--title", "Task Test Track", "--description", "Test")
	s.NoError(err, "track create should work\nOutput:\n%s", output)

	trackID := s.parseID(output, "-track-")
	s.NotEmpty(trackID, "should extract track ID from output")

	// Create a task using headless binary
	output, err = s.runHeadless("task", "create", "--track", trackID, "--title", "Headless Test Task")
	s.NoError(err, "task create should work in headless binary\nOutput:\n%s", output)
	s.Contains(output, "Headless Test Task", "output should contain task title")
}

// TestHeadlessIterationCommand verifies iteration commands work in headless build
func (s *HeadlessTestSuite) TestHeadlessIterationCommand() {
	// Create iteration using headless binary
	output, err := s.runHeadless("iteration", "create", "--name", "Headless Test Iteration", "--goal", "Test", "--deliverable", "Test Deliverable")
	s.NoError(err, "iteration create should work in headless binary\nOutput:\n%s", output)
	s.Contains(output, "Headless Test Iteration", "output should contain iteration name")

	// List iterations
	output, err = s.runHeadless("iteration", "list")
	s.NoError(err, "iteration list should work in headless binary\nOutput:\n%s", output)
}

// TestHeadlessACCommand verifies AC commands work in headless build
func (s *HeadlessTestSuite) TestHeadlessACCommand() {
	// Create a task for AC
	trackOutput, err := s.runHeadless("track", "create", "--title", "AC Test Track", "--description", "Test")
	s.NoError(err, "track create should work\nOutput:\n%s", trackOutput)

	trackID := s.parseID(trackOutput, "-track-")
	s.NotEmpty(trackID, "should extract track ID")

	taskOutput, err := s.runHeadless("task", "create", "--track", trackID, "--title", "AC Test Task")
	s.NoError(err, "task create should work\nOutput:\n%s", taskOutput)

	taskID := s.parseID(taskOutput, "-task-")
	s.NotEmpty(taskID, "should extract task ID")

	// Add AC using headless binary
	output, err := s.runHeadless("ac", "add", taskID, "--description", "Headless AC Test", "--testing-instructions", "1. Test step")
	s.NoError(err, "ac add should work in headless binary\nOutput:\n%s", output)

	// List ACs
	output, err = s.runHeadless("ac", "list", taskID)
	s.NoError(err, "ac list should work in headless binary\nOutput:\n%s", output)
	s.Contains(output, "Headless AC Test", "ac list should show created AC")
}

// TestHeadlessUICommandFails verifies ui command is unavailable in headless build
func (s *HeadlessTestSuite) TestHeadlessUICommandFails() {
	// Run ui command with headless binary - should fail
	output, err := s.runHeadless("ui")

	// Command should fail with error
	s.Error(err, "ui command should fail in headless binary")

	// Error message should mention headless or TUI not available
	outputLower := strings.ToLower(output)
	s.True(
		strings.Contains(outputLower, "headless") || strings.Contains(outputLower, "not available") || strings.Contains(outputLower, "tui"),
		fmt.Sprintf("error message should mention headless or TUI not available. Got: %s", output),
	)
}

// TestHeadlessMultipleCommands verifies diverse CLI commands in headless build
func (s *HeadlessTestSuite) TestHeadlessMultipleCommands() {
	// Test roadmap commands
	output, err := s.runHeadless("roadmap", "show")
	s.NoError(err, "roadmap show should work\nOutput:\n%s", output)

	// Test project show
	output, err = s.runHeadless("project", "show")
	s.NoError(err, "project show should work\nOutput:\n%s", output)

	// Test help command
	output, err = s.runHeadless("--help")
	s.NoError(err, "help should work\nOutput:\n%s", output)
	s.Contains(output, "tm", "help output should show tm usage")

	// Test version
	output, err = s.runHeadless("--version")
	s.NoError(err, "version should work\nOutput:\n%s", output)
	s.Contains(output, "version", "version output should show version")
}

func TestHeadlessSuite(t *testing.T) {
	suite.Run(t, new(HeadlessTestSuite))
}
