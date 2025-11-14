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

// TaskApplicationService handles all task-related operations.
// It orchestrates domain validation and repository persistence.
type TaskApplicationService struct {
	taskRepo      repositories.TaskRepository
	trackRepo     repositories.TrackRepository
	aggregateRepo repositories.AggregateRepository
	validationSvc *services.ValidationService
}

// NewTaskApplicationService creates a new task application service
func NewTaskApplicationService(
	taskRepo repositories.TaskRepository,
	trackRepo repositories.TrackRepository,
	aggregateRepo repositories.AggregateRepository,
	validationSvc *services.ValidationService,
) *TaskApplicationService {
	return &TaskApplicationService{
		taskRepo:      taskRepo,
		trackRepo:     trackRepo,
		aggregateRepo: aggregateRepo,
		validationSvc: validationSvc,
	}
}

// CreateTask creates a new task with validation
func (s *TaskApplicationService) CreateTask(ctx context.Context, input dto.CreateTaskDTO) (*entities.TaskEntity, error) {
	// Generate task ID
	projectCode := s.aggregateRepo.GetProjectCode(ctx)
	nextNum, err := s.aggregateRepo.GetNextSequenceNumber(ctx, "task")
	if err != nil {
		return nil, fmt.Errorf("failed to generate task ID: %w", err)
	}
	id := fmt.Sprintf("%s-task-%d", projectCode, nextNum)

	// Validate title is non-empty
	if err := s.validationSvc.ValidateNonEmpty("title", input.Title); err != nil {
		return nil, err
	}

	// Validate rank is in valid range
	if err := s.validationSvc.ValidateRank(input.Rank); err != nil {
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
		status = string(entities.TaskStatusTodo)
	}

	// Validate status
	if !entities.IsValidTaskStatus(status) {
		return nil, fmt.Errorf("%w: invalid task status: %s", pluginsdk.ErrInvalidArgument, status)
	}

	// Create task entity
	now := time.Now().UTC()
	task, err := entities.NewTaskEntity(
		id,
		input.TrackID,
		input.Title,
		input.Description,
		status,
		input.Rank,
		"", // No branch initially
		now,
		now,
	)
	if err != nil {
		return nil, err
	}

	// Persist task
	if err := s.taskRepo.SaveTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask updates an existing task
func (s *TaskApplicationService) UpdateTask(ctx context.Context, input dto.UpdateTaskDTO) (*entities.TaskEntity, error) {
	// Fetch existing task
	task, err := s.taskRepo.GetTask(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Title != nil {
		if err := s.validationSvc.ValidateNonEmpty("title", *input.Title); err != nil {
			return nil, err
		}
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Status != nil {
		if err := task.TransitionTo(*input.Status); err != nil {
			return nil, err
		}
	}

	if input.Rank != nil {
		if err := s.validationSvc.ValidateRank(*input.Rank); err != nil {
			return nil, err
		}
		task.Rank = *input.Rank
	}

	if input.TrackID != nil {
		// Verify new track exists
		_, err := s.trackRepo.GetTrack(ctx, *input.TrackID)
		if err != nil {
			return nil, fmt.Errorf("track not found: %w", err)
		}
		task.TrackID = *input.TrackID
	}

	// Update timestamp
	task.UpdatedAt = time.Now().UTC()

	// Persist changes
	if err := s.taskRepo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask removes a task
func (s *TaskApplicationService) DeleteTask(ctx context.Context, taskID string) error {
	// Verify task exists before deleting
	_, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	return s.taskRepo.DeleteTask(ctx, taskID)
}

// MoveTask moves a task to a different track
func (s *TaskApplicationService) MoveTask(ctx context.Context, taskID, newTrackID string) error {
	// Verify task exists
	_, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Verify new track exists
	_, err = s.trackRepo.GetTrack(ctx, newTrackID)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}

	// Move task using repository method
	return s.taskRepo.MoveTaskToTrack(ctx, taskID, newTrackID)
}

// GetTask retrieves a task by ID
func (s *TaskApplicationService) GetTask(ctx context.Context, taskID string) (*entities.TaskEntity, error) {
	return s.taskRepo.GetTask(ctx, taskID)
}

// ListTasks returns all tasks, optionally filtered
func (s *TaskApplicationService) ListTasks(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error) {
	return s.taskRepo.ListTasks(ctx, filters)
}

// GetBacklogTasks returns all tasks with status "todo"
func (s *TaskApplicationService) GetBacklogTasks(ctx context.Context) ([]*entities.TaskEntity, error) {
	return s.taskRepo.GetBacklogTasks(ctx)
}
