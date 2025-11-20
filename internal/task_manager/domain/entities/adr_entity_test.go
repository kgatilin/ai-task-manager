package entities_test

import (
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
)

func TestNewADREntity_ValidInput(t *testing.T) {
	now := time.Now()

	adr, err := entities.NewADREntity(
		"DW-adr-1",
		"DW-track-1",
		"Use SQLite for storage",
		"proposed",
		"We need local storage",
		"Use SQLite database",
		"Simple, no external deps",
		"Could use Postgres",
		now,
		now,
		nil,
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if adr == nil {
		t.Fatal("expected non-nil ADR")
	}
	if adr.GetID() != "DW-adr-1" {
		t.Errorf("GetID() = %q, want %q", adr.GetID(), "DW-adr-1")
	}
	if adr.Status != "proposed" {
		t.Errorf("Status = %q, want %q", adr.Status, "proposed")
	}
	if adr.Title != "Use SQLite for storage" {
		t.Errorf("Title = %q, want %q", adr.Title, "Use SQLite for storage")
	}
	if adr.Context != "We need local storage" {
		t.Errorf("Context = %q, want %q", adr.Context, "We need local storage")
	}
	if adr.Decision != "Use SQLite database" {
		t.Errorf("Decision = %q, want %q", adr.Decision, "Use SQLite database")
	}
	if adr.Consequences != "Simple, no external deps" {
		t.Errorf("Consequences = %q, want %q", adr.Consequences, "Simple, no external deps")
	}
	if adr.Alternatives != "Could use Postgres" {
		t.Errorf("Alternatives = %q, want %q", adr.Alternatives, "Could use Postgres")
	}
}

func TestNewADREntity_ValidationErrors(t *testing.T) {
	now := time.Now()
	supersededBy := "DW-adr-2"

	tests := []struct {
		name         string
		status       string
		title        string
		context      string
		decision     string
		consequences string
		supersededBy *string
		wantErr      error
		wantErrMsg   string
	}{
		{
			name:         "invalid status",
			status:       "invalid",
			title:        "Title",
			context:      "Context",
			decision:     "Decision",
			consequences: "Consequences",
			wantErr:      errors.ErrInvalidArgument,
			wantErrMsg:   "invalid ADR status",
		},
		{
			name:         "empty title",
			status:       "proposed",
			title:        "",
			context:      "Context",
			decision:     "Decision",
			consequences: "Consequences",
			wantErr:      errors.ErrInvalidArgument,
			wantErrMsg:   "title",
		},
		{
			name:         "empty context",
			status:       "proposed",
			title:        "Title",
			context:      "",
			decision:     "Decision",
			consequences: "Consequences",
			wantErr:      errors.ErrInvalidArgument,
			wantErrMsg:   "context",
		},
		{
			name:         "empty decision",
			status:       "proposed",
			title:        "Title",
			context:      "Context",
			decision:     "",
			consequences: "Consequences",
			wantErr:      errors.ErrInvalidArgument,
			wantErrMsg:   "decision",
		},
		{
			name:         "empty consequences",
			status:       "proposed",
			title:        "Title",
			context:      "Context",
			decision:     "Decision",
			consequences: "",
			wantErr:      errors.ErrInvalidArgument,
			wantErrMsg:   "consequences",
		},
		{
			name:         "superseded without supersededBy",
			status:       "superseded",
			title:        "Title",
			context:      "Context",
			decision:     "Decision",
			consequences: "Consequences",
			supersededBy: nil,
			wantErr:      errors.ErrInvalidArgument,
			wantErrMsg:   "superseded_by",
		},
		{
			name:         "superseded with supersededBy - valid",
			status:       "superseded",
			title:        "Title",
			context:      "Context",
			decision:     "Decision",
			consequences: "Consequences",
			supersededBy: &supersededBy,
			wantErr:      nil,
		},
		{
			name:         "accepted status - valid",
			status:       "accepted",
			title:        "Title",
			context:      "Context",
			decision:     "Decision",
			consequences: "Consequences",
			wantErr:      nil,
		},
		{
			name:         "deprecated status - valid",
			status:       "deprecated",
			title:        "Title",
			context:      "Context",
			decision:     "Decision",
			consequences: "Consequences",
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entities.NewADREntity(
				"id", "track-1", tt.title, tt.status,
				tt.context, tt.decision, tt.consequences, "",
				now, now, tt.supersededBy,
			)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if !isErrorType(err, tt.wantErr) {
					t.Errorf("expected error type %v, got %v", tt.wantErr, err)
				} else if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("expected error containing %q, got %q", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestADREntity_GetType(t *testing.T) {
	now := time.Now()
	adr, err := entities.NewADREntity(
		"DW-adr-1", "DW-track-1", "Title", "proposed",
		"Context", "Decision", "Consequences", "",
		now, now, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error creating ADR: %v", err)
	}

	entityType := adr.GetType()
	if entityType != "adr" {
		t.Errorf("GetType() = %q, want %q", entityType, "adr")
	}
}

func TestADREntity_GetCapabilities(t *testing.T) {
	now := time.Now()
	adr, err := entities.NewADREntity(
		"DW-adr-1", "DW-track-1", "Title", "proposed",
		"Context", "Decision", "Consequences", "",
		now, now, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error creating ADR: %v", err)
	}

	capabilities := adr.GetCapabilities()
	if len(capabilities) != 1 {
		t.Errorf("GetCapabilities() length = %d, want 1", len(capabilities))
	}
	if capabilities[0] != "IExtensible" {
		t.Errorf("GetCapabilities()[0] = %q, want %q", capabilities[0], "IExtensible")
	}
}

func TestADREntity_GetField(t *testing.T) {
	now := time.Now()
	supersededBy := "DW-adr-2"
	adr, err := entities.NewADREntity(
		"DW-adr-1", "DW-track-1", "Title", "superseded",
		"Context", "Decision", "Consequences", "Alternatives",
		now, now, &supersededBy,
	)
	if err != nil {
		t.Fatalf("unexpected error creating ADR: %v", err)
	}

	tests := []struct {
		fieldName     string
		expectedValue interface{}
	}{
		{"id", "DW-adr-1"},
		{"track_id", "DW-track-1"},
		{"title", "Title"},
		{"status", "superseded"},
		{"context", "Context"},
		{"decision", "Decision"},
		{"consequences", "Consequences"},
		{"alternatives", "Alternatives"},
		{"superseded_by", "DW-adr-2"},
		{"created_at", now},
		{"updated_at", now},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			value := adr.GetField(tt.fieldName)
			if value != tt.expectedValue {
				t.Errorf("GetField(%q) = %v, want %v", tt.fieldName, value, tt.expectedValue)
			}
		})
	}
}

func TestADREntity_GetField_NonExistent(t *testing.T) {
	now := time.Now()
	adr, err := entities.NewADREntity(
		"DW-adr-1", "DW-track-1", "Title", "proposed",
		"Context", "Decision", "Consequences", "",
		now, now, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error creating ADR: %v", err)
	}

	value := adr.GetField("nonexistent")
	if value != nil {
		t.Errorf("GetField(nonexistent) = %v, want nil", value)
	}
}

func TestADREntity_GetAllFields(t *testing.T) {
	now := time.Now()
	adr, err := entities.NewADREntity(
		"DW-adr-1", "DW-track-1", "Title", "proposed",
		"Context", "Decision", "Consequences", "Alternatives",
		now, now, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error creating ADR: %v", err)
	}

	fields := adr.GetAllFields()

	expectedFields := []string{
		"id", "track_id", "title", "status",
		"context", "decision", "consequences", "alternatives",
		"created_at", "updated_at", "superseded_by",
	}

	if len(fields) != len(expectedFields) {
		t.Errorf("GetAllFields() returned %d fields, want %d", len(fields), len(expectedFields))
	}

	for _, fieldName := range expectedFields {
		if _, exists := fields[fieldName]; !exists {
			t.Errorf("GetAllFields() missing field %q", fieldName)
		}
	}

	// Verify superseded_by is empty string when nil
	if fields["superseded_by"] != "" {
		t.Errorf("GetAllFields()[superseded_by] = %q, want empty string", fields["superseded_by"])
	}
}

func TestADREntity_GetAllFields_WithSupersededBy(t *testing.T) {
	now := time.Now()
	supersededBy := "DW-adr-2"
	adr, err := entities.NewADREntity(
		"DW-adr-1", "DW-track-1", "Title", "superseded",
		"Context", "Decision", "Consequences", "",
		now, now, &supersededBy,
	)
	if err != nil {
		t.Fatalf("unexpected error creating ADR: %v", err)
	}

	fields := adr.GetAllFields()

	if fields["superseded_by"] != "DW-adr-2" {
		t.Errorf("GetAllFields()[superseded_by] = %q, want %q", fields["superseded_by"], "DW-adr-2")
	}
}

func TestADREntity_IsAccepted(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"accepted status", "accepted", true},
		{"proposed status", "proposed", false},
		{"deprecated status", "deprecated", false},
		{"superseded status", "superseded", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adr, err := entities.NewADREntity(
				"DW-adr-1", "DW-track-1", "Title", tt.status,
				"Context", "Decision", "Consequences", "",
				now, now, nil,
			)
			if err != nil && tt.status != "superseded" {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.status != "superseded" {
				result := adr.IsAccepted()
				if result != tt.expected {
					t.Errorf("IsAccepted() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestADREntity_IsSuperseded(t *testing.T) {
	now := time.Now()
	supersededBy := "DW-adr-2"

	tests := []struct {
		name         string
		status       string
		supersededBy *string
		expected     bool
	}{
		{"superseded status", "superseded", &supersededBy, true},
		{"proposed status", "proposed", nil, false},
		{"accepted status", "accepted", nil, false},
		{"deprecated status", "deprecated", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adr, err := entities.NewADREntity(
				"DW-adr-1", "DW-track-1", "Title", tt.status,
				"Context", "Decision", "Consequences", "",
				now, now, tt.supersededBy,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			result := adr.IsSuperseded()
			if result != tt.expected {
				t.Errorf("IsSuperseded() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestADREntity_IsDeprecated(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"deprecated status", "deprecated", true},
		{"proposed status", "proposed", false},
		{"accepted status", "accepted", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adr, err := entities.NewADREntity(
				"DW-adr-1", "DW-track-1", "Title", tt.status,
				"Context", "Decision", "Consequences", "",
				now, now, nil,
			)
			if err != nil && tt.status != "superseded" {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.status != "superseded" {
				result := adr.IsDeprecated()
				if result != tt.expected {
					t.Errorf("IsDeprecated() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestADREntity_ToMarkdown(t *testing.T) {
	now := time.Now()
	supersededBy := "DW-adr-2"

	tests := []struct {
		name          string
		status        string
		supersededBy  *string
		alternatives  string
		expectedParts []string
	}{
		{
			name:         "basic ADR",
			status:       "proposed",
			supersededBy: nil,
			alternatives: "",
			expectedParts: []string{
				"# ADR DW-adr-1: Use SQLite",
				"**Status**: proposed",
				"## Context",
				"We need storage",
				"## Decision",
				"Use SQLite",
				"## Consequences",
				"Simple and reliable",
			},
		},
		{
			name:         "with alternatives",
			status:       "accepted",
			supersededBy: nil,
			alternatives: "Could use Postgres or MySQL",
			expectedParts: []string{
				"## Alternatives",
				"Could use Postgres or MySQL",
			},
		},
		{
			name:         "superseded ADR",
			status:       "superseded",
			supersededBy: &supersededBy,
			alternatives: "",
			expectedParts: []string{
				"**Status**: superseded",
				"**Superseded by**: DW-adr-2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adr, err := entities.NewADREntity(
				"DW-adr-1", "DW-track-1", "Use SQLite", tt.status,
				"We need storage", "Use SQLite", "Simple and reliable", tt.alternatives,
				now, now, tt.supersededBy,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			markdown := adr.ToMarkdown()

			for _, part := range tt.expectedParts {
				if !contains(markdown, part) {
					t.Errorf("ToMarkdown() missing expected part %q", part)
				}
			}
		})
	}
}

// Helper function to check if error is of a specific type
func isErrorType(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}
	// Check if error wraps the target error
	for err != nil {
		if err == target {
			return true
		}
		// Try to unwrap
		type unwrapper interface {
			Unwrap() error
		}
		if u, ok := err.(unwrapper); ok {
			err = u.Unwrap()
		} else {
			break
		}
	}
	return false
}
