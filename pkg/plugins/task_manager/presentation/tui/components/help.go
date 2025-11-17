package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// Help wraps bubbles/help with consistent width handling
type Help struct {
	model help.Model
	width int
}

// NewHelp creates a new help component
func NewHelp() Help {
	return Help{
		model: help.New(),
		width: 0,
	}
}

// SetWidth updates the help width
func (h *Help) SetWidth(width int) {
	h.width = width
	h.model.Width = width
}

// ShortHelpView renders short help (4-5 key bindings)
func (h Help) ShortHelpView(groups []key.Binding) string {
	return h.model.ShortHelpView(groups)
}

// FullHelpView renders full help (all key bindings)
func (h Help) FullHelpView(groups [][]key.Binding) string {
	return h.model.FullHelpView(groups)
}

// ShowAll toggles between short and full help
func (h *Help) ShowAll(show bool) {
	h.model.ShowAll = show
}
