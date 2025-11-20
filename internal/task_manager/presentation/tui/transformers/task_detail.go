package transformers

import (
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
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

	// Pre-compute display fields for task
	vm.StatusLabel = GetTaskStatusLabel(task.Status)
	vm.StatusColor = GetTaskColor(task.Status)
	vm.Icon = GetTaskIcon(task.Status)

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
			// Pre-computed display fields
			StatusLabel: GetTrackStatusLabel(track.Status),
			StatusColor: GetTrackColor(track.Status),
			Icon:        GetTrackIcon(track.Status),
		}
	}

	// Iteration membership
	for _, iter := range iterations {
		vm.Iterations = append(vm.Iterations, &viewmodels.IterationMembershipViewModel{
			Number: iter.Number,
			Name:   iter.Name,
			Status: iter.Status,
			// Pre-computed display fields
			StatusLabel: GetIterationStatusLabel(iter.Status),
			StatusColor: GetIterationColor(iter.Status),
			Icon:        GetIterationIcon(iter.Status),
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
			// Pre-computed display fields
			StatusLabel: GetACStatusLabel(ac.Status),
			StatusColor: GetACColor(ac.Status),
			IsFailed:    ac.Status == entities.ACStatusFailed,
		}
		vm.AcceptanceCriteria = append(vm.AcceptanceCriteria, acVM)
	}

	return vm
}
