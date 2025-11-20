package entities_test

import (
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

func TestNewIterationEntity(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		number      int
		iterName    string
		goal        string
		deliverable string
		taskIDs     []string
		status      string
		rank        float64
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid iteration",
			number:      1,
			iterName:    "Sprint 1",
			goal:        "Foundation",
			deliverable: "Core framework",
			taskIDs:     []string{},
			status:      "planned",
			rank:        500,
			wantErr:     false,
		},
		{
			name:        "valid with tasks",
			number:      2,
			iterName:    "Sprint 2",
			goal:        "Features",
			deliverable: "Feature complete",
			taskIDs:     []string{"DW-task-1", "DW-task-2"},
			status:      "current",
			rank:        100,
			wantErr:     false,
		},
		{
			name:        "number zero",
			number:      0,
			iterName:    "Sprint 0",
			goal:        "Setup",
			deliverable: "Initial setup",
			taskIDs:     []string{},
			status:      "planned",
			rank:        500,
			wantErr:     true,
			errContains: "iteration number must be positive",
		},
		{
			name:        "number negative",
			number:      -1,
			iterName:    "Sprint -1",
			goal:        "Back to the past",
			deliverable: "Time machine",
			taskIDs:     []string{},
			status:      "planned",
			rank:        500,
			wantErr:     true,
			errContains: "iteration number must be positive",
		},
		{
			name:        "invalid status",
			number:      1,
			iterName:    "Sprint 1",
			goal:        "Foundation",
			deliverable: "Core framework",
			taskIDs:     []string{},
			status:      "invalid-status",
			rank:        500,
			wantErr:     true,
			errContains: "invalid iteration status",
		},
		{
			name:        "rank too low",
			number:      1,
			iterName:    "Sprint 1",
			goal:        "Foundation",
			deliverable: "Core framework",
			taskIDs:     []string{},
			status:      "planned",
			rank:        0,
			wantErr:     true,
			errContains: "invalid iteration rank",
		},
		{
			name:        "rank too high",
			number:      1,
			iterName:    "Sprint 1",
			goal:        "Foundation",
			deliverable: "Core framework",
			taskIDs:     []string{},
			status:      "planned",
			rank:        1001,
			wantErr:     true,
			errContains: "invalid iteration rank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iteration, err := entities.NewIterationEntity(
				tt.number, tt.iterName, tt.goal, tt.deliverable,
				tt.taskIDs, tt.status, tt.rank,
				time.Time{}, time.Time{}, now, now,
			)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if iteration == nil {
					t.Fatal("expected non-nil iteration")
				}
				if iteration.Number != tt.number {
					t.Errorf("Number = %d, want %d", iteration.Number, tt.number)
				}
				if iteration.Status != tt.status {
					t.Errorf("Status = %q, want %q", iteration.Status, tt.status)
				}
				if iteration.Rank != tt.rank {
					t.Errorf("Rank = %f, want %f", iteration.Rank, tt.rank)
				}
			}
		})
	}
}

func TestIterationEntity_TransitionTo(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  string
		toStatus    string
		wantErr     bool
		errContains string
	}{
		// Valid transitions
		{"planned to current", "planned", "current", false, ""},
		{"current to complete", "current", "complete", false, ""},
		{"planned to planned (no-op)", "planned", "planned", false, ""},
		{"current to current (no-op)", "current", "current", false, ""},
		{"complete to complete (no-op)", "complete", "complete", false, ""},

		// Invalid transitions
		{"planned to complete (skip current)", "planned", "complete", true, "can only transition from planned to current"},
		{"current to planned (backward)", "current", "planned", true, "can only transition from current to complete"},
		{"complete to planned (reopen)", "complete", "planned", true, "cannot transition from complete"},
		{"complete to current (reopen)", "complete", "current", true, "cannot transition from complete"},

		// Invalid status
		{"to invalid status", "planned", "invalid-status", true, "invalid iteration status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iteration := &entities.IterationEntity{
				Number:      1,
				Name:        "Sprint 1",
				Goal:        "Foundation",
				Status:      tt.fromStatus,
				Rank:        500,
				Deliverable: "Core framework",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			err := iteration.TransitionTo(tt.toStatus)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if iteration.Status != tt.toStatus {
					t.Errorf("Status = %q, want %q", iteration.Status, tt.toStatus)
				}

				// Verify timestamps are set
				if tt.toStatus == "current" && iteration.StartedAt == nil {
					t.Error("expected StartedAt to be set when transitioning to current")
				}
				if tt.toStatus == "complete" && iteration.CompletedAt == nil {
					t.Error("expected CompletedAt to be set when transitioning to complete")
				}
			}
		})
	}
}

func TestIterationEntity_AddTask(t *testing.T) {
	iteration := &entities.IterationEntity{
		Number:  1,
		Name:    "Sprint 1",
		TaskIDs: []string{},
	}

	// Add first task
	err := iteration.AddTask("DW-task-1")
	if err != nil {
		t.Errorf("unexpected error adding first task: %v", err)
	}
	if len(iteration.TaskIDs) != 1 {
		t.Errorf("expected 1 task, got %d", len(iteration.TaskIDs))
	}

	// Add second task
	err = iteration.AddTask("DW-task-2")
	if err != nil {
		t.Errorf("unexpected error adding second task: %v", err)
	}
	if len(iteration.TaskIDs) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(iteration.TaskIDs))
	}

	// Attempt to add duplicate task
	err = iteration.AddTask("DW-task-1")
	if err == nil {
		t.Error("expected error for duplicate task, got nil")
	} else if !contains(err.Error(), "task already in iteration") {
		t.Errorf("expected error about duplicate task, got %q", err.Error())
	}
}

