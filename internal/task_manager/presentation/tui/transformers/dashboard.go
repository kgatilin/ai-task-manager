package transformers

import (
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

// TransformToRoadmapListViewModel transforms domain entities to the dashboard view model
// Applies filtering rules:
// - ActiveIterations: status != "complete"
// - ActiveTracks: status != "complete"
// - BacklogTasks: status != "done" AND status != "cancelled"
func TransformToRoadmapListViewModel(
	roadmap *entities.RoadmapEntity,
	iterations []*entities.IterationEntity,
	tracks []*entities.TrackEntity,
	tasks []*entities.TaskEntity,
) *viewmodels.RoadmapListViewModel {
	vm := viewmodels.NewRoadmapListViewModel()

	// Include roadmap vision and success criteria
	if roadmap != nil {
		vm.Vision = roadmap.Vision
		vm.SuccessCriteria = roadmap.SuccessCriteria
	}

	// Filter and transform iterations (exclude complete)
	for _, iter := range iterations {
		if iter.Status != string(entities.IterationStatusComplete) {
			isCurrent := iter.Status == string(entities.IterationStatusCurrent)
			vm.ActiveIterations = append(vm.ActiveIterations, &viewmodels.IterationCardViewModel{
				Number:      iter.Number,
				Name:        iter.Name,
				Goal:        iter.Goal,
				Status:      iter.Status,
				TaskCount:   len(iter.TaskIDs),
				Deliverable: iter.Deliverable,
				// Pre-computed display fields
				StatusLabel: GetIterationStatusLabel(iter.Status),
				StatusColor: GetIterationColor(iter.Status),
				Icon:        GetIterationIcon(iter.Status),
				IsCurrent:   isCurrent,
			})
		}
	}

	// Build task count map for tracks
	// Only count active tasks (exclude done and cancelled) to match backlog display
	trackTaskCounts := make(map[string]int)
	for _, task := range tasks {
		if task.Status != string(entities.TaskStatusDone) && task.Status != string(entities.TaskStatusCancelled) {
			trackTaskCounts[task.TrackID]++
		}
	}

	// Filter and transform tracks (exclude complete)
	for _, track := range tracks {
		if track.Status != string(entities.TrackStatusComplete) {
			vm.ActiveTracks = append(vm.ActiveTracks, &viewmodels.TrackCardViewModel{
				ID:          track.ID,
				Title:       track.Title,
				Description: track.Description,
				Status:      track.Status,
				TaskCount:   trackTaskCounts[track.ID],
				// Pre-computed display fields
				StatusLabel: GetTrackStatusLabel(track.Status),
				StatusColor: GetTrackColor(track.Status),
				Icon:        GetTrackIcon(track.Status),
			})
		}
	}

	// Filter and transform backlog tasks (exclude done and cancelled)
	for _, task := range tasks {
		if task.Status != string(entities.TaskStatusDone) && task.Status != string(entities.TaskStatusCancelled) {
			vm.BacklogTasks = append(vm.BacklogTasks, &viewmodels.BacklogTaskViewModel{
				ID:          task.ID,
				Title:       task.Title,
				Status:      task.Status,
				TrackID:     task.TrackID,
				Description: task.Description,
				// Pre-computed display fields
				StatusLabel: GetTaskStatusLabel(task.Status),
				StatusColor: GetTaskColor(task.Status),
				Icon:        GetTaskIcon(task.Status),
			})
		}
	}

	return vm
}

// FilterActiveIterations returns iterations with status != "complete"
func FilterActiveIterations(iterations []*entities.IterationEntity) []*entities.IterationEntity {
	active := []*entities.IterationEntity{}
	for _, iter := range iterations {
		if iter.Status != string(entities.IterationStatusComplete) {
			active = append(active, iter)
		}
	}
	return active
}

// FilterActiveTracks returns tracks with status != "complete"
func FilterActiveTracks(tracks []*entities.TrackEntity) []*entities.TrackEntity {
	active := []*entities.TrackEntity{}
	for _, track := range tracks {
		if track.Status != string(entities.TrackStatusComplete) {
			active = append(active, track)
		}
	}
	return active
}

// FilterBacklogTasks returns tasks with status != "done" AND status != "cancelled"
func FilterBacklogTasks(tasks []*entities.TaskEntity) []*entities.TaskEntity {
	backlog := []*entities.TaskEntity{}
	for _, task := range tasks {
		if task.Status != string(entities.TaskStatusDone) && task.Status != string(entities.TaskStatusCancelled) {
			backlog = append(backlog, task)
		}
	}
	return backlog
}
