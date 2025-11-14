package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/dto"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/mocks"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// setupDocumentTestService creates a test service with mock repositories
func setupDocumentTestService(t *testing.T) (*application.DocumentApplicationService, context.Context, *mocks.MockDocumentRepository, *mocks.MockTrackRepository, *mocks.MockIterationRepository) {
	mockDocRepo := &mocks.MockDocumentRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockIterationRepo := &mocks.MockIterationRepository{}

	service := application.NewDocumentApplicationService(mockDocRepo, mockTrackRepo, mockIterationRepo)
	ctx := context.Background()

	return service, ctx, mockDocRepo, mockTrackRepo, mockIterationRepo
}

// createTestDocument creates a test document entity
func createTestDocument(t *testing.T, id, title, docType, status, content string) *entities.DocumentEntity {
	now := time.Now().UTC()
	dt, err := entities.NewDocumentType(docType)
	if err != nil {
		t.Fatalf("failed to create document type: %v", err)
	}
	ds, err := entities.NewDocumentStatus(status)
	if err != nil {
		t.Fatalf("failed to create document status: %v", err)
	}
	doc, err := entities.NewDocumentEntity(id, title, dt, ds, content, nil, nil, now, now)
	if err != nil {
		t.Fatalf("failed to create test document: %v", err)
	}
	return doc
}

// createTestTrack creates a test track entity
func createTestTrack(t *testing.T, id, roadmapID, title string) *entities.TrackEntity {
	now := time.Now().UTC()
	track, err := entities.NewTrackEntity(id, roadmapID, title, "Test", "not-started", 100, []string{}, now, now)
	if err != nil {
		t.Fatalf("failed to create test track: %v", err)
	}
	return track
}

// TestDocumentService_CreateDocument_Success tests successful document creation
func TestDocumentService_CreateDocument_Success(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.SaveDocumentFunc = func(ctx context.Context, doc *entities.DocumentEntity) error {
		return nil
	}

	input := dto.CreateDocumentDTO{
		Title:           "Test Document",
		Type:            "plan",
		Status:          "draft",
		Content:         "# Test Content",
		TrackID:         nil,
		IterationNumber: nil,
	}

	id, err := service.CreateDocument(ctx, input)
	if err != nil {
		t.Fatalf("CreateDocument() failed: %v", err)
	}

	if id == "" {
		t.Error("document ID should not be empty")
	}
}

// TestDocumentService_CreateDocument_EmptyTitle tests creation with empty title
func TestDocumentService_CreateDocument_EmptyTitle(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	input := dto.CreateDocumentDTO{
		Title:   "",
		Type:    "plan",
		Status:  "draft",
		Content: "# Test Content",
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with empty title")
	}
}

// TestDocumentService_CreateDocument_TitleTooLong tests creation with title exceeding max length
func TestDocumentService_CreateDocument_TitleTooLong(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	// Create title longer than 200 chars
	longTitle := ""
	for i := 0; i < 201; i++ {
		longTitle += "a"
	}

	input := dto.CreateDocumentDTO{
		Title:   longTitle,
		Type:    "plan",
		Status:  "draft",
		Content: "# Test Content",
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with title exceeding 200 chars")
	}
}

// TestDocumentService_CreateDocument_InvalidType tests creation with invalid document type
func TestDocumentService_CreateDocument_InvalidType(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	input := dto.CreateDocumentDTO{
		Title:   "Test Document",
		Type:    "invalid-type",
		Status:  "draft",
		Content: "# Test Content",
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with invalid type")
	}
}

// TestDocumentService_CreateDocument_InvalidStatus tests creation with invalid status
func TestDocumentService_CreateDocument_InvalidStatus(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	input := dto.CreateDocumentDTO{
		Title:   "Test Document",
		Type:    "plan",
		Status:  "invalid-status",
		Content: "# Test Content",
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with invalid status")
	}
}

// TestDocumentService_CreateDocument_EmptyContent tests creation with empty content
func TestDocumentService_CreateDocument_EmptyContent(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	input := dto.CreateDocumentDTO{
		Title:   "Test Document",
		Type:    "plan",
		Status:  "draft",
		Content: "",
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with empty content")
	}
}

// TestDocumentService_CreateDocument_BothTrackAndIteration tests XOR validation
func TestDocumentService_CreateDocument_BothTrackAndIteration(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	trackID := "TM-track-123"
	iterationNum := 1

	input := dto.CreateDocumentDTO{
		Title:           "Test Document",
		Type:            "plan",
		Status:          "draft",
		Content:         "# Test Content",
		TrackID:         &trackID,
		IterationNumber: &iterationNum,
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail when both TrackID and IterationNumber provided")
	}
}

