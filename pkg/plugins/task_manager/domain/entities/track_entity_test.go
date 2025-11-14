package entities_test

import (
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

func TestNewTrackEntity(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		id            string
		roadmapID     string
		title         string
		description   string
		status        string
		rank          int
		dependencies  []string
		wantErr       bool
		errContains   string
	}{
		{
			name:        "valid new format ID",
			id:          "DW-track-1",
			roadmapID:   "roadmap-1",
			title:       "Test Track",
			description: "Test Description",
			status:      "not-started",
			rank:        500,
			dependencies: []string{},
			wantErr:     false,
		},
		{
			name:        "valid old format ID",
			id:          "track-core-framework",
			roadmapID:   "roadmap-1",
			title:       "Test Track",
			description: "Test Description",
			status:      "in-progress",
			rank:        100,
			dependencies: []string{},
			wantErr:     false,
		},
		{
			name:         "invalid track ID format",
			id:           "invalid_id",
			roadmapID:    "roadmap-1",
			title:        "Test Track",
			description:  "Test Description",
			status:       "not-started",
			rank:         500,
			dependencies: []string{},
			wantErr:      true,
			errContains:  "track ID must follow convention",
		},
		{
			name:         "invalid status",
			id:           "DW-track-1",
			roadmapID:    "roadmap-1",
			title:        "Test Track",
			description:  "Test Description",
			status:       "invalid-status",
			rank:         500,
			dependencies: []string{},
			wantErr:      true,
			errContains:  "invalid track status",
		},
		{
			name:         "rank too low",
			id:           "DW-track-1",
			roadmapID:    "roadmap-1",
			title:        "Test Track",
			description:  "Test Description",
			status:       "not-started",
			rank:         0,
			dependencies: []string{},
			wantErr:      true,
			errContains:  "invalid track rank",
		},
		{
			name:         "rank too high",
			id:           "DW-track-1",
			roadmapID:    "roadmap-1",
			title:        "Test Track",
			description:  "Test Description",
			status:       "not-started",
			rank:         1001,
			dependencies: []string{},
			wantErr:      true,
			errContains:  "invalid track rank",
		},
		{
			name:         "self-dependency",
			id:           "DW-track-1",
			roadmapID:    "roadmap-1",
			title:        "Test Track",
			description:  "Test Description",
			status:       "not-started",
			rank:         500,
			dependencies: []string{"DW-track-1"},
			wantErr:      true,
			errContains:  "track cannot depend on itself",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track, err := entities.NewTrackEntity(
				tt.id, tt.roadmapID, tt.title, tt.description,
				tt.status, tt.rank, tt.dependencies, now, now,
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
				if track == nil {
					t.Fatal("expected non-nil track")
				}
				if track.ID != tt.id {
					t.Errorf("ID = %q, want %q", track.ID, tt.id)
				}
				if track.Status != tt.status {
					t.Errorf("Status = %q, want %q", track.Status, tt.status)
				}
				if track.Rank != tt.rank {
					t.Errorf("Rank = %d, want %d", track.Rank, tt.rank)
				}
			}
		})
	}
}

