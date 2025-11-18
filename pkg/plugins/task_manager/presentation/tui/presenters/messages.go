package presenters

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// IterationSelectedMsg is sent when a user selects an iteration on the dashboard
type IterationSelectedMsg struct {
	IterationNumber int
	SelectedIndex   int // Dashboard selected index (for restoring focus on return)
}

// TrackSelectedMsg is sent when a user selects a track on the dashboard
type TrackSelectedMsg struct {
	TrackID       string
	SelectedIndex int // Dashboard selected index (for restoring focus on return)
}

// TaskSelectedMsg is sent when a user selects a task
type TaskSelectedMsg struct {
	TaskID        string
	SelectedIndex int // Dashboard selected index (for restoring focus on return)
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

// TaskTransitionCompletedMsg is sent after a successful task status transition
type TaskTransitionCompletedMsg struct {
	ActiveTab     IterationDetailTab // Preserve active tab (Tasks=0, ACs=1)
	SelectedIndex int                // Preserve selected index across reload
}

// ReorderCompletedMsg is sent after iterations are successfully reordered
type ReorderCompletedMsg struct {
	SelectedIterationNumber int
}

// RefreshDashboardMsg is sent when user requests dashboard refresh (r key)
type RefreshDashboardMsg struct {
	SelectedIndex int // Preserve selected index across reload
}

// DocumentLoadedMsg is sent when a document has been loaded from repository
type DocumentLoadedMsg struct {
	ViewModel *viewmodels.DocumentViewModel
	Error     error
}

// DocumentActionCompletedMsg is sent after a successful document action (approve/disapprove)
type DocumentActionCompletedMsg struct{}

// DrillIntoDocumentMsg is sent when a user navigates to view a document
type DrillIntoDocumentMsg struct {
	DocumentID string
}

// Ensure these are valid Bubble Tea messages
var (
	_ tea.Msg = IterationSelectedMsg{}
	_ tea.Msg = TrackSelectedMsg{}
	_ tea.Msg = TaskSelectedMsg{}
	_ tea.Msg = ErrorMsg{}
	_ tea.Msg = ACActionCompletedMsg{}
	_ tea.Msg = TaskTransitionCompletedMsg{}
	_ tea.Msg = ReorderCompletedMsg{}
	_ tea.Msg = RefreshDashboardMsg{}
	_ tea.Msg = DocumentLoadedMsg{}
	_ tea.Msg = DocumentActionCompletedMsg{}
	_ tea.Msg = DrillIntoDocumentMsg{}
	_ tea.Msg = BackMsgNew{}
)
