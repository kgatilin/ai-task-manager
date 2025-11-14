package entities_test

import (
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

func TestNewRoadmapEntity_ValidInput(t *testing.T) {
	now := time.Now()

	roadmap, err := entities.NewRoadmapEntity(
		"roadmap-1",
		"Build amazing product",
		"Launch with 1000 users",
		now,
		now,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if roadmap == nil {
		t.Fatal("expected non-nil roadmap")
	}
	if roadmap.GetID() != "roadmap-1" {
		t.Errorf("GetID() = %q, want %q", roadmap.GetID(), "roadmap-1")
	}
	if roadmap.Vision != "Build amazing product" {
		t.Errorf("Vision = %q, want %q", roadmap.Vision, "Build amazing product")
	}
	if roadmap.SuccessCriteria != "Launch with 1000 users" {
		t.Errorf("SuccessCriteria = %q, want %q", roadmap.SuccessCriteria, "Launch with 1000 users")
	}
	if roadmap.CreatedAt.Unix() != now.Unix() {
		t.Errorf("CreatedAt = %v, want %v", roadmap.CreatedAt.Unix(), now.Unix())
	}
	if roadmap.UpdatedAt.Unix() != now.Unix() {
		t.Errorf("UpdatedAt = %v, want %v", roadmap.UpdatedAt.Unix(), now.Unix())
	}
}

func TestNewRoadmapEntity_ValidationErrors(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		id              string
		vision          string
		successCriteria string
		wantErr         bool
		errContains     string
	}{
		{
			name:            "empty vision",
			id:              "roadmap-1",
			vision:          "",
			successCriteria: "Some criteria",
			wantErr:         true,
			errContains:     "vision must be non-empty",
		},
		{
			name:            "empty success criteria",
			id:              "roadmap-1",
			vision:          "Some vision",
			successCriteria: "",
			wantErr:         true,
			errContains:     "success criteria must be non-empty",
		},
		{
			name:            "both empty",
			id:              "roadmap-1",
			vision:          "",
			successCriteria: "",
			wantErr:         true,
			errContains:     "vision must be non-empty",
		},
		{
			name:            "whitespace only vision",
			id:              "roadmap-1",
			vision:          "   ",
			successCriteria: "Some criteria",
			wantErr:         false, // Constructor doesn't trim, so this passes
		},
		{
			name:            "whitespace only criteria",
			id:              "roadmap-1",
			vision:          "Some vision",
			successCriteria: "   ",
			wantErr:         false, // Constructor doesn't trim, so this passes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roadmap, err := entities.NewRoadmapEntity(tt.id, tt.vision, tt.successCriteria, now, now)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				if roadmap != nil {
					t.Errorf("expected nil roadmap on error, got %+v", roadmap)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if roadmap == nil {
					t.Fatal("expected non-nil roadmap")
				}
			}
		})
	}
}

func TestNewRoadmapEntity_ErrorIsInvalidArgument(t *testing.T) {
	now := time.Now()

	// Test empty vision
	_, err := entities.NewRoadmapEntity("roadmap-1", "", "criteria", now, now)
	if err == nil {
		t.Fatal("expected error for empty vision, got nil")
	}
	if !contains(err.Error(), pluginsdk.ErrInvalidArgument.Error()) {
		t.Errorf("expected error to wrap ErrInvalidArgument, got %v", err)
	}

	// Test empty success criteria
	_, err = entities.NewRoadmapEntity("roadmap-1", "vision", "", now, now)
	if err == nil {
		t.Fatal("expected error for empty success criteria, got nil")
	}
	if !contains(err.Error(), pluginsdk.ErrInvalidArgument.Error()) {
		t.Errorf("expected error to wrap ErrInvalidArgument, got %v", err)
	}
}

