package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/persistence"
)

// TestSQLiteDocumentRepository_SaveDocument tests saving a new document
func TestSQLiteDocumentRepository_SaveDocument(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Create a document
	doc, err := entities.NewDocumentEntity(
		"TM-doc-test001",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"This is test content",
		nil, // No track attachment
		nil, // No iteration attachment
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document entity: %v", err)
	}

	// Save document
	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Verify it was saved
	found, err := repo.FindDocumentByID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("failed to find document: %v", err)
	}
	if found.ID != doc.ID {
		t.Errorf("document ID mismatch: got %s, want %s", found.ID, doc.ID)
	}
	if found.Title != doc.Title {
		t.Errorf("document title mismatch: got %s, want %s", found.Title, doc.Title)
	}
}

// TestSQLiteDocumentRepository_SaveDocument_WithTrackAttachment tests saving a document with track attachment
func TestSQLiteDocumentRepository_SaveDocument_WithTrackAttachment(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Setup: Create track first
	createTrack(t, db, "TM-track-001")

	trackID := "TM-track-001"
	doc, err := entities.NewDocumentEntity(
		"TM-doc-test002",
		"Plan Document",
		entities.DocumentTypePlan,
		entities.DocumentStatusPublished,
		"Plan content",
		&trackID,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document entity: %v", err)
	}

	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Verify track attachment
	found, err := repo.FindDocumentByID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("failed to find document: %v", err)
	}
	if found.TrackID == nil || *found.TrackID != trackID {
		t.Errorf("document track attachment failed: got %v, want %s", found.TrackID, trackID)
	}
	if found.IterationNumber != nil {
		t.Errorf("iteration number should be nil, got %v", found.IterationNumber)
	}
}

// TestSQLiteDocumentRepository_SaveDocument_WithIterationAttachment tests saving a document with iteration attachment
func TestSQLiteDocumentRepository_SaveDocument_WithIterationAttachment(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Setup: Create iteration first
	createIteration(t, db, 1)

	iterNum := 1
	doc, err := entities.NewDocumentEntity(
		"TM-doc-test003",
		"Retrospective",
		entities.DocumentTypeRetrospective,
		entities.DocumentStatusPublished,
		"Retrospective content",
		nil,
		&iterNum,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document entity: %v", err)
	}

	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Verify iteration attachment
	found, err := repo.FindDocumentByID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("failed to find document: %v", err)
	}
	if found.IterationNumber == nil || *found.IterationNumber != iterNum {
		t.Errorf("document iteration attachment failed: got %v, want %d", found.IterationNumber, iterNum)
	}
	if found.TrackID != nil {
		t.Errorf("track ID should be nil, got %v", found.TrackID)
	}
}

// TestSQLiteDocumentRepository_SaveDocument_XORConstraint tests XOR constraint (neither or one attachment)
func TestSQLiteDocumentRepository_SaveDocument_XORConstraint(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	// Setup
	createTrack(t, db, "TM-track-002")
	createIteration(t, db, 2)

	trackID := "TM-track-002"
	iterNum := 2

	// Try to create document with both attachments (should fail in entity constructor)
	_, err := entities.NewDocumentEntity(
		"TM-doc-test004",
		"Invalid Document",
		entities.DocumentTypeOther,
		entities.DocumentStatusDraft,
		"content",
		&trackID,
		&iterNum,
		time.Now(),
		time.Now(),
	)
	if err == nil {
		t.Fatalf("expected error when creating document with both attachments, got nil")
	}
}

