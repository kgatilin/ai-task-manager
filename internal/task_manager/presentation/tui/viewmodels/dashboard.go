package viewmodels

// IterationCardViewModel represents an iteration card in the dashboard
type IterationCardViewModel struct {
	Number      int
	Name        string
	Goal        string
	Status      string
	TaskCount   int
	Deliverable string
	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	Icon        string // Status icon
	IsCurrent   bool   // True if this is the current iteration
}

// TrackCardViewModel represents a track card in the dashboard
type TrackCardViewModel struct {
	ID          string
	Title       string
	Description string
	Status      string
	TaskCount   int
	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	Icon        string // Status icon
}

// BacklogTaskViewModel represents a backlog task in the dashboard
type BacklogTaskViewModel struct {
	ID          string
	Title       string
	Status      string
	TrackID     string
	Description string
	// Display fields (pre-computed by transformer)
	StatusLabel string // Human-readable status label
	StatusColor string // Color name for status styling
	Icon        string // Status icon
}

// RoadmapListViewModel represents the dashboard view with filtered data
type RoadmapListViewModel struct {
	Vision           string
	SuccessCriteria  string
	ActiveIterations []*IterationCardViewModel
	ActiveTracks     []*TrackCardViewModel
	BacklogTasks     []*BacklogTaskViewModel
}

// NewRoadmapListViewModel creates a new dashboard view model
func NewRoadmapListViewModel() *RoadmapListViewModel {
	return &RoadmapListViewModel{
		ActiveIterations: []*IterationCardViewModel{},
		ActiveTracks:     []*TrackCardViewModel{},
		BacklogTasks:     []*BacklogTaskViewModel{},
	}
}
