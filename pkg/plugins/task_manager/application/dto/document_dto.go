package dto

import "time"

// CreateDocumentDTO represents input for creating a new document
type CreateDocumentDTO struct {
	Title            string // Required, max 200 chars
	Type             string // Required, one of: adr, plan, retrospective, other
	Status           string // Required, one of: draft, published, archived
	Content          string // Required, markdown text
	TrackID          *string // Optional, validates format TM-track-X
	IterationNumber  *int   // Optional, validates >= 1
}

// UpdateDocumentDTO represents input for updating a document
type UpdateDocumentDTO struct {
	ID               string  // Required
	Content          *string // Optional, markdown text to update
	Status           *string // Optional, new status
	TrackID          *string // Optional, new track attachment
	IterationNumber  *int    // Optional, new iteration attachment
	Detach           bool    // If true, remove all attachments
}

// DocumentViewDTO represents the output representation of a document
type DocumentViewDTO struct {
	ID               string     // Document ID (TM-doc-X)
	Title            string     // Document title
	Type             string     // Document type (adr, plan, retrospective, other)
	Status           string     // Document status (draft, published, archived)
	Content          string     // Markdown content
	TrackID          *string    // Optional track attachment
	IterationNumber  *int       // Optional iteration attachment
	CreatedAt        time.Time  // Creation timestamp
	UpdatedAt        time.Time  // Last update timestamp
}
