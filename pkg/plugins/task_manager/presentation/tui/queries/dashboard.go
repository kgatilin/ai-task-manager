package queries

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/transformers"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// LoadRoadmapListData loads all data needed for the dashboard view.
// Returns filtered and transformed view model ready for presentation.
//
// Pre-loads:
// - All iterations
// - Active roadmap
// - All tracks for the roadmap
// - All backlog tasks (not in any iteration)
//
// Eliminates N+1 queries by loading all related data upfront.
func LoadRoadmapListData(
	ctx context.Context,
	repo domain.RoadmapRepository,
) (*viewmodels.RoadmapListViewModel, error) {
	// Fetch all iterations
	iterations, err := repo.ListIterations(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch active roadmap
	roadmap, err := repo.GetActiveRoadmap(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch all tracks for the roadmap
	tracks, err := repo.ListTracks(ctx, roadmap.ID, entities.TrackFilters{})
	if err != nil {
		return nil, err
	}

	// Fetch backlog tasks (not in any iteration)
	backlogTasks, err := repo.GetBacklogTasks(ctx)
	if err != nil {
		return nil, err
	}

	// Transform to view model with filtering
	vm := transformers.TransformToRoadmapListViewModel(roadmap, iterations, tracks, backlogTasks)

	return vm, nil
}
