package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// IterationService handles iteration lifecycle management
type IterationService struct{}

// NewIterationService creates a new iteration service
func NewIterationService() *IterationService {
	return &IterationService{}
}

// CanStartIteration validates if an iteration can be started
// Returns error if:
// - Iteration is not in "planned" status
// - Another iteration is already "current"
func (s *IterationService) CanStartIteration(
	ctx context.Context,
	iteration *entities.IterationEntity,
	getCurrentIteration func(context.Context) (*entities.IterationEntity, error),
) error {
	// Check if iteration is in planned status
	if iteration.Status != string(entities.IterationStatusPlanned) {
		return fmt.Errorf("%w: iteration must be in planned status to start (current: %s)",
			pluginsdk.ErrInvalidArgument, iteration.Status)
	}

	// Check if another iteration is already current
	currentIter, err := getCurrentIteration(ctx)
	if err != nil {
		// ErrNotFound is OK (no current iteration)
		if !errors.Is(err, pluginsdk.ErrNotFound) {
			return fmt.Errorf("failed to check for current iteration: %w", err)
		}
	} else if currentIter != nil && currentIter.Number != iteration.Number {
		return fmt.Errorf("%w: iteration %d is already current",
			pluginsdk.ErrInvalidArgument, currentIter.Number)
	}

	return nil
}

// CanCompleteIteration validates if an iteration can be completed
// Returns error if iteration is not in "current" status
func (s *IterationService) CanCompleteIteration(iteration *entities.IterationEntity) error {
	// Check if iteration is in current status
	if iteration.Status != string(entities.IterationStatusCurrent) {
		return fmt.Errorf("%w: iteration must be in current status to complete (current: %s)",
			pluginsdk.ErrInvalidArgument, iteration.Status)
	}

	return nil
}
