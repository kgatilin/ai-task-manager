package task_manager

// Event type constants for task-manager plugin
// These events are emitted when roadmap, track, task, and iteration operations occur

// Roadmap Events
const (
	EventRoadmapCreated = "task-manager.roadmap.created"
	EventRoadmapUpdated = "task-manager.roadmap.updated"
)

// Track Events
const (
	EventTrackCreated       = "task-manager.track.created"
	EventTrackUpdated       = "task-manager.track.updated"
	EventTrackStatusChanged = "task-manager.track.status_changed"
	EventTrackCompleted     = "task-manager.track.completed"
	EventTrackBlocked       = "task-manager.track.blocked"
)

// Task Events
const (
	EventTaskCreated       = "task-manager.task.created"
	EventTaskUpdated       = "task-manager.task.updated"
	EventTaskStatusChanged = "task-manager.task.status_changed"
	EventTaskCompleted     = "task-manager.task.completed"

	// Deprecated: use EventTaskCreated instead
	EventTaskDeleted = "task.deleted"
)

// Iteration Events
const (
	EventIterationCreated   = "task-manager.iteration.created"
	EventIterationStarted   = "task-manager.iteration.started"
	EventIterationCompleted = "task-manager.iteration.completed"
	EventIterationUpdated   = "task-manager.iteration.updated"
)

// Plugin source name
const PluginSourceName = "task-manager"
