package viewmodels_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

func TestNewDocumentViewModel(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		title            string
		docType          string
		status           string
		content          string
		trackID          *string
		iterationNumber  *int
		createdAt        string
		updatedAt        string
		expectedID       string
		expectedTitle    string
		expectedType     string
		expectedStatus   string
		expectedContent  string
		expectNilTrackID bool
		expectNilIterNum bool
	}{
		{
			name:             "Basic document with all fields",
			id:               "TM-doc-abc123",
			title:            "Architecture Decision Record",
			docType:          "adr",
			status:           "published",
			content:          "# ADR: Choose Framework\nWe chose Bubble Tea.",
			trackID:          nil,
			iterationNumber:  nil,
			createdAt:        "2025-01-15 10:30",
			updatedAt:        "2025-01-15 10:30",
			expectedID:       "TM-doc-abc123",
			expectedTitle:    "Architecture Decision Record",
			expectedType:     "adr",
			expectedStatus:   "published",
			expectedContent:  "# ADR: Choose Framework\nWe chose Bubble Tea.",
			expectNilTrackID: true,
			expectNilIterNum: true,
		},
		{
			name:             "Document attached to track",
			id:               "TM-doc-def456",
			title:            "Planning Document",
			docType:          "plan",
			status:           "draft",
			content:          "# Phase 1 Plan",
			trackID:          strPtr("TM-track-xyz789"),
			iterationNumber:  nil,
			createdAt:        "2025-01-10 14:20",
			updatedAt:        "2025-01-12 09:15",
			expectedID:       "TM-doc-def456",
			expectedTitle:    "Planning Document",
			expectedType:     "plan",
			expectedStatus:   "draft",
			expectedContent:  "# Phase 1 Plan",
			expectNilTrackID: false,
			expectNilIterNum: true,
		},
		{
			name:             "Document attached to iteration",
			id:               "TM-doc-ghi789",
			title:            "Retrospective",
			docType:          "retrospective",
			status:           "published",
			content:          "# What went well",
			trackID:          nil,
			iterationNumber:  intPtr(5),
			createdAt:        "2025-01-20 16:45",
			updatedAt:        "2025-01-21 11:00",
			expectedID:       "TM-doc-ghi789",
			expectedTitle:    "Retrospective",
			expectedType:     "retrospective",
			expectedStatus:   "published",
			expectedContent:  "# What went well",
			expectNilTrackID: true,
			expectNilIterNum: false,
		},
		{
			name:             "Other document type",
			id:               "TM-doc-jkl012",
			title:            "Technical Specification",
			docType:          "other",
			status:           "archived",
			content:          "Archived spec",
			trackID:          nil,
			iterationNumber:  nil,
			createdAt:        "2024-12-01 08:00",
			updatedAt:        "2024-12-15 17:30",
			expectedID:       "TM-doc-jkl012",
			expectedTitle:    "Technical Specification",
			expectedType:     "other",
			expectedStatus:   "archived",
			expectedContent:  "Archived spec",
			expectNilTrackID: true,
			expectNilIterNum: true,
		},
		{
			name:             "Document with empty content",
			id:               "TM-doc-mno345",
			title:            "Empty Document",
			docType:          "plan",
			status:           "draft",
			content:          "",
			trackID:          nil,
			iterationNumber:  nil,
			createdAt:        "2025-01-15 10:00",
			updatedAt:        "2025-01-15 10:00",
			expectedID:       "TM-doc-mno345",
			expectedTitle:    "Empty Document",
			expectedType:     "plan",
			expectedStatus:   "draft",
			expectedContent:  "",
			expectNilTrackID: true,
			expectNilIterNum: true,
		},
		{
			name:             "Document with long title",
			id:               "TM-doc-pqr678",
			title:            "This is a very long document title that describes in detail the nature and purpose of the document for a complex project feature implementation",
			docType:          "adr",
			status:           "published",
			content:          "Content",
			trackID:          nil,
			iterationNumber:  nil,
			createdAt:        "2025-01-15 10:00",
			updatedAt:        "2025-01-15 10:00",
			expectedID:       "TM-doc-pqr678",
			expectedTitle:    "This is a very long document title that describes in detail the nature and purpose of the document for a complex project feature implementation",
			expectedType:     "adr",
			expectedStatus:   "published",
			expectedContent:  "Content",
			expectNilTrackID: true,
			expectNilIterNum: true,
		},
		{
			name:             "Document with multiline content",
			id:               "TM-doc-stu901",
			title:            "Multiline Document",
			docType:          "retrospective",
			status:           "draft",
			content:          "Line 1\nLine 2\nLine 3\nLine 4",
			trackID:          nil,
			iterationNumber:  intPtr(3),
			createdAt:        "2025-01-15 10:00",
			updatedAt:        "2025-01-15 10:00",
			expectedID:       "TM-doc-stu901",
			expectedTitle:    "Multiline Document",
			expectedType:     "retrospective",
			expectedStatus:   "draft",
			expectedContent:  "Line 1\nLine 2\nLine 3\nLine 4",
			expectNilTrackID: true,
			expectNilIterNum: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := viewmodels.NewDocumentViewModel(
				tt.id,
				tt.title,
				tt.docType,
				tt.status,
				tt.content,
				tt.trackID,
				tt.iterationNumber,
				tt.createdAt,
				tt.updatedAt,
			)

			// Verify all fields are set correctly
			if vm.ID != tt.expectedID {
				t.Errorf("ID: expected %q, got %q", tt.expectedID, vm.ID)
			}
			if vm.Title != tt.expectedTitle {
				t.Errorf("Title: expected %q, got %q", tt.expectedTitle, vm.Title)
			}
			if vm.Type != tt.expectedType {
				t.Errorf("Type: expected %q, got %q", tt.expectedType, vm.Type)
			}
			if vm.Status != tt.expectedStatus {
				t.Errorf("Status: expected %q, got %q", tt.expectedStatus, vm.Status)
			}
			if vm.Content != tt.expectedContent {
				t.Errorf("Content: expected %q, got %q", tt.expectedContent, vm.Content)
			}

			// Verify TrackID
			if tt.expectNilTrackID {
				if vm.TrackID != nil {
					t.Errorf("TrackID: expected nil, got %v", vm.TrackID)
				}
			} else {
				if vm.TrackID == nil {
					t.Errorf("TrackID: expected non-nil, got nil")
				} else if *vm.TrackID != *tt.trackID {
					t.Errorf("TrackID: expected %q, got %q", *tt.trackID, *vm.TrackID)
				}
			}

			// Verify IterationNumber
			if tt.expectNilIterNum {
				if vm.IterationNumber != nil {
					t.Errorf("IterationNumber: expected nil, got %v", vm.IterationNumber)
				}
			} else {
				if vm.IterationNumber == nil {
					t.Errorf("IterationNumber: expected non-nil, got nil")
				} else if *vm.IterationNumber != *tt.iterationNumber {
					t.Errorf("IterationNumber: expected %d, got %d", *tt.iterationNumber, *vm.IterationNumber)
				}
			}

			// Verify timestamps
			if vm.CreatedAt != tt.createdAt {
				t.Errorf("CreatedAt: expected %q, got %q", tt.createdAt, vm.CreatedAt)
			}
			if vm.UpdatedAt != tt.updatedAt {
				t.Errorf("UpdatedAt: expected %q, got %q", tt.updatedAt, vm.UpdatedAt)
			}

			// Verify display fields are zero-initialized (transformer sets these)
			if vm.StatusLabel != "" {
				t.Errorf("StatusLabel: expected empty string, got %q", vm.StatusLabel)
			}
			if vm.StatusColor != "" {
				t.Errorf("StatusColor: expected empty string, got %q", vm.StatusColor)
			}
			if vm.TypeLabel != "" {
				t.Errorf("TypeLabel: expected empty string, got %q", vm.TypeLabel)
			}
			if vm.Icon != "" {
				t.Errorf("Icon: expected empty string, got %q", vm.Icon)
			}
		})
	}
}

