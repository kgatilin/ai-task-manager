package dto

// CreateTrackDTO represents input for creating a new track
type CreateTrackDTO struct {
	RoadmapID   string
	Title       string
	Description string
	Status      string
	Rank        int
}

// UpdateTrackDTO represents input for updating a track
// Pointer fields indicate optional updates (nil = no change)
type UpdateTrackDTO struct {
	ID          string
	Title       *string
	Description *string
	Status      *string
	Rank        *int
}

// TrackListFilters represents filters for listing tracks
type TrackListFilters struct {
	Status []string
}