func TestRoadmapEntity_GetType(t *testing.T) {
	now := time.Now()
	roadmap, err := entities.NewRoadmapEntity("r1", "vision", "criteria", now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entityType := roadmap.GetType()
	expectedType := "roadmap"

	if entityType != expectedType {
		t.Errorf("GetType() = %q, want %q", entityType, expectedType)
	}
}

func TestRoadmapEntity_GetCapabilities(t *testing.T) {
	now := time.Now()
	roadmap, err := entities.NewRoadmapEntity("r1", "vision", "criteria", now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	caps := roadmap.GetCapabilities()

	if len(caps) != 1 {
		t.Errorf("GetCapabilities() returned %d capabilities, want 1", len(caps))
	}

	if len(caps) > 0 && caps[0] != "IExtensible" {
		t.Errorf("GetCapabilities()[0] = %q, want %q", caps[0], "IExtensible")
	}

	// Verify it contains IExtensible
	hasIExtensible := false
	for _, cap := range caps {
		if cap == "IExtensible" {
			hasIExtensible = true
			break
		}
	}
	if !hasIExtensible {
		t.Errorf("GetCapabilities() = %v, missing 'IExtensible'", caps)
	}
}

func TestRoadmapEntity_GetField(t *testing.T) {
	now := time.Now()
	roadmap, err := entities.NewRoadmapEntity("r1", "My Vision", "My Criteria", now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		fieldName     string
		expectedValue interface{}
		checkType     bool // whether to check the type instead of exact value
	}{
		{
			fieldName:     "id",
			expectedValue: "r1",
			checkType:     false,
		},
		{
			fieldName:     "vision",
			expectedValue: "My Vision",
			checkType:     false,
		},
		{
			fieldName:     "success_criteria",
			expectedValue: "My Criteria",
			checkType:     false,
		},
		{
			fieldName:     "created_at",
			expectedValue: now,
			checkType:     true, // Check type only, value will be time.Time
		},
		{
			fieldName:     "updated_at",
			expectedValue: now,
			checkType:     true, // Check type only, value will be time.Time
		},
		{
			fieldName:     "non_existent",
			expectedValue: nil,
			checkType:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			value := roadmap.GetField(tt.fieldName)

			if tt.expectedValue == nil {
				if value != nil {
					t.Errorf("GetField(%q) = %v, want nil", tt.fieldName, value)
				}
			} else if tt.checkType {
				// For time fields, just verify they're not nil
				if value == nil {
					t.Errorf("GetField(%q) returned nil, expected time.Time", tt.fieldName)
				}
				if _, ok := value.(time.Time); !ok {
					t.Errorf("GetField(%q) returned %T, expected time.Time", tt.fieldName, value)
				}
			} else {
				if value != tt.expectedValue {
					t.Errorf("GetField(%q) = %v, want %v", tt.fieldName, value, tt.expectedValue)
				}
			}
		})
	}
}

func TestRoadmapEntity_GetAllFields(t *testing.T) {
	now := time.Now()
	roadmap, err := entities.NewRoadmapEntity("r1", "test vision", "test criteria", now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields := roadmap.GetAllFields()

	// Verify all expected fields are present
	expectedFields := []string{"id", "vision", "success_criteria", "created_at", "updated_at"}
	for _, fieldName := range expectedFields {
		if _, exists := fields[fieldName]; !exists {
			t.Errorf("GetAllFields() missing field %q", fieldName)
		}
	}

	// Verify field values
	if fields["id"] != "r1" {
		t.Errorf("fields[id] = %v, want %q", fields["id"], "r1")
	}
	if fields["vision"] != "test vision" {
		t.Errorf("fields[vision] = %v, want %q", fields["vision"], "test vision")
	}
	if fields["success_criteria"] != "test criteria" {
		t.Errorf("fields[success_criteria] = %v, want %q", fields["success_criteria"], "test criteria")
	}

	// Verify time fields are present and correct type
	if createdAt, ok := fields["created_at"].(time.Time); !ok {
		t.Errorf("fields[created_at] is not time.Time, got %T", fields["created_at"])
	} else if createdAt.Unix() != now.Unix() {
		t.Errorf("fields[created_at] = %v, want %v", createdAt.Unix(), now.Unix())
	}

	if updatedAt, ok := fields["updated_at"].(time.Time); !ok {
		t.Errorf("fields[updated_at] is not time.Time, got %T", fields["updated_at"])
	} else if updatedAt.Unix() != now.Unix() {
		t.Errorf("fields[updated_at] = %v, want %v", updatedAt.Unix(), now.Unix())
	}
}

func TestRoadmapEntity_GetAllFields_Independence(t *testing.T) {
	now := time.Now()
	roadmap, err := entities.NewRoadmapEntity("r1", "vision", "criteria", now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get fields map
	fields1 := roadmap.GetAllFields()
	fields2 := roadmap.GetAllFields()

	// Modify one map
	fields1["id"] = "modified"

	// Verify the other map is unaffected (tests if map is copied)
	if fields2["id"] == "modified" {
		t.Error("GetAllFields() returns same map instance, expected independent copies")
	}
	if fields2["id"] != "r1" {
		t.Errorf("fields2[id] = %v, want %q", fields2["id"], "r1")
	}
}

func TestRoadmapEntity_GetID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		id         string
		expectedID string
	}{
		{
			name:       "simple id",
			id:         "roadmap-1",
			expectedID: "roadmap-1",
		},
		{
			name:       "complex id",
			id:         "roadmap-prod-2025",
			expectedID: "roadmap-prod-2025",
		},
		{
			name:       "short id",
			id:         "r1",
			expectedID: "r1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roadmap, err := entities.NewRoadmapEntity(tt.id, "vision", "criteria", now, now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			id := roadmap.GetID()
			if id != tt.expectedID {
				t.Errorf("GetID() = %q, want %q", id, tt.expectedID)
			}
		})
	}
}

