package transformers_test

import (
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/transformers"
)

func TestTransformToTaskDetailViewModel_Basic(t *testing.T) {
	now := time.Now()

	task := mustCreateTask("TM-task-1", "TM-track-1", "Test Task", "Task description", "todo", 100, "feature/test", now, now)
	track := mustCreateTrack("TM-track-1", "roadmap-1", "Test Track", "Track description", "in-progress", 100, []string{}, now, now)
	acs := []*entities.AcceptanceCriteriaEntity{}
	iterations := []*entities.IterationEntity{}

	vm := transformers.TransformToTaskDetailViewModel(task, acs, track, iterations)

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

	if vm.TrackInfo == nil {
		t.Fatal("expected non-nil TrackInfo")
	}

	if vm.TrackInfo.ID != "TM-track-1" {
		t.Errorf("expected TrackInfo ID 'TM-track-1', got %q", vm.TrackInfo.ID)
	}

	if vm.TrackInfo.Title != "Test Track" {
		t.Errorf("expected TrackInfo Title 'Test Track', got %q", vm.TrackInfo.Title)
	}

	if len(vm.Iterations) != 0 {
		t.Errorf("expected 0 iterations, got %d", len(vm.Iterations))
	}

	if len(vm.AcceptanceCriteria) != 0 {
		t.Errorf("expected 0 ACs, got %d", len(vm.AcceptanceCriteria))
	}

	// Verify timestamps are formatted
	if vm.CreatedAt == "" {
		t.Error("expected non-empty CreatedAt")
	}

	if vm.UpdatedAt == "" {
		t.Error("expected non-empty UpdatedAt")
	}
}

func TestTransformToTaskDetailViewModel_WithIterations(t *testing.T) {
	now := time.Now()

	task := mustCreateTask("TM-task-1", "TM-track-1", "Test Task", "Description", "in-progress", 100, "", now, now)
	track := mustCreateTrack("TM-track-1", "roadmap-1", "Test Track", "Description", "in-progress", 100, []string{}, now, now)
	acs := []*entities.AcceptanceCriteriaEntity{}

	iterations := []*entities.IterationEntity{
		mustCreateIteration(1, "Sprint 1", "Goal 1", "Deliverable 1", []string{}, "current", 100, now, now),
		mustCreateIteration(2, "Sprint 2", "Goal 2", "Deliverable 2", []string{}, "planned", 200, now, now),
		mustCreateIteration(3, "Sprint 3", "Goal 3", "Deliverable 3", []string{}, "complete", 300, now, now),
	}

	vm := transformers.TransformToTaskDetailViewModel(task, acs, track, iterations)

	if len(vm.Iterations) != 3 {
		t.Errorf("expected 3 iterations, got %d", len(vm.Iterations))
	}

	// Verify first iteration
	if vm.Iterations[0].Number != 1 {
		t.Errorf("expected first iteration Number 1, got %d", vm.Iterations[0].Number)
	}

	if vm.Iterations[0].Name != "Sprint 1" {
		t.Errorf("expected first iteration Name 'Sprint 1', got %q", vm.Iterations[0].Name)
	}

	if vm.Iterations[0].Status != "current" {
		t.Errorf("expected first iteration Status 'current', got %q", vm.Iterations[0].Status)
	}

	// Verify complete iteration is included (no filtering at this level)
	if vm.Iterations[2].Status != "complete" {
		t.Errorf("expected third iteration Status 'complete', got %q", vm.Iterations[2].Status)
	}
}

