package entities_test

import (
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

func TestNewDocumentEntity(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		id              string
		title           string
		docType         entities.DocumentType
		status          entities.DocumentStatus
		content         string
		trackID         *string
		iterationNumber *int
		wantErr         bool
		errContains     string
	}{
		{
			name:            "valid document with no attachment",
			id:              "TM-doc-abc123",
			title:           "Test Document",
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Test Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         false,
		},
		{
			name:            "valid document attached to track",
			id:              "TM-doc-def456",
			title:           "Track Document",
			docType:         entities.DocumentTypePlan,
			status:          entities.DocumentStatusPublished,
			content:         "# Plan Content",
			trackID:         strPtr("DW-track-1"),
			iterationNumber: nil,
			wantErr:         false,
		},
		{
			name:            "valid document attached to iteration",
			id:              "TM-doc-ghi789",
			title:           "Iteration Document",
			docType:         entities.DocumentTypeRetrospective,
			status:          entities.DocumentStatusPublished,
			content:         "# Retro Content",
			trackID:         nil,
			iterationNumber: intPtr(1),
			wantErr:         false,
		},
		{
			name:            "invalid document ID format",
			id:              "invalid-id",
			title:           "Test Document",
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         true,
			errContains:     "document ID must follow convention",
		},
		{
			name:            "empty title",
			id:              "TM-doc-abc123",
			title:           "",
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         true,
			errContains:     "document title is required",
		},
		{
			name:            "title too long (> 200 chars)",
			id:              "TM-doc-abc123",
			title:           string(make([]byte, 201)),
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         true,
			errContains:     "must be 200 characters or less",
		},
		{
			name:            "invalid document type",
			id:              "TM-doc-abc123",
			title:           "Test Document",
			docType:         entities.DocumentType("invalid"),
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         true,
			errContains:     "invalid document type",
		},
		{
			name:            "invalid document status",
			id:              "TM-doc-abc123",
			title:           "Test Document",
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatus("invalid"),
			content:         "# Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         true,
			errContains:     "invalid document status",
		},
		{
			name:            "both TrackID and IterationNumber (invalid)",
			id:              "TM-doc-abc123",
			title:           "Test Document",
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         strPtr("DW-track-1"),
			iterationNumber: intPtr(1),
			wantErr:         true,
			errContains:     "cannot have both TrackID and IterationNumber",
		},
		{
			name:            "invalid track ID format",
			id:              "TM-doc-abc123",
			title:           "Test Document",
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         strPtr("invalid-track"),
			iterationNumber: nil,
			wantErr:         true,
			errContains:     "invalid track ID format",
		},
		{
			name:            "invalid iteration number (< 1)",
			id:              "TM-doc-abc123",
			title:           "Test Document",
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         nil,
			iterationNumber: intPtr(0),
			wantErr:         true,
			errContains:     "iteration number must be >= 1",
		},
		{
			name:            "title exactly 200 chars",
			id:              "TM-doc-abc123",
			title:           string(make([]byte, 200)),
			docType:         entities.DocumentTypeADR,
			status:          entities.DocumentStatusDraft,
			content:         "# Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         false,
		},
		{
			name:            "all document types valid",
			id:              "TM-doc-abc123",
			title:           "Test Document",
			docType:         entities.DocumentTypeOther,
			status:          entities.DocumentStatusArchived,
			content:         "# Content",
			trackID:         nil,
			iterationNumber: nil,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := entities.NewDocumentEntity(
				tt.id,
				tt.title,
				tt.docType,
				tt.status,
				tt.content,
				tt.trackID,
				tt.iterationNumber,
				now,
				now,
			)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain expected string %q", err.Error(), tt.errContains)
				}
				if doc != nil {
					t.Errorf("expected nil document but got %v", doc)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if doc == nil {
					t.Fatal("expected document but got nil")
				}
				if doc.GetID() != tt.id {
					t.Errorf("ID mismatch: got %q, want %q", doc.GetID(), tt.id)
				}
				if doc.GetTitle() != tt.title {
					t.Errorf("Title mismatch: got %q, want %q", doc.GetTitle(), tt.title)
				}
			}
		})
	}
}

func TestDocumentEntity_IsAttached(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		trackID         *string
		iterationNumber *int
		wantAttached    bool
	}{
		{
			name:            "unattached (both nil)",
			trackID:         nil,
			iterationNumber: nil,
			wantAttached:    false,
		},
		{
			name:            "attached to track",
			trackID:         strPtr("DW-track-1"),
			iterationNumber: nil,
			wantAttached:    true,
		},
		{
			name:            "attached to iteration",
			trackID:         nil,
			iterationNumber: intPtr(1),
			wantAttached:    true,
		},
		{
			name:            "empty track ID (unattached)",
			trackID:         strPtr(""),
			iterationNumber: nil,
			wantAttached:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := entities.NewDocumentEntity(
				"TM-doc-test",
				"Test Document",
				entities.DocumentTypeADR,
				entities.DocumentStatusDraft,
				"# Content",
				tt.trackID,
				tt.iterationNumber,
				now,
				now,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if doc.IsAttached() != tt.wantAttached {
				t.Errorf("IsAttached() = %v, want %v", doc.IsAttached(), tt.wantAttached)
			}
		})
	}
}

