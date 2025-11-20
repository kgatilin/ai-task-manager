package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
)

// Compile-time check that SQLiteDocumentRepository implements repositories.DocumentRepository
var _ repositories.DocumentRepository = (*SQLiteDocumentRepository)(nil)

// SQLiteDocumentRepository implements repositories.DocumentRepository using SQLite as the backend.
type SQLiteDocumentRepository struct {
	DB *sql.DB
}

// NewSQLiteDocumentRepository creates a new SQLite-backed document repository.
func NewSQLiteDocumentRepository(db *sql.DB) *SQLiteDocumentRepository {
	return &SQLiteDocumentRepository{
		DB: db,
	}
}

// ============================================================================
// Document Operations
// ============================================================================

// SaveDocument persists a new document to storage.
// Returns ErrAlreadyExists if a document with the same ID already exists.
func (r *SQLiteDocumentRepository) SaveDocument(ctx context.Context, doc *entities.DocumentEntity) error {
	// Check if document already exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM documents WHERE id = ?", doc.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check document existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w: document %s already exists", tmerrors.ErrAlreadyExists, doc.ID)
	}

	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(doc.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Insert document with NULL handling for optional attachments
	_, err = r.DB.ExecContext(
		ctx,
		`INSERT INTO documents (id, title, type, status, content, track_id, iteration_number,
			created_at, updated_at, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		doc.ID, doc.Title, doc.Type.String(), doc.Status.String(), doc.Content,
		doc.TrackID, doc.IterationNumber, doc.CreatedAt, doc.UpdatedAt, string(metadataJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}

// FindDocumentByID retrieves a document by its ID.
// Returns ErrNotFound if the document doesn't exist.
func (r *SQLiteDocumentRepository) FindDocumentByID(ctx context.Context, id string) (*entities.DocumentEntity, error) {
	row := r.DB.QueryRowContext(
		ctx,
		`SELECT id, title, type, status, content, track_id, iteration_number,
			created_at, updated_at, metadata FROM documents WHERE id = ?`,
		id,
	)

	doc, err := r.scanDocument(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: document %s not found", tmerrors.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to scan document: %w", err)
	}

	return doc, nil
}

// FindAllDocuments returns all documents in storage.
// Returns empty slice if no documents exist.
func (r *SQLiteDocumentRepository) FindAllDocuments(ctx context.Context) ([]*entities.DocumentEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT id, title, type, status, content, track_id, iteration_number,
			created_at, updated_at, metadata FROM documents ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	documents := []*entities.DocumentEntity{}
	for rows.Next() {
		doc, err := r.scanDocumentFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}

	return documents, nil
}

// FindDocumentsByTrack returns all documents attached to a specific track.
// Returns empty slice if no documents are attached to the track.
func (r *SQLiteDocumentRepository) FindDocumentsByTrack(ctx context.Context, trackID string) ([]*entities.DocumentEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT id, title, type, status, content, track_id, iteration_number,
			created_at, updated_at, metadata FROM documents
		WHERE track_id = ? ORDER BY created_at DESC`,
		trackID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents by track: %w", err)
	}
	defer rows.Close()

	documents := []*entities.DocumentEntity{}
	for rows.Next() {
		doc, err := r.scanDocumentFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}

	return documents, nil
}

// FindDocumentsByIteration returns all documents attached to a specific iteration.
// Returns empty slice if no documents are attached to the iteration.
func (r *SQLiteDocumentRepository) FindDocumentsByIteration(ctx context.Context, iterationNumber int) ([]*entities.DocumentEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT id, title, type, status, content, track_id, iteration_number,
			created_at, updated_at, metadata FROM documents
		WHERE iteration_number = ? ORDER BY created_at DESC`,
		iterationNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents by iteration: %w", err)
	}
	defer rows.Close()

	documents := []*entities.DocumentEntity{}
	for rows.Next() {
		doc, err := r.scanDocumentFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}

	return documents, nil
}

// FindDocumentsByType returns all documents of a specific type.
// Returns empty slice if no documents of that type exist.
func (r *SQLiteDocumentRepository) FindDocumentsByType(ctx context.Context, docType entities.DocumentType) ([]*entities.DocumentEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT id, title, type, status, content, track_id, iteration_number,
			created_at, updated_at, metadata FROM documents
		WHERE type = ? ORDER BY created_at DESC`,
		docType.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents by type: %w", err)
	}
	defer rows.Close()

	documents := []*entities.DocumentEntity{}
	for rows.Next() {
		doc, err := r.scanDocumentFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}

	return documents, nil
}

// UpdateDocument updates an existing document.
// Returns ErrNotFound if the document doesn't exist.
func (r *SQLiteDocumentRepository) UpdateDocument(ctx context.Context, doc *entities.DocumentEntity) error {
	// Convert metadata to JSON
	metadataJSON, err := json.Marshal(doc.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := r.DB.ExecContext(
		ctx,
		`UPDATE documents SET title = ?, type = ?, status = ?, content = ?,
			track_id = ?, iteration_number = ?, updated_at = ?, metadata = ?
		WHERE id = ?`,
		doc.Title, doc.Type.String(), doc.Status.String(), doc.Content,
		doc.TrackID, doc.IterationNumber, doc.UpdatedAt, string(metadataJSON), doc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: document %s not found", tmerrors.ErrNotFound, doc.ID)
	}

	return nil
}

// DeleteDocument removes a document from storage.
// Returns ErrNotFound if the document doesn't exist.
func (r *SQLiteDocumentRepository) DeleteDocument(ctx context.Context, id string) error {
	result, err := r.DB.ExecContext(ctx, "DELETE FROM documents WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: document %s not found", tmerrors.ErrNotFound, id)
	}

	return nil
}

// ============================================================================
// Helper Methods
// ============================================================================

// scanDocument scans a single row into a DocumentEntity using a Row interface
func (r *SQLiteDocumentRepository) scanDocument(row *sql.Row) (*entities.DocumentEntity, error) {
	var doc entities.DocumentEntity
	var metadataJSON string

	err := row.Scan(
		&doc.ID, &doc.Title, &doc.Type, &doc.Status, &doc.Content,
		&doc.TrackID, &doc.IterationNumber, &doc.CreatedAt, &doc.UpdatedAt, &metadataJSON,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal metadata JSON
	if metadataJSON != "" {
		err = json.Unmarshal([]byte(metadataJSON), &doc.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	} else {
		doc.Metadata = make(map[string]interface{})
	}

	return &doc, nil
}

// scanDocumentFromRows scans a document from a Rows iterator
func (r *SQLiteDocumentRepository) scanDocumentFromRows(rows *sql.Rows) (*entities.DocumentEntity, error) {
	var doc entities.DocumentEntity
	var metadataJSON string

	err := rows.Scan(
		&doc.ID, &doc.Title, &doc.Type, &doc.Status, &doc.Content,
		&doc.TrackID, &doc.IterationNumber, &doc.CreatedAt, &doc.UpdatedAt, &metadataJSON,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal metadata JSON
	if metadataJSON != "" {
		err = json.Unmarshal([]byte(metadataJSON), &doc.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	} else {
		doc.Metadata = make(map[string]interface{})
	}

	return &doc, nil
}