func TestIterationEntity_RemoveTask(t *testing.T) {
	iteration := &entities.IterationEntity{
		Number:  1,
		Name:    "Sprint 1",
		TaskIDs: []string{"DW-task-1", "DW-task-2", "DW-task-3"},
	}

	// Remove existing task
	err := iteration.RemoveTask("DW-task-2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(iteration.TaskIDs) != 2 {
		t.Errorf("expected 2 tasks after removal, got %d", len(iteration.TaskIDs))
	}
	if iteration.HasTask("DW-task-2") {
		t.Error("expected task to be removed")
	}

	// Attempt to remove non-existent task
	err = iteration.RemoveTask("DW-task-999")
	if err == nil {
		t.Error("expected error for non-existent task, got nil")
	} else if !contains(err.Error(), "task not in iteration") {
		t.Errorf("expected error about missing task, got %q", err.Error())
	}
}

func TestIterationEntity_HasTask(t *testing.T) {
	iteration := &entities.IterationEntity{
		Number:  1,
		Name:    "Sprint 1",
		TaskIDs: []string{"DW-task-1", "DW-task-2"},
	}

	if !iteration.HasTask("DW-task-1") {
		t.Error("expected HasTask to return true for DW-task-1")
	}
	if !iteration.HasTask("DW-task-2") {
		t.Error("expected HasTask to return true for DW-task-2")
	}
	if iteration.HasTask("DW-task-999") {
		t.Error("expected HasTask to return false for DW-task-999")
	}
}

func TestIterationEntity_GetTaskCount(t *testing.T) {
	tests := []struct {
		name     string
		taskIDs  []string
		expected int
	}{
		{"no tasks", []string{}, 0},
		{"one task", []string{"DW-task-1"}, 1},
		{"multiple tasks", []string{"DW-task-1", "DW-task-2", "DW-task-3"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iteration := &entities.IterationEntity{TaskIDs: tt.taskIDs}
			count := iteration.GetTaskCount()
			if count != tt.expected {
				t.Errorf("GetTaskCount() = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestIterationEntity_GetID(t *testing.T) {
	iteration := &entities.IterationEntity{Number: 42}
	id := iteration.GetID()
	if id != "42" {
		t.Errorf("GetID() = %q, want \"42\"", id)
	}
}

// SDK Interface Tests

func TestIterationEntity_GetType(t *testing.T) {
	iteration := &entities.IterationEntity{}
	if got := iteration.GetType(); got != "iteration" {
		t.Errorf("GetType() = %q, want %q", got, "iteration")
	}
}

func TestIterationEntity_GetCapabilities(t *testing.T) {
	iteration := &entities.IterationEntity{}
	capabilities := iteration.GetCapabilities()

	expected := []string{"IExtensible"}
	if len(capabilities) != len(expected) {
		t.Errorf("GetCapabilities() length = %d, want %d", len(capabilities), len(expected))
		return
	}

	for i, cap := range capabilities {
		if cap != expected[i] {
			t.Errorf("GetCapabilities()[%d] = %q, want %q", i, cap, expected[i])
		}
	}
}

func TestIterationEntity_GetField(t *testing.T) {
	now := time.Now()
	iteration := &entities.IterationEntity{
		Number:      1,
		Name:        "Sprint 1",
		Goal:        "Foundation",
		Deliverable: "Core framework",
		Status:      "current",
		Rank:        500,
		TaskIDs:     []string{"DW-task-1"},
		StartedAt:   &now,
	}

	tests := []struct {
		field    string
		expected interface{}
	}{
		{"number", 1},
		{"name", "Sprint 1"},
		{"goal", "Foundation"},
		{"deliverable", "Core framework"},
		{"status", "current"},
		{"rank", 500.0},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			got := iteration.GetField(tt.field)
			if got != tt.expected {
				t.Errorf("GetField(%q) = %v, want %v", tt.field, got, tt.expected)
			}
		})
	}
}

func TestIterationEntity_GetAllFields(t *testing.T) {
	now := time.Now()
	iteration := &entities.IterationEntity{
		Number:      1,
		Name:        "Sprint 1",
		Goal:        "Foundation",
		Deliverable: "Core framework",
		Status:      "current",
		Rank:        500.0,
		TaskIDs:     []string{"DW-task-1", "DW-task-2"},
		StartedAt:   &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	fields := iteration.GetAllFields()

	// Verify all expected fields are present
	expectedFields := []string{
		"number", "name", "goal", "deliverable",
		"status", "rank", "task_ids",
		"started_at", "completed_at", "created_at", "updated_at",
	}

	for _, field := range expectedFields {
		if _, exists := fields[field]; !exists {
			t.Errorf("GetAllFields() missing field %q", field)
		}
	}

	// Verify some key values
	if fields["number"] != 1 {
		t.Errorf("GetAllFields()[\"number\"] = %v, want %v", fields["number"], 1)
	}
	if fields["status"] != "current" {
		t.Errorf("GetAllFields()[\"status\"] = %v, want %v", fields["status"], "current")
	}
	taskIDs, ok := fields["task_ids"].([]string)
	if !ok || len(taskIDs) != 2 {
		t.Errorf("GetAllFields()[\"task_ids\"] = %v, want slice of 2 task IDs", fields["task_ids"])
	}
}
