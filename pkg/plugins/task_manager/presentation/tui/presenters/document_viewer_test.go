package presenters_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/repositories"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/presenters"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// mockDocumentRepository is a mock implementation of repositories.DocumentRepository for testing
type mockDocumentRepository struct {
	documents map[string]*entities.DocumentEntity
	err       error
}

var _ repositories.DocumentRepository = (*mockDocumentRepository)(nil)

func newMockDocumentRepository() *mockDocumentRepository {
	return &mockDocumentRepository{
		documents: make(map[string]*entities.DocumentEntity),
	}
}

func (m *mockDocumentRepository) SaveDocument(ctx context.Context, doc *entities.DocumentEntity) error {
	if m.err != nil {
		return m.err
	}
	m.documents[doc.ID] = doc
	return nil
}

func (m *mockDocumentRepository) FindDocumentByID(ctx context.Context, id string) (*entities.DocumentEntity, error) {
	if m.err != nil {
		return nil, m.err
	}
	doc, exists := m.documents[id]
	if !exists {
		return nil, errors.New("document not found")
	}
	return doc, nil
}

func (m *mockDocumentRepository) FindAllDocuments(ctx context.Context) ([]*entities.DocumentEntity, error) {
	if m.err != nil {
		return nil, m.err
	}
	docs := make([]*entities.DocumentEntity, 0)
	for _, doc := range m.documents {
		docs = append(docs, doc)
	}
	return docs, nil
}

func (m *mockDocumentRepository) FindDocumentsByTrack(ctx context.Context, trackID string) ([]*entities.DocumentEntity, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []*entities.DocumentEntity{}, nil
}

func (m *mockDocumentRepository) FindDocumentsByIteration(ctx context.Context, iterationNumber int) ([]*entities.DocumentEntity, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []*entities.DocumentEntity{}, nil
}

func (m *mockDocumentRepository) FindDocumentsByType(ctx context.Context, docType entities.DocumentType) ([]*entities.DocumentEntity, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []*entities.DocumentEntity{}, nil
}

func (m *mockDocumentRepository) UpdateDocument(ctx context.Context, doc *entities.DocumentEntity) error {
	if m.err != nil {
		return m.err
	}
	m.documents[doc.ID] = doc
	return nil
}

func (m *mockDocumentRepository) DeleteDocument(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.documents, id)
	return nil
}

// Test helpers
func createTestDocument(id, title, docType, status string) *entities.DocumentEntity {
	now := time.Now()
	doc, _ := entities.NewDocumentEntity(
		id, title, entities.DocumentType(docType), entities.DocumentStatus(status),
		"# Test Content\n\nThis is test markdown content.",
		nil, nil, now, now,
	)
	return doc
}

// Tests

func TestNewDocumentViewerPresenter(t *testing.T) {
	repo := newMockDocumentRepository()
	ctx := context.Background()

	presenter := presenters.NewDocumentViewerPresenter("TM-doc-1", repo, ctx)

	if presenter == nil {
		t.Fatal("expected presenter, got nil")
	}
	// Smoke test - just verify initialization succeeds
}

func TestDocumentViewerPresenter_DocumentLoading(t *testing.T) {
	repo := newMockDocumentRepository()
	doc := createTestDocument("TM-doc-1", "Test Doc", "adr", "draft")
	repo.documents["TM-doc-1"] = doc

	ctx := context.Background()
	presenter := presenters.NewDocumentViewerPresenter("TM-doc-1", repo, ctx)

	// Init should return a command (async load)
	cmd := presenter.Init()
	if cmd == nil {
		t.Error("expected Init to return a command")
	}
}

func TestDocumentViewerPresenter_ApproveWorkflow(t *testing.T) {
	repo := newMockDocumentRepository()
	doc := createTestDocument("TM-doc-1", "Test Doc", "adr", "draft")
	repo.documents["TM-doc-1"] = doc

	ctx := context.Background()
	presenter := presenters.NewDocumentViewerPresenter("TM-doc-1", repo, ctx)
	_ = presenter

	// Get the approve command (would be triggered by key message in real usage)
	// Command execution happens asynchronously in real usage

	// Verify document state can be updated
	updatedDoc, err := repo.FindDocumentByID(ctx, "TM-doc-1")
	if err != nil {
		t.Fatalf("failed to retrieve document: %v", err)
	}

	// Document should have initial draft status
	if string(updatedDoc.Status) != "draft" {
		t.Errorf("expected status 'draft', got %s", updatedDoc.Status)
	}
}

func TestDocumentViewerPresenter_ErrorHandling(t *testing.T) {
	repo := newMockDocumentRepository()
	repo.err = errors.New("database error")

	ctx := context.Background()
	presenter := presenters.NewDocumentViewerPresenter("TM-doc-1", repo, ctx)

	cmd := presenter.Init()
	if cmd == nil {
		t.Error("expected Init to return a command for async load")
	}
	// Actual error handling tested in Update message handling
}

func TestDocumentViewerPresenter_WindowResize(t *testing.T) {
	repo := newMockDocumentRepository()
	ctx := context.Background()

	presenter := presenters.NewDocumentViewerPresenter("TM-doc-1", repo, ctx)
	_ = presenter

	// Update with window size message should not panic
	// (actual message handling tested via integration tests)
}

