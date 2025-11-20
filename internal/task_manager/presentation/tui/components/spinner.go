package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner wraps bubbles/spinner with project-specific styling.
// This provides a centralized, consistent spinner implementation across the TUI.
type Spinner struct {
	model spinner.Model
}

// NewSpinner creates a spinner with preset dot style and accent color.
// The spinner is pre-configured with:
// - Dot spinner style (animated dots)
// - Accent color (magenta/pink #205)
func NewSpinner() Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Spinner{model: s}
}

// Tick returns the spinner tick command.
// This should be called in Init() to start the animation.
func (s Spinner) Tick() tea.Cmd {
	return s.model.Tick
}

// Update handles spinner tick messages.
// This advances the animation to the next frame.
func (s *Spinner) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return cmd
}

// View renders the current spinner animation frame.
func (s Spinner) View() string {
	return s.model.View()
}
