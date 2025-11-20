package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
)

// RoadmapApplicationService handles all roadmap-related operations.
// It orchestrates domain validation and repository persistence.
type RoadmapApplicationService struct {
	roadmapRepo   repositories.RoadmapRepository
	trackRepo     repositories.TrackRepository
	taskRepo      repositories.TaskRepository
	iterationRepo repositories.IterationRepository
	validationSvc *services.ValidationService
}

// NewRoadmapApplicationService creates a new roadmap application service
func NewRoadmapApplicationService(
	roadmapRepo repositories.RoadmapRepository,
	trackRepo repositories.TrackRepository,
	taskRepo repositories.TaskRepository,
	iterationRepo repositories.IterationRepository,
	validationSvc *services.ValidationService,
) *RoadmapApplicationService {
	return &RoadmapApplicationService{
		roadmapRepo:   roadmapRepo,
		trackRepo:     trackRepo,
		taskRepo:      taskRepo,
		iterationRepo: iterationRepo,
		validationSvc: validationSvc,
	}
}

// InitRoadmap creates a new roadmap with validation
func (s *RoadmapApplicationService) InitRoadmap(ctx context.Context, input dto.CreateRoadmapDTO) (*entities.RoadmapEntity, error) {
	// Validate required fields
	if err := s.validationSvc.ValidateNonEmpty("vision", input.Vision); err != nil {
		return nil, err
	}
	if err := s.validationSvc.ValidateNonEmpty("success_criteria", input.SuccessCriteria); err != nil {
		return nil, err
	}

	// Check if roadmap already exists
	existing, err := s.roadmapRepo.GetActiveRoadmap(ctx)
	if err != nil && !errors.Is(err, tmerrors.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing roadmap: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("roadmap already exists: %s - delete it first to create a new one", existing.ID)
	}

	// Generate roadmap ID
	roadmapID := fmt.Sprintf("roadmap-%d", time.Now().UnixNano())
	now := time.Now().UTC()

	// Create roadmap entity
	roadmap, err := entities.NewRoadmapEntity(roadmapID, input.Vision, input.SuccessCriteria, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create roadmap entity: %w", err)
	}

	// Persist roadmap
	if err := s.roadmapRepo.SaveRoadmap(ctx, roadmap); err != nil {
		return nil, fmt.Errorf("failed to save roadmap: %w", err)
	}

	return roadmap, nil
}

// GetRoadmap retrieves the active roadmap
func (s *RoadmapApplicationService) GetRoadmap(ctx context.Context) (*entities.RoadmapEntity, error) {
	return s.roadmapRepo.GetActiveRoadmap(ctx)
}

// UpdateRoadmap updates an existing roadmap
func (s *RoadmapApplicationService) UpdateRoadmap(ctx context.Context, input dto.UpdateRoadmapDTO) (*entities.RoadmapEntity, error) {
	// At least one field must be provided
	if input.Vision == nil && input.SuccessCriteria == nil {
		return nil, fmt.Errorf("at least one field must be provided (vision or success criteria)")
	}

	// Get active roadmap
	roadmap, err := s.roadmapRepo.GetActiveRoadmap(ctx)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Vision != nil {
		if err := s.validationSvc.ValidateNonEmpty("vision", *input.Vision); err != nil {
			return nil, err
		}
		roadmap.Vision = *input.Vision
	}

	if input.SuccessCriteria != nil {
		if err := s.validationSvc.ValidateNonEmpty("success_criteria", *input.SuccessCriteria); err != nil {
			return nil, err
		}
		roadmap.SuccessCriteria = *input.SuccessCriteria
	}

	// Update timestamp
	roadmap.UpdatedAt = time.Now().UTC()

	// Persist changes
	if err := s.roadmapRepo.UpdateRoadmap(ctx, roadmap); err != nil {
		return nil, fmt.Errorf("failed to update roadmap: %w", err)
	}

	return roadmap, nil
}

// GetFullOverview retrieves a complete roadmap overview with all related entities
func (s *RoadmapApplicationService) GetFullOverview(ctx context.Context, options dto.RoadmapOverviewOptions) (*dto.RoadmapOverviewDTO, error) {
	// Get active roadmap
	roadmap, err := s.roadmapRepo.GetActiveRoadmap(ctx)
	if err != nil {
		return nil, err
	}

	// Get all tracks
	tracks, err := s.trackRepo.ListTracks(ctx, roadmap.ID, entities.TrackFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tracks: %w", err)
	}

	// Get all tasks
	tasks, err := s.taskRepo.ListTasks(ctx, entities.TaskFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	// Get all iterations
	iterations, err := s.iterationRepo.ListIterations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list iterations: %w", err)
	}

	// Convert to interfaces for DTO
	trackInterfaces := make([]interface{}, len(tracks))
	for i, t := range tracks {
		trackInterfaces[i] = t
	}

	taskInterfaces := make([]interface{}, len(tasks))
	for i, t := range tasks {
		taskInterfaces[i] = t
	}

	iterationInterfaces := make([]interface{}, len(iterations))
	for i, it := range iterations {
		iterationInterfaces[i] = it
	}

	return &dto.RoadmapOverviewDTO{
		Roadmap:    roadmap,
		Tracks:     trackInterfaces,
		Tasks:      taskInterfaces,
		Iterations: iterationInterfaces,
		ADRs:       []interface{}{}, // ADRs can be added later if needed
	}, nil
}
