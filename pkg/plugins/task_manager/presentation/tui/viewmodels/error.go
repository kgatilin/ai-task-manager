package viewmodels

// ErrorViewModel represents the view state for an error screen
type ErrorViewModel struct {
	ErrorMessage string
	Details      string
	CanGoBack    bool
	RetryAction  string
}

// NewErrorViewModel creates a new error view model with the given message
func NewErrorViewModel(errorMessage string) *ErrorViewModel {
	return &ErrorViewModel{
		ErrorMessage: errorMessage,
		CanGoBack:    false,
		Details:      "",
		RetryAction:  "",
	}
}
