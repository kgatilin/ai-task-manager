package task_manager_e2e_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TaskTestSuite tests task commands end-to-end
type TaskTestSuite struct {
	E2ETestSuite
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}

// TestTaskCreate tests creating a task with valid track
func (s *TaskTestSuite) TestTaskCreate() {
	// Create track first (tasks require track)
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")
	s.NotEmpty(trackID, "track ID should be extracted")

	// Create task
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	taskID := s.parseID(taskOutput, "task")
	s.NotEmpty(taskID, "task ID should be extracted")
}

// TestTaskCreateWithoutTrack tests that task creation fails without track
func (s *TaskTestSuite) TestTaskCreateWithoutTrack() {
	taskOutput, err := s.run("task", "create", "--title", "Test Task")
	s.requireError(err, "task creation without track should fail")
	s.NotEmpty(taskOutput, "error output should be provided")
}

// TestTaskList tests listing tasks
func (s *TaskTestSuite) TestTaskList() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create task
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	taskID := s.parseID(taskOutput, "task")

	// List all tasks
	listOutput, err := s.run("task", "list")
	s.requireSuccess(listOutput, err, "failed to list tasks")
	s.Contains(listOutput, taskID, "created task should appear in list")
}

// TestTaskListByStatus tests filtering tasks by status
func (s *TaskTestSuite) TestTaskListByStatus() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create task with todo status
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")

	// List tasks with status=todo filter
	listOutput, err := s.run("task", "list", "--status", "todo")
	s.requireSuccess(listOutput, err, "failed to list tasks with status filter")
	// Should contain at least one task
	s.NotEmpty(listOutput, "task list output should not be empty")
}

// TestTaskShow tests showing task details
func (s *TaskTestSuite) TestTaskShow() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create task
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	taskID := s.parseID(taskOutput, "task")

	// Show task details
	showOutput, err := s.run("task", "show", taskID)
	s.requireSuccess(showOutput, err, "failed to show task")
	s.Contains(showOutput, taskID, "task ID should appear in show output")
	s.Contains(showOutput, "Test Task", "task title should appear in show output")
}

// TestTaskUpdate tests updating task properties
func (s *TaskTestSuite) TestTaskUpdate() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create task
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	taskID := s.parseID(taskOutput, "task")

	// Update task title
	updateOutput, err := s.run("task", "update", taskID, "--title", "Updated Task Title")
	s.requireSuccess(updateOutput, err, "failed to update task title")

	// Verify update
	showOutput, err := s.run("task", "show", taskID)
	s.requireSuccess(showOutput, err, "failed to show updated task")
	s.Contains(showOutput, "Updated Task Title", "updated title should appear in show output")
}

// TestTaskStatusTransition tests status transitions (todo → in-progress → done)
func (s *TaskTestSuite) TestTaskStatusTransition() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create task (starts in "todo" status)
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	taskID := s.parseID(taskOutput, "task")

	// Transition to in-progress
	updateOutput, err := s.run("task", "update", taskID, "--status", "in-progress")
	s.requireSuccess(updateOutput, err, "failed to transition task to in-progress")

	// Verify in-progress status
	showOutput, err := s.run("task", "show", taskID)
	s.requireSuccess(showOutput, err, "failed to show task after status update")
	s.Contains(showOutput, "in-progress", "task should show in-progress status")

	// Transition to done
	updateOutput, err = s.run("task", "update", taskID, "--status", "done")
	s.requireSuccess(updateOutput, err, "failed to transition task to done")

	// Verify done status
	showOutput, err = s.run("task", "show", taskID)
	s.requireSuccess(showOutput, err, "failed to show task after final update")
	s.Contains(showOutput, "done", "task should show done status")
}

// TestTaskMove tests moving task between tracks
func (s *TaskTestSuite) TestTaskMove() {
	// Create first track
	track1Output, err := s.run("track", "create", "--title", "Track 1", "--rank", "100")
	s.requireSuccess(track1Output, err, "failed to create track 1")
	track1ID := s.parseID(track1Output, "track")

	// Create second track
	track2Output, err := s.run("track", "create", "--title", "Track 2", "--rank", "200")
	s.requireSuccess(track2Output, err, "failed to create track 2")
	track2ID := s.parseID(track2Output, "track")

	// Create task in track 1
	taskOutput, err := s.run("task", "create", "--track", track1ID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	taskID := s.parseID(taskOutput, "task")

	// Move task to track 2
	moveOutput, err := s.run("task", "move", taskID, "--track", track2ID)
	s.requireSuccess(moveOutput, err, "failed to move task to track 2")

	// Verify task was moved
	showOutput, err := s.run("task", "show", taskID)
	s.requireSuccess(showOutput, err, "failed to show task after move")
	s.Contains(showOutput, track2ID, "task should show new track in details")
}

// TestTaskBacklog tests listing backlog (todo status) tasks
func (s *TaskTestSuite) TestTaskBacklog() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create task (defaults to todo status)
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	_ = s.parseID(taskOutput, "task")

	// Get backlog (should list todo status tasks)
	backlogOutput, err := s.run("task", "backlog")
	s.requireSuccess(backlogOutput, err, "failed to get task backlog")
	// Backlog output should contain reference to tasks or confirmation
	s.NotEmpty(backlogOutput, "backlog output should not be empty")
}

// TestTaskDelete tests deleting a task
func (s *TaskTestSuite) TestTaskDelete() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create task
	taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Test Task", "--rank", "100")
	s.requireSuccess(taskOutput, err, "failed to create task")
	taskID := s.parseID(taskOutput, "task")

	// Delete task
	deleteOutput, err := s.run("task", "delete", taskID, "--force")
	s.requireSuccess(deleteOutput, err, "failed to delete task")

	// Verify task is deleted by attempting to show it
	showOutput, err := s.run("task", "show", taskID)
	s.requireError(err, "task should not be found after deletion")
	s.NotEmpty(showOutput, "error message should be provided")
}

// TestTaskMultiple tests creating and managing multiple tasks
func (s *TaskTestSuite) TestTaskMultiple() {
	// Create track
	trackOutput, err := s.run("track", "create", "--title", "Test Track", "--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "track")

	// Create multiple tasks
	taskIDs := []string{}
	for i := 1; i <= 3; i++ {
		taskOutput, err := s.run("task", "create", "--track", trackID, "--title", "Task "+string(rune('0'+byte(i))), "--rank", "100")
		s.requireSuccess(taskOutput, err, "failed to create task")
		taskID := s.parseID(taskOutput, "task")
		s.NotEmpty(taskID, "task ID should be extracted")
		taskIDs = append(taskIDs, taskID)
	}

	// List all tasks
	listOutput, err := s.run("task", "list")
	s.requireSuccess(listOutput, err, "failed to list tasks")

	// Verify all tasks appear in list
	for _, taskID := range taskIDs {
		s.Contains(listOutput, taskID, "task "+taskID+" should appear in list")
	}
}