func TestDocumentEntity_GetterMethods(t *testing.T) {
	now := time.Now()
	trackID := "DW-track-1"

	doc, err := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusPublished,
		"# Test Content",
		&trackID,
		nil,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		got      interface{}
		want     interface{}
		field    string
	}{
		{
			name:  "GetID",
			got:   doc.GetID(),
			want:  "TM-doc-test",
			field: "id",
		},
		{
			name:  "GetType",
			got:   doc.GetType(),
			want:  "document",
			field: "type",
		},
		{
			name:  "GetTitle",
			got:   doc.GetTitle(),
			want:  "Test Document",
			field: "title",
		},
		{
			name:  "GetDocType",
			got:   doc.GetDocType(),
			want:  entities.DocumentTypeADR,
			field: "doc_type",
		},
		{
			name:  "GetDocStatus",
			got:   doc.GetDocStatus(),
			want:  entities.DocumentStatusPublished,
			field: "doc_status",
		},
		{
			name:  "GetContent",
			got:   doc.GetContent(),
			want:  "# Test Content",
			field: "content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s: got %v, want %v", tt.field, tt.got, tt.want)
			}
		})
	}

	// Test pointer getters
	if doc.GetTrackID() == nil || *doc.GetTrackID() != trackID {
		t.Errorf("GetTrackID: got %v, want %s", doc.GetTrackID(), trackID)
	}
	if doc.GetIterationNumber() != nil {
		t.Errorf("GetIterationNumber: got %v, want nil", doc.GetIterationNumber())
	}
}

func TestDocumentEntity_UpdateContent(t *testing.T) {
	now := time.Now()
	doc, _ := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"# Original Content",
		nil,
		nil,
		now,
		now,
	)

	originalUpdatedAt := doc.UpdatedAt
	doc.UpdateContent("# New Content")

	if doc.GetContent() != "# New Content" {
		t.Errorf("Content not updated: got %q", doc.GetContent())
	}
	if doc.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt was not modified")
	}
}

func TestDocumentEntity_UpdateStatus(t *testing.T) {
	now := time.Now()
	doc, _ := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"# Content",
		nil,
		nil,
		now,
		now,
	)

	// Valid status update
	err := doc.UpdateStatus(entities.DocumentStatusPublished)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if doc.GetDocStatus() != entities.DocumentStatusPublished {
		t.Errorf("Status not updated: got %s", doc.GetDocStatus())
	}

	// Invalid status update
	err = doc.UpdateStatus(entities.DocumentStatus("invalid"))
	if err == nil {
		t.Error("expected error for invalid status but got nil")
	}
}

func TestDocumentEntity_AttachToTrack(t *testing.T) {
	now := time.Now()
	doc, _ := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"# Content",
		nil,
		nil,
		now,
		now,
	)

	// Valid attachment
	err := doc.AttachToTrack("DW-track-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !doc.IsAttached() {
		t.Error("document should be attached")
	}
	if doc.GetTrackID() == nil || *doc.GetTrackID() != "DW-track-1" {
		t.Errorf("TrackID not set correctly: got %v", doc.GetTrackID())
	}

	// Invalid track ID format
	err = doc.AttachToTrack("invalid-track")
	if err == nil {
		t.Error("expected error for invalid track ID")
	}
}

func TestDocumentEntity_AttachToIteration(t *testing.T) {
	now := time.Now()
	doc, _ := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"# Content",
		nil,
		nil,
		now,
		now,
	)

	// Valid attachment
	err := doc.AttachToIteration(5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !doc.IsAttached() {
		t.Error("document should be attached")
	}
	if doc.GetIterationNumber() == nil || *doc.GetIterationNumber() != 5 {
		t.Errorf("IterationNumber not set correctly: got %v", doc.GetIterationNumber())
	}

	// Invalid iteration number
	err = doc.AttachToIteration(0)
	if err == nil {
		t.Error("expected error for invalid iteration number")
	}
}

func TestDocumentEntity_AttachmentMutualExclusivity(t *testing.T) {
	now := time.Now()
	trackID := "DW-track-1"
	doc, _ := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"# Content",
		&trackID,
		nil,
		now,
		now,
	)

	// Try to attach to iteration when already attached to track
	err := doc.AttachToIteration(1)
	if err == nil {
		t.Error("expected error when attaching to iteration while attached to track")
	}
}

func TestDocumentEntity_Detach(t *testing.T) {
	now := time.Now()
	trackID := "DW-track-1"
	doc, _ := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"# Content",
		&trackID,
		nil,
		now,
		now,
	)

	doc.Detach()

	if doc.IsAttached() {
		t.Error("document should not be attached after detach")
	}
	if doc.GetTrackID() != nil {
		t.Errorf("TrackID should be nil after detach: got %v", doc.GetTrackID())
	}
}

func TestDocumentEntity_IExtensibleImplementation(t *testing.T) {
	now := time.Now()
	doc, _ := entities.NewDocumentEntity(
		"TM-doc-test",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusPublished,
		"# Content",
		nil,
		nil,
		now,
		now,
	)

	// Test GetType
	if doc.GetType() != "document" {
		t.Errorf("GetType: got %q, want %q", doc.GetType(), "document")
	}

	// Test GetCapabilities
	caps := doc.GetCapabilities()
	if len(caps) == 0 {
		t.Error("GetCapabilities returned empty list")
	}

	// Test GetField
	if doc.GetField("title") != "Test Document" {
		t.Errorf("GetField(title) failed")
	}

	// Test GetAllFields
	fields := doc.GetAllFields()
	if fields["id"] != "TM-doc-test" {
		t.Errorf("GetAllFields: id field mismatch")
	}
	if fields["title"] != "Test Document" {
		t.Errorf("GetAllFields: title field mismatch")
	}
	if fields["type"] != "adr" {
		t.Errorf("GetAllFields: type field mismatch")
	}
	if fields["status"] != "published" {
		t.Errorf("GetAllFields: status field mismatch")
	}
}

// Helper functions

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