func TestTrackEntity_TransitionTo(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  string
		toStatus    string
		wantErr     bool
		errContains string
	}{
		// Valid transitions
		{"not-started to in-progress", "not-started", "in-progress", false, ""},
		{"not-started to blocked", "not-started", "blocked", false, ""},
		{"not-started to waiting", "not-started", "waiting", false, ""},
		{"in-progress to complete", "in-progress", "complete", false, ""},
		{"in-progress to blocked", "in-progress", "blocked", false, ""},
		{"in-progress to waiting", "in-progress", "waiting", false, ""},
		{"blocked to in-progress", "blocked", "in-progress", false, ""},
		{"blocked to waiting", "blocked", "waiting", false, ""},
		{"waiting to in-progress", "waiting", "in-progress", false, ""},
		{"complete to complete", "complete", "complete", false, ""},

		// Invalid transitions - complete is terminal
		{"complete to in-progress", "complete", "in-progress", true, "complete is terminal"},
		{"complete to blocked", "complete", "blocked", true, "complete is terminal"},
		{"complete to waiting", "complete", "waiting", true, "complete is terminal"},
		{"complete to not-started", "complete", "not-started", true, "complete is terminal"},

		// Invalid status
		{"to invalid status", "not-started", "invalid-status", true, "invalid track status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track := &entities.TrackEntity{
				ID:          "DW-track-1",
				RoadmapID:   "roadmap-1",
				Title:       "Test Track",
				Description: "Test Description",
				Status:      tt.fromStatus,
				Rank:        500,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			err := track.TransitionTo(tt.toStatus)

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
				if track.Status != tt.toStatus {
					t.Errorf("Status = %q, want %q", track.Status, tt.toStatus)
				}
			}
		})
	}
}

func TestTrackEntity_AddDependency(t *testing.T) {
	track := &entities.TrackEntity{
		ID:           "DW-track-1",
		RoadmapID:    "roadmap-1",
		Title:        "Test Track",
		Status:       "not-started",
		Rank:         500,
		Dependencies: []string{},
	}

	// Add first dependency
	err := track.AddDependency("DW-track-2")
	if err != nil {
		t.Errorf("unexpected error adding first dependency: %v", err)
	}
	if len(track.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(track.Dependencies))
	}

	// Add second dependency
	err = track.AddDependency("DW-track-3")
	if err != nil {
		t.Errorf("unexpected error adding second dependency: %v", err)
	}
	if len(track.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(track.Dependencies))
	}

	// Attempt to add duplicate dependency
	err = track.AddDependency("DW-track-2")
	if err == nil {
		t.Error("expected error for duplicate dependency, got nil")
	}

	// Attempt to add self-dependency
	err = track.AddDependency("DW-track-1")
	if err == nil {
		t.Error("expected error for self-dependency, got nil")
	} else if !contains(err.Error(), "cannot depend on itself") {
		t.Errorf("expected error about self-dependency, got %q", err.Error())
	}
}

func TestTrackEntity_RemoveDependency(t *testing.T) {
	track := &entities.TrackEntity{
		ID:           "DW-track-1",
		RoadmapID:    "roadmap-1",
		Title:        "Test Track",
		Status:       "not-started",
		Rank:         500,
		Dependencies: []string{"DW-track-2", "DW-track-3"},
	}

	// Remove existing dependency
	err := track.RemoveDependency("DW-track-2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(track.Dependencies) != 1 {
		t.Errorf("expected 1 dependency after removal, got %d", len(track.Dependencies))
	}
	if track.Dependencies[0] != "DW-track-3" {
		t.Errorf("expected DW-track-3, got %q", track.Dependencies[0])
	}

	// Attempt to remove non-existent dependency
	err = track.RemoveDependency("DW-track-999")
	if err == nil {
		t.Error("expected error for non-existent dependency, got nil")
	}
}

func TestTrackEntity_HasDependency(t *testing.T) {
	track := &entities.TrackEntity{
		ID:           "DW-track-1",
		Dependencies: []string{"DW-track-2", "DW-track-3"},
	}

	if !track.HasDependency("DW-track-2") {
		t.Error("expected HasDependency to return true for DW-track-2")
	}
	if !track.HasDependency("DW-track-3") {
		t.Error("expected HasDependency to return true for DW-track-3")
	}
	if track.HasDependency("DW-track-999") {
		t.Error("expected HasDependency to return false for DW-track-999")
	}
}

func TestTrackEntity_GetProgress(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected float64
	}{
		{"complete", "complete", 1.0},
		{"in-progress", "in-progress", 0.5},
		{"not-started", "not-started", 0.0},
		{"blocked", "blocked", 0.0},
		{"waiting", "waiting", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track := &entities.TrackEntity{Status: tt.status}
			progress := track.GetProgress()
			if progress != tt.expected {
				t.Errorf("GetProgress() = %v, want %v", progress, tt.expected)
			}
		})
	}
}

