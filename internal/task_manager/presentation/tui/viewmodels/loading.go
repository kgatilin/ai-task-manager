package viewmodels

// LoadingViewModel represents the view state for a loading screen
type LoadingViewModel struct {
	Message     string
	ShowSpinner bool
}

// NewLoadingViewModel creates a new loading view model with default values
func NewLoadingViewModel(message string) *LoadingViewModel {
	return &LoadingViewModel{
		Message:     message,
		ShowSpinner: true,
	}
}
