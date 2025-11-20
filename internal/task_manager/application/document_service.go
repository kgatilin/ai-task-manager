package application

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
)

// DocumentApplicationService handles all document-related operations.
// It orchestrates domain validation and repository persistence.
type DocumentApplicationService struct {
	documentRepo  repositories.DocumentRepository
	trackRepo     repositories.TrackRepository
	iterationRepo repositories.IterationRepository
}

// NewDocumentApplicationService creates a new document application service
func NewDocumentApplicationService(
	documentRepo repositories.DocumentRepository,
	trackRepo repositories.TrackRepository,
	iterationRepo repositories.IterationRepository,
) *DocumentApplicationService {
	return &DocumentApplicationService{
		documentRepo:  documentRepo,
		trackRepo:     trackRepo,
		iterationRepo: iterationRepo,
	}
}

// CreateDocument creates a new document with validation
func (s *DocumentApplicationService) CreateDocument(ctx context.Context, input dto.CreateDocumentDTO) (string, error) {
	// Validate title
	if input.Title == "" {
		return "", fmt.Errorf("%w: document title is required", tmerrors.ErrInvalidArgument)
	}
	if len(input.Title) > 200 {
		return "", fmt.Errorf("%w: document title must be 200 characters or less", tmerrors.ErrInvalidArgument)
	}

	// Validate type
	docType, err := entities.NewDocumentType(input.Type)
	if err != nil {
		return "", err
	}

	// Validate status
	docStatus, err := entities.NewDocumentStatus(input.Status)
	if err != nil {
		return "", err
	}

	// Validate content
	if input.Content == "" {
		return "", fmt.Errorf("%w: document content is required", tmerrors.ErrInvalidArgument)
	}

	// Validate XOR for TrackID and IterationNumber
	if input.TrackID != nil && input.IterationNumber != nil {
		return "", fmt.Errorf("%w: document cannot have both TrackID and IterationNumber (choose one or neither)", tmerrors.ErrInvalidArgument)
	}

	// Validate TrackID format if provided
	if input.TrackID != nil && *input.TrackID != "" {
		if !isValidTrackIDFormat(*input.TrackID) {
			return "", fmt.Errorf("%w: invalid track ID format: %s", tmerrors.ErrInvalidArgument, *input.TrackID)
		}
		// Verify track exists
		_, err := s.trackRepo.GetTrack(ctx, *input.TrackID)
		if err != nil {
			return "", fmt.Errorf("track not found: %w", err)
		}
	}

	// Validate IterationNumber if provided
	if input.IterationNumber != nil && *input.IterationNumber < 1 {
		return "", fmt.Errorf("%w: iteration number must be >= 1", tmerrors.ErrInvalidArgument)
	}

	// Generate document ID
	id := generateDocumentID()

	// Create document entity
	now := time.Now().UTC()
	doc, err := entities.NewDocumentEntity(
		id,
		input.Title,
		docType,
		docStatus,
		input.Content,
		input.TrackID,
		input.IterationNumber,
		now,
		now,
	)
	if err != nil {
		return "", err
	}

	// Persist document
	if err := s.documentRepo.SaveDocument(ctx, doc); err != nil {
		return "", err
	}

	return id, nil
}

// UpdateDocument updates an existing document
func (s *DocumentApplicationService) UpdateDocument(ctx context.Context, input dto.UpdateDocumentDTO) error {
	// Load existing document
	doc, err := s.documentRepo.FindDocumentByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Handle detach
	if input.Detach {
		doc.Detach()
		return s.documentRepo.UpdateDocument(ctx, doc)
	}

	// Update content if provided
	if input.Content != nil {
		doc.UpdateContent(*input.Content)
	}

	// Update status if provided
	if input.Status != nil {
		docStatus, err := entities.NewDocumentStatus(*input.Status)
		if err != nil {
			return err
		}
		if err := doc.UpdateStatus(docStatus); err != nil {
			return err
		}
	}

	// Handle attachment changes
	if input.TrackID != nil || input.IterationNumber != nil {
		// Validate XOR constraint
		if input.TrackID != nil && input.IterationNumber != nil {
			return fmt.Errorf("%w: document cannot have both TrackID and IterationNumber (choose one or neither)", tmerrors.ErrInvalidArgument)
		}

		// Attach to track
		if input.TrackID != nil && *input.TrackID != "" {
			if !isValidTrackIDFormat(*input.TrackID) {
				return fmt.Errorf("%w: invalid track ID format: %s", tmerrors.ErrInvalidArgument, *input.TrackID)
			}
			// Verify track exists
			_, err := s.trackRepo.GetTrack(ctx, *input.TrackID)
			if err != nil {
				return fmt.Errorf("track not found: %w", err)
			}
			if err := doc.AttachToTrack(*input.TrackID); err != nil {
				return err
			}
		}

		// Attach to iteration
		if input.IterationNumber != nil && *input.IterationNumber > 0 {
			if err := doc.AttachToIteration(*input.IterationNumber); err != nil {
				return err
			}
		}
	}

	// Persist updates
	if err := s.documentRepo.UpdateDocument(ctx, doc); err != nil {
		return err
	}

	return nil
}

