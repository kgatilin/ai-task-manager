package application

import (
	"context"
	"fmt"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/dto"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/repositories"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/services"
)

// TrackApplicationService handles all track-related operations.
// It orchestrates domain validation and repository persistence.
type TrackApplicationService struct {
	trackRepo      repositories.TrackRepository
	roadmapRepo    repositories.RoadmapRepository
	aggregateRepo  repositories.AggregateRepository
	validationSvc  *services.ValidationService
}

// NewTrackApplicationService creates a new track application service
func NewTrackApplicationService(
	trackRepo repositories.TrackRepository,
	roadmapRepo repositories.RoadmapRepository,
	aggregateRepo repositories.AggregateRepository,
	validationSvc *services.ValidationService,
) *TrackApplicationService {
	return &TrackApplicationService{
		trackRepo:     trackRepo,
		roadmapRepo:   roadmapRepo,
		aggregateRepo: aggregateRepo,
		validationSvc: validationSvc,
	}
}

// CreateTrack creates a new track with validation
func (s *TrackApplicationService) CreateTrack(ctx context.Context, input dto.CreateTrackDTO) (*entities.TrackEntity, error) {
	// Generate track ID
	projectCode := s.aggregateRepo.GetProjectCode(ctx)
	nextNum, err := s.aggregateRepo.GetNextSequenceNumber(ctx, "track")
	if err != nil {
		return nil, fmt.Errorf("failed to generate track ID: %w", err)
	}
	id := fmt.Sprintf("%s-track-%d", projectCode, nextNum)

	// Validate track ID format
	if err := s.validationSvc.ValidateTrackID(id); err != nil {
		return nil, err
	}

	// Validate title is non-empty
	if err := s.validationSvc.ValidateNonEmpty("title", input.Title); err != nil {
		return nil, err
	}

	// Validate rank is in valid range
	if err := s.validationSvc.ValidateRank(input.Rank); err != nil {
		return nil, err
	}

	// Verify roadmap exists
	_, err = s.roadmapRepo.GetRoadmap(ctx, input.RoadmapID)
	if err != nil {
		return nil, fmt.Errorf("roadmap not found: %w", err)
	}

	// Set default status if not provided
	status := input.Status
	if status == "" {
		status = string(entities.TrackStatusNotStarted)
	}

	// Create track entity
	now := time.Now().UTC()
	track, err := entities.NewTrackEntity(
		id,
		input.RoadmapID,
		input.Title,
		input.Description,
		status,
		input.Rank,
		[]string{}, // No dependencies initially
		now,
		now,
	)
	if err != nil {
		return nil, err
	}

	// Persist track
	if err := s.trackRepo.SaveTrack(ctx, track); err != nil {
		return nil, err
	}

	return track, nil
}

// UpdateTrack updates an existing track
func (s *TrackApplicationService) UpdateTrack(ctx context.Context, input dto.UpdateTrackDTO) (*entities.TrackEntity, error) {
	// Fetch existing track
	track, err := s.trackRepo.GetTrack(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Title != nil {
		if err := s.validationSvc.ValidateNonEmpty("title", *input.Title); err != nil {
			return nil, err
		}
		track.Title = *input.Title
	}

	if input.Description != nil {
		track.Description = *input.Description
	}

	if input.Status != nil {
		if err := track.TransitionTo(*input.Status); err != nil {
			return nil, err
		}
	}

	if input.Rank != nil {
		if err := s.validationSvc.ValidateRank(*input.Rank); err != nil {
			return nil, err
		}
		track.Rank = *input.Rank
	}

	// Update timestamp
	track.UpdatedAt = time.Now().UTC()

	// Persist changes
	if err := s.trackRepo.UpdateTrack(ctx, track); err != nil {
		return nil, err
	}

	return track, nil
}

// DeleteTrack removes a track
func (s *TrackApplicationService) DeleteTrack(ctx context.Context, trackID string) error {
	// Verify track exists before deleting
	_, err := s.trackRepo.GetTrack(ctx, trackID)
	if err != nil {
		return err
	}

	return s.trackRepo.DeleteTrack(ctx, trackID)
}

// GetTrack retrieves a track by ID
func (s *TrackApplicationService) GetTrack(ctx context.Context, trackID string) (*entities.TrackEntity, error) {
	return s.trackRepo.GetTrack(ctx, trackID)
}

// ListTracks returns all tracks for a roadmap, optionally filtered
func (s *TrackApplicationService) ListTracks(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
	// Verify roadmap exists
	_, err := s.roadmapRepo.GetRoadmap(ctx, roadmapID)
	if err != nil {
		return nil, fmt.Errorf("roadmap not found: %w", err)
	}

	return s.trackRepo.ListTracks(ctx, roadmapID, filters)
}

// GetTrackWithTasks retrieves a track with all its tasks
func (s *TrackApplicationService) GetTrackWithTasks(ctx context.Context, trackID string) (*entities.TrackEntity, error) {
	return s.trackRepo.GetTrackWithTasks(ctx, trackID)
}

// AddDependency adds a dependency from trackID to dependsOnID
func (s *TrackApplicationService) AddDependency(ctx context.Context, trackID, dependsOnID string) error {
	// Validate both tracks exist
	_, err := s.trackRepo.GetTrack(ctx, trackID)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}

	_, err = s.trackRepo.GetTrack(ctx, dependsOnID)
	if err != nil {
		return fmt.Errorf("dependency track not found: %w", err)
	}

	// Prevent self-dependency
	if trackID == dependsOnID {
		return fmt.Errorf("%w: track cannot depend on itself", pluginsdk.ErrInvalidArgument)
	}

	// Add dependency
	if err := s.trackRepo.AddTrackDependency(ctx, trackID, dependsOnID); err != nil {
		return err
	}

	// Check for cycles
	if err := s.trackRepo.ValidateNoCycles(ctx, trackID); err != nil {
		// Rollback by removing the dependency
		_ = s.trackRepo.RemoveTrackDependency(ctx, trackID, dependsOnID)
		return fmt.Errorf("circular dependency detected: %w", err)
	}

	return nil
}

// RemoveDependency removes a dependency from trackID to dependsOnID
func (s *TrackApplicationService) RemoveDependency(ctx context.Context, trackID, dependsOnID string) error {
	return s.trackRepo.RemoveTrackDependency(ctx, trackID, dependsOnID)
}

// GetDependencies returns the IDs of all tracks that trackID depends on
func (s *TrackApplicationService) GetDependencies(ctx context.Context, trackID string) ([]string, error) {
	// Verify track exists
	_, err := s.trackRepo.GetTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}

	return s.trackRepo.GetTrackDependencies(ctx, trackID)
}

// GetActiveRoadmap returns the active roadmap for the current project
func (s *TrackApplicationService) GetActiveRoadmap(ctx context.Context) (*entities.RoadmapEntity, error) {
	return s.roadmapRepo.GetActiveRoadmap(ctx)
}