// TestSQLiteDocumentRepository_SaveDocument_AlreadyExists tests duplicate document handling
func TestSQLiteDocumentRepository_SaveDocument_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	doc, err := entities.NewDocumentEntity(
		"TM-doc-test005",
		"Test Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document entity: %v", err)
	}

	// Save first time
	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Try to save again
	err = repo.SaveDocument(ctx, doc)
	if err == nil {
		t.Fatalf("expected ErrAlreadyExists, got nil")
	}
	if !errors.Is(err, tmerrors.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

// TestSQLiteDocumentRepository_FindDocumentByID tests finding document by ID
func TestSQLiteDocumentRepository_FindDocumentByID(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Save a document
	doc, err := entities.NewDocumentEntity(
		"TM-doc-test006",
		"Find Me",
		entities.DocumentTypePlan,
		entities.DocumentStatusPublished,
		"content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}
	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Find it
	found, err := repo.FindDocumentByID(ctx, "TM-doc-test006")
	if err != nil {
		t.Fatalf("failed to find document: %v", err)
	}
	if found.Title != "Find Me" {
		t.Errorf("wrong document found: got %s, want Find Me", found.Title)
	}
}

// TestSQLiteDocumentRepository_FindDocumentByID_NotFound tests not found case
func TestSQLiteDocumentRepository_FindDocumentByID_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	_, err := repo.FindDocumentByID(ctx, "TM-doc-nonexistent")
	if err == nil {
		t.Fatalf("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, tmerrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// TestSQLiteDocumentRepository_FindAllDocuments tests finding all documents
func TestSQLiteDocumentRepository_FindAllDocuments(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Save multiple documents
	for i := 1; i <= 3; i++ {
		suffix := "001"
		if i == 2 {
			suffix = "002"
		} else if i == 3 {
			suffix = "003"
		}
		doc, err := entities.NewDocumentEntity(
			"TM-doc-all"+suffix,
			"Document "+suffix,
			entities.DocumentTypeADR,
			entities.DocumentStatusDraft,
			"content",
			nil,
			nil,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			t.Fatalf("failed to create document: %v", err)
		}
		err = repo.SaveDocument(ctx, doc)
		if err != nil {
			t.Fatalf("failed to save document: %v", err)
		}
	}

	// Find all
	all, err := repo.FindAllDocuments(ctx)
	if err != nil {
		t.Fatalf("failed to find all documents: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 documents, got %d", len(all))
	}
}

// TestSQLiteDocumentRepository_FindAllDocuments_Empty tests empty result
func TestSQLiteDocumentRepository_FindAllDocuments_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	all, err := repo.FindAllDocuments(ctx)
	if err != nil {
		t.Fatalf("failed to find all documents: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("expected 0 documents, got %d", len(all))
	}
}

// TestSQLiteDocumentRepository_FindDocumentsByTrack tests finding documents by track
func TestSQLiteDocumentRepository_FindDocumentsByTrack(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Setup
	createTrack(t, db, "TM-track-003")
	createTrack(t, db, "TM-track-004")

	// Save documents for track 003
	for i := 1; i <= 2; i++ {
		trackID := "TM-track-003"
		suffix := "001"
		if i == 2 {
			suffix = "002"
		}
		doc, err := entities.NewDocumentEntity(
			"TM-doc-track"+suffix,
			"Track Document",
			entities.DocumentTypePlan,
			entities.DocumentStatusPublished,
			"content",
			&trackID,
			nil,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			t.Fatalf("failed to create document: %v", err)
		}
		err = repo.SaveDocument(ctx, doc)
		if err != nil {
			t.Fatalf("failed to save document: %v", err)
		}
	}

	// Save document for track 004
	trackID := "TM-track-004"
	doc, err := entities.NewDocumentEntity(
		"TM-doc-othertrack",
		"Other Track Document",
		entities.DocumentTypePlan,
		entities.DocumentStatusPublished,
		"content",
		&trackID,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}
	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Find documents for track 003
	found, err := repo.FindDocumentsByTrack(ctx, "TM-track-003")
	if err != nil {
		t.Fatalf("failed to find documents by track: %v", err)
	}
	if len(found) != 2 {
		t.Errorf("expected 2 documents for track, got %d", len(found))
	}

	// Verify all are for correct track
	for _, doc := range found {
		if doc.TrackID == nil || *doc.TrackID != "TM-track-003" {
			t.Errorf("document has wrong track: got %v, want TM-track-003", doc.TrackID)
		}
	}
}

// TestSQLiteDocumentRepository_FindDocumentsByIteration tests finding documents by iteration
func TestSQLiteDocumentRepository_FindDocumentsByIteration(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Setup
	createIteration(t, db, 3)
	createIteration(t, db, 4)

	// Save documents for iteration 3
	for i := 1; i <= 2; i++ {
		iterNum := 3
		suffix := "001"
		if i == 2 {
			suffix = "002"
		}
		doc, err := entities.NewDocumentEntity(
			"TM-doc-iter"+suffix,
			"Iteration Document",
			entities.DocumentTypeRetrospective,
			entities.DocumentStatusPublished,
			"content",
			nil,
			&iterNum,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			t.Fatalf("failed to create document: %v", err)
		}
		err = repo.SaveDocument(ctx, doc)
		if err != nil {
			t.Fatalf("failed to save document: %v", err)
		}
	}

	// Find documents for iteration 3
	found, err := repo.FindDocumentsByIteration(ctx, 3)
	if err != nil {
		t.Fatalf("failed to find documents by iteration: %v", err)
	}
	if len(found) != 2 {
		t.Errorf("expected 2 documents for iteration, got %d", len(found))
	}

	// Verify all are for correct iteration
	for _, doc := range found {
		if doc.IterationNumber == nil || *doc.IterationNumber != 3 {
			t.Errorf("document has wrong iteration: got %v, want 3", doc.IterationNumber)
		}
	}
}

// TestSQLiteDocumentRepository_FindDocumentsByType tests finding documents by type
func TestSQLiteDocumentRepository_FindDocumentsByType(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Save documents of different types
	doc1, err := entities.NewDocumentEntity(
		"TM-doc-adr001",
		"ADR Document",
		entities.DocumentTypeADR,
		entities.DocumentStatusPublished,
		"content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}
	err = repo.SaveDocument(ctx, doc1)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	doc2, err := entities.NewDocumentEntity(
		"TM-doc-plan001",
		"Plan Document",
		entities.DocumentTypePlan,
		entities.DocumentStatusPublished,
		"content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}
	err = repo.SaveDocument(ctx, doc2)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	doc3, err := entities.NewDocumentEntity(
		"TM-doc-adr002",
		"Another ADR",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}
	err = repo.SaveDocument(ctx, doc3)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Find all ADR documents
	found, err := repo.FindDocumentsByType(ctx, entities.DocumentTypeADR)
	if err != nil {
		t.Fatalf("failed to find documents by type: %v", err)
	}
	if len(found) != 2 {
		t.Errorf("expected 2 ADR documents, got %d", len(found))
	}

	// Verify all are ADR type
	for _, doc := range found {
		if doc.Type != entities.DocumentTypeADR {
			t.Errorf("document has wrong type: got %v, want adr", doc.Type)
		}
	}
}

// TestSQLiteDocumentRepository_UpdateDocument tests updating document
func TestSQLiteDocumentRepository_UpdateDocument(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Create and save document
	doc, err := entities.NewDocumentEntity(
		"TM-doc-update",
		"Original Title",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"original content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}
	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Update document
	doc.Title = "Updated Title"
	doc.Content = "updated content"
	doc.Status = entities.DocumentStatusPublished
	doc.UpdatedAt = time.Now()

	updateErr := repo.UpdateDocument(ctx, doc)
	if updateErr != nil {
		t.Fatalf("failed to update document: %v", updateErr)
	}

	// Verify update
	found, _ := repo.FindDocumentByID(ctx, doc.ID)
	if found.Title != "Updated Title" {
		t.Errorf("title not updated: got %s, want Updated Title", found.Title)
	}
	if found.Content != "updated content" {
		t.Errorf("content not updated: got %s, want updated content", found.Content)
	}
	if found.Status != entities.DocumentStatusPublished {
		t.Errorf("status not updated: got %v, want published", found.Status)
	}
}

// TestSQLiteDocumentRepository_UpdateDocument_NotFound tests updating non-existent document
func TestSQLiteDocumentRepository_UpdateDocument_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	doc, err := entities.NewDocumentEntity(
		"TM-doc-nonexistent2",
		"Title",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	updateErr := repo.UpdateDocument(ctx, doc)
	if updateErr == nil {
		t.Fatalf("expected ErrNotFound, got nil")
	}
	if !errors.Is(updateErr, tmerrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", updateErr)
	}
}

// TestSQLiteDocumentRepository_DeleteDocument tests deleting document
func TestSQLiteDocumentRepository_DeleteDocument(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	// Create and save document
	doc, err := entities.NewDocumentEntity(
		"TM-doc-delete",
		"To Delete",
		entities.DocumentTypeADR,
		entities.DocumentStatusDraft,
		"content",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}
	err = repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Fatalf("failed to save document: %v", err)
	}

	// Delete it
	err = repo.DeleteDocument(ctx, doc.ID)
	if err != nil {
		t.Fatalf("failed to delete document: %v", err)
	}

	// Verify it's gone
	_, err = repo.FindDocumentByID(ctx, doc.ID)
	if err == nil {
		t.Fatalf("expected document to be deleted")
	}
}

// TestSQLiteDocumentRepository_DeleteDocument_NotFound tests deleting non-existent document
func TestSQLiteDocumentRepository_DeleteDocument_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	db := setupTestDB(t, tmpDir)
	defer db.Close()

	repo := persistence.NewSQLiteDocumentRepository(db)
	ctx := context.Background()

	err := repo.DeleteDocument(ctx, "TM-doc-nonexistent-3")
	if err == nil {
		t.Fatalf("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, tmerrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// setupTestDB creates a test database with schema initialized
func setupTestDB(t *testing.T, tmpDir string) *sql.DB {
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Initialize schema
	err = persistence.InitSchema(db)
	if err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return db
}

// createTrack is a helper to create a test track
func createTrack(t *testing.T, db *sql.DB, trackID string) {
	ctx := context.Background()

	// Create roadmap first
	roadmapID := "roadmap-default"
	_, err := db.ExecContext(
		ctx,
		"INSERT OR IGNORE INTO roadmaps (id, vision, success_criteria, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		roadmapID, "Test Vision", "Test Success", time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap: %v", err)
	}

	// Create track
	_, err = db.ExecContext(
		ctx,
		"INSERT INTO tracks (id, roadmap_id, title, description, status, rank, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		trackID, roadmapID, "Track", "Track description", "not-started", 500, time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}
}

// createIteration is a helper to create a test iteration
func createIteration(t *testing.T, db *sql.DB, iterNum int) {
	ctx := context.Background()
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO iterations (number, name, goal, status, rank, deliverable, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		iterNum, "Iteration "+string(rune(iterNum+48)), "Goal", "planned", 500, "Deliverable", time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create iteration: %v", err)
	}
}
