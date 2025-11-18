package transformers_test

import (
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/transformers"
)

func TestTransformToIterationDetailViewModel(t *testing.T) {
	now := time.Now()
	startedAt := now.Add(-24 * time.Hour)

	// Create iteration with StartedAt
	iteration, err := entities.NewIterationEntity(1, "Sprint 1", "Complete features", "Feature set", []string{"TM-task-1", "TM-task-2", "TM-task-3"}, "current", 100, startedAt, time.Time{}, now, now)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}

	// Create tasks with different statuses
	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "in-progress", 200, "", now, now),
		mustCreateTask("TM-task-3", "TM-track-1", "Task 3", "Description 3", "done", 300, "", now, now),
		mustCreateTask("TM-task-4", "TM-track-1", "Task 4", "Description 4", "review", 400, "", now, now),
	}

	// Create ACs with different statuses
	acs := []*entities.AcceptanceCriteriaEntity{
		entities.NewAcceptanceCriteriaEntity("TM-ac-1", "TM-task-1", "AC 1", entities.VerificationTypeManual, "Test 1", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-2", "TM-task-2", "AC 2", entities.VerificationTypeManual, "Test 2", now, now),
	}
	// Set AC statuses
	acs[0].Status = entities.ACStatusVerified
	acs[1].Status = entities.ACStatusSkipped

	vm := transformers.TransformToIterationDetailViewModel(iteration, tasks, acs, []*entities.DocumentEntity{})

	// Verify iteration metadata
	if vm.Number != 1 {
		t.Errorf("expected Number 1, got %d", vm.Number)
	}

	if vm.Name != "Sprint 1" {
		t.Errorf("expected Name 'Sprint 1', got %q", vm.Name)
	}

	if vm.Goal != "Complete features" {
		t.Errorf("expected Goal 'Complete features', got %q", vm.Goal)
	}

	if vm.Deliverable != "Feature set" {
		t.Errorf("expected Deliverable 'Feature set', got %q", vm.Deliverable)
	}

	if vm.Status != "current" {
		t.Errorf("expected Status 'current', got %q", vm.Status)
	}

	// Verify timestamps are formatted
	if vm.StartedAt == "" {
		t.Error("expected non-empty StartedAt")
	}

	// Verify task grouping
	if len(vm.TODOTasks) != 1 {
		t.Errorf("expected 1 TODO task, got %d", len(vm.TODOTasks))
	}

	if len(vm.InProgressTasks) != 1 {
		t.Errorf("expected 1 in-progress task, got %d", len(vm.InProgressTasks))
	}

	if len(vm.ReviewTasks) != 1 {
		t.Errorf("expected 1 review task, got %d", len(vm.ReviewTasks))
	}

	if len(vm.DoneTasks) != 1 {
		t.Errorf("expected 1 done task, got %d", len(vm.DoneTasks))
	}

	// Verify task details
	if vm.TODOTasks[0].ID != "TM-task-1" {
		t.Errorf("expected TODO task ID 'TM-task-1', got %q", vm.TODOTasks[0].ID)
	}

	if vm.InProgressTasks[0].ID != "TM-task-2" {
		t.Errorf("expected in-progress task ID 'TM-task-2', got %q", vm.InProgressTasks[0].ID)
	}

	if vm.ReviewTasks[0].ID != "TM-task-4" {
		t.Errorf("expected review task ID 'TM-task-4', got %q", vm.ReviewTasks[0].ID)
	}

	// Verify ACs
	if len(vm.AcceptanceCriteria) != 2 {
		t.Errorf("expected 2 ACs, got %d", len(vm.AcceptanceCriteria))
	}

	if vm.AcceptanceCriteria[0].ID != "TM-ac-1" {
		t.Errorf("expected first AC ID 'TM-ac-1', got %q", vm.AcceptanceCriteria[0].ID)
	}

	if vm.AcceptanceCriteria[0].StatusIcon != "✓" {
		t.Errorf("expected first AC icon '✓', got %q", vm.AcceptanceCriteria[0].StatusIcon)
	}

	if vm.AcceptanceCriteria[1].StatusIcon != "⊘" {
		t.Errorf("expected second AC icon '⊘', got %q", vm.AcceptanceCriteria[1].StatusIcon)
	}

	// Verify progress (1 done out of 4 total)
	if vm.Progress.Completed != 1 {
		t.Errorf("expected progress completed 1, got %d", vm.Progress.Completed)
	}

	if vm.Progress.Total != 4 {
		t.Errorf("expected progress total 4, got %d", vm.Progress.Total)
	}

	if vm.Progress.Percent != 0.25 {
		t.Errorf("expected progress percent 0.25, got %f", vm.Progress.Percent)
	}
}

