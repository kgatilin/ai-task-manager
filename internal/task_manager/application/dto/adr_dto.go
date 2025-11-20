package dto

// CreateADRDTO represents input for creating a new ADR
type CreateADRDTO struct {
	TrackID      string
	Title        string
	Context      string
	Decision     string
	Consequences string
	Alternatives string
	Status       string
}

// UpdateADRDTO represents input for updating an ADR
type UpdateADRDTO struct {
	ID           string
	Title        *string
	Context      *string
	Decision     *string
	Consequences *string
	Alternatives *string
	Status       *string
}

// ADRFilters represents filters for listing ADRs
type ADRFilters struct {
	TrackID *string
	Status  []string
}
