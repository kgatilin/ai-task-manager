package viewmodels

// DocumentListItemViewModel represents a document in a list view.
// Minimal data needed for list row rendering.
type DocumentListItemViewModel struct {
	ID         string // Document ID
	Title      string // Document title
	Type       string // Document type label
	StatusIcon string // Status icon (✓ or ○)
}

// DocumentViewModel represents a document for display in the viewer
type DocumentViewModel struct {
	ID               string
	Title            string
	Type             string // e.g., "adr", "plan", "retrospective", "other"
	Status           string // e.g., "draft", "published", "archived"
	Content          string // Markdown content
	TrackID          *string
	IterationNumber  *int
	CreatedAt        string // Pre-formatted timestamp
	UpdatedAt        string // Pre-formatted timestamp

	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	TypeLabel   string // Human-readable type label
	Icon        string // Status icon
}

// NewDocumentViewModel creates a new document view model
func NewDocumentViewModel(
	id, title, docType, status, content string,
	trackID *string,
	iterationNumber *int,
	createdAt, updatedAt string,
) *DocumentViewModel {
	return &DocumentViewModel{
		ID:              id,
		Title:           title,
		Type:            docType,
		Status:          status,
		Content:         content,
		TrackID:         trackID,
		IterationNumber: iterationNumber,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}