// TestDocumentService_CreateDocument_InvalidTrackID tests creation with invalid track ID format
func TestDocumentService_CreateDocument_InvalidTrackID(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	invalidTrackID := "invalid-track"

	input := dto.CreateDocumentDTO{
		Title:           "Test Document",
		Type:            "plan",
		Status:          "draft",
		Content:         "# Test Content",
		TrackID:         &invalidTrackID,
		IterationNumber: nil,
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with invalid track ID format")
	}
}

// TestDocumentService_CreateDocument_TrackNotFound tests creation with non-existent track
func TestDocumentService_CreateDocument_TrackNotFound(t *testing.T) {
	service, ctx, _, mockTrackRepo, _ := setupDocumentTestService(t)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	trackID := "TM-track-123"

	input := dto.CreateDocumentDTO{
		Title:           "Test Document",
		Type:            "plan",
		Status:          "draft",
		Content:         "# Test Content",
		TrackID:         &trackID,
		IterationNumber: nil,
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with non-existent track")
	}
}

// TestDocumentService_CreateDocument_InvalidIterationNumber tests creation with invalid iteration
func TestDocumentService_CreateDocument_InvalidIterationNumber(t *testing.T) {
	service, ctx, _, _, _ := setupDocumentTestService(t)

	invalidIterationNum := 0

	input := dto.CreateDocumentDTO{
		Title:           "Test Document",
		Type:            "plan",
		Status:          "draft",
		Content:         "# Test Content",
		TrackID:         nil,
		IterationNumber: &invalidIterationNum,
	}

	_, err := service.CreateDocument(ctx, input)
	if err == nil {
		t.Fatal("CreateDocument() should fail with invalid iteration number")
	}
}

