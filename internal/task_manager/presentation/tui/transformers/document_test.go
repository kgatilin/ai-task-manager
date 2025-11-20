package transformers_test

import (
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/transformers"
)

// Test helper to create a DocumentEntity
func mustCreateDocument(
	id, title string,
	docType entities.DocumentType,
	status entities.DocumentStatus,
	content string,
	trackID *string,
	iterationNumber *int,
	createdAt, updatedAt time.Time,
) *entities.DocumentEntity {
	doc, err := entities.NewDocumentEntity(
		id, title, docType, status, content,
		trackID, iterationNumber, createdAt, updatedAt,
	)
	if err != nil {
		panic(err)
	}
	return doc
}

func TestFormatDocumentType(t *testing.T) {
	tests := []struct {
		name     string
		docType  entities.DocumentType
		expected string
	}{
		{"ADR type", entities.DocumentTypeADR, "ADR"},
		{"Plan type", entities.DocumentTypePlan, "Plan"},
		{"Retrospective type", entities.DocumentTypeRetrospective, "Retrospective"},
		{"Other type", entities.DocumentTypeOther, "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: formatDocumentType is not exported, testing via TransformToDocumentViewModel
			doc := mustCreateDocument(
				"TM-doc-test", "Test", tt.docType, entities.DocumentStatusDraft, "content",
				nil, nil, time.Now(), time.Now(),
			)
			vm := transformers.TransformToDocumentViewModel(doc)
			if vm.TypeLabel != tt.expected {
				t.Errorf("expected TypeLabel %q, got %q", tt.expected, vm.TypeLabel)
			}
		})
	}
}

func TestGetDocumentStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Draft status", string(entities.DocumentStatusDraft), "📝"},
		{"Published status", string(entities.DocumentStatusPublished), "✓"},
		{"Archived status", string(entities.DocumentStatusArchived), "📦"},
		{"Unknown status", "unknown", "📄"},
		{"Empty string", "", "📄"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetDocumentStatusIcon(tt.status)
			if result != tt.expected {
				t.Errorf("GetDocumentStatusIcon(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetDocumentStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Draft status", string(entities.DocumentStatusDraft), "Draft"},
		{"Published status", string(entities.DocumentStatusPublished), "Published"},
		{"Archived status", string(entities.DocumentStatusArchived), "Archived"},
		{"Unknown status returns as-is", "unknown_status", "unknown_status"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetDocumentStatusLabel(tt.status)
			if result != tt.expected {
				t.Errorf("GetDocumentStatusLabel(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetDocumentStatusColor(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Draft status", string(entities.DocumentStatusDraft), "11"},
		{"Published status", string(entities.DocumentStatusPublished), "10"},
		{"Archived status", string(entities.DocumentStatusArchived), "240"},
		{"Unknown status defaults to gray", "unknown", "240"},
		{"Empty string defaults to gray", "", "240"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetDocumentStatusColor(tt.status)
			if result != tt.expected {
				t.Errorf("GetDocumentStatusColor(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetDocumentTypeLabel(t *testing.T) {
	tests := []struct {
		name     string
		docType  string
		expected string
	}{
		{"ADR type", string(entities.DocumentTypeADR), "ADR"},
		{"Plan type", string(entities.DocumentTypePlan), "Plan"},
		{"Retrospective type", string(entities.DocumentTypeRetrospective), "Retrospective"},
		{"Other type", string(entities.DocumentTypeOther), "Other"},
		{"Unknown type returns as-is", "custom_type", "custom_type"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.GetDocumentTypeLabel(tt.docType)
			if result != tt.expected {
				t.Errorf("GetDocumentTypeLabel(%q) = %q, want %q", tt.docType, result, tt.expected)
			}
		})
	}
}

func TestTransformToDocumentViewModel(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	later := time.Date(2025, 1, 20, 14, 45, 0, 0, time.UTC)

	tests := []struct {
		name                string
		id                  string
		title               string
		docType             entities.DocumentType
		status              entities.DocumentStatus
		content             string
		trackID             *string
		iterationNumber     *int
		expectedID          string
		expectedTitle       string
		expectedType        string
		expectedTypeLabel   string
		expectedStatus      string
		expectedStatusLabel string
		expectedStatusColor string
		expectedIcon        string
		expectedTrackIDNil  bool
		expectedIterNumNil  bool
	}{
		{
			name:                "ADR document - published",
			id:                  "TM-doc-adr001",
			title:               "Architecture Decision Record",
			docType:             entities.DocumentTypeADR,
			status:              entities.DocumentStatusPublished,
			content:             "# Context\nWe need...\n\n# Decision\nWe chose...",
			trackID:             nil,
			iterationNumber:     nil,
			expectedID:          "TM-doc-adr001",
			expectedTitle:       "Architecture Decision Record",
			expectedType:        "adr",
			expectedTypeLabel:   "ADR",
			expectedStatus:      "published",
			expectedStatusLabel: "Published",
			expectedStatusColor: "10",
			expectedIcon:        "✓",
			expectedTrackIDNil:  true,
			expectedIterNumNil:  true,
		},
		{
			name:                "Plan document - draft",
			id:                  "TM-doc-plan001",
			title:               "Q1 2025 Planning",
			docType:             entities.DocumentTypePlan,
			status:              entities.DocumentStatusDraft,
			content:             "# Phase 1",
			trackID:             strPtr("TM-track-1"),
			iterationNumber:     nil,
			expectedID:          "TM-doc-plan001",
			expectedTitle:       "Q1 2025 Planning",
			expectedType:        "plan",
			expectedTypeLabel:   "Plan",
			expectedStatus:      "draft",
			expectedStatusLabel: "Draft",
			expectedStatusColor: "11",
			expectedIcon:        "📝",
			expectedTrackIDNil:  false,
			expectedIterNumNil:  true,
		},
		{
			name:                "Retrospective - archived",
			id:                  "TM-doc-retro001",
			title:               "Sprint 10 Retrospective",
			docType:             entities.DocumentTypeRetrospective,
			status:              entities.DocumentStatusArchived,
			content:             "# What went well",
			trackID:             nil,
			iterationNumber:     intPtr(10),
			expectedID:          "TM-doc-retro001",
			expectedTitle:       "Sprint 10 Retrospective",
			expectedType:        "retrospective",
			expectedTypeLabel:   "Retrospective",
			expectedStatus:      "archived",
			expectedStatusLabel: "Archived",
			expectedStatusColor: "240",
			expectedIcon:        "📦",
			expectedTrackIDNil:  true,
			expectedIterNumNil:  false,
		},
		{
			name:                "Other document type",
			id:                  "TM-doc-other001",
			title:               "Technical Notes",
			docType:             entities.DocumentTypeOther,
			status:              entities.DocumentStatusPublished,
			content:             "Some notes",
			trackID:             nil,
			iterationNumber:     nil,
			expectedID:          "TM-doc-other001",
			expectedTitle:       "Technical Notes",
			expectedType:        "other",
			expectedTypeLabel:   "Other",
			expectedStatus:      "published",
			expectedStatusLabel: "Published",
			expectedStatusColor: "10",
			expectedIcon:        "✓",
			expectedTrackIDNil:  true,
			expectedIterNumNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := mustCreateDocument(
				tt.id, tt.title, tt.docType, tt.status, tt.content,
				tt.trackID, tt.iterationNumber, now, later,
			)

			vm := transformers.TransformToDocumentViewModel(doc)

			// Verify basic fields
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
			if vm.Content != tt.content {
				t.Errorf("Content: expected %q, got %q", tt.content, vm.Content)
			}

			// Verify display fields
			if vm.TypeLabel != tt.expectedTypeLabel {
				t.Errorf("TypeLabel: expected %q, got %q", tt.expectedTypeLabel, vm.TypeLabel)
			}
			if vm.StatusLabel != tt.expectedStatusLabel {
				t.Errorf("StatusLabel: expected %q, got %q", tt.expectedStatusLabel, vm.StatusLabel)
			}
			if vm.StatusColor != tt.expectedStatusColor {
				t.Errorf("StatusColor: expected %q, got %q", tt.expectedStatusColor, vm.StatusColor)
			}
			if vm.Icon != tt.expectedIcon {
				t.Errorf("Icon: expected %q, got %q", tt.expectedIcon, vm.Icon)
			}

			// Verify TrackID
			if tt.expectedTrackIDNil {
				if vm.TrackID != nil {
					t.Errorf("TrackID: expected nil, got %v", vm.TrackID)
				}
			} else {
				if vm.TrackID == nil {
					t.Errorf("TrackID: expected non-nil")
				}
			}

			// Verify IterationNumber
			if tt.expectedIterNumNil {
				if vm.IterationNumber != nil {
					t.Errorf("IterationNumber: expected nil, got %v", vm.IterationNumber)
				}
			} else {
				if vm.IterationNumber == nil {
					t.Errorf("IterationNumber: expected non-nil")
				}
			}

			// Verify timestamps are formatted
			if vm.CreatedAt == "" {
				t.Errorf("CreatedAt: expected non-empty formatted string")
			}
			if vm.UpdatedAt == "" {
				t.Errorf("UpdatedAt: expected non-empty formatted string")
			}
		})
	}
}

func TestTransformDocumentsToListItems_WithDocuments(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		documents     []*entities.DocumentEntity
		expectedCount int
		expectedItems []struct {
			id         string
			title      string
			docType    string
			statusIcon string
		}
	}{
		{
			name: "Single ADR document",
			documents: []*entities.DocumentEntity{
				mustCreateDocument(
					"TM-doc-adr001", "Architecture Decision",
					entities.DocumentTypeADR, entities.DocumentStatusPublished,
					"Decision content", nil, nil, now, now,
				),
			},
			expectedCount: 1,
			expectedItems: []struct {
				id         string
				title      string
				docType    string
				statusIcon string
			}{
				{"TM-doc-adr001", "Architecture Decision", "ADR", "✓"},
			},
		},
		{
			name: "Multiple documents with different types",
			documents: []*entities.DocumentEntity{
				mustCreateDocument(
					"TM-doc-1", "ADR Title",
					entities.DocumentTypeADR, entities.DocumentStatusPublished,
					"", nil, nil, now, now,
				),
				mustCreateDocument(
					"TM-doc-2", "Plan Title",
					entities.DocumentTypePlan, entities.DocumentStatusDraft,
					"", nil, nil, now, now,
				),
				mustCreateDocument(
					"TM-doc-3", "Retro Title",
					entities.DocumentTypeRetrospective, entities.DocumentStatusArchived,
					"", nil, nil, now, now,
				),
			},
			expectedCount: 3,
			expectedItems: []struct {
				id         string
				title      string
				docType    string
				statusIcon string
			}{
				{"TM-doc-1", "ADR Title", "ADR", "✓"},
				{"TM-doc-2", "Plan Title", "Plan", "○"},
				{"TM-doc-3", "Retro Title", "Retrospective", "✗"},
			},
		},
		{
			name: "Documents with same type, different status",
			documents: []*entities.DocumentEntity{
				mustCreateDocument(
					"TM-doc-draft", "Draft Plan",
					entities.DocumentTypePlan, entities.DocumentStatusDraft,
					"", nil, nil, now, now,
				),
				mustCreateDocument(
					"TM-doc-pub", "Published Plan",
					entities.DocumentTypePlan, entities.DocumentStatusPublished,
					"", nil, nil, now, now,
				),
				mustCreateDocument(
					"TM-doc-arch", "Archived Plan",
					entities.DocumentTypePlan, entities.DocumentStatusArchived,
					"", nil, nil, now, now,
				),
			},
			expectedCount: 3,
			expectedItems: []struct {
				id         string
				title      string
				docType    string
				statusIcon string
			}{
				{"TM-doc-draft", "Draft Plan", "Plan", "○"},
				{"TM-doc-pub", "Published Plan", "Plan", "✓"},
				{"TM-doc-arch", "Archived Plan", "Plan", "✗"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := transformers.TransformDocumentsToListItems(tt.documents)

			if len(items) != tt.expectedCount {
				t.Errorf("expected %d items, got %d", tt.expectedCount, len(items))
			}

			for i, expectedItem := range tt.expectedItems {
				if i >= len(items) {
					t.Fatalf("item %d: list too short", i)
				}

				item := items[i]
				if item.ID != expectedItem.id {
					t.Errorf("item[%d].ID: expected %q, got %q", i, expectedItem.id, item.ID)
				}
				if item.Title != expectedItem.title {
					t.Errorf("item[%d].Title: expected %q, got %q", i, expectedItem.title, item.Title)
				}
				if item.Type != expectedItem.docType {
					t.Errorf("item[%d].Type: expected %q, got %q", i, expectedItem.docType, item.Type)
				}
				if item.StatusIcon != expectedItem.statusIcon {
					t.Errorf("item[%d].StatusIcon: expected %q, got %q", i, expectedItem.statusIcon, item.StatusIcon)
				}
			}
		})
	}
}

func TestTransformDocumentsToListItems_EmptySlice(t *testing.T) {
	items := transformers.TransformDocumentsToListItems([](*entities.DocumentEntity){})

	if len(items) != 0 {
		t.Errorf("expected empty slice, got %d items", len(items))
	}

	// Verify it returns an empty slice, not nil
	if items == nil {
		t.Error("expected empty slice (not nil)")
	}
}

func TestTransformDocumentsToListItems_WithNilElements(t *testing.T) {
	now := time.Now()

	documents := []*entities.DocumentEntity{
		mustCreateDocument(
			"TM-doc-1", "First",
			entities.DocumentTypeADR, entities.DocumentStatusPublished,
			"", nil, nil, now, now,
		),
		nil, // Nil element should be skipped
		mustCreateDocument(
			"TM-doc-2", "Second",
			entities.DocumentTypePlan, entities.DocumentStatusDraft,
			"", nil, nil, now, now,
		),
		nil, // Another nil element
	}

	items := transformers.TransformDocumentsToListItems(documents)

	// Should only have 2 items (nil elements skipped)
	if len(items) != 2 {
		t.Errorf("expected 2 items (nil skipped), got %d", len(items))
	}

	if len(items) >= 1 {
		if items[0].ID != "TM-doc-1" {
			t.Errorf("item[0].ID: expected TM-doc-1, got %q", items[0].ID)
		}
	}

	if len(items) >= 2 {
		if items[1].ID != "TM-doc-2" {
			t.Errorf("item[1].ID: expected TM-doc-2, got %q", items[1].ID)
		}
	}
}

func TestTransformDocumentsToListItems_LargeList(t *testing.T) {
	now := time.Now()

	// Create 100 documents
	documents := make([]*entities.DocumentEntity, 100)
	for i := 0; i < 100; i++ {
		docType := entities.DocumentTypeADR
		if i%3 == 1 {
			docType = entities.DocumentTypePlan
		} else if i%3 == 2 {
			docType = entities.DocumentTypeRetrospective
		}

		status := entities.DocumentStatusPublished
		if i%2 == 0 {
			status = entities.DocumentStatusDraft
		}

		id := generateID(i)
		documents[i] = mustCreateDocument(
			id, "Document "+id,
			docType, status, "", nil, nil, now, now,
		)
	}

	items := transformers.TransformDocumentsToListItems(documents)

	if len(items) != 100 {
		t.Errorf("expected 100 items, got %d", len(items))
	}

	// Verify first and last items
	if len(items) > 0 {
		if !startsWith(items[0].ID, "TM-doc-") {
			t.Errorf("first item ID has wrong format: %q", items[0].ID)
		}
	}

	if len(items) > 99 {
		if !startsWith(items[99].ID, "TM-doc-") {
			t.Errorf("last item ID has wrong format: %q", items[99].ID)
		}
	}
}

func TestTransformDocumentsToListItems_AllDocumentTypes(t *testing.T) {
	now := time.Now()

	documents := []*entities.DocumentEntity{
		mustCreateDocument(
			"TM-doc-adr", "ADR",
			entities.DocumentTypeADR, entities.DocumentStatusPublished, "", nil, nil, now, now,
		),
		mustCreateDocument(
			"TM-doc-plan", "Plan",
			entities.DocumentTypePlan, entities.DocumentStatusDraft, "", nil, nil, now, now,
		),
		mustCreateDocument(
			"TM-doc-retro", "Retro",
			entities.DocumentTypeRetrospective, entities.DocumentStatusArchived, "", nil, nil, now, now,
		),
		mustCreateDocument(
			"TM-doc-other", "Other",
			entities.DocumentTypeOther, entities.DocumentStatusPublished, "", nil, nil, now, now,
		),
	}

	items := transformers.TransformDocumentsToListItems(documents)

	typeMap := map[string]string{
		"TM-doc-adr":   "ADR",
		"TM-doc-plan":  "Plan",
		"TM-doc-retro": "Retrospective",
		"TM-doc-other": "Other",
	}

	for _, item := range items {
		expectedType, exists := typeMap[item.ID]
		if !exists {
			t.Errorf("unexpected document ID: %q", item.ID)
			continue
		}
		if item.Type != expectedType {
			t.Errorf("document %q: expected Type %q, got %q", item.ID, expectedType, item.Type)
		}
	}
}

func TestTransformDocumentsToListItems_AllDocumentStatuses(t *testing.T) {
	now := time.Now()

	documents := []*entities.DocumentEntity{
		mustCreateDocument(
			"TM-doc-draft", "Draft",
			entities.DocumentTypePlan, entities.DocumentStatusDraft, "", nil, nil, now, now,
		),
		mustCreateDocument(
			"TM-doc-pub", "Published",
			entities.DocumentTypePlan, entities.DocumentStatusPublished, "", nil, nil, now, now,
		),
		mustCreateDocument(
			"TM-doc-arch", "Archived",
			entities.DocumentTypePlan, entities.DocumentStatusArchived, "", nil, nil, now, now,
		),
	}

	items := transformers.TransformDocumentsToListItems(documents)

	statusIconMap := map[string]string{
		"TM-doc-draft": "○",
		"TM-doc-pub":   "✓",
		"TM-doc-arch":  "✗",
	}

	for _, item := range items {
		expectedIcon, exists := statusIconMap[item.ID]
		if !exists {
			t.Errorf("unexpected document ID: %q", item.ID)
			continue
		}
		if item.StatusIcon != expectedIcon {
			t.Errorf("document %q: expected StatusIcon %q, got %q", item.ID, expectedIcon, item.StatusIcon)
		}
	}
}

func TestFormatDocumentType_DefaultCase(t *testing.T) {
	// Test the default case via TransformToDocumentViewModel which calls formatDocumentType
	// formatDocumentType is not exported, but we can test through the visible behavior
	doc := mustCreateDocument(
		"TM-doc-test", "Test",
		entities.DocumentTypeOther, entities.DocumentStatusDraft,
		"", nil, nil, time.Now(), time.Now(),
	)
	vm := transformers.TransformToDocumentViewModel(doc)
	if vm.TypeLabel != "Other" {
		t.Errorf("expected TypeLabel Other, got %q", vm.TypeLabel)
	}
}

func TestGetDocumentStatusIcon_DefaultCase(t *testing.T) {
	// getDocumentStatusIcon is not exported, test through GetDocumentStatusIcon
	result := transformers.GetDocumentStatusIcon("invalid_status")
	if result != "📄" {
		t.Errorf("expected default icon 📄, got %q", result)
	}
}

func TestFormatTime_ZeroTime(t *testing.T) {
	// formatTime is not exported, but we can test it by creating a doc with zero timestamps
	// Since NewDocumentEntity validates, we create with valid times and check formatter through VM
	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	zeroTime := time.Time{}

	doc, err := entities.NewDocumentEntity(
		"TM-doc-test", "Test",
		entities.DocumentTypeADR, entities.DocumentStatusDraft,
		"", nil, nil, zeroTime, now,
	)
	if err == nil {
		// If creation succeeded, check that zero time formats as empty
		vm := transformers.TransformToDocumentViewModel(doc)
		// CreatedAt should be empty string for zero time
		if vm.CreatedAt != "" {
			t.Errorf("expected empty CreatedAt for zero time, got %q", vm.CreatedAt)
		}
		if vm.UpdatedAt == "" {
			t.Errorf("expected non-empty UpdatedAt for non-zero time")
		}
	}
}

func TestFormatTime_NonZeroTime(t *testing.T) {
	createdAt := time.Date(2025, 1, 15, 10, 30, 45, 0, time.UTC)
	updatedAt := time.Date(2025, 1, 20, 14, 45, 30, 0, time.UTC)

	doc := mustCreateDocument(
		"TM-doc-test", "Test",
		entities.DocumentTypeADR, entities.DocumentStatusPublished,
		"", nil, nil, createdAt, updatedAt,
	)

	vm := transformers.TransformToDocumentViewModel(doc)

	// Verify format is "2006-01-02 15:04"
	expectedCreated := "2025-01-15 10:30"
	expectedUpdated := "2025-01-20 14:45"

	if vm.CreatedAt != expectedCreated {
		t.Errorf("CreatedAt: expected %q, got %q", expectedCreated, vm.CreatedAt)
	}
	if vm.UpdatedAt != expectedUpdated {
		t.Errorf("UpdatedAt: expected %q, got %q", expectedUpdated, vm.UpdatedAt)
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func generateID(i int) string {
	return "TM-doc-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