func TestTransformToTaskDetailViewModel_WithACs(t *testing.T) {
	now := time.Now()

	task := mustCreateTask("TM-task-1", "TM-track-1", "Test Task", "Description", "todo", 100, "", now, now)
	track := mustCreateTrack("TM-track-1", "roadmap-1", "Test Track", "Description", "in-progress", 100, []string{}, now, now)

	acs := []*entities.AcceptanceCriteriaEntity{
		entities.NewAcceptanceCriteriaEntity("TM-ac-1", "TM-task-1", "AC 1 description", entities.VerificationTypeManual, "Test instructions 1", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-2", "TM-task-1", "AC 2 description", entities.VerificationTypeManual, "Test instructions 2", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-3", "TM-task-1", "AC 3 description", entities.VerificationTypeManual, "Test instructions 3", now, now),
	}

	// Set different statuses
	acs[0].Status = entities.ACStatusVerified
	acs[1].Status = entities.ACStatusFailed
	acs[1].Notes = "Failed due to X"
	acs[2].Status = entities.ACStatusSkipped

	iterations := []*entities.IterationEntity{}

	vm := transformers.TransformToTaskDetailViewModel(task, acs, track, iterations)

	if len(vm.AcceptanceCriteria) != 3 {
		t.Errorf("expected 3 ACs, got %d", len(vm.AcceptanceCriteria))
	}

	// Verify first AC
	if vm.AcceptanceCriteria[0].ID != "TM-ac-1" {
		t.Errorf("expected first AC ID 'TM-ac-1', got %q", vm.AcceptanceCriteria[0].ID)
	}

	if vm.AcceptanceCriteria[0].Description != "AC 1 description" {
		t.Errorf("expected first AC Description 'AC 1 description', got %q", vm.AcceptanceCriteria[0].Description)
	}

	if vm.AcceptanceCriteria[0].Status != "verified" {
		t.Errorf("expected first AC Status 'verified', got %q", vm.AcceptanceCriteria[0].Status)
	}

	if vm.AcceptanceCriteria[0].StatusIcon != "✓" {
		t.Errorf("expected first AC StatusIcon '✓', got %q", vm.AcceptanceCriteria[0].StatusIcon)
	}

	if vm.AcceptanceCriteria[0].TestingInstructions != "Test instructions 1" {
		t.Errorf("expected first AC TestingInstructions 'Test instructions 1', got %q", vm.AcceptanceCriteria[0].TestingInstructions)
	}

	// Verify all ACs are initially collapsed
	for i, ac := range vm.AcceptanceCriteria {
		if ac.IsExpanded {
			t.Errorf("expected AC %d to be collapsed, got expanded", i+1)
		}
	}

	// Verify second AC has notes
	if vm.AcceptanceCriteria[1].Notes != "Failed due to X" {
		t.Errorf("expected second AC Notes 'Failed due to X', got %q", vm.AcceptanceCriteria[1].Notes)
	}

	// Verify status icons
	expectedIcons := []string{"✓", "✗", "⊘"}
	for i, expected := range expectedIcons {
		if vm.AcceptanceCriteria[i].StatusIcon != expected {
			t.Errorf("expected AC %d icon %q, got %q", i+1, expected, vm.AcceptanceCriteria[i].StatusIcon)
		}
	}
}

func TestTransformToTaskDetailViewModel_NilTrack(t *testing.T) {
	now := time.Now()

	task := mustCreateTask("TM-task-1", "TM-track-1", "Test Task", "Description", "todo", 100, "", now, now)
	acs := []*entities.AcceptanceCriteriaEntity{}
	iterations := []*entities.IterationEntity{}

	vm := transformers.TransformToTaskDetailViewModel(task, acs, nil, iterations)

	if vm.TrackInfo != nil {
		t.Error("expected nil TrackInfo when track is nil")
	}
}

func TestTransformToTaskDetailViewModel_EmptyBranch(t *testing.T) {
	now := time.Now()

	task := mustCreateTask("TM-task-1", "TM-track-1", "Test Task", "Description", "done", 100, "", now, now)
	track := mustCreateTrack("TM-track-1", "roadmap-1", "Test Track", "Description", "complete", 100, []string{}, now, now)
	acs := []*entities.AcceptanceCriteriaEntity{}
	iterations := []*entities.IterationEntity{}

	vm := transformers.TransformToTaskDetailViewModel(task, acs, track, iterations)

	if vm.Branch != "" {
		t.Errorf("expected empty Branch, got %q", vm.Branch)
	}
}

func TestTransformToTaskDetailViewModel_TimestampFormat(t *testing.T) {
	createdAt := time.Date(2025, 11, 14, 10, 30, 45, 0, time.UTC)
	updatedAt := time.Date(2025, 11, 14, 14, 15, 30, 0, time.UTC)

	task := mustCreateTask("TM-task-1", "TM-track-1", "Test Task", "Description", "todo", 100, "", createdAt, updatedAt)
	track := mustCreateTrack("TM-track-1", "roadmap-1", "Test Track", "Description", "in-progress", 100, []string{}, createdAt, updatedAt)
	acs := []*entities.AcceptanceCriteriaEntity{}
	iterations := []*entities.IterationEntity{}

	vm := transformers.TransformToTaskDetailViewModel(task, acs, track, iterations)

	// Check timestamp format (YYYY-MM-DD HH:MM:SS)
	expectedCreatedAt := "2025-11-14 10:30:45"
	expectedUpdatedAt := "2025-11-14 14:15:30"

	if vm.CreatedAt != expectedCreatedAt {
		t.Errorf("expected CreatedAt %q, got %q", expectedCreatedAt, vm.CreatedAt)
	}

	if vm.UpdatedAt != expectedUpdatedAt {
		t.Errorf("expected UpdatedAt %q, got %q", expectedUpdatedAt, vm.UpdatedAt)
	}
}
