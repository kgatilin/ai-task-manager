package dto

// CreateIterationDTO represents input for creating a new iteration
type CreateIterationDTO struct {
	Number      int
	Name        string
	Goal        string
	Deliverable string
	Status      string
}

// UpdateIterationDTO represents input for updating an iteration
type UpdateIterationDTO struct {
	Number      int
	Name        *string
	Goal        *string
	Deliverable *string
}

// IterationFilters represents filters for listing iterations
type IterationFilters struct {
	Status []string
}