func TestTransformToIterationDetailViewModel_EmptyTasks(t *testing.T) {
	now := time.Now()

	iteration, err := entities.NewIterationEntity(1, "Sprint 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}

	vm := transformers.TransformToIterationDetailViewModel(iteration, []*entities.TaskEntity{}, []*entities.AcceptanceCriteriaEntity{}, []*entities.DocumentEntity{})

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

	// Progress should be 0/0
	if vm.Progress.Completed != 0 || vm.Progress.Total != 0 {
		t.Errorf("expected progress 0/0, got %d/%d", vm.Progress.Completed, vm.Progress.Total)
	}
}

func TestTransformToIterationDetailViewModel_AllTasksDone(t *testing.T) {
	now := time.Now()

	iteration, err := entities.NewIterationEntity(1, "Sprint 1", "Goal", "Deliverable", []string{"TM-task-1", "TM-task-2"}, "complete", 100, time.Time{}, time.Time{}, now, now)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}

	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "done", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "done", 200, "", now, now),
	}

	vm := transformers.TransformToIterationDetailViewModel(iteration, tasks, []*entities.AcceptanceCriteriaEntity{}, []*entities.DocumentEntity{})

	if len(vm.TODOTasks) != 0 {
		t.Errorf("expected 0 TODO tasks, got %d", len(vm.TODOTasks))
	}

	if len(vm.InProgressTasks) != 0 {
		t.Errorf("expected 0 in-progress tasks, got %d", len(vm.InProgressTasks))
	}

	if len(vm.DoneTasks) != 2 {
		t.Errorf("expected 2 done tasks, got %d", len(vm.DoneTasks))
	}

	// Progress should be 2/2 (100%)
	if vm.Progress.Completed != 2 {
		t.Errorf("expected progress completed 2, got %d", vm.Progress.Completed)
	}

	if vm.Progress.Total != 2 {
		t.Errorf("expected progress total 2, got %d", vm.Progress.Total)
	}

	if vm.Progress.Percent != 1.0 {
		t.Errorf("expected progress percent 1.0, got %f", vm.Progress.Percent)
	}
}

func TestTransformToIterationDetailViewModel_CompletedAtTimestamp(t *testing.T) {
	now := time.Now()
	completedAt := now.Add(-1 * time.Hour)

	iteration, err := entities.NewIterationEntity(1, "Sprint 1", "Goal", "Deliverable", []string{}, "complete", 100, time.Time{}, completedAt, now, now)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}

	vm := transformers.TransformToIterationDetailViewModel(iteration, []*entities.TaskEntity{}, []*entities.AcceptanceCriteriaEntity{}, []*entities.DocumentEntity{})

	if vm.CompletedAt == "" {
		t.Error("expected non-empty CompletedAt")
	}

	// Verify timestamp format
	expectedFormat := completedAt.Format("2006-01-02 15:04:05")
	if vm.CompletedAt != expectedFormat {
		t.Errorf("expected CompletedAt %q, got %q", expectedFormat, vm.CompletedAt)
	}
}

func TestTransformToIterationDetailViewModel_ACStatusIcons(t *testing.T) {
	now := time.Now()

	iteration, err := entities.NewIterationEntity(1, "Sprint 1", "Goal", "Deliverable", []string{}, "current", 100, time.Time{}, time.Time{}, now, now)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}

	// Create ACs with all status types
	acs := []*entities.AcceptanceCriteriaEntity{
		entities.NewAcceptanceCriteriaEntity("TM-ac-1", "TM-task-1", "AC 1", entities.VerificationTypeManual, "Test 1", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-2", "TM-task-1", "AC 2", entities.VerificationTypeManual, "Test 2", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-3", "TM-task-1", "AC 3", entities.VerificationTypeManual, "Test 3", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-4", "TM-task-1", "AC 4", entities.VerificationTypeManual, "Test 4", now, now),
	}

	// Set different statuses
	acs[0].Status = entities.ACStatusVerified
	acs[1].Status = entities.ACStatusSkipped
	acs[2].Status = entities.ACStatusFailed
	acs[3].Status = entities.ACStatusNotStarted

	vm := transformers.TransformToIterationDetailViewModel(iteration, []*entities.TaskEntity{}, acs, []*entities.DocumentEntity{})

	// Verify status icons
	expectedIcons := []string{"✓", "⊘", "✗", "○"}
	for i, expected := range expectedIcons {
		if vm.AcceptanceCriteria[i].StatusIcon != expected {
			t.Errorf("expected AC %d icon %q, got %q", i+1, expected, vm.AcceptanceCriteria[i].StatusIcon)
		}
	}
}

