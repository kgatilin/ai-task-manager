package viewmodels_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

func TestNewTaskDetailViewModel(t *testing.T) {
	vm := viewmodels.NewTaskDetailViewModel("TM-task-1", "Test Task", "Task description", "todo", "feature/test")

	if vm == nil {
		t.Fatal("expected non-nil view model")
	}

	if vm.ID != "TM-task-1" {
		t.Errorf("expected ID 'TM-task-1', got %q", vm.ID)
	}

	if vm.Title != "Test Task" {
		t.Errorf("expected Title 'Test Task', got %q", vm.Title)
	}

	if vm.Description != "Task description" {
		t.Errorf("expected Description 'Task description', got %q", vm.Description)
	}

	if vm.Status != "todo" {
		t.Errorf("expected Status 'todo', got %q", vm.Status)
	}

	if vm.Branch != "feature/test" {
		t.Errorf("expected Branch 'feature/test', got %q", vm.Branch)
	}

	if vm.Iterations == nil {
		t.Error("expected non-nil Iterations slice")
	}

	if len(vm.Iterations) != 0 {
		t.Errorf("expected empty Iterations, got %d items", len(vm.Iterations))
	}

	if vm.AcceptanceCriteria == nil {
		t.Error("expected non-nil AcceptanceCriteria slice")
	}

	if len(vm.AcceptanceCriteria) != 0 {
		t.Errorf("expected empty AcceptanceCriteria, got %d items", len(vm.AcceptanceCriteria))
	}
}

func TestTaskDetailViewModel_PopulatedData(t *testing.T) {
	vm := viewmodels.NewTaskDetailViewModel("TM-task-1", "Test Task", "Description", "in-progress", "main")

	// Add track info
	vm.TrackInfo = &viewmodels.TrackInfoViewModel{
		ID:          "TM-track-1",
		Title:       "Test Track",
		Description: "Track description",
		Status:      "in-progress",
	}

	// Add iterations
	vm.Iterations = append(vm.Iterations, &viewmodels.IterationMembershipViewModel{
		Number: 1,
		Name:   "Sprint 1",
		Status: "current",
	})

	vm.Iterations = append(vm.Iterations, &viewmodels.IterationMembershipViewModel{
		Number: 2,
		Name:   "Sprint 2",
		Status: "planned",
	})

	// Add ACs
	vm.AcceptanceCriteria = append(vm.AcceptanceCriteria, &viewmodels.ACDetailViewModel{
		ID:                  "TM-ac-1",
		Description:         "AC description",
		Status:              "verified",
		StatusIcon:          "✓",
		TestingInstructions: "Test instructions",
		Notes:               "AC notes",
		IsExpanded:          false,
	})

	vm.AcceptanceCriteria = append(vm.AcceptanceCriteria, &viewmodels.ACDetailViewModel{
		ID:                  "TM-ac-2",
		Description:         "Another AC",
		Status:              "not_started",
		StatusIcon:          "○",
		TestingInstructions: "More instructions",
		Notes:               "",
		IsExpanded:          true,
	})

	if vm.TrackInfo == nil {
		t.Fatal("expected non-nil TrackInfo")
	}

	if vm.TrackInfo.ID != "TM-track-1" {
		t.Errorf("expected TrackInfo ID 'TM-track-1', got %q", vm.TrackInfo.ID)
	}

	if len(vm.Iterations) != 2 {
		t.Errorf("expected 2 iterations, got %d", len(vm.Iterations))
	}

	if vm.Iterations[0].Number != 1 {
		t.Errorf("expected first iteration number 1, got %d", vm.Iterations[0].Number)
	}

	if len(vm.AcceptanceCriteria) != 2 {
		t.Errorf("expected 2 ACs, got %d", len(vm.AcceptanceCriteria))
	}

	if vm.AcceptanceCriteria[0].IsExpanded {
		t.Error("expected first AC to be collapsed")
	}

	if !vm.AcceptanceCriteria[1].IsExpanded {
		t.Error("expected second AC to be expanded")
	}
}

func TestACDetailViewModel_Fields(t *testing.T) {
	vm := &viewmodels.ACDetailViewModel{
		ID:                  "TM-ac-1",
		Description:         "AC description",
		Status:              "failed",
		StatusIcon:          "✗",
		TestingInstructions: "Testing instructions",
		Notes:               "Failure notes",
		IsExpanded:          true,
	}

	if vm.ID != "TM-ac-1" {
		t.Errorf("expected ID 'TM-ac-1', got %q", vm.ID)
	}

	if vm.Description != "AC description" {
		t.Errorf("expected Description 'AC description', got %q", vm.Description)
	}

	if vm.Status != "failed" {
		t.Errorf("expected Status 'failed', got %q", vm.Status)
	}

	if vm.StatusIcon != "✗" {
		t.Errorf("expected StatusIcon '✗', got %q", vm.StatusIcon)
	}

	if vm.TestingInstructions != "Testing instructions" {
		t.Errorf("expected TestingInstructions 'Testing instructions', got %q", vm.TestingInstructions)
	}

	if vm.Notes != "Failure notes" {
		t.Errorf("expected Notes 'Failure notes', got %q", vm.Notes)
	}

	if !vm.IsExpanded {
		t.Error("expected IsExpanded to be true")
	}
}

func TestTrackInfoViewModel_Fields(t *testing.T) {
	vm := &viewmodels.TrackInfoViewModel{
		ID:          "TM-track-1",
		Title:       "Track Title",
		Description: "Track Description",
		Status:      "in-progress",
	}

	if vm.ID != "TM-track-1" {
		t.Errorf("expected ID 'TM-track-1', got %q", vm.ID)
	}

	if vm.Title != "Track Title" {
		t.Errorf("expected Title 'Track Title', got %q", vm.Title)
	}

	if vm.Description != "Track Description" {
		t.Errorf("expected Description 'Track Description', got %q", vm.Description)
	}

	if vm.Status != "in-progress" {
		t.Errorf("expected Status 'in-progress', got %q", vm.Status)
	}
}

func TestIterationMembershipViewModel_Fields(t *testing.T) {
	vm := &viewmodels.IterationMembershipViewModel{
		Number: 5,
		Name:   "Iteration Name",
		Status: "complete",
	}

	if vm.Number != 5 {
		t.Errorf("expected Number 5, got %d", vm.Number)
	}

	if vm.Name != "Iteration Name" {
		t.Errorf("expected Name 'Iteration Name', got %q", vm.Name)
	}

	if vm.Status != "complete" {
		t.Errorf("expected Status 'complete', got %q", vm.Status)
	}
}

func TestTaskDetailViewModel_EmptyBranch(t *testing.T) {
	vm := viewmodels.NewTaskDetailViewModel("TM-task-1", "Test Task", "Description", "todo", "")

	if vm.Branch != "" {
		t.Errorf("expected empty Branch, got %q", vm.Branch)
	}
}

func TestTaskDetailViewModel_TimestampFields(t *testing.T) {
	vm := viewmodels.NewTaskDetailViewModel("TM-task-1", "Test Task", "Description", "todo", "main")

	vm.CreatedAt = "2025-11-14 10:00:00"
	vm.UpdatedAt = "2025-11-14 12:00:00"

	if vm.CreatedAt != "2025-11-14 10:00:00" {
		t.Errorf("expected CreatedAt '2025-11-14 10:00:00', got %q", vm.CreatedAt)
	}

	if vm.UpdatedAt != "2025-11-14 12:00:00" {
		t.Errorf("expected UpdatedAt '2025-11-14 12:00:00', got %q", vm.UpdatedAt)
	}
}