func TestTrackEntity_IsBlocked(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"blocked", "blocked", true},
		{"not-started", "not-started", false},
		{"in-progress", "in-progress", false},
		{"complete", "complete", false},
		{"waiting", "waiting", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track := &entities.TrackEntity{Status: tt.status}
			isBlocked := track.IsBlocked()
			if isBlocked != tt.expected {
				t.Errorf("IsBlocked() = %v, want %v", isBlocked, tt.expected)
			}
		})
	}
}

// SDK Interface Tests

func TestTrackEntity_GetID(t *testing.T) {
	track := &entities.TrackEntity{ID: "DW-track-42"}
	if got := track.GetID(); got != "DW-track-42" {
		t.Errorf("GetID() = %q, want %q", got, "DW-track-42")
	}
}

func TestTrackEntity_GetType(t *testing.T) {
	track := &entities.TrackEntity{}
	if got := track.GetType(); got != "track" {
		t.Errorf("GetType() = %q, want %q", got, "track")
	}
}

func TestTrackEntity_GetCapabilities(t *testing.T) {
	track := &entities.TrackEntity{}
	capabilities := track.GetCapabilities()

	expected := []string{"IExtensible", "ITrackable"}
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

func TestTrackEntity_GetField(t *testing.T) {
	track := &entities.TrackEntity{
		ID:          "DW-track-1",
		RoadmapID:   "roadmap-1",
		Title:       "Test Track",
		Description: "Test Description",
		Status:      "in-progress",
		Rank:        500,
	}

	tests := []struct {
		field    string
		expected interface{}
	}{
		{"id", "DW-track-1"},
		{"roadmap_id", "roadmap-1"},
		{"title", "Test Track"},
		{"description", "Test Description"},
		{"status", "in-progress"},
		{"rank", 500},
		{"progress", 0.5},
		{"is_blocked", false},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			got := track.GetField(tt.field)
			if got != tt.expected {
				t.Errorf("GetField(%q) = %v, want %v", tt.field, got, tt.expected)
			}
		})
	}
}

func TestTrackEntity_GetAllFields(t *testing.T) {
	now := time.Now()
	track := &entities.TrackEntity{
		ID:           "DW-track-1",
		RoadmapID:    "roadmap-1",
		Title:        "Test Track",
		Description:  "Test Description",
		Status:       "in-progress",
		Rank:         500,
		Dependencies: []string{"DW-track-2"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	fields := track.GetAllFields()

	// Verify all expected fields are present
	expectedFields := []string{
		"id", "roadmap_id", "title", "description",
		"status", "rank", "dependencies",
		"created_at", "updated_at", "progress", "is_blocked",
	}

	for _, field := range expectedFields {
		if _, exists := fields[field]; !exists {
			t.Errorf("GetAllFields() missing field %q", field)
		}
	}

	// Verify some key values
	if fields["id"] != "DW-track-1" {
		t.Errorf("GetAllFields()[\"id\"] = %v, want %v", fields["id"], "DW-track-1")
	}
	if fields["status"] != "in-progress" {
		t.Errorf("GetAllFields()[\"status\"] = %v, want %v", fields["status"], "in-progress")
	}
	if fields["progress"] != 0.5 {
		t.Errorf("GetAllFields()[\"progress\"] = %v, want %v", fields["progress"], 0.5)
	}
}

func TestTrackEntity_GetStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"not-started", "not-started"},
		{"in-progress", "in-progress"},
		{"complete", "complete"},
		{"blocked", "blocked"},
		{"waiting", "waiting"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track := &entities.TrackEntity{Status: tt.status}
			if got := track.GetStatus(); got != tt.status {
				t.Errorf("GetStatus() = %q, want %q", got, tt.status)
			}
		})
	}
}