func TestTransformToIterationDetailViewModel_TaskACGrouping(t *testing.T) {
	now := time.Now()

	// Create iteration with 3 tasks
	iteration, err := entities.NewIterationEntity(1, "Sprint 1", "Complete features", "Feature set", []string{"TM-task-1", "TM-task-2", "TM-task-3"}, "current", 100, time.Time{}, time.Time{}, now, now)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}

	// Create tasks
	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "in-progress", 200, "", now, now),
		mustCreateTask("TM-task-3", "TM-track-1", "Task 3", "Description 3", "done", 300, "", now, now),
	}

	// Create ACs for different tasks
	acs := []*entities.AcceptanceCriteriaEntity{
		// Task 1 has 2 ACs
		entities.NewAcceptanceCriteriaEntity("TM-ac-1", "TM-task-1", "AC 1 for task 1", entities.VerificationTypeManual, "Test 1", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-2", "TM-task-1", "AC 2 for task 1", entities.VerificationTypeManual, "Test 2", now, now),
		// Task 2 has 1 AC
		entities.NewAcceptanceCriteriaEntity("TM-ac-3", "TM-task-2", "AC 1 for task 2", entities.VerificationTypeManual, "Test 3", now, now),
		// Task 3 has 3 ACs
		entities.NewAcceptanceCriteriaEntity("TM-ac-4", "TM-task-3", "AC 1 for task 3", entities.VerificationTypeManual, "Test 4", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-5", "TM-task-3", "AC 2 for task 3", entities.VerificationTypeManual, "Test 5", now, now),
		entities.NewAcceptanceCriteriaEntity("TM-ac-6", "TM-task-3", "AC 3 for task 3", entities.VerificationTypeManual, "Test 6", now, now),
	}

	vm := transformers.TransformToIterationDetailViewModel(iteration, tasks, acs, []*entities.DocumentEntity{})

	// Verify TaskACs grouping
	if len(vm.TaskACs) != 3 {
		t.Errorf("expected 3 task groups, got %d", len(vm.TaskACs))
	}

	// Verify Task 1 group
	if vm.TaskACs[0].Task.ID != "TM-task-1" {
		t.Errorf("expected task 1 ID 'TM-task-1', got %q", vm.TaskACs[0].Task.ID)
	}
	if len(vm.TaskACs[0].ACs) != 2 {
		t.Errorf("expected 2 ACs for task 1, got %d", len(vm.TaskACs[0].ACs))
	}

	// Verify Task 2 group
	if vm.TaskACs[1].Task.ID != "TM-task-2" {
		t.Errorf("expected task 2 ID 'TM-task-2', got %q", vm.TaskACs[1].Task.ID)
	}
	if len(vm.TaskACs[1].ACs) != 1 {
		t.Errorf("expected 1 AC for task 2, got %d", len(vm.TaskACs[1].ACs))
	}

	// Verify Task 3 group
	if vm.TaskACs[2].Task.ID != "TM-task-3" {
		t.Errorf("expected task 3 ID 'TM-task-3', got %q", vm.TaskACs[2].Task.ID)
	}
	if len(vm.TaskACs[2].ACs) != 3 {
		t.Errorf("expected 3 ACs for task 3, got %d", len(vm.TaskACs[2].ACs))
	}

	// Verify AC content
	if vm.TaskACs[0].ACs[0].ID != "TM-ac-1" {
		t.Errorf("expected AC ID 'TM-ac-1', got %q", vm.TaskACs[0].ACs[0].ID)
	}
	if vm.TaskACs[0].ACs[1].ID != "TM-ac-2" {
		t.Errorf("expected AC ID 'TM-ac-2', got %q", vm.TaskACs[0].ACs[1].ID)
	}
	if vm.TaskACs[1].ACs[0].ID != "TM-ac-3" {
		t.Errorf("expected AC ID 'TM-ac-3', got %q", vm.TaskACs[1].ACs[0].ID)
	}
}

func TestTransformToIterationDetailViewModel_TaskWithoutACs(t *testing.T) {
	now := time.Now()

	// Create iteration with 2 tasks
	iteration, err := entities.NewIterationEntity(1, "Sprint 1", "Goal", "Deliverable", []string{"TM-task-1", "TM-task-2"}, "current", 100, time.Time{}, time.Time{}, now, now)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}

	// Create tasks
	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "done", 200, "", now, now),
	}

	// Create ACs only for Task 1
	acs := []*entities.AcceptanceCriteriaEntity{
		entities.NewAcceptanceCriteriaEntity("TM-ac-1", "TM-task-1", "AC 1", entities.VerificationTypeManual, "Test 1", now, now),
	}

	vm := transformers.TransformToIterationDetailViewModel(iteration, tasks, acs, []*entities.DocumentEntity{})

	// Verify TaskACs - only Task 1 should have a group (Task 2 has no ACs)
	if len(vm.TaskACs) != 1 {
		t.Errorf("expected 1 task group (only task 1 has ACs), got %d", len(vm.TaskACs))
	}

	if vm.TaskACs[0].Task.ID != "TM-task-1" {
		t.Errorf("expected task group for 'TM-task-1', got %q", vm.TaskACs[0].Task.ID)
	}

	if len(vm.TaskACs[0].ACs) != 1 {
		t.Errorf("expected 1 AC in group, got %d", len(vm.TaskACs[0].ACs))
	}
}
