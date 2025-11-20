package presenters

// ACListComponent provides reusable acceptance criteria list management for TUI presenters.
//
// This component eliminates ~100-120 lines of duplication between IterationDetailPresenter
// and TaskDetailPresenter by centralizing:
// 1. AC action methods (verify, skip, fail) - ~60 lines
// 2. AC rendering logic - ~40 lines
// 3. Feedback input delegation to FeedbackInputComponent
//
// Architecture:
// - Stateless rendering: RenderACList accepts ACs and selectedIndex as parameters
// - Stateful actions: Repository operations encapsulated in VerifyAC/SkipAC/FailAC
// - Interface-based: Works with both ACDetailViewModel and IterationACViewModel via ACViewModel interface
// - Delegate feedback: Uses FeedbackInputComponent for failure reason input
//
// Usage in TaskDetailPresenter:
//   component := NewACListComponent(repo, ctx, true)  // enableExpand=true for testing instructions
//   acVMs := WrapACDetailViewModels(presenter.viewModel.AcceptanceCriteria)
//   component.RenderACList(&b, acVMs, selectedIndex, width)
//
// Usage in IterationDetailPresenter:
//   component := NewACListComponent(repo, ctx, true)  // enableExpand=true (same as TaskDetail)
//   acVMs := WrapIterationACViewModels(presenter.viewModel.AllAcceptanceCriteria)
//   component.RenderACList(&b, acVMs, selectedIndex, width)
//
// Key Design Decisions:
// - ACViewModel interface allows working with different ViewModel types
// - Wrappers (ACDetailViewModelWrapper, IterationACViewModelWrapper) adapt ViewModels to interface
// - Component returns ACActionCompletedMsg with activeTab and selectedIndex to preserve state
// - Presenters handle grouping logic (e.g., IterationDetail groups by task)
// - Component only renders flat lists (grouping is presenter's responsibility)

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/components"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
	"github.com/muesli/reflow/wordwrap"
)

// ACViewModel is an interface satisfied by both ACDetailViewModel and IterationACViewModel.
// This allows the component to work with AC view models from different presenters.
type ACViewModel interface {
	GetID() string
	GetDescription() string
	GetStatus() string
	GetStatusIcon() string
	GetTestingInstructions() string
	GetNotes() string
	GetIsExpanded() bool
	SetIsExpanded(bool)
	GetStatusColor() string
}

// Ensure ACDetailViewModel implements ACViewModel
var _ ACViewModel = (*ACDetailViewModelWrapper)(nil)

// ACDetailViewModelWrapper wraps viewmodels.ACDetailViewModel to implement ACViewModel
type ACDetailViewModelWrapper struct {
	*viewmodels.ACDetailViewModel
}

func (w *ACDetailViewModelWrapper) GetID() string                  { return w.ID }
func (w *ACDetailViewModelWrapper) GetDescription() string         { return w.Description }
func (w *ACDetailViewModelWrapper) GetStatus() string              { return w.Status }
func (w *ACDetailViewModelWrapper) GetStatusIcon() string          { return w.StatusIcon }
func (w *ACDetailViewModelWrapper) GetTestingInstructions() string { return w.TestingInstructions }
func (w *ACDetailViewModelWrapper) GetNotes() string               { return w.Notes }
func (w *ACDetailViewModelWrapper) GetIsExpanded() bool            { return w.IsExpanded }
func (w *ACDetailViewModelWrapper) SetIsExpanded(expanded bool)    { w.IsExpanded = expanded }
func (w *ACDetailViewModelWrapper) GetStatusColor() string         { return w.StatusColor }

// IterationACViewModelWrapper wraps viewmodels.IterationACViewModel to implement ACViewModel
type IterationACViewModelWrapper struct {
	*viewmodels.IterationACViewModel
}

func (w *IterationACViewModelWrapper) GetID() string                  { return w.ID }
func (w *IterationACViewModelWrapper) GetDescription() string         { return w.Description }
func (w *IterationACViewModelWrapper) GetStatus() string              { return w.Status }
func (w *IterationACViewModelWrapper) GetStatusIcon() string          { return w.StatusIcon }
func (w *IterationACViewModelWrapper) GetTestingInstructions() string { return w.TestingInstructions }
func (w *IterationACViewModelWrapper) GetNotes() string               { return w.Notes }
func (w *IterationACViewModelWrapper) GetIsExpanded() bool            { return w.IsExpanded }
func (w *IterationACViewModelWrapper) SetIsExpanded(expanded bool)    { w.IsExpanded = expanded }
func (w *IterationACViewModelWrapper) GetStatusColor() string         { return w.StatusColor }

