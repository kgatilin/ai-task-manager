package transformers

import (
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// TransformToTaskDetailViewModel transforms task + ACs + track + iterations to task detail view model
func TransformToTaskDetailViewModel(
	task *entities.TaskEntity,
	acs []*entities.AcceptanceCriteriaEntity,
	track *entities.TrackEntity,
	iterations []*entities.IterationEntity,
) *viewmodels.TaskDetailViewModel {
	vm := viewmodels.NewTaskDetailViewModel(
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Branch,
	)

	// Format timestamps
	vm.CreatedAt = task.CreatedAt.Format("2006-01-02 15:04:05")
	vm.UpdatedAt = task.UpdatedAt.Format("2006-01-02 15:04:05")

	// Track info
	if track != nil {
		vm.TrackInfo = &viewmodels.TrackInfoViewModel{
			ID:          track.ID,
			Title:       track.Title,
			Description: track.Description,
			Status:      track.Status,
		}
	}

	// Iteration membership
	for _, iter := range iterations {
		vm.Iterations = append(vm.Iterations, &viewmodels.IterationMembershipViewModel{
			Number: iter.Number,
			Name:   iter.Name,
			Status: iter.Status,
		})
	}

	// Transform ACs (initially collapsed)
	for _, ac := range acs {
		acVM := &viewmodels.ACDetailViewModel{
			ID:                  ac.ID,
			Description:         ac.Description,
			Status:              string(ac.Status),
			StatusIcon:          ac.StatusIndicator(),
			TestingInstructions: ac.TestingInstructions,
			Notes:               ac.Notes,
			IsExpanded:          false, // Initially collapsed
		}
		vm.AcceptanceCriteria = append(vm.AcceptanceCriteria, acVM)
	}

	return vm
}
