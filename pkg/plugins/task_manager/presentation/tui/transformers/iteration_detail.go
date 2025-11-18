package transformers

import (
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// TransformToIterationDetailViewModel transforms iteration + tasks + ACs + documents to iteration detail view model
func TransformToIterationDetailViewModel(
	iteration *entities.IterationEntity,
	tasks []*entities.TaskEntity,
	acs []*entities.AcceptanceCriteriaEntity,
	documents []*entities.DocumentEntity,
) *viewmodels.IterationDetailViewModel {
	vm := viewmodels.NewIterationDetailViewModel(
		iteration.Number,
		iteration.Name,
		iteration.Goal,
		iteration.Deliverable,
		iteration.Status,
	)

	// Pre-compute display fields for iteration
	vm.StatusLabel = GetIterationStatusLabel(iteration.Status)
	vm.StatusColor = GetIterationColor(iteration.Status)
	vm.Icon = GetIterationIcon(iteration.Status)

	// Format timestamps
	if iteration.StartedAt != nil {
		vm.StartedAt = iteration.StartedAt.Format("2006-01-02 15:04:05")
	}
	if iteration.CompletedAt != nil {
		vm.CompletedAt = iteration.CompletedAt.Format("2006-01-02 15:04:05")
	}

	// Group tasks by status and create task map
	taskMap := make(map[string]*viewmodels.TaskRowViewModel)
	for _, task := range tasks {
		taskRow := &viewmodels.TaskRowViewModel{
			ID:          task.ID,
			Title:       task.Title,
			Status:      task.Status,
			Description: task.Description,
			// Pre-computed display fields
			StatusLabel: GetTaskStatusLabel(task.Status),
			StatusColor: GetTaskColor(task.Status),
			Icon:        GetTaskIcon(task.Status),
		}

		// Store in map for AC grouping
		taskMap[task.ID] = taskRow

		switch task.Status {
		case string(entities.TaskStatusTodo):
			vm.TODOTasks = append(vm.TODOTasks, taskRow)
		case string(entities.TaskStatusInProgress):
			vm.InProgressTasks = append(vm.InProgressTasks, taskRow)
		case string(entities.TaskStatusReview):
			vm.ReviewTasks = append(vm.ReviewTasks, taskRow)
		case string(entities.TaskStatusDone):
			vm.DoneTasks = append(vm.DoneTasks, taskRow)
		}
	}

	// Transform ACs
	acMap := make(map[string][]*viewmodels.IterationACViewModel)
	for _, ac := range acs {
		acVM := &viewmodels.IterationACViewModel{
			ID:                  ac.ID,
			Description:         ac.Description,
			Status:              string(ac.Status),
			StatusIcon:          ac.StatusIndicator(),
			TestingInstructions: ac.TestingInstructions,
			Notes:               ac.Notes,
			// Pre-computed display fields
			StatusLabel: GetACStatusLabel(ac.Status),
			StatusColor: GetACColor(ac.Status),
			IsFailed:    ac.Status == entities.ACStatusFailed,
		}
		vm.AcceptanceCriteria = append(vm.AcceptanceCriteria, acVM)

		// Group ACs by their task ID
		acMap[ac.TaskID] = append(acMap[ac.TaskID], acVM)
	}

	// Build TaskACs groups in task order (preserve order from tasks slice)
	for _, task := range tasks {
		if taskACs, exists := acMap[task.ID]; exists {
			group := &viewmodels.TaskACGroupViewModel{
				Task: taskMap[task.ID],
				ACs:  taskACs,
			}
			vm.TaskACs = append(vm.TaskACs, group)
		}
	}

	// Transform and add documents
	vm.Documents = TransformDocumentsToListItems(documents)

	// Calculate progress (done tasks / total tasks)
	totalTasks := len(tasks)
	doneTasks := len(vm.DoneTasks)
	vm.Progress = viewmodels.NewProgressViewModel(doneTasks, totalTasks)

	return vm
}