// getACStyleForStatus returns the appropriate style for an AC based on its status color
func getACStyleForStatus(statusColor string) lipgloss.Style {
	switch statusColor {
	case "failed":
		return components.Styles.ACFailedStyle
	case "success":
		return components.Styles.ACVerifiedStyle
	case "warning":
		return components.Styles.ACPendingStyle
	case "skipped":
		return components.Styles.ACSkippedStyle
	default:
		return lipgloss.NewStyle()
	}
}

// ACListComponent handles rendering and actions for acceptance criteria lists.
// This component is shared between IterationDetailPresenter and TaskDetailPresenter
// to eliminate duplication in AC management.
//
// Design:
// - Stateless for rendering: accepts ACs and selectedIndex as parameters
// - Stateful for actions: encapsulates repository operations (verify/skip/fail)
// - Delegates feedback input to FeedbackInputComponent
// - Preserves selection after actions via ACActionCompletedMsg
type ACListComponent struct {
	repo          domain.RoadmapRepository
	ctx           context.Context
	feedbackInput *FeedbackInputComponent
	enableExpand  bool // Whether to allow expanding ACs (TaskDetail only)
}

// NewACListComponent creates a new AC list component
func NewACListComponent(repo domain.RoadmapRepository, ctx context.Context, enableExpand bool) *ACListComponent {
	return &ACListComponent{
		repo:          repo,
		ctx:           ctx,
		feedbackInput: NewFeedbackInputComponent(),
		enableExpand:  enableExpand,
	}
}

