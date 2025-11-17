package presenters

import tea "github.com/charmbracelet/bubbletea"

// Presenter is the base interface for all presenters in the TUI
type Presenter interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Presenter, tea.Cmd)
	View() string
}

// BackMsgNew is sent when the user wants to go back in the TUI
type BackMsgNew struct{}
