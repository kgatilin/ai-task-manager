package viewmodels

// ProgressViewModel represents a progress bar view model
type ProgressViewModel struct {
	Completed int
	Total     int
	Percent   float64
}

// NewProgressViewModel creates a new progress view model
func NewProgressViewModel(completed, total int) *ProgressViewModel {
	percent := 0.0
	if total > 0 {
		percent = float64(completed) / float64(total)
	}
	return &ProgressViewModel{
		Completed: completed,
		Total:     total,
		Percent:   percent,
	}
}

// TaskRowViewModel represents a task row in the iteration detail view
type TaskRowViewModel struct {
	ID          string
	Title       string
	Status      string
	Description string
	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	Icon        string // Status icon
}

// IterationACViewModel represents an AC row with skipped status support
type IterationACViewModel struct {
	ID                  string
	Description         string
	Status              string
	StatusIcon          string
	TestingInstructions string
	Notes               string
	IsExpanded          bool // Whether testing instructions are visible (same as ACDetailViewModel)
	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	IsFailed    bool   // True if status is "failed" (for highlighting)
}

// TaskACGroupViewModel represents a task with its ACs grouped together
type TaskACGroupViewModel struct {
	Task *TaskRowViewModel
	ACs  []*IterationACViewModel
}

// IterationDetailViewModel represents the iteration detail view with tasks and ACs
type IterationDetailViewModel struct {
	Number      int
	Name        string
	Goal        string
	Deliverable string
	Status      string
	StartedAt   string
	CompletedAt string

	// Task grouping by status
	TODOTasks       []*TaskRowViewModel
	InProgressTasks []*TaskRowViewModel
	ReviewTasks     []*TaskRowViewModel
	DoneTasks       []*TaskRowViewModel

	// All ACs for the iteration
	AcceptanceCriteria []*IterationACViewModel

	// ACs grouped by task (for AC view in ACs tab)
	TaskACs []*TaskACGroupViewModel

	// Progress tracking
	Progress *ProgressViewModel

	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	Icon        string // Status icon
}

// NewIterationDetailViewModel creates a new iteration detail view model
func NewIterationDetailViewModel(number int, name, goal, deliverable, status string) *IterationDetailViewModel {
	return &IterationDetailViewModel{
		Number:             number,
		Name:               name,
		Goal:               goal,
		Deliverable:        deliverable,
		Status:             status,
		TODOTasks:          []*TaskRowViewModel{},
		InProgressTasks:    []*TaskRowViewModel{},
		ReviewTasks:        []*TaskRowViewModel{},
		DoneTasks:          []*TaskRowViewModel{},
		AcceptanceCriteria: []*IterationACViewModel{},
		TaskACs:            []*TaskACGroupViewModel{},
		Progress:           NewProgressViewModel(0, 0),
	}
}
