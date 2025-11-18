package queries

import (
	"context"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/transformers"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// LoadIterationDetailData loads iteration detail data for a specific iteration.
// Returns iteration + tasks + ACs + documents transformed into view model ready for presentation.
//
// Pre-loads:
// - Iteration entity
// - All tasks in the iteration
// - All acceptance criteria for all tasks in the iteration
// - All documents attached to the iteration
//
// Eliminates N+1 queries by loading all related data upfront.
func LoadIterationDetailData(
	ctx context.Context,
	repo domain.RoadmapRepository,
	iterationNumber int,
) (*viewmodels.IterationDetailViewModel, error) {
	// Fetch iteration
	iteration, err := repo.GetIteration(ctx, iterationNumber)
	if err != nil {
		return nil, err
	}

	// Fetch iteration tasks
	tasks, err := repo.GetIterationTasks(ctx, iterationNumber)
	if err != nil {
		return nil, err
	}

	// Fetch ACs for all tasks in the iteration
	acs, err := repo.ListACByIteration(ctx, iterationNumber)
	if err != nil {
		return nil, err
	}

	// Fetch documents attached to the iteration
	documents, err := repo.FindDocumentsByIteration(ctx, iterationNumber)
	if err != nil {
		// Log error but continue (documents are non-critical)
		documents = []*entities.DocumentEntity{}
	}

	// Transform to view model
	vm := transformers.TransformToIterationDetailViewModel(iteration, tasks, acs, documents)

	return vm, nil
}
