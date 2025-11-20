package dto

// CreateTaskDTO represents input for creating a new task
type CreateTaskDTO struct {
	TrackID     string
	Title       string
	Description string
	Status      string
	Rank        int
}

// UpdateTaskDTO represents input for updating a task
type UpdateTaskDTO struct {
	ID          string
	Title       *string
	Description *string
	Status      *string
	Rank        *int
	TrackID     *string
}

// TaskListFilters represents filters for listing tasks
type TaskListFilters struct {
	Status  []string
	TrackID *string
}
