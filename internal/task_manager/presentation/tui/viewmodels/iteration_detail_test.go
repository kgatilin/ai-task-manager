package viewmodels_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

func TestNewProgressViewModel(t *testing.T) {
	tests := []struct {
		name      string
		completed int
		total     int
		expected  float64
	}{
		{"zero total", 0, 0, 0.0},
		{"zero completed", 0, 10, 0.0},
		{"half completed", 5, 10, 0.5},
		{"all completed", 10, 10, 1.0},
		{"partial", 3, 7, 3.0 / 7.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := viewmodels.NewProgressViewModel(tt.completed, tt.total)

			if vm.Completed != tt.completed {
				t.Errorf("expected Completed %d, got %d", tt.completed, vm.Completed)
			}

			if vm.Total != tt.total {
				t.Errorf("expected Total %d, got %d", tt.total, vm.Total)
			}

			if vm.Percent != tt.expected {
				t.Errorf("expected Percent %f, got %f", tt.expected, vm.Percent)
			}
		})
	}
}

func TestTaskRowViewModel_Fields(t *testing.T) {
	vm := &viewmodels.TaskRowViewModel{
		ID:          "TM-task-1",
		Title:       "Task title",
		Status:      "todo",
		Description: "Task description",
	}

	if vm.ID != "TM-task-1" {
		t.Errorf("expected ID 'TM-task-1', got %q", vm.ID)
	}

	if vm.Title != "Task title" {
		t.Errorf("expected Title 'Task title', got %q", vm.Title)
	}

	if vm.Status != "todo" {
		t.Errorf("expected Status 'todo', got %q", vm.Status)
	}

	if vm.Description != "Task description" {
		t.Errorf("expected Description 'Task description', got %q", vm.Description)
	}
}

func TestIterationACViewModel_Fields(t *testing.T) {
	vm := &viewmodels.IterationACViewModel{
		ID:                  "TM-ac-1",
		Description:         "AC description",
		Status:              "verified",
		StatusIcon:          "✓",
		TestingInstructions: "Test instructions",
		Notes:               "AC notes",
	}

	if vm.ID != "TM-ac-1" {
		t.Errorf("expected ID 'TM-ac-1', got %q", vm.ID)
	}

	if vm.Description != "AC description" {
		t.Errorf("expected Description 'AC description', got %q", vm.Description)
	}

	if vm.Status != "verified" {
		t.Errorf("expected Status 'verified', got %q", vm.Status)
	}

	if vm.StatusIcon != "✓" {
		t.Errorf("expected StatusIcon '✓', got %q", vm.StatusIcon)
	}

	if vm.TestingInstructions != "Test instructions" {
		t.Errorf("expected TestingInstructions 'Test instructions', got %q", vm.TestingInstructions)
	}

	if vm.Notes != "AC notes" {
		t.Errorf("expected Notes 'AC notes', got %q", vm.Notes)
	}
}

func TestNewIterationDetailViewModel(t *testing.T) {
	vm := viewmodels.NewIterationDetailViewModel(1, "Test Iteration", "Test goal", "Test deliverable", "current")

	if vm == nil {
		t.Fatal("expected non-nil view model")
	}

	if vm.Number != 1 {
		t.Errorf("expected Number 1, got %d", vm.Number)
	}

	if vm.Name != "Test Iteration" {
		t.Errorf("expected Name 'Test Iteration', got %q", vm.Name)
	}

	if vm.Goal != "Test goal" {
		t.Errorf("expected Goal 'Test goal', got %q", vm.Goal)
	}

	if vm.Deliverable != "Test deliverable" {
		t.Errorf("expected Deliverable 'Test deliverable', got %q", vm.Deliverable)
	}

	if vm.Status != "current" {
		t.Errorf("expected Status 'current', got %q", vm.Status)
	}

	if vm.TODOTasks == nil {
		t.Error("expected non-nil TODOTasks slice")
	}

	if vm.InProgressTasks == nil {
		t.Error("expected non-nil InProgressTasks slice")
	}

	if vm.DoneTasks == nil {
		t.Error("expected non-nil DoneTasks slice")
	}

	if vm.AcceptanceCriteria == nil {
		t.Error("expected non-nil AcceptanceCriteria slice")
	}

	if vm.Progress == nil {
		t.Error("expected non-nil Progress")
	}

	if vm.Progress.Completed != 0 || vm.Progress.Total != 0 {
		t.Errorf("expected zero progress, got %d/%d", vm.Progress.Completed, vm.Progress.Total)
	}
}

