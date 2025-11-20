package repositories

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// DocumentRepository defines the contract for persistent storage of document entities.
// It handles CRUD operations and filtered queries for document management.
type DocumentRepository interface {
	// SaveDocument persists a new document to storage.
	// Returns ErrAlreadyExists if a document with the same ID already exists.
	SaveDocument(ctx context.Context, doc *entities.DocumentEntity) error

	// FindDocumentByID retrieves a document by its ID.
	// Returns ErrNotFound if the document doesn't exist.
	FindDocumentByID(ctx context.Context, id string) (*entities.DocumentEntity, error)

	// FindAllDocuments returns all documents in storage.
	// Returns empty slice if no documents exist.
	FindAllDocuments(ctx context.Context) ([]*entities.DocumentEntity, error)

	// FindDocumentsByTrack returns all documents attached to a specific track.
	// Returns empty slice if no documents are attached to the track.
	FindDocumentsByTrack(ctx context.Context, trackID string) ([]*entities.DocumentEntity, error)

	// FindDocumentsByIteration returns all documents attached to a specific iteration.
	// Returns empty slice if no documents are attached to the iteration.
	FindDocumentsByIteration(ctx context.Context, iterationNumber int) ([]*entities.DocumentEntity, error)

	// FindDocumentsByType returns all documents of a specific type.
	// Returns empty slice if no documents of that type exist.
	FindDocumentsByType(ctx context.Context, docType entities.DocumentType) ([]*entities.DocumentEntity, error)

	// UpdateDocument updates an existing document.
	// Returns ErrNotFound if the document doesn't exist.
	UpdateDocument(ctx context.Context, doc *entities.DocumentEntity) error

	// DeleteDocument removes a document from storage.
	// Returns ErrNotFound if the document doesn't exist.
	DeleteDocument(ctx context.Context, id string) error
}
