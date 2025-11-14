package dto

// CreateRoadmapDTO represents input for creating a new roadmap
type CreateRoadmapDTO struct {
	Vision          string
	SuccessCriteria string
}

// UpdateRoadmapDTO represents input for updating a roadmap
// Pointer fields indicate optional updates (nil = no change)
type UpdateRoadmapDTO struct {
	Vision          *string
	SuccessCriteria *string
}

// RoadmapOverviewDTO represents a complete roadmap overview for display
type RoadmapOverviewDTO struct {
	Roadmap    interface{} // RoadmapEntity
	Tracks     []interface{}
	Tasks      []interface{}
	Iterations []interface{}
	ADRs       []interface{}
}

// RoadmapOverviewOptions represents options for retrieving roadmap overview
type RoadmapOverviewOptions struct {
	Verbose  bool
	Sections []string // vision, tracks, iterations, backlog
}

// ShouldShowSection checks if a section should be shown based on options
func (o *RoadmapOverviewOptions) ShouldShowSection(section string) bool {
	if len(o.Sections) == 0 {
		return true
	}
	for _, s := range o.Sections {
		if s == section {
			return true
		}
	}
	return false
}
