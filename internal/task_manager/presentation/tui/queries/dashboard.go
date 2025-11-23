package queries

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/transformers"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

// LoadRoadmapListData loads all data needed for the dashboard view.
// Returns filtered and transformed view model ready for presentation.
//
// Pre-loads:
// - All iterations
// - Active roadmap
// - All tracks for the roadmap
// - All tasks (for accurate track counts; transformer filters for display)
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

	// Fetch ALL tasks (for accurate track counts)
	// The transformer will filter appropriately for display and counting
	allTasks, err := repo.ListTasks(ctx, entities.TaskFilters{})
	if err != nil {
		return nil, err
	}

	// Transform to view model with filtering
	vm := transformers.TransformToRoadmapListViewModel(roadmap, iterations, tracks, allTasks)

	return vm, nil
}
