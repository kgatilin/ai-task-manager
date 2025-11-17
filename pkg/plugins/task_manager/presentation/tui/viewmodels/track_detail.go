package viewmodels

// TrackDetailTaskViewModel represents a task row in the track detail view
type TrackDetailTaskViewModel struct {
	ID          string
	Title       string
	Status      string
	Description string
	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	Icon        string // Status icon
}

// TrackDetailViewModel represents the track detail view with track info and tasks
type TrackDetailViewModel struct {
	ID               string
	Title            string
	Description      string
	Status           string
	StatusLabel      string // Display-friendly status (e.g., "In Progress", "Complete")
	Rank             int
	Dependencies     []string
	DependencyLabels []string // Display-friendly dependency names

	// Task grouping by status
	TODOTasks       []*TrackDetailTaskViewModel
	InProgressTasks []*TrackDetailTaskViewModel
	DoneTasks       []*TrackDetailTaskViewModel

	// Progress tracking
	Progress *ProgressViewModel

	// Display fields (pre-computed by transformer)
	StatusColor string // Color name for status styling
	Icon        string // Status icon
}

// NewTrackDetailViewModel creates a new track detail view model
func NewTrackDetailViewModel(id, title, description, status, statusLabel string, rank int, dependencies []string, dependencyLabels []string) *TrackDetailViewModel {
	return &TrackDetailViewModel{
		ID:               id,
		Title:            title,
		Description:      description,
		Status:           status,
		StatusLabel:      statusLabel,
		Rank:             rank,
		Dependencies:     dependencies,
		DependencyLabels: dependencyLabels,
		TODOTasks:        []*TrackDetailTaskViewModel{},
		InProgressTasks:  []*TrackDetailTaskViewModel{},
		DoneTasks:        []*TrackDetailTaskViewModel{},
		Progress:         NewProgressViewModel(0, 0),
	}
}
