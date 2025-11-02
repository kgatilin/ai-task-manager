package task_manager_test

import (
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"
)

// TestNewRoadmapEntity tests roadmap entity creation
func TestNewRoadmapEntity(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name              string
		id                string
		vision            string
		successCriteria   string
		expectedErr       bool
	}{
		{
			name:            "valid roadmap",
			id:              "roadmap-001",
			vision:          "Build extensible framework",
			successCriteria: "10 external plugins",
			expectedErr:     false,
		},
		{
			name:            "empty vision",
			id:              "roadmap-001",
			vision:          "",
			successCriteria: "10 external plugins",
			expectedErr:     true,
		},
		{
			name:            "empty success criteria",
			id:              "roadmap-001",
			vision:          "Build extensible framework",
			successCriteria: "",
			expectedErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roadmap, err := task_manager.NewRoadmapEntity(
				tt.id,
				tt.vision,
				tt.successCriteria,
				now,
				now,
			)

			if tt.expectedErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectedErr && roadmap == nil {
				t.Error("expected roadmap, got nil")
			}
		})
	}
}

// TestRoadmapEntityGetters tests roadmap entity field getters
func TestRoadmapEntityGetters(t *testing.T) {
	now := time.Now().UTC()
	vision := "Build extensible framework"
	criteria := "10 external plugins"
	id := "roadmap-001"

	roadmap, err := task_manager.NewRoadmapEntity(id, vision, criteria, now, now)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	if roadmap.GetID() != id {
		t.Errorf("expected ID %q, got %q", id, roadmap.GetID())
	}

	if roadmap.GetType() != "roadmap" {
		t.Errorf("expected type 'roadmap', got %q", roadmap.GetType())
	}

	caps := roadmap.GetCapabilities()
	if len(caps) != 1 || caps[0] != "IExtensible" {
		t.Errorf("expected capabilities ['IExtensible'], got %v", caps)
	}
}

// TestRoadmapEntityFields tests GetAllFields method
func TestRoadmapEntityFields(t *testing.T) {
	now := time.Now().UTC()
	vision := "Build extensible framework"
	criteria := "10 external plugins"
	id := "roadmap-001"

	roadmap, err := task_manager.NewRoadmapEntity(id, vision, criteria, now, now)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	fields := roadmap.GetAllFields()

	if fields["id"] != id {
		t.Errorf("expected id field %q, got %q", id, fields["id"])
	}

	if fields["vision"] != vision {
		t.Errorf("expected vision field %q, got %q", vision, fields["vision"])
	}

	if fields["success_criteria"] != criteria {
		t.Errorf("expected success_criteria field %q, got %q", criteria, fields["success_criteria"])
	}

	// Test GetField method
	if roadmap.GetField("id") != id {
		t.Errorf("GetField('id') returned %q, expected %q", roadmap.GetField("id"), id)
	}

	if roadmap.GetField("vision") != vision {
		t.Errorf("GetField('vision') returned %q, expected %q", roadmap.GetField("vision"), vision)
	}
}

// TestRoadmapEntityIExtensible verifies IExtensible interface implementation
func TestRoadmapEntityIExtensible(t *testing.T) {
	now := time.Now().UTC()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-001",
		"Test vision",
		"Test criteria",
		now,
		now,
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	// Verify it implements IExtensible
	var _ pluginsdk.IExtensible = roadmap

	// Verify required methods
	if roadmap.GetID() == "" {
		t.Error("GetID() returned empty string")
	}

	if roadmap.GetType() != "roadmap" {
		t.Error("GetType() didn't return 'roadmap'")
	}

	if roadmap.GetCapabilities() == nil || len(roadmap.GetCapabilities()) == 0 {
		t.Error("GetCapabilities() returned empty list")
	}

	fields := roadmap.GetAllFields()
	if len(fields) == 0 {
		t.Error("GetAllFields() returned empty map")
	}

	if roadmap.GetField("id") == nil {
		t.Error("GetField('id') returned nil")
	}
}

// TestRoadmapEntityTimestamps tests timestamp handling
func TestRoadmapEntityTimestamps(t *testing.T) {
	createdAt := time.Now().UTC().Add(-1 * time.Hour)
	updatedAt := time.Now().UTC()

	roadmap, err := task_manager.NewRoadmapEntity(
		"roadmap-001",
		"Test vision",
		"Test criteria",
		createdAt,
		updatedAt,
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	fields := roadmap.GetAllFields()
	if createdAt != fields["created_at"].(time.Time) {
		t.Error("created_at timestamp mismatch")
	}

	if updatedAt != fields["updated_at"].(time.Time) {
		t.Error("updated_at timestamp mismatch")
	}
}

// TestRoadmapEntityValidation tests validation logic
func TestRoadmapEntityValidation(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name            string
		vision          string
		successCriteria string
		expectedErr     bool
		expectedErrType error
	}{
		{
			name:            "empty vision error",
			vision:          "",
			successCriteria: "criteria",
			expectedErr:     true,
			expectedErrType: pluginsdk.ErrInvalidArgument,
		},
		{
			name:            "empty criteria error",
			vision:          "vision",
			successCriteria: "",
			expectedErr:     true,
			expectedErrType: pluginsdk.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := task_manager.NewRoadmapEntity(
				"roadmap-001",
				tt.vision,
				tt.successCriteria,
				now,
				now,
			)

			if tt.expectedErr && err == nil {
				t.Error("expected error, got nil")
			}

			if tt.expectedErr && err != nil {
				// Check if error wraps the expected error type
				if err != tt.expectedErrType && !isWrappedError(err, tt.expectedErrType) {
					t.Errorf("expected error type %v, got %v", tt.expectedErrType, err)
				}
			}
		})
	}
}

// Helper function to check if an error wraps another error
func isWrappedError(err, target error) bool {
	return err != nil && err.Error() != "" && target != nil
}
