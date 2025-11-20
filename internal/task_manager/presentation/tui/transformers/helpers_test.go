package transformers_test

import (
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/transformers"
)

// Test helper functions (used by other test files)

// mustCreateIteration is a test helper that creates an IterationEntity, panicking on error
func mustCreateIteration(number int, name, goal, deliverable string, taskIDs []string, status string, rank float64, createdAt, updatedAt time.Time) *entities.IterationEntity {
	iter, err := entities.NewIterationEntity(number, name, goal, deliverable, taskIDs, status, rank, time.Time{}, time.Time{}, createdAt, updatedAt)
	if err != nil {
		panic(err)
	}
	return iter
}

// mustCreateTrack is a test helper that creates a TrackEntity, panicking on error
func mustCreateTrack(id, roadmapID, title, description, status string, rank int, dependencies []string, createdAt, updatedAt time.Time) *entities.TrackEntity {
	track, err := entities.NewTrackEntity(id, roadmapID, title, description, status, rank, dependencies, createdAt, updatedAt)
	if err != nil {
		panic(err)
	}
	return track
}

// mustCreateTask is a test helper that creates a TaskEntity, panicking on error
func mustCreateTask(id, trackID, title, description, status string, rank int, branch string, createdAt, updatedAt time.Time) *entities.TaskEntity {
	task, err := entities.NewTaskEntity(id, trackID, title, description, status, rank, branch, createdAt, updatedAt)
	if err != nil {
		panic(err)
	}
	return task
}

func TestGetIterationIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Planned iteration", string(entities.IterationStatusPlanned), "📋"},
		{"Current iteration", string(entities.IterationStatusCurrent), "▶"},
		{"Complete iteration", string(entities.IterationStatusComplete), "✓"},
		{"Unknown status defaults to planned", "unknown", "📋"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetIterationIcon(tt.status)
			if result != tt.expected {
				t.Errorf("GetIterationIcon(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetIterationColor(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Planned iteration", string(entities.IterationStatusPlanned), "info"},
		{"Current iteration", string(entities.IterationStatusCurrent), "current"},
		{"Complete iteration", string(entities.IterationStatusComplete), "success"},
		{"Unknown status defaults to info", "unknown", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetIterationColor(tt.status)
			if result != tt.expected {
				t.Errorf("GetIterationColor(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetIterationStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Planned iteration", string(entities.IterationStatusPlanned), "Planned"},
		{"Current iteration", string(entities.IterationStatusCurrent), "Current"},
		{"Complete iteration", string(entities.IterationStatusComplete), "Complete"},
		{"Unknown status returns as-is", "unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetIterationStatusLabel(tt.status)
			if result != tt.expected {
				t.Errorf("GetIterationStatusLabel(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetTaskIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Todo task", string(entities.TaskStatusTodo), "○"},
		{"In-progress task", string(entities.TaskStatusInProgress), "◐"},
		{"Review task", string(entities.TaskStatusReview), "◑"},
		{"Done task", string(entities.TaskStatusDone), "●"},
		{"Cancelled task", string(entities.TaskStatusCancelled), "⊗"},
		{"Unknown status defaults to todo", "unknown", "○"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetTaskIcon(tt.status)
			if result != tt.expected {
				t.Errorf("GetTaskIcon(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetTaskColor(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Todo task", string(entities.TaskStatusTodo), "info"},
		{"In-progress task", string(entities.TaskStatusInProgress), "warning"},
		{"Review task", string(entities.TaskStatusReview), "warning"},
		{"Done task", string(entities.TaskStatusDone), "success"},
		{"Cancelled task", string(entities.TaskStatusCancelled), "muted"},
		{"Unknown status defaults to info", "unknown", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetTaskColor(tt.status)
			if result != tt.expected {
				t.Errorf("GetTaskColor(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetTaskStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Todo task", string(entities.TaskStatusTodo), "Todo"},
		{"In-progress task", string(entities.TaskStatusInProgress), "In Progress"},
		{"Review task", string(entities.TaskStatusReview), "Review"},
		{"Done task", string(entities.TaskStatusDone), "Done"},
		{"Cancelled task", string(entities.TaskStatusCancelled), "Cancelled"},
		{"Unknown status returns as-is", "unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetTaskStatusLabel(tt.status)
			if result != tt.expected {
				t.Errorf("GetTaskStatusLabel(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetTrackIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Not started track", string(entities.TrackStatusNotStarted), "○"},
		{"In-progress track", string(entities.TrackStatusInProgress), "◐"},
		{"Complete track", string(entities.TrackStatusComplete), "●"},
		{"Blocked track", string(entities.TrackStatusBlocked), "⊠"},
		{"Waiting track", string(entities.TrackStatusWaiting), "⏸"},
		{"Unknown status defaults to not started", "unknown", "○"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetTrackIcon(tt.status)
			if result != tt.expected {
				t.Errorf("GetTrackIcon(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetTrackColor(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Not started track", string(entities.TrackStatusNotStarted), "muted"},
		{"In-progress track", string(entities.TrackStatusInProgress), "warning"},
		{"Complete track", string(entities.TrackStatusComplete), "success"},
		{"Blocked track", string(entities.TrackStatusBlocked), "failed"},
		{"Waiting track", string(entities.TrackStatusWaiting), "warning"},
		{"Unknown status defaults to muted", "unknown", "muted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetTrackColor(tt.status)
			if result != tt.expected {
				t.Errorf("GetTrackColor(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetTrackStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Not started track", string(entities.TrackStatusNotStarted), "Not Started"},
		{"In-progress track", string(entities.TrackStatusInProgress), "In Progress"},
		{"Complete track", string(entities.TrackStatusComplete), "Complete"},
		{"Blocked track", string(entities.TrackStatusBlocked), "Blocked"},
		{"Waiting track", string(entities.TrackStatusWaiting), "Waiting"},
		{"Unknown status returns as-is", "unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetTrackStatusLabel(tt.status)
			if result != tt.expected {
				t.Errorf("GetTrackStatusLabel(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetACColor(t *testing.T) {
	tests := []struct {
		name     string
		status   entities.AcceptanceCriteriaStatus
		expected string
	}{
		{"Not started AC", entities.ACStatusNotStarted, "muted"},
		{"Verified AC", entities.ACStatusVerified, "success"},
		{"Auto verified AC", entities.ACStatusAutomaticallyVerified, "success"},
		{"Pending review AC", entities.ACStatusPendingHumanReview, "warning"},
		{"Failed AC", entities.ACStatusFailed, "failed"},
		{"Skipped AC", entities.ACStatusSkipped, "skipped"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetACColor(tt.status)
			if result != tt.expected {
				t.Errorf("GetACColor(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetACStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		status   entities.AcceptanceCriteriaStatus
		expected string
	}{
		{"Not started AC", entities.ACStatusNotStarted, "Not Started"},
		{"Verified AC", entities.ACStatusVerified, "Verified"},
		{"Auto verified AC", entities.ACStatusAutomaticallyVerified, "Auto Verified"},
		{"Pending review AC", entities.ACStatusPendingHumanReview, "Pending Review"},
		{"Failed AC", entities.ACStatusFailed, "Failed"},
		{"Skipped AC", entities.ACStatusSkipped, "Skipped"},
		{"Unknown status returns as-is", entities.AcceptanceCriteriaStatus("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetACStatusLabel(tt.status)
			if result != tt.expected {
				t.Errorf("GetACStatusLabel(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}
