package viewmodels

// ACDetailViewModel represents an AC row with expandable testing instructions
type ACDetailViewModel struct {
	ID                  string
	Description         string
	Status              string
	StatusIcon          string
	TestingInstructions string
	Notes               string
	IsExpanded          bool // Whether testing instructions are visible
}

// TrackInfoViewModel represents track context for task detail view
type TrackInfoViewModel struct {
	ID          string
	Title       string
	Description string
	Status      string
}

// IterationMembershipViewModel represents iteration membership for a task
type IterationMembershipViewModel struct {
	Number int
	Name   string
	Status string
}

// TaskDetailViewModel represents the task detail view with expandable ACs
type TaskDetailViewModel struct {
	// Task metadata
	ID          string
	Title       string
	Description string
	Status      string
	Branch      string
	CreatedAt   string
	UpdatedAt   string

	// Track info
	TrackInfo *TrackInfoViewModel

	// Iteration membership
	Iterations []*IterationMembershipViewModel

	// Acceptance criteria with expandable testing instructions
	AcceptanceCriteria []*ACDetailViewModel
}

// NewTaskDetailViewModel creates a new task detail view model
func NewTaskDetailViewModel(id, title, description, status, branch string) *TaskDetailViewModel {
	return &TaskDetailViewModel{
		ID:                 id,
		Title:              title,
		Description:        description,
		Status:             status,
		Branch:             branch,
		Iterations:         []*IterationMembershipViewModel{},
		AcceptanceCriteria: []*ACDetailViewModel{},
	}
}