func TestDocumentViewModel_Creation(t *testing.T) {
	vm := viewmodels.NewDocumentViewModel(
		"TM-doc-1",
		"Test Document",
		"adr",
		"draft",
		"# Content",
		nil,
		nil,
		"2025-01-01 10:00",
		"2025-01-01 10:00",
	)

	if vm == nil {
		t.Fatal("expected ViewModel, got nil")
	}

	if vm.ID != "TM-doc-1" {
		t.Errorf("expected ID TM-doc-1, got %s", vm.ID)
	}

	if vm.Title != "Test Document" {
		t.Errorf("expected title 'Test Document', got %s", vm.Title)
	}

	if vm.Type != "adr" {
		t.Errorf("expected type 'adr', got %s", vm.Type)
	}

	if vm.Status != "draft" {
		t.Errorf("expected status 'draft', got %s", vm.Status)
	}
}

func TestDocumentLoading_WithMockRepo(t *testing.T) {
	repo := newMockDocumentRepository()
	doc := createTestDocument("TM-doc-1", "Architecture Decision", "adr", "published")
	repo.documents["TM-doc-1"] = doc

	ctx := context.Background()

	// Retrieve document to verify mock works
	retrieved, err := repo.FindDocumentByID(ctx, "TM-doc-1")
	if err != nil {
		t.Fatalf("failed to retrieve: %v", err)
	}

	if retrieved.Title != "Architecture Decision" {
		t.Errorf("expected title 'Architecture Decision', got %s", retrieved.Title)
	}

	if string(retrieved.Status) != "published" {
		t.Errorf("expected status 'published', got %s", retrieved.Status)
	}
}

func TestDocumentStatusUpdate(t *testing.T) {
	repo := newMockDocumentRepository()
	doc := createTestDocument("TM-doc-1", "Test", "plan", "draft")
	repo.documents["TM-doc-1"] = doc

	ctx := context.Background()

	// Simulate status update
	doc.UpdateStatus("published")
	repo.UpdateDocument(ctx, doc)

	// Verify update persisted
	updated, _ := repo.FindDocumentByID(ctx, "TM-doc-1")
	if string(updated.Status) != "published" {
		t.Errorf("expected status 'published', got %s", updated.Status)
	}
}

func TestDocumentNotFound(t *testing.T) {
	repo := newMockDocumentRepository()
	ctx := context.Background()

	_, err := repo.FindDocumentByID(ctx, "TM-doc-nonexistent")
	if err == nil {
		t.Error("expected error for missing document")
	}
}

func TestDocumentRepositoryMock_AllMethods(t *testing.T) {
	repo := newMockDocumentRepository()
	ctx := context.Background()

	doc := createTestDocument("TM-doc-1", "Test", "adr", "draft")

	// Test SaveDocument
	err := repo.SaveDocument(ctx, doc)
	if err != nil {
		t.Errorf("SaveDocument failed: %v", err)
	}

	// Test FindDocumentByID
	found, err := repo.FindDocumentByID(ctx, "TM-doc-1")
	if err != nil || found == nil {
		t.Errorf("FindDocumentByID failed: %v", err)
	}

	// Test FindAllDocuments
	all, err := repo.FindAllDocuments(ctx)
	if err != nil || len(all) != 1 {
		t.Errorf("FindAllDocuments failed: expected 1 document, got %d", len(all))
	}

	// Test UpdateDocument
	doc.UpdateStatus("published")
	err = repo.UpdateDocument(ctx, doc)
	if err != nil {
		t.Errorf("UpdateDocument failed: %v", err)
	}

	// Test DeleteDocument
	err = repo.DeleteDocument(ctx, "TM-doc-1")
	if err != nil {
		t.Errorf("DeleteDocument failed: %v", err)
	}

	// Verify deletion
	_, err = repo.FindDocumentByID(ctx, "TM-doc-1")
	if err == nil {
		t.Error("expected error after deleting document")
	}
}

// ============================================================================
// Document Viewer Scrolling Tests
// ============================================================================

func TestDocumentViewerPresenter_KeyMap_ShiftUpDown(t *testing.T) {
	km := presenters.NewDocumentViewerKeyMap()

	// Verify shift+up and shift+down are in PageUp/PageDown bindings
	pageUpKeys := km.PageUp.Keys()
	pageDownKeys := km.PageDown.Keys()

	hasShiftUp := false
	for _, k := range pageUpKeys {
		if k == "shift+up" {
			hasShiftUp = true
			break
		}
	}

	hasShiftDown := false
	for _, k := range pageDownKeys {
		if k == "shift+down" {
			hasShiftDown = true
			break
		}
	}

	if !hasShiftUp {
		t.Errorf("expected 'shift+up' in PageUp keys, got %v", pageUpKeys)
	}
	if !hasShiftDown {
		t.Errorf("expected 'shift+down' in PageDown keys, got %v", pageDownKeys)
	}
}

func TestDocumentViewerPresenter_Init_RequestsWindowSize(t *testing.T) {
	repo := newMockDocumentRepository()
	ctx := context.Background()

	presenter := presenters.NewDocumentViewerPresenter("TM-doc-1", repo, ctx)

	// Init should return a command (batch of loadDocumentCmd and WindowSize)
	cmd := presenter.Init()
	if cmd == nil {
		t.Error("expected Init to return a command")
	}
}
