package queries

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/transformers"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// LoadTaskDetailData loads task detail data for a specific task.
// Returns task + ACs + track + iteration membership transformed into view model ready for presentation.
//
// Pre-loads:
// - Task entity
// - All acceptance criteria for the task
// - Track entity that owns the task
// - All iterations the task belongs to
//
// Eliminates N+1 queries by loading all related data upfront.
func LoadTaskDetailData(
	ctx context.Context,
	repo domain.RoadmapRepository,
	taskID string,
) (*viewmodels.TaskDetailViewModel, error) {
	// Fetch task
	task, err := repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Fetch ACs for the task
	acs, err := repo.ListAC(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Fetch track info
	track, err := repo.GetTrack(ctx, task.TrackID)
	if err != nil {
		return nil, err
	}

	// Fetch iteration membership
	iterations, err := repo.GetIterationsForTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Transform to view model
	vm := transformers.TransformToTaskDetailViewModel(task, acs, track, iterations)

	return vm, nil
}