// TestDocumentService_CreateDocument_WithTrack tests successful creation with track attachment
func TestDocumentService_CreateDocument_WithTrack(t *testing.T) {
	service, ctx, mockDocRepo, mockTrackRepo, _ := setupDocumentTestService(t)

	track := createTestTrack(t, "TM-track-123", "roadmap-1", "Test Track")

	mockDocRepo.SaveDocumentFunc = func(ctx context.Context, doc *entities.DocumentEntity) error {
		return nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track.ID {
			return track, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	trackID := track.ID

	input := dto.CreateDocumentDTO{
		Title:           "Test Document",
		Type:            "adr",
		Status:          "draft",
		Content:         "# ADR Content",
		TrackID:         &trackID,
		IterationNumber: nil,
	}

	id, err := service.CreateDocument(ctx, input)
	if err != nil {
		t.Fatalf("CreateDocument() failed: %v", err)
	}

	if id == "" {
		t.Error("document ID should not be empty")
	}
}

// TestDocumentService_UpdateDocument_Success tests successful document update
func TestDocumentService_UpdateDocument_Success(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Original content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		if id == doc.ID {
			return doc, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	newContent := "Updated content"
	input := dto.UpdateDocumentDTO{
		ID:      doc.ID,
		Content: &newContent,
		Status:  nil,
	}

	err := service.UpdateDocument(ctx, input)
	if err != nil {
		t.Fatalf("UpdateDocument() failed: %v", err)
	}
}

// TestDocumentService_UpdateDocument_NotFound tests update with non-existent document
func TestDocumentService_UpdateDocument_NotFound(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	input := dto.UpdateDocumentDTO{
		ID: "TM-doc-nonexistent",
	}

	err := service.UpdateDocument(ctx, input)
	if err == nil {
		t.Fatal("UpdateDocument() should fail with non-existent document")
	}
}

// TestDocumentService_UpdateDocument_InvalidStatus tests update with invalid status
func TestDocumentService_UpdateDocument_InvalidStatus(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	invalidStatus := "invalid-status"
	input := dto.UpdateDocumentDTO{
		ID:     doc.ID,
		Status: &invalidStatus,
	}

	err := service.UpdateDocument(ctx, input)
	if err == nil {
		t.Fatal("UpdateDocument() should fail with invalid status")
	}
}

// TestDocumentService_UpdateDocument_Detach tests detaching document
func TestDocumentService_UpdateDocument_Detach(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	trackID := "TM-track-123"
	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")
	doc.TrackID = &trackID

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	input := dto.UpdateDocumentDTO{
		ID:     doc.ID,
		Detach: true,
	}

	err := service.UpdateDocument(ctx, input)
	if err != nil {
		t.Fatalf("UpdateDocument() failed: %v", err)
	}

	if doc.TrackID != nil {
		t.Error("document should be detached from track")
	}
}

// TestDocumentService_GetDocument_Success tests successful document retrieval
func TestDocumentService_GetDocument_Success(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		if id == doc.ID {
			return doc, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	result, err := service.GetDocument(ctx, doc.ID)
	if err != nil {
		t.Fatalf("GetDocument() failed: %v", err)
	}

	if result.ID != doc.ID {
		t.Errorf("result.ID = %q, want %q", result.ID, doc.ID)
	}
	if result.Title != doc.Title {
		t.Errorf("result.Title = %q, want %q", result.Title, doc.Title)
	}
}

// TestDocumentService_GetDocument_NotFound tests retrieval of non-existent document
func TestDocumentService_GetDocument_NotFound(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	_, err := service.GetDocument(ctx, "TM-doc-nonexistent")
	if err == nil {
		t.Fatal("GetDocument() should fail with non-existent document")
	}
}

// TestDocumentService_ListDocuments_NoFilter tests listing all documents
func TestDocumentService_ListDocuments_NoFilter(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	docs := []*entities.DocumentEntity{
		createTestDocument(t, "TM-doc-1", "Doc 1", "plan", "draft", "Content 1"),
		createTestDocument(t, "TM-doc-2", "Doc 2", "adr", "published", "Content 2"),
	}

	mockDocRepo.FindAllDocumentsFunc = func(ctx context.Context) ([]*entities.DocumentEntity, error) {
		return docs, nil
	}

	result, err := service.ListDocuments(ctx, nil, nil, nil)
	if err != nil {
		t.Fatalf("ListDocuments() failed: %v", err)
	}

	if len(result) != len(docs) {
		t.Errorf("result length = %d, want %d", len(result), len(docs))
	}
}

// TestDocumentService_ListDocuments_FilterByTrack tests listing documents filtered by track
func TestDocumentService_ListDocuments_FilterByTrack(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	trackID := "TM-track-123"
	docs := []*entities.DocumentEntity{
		createTestDocument(t, "TM-doc-1", "Doc 1", "plan", "draft", "Content 1"),
	}
	docs[0].TrackID = &trackID

	mockDocRepo.FindDocumentsByTrackFunc = func(ctx context.Context, tid string) ([]*entities.DocumentEntity, error) {
		if tid == trackID {
			return docs, nil
		}
		return []*entities.DocumentEntity{}, nil
	}

	result, err := service.ListDocuments(ctx, &trackID, nil, nil)
	if err != nil {
		t.Fatalf("ListDocuments() failed: %v", err)
	}

	if len(result) != len(docs) {
		t.Errorf("result length = %d, want %d", len(result), len(docs))
	}
}

// TestDocumentService_ListDocuments_FilterByIteration tests listing documents filtered by iteration
func TestDocumentService_ListDocuments_FilterByIteration(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	iterationNum := 1
	docs := []*entities.DocumentEntity{
		createTestDocument(t, "TM-doc-1", "Doc 1", "plan", "draft", "Content 1"),
	}
	docs[0].IterationNumber = &iterationNum

	mockDocRepo.FindDocumentsByIterationFunc = func(ctx context.Context, num int) ([]*entities.DocumentEntity, error) {
		if num == iterationNum {
			return docs, nil
		}
		return []*entities.DocumentEntity{}, nil
	}

	result, err := service.ListDocuments(ctx, nil, &iterationNum, nil)
	if err != nil {
		t.Fatalf("ListDocuments() failed: %v", err)
	}

	if len(result) != len(docs) {
		t.Errorf("result length = %d, want %d", len(result), len(docs))
	}
}

// TestDocumentService_ListDocuments_FilterByType tests listing documents filtered by type
func TestDocumentService_ListDocuments_FilterByType(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	docType := "adr"
	docs := []*entities.DocumentEntity{
		createTestDocument(t, "TM-doc-1", "Doc 1", "adr", "draft", "Content 1"),
		createTestDocument(t, "TM-doc-2", "Doc 2", "adr", "published", "Content 2"),
	}

	mockDocRepo.FindDocumentsByTypeFunc = func(ctx context.Context, dt entities.DocumentType) ([]*entities.DocumentEntity, error) {
		if dt.String() == docType {
			return docs, nil
		}
		return []*entities.DocumentEntity{}, nil
	}

	result, err := service.ListDocuments(ctx, nil, nil, &docType)
	if err != nil {
		t.Fatalf("ListDocuments() failed: %v", err)
	}

	if len(result) != len(docs) {
		t.Errorf("result length = %d, want %d", len(result), len(docs))
	}
}

// TestDocumentService_AttachDocument_ToTrack tests attaching document to a track
func TestDocumentService_AttachDocument_ToTrack(t *testing.T) {
	service, ctx, mockDocRepo, mockTrackRepo, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")
	track := createTestTrack(t, "TM-track-123", "roadmap-1", "Test Track")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	err := service.AttachDocument(ctx, doc.ID, &track.ID, nil)
	if err != nil {
		t.Fatalf("AttachDocument() failed: %v", err)
	}

	if doc.TrackID == nil || *doc.TrackID != track.ID {
		t.Errorf("document should be attached to track %q", track.ID)
	}
}

// TestDocumentService_AttachDocument_ToIteration tests attaching document to an iteration
func TestDocumentService_AttachDocument_ToIteration(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")
	iterationNum := 1

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	err := service.AttachDocument(ctx, doc.ID, nil, &iterationNum)
	if err != nil {
		t.Fatalf("AttachDocument() failed: %v", err)
	}

	if doc.IterationNumber == nil || *doc.IterationNumber != iterationNum {
		t.Errorf("document should be attached to iteration %d", iterationNum)
	}
}

// TestDocumentService_AttachDocument_BothTrackAndIteration tests XOR validation in attach
func TestDocumentService_AttachDocument_BothTrackAndIteration(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	trackID := "TM-track-123"
	iterationNum := 1

	err := service.AttachDocument(ctx, doc.ID, &trackID, &iterationNum)
	if err == nil {
		t.Fatal("AttachDocument() should fail when both trackID and iterationNumber provided")
	}
}

// TestDocumentService_DetachDocument_Success tests successful detachment
func TestDocumentService_DetachDocument_Success(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	trackID := "TM-track-123"
	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")
	doc.TrackID = &trackID

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	err := service.DetachDocument(ctx, doc.ID)
	if err != nil {
		t.Fatalf("DetachDocument() failed: %v", err)
	}

	if doc.TrackID != nil {
		t.Error("document should be detached")
	}
}

// TestDocumentService_DetachDocument_NotFound tests detach with non-existent document
func TestDocumentService_DetachDocument_NotFound(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	err := service.DetachDocument(ctx, "TM-doc-nonexistent")
	if err == nil {
		t.Fatal("DetachDocument() should fail with non-existent document")
	}
}

// TestDocumentService_DeleteDocument_Success tests successful deletion
func TestDocumentService_DeleteDocument_Success(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.DeleteDocumentFunc = func(ctx context.Context, id string) error {
		return nil
	}

	err := service.DeleteDocument(ctx, "TM-doc-123")
	if err != nil {
		t.Fatalf("DeleteDocument() failed: %v", err)
	}
}

// TestDocumentService_DeleteDocument_NotFound tests deletion of non-existent document
func TestDocumentService_DeleteDocument_NotFound(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.DeleteDocumentFunc = func(ctx context.Context, id string) error {
		return pluginsdk.ErrNotFound
	}

	err := service.DeleteDocument(ctx, "TM-doc-nonexistent")
	if err == nil {
		t.Fatal("DeleteDocument() should fail with non-existent document")
	}
}

// TestDocumentService_UpdateDocument_StatusUpdate tests updating only status
func TestDocumentService_UpdateDocument_StatusUpdate(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	newStatus := "published"
	input := dto.UpdateDocumentDTO{
		ID:     doc.ID,
		Status: &newStatus,
	}

	err := service.UpdateDocument(ctx, input)
	if err != nil {
		t.Fatalf("UpdateDocument() failed: %v", err)
	}

	if doc.Status.String() != newStatus {
		t.Errorf("document status = %q, want %q", doc.Status.String(), newStatus)
	}
}

// TestDocumentService_ListDocuments_EmptyResult tests listing with empty results
func TestDocumentService_ListDocuments_EmptyResult(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.FindAllDocumentsFunc = func(ctx context.Context) ([]*entities.DocumentEntity, error) {
		return []*entities.DocumentEntity{}, nil
	}

	result, err := service.ListDocuments(ctx, nil, nil, nil)
	if err != nil {
		t.Fatalf("ListDocuments() failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("result length = %d, want 0", len(result))
	}
}

// TestDocumentService_UpdateDocument_AttachToTrack tests attaching document via update
func TestDocumentService_UpdateDocument_AttachToTrack(t *testing.T) {
	service, ctx, mockDocRepo, mockTrackRepo, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")
	track := createTestTrack(t, "TM-track-123", "roadmap-1", "Test Track")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	trackID := track.ID
	input := dto.UpdateDocumentDTO{
		ID:      doc.ID,
		TrackID: &trackID,
	}

	err := service.UpdateDocument(ctx, input)
	if err != nil {
		t.Fatalf("UpdateDocument() failed: %v", err)
	}

	if doc.TrackID == nil || *doc.TrackID != track.ID {
		t.Errorf("document should be attached to track %q", track.ID)
	}
}

// TestDocumentService_UpdateDocument_AttachToIteration tests attaching document via update
func TestDocumentService_UpdateDocument_AttachToIteration(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")
	iterationNum := 2

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	input := dto.UpdateDocumentDTO{
		ID:              doc.ID,
		IterationNumber: &iterationNum,
	}

	err := service.UpdateDocument(ctx, input)
	if err != nil {
		t.Fatalf("UpdateDocument() failed: %v", err)
	}

	if doc.IterationNumber == nil || *doc.IterationNumber != iterationNum {
		t.Errorf("document should be attached to iteration %d", iterationNum)
	}
}

// TestDocumentService_UpdateDocument_BothTrackAndIterationUpdate tests XOR validation in update
func TestDocumentService_UpdateDocument_BothTrackAndIterationUpdate(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	trackID := "TM-track-123"
	iterationNum := 1

	input := dto.UpdateDocumentDTO{
		ID:              doc.ID,
		TrackID:         &trackID,
		IterationNumber: &iterationNum,
	}

	err := service.UpdateDocument(ctx, input)
	if err == nil {
		t.Fatal("UpdateDocument() should fail when both trackID and iterationNumber provided")
	}
}

// TestDocumentService_UpdateDocument_ContentUpdate tests updating only content
func TestDocumentService_UpdateDocument_ContentUpdate(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Original content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockDocRepo.UpdateDocumentFunc = func(ctx context.Context, d *entities.DocumentEntity) error {
		return nil
	}

	newContent := "Updated content"
	input := dto.UpdateDocumentDTO{
		ID:      doc.ID,
		Content: &newContent,
	}

	err := service.UpdateDocument(ctx, input)
	if err != nil {
		t.Fatalf("UpdateDocument() failed: %v", err)
	}

	if doc.Content != newContent {
		t.Errorf("document content = %q, want %q", doc.Content, newContent)
	}
}

// TestDocumentService_AttachDocument_NotFound tests attaching to non-existent document
func TestDocumentService_AttachDocument_NotFound(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	trackID := "TM-track-123"

	err := service.AttachDocument(ctx, "TM-doc-nonexistent", &trackID, nil)
	if err == nil {
		t.Fatal("AttachDocument() should fail with non-existent document")
	}
}

// TestDocumentService_AttachDocument_InvalidTrackID tests attaching with invalid track ID
func TestDocumentService_AttachDocument_InvalidTrackID(t *testing.T) {
	service, ctx, mockDocRepo, _, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	invalidTrackID := "invalid-track"

	err := service.AttachDocument(ctx, doc.ID, &invalidTrackID, nil)
	if err == nil {
		t.Fatal("AttachDocument() should fail with invalid track ID format")
	}
}

// TestDocumentService_AttachDocument_TrackNotFound tests attaching to non-existent track
func TestDocumentService_AttachDocument_TrackNotFound(t *testing.T) {
	service, ctx, mockDocRepo, mockTrackRepo, _ := setupDocumentTestService(t)

	doc := createTestDocument(t, "TM-doc-123", "Test Doc", "plan", "draft", "Content")

	mockDocRepo.FindDocumentByIDFunc = func(ctx context.Context, id string) (*entities.DocumentEntity, error) {
		return doc, nil
	}

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	trackID := "TM-track-nonexistent"

	err := service.AttachDocument(ctx, doc.ID, &trackID, nil)
	if err == nil {
		t.Fatal("AttachDocument() should fail with non-existent track")
	}
}