func TestRoadmapEntity_IExtensibleInterface(t *testing.T) {
	now := time.Now()
	roadmap, err := entities.NewRoadmapEntity("r1", "vision", "criteria", now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it implements IExtensible interface
	var _ pluginsdk.IExtensible = roadmap

	// Test all IExtensible methods
	if roadmap.GetID() == "" {
		t.Error("GetID() returned empty string")
	}
	if roadmap.GetType() == "" {
		t.Error("GetType() returned empty string")
	}
	if len(roadmap.GetCapabilities()) == 0 {
		t.Error("GetCapabilities() returned empty slice")
	}
	if roadmap.GetField("id") == nil {
		t.Error("GetField(id) returned nil")
	}
	if len(roadmap.GetAllFields()) == 0 {
		t.Error("GetAllFields() returned empty map")
	}
}

func TestRoadmapEntity_ZeroTimeValues(t *testing.T) {
	zeroTime := time.Time{}

	roadmap, err := entities.NewRoadmapEntity("r1", "vision", "criteria", zeroTime, zeroTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !roadmap.CreatedAt.IsZero() {
		t.Errorf("expected zero time for CreatedAt, got %v", roadmap.CreatedAt)
	}
	if !roadmap.UpdatedAt.IsZero() {
		t.Errorf("expected zero time for UpdatedAt, got %v", roadmap.UpdatedAt)
	}
}

func TestRoadmapEntity_DifferentTimestamps(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	roadmap, err := entities.NewRoadmapEntity("r1", "vision", "criteria", createdAt, updatedAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if roadmap.CreatedAt.Unix() != createdAt.Unix() {
		t.Errorf("CreatedAt = %v, want %v", roadmap.CreatedAt.Unix(), createdAt.Unix())
	}
	if roadmap.UpdatedAt.Unix() != updatedAt.Unix() {
		t.Errorf("UpdatedAt = %v, want %v", roadmap.UpdatedAt.Unix(), updatedAt.Unix())
	}
	if !roadmap.UpdatedAt.After(roadmap.CreatedAt) {
		t.Error("expected UpdatedAt to be after CreatedAt")
	}
}

func TestRoadmapEntity_LongStrings(t *testing.T) {
	now := time.Now()

	longVision := "This is a very long vision statement that describes the future state we want to achieve with extensive detail about all aspects of the product and its impact on users and the market. " +
		"It includes multiple paragraphs and covers various strategic objectives, market positioning, and long-term goals that span several years into the future."

	longCriteria := "These are comprehensive success criteria that include multiple measurable outcomes: 1) Achieve 1 million active users within 12 months, " +
		"2) Reach 95% uptime SLA, 3) Generate $10M ARR, 4) Expand to 50 countries, 5) Achieve Net Promoter Score of 70+, and many more detailed metrics."

	roadmap, err := entities.NewRoadmapEntity("r1", longVision, longCriteria, now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if roadmap.Vision != longVision {
		t.Error("long vision was not preserved correctly")
	}
	if roadmap.SuccessCriteria != longCriteria {
		t.Error("long success criteria was not preserved correctly")
	}
}

func TestRoadmapEntity_SpecialCharacters(t *testing.T) {
	now := time.Now()

	specialVision := "Vision with special chars: émojis 🚀, symbols @#$%, quotes \"test\", and newlines\n\nlike this"
	specialCriteria := "Criteria with unicode: 日本語, math: ∑∫∂, and emojis: 👍✅"

	roadmap, err := entities.NewRoadmapEntity("r1", specialVision, specialCriteria, now, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if roadmap.Vision != specialVision {
		t.Error("special characters in vision were not preserved correctly")
	}
	if roadmap.SuccessCriteria != specialCriteria {
		t.Error("special characters in success criteria were not preserved correctly")
	}
}