// GetDocument retrieves a document by ID
func (s *DocumentApplicationService) GetDocument(ctx context.Context, id string) (*dto.DocumentViewDTO, error) {
	doc, err := s.documentRepo.FindDocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	return &dto.DocumentViewDTO{
		ID:              doc.ID,
		Title:           doc.Title,
		Type:            doc.Type.String(),
		Status:          doc.Status.String(),
		Content:         doc.Content,
		TrackID:         doc.TrackID,
		IterationNumber: doc.IterationNumber,
		CreatedAt:       doc.CreatedAt,
		UpdatedAt:       doc.UpdatedAt,
	}, nil
}

// ListDocuments lists documents with optional filters
func (s *DocumentApplicationService) ListDocuments(ctx context.Context, trackID *string, iterationNumber *int, docType *string) ([]*dto.DocumentViewDTO, error) {
	var docs []*entities.DocumentEntity
	var err error

	// Apply filters
	if trackID != nil && *trackID != "" {
		docs, err = s.documentRepo.FindDocumentsByTrack(ctx, *trackID)
	} else if iterationNumber != nil && *iterationNumber > 0 {
		docs, err = s.documentRepo.FindDocumentsByIteration(ctx, *iterationNumber)
	} else if docType != nil && *docType != "" {
		docTypeEnum, typeErr := entities.NewDocumentType(*docType)
		if typeErr != nil {
			return nil, typeErr
		}
		docs, err = s.documentRepo.FindDocumentsByType(ctx, docTypeEnum)
	} else {
		docs, err = s.documentRepo.FindAllDocuments(ctx)
	}

	if err != nil {
		return nil, err
	}

	// Convert entities to DTOs
	result := make([]*dto.DocumentViewDTO, len(docs))
	for i, doc := range docs {
		result[i] = &dto.DocumentViewDTO{
			ID:              doc.ID,
			Title:           doc.Title,
			Type:            doc.Type.String(),
			Status:          doc.Status.String(),
			Content:         doc.Content,
			TrackID:         doc.TrackID,
			IterationNumber: doc.IterationNumber,
			CreatedAt:       doc.CreatedAt,
			UpdatedAt:       doc.UpdatedAt,
		}
	}

	return result, nil
}

// AttachDocument attaches a document to a track or iteration
func (s *DocumentApplicationService) AttachDocument(ctx context.Context, id string, trackID *string, iterationNumber *int) error {
	// Load document
	doc, err := s.documentRepo.FindDocumentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Validate XOR
	if trackID != nil && iterationNumber != nil {
		return fmt.Errorf("%w: document cannot have both TrackID and IterationNumber (choose one or neither)", tmerrors.ErrInvalidArgument)
	}

	// Attach to track
	if trackID != nil && *trackID != "" {
		if !isValidTrackIDFormat(*trackID) {
			return fmt.Errorf("%w: invalid track ID format: %s", tmerrors.ErrInvalidArgument, *trackID)
		}
		// Verify track exists
		_, err := s.trackRepo.GetTrack(ctx, *trackID)
		if err != nil {
			return fmt.Errorf("track not found: %w", err)
		}
		if err := doc.AttachToTrack(*trackID); err != nil {
			return err
		}
	}

	// Attach to iteration
	if iterationNumber != nil && *iterationNumber > 0 {
		if err := doc.AttachToIteration(*iterationNumber); err != nil {
			return err
		}
	}

	// Persist
	if err := s.documentRepo.UpdateDocument(ctx, doc); err != nil {
		return err
	}

	return nil
}

// DetachDocument removes a document from track or iteration attachment
func (s *DocumentApplicationService) DetachDocument(ctx context.Context, id string) error {
	// Load document
	doc, err := s.documentRepo.FindDocumentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Detach
	doc.Detach()

	// Persist
	if err := s.documentRepo.UpdateDocument(ctx, doc); err != nil {
		return err
	}

	return nil
}

// DeleteDocument deletes a document
func (s *DocumentApplicationService) DeleteDocument(ctx context.Context, id string) error {
	if err := s.documentRepo.DeleteDocument(ctx, id); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// isValidTrackIDFormat validates track ID format (TM-track-X)
func isValidTrackIDFormat(trackID string) bool {
	pattern := `^[A-Z]+-track-[a-z0-9]+$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(trackID)
}

// generateDocumentID generates a new document ID
func generateDocumentID() string {
	// Use timestamp-based ID generation for simplicity
	// In real implementation, could use UUID or sequence number
	return fmt.Sprintf("TM-doc-%d", time.Now().UnixNano())
}
