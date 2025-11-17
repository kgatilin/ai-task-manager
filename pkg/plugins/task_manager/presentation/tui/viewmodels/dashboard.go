package viewmodels

// IterationCardViewModel represents an iteration card in the dashboard
type IterationCardViewModel struct {
	Number      int
	Name        string
	Goal        string
	Status      string
	TaskCount   int
	Deliverable string
}

// TrackCardViewModel represents a track card in the dashboard
type TrackCardViewModel struct {
	ID          string
	Title       string
	Description string
	Status      string
	TaskCount   int
}

// BacklogTaskViewModel represents a backlog task in the dashboard
type BacklogTaskViewModel struct {
	ID          string
	Title       string
	Status      string
	TrackID     string
	Description string
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
