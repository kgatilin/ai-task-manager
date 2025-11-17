package viewmodels_test

import (
	"testing"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

func TestNewRoadmapListViewModel(t *testing.T) {
	vm := viewmodels.NewRoadmapListViewModel()

	if vm == nil {
		t.Fatal("expected non-nil view model")
	}

	if vm.Vision != "" {
		t.Errorf("expected empty Vision, got %q", vm.Vision)
	}

	if vm.SuccessCriteria != "" {
		t.Errorf("expected empty SuccessCriteria, got %q", vm.SuccessCriteria)
	}

	if vm.ActiveIterations == nil {
		t.Error("expected non-nil ActiveIterations slice")
	}

	if len(vm.ActiveIterations) != 0 {
		t.Errorf("expected empty ActiveIterations, got %d items", len(vm.ActiveIterations))
	}

	if vm.ActiveTracks == nil {
		t.Error("expected non-nil ActiveTracks slice")
	}

	if len(vm.ActiveTracks) != 0 {
		t.Errorf("expected empty ActiveTracks, got %d items", len(vm.ActiveTracks))
	}

	if vm.BacklogTasks == nil {
		t.Error("expected non-nil BacklogTasks slice")
	}

	if len(vm.BacklogTasks) != 0 {
		t.Errorf("expected empty BacklogTasks, got %d items", len(vm.BacklogTasks))
	}
}

func TestRoadmapListViewModel_PopulatedData(t *testing.T) {
	vm := viewmodels.NewRoadmapListViewModel()

	// Set roadmap fields
	vm.Vision = "Build world-class platform"
	vm.SuccessCriteria = "Launch in Q1 2025"

	if vm.Vision != "Build world-class platform" {
		t.Errorf("expected Vision 'Build world-class platform', got %q", vm.Vision)
	}

	if vm.SuccessCriteria != "Launch in Q1 2025" {
		t.Errorf("expected SuccessCriteria 'Launch in Q1 2025', got %q", vm.SuccessCriteria)
	}

	// Add iterations
	vm.ActiveIterations = append(vm.ActiveIterations, &viewmodels.IterationCardViewModel{
		Number:      1,
		Name:        "Sprint 1",
		Goal:        "Build foundation",
		Status:      "current",
		TaskCount:   5,
		Deliverable: "Working prototype",
	})

	if len(vm.ActiveIterations) != 1 {
		t.Errorf("expected 1 iteration, got %d", len(vm.ActiveIterations))
	}

	if vm.ActiveIterations[0].Number != 1 {
		t.Errorf("expected iteration number 1, got %d", vm.ActiveIterations[0].Number)
	}

	if vm.ActiveIterations[0].Name != "Sprint 1" {
		t.Errorf("expected iteration name 'Sprint 1', got %q", vm.ActiveIterations[0].Name)
	}

	// Add tracks
	vm.ActiveTracks = append(vm.ActiveTracks, &viewmodels.TrackCardViewModel{
		ID:          "TM-track-1",
		Title:       "Core Features",
		Description: "Implement core functionality",
		Status:      "in-progress",
		TaskCount:   3,
	})

	if len(vm.ActiveTracks) != 1 {
		t.Errorf("expected 1 track, got %d", len(vm.ActiveTracks))
	}

	if vm.ActiveTracks[0].ID != "TM-track-1" {
		t.Errorf("expected track ID 'TM-track-1', got %q", vm.ActiveTracks[0].ID)
	}

	// Add backlog tasks
	vm.BacklogTasks = append(vm.BacklogTasks, &viewmodels.BacklogTaskViewModel{
		ID:          "TM-task-1",
		Title:       "Implement feature X",
		Status:      "todo",
		TrackID:     "TM-track-1",
		Description: "Feature X description",
	})

	if len(vm.BacklogTasks) != 1 {
		t.Errorf("expected 1 backlog task, got %d", len(vm.BacklogTasks))
	}

	if vm.BacklogTasks[0].ID != "TM-task-1" {
		t.Errorf("expected task ID 'TM-task-1', got %q", vm.BacklogTasks[0].ID)
	}
}

func TestIterationCardViewModel_Fields(t *testing.T) {
	vm := &viewmodels.IterationCardViewModel{
		Number:      42,
		Name:        "Test Iteration",
		Goal:        "Test goal",
		Status:      "planned",
		TaskCount:   10,
		Deliverable: "Test deliverable",
	}

	if vm.Number != 42 {
		t.Errorf("expected Number 42, got %d", vm.Number)
	}

	if vm.Name != "Test Iteration" {
		t.Errorf("expected Name 'Test Iteration', got %q", vm.Name)
	}

	if vm.Goal != "Test goal" {
		t.Errorf("expected Goal 'Test goal', got %q", vm.Goal)
	}

	if vm.Status != "planned" {
		t.Errorf("expected Status 'planned', got %q", vm.Status)
	}

	if vm.TaskCount != 10 {
		t.Errorf("expected TaskCount 10, got %d", vm.TaskCount)
	}

	if vm.Deliverable != "Test deliverable" {
		t.Errorf("expected Deliverable 'Test deliverable', got %q", vm.Deliverable)
	}
}

func TestTrackCardViewModel_Fields(t *testing.T) {
	vm := &viewmodels.TrackCardViewModel{
		ID:          "test-track-1",
		Title:       "Test Track",
		Description: "Test description",
		Status:      "in-progress",
		TaskCount:   7,
	}

	if vm.ID != "test-track-1" {
		t.Errorf("expected ID 'test-track-1', got %q", vm.ID)
	}

	if vm.Title != "Test Track" {
		t.Errorf("expected Title 'Test Track', got %q", vm.Title)
	}

	if vm.Description != "Test description" {
		t.Errorf("expected Description 'Test description', got %q", vm.Description)
	}

	if vm.Status != "in-progress" {
		t.Errorf("expected Status 'in-progress', got %q", vm.Status)
	}

	if vm.TaskCount != 7 {
		t.Errorf("expected TaskCount 7, got %d", vm.TaskCount)
	}
}

func TestBacklogTaskViewModel_Fields(t *testing.T) {
	vm := &viewmodels.BacklogTaskViewModel{
		ID:          "test-task-1",
		Title:       "Test Task",
		Status:      "todo",
		TrackID:     "test-track-1",
		Description: "Task description",
	}

	if vm.ID != "test-task-1" {
		t.Errorf("expected ID 'test-task-1', got %q", vm.ID)
	}

	if vm.Title != "Test Task" {
		t.Errorf("expected Title 'Test Task', got %q", vm.Title)
	}

	if vm.Status != "todo" {
		t.Errorf("expected Status 'todo', got %q", vm.Status)
	}

	if vm.TrackID != "test-track-1" {
		t.Errorf("expected TrackID 'test-track-1', got %q", vm.TrackID)
	}

	if vm.Description != "Task description" {
		t.Errorf("expected Description 'Task description', got %q", vm.Description)
	}
}
