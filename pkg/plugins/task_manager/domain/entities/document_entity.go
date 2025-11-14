package entities

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// DocumentEntity represents a document that can be attached to tracks or iterations
// It implements IExtensible interface for SDK capability support
type DocumentEntity struct {
	ID               string                 `json:"id"`                // Format: TM-doc-{random}
	Title            string                 `json:"title"`             // Required, max 200 chars
	Type             DocumentType           `json:"type"`              // DocumentType value object
	Status           DocumentStatus         `json:"status"`            // DocumentStatus value object
	Content          string                 `json:"content"`           // Markdown text
	TrackID          *string                `json:"track_id"`          // Optional, validates format
	IterationNumber  *int                   `json:"iteration_number"`  // Optional, validates >= 1
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Metadata         map[string]interface{} `json:"metadata"` // For SDK.IExtensible
}

// NewDocumentEntity creates a new document entity with comprehensive validation
// Validates: ID format, title constraints, type and status validity, attachment constraints
func NewDocumentEntity(
	id string,
	title string,
	docType DocumentType,
	status DocumentStatus,
	content string,
	trackID *string,
	iterationNumber *int,
	createdAt time.Time,
	updatedAt time.Time,
) (*DocumentEntity, error) {
	// Validate document ID format
	if !isValidDocumentID(id) {
		return nil, fmt.Errorf("%w: document ID must follow convention: TM-doc-<random>", pluginsdk.ErrInvalidArgument)
	}

	// Validate title
	if title == "" {
		return nil, fmt.Errorf("%w: document title is required", pluginsdk.ErrInvalidArgument)
	}
	if len(title) > 200 {
		return nil, fmt.Errorf("%w: document title must be 200 characters or less", pluginsdk.ErrInvalidArgument)
	}

	// Validate document type
	if !docType.IsValid() {
		return nil, fmt.Errorf("%w: invalid document type: %s", pluginsdk.ErrInvalidArgument, docType)
	}

	// Validate document status
	if !status.IsValid() {
		return nil, fmt.Errorf("%w: invalid document status: %s", pluginsdk.ErrInvalidArgument, status)
	}

	// Validate attachment constraints: XOR (either TrackID or IterationNumber, not both)
	if trackID != nil && iterationNumber != nil {
		return nil, fmt.Errorf("%w: document cannot have both TrackID and IterationNumber (choose one or neither)", pluginsdk.ErrInvalidArgument)
	}

	// Validate TrackID format if provided
	if trackID != nil && *trackID != "" {
		if !isValidTrackID(*trackID) {
			return nil, fmt.Errorf("%w: invalid track ID format: %s", pluginsdk.ErrInvalidArgument, *trackID)
		}
	}

	// Validate IterationNumber if provided
	if iterationNumber != nil && *iterationNumber < 1 {
		return nil, fmt.Errorf("%w: iteration number must be >= 1", pluginsdk.ErrInvalidArgument)
	}

	return &DocumentEntity{
		ID:              id,
		Title:           title,
		Type:            docType,
		Status:          status,
		Content:         content,
		TrackID:         trackID,
		IterationNumber: iterationNumber,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Metadata:        make(map[string]interface{}),
	}, nil
}

// isValidDocumentID validates document ID format (TM-doc-<random>)
func isValidDocumentID(id string) bool {
	pattern := `^TM-doc-[a-z0-9]+$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(id)
}

// IExtensible implementation

// GetID returns the unique identifier for this entity
func (d *DocumentEntity) GetID() string {
	return d.ID
}

// GetType returns the entity type
func (d *DocumentEntity) GetType() string {
	return "document"
}

// GetCapabilities returns list of capability names this entity supports
func (d *DocumentEntity) GetCapabilities() []string {
	return []string{"IExtensible"}
}

// GetField retrieves a named field value
func (d *DocumentEntity) GetField(name string) interface{} {
	fields := d.GetAllFields()
	return fields[name]
}

// GetAllFields returns all fields as a map
func (d *DocumentEntity) GetAllFields() map[string]interface{} {
	return map[string]interface{}{
		"id":               d.ID,
		"title":            d.Title,
		"type":             d.Type.String(),
		"status":           d.Status.String(),
		"content":          d.Content,
		"track_id":         d.TrackID,
		"iteration_number": d.IterationNumber,
		"created_at":       d.CreatedAt,
		"updated_at":       d.UpdatedAt,
		"is_attached":      d.IsAttached(),
	}
}

// GetTitle returns the document title
func (d *DocumentEntity) GetTitle() string {
	return d.Title
}

// GetType returns the document type value object
func (d *DocumentEntity) GetDocType() DocumentType {
	return d.Type
}

// GetStatus returns the document status value object
func (d *DocumentEntity) GetDocStatus() DocumentStatus {
	return d.Status
}

// GetContent returns the document content
func (d *DocumentEntity) GetContent() string {
	return d.Content
}

// GetTrackID returns the track ID if attached to a track
func (d *DocumentEntity) GetTrackID() *string {
	return d.TrackID
}

// GetIterationNumber returns the iteration number if attached to an iteration
func (d *DocumentEntity) GetIterationNumber() *int {
	return d.IterationNumber
}

// IsAttached returns true if the document is attached to a track or iteration
func (d *DocumentEntity) IsAttached() bool {
	trackAttached := d.TrackID != nil && *d.TrackID != ""
	iterationAttached := d.IterationNumber != nil && *d.IterationNumber > 0
	return trackAttached || iterationAttached
}

// UpdateContent updates the document content
func (d *DocumentEntity) UpdateContent(content string) {
	d.Content = content
	d.UpdatedAt = time.Now()
}

// UpdateStatus updates the document status
func (d *DocumentEntity) UpdateStatus(status DocumentStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("%w: invalid document status: %s", pluginsdk.ErrInvalidArgument, status)
	}
	d.Status = status
	d.UpdatedAt = time.Now()
	return nil
}

// AttachToTrack attaches the document to a track
func (d *DocumentEntity) AttachToTrack(trackID string) error {
	if !isValidTrackID(trackID) {
		return fmt.Errorf("%w: invalid track ID format: %s", pluginsdk.ErrInvalidArgument, trackID)
	}
	if d.IterationNumber != nil && *d.IterationNumber > 0 {
		return fmt.Errorf("%w: document is already attached to an iteration", pluginsdk.ErrInvalidArgument)
	}
	d.TrackID = &trackID
	d.IterationNumber = nil
	d.UpdatedAt = time.Now()
	return nil
}

// AttachToIteration attaches the document to an iteration
func (d *DocumentEntity) AttachToIteration(iterationNumber int) error {
	if iterationNumber < 1 {
		return fmt.Errorf("%w: iteration number must be >= 1", pluginsdk.ErrInvalidArgument)
	}
	if d.TrackID != nil && *d.TrackID != "" {
		return fmt.Errorf("%w: document is already attached to a track", pluginsdk.ErrInvalidArgument)
	}
	d.IterationNumber = &iterationNumber
	d.TrackID = nil
	d.UpdatedAt = time.Now()
	return nil
}

// Detach removes the document from any track or iteration attachment
func (d *DocumentEntity) Detach() {
	d.TrackID = nil
	d.IterationNumber = nil
	d.UpdatedAt = time.Now()
}