// VerifyAC marks an AC as verified.
// Returns ACActionCompletedMsg to preserve selection and active tab.
func (c *ACListComponent) VerifyAC(acID string, activeTab IterationDetailTab, currentSelectedIndex int) tea.Cmd {
	return func() tea.Msg {
		// Fetch AC
		ac, err := c.repo.GetAC(c.ctx, acID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		// Update status to verified
		ac.Status = entities.ACStatusVerified
		ac.UpdatedAt = time.Now()

		// Save
		err = c.repo.UpdateAC(c.ctx, ac)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return ACActionCompletedMsg{ActiveTab: activeTab, SelectedIndex: currentSelectedIndex}
	}
}

// SkipAC marks an AC as skipped.
// Returns ACActionCompletedMsg to preserve selection and active tab.
func (c *ACListComponent) SkipAC(acID string, activeTab IterationDetailTab, currentSelectedIndex int) tea.Cmd {
	return func() tea.Msg {
		// Fetch AC
		ac, err := c.repo.GetAC(c.ctx, acID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		// Update status to skipped
		ac.Status = entities.ACStatusSkipped
		ac.Notes = "Skipped via TUI"
		ac.UpdatedAt = time.Now()

		// Save
		err = c.repo.UpdateAC(c.ctx, ac)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return ACActionCompletedMsg{ActiveTab: activeTab, SelectedIndex: currentSelectedIndex}
	}
}

// FailAC marks an AC as failed with feedback.
// Returns ACActionCompletedMsg to preserve selection and active tab.
func (c *ACListComponent) FailAC(acID, feedback string, activeTab IterationDetailTab, currentSelectedIndex int) tea.Cmd {
	return func() tea.Msg {
		// Fetch AC
		ac, err := c.repo.GetAC(c.ctx, acID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		// Update status to failed with feedback
		ac.Status = entities.ACStatusFailed
		ac.Notes = feedback
		ac.UpdatedAt = time.Now()

		// Save
		err = c.repo.UpdateAC(c.ctx, ac)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return ACActionCompletedMsg{ActiveTab: activeTab, SelectedIndex: currentSelectedIndex}
	}
}

// StartFeedback enters feedback mode for the given AC
func (c *ACListComponent) StartFeedback(acID string) tea.Cmd {
	return c.feedbackInput.StartFeedback(acID)
}

// UpdateFeedback delegates to FeedbackInputComponent.
// Returns true if the message was handled by the feedback component.
func (c *ACListComponent) UpdateFeedback(msg tea.Msg) (handled bool, cmd tea.Cmd) {
	return c.feedbackInput.Update(msg)
}

// ViewFeedback delegates to FeedbackInputComponent.
// Returns the feedback input UI or empty string if not active.
func (c *ACListComponent) ViewFeedback(width int) string {
	return c.feedbackInput.View(width)
}

// IsFeedbackActive returns whether feedback mode is currently active
func (c *ACListComponent) IsFeedbackActive() bool {
	return c.feedbackInput.IsActive()
}

// SubmitFeedback returns the AC ID and feedback, then resets the component.
func (c *ACListComponent) SubmitFeedback() (acID string, feedback string) {
	return c.feedbackInput.SubmitFeedback()
}

// RenderACList renders a flat list of acceptance criteria.
// The presenter is responsible for grouping (IterationDetail groups by task).
// This method accepts:
// - acViewModels: The ACs to render (already flattened/filtered by presenter)
// - selectedIndex: The currently selected AC index
// - availableWidth: Width for text wrapping
//
// Features:
// - Status icon + ID + description
// - Highlight selected AC
// - Show testing instructions if expanded (when enableExpand=true)
// - Show notes if present
// - Visual indicator (📋) for ACs with testing instructions
func (c *ACListComponent) RenderACList(b *strings.Builder, acViewModels []ACViewModel, selectedIndex int, availableWidth int) {
	for i, ac := range acViewModels {
		// AC header with icon and description
		hasInstructions := ""
		if ac.GetTestingInstructions() != "" && c.enableExpand {
			hasInstructions = " 📋"
		}

		headerText := fmt.Sprintf("  %s %s: %s%s", ac.GetStatusIcon(), ac.GetID(), ac.GetDescription(), hasInstructions)
		wrappedHeaderText := lipgloss.NewStyle().Width(availableWidth).Render(headerText)
		if i == selectedIndex {
			b.WriteString(components.Styles.SelectedStyle.Render(wrappedHeaderText))
		} else {
			// Apply status-based styling
			statusStyle := getACStyleForStatus(ac.GetStatusColor())
			b.WriteString(statusStyle.Render(wrappedHeaderText))
		}
		b.WriteString("\n")

		// If expanded and has testing instructions, show them (only when enableExpand=true)
		if c.enableExpand && ac.GetIsExpanded() && ac.GetTestingInstructions() != "" {
			b.WriteString(components.Styles.TestingStyle.Render("    Testing Instructions:"))
			b.WriteString("\n")
			// Format testing instructions with proper indentation and wrapping
			instructions := strings.Split(ac.GetTestingInstructions(), "\n")
			for _, line := range instructions {
				if line != "" {
					wrappedLine := lipgloss.NewStyle().Width(availableWidth - 6).Render(line)
					b.WriteString(components.Styles.TestingStyle.Render("      " + wrappedLine))
					b.WriteString("\n")
				}
			}
		}

		// Show notes/failure reason if present (skip for verified ACs)
		if ac.GetNotes() != "" &&
			ac.GetStatus() != string(entities.ACStatusVerified) &&
			ac.GetStatus() != string(entities.ACStatusAutomaticallyVerified) {
			// Use "Failure Reason:" label for failed ACs, otherwise "Notes:"
			label := "Notes:"
			style := components.Styles.TestingStyle
			if ac.GetStatusColor() == "failed" {
				label = "Failure Reason:"
				style = components.Styles.ACFailedStyle
			}

			// Calculate available width for notes (total width - indentation - label length - margins)
			availableNoteWidth := availableWidth - 4 - len(label) - 2
			if availableNoteWidth < 20 {
				availableNoteWidth = 20 // Minimum width
			}
			// Wrap the notes text
			wrappedNotes := wordwrap.String(ac.GetNotes(), availableNoteWidth)
			// Split into lines and render with indentation
			notesLines := strings.Split(wrappedNotes, "\n")
			for i, line := range notesLines {
				if i == 0 {
					// First line with label
					notesLine := fmt.Sprintf("    %s %s", label, line)
					b.WriteString(style.Render(notesLine))
				} else {
					// Subsequent lines with matching indentation
					indent := strings.Repeat(" ", 4+len(label)+1)
					notesLine := fmt.Sprintf("%s%s", indent, line)
					b.WriteString(style.Render(notesLine))
				}
				b.WriteString("\n")
			}
		}
	}
}

// WrapACDetailViewModels wraps ACDetailViewModel slice into ACViewModel interface slice
func WrapACDetailViewModels(acs []*viewmodels.ACDetailViewModel) []ACViewModel {
	wrapped := make([]ACViewModel, len(acs))
	for i, ac := range acs {
		wrapped[i] = &ACDetailViewModelWrapper{ACDetailViewModel: ac}
	}
	return wrapped
}

// WrapIterationACViewModels wraps IterationACViewModel slice into ACViewModel interface slice
func WrapIterationACViewModels(acs []*viewmodels.IterationACViewModel) []ACViewModel {
	wrapped := make([]ACViewModel, len(acs))
	for i, ac := range acs {
		wrapped[i] = &IterationACViewModelWrapper{IterationACViewModel: ac}
	}
	return wrapped
}
