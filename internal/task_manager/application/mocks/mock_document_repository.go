package mocks

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// MockDocumentRepository is a mock implementation of DocumentRepository for testing
type MockDocumentRepository struct {
	SaveDocumentFunc             func(ctx context.Context, doc *entities.DocumentEntity) error
	FindDocumentByIDFunc         func(ctx context.Context, id string) (*entities.DocumentEntity, error)
	FindAllDocumentsFunc         func(ctx context.Context) ([]*entities.DocumentEntity, error)
	FindDocumentsByTrackFunc     func(ctx context.Context, trackID string) ([]*entities.DocumentEntity, error)
	FindDocumentsByIterationFunc func(ctx context.Context, iterationNumber int) ([]*entities.DocumentEntity, error)
	FindDocumentsByTypeFunc      func(ctx context.Context, docType entities.DocumentType) ([]*entities.DocumentEntity, error)
	UpdateDocumentFunc           func(ctx context.Context, doc *entities.DocumentEntity) error
	DeleteDocumentFunc           func(ctx context.Context, id string) error
}

// SaveDocument implements DocumentRepository.SaveDocument
func (m *MockDocumentRepository) SaveDocument(ctx context.Context, doc *entities.DocumentEntity) error {
	if m.SaveDocumentFunc != nil {
		return m.SaveDocumentFunc(ctx, doc)
	}
	return nil
}

// FindDocumentByID implements DocumentRepository.FindDocumentByID
func (m *MockDocumentRepository) FindDocumentByID(ctx context.Context, id string) (*entities.DocumentEntity, error) {
	if m.FindDocumentByIDFunc != nil {
		return m.FindDocumentByIDFunc(ctx, id)
	}
	return nil, nil
}

// FindAllDocuments implements DocumentRepository.FindAllDocuments
func (m *MockDocumentRepository) FindAllDocuments(ctx context.Context) ([]*entities.DocumentEntity, error) {
	if m.FindAllDocumentsFunc != nil {
		return m.FindAllDocumentsFunc(ctx)
	}
	return []*entities.DocumentEntity{}, nil
}

// FindDocumentsByTrack implements DocumentRepository.FindDocumentsByTrack
func (m *MockDocumentRepository) FindDocumentsByTrack(ctx context.Context, trackID string) ([]*entities.DocumentEntity, error) {
	if m.FindDocumentsByTrackFunc != nil {
		return m.FindDocumentsByTrackFunc(ctx, trackID)
	}
	return []*entities.DocumentEntity{}, nil
}

// FindDocumentsByIteration implements DocumentRepository.FindDocumentsByIteration
func (m *MockDocumentRepository) FindDocumentsByIteration(ctx context.Context, iterationNumber int) ([]*entities.DocumentEntity, error) {
	if m.FindDocumentsByIterationFunc != nil {
		return m.FindDocumentsByIterationFunc(ctx, iterationNumber)
	}
	return []*entities.DocumentEntity{}, nil
}

// FindDocumentsByType implements DocumentRepository.FindDocumentsByType
func (m *MockDocumentRepository) FindDocumentsByType(ctx context.Context, docType entities.DocumentType) ([]*entities.DocumentEntity, error) {
	if m.FindDocumentsByTypeFunc != nil {
		return m.FindDocumentsByTypeFunc(ctx, docType)
	}
	return []*entities.DocumentEntity{}, nil
}

// UpdateDocument implements DocumentRepository.UpdateDocument
func (m *MockDocumentRepository) UpdateDocument(ctx context.Context, doc *entities.DocumentEntity) error {
	if m.UpdateDocumentFunc != nil {
		return m.UpdateDocumentFunc(ctx, doc)
	}
	return nil
}

// DeleteDocument implements DocumentRepository.DeleteDocument
func (m *MockDocumentRepository) DeleteDocument(ctx context.Context, id string) error {
	if m.DeleteDocumentFunc != nil {
		return m.DeleteDocumentFunc(ctx, id)
	}
	return nil
}