func TestDocumentViewModel_FieldMutation(t *testing.T) {
	vm := viewmodels.NewDocumentViewModel(
		"TM-doc-test",
		"Original Title",
		"adr",
		"draft",
		"Original content",
		nil,
		nil,
		"2025-01-15 10:00",
		"2025-01-15 10:00",
	)

	// Verify initial state
	if vm.Title != "Original Title" {
		t.Errorf("expected Title %q, got %q", "Original Title", vm.Title)
	}

	// Test field mutations (for transformer use)
	vm.StatusLabel = "Draft"
	vm.StatusColor = "11"
	vm.TypeLabel = "ADR"
	vm.Icon = "📝"

	if vm.StatusLabel != "Draft" {
		t.Errorf("expected StatusLabel %q, got %q", "Draft", vm.StatusLabel)
	}
	if vm.StatusColor != "11" {
		t.Errorf("expected StatusColor %q, got %q", "11", vm.StatusColor)
	}
	if vm.TypeLabel != "ADR" {
		t.Errorf("expected TypeLabel %q, got %q", "ADR", vm.TypeLabel)
	}
	if vm.Icon != "📝" {
		t.Errorf("expected Icon %q, got %q", "📝", vm.Icon)
	}
}

func TestDocumentViewModel_EdgCases(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{"Single character title", "X"},
		{"Title with special chars", "ADR: Strategy & Roadmap!"},
		{"Title with quotes", "Document: \"The Plan\""},
		{"Title with markdown", "# Plan for **Q1** 2025"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := viewmodels.NewDocumentViewModel(
				"TM-doc-test",
				tt.title,
				"plan",
				"published",
				"content",
				nil,
				nil,
				"2025-01-15 10:00",
				"2025-01-15 10:00",
			)

			if vm.Title != tt.title {
				t.Errorf("expected Title %q, got %q", tt.title, vm.Title)
			}
		})
	}
}

func TestDocumentViewModel_ZeroValues(t *testing.T) {
	vm := viewmodels.NewDocumentViewModel(
		"",
		"",
		"",
		"",
		"",
		nil,
		nil,
		"",
		"",
	)

	if vm.ID != "" {
		t.Errorf("expected empty ID, got %q", vm.ID)
	}
	if vm.Title != "" {
		t.Errorf("expected empty Title, got %q", vm.Title)
	}
	if vm.Type != "" {
		t.Errorf("expected empty Type, got %q", vm.Type)
	}
	if vm.Status != "" {
		t.Errorf("expected empty Status, got %q", vm.Status)
	}
	if vm.Content != "" {
		t.Errorf("expected empty Content, got %q", vm.Content)
	}
	if vm.CreatedAt != "" {
		t.Errorf("expected empty CreatedAt, got %q", vm.CreatedAt)
	}
	if vm.UpdatedAt != "" {
		t.Errorf("expected empty UpdatedAt, got %q", vm.UpdatedAt)
	}
	if vm.TrackID != nil {
		t.Errorf("expected nil TrackID, got %v", vm.TrackID)
	}
	if vm.IterationNumber != nil {
		t.Errorf("expected nil IterationNumber, got %v", vm.IterationNumber)
	}
}

// Helper functions for pointer creation
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
