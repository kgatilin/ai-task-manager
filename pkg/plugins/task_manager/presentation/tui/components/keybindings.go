package components

import "github.com/charmbracelet/bubbles/key"

// Common key binding factory functions for consistent key handling across presenters

// NewQuitKey creates a quit key binding (q, ctrl+c)
func NewQuitKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)
}

// NewBackKey creates a back key binding (esc)
func NewBackKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	)
}

// NewHelpKey creates a help toggle key binding (?)
func NewHelpKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
}

// NewUpKey creates an up navigation key binding (↑, k)
func NewUpKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	)
}

// NewDownKey creates a down navigation key binding (↓, j)
func NewDownKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	)
}

// NewEnterKey creates an enter/select key binding
func NewEnterKey() key.Binding {
	return key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	)
}