func TestIterationDetailViewModel_PopulatedData(t *testing.T) {
	vm := viewmodels.NewIterationDetailViewModel(2, "Sprint 2", "Complete features", "Feature set", "planned")

	// Add tasks
	vm.TODOTasks = append(vm.TODOTasks, &viewmodels.TaskRowViewModel{
		ID:     "TM-task-1",
		Title:  "Task 1",
		Status: "todo",
	})

	vm.InProgressTasks = append(vm.InProgressTasks, &viewmodels.TaskRowViewModel{
		ID:     "TM-task-2",
		Title:  "Task 2",
		Status: "in-progress",
	})

	vm.DoneTasks = append(vm.DoneTasks, &viewmodels.TaskRowViewModel{
		ID:     "TM-task-3",
		Title:  "Task 3",
		Status: "done",
	})

	// Add ACs
	vm.AcceptanceCriteria = append(vm.AcceptanceCriteria, &viewmodels.IterationACViewModel{
		ID:         "TM-ac-1",
		Status:     "verified",
		StatusIcon: "✓",
	})

	// Update progress
	vm.Progress = viewmodels.NewProgressViewModel(1, 3)

	if len(vm.TODOTasks) != 1 {
		t.Errorf("expected 1 TODO task, got %d", len(vm.TODOTasks))
	}

	if len(vm.InProgressTasks) != 1 {
		t.Errorf("expected 1 in-progress task, got %d", len(vm.InProgressTasks))
	}

	if len(vm.DoneTasks) != 1 {
		t.Errorf("expected 1 done task, got %d", len(vm.DoneTasks))
	}

	if len(vm.AcceptanceCriteria) != 1 {
		t.Errorf("expected 1 AC, got %d", len(vm.AcceptanceCriteria))
	}

	if vm.Progress.Completed != 1 || vm.Progress.Total != 3 {
		t.Errorf("expected progress 1/3, got %d/%d", vm.Progress.Completed, vm.Progress.Total)
	}
}

func TestIterationDetailViewModel_EmptyCollections(t *testing.T) {
	vm := viewmodels.NewIterationDetailViewModel(1, "Test", "Goal", "Deliverable", "current")

	if len(vm.TODOTasks) != 0 {
		t.Errorf("expected 0 TODO tasks, got %d", len(vm.TODOTasks))
	}

	if len(vm.InProgressTasks) != 0 {
		t.Errorf("expected 0 in-progress tasks, got %d", len(vm.InProgressTasks))
	}

	if len(vm.DoneTasks) != 0 {
		t.Errorf("expected 0 done tasks, got %d", len(vm.DoneTasks))
	}

	if len(vm.AcceptanceCriteria) != 0 {
		t.Errorf("expected 0 ACs, got %d", len(vm.AcceptanceCriteria))
	}
}

func TestIterationDetailViewModel_TimestampFields(t *testing.T) {
	vm := viewmodels.NewIterationDetailViewModel(1, "Test", "Goal", "Deliverable", "current")

	vm.StartedAt = "2025-11-14 10:00:00"
	vm.CompletedAt = "2025-11-14 12:00:00"

	if vm.StartedAt != "2025-11-14 10:00:00" {
		t.Errorf("expected StartedAt '2025-11-14 10:00:00', got %q", vm.StartedAt)
	}

	if vm.CompletedAt != "2025-11-14 12:00:00" {
		t.Errorf("expected CompletedAt '2025-11-14 12:00:00', got %q", vm.CompletedAt)
	}
}
