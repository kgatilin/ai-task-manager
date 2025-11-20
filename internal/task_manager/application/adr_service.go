package application

import (
	"context"
	"fmt"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
)

// ADRApplicationService handles all ADR operations
type ADRApplicationService struct {
	adrRepo           repositories.ADRRepository
	trackRepo         repositories.TrackRepository
	aggregateRepo     repositories.AggregateRepository
	validationService *services.ValidationService
}

// NewADRApplicationService creates a new ADR service
func NewADRApplicationService(
	adrRepo repositories.ADRRepository,
	trackRepo repositories.TrackRepository,
	aggregateRepo repositories.AggregateRepository,
	validationService *services.ValidationService,
) *ADRApplicationService {
	return &ADRApplicationService{
		adrRepo:           adrRepo,
		trackRepo:         trackRepo,
		aggregateRepo:     aggregateRepo,
		validationService: validationService,
	}
}

// CreateADR creates a new ADR
func (s *ADRApplicationService) CreateADR(ctx context.Context, input dto.CreateADRDTO) (*entities.ADREntity, error) {
	// Generate ADR ID
	projectCode := s.aggregateRepo.GetProjectCode(ctx)
	nextNum, err := s.aggregateRepo.GetNextSequenceNumber(ctx, "adr")
	if err != nil {
		return nil, fmt.Errorf("failed to generate ADR ID: %w", err)
	}
	id := fmt.Sprintf("%s-adr-%d", projectCode, nextNum)

	// Validate ADR ID format
	if err := s.validationService.ValidateNonEmpty("ADR ID", id); err != nil {
		return nil, err
	}

	// Validate title
	if err := s.validationService.ValidateNonEmpty("title", input.Title); err != nil {
		return nil, err
	}

	// Validate required fields
	if err := s.validationService.ValidateNonEmpty("context", input.Context); err != nil {
		return nil, err
	}
	if err := s.validationService.ValidateNonEmpty("decision", input.Decision); err != nil {
		return nil, err
	}
	if err := s.validationService.ValidateNonEmpty("consequences", input.Consequences); err != nil {
		return nil, err
	}

	// Verify track exists
	_, err = s.trackRepo.GetTrack(ctx, input.TrackID)
	if err != nil {
		return nil, fmt.Errorf("track not found: %w", err)
	}

	// Set default status if not provided
	status := input.Status
	if status == "" {
		status = string(entities.ADRStatusProposed)
	}

	// Validate status
	if !entities.IsValidADRStatus(status) {
		return nil, fmt.Errorf("%w: invalid ADR status: %s", tmerrors.ErrInvalidArgument, status)
	}

	now := time.Now().UTC()

	// Create ADR entity
	adr, err := entities.NewADREntity(
		id,
		input.TrackID,
		input.Title,
		status,
		input.Context,
		input.Decision,
		input.Consequences,
		input.Alternatives,
		now,
		now,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ADR entity: %w", err)
	}

	// Persist ADR
	if err := s.adrRepo.SaveADR(ctx, adr); err != nil {
		return nil, fmt.Errorf("failed to save ADR: %w", err)
	}

	return adr, nil
}

// UpdateADR updates an existing ADR
func (s *ADRApplicationService) UpdateADR(ctx context.Context, input dto.UpdateADRDTO) (*entities.ADREntity, error) {
	// Fetch existing ADR
	adr, err := s.adrRepo.GetADR(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ADR: %w", err)
	}

	// Apply updates
	if input.Title != nil {
		if err := s.validationService.ValidateNonEmpty("title", *input.Title); err != nil {
			return nil, err
		}
		adr.Title = *input.Title
	}

	if input.Context != nil {
		if err := s.validationService.ValidateNonEmpty("context", *input.Context); err != nil {
			return nil, err
		}
		adr.Context = *input.Context
	}

	if input.Decision != nil {
		if err := s.validationService.ValidateNonEmpty("decision", *input.Decision); err != nil {
			return nil, err
		}
		adr.Decision = *input.Decision
	}

	if input.Consequences != nil {
		if err := s.validationService.ValidateNonEmpty("consequences", *input.Consequences); err != nil {
			return nil, err
		}
		adr.Consequences = *input.Consequences
	}

	if input.Alternatives != nil {
		adr.Alternatives = *input.Alternatives
	}

	if input.Status != nil {
		if !entities.IsValidADRStatus(*input.Status) {
			return nil, fmt.Errorf("%w: invalid ADR status: %s", tmerrors.ErrInvalidArgument, *input.Status)
		}
		adr.Status = *input.Status
	}

	// Update timestamp
	adr.UpdatedAt = time.Now().UTC()

	// Persist updates
	if err := s.adrRepo.UpdateADR(ctx, adr); err != nil {
		return nil, fmt.Errorf("failed to update ADR: %w", err)
	}

	return adr, nil
}

// SupersedeADR marks an ADR as superseded by another ADR
func (s *ADRApplicationService) SupersedeADR(ctx context.Context, adrID, supersededByID string) error {
	// Validate both ADRs exist
	adr, err := s.adrRepo.GetADR(ctx, adrID)
	if err != nil {
		return fmt.Errorf("ADR not found: %w", err)
	}

	_, err = s.adrRepo.GetADR(ctx, supersededByID)
	if err != nil {
		return fmt.Errorf("superseding ADR not found: %w", err)
	}

	// Update status and superseded_by
	adr.Status = string(entities.ADRStatusSuperseded)
	adr.SupersededBy = &supersededByID
	adr.UpdatedAt = time.Now().UTC()

	// Persist updates
	if err := s.adrRepo.UpdateADR(ctx, adr); err != nil {
		return fmt.Errorf("failed to supersede ADR: %w", err)
	}

	return nil
}

// DeprecateADR marks an ADR as deprecated
func (s *ADRApplicationService) DeprecateADR(ctx context.Context, adrID string) error {
	// Validate ADR exists
	adr, err := s.adrRepo.GetADR(ctx, adrID)
	if err != nil {
		return fmt.Errorf("ADR not found: %w", err)
	}

	// Update status
	adr.Status = string(entities.ADRStatusDeprecated)
	adr.UpdatedAt = time.Now().UTC()

	// Persist updates
	if err := s.adrRepo.UpdateADR(ctx, adr); err != nil {
		return fmt.Errorf("failed to deprecate ADR: %w", err)
	}

	return nil
}

// GetADR retrieves an ADR by ID
func (s *ADRApplicationService) GetADR(ctx context.Context, adrID string) (*entities.ADREntity, error) {
	adr, err := s.adrRepo.GetADR(ctx, adrID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ADR: %w", err)
	}
	return adr, nil
}

// ListADRs returns all ADRs, optionally filtered by track
func (s *ADRApplicationService) ListADRs(ctx context.Context, trackID *string) ([]*entities.ADREntity, error) {
	adrs, err := s.adrRepo.ListADRs(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to list ADRs: %w", err)
	}
	return adrs, nil
}

// GetADRsByTrack returns all ADRs for a specific track
func (s *ADRApplicationService) GetADRsByTrack(ctx context.Context, trackID string) ([]*entities.ADREntity, error) {
	adrs, err := s.adrRepo.GetADRsByTrack(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ADRs by track: %w", err)
	}
	return adrs, nil
}
