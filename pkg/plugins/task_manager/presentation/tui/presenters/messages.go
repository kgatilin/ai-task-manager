package presenters

import tea "github.com/charmbracelet/bubbletea"

// IterationSelectedMsg is sent when a user selects an iteration on the dashboard
type IterationSelectedMsg struct {
	IterationNumber int
}

// TaskSelectedMsg is sent when a user selects a task
type TaskSelectedMsg struct {
	TaskID string
}

// ErrorMsg is sent when an error occurs during loading or operations
type ErrorMsg struct {
	Err error
}

// ACActionCompletedMsg is sent after a successful AC action (verify/skip/fail)
type ACActionCompletedMsg struct {
	ActiveTab     IterationDetailTab // Preserve active tab (Tasks=0, ACs=1)
	SelectedIndex int                // Preserve selected index across reload
}

// ReorderCompletedMsg is sent after iterations are successfully reordered
type ReorderCompletedMsg struct {
	SelectedIterationNumber int
}

// Ensure these are valid Bubble Tea messages
var (
	_ tea.Msg = IterationSelectedMsg{}
	_ tea.Msg = TaskSelectedMsg{}
	_ tea.Msg = ErrorMsg{}
	_ tea.Msg = ACActionCompletedMsg{}
	_ tea.Msg = ReorderCompletedMsg{}
)
