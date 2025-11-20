package entities_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

func TestIsValidTrackStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"valid not-started", "not-started", true},
		{"valid in-progress", "in-progress", true},
		{"valid complete", "complete", true},
		{"valid blocked", "blocked", true},
		{"valid waiting", "waiting", true},
		{"invalid empty", "", false},
		{"invalid unknown", "unknown", false},
		{"invalid case", "COMPLETE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entities.IsValidTrackStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidTrackStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestIsValidTaskStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"valid todo", "todo", true},
		{"valid in-progress", "in-progress", true},
		{"valid done", "done", true},
		{"invalid empty", "", false},
		{"invalid unknown", "unknown", false},
		{"invalid case", "DONE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entities.IsValidTaskStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidTaskStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestIsValidIterationStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"valid planned", "planned", true},
		{"valid current", "current", true},
		{"valid complete", "complete", true},
		{"invalid empty", "", false},
		{"invalid unknown", "unknown", false},
		{"invalid case", "PLANNED", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entities.IsValidIterationStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidIterationStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestIsValidADRStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"valid proposed", "proposed", true},
		{"valid accepted", "accepted", true},
		{"valid deprecated", "deprecated", true},
		{"valid superseded", "superseded", true},
		{"invalid empty", "", false},
		{"invalid unknown", "unknown", false},
		{"invalid case", "PROPOSED", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entities.IsValidADRStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidADRStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
