package transformers

import (
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// Icon constants (duplicated from components to avoid circular dependency)
const (
	iterationPlannedIcon  = "📋"
	iterationCurrentIcon  = "▶"
	iterationCompleteIcon = "✓"
	taskTodoIcon          = "○"
	taskInProgressIcon    = "◐"
	taskReviewIcon        = "◑"
	taskDoneIcon          = "●"
	taskCancelledIcon     = "⊗"
	trackNotStartedIcon   = "○"
	trackInProgressIcon   = "◐"
	trackCompleteIcon     = "●"
	trackBlockedIcon      = "⊠"
	trackWaitingIcon      = "⏸"
)

// GetIterationIcon returns the icon for an iteration status
func GetIterationIcon(status string) string {
	switch status {
	case string(entities.IterationStatusPlanned):
		return iterationPlannedIcon
	case string(entities.IterationStatusCurrent):
		return iterationCurrentIcon
	case string(entities.IterationStatusComplete):
		return iterationCompleteIcon
	default:
		return iterationPlannedIcon
	}
}

// GetIterationColor returns the color name for an iteration status
func GetIterationColor(status string) string {
	switch status {
	case string(entities.IterationStatusPlanned):
		return "info"
	case string(entities.IterationStatusCurrent):
		return "current"
	case string(entities.IterationStatusComplete):
		return "success"
	default:
		return "info"
	}
}

// GetIterationStatusLabel returns a human-readable label for iteration status
func GetIterationStatusLabel(status string) string {
	switch status {
	case string(entities.IterationStatusPlanned):
		return "Planned"
	case string(entities.IterationStatusCurrent):
		return "Current"
	case string(entities.IterationStatusComplete):
		return "Complete"
	default:
		return status
	}
}

// GetTaskIcon returns the icon for a task status
func GetTaskIcon(status string) string {
	switch status {
	case string(entities.TaskStatusTodo):
		return taskTodoIcon
	case string(entities.TaskStatusInProgress):
		return taskInProgressIcon
	case string(entities.TaskStatusReview):
		return taskReviewIcon
	case string(entities.TaskStatusDone):
		return taskDoneIcon
	case string(entities.TaskStatusCancelled):
		return taskCancelledIcon
	default:
		return taskTodoIcon
	}
}

// GetTaskColor returns the color name for a task status
func GetTaskColor(status string) string {
	switch status {
	case string(entities.TaskStatusTodo):
		return "info"
	case string(entities.TaskStatusInProgress):
		return "warning"
	case string(entities.TaskStatusReview):
		return "warning"
	case string(entities.TaskStatusDone):
		return "success"
	case string(entities.TaskStatusCancelled):
		return "muted"
	default:
		return "info"
	}
}

// GetTaskStatusLabel returns a human-readable label for task status
func GetTaskStatusLabel(status string) string {
	switch status {
	case string(entities.TaskStatusTodo):
		return "Todo"
	case string(entities.TaskStatusInProgress):
		return "In Progress"
	case string(entities.TaskStatusReview):
		return "Review"
	case string(entities.TaskStatusDone):
		return "Done"
	case string(entities.TaskStatusCancelled):
		return "Cancelled"
	default:
		return status
	}
}

// GetTrackIcon returns the icon for a track status
func GetTrackIcon(status string) string {
	switch status {
	case string(entities.TrackStatusNotStarted):
		return trackNotStartedIcon
	case string(entities.TrackStatusInProgress):
		return trackInProgressIcon
	case string(entities.TrackStatusComplete):
		return trackCompleteIcon
	case string(entities.TrackStatusBlocked):
		return trackBlockedIcon
	case string(entities.TrackStatusWaiting):
		return trackWaitingIcon
	default:
		return trackNotStartedIcon
	}
}

// GetTrackColor returns the color name for a track status
func GetTrackColor(status string) string {
	switch status {
	case string(entities.TrackStatusNotStarted):
		return "muted"
	case string(entities.TrackStatusInProgress):
		return "warning"
	case string(entities.TrackStatusComplete):
		return "success"
	case string(entities.TrackStatusBlocked):
		return "failed"
	case string(entities.TrackStatusWaiting):
		return "warning"
	default:
		return "muted"
	}
}

// GetTrackStatusLabel returns a human-readable label for track status
func GetTrackStatusLabel(status string) string {
	switch status {
	case string(entities.TrackStatusNotStarted):
		return "Not Started"
	case string(entities.TrackStatusInProgress):
		return "In Progress"
	case string(entities.TrackStatusComplete):
		return "Complete"
	case string(entities.TrackStatusBlocked):
		return "Blocked"
	case string(entities.TrackStatusWaiting):
		return "Waiting"
	default:
		return status
	}
}

// GetACColor returns the color name for an AC status
func GetACColor(status entities.AcceptanceCriteriaStatus) string {
	switch status {
	case entities.ACStatusVerified, entities.ACStatusAutomaticallyVerified:
		return "success"
	case entities.ACStatusFailed:
		return "failed"
	case entities.ACStatusPendingHumanReview:
		return "warning"
	case entities.ACStatusSkipped:
		return "skipped"
	default: // not_started
		return "muted"
	}
}

// GetACStatusLabel returns a human-readable label for AC status
func GetACStatusLabel(status entities.AcceptanceCriteriaStatus) string {
	switch status {
	case entities.ACStatusNotStarted:
		return "Not Started"
	case entities.ACStatusAutomaticallyVerified:
		return "Auto Verified"
	case entities.ACStatusPendingHumanReview:
		return "Pending Review"
	case entities.ACStatusVerified:
		return "Verified"
	case entities.ACStatusFailed:
		return "Failed"
	case entities.ACStatusSkipped:
		return "Skipped"
	default:
		return string(status)
	}
}
