package presenters

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/components"
	"github.com/muesli/reflow/wordwrap"
)

// FeedbackInputComponent handles failure feedback input for acceptance criteria.
// This component is shared between IterationDetailPresenter and TaskDetailPresenter
// to ensure consistent UX when failing ACs.
type FeedbackInputComponent struct {
	inFeedbackMode bool
	input          textinput.Model
	failingACID    string
}

// NewFeedbackInputComponent creates a new feedback input component
func NewFeedbackInputComponent() *FeedbackInputComponent {
	ti := textinput.New()
	ti.Placeholder = "Enter failure reason..."
	ti.CharLimit = 500

	return &FeedbackInputComponent{
		input: ti,
	}
}

// StartFeedback enters feedback mode for the given AC
func (c *FeedbackInputComponent) StartFeedback(acID string) tea.Cmd {
	c.inFeedbackMode = true
	c.failingACID = acID
	c.input.Focus()
	c.input.SetValue("")
	return textinput.Blink
}

// CancelFeedback exits feedback mode without submitting
func (c *FeedbackInputComponent) CancelFeedback() {
	c.inFeedbackMode = false
	c.failingACID = ""
	c.input.SetValue("")
	c.input.Blur()
}

// SubmitFeedback returns the AC ID and feedback, then resets the component.
// If feedback is empty, returns "Failed via TUI" as default.
func (c *FeedbackInputComponent) SubmitFeedback() (acID string, feedback string) {
	feedback = c.input.Value()
	if feedback == "" {
		feedback = "Failed via TUI"
	}
	acID = c.failingACID
	c.CancelFeedback()
	return acID, feedback
}

// Update handles keyboard input when in feedback mode.
// Returns true if the message was handled by this component, false otherwise.
// On Enter key, caller should handle submission by calling SubmitFeedback().
func (c *FeedbackInputComponent) Update(msg tea.Msg) (handled bool, cmd tea.Cmd) {
	if !c.inFeedbackMode {
		return false, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			c.CancelFeedback()
			return true, nil
		case tea.KeyEnter:
			// Caller should handle submission (call SubmitFeedback() and perform AC fail action)
			return true, nil
		default:
			// Pass typing to textinput
			var cmd tea.Cmd
			c.input, cmd = c.input.Update(msg)
			return true, cmd
		}
	}

	return false, nil
}

// View renders the feedback input UI inline at the bottom of the view.
// Returns empty string if not in feedback mode.
func (c *FeedbackInputComponent) View(width int) string {
	if !c.inFeedbackMode {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString(components.Styles.SectionStyle.Render("Failure Reason"))
	b.WriteString("\n")

	// Apply width constraint to input for text wrapping
	availableWidth := width - 4
	if availableWidth < 40 {
		availableWidth = 40
	}
	// Wrap the input text to constrain it within available width
	inputText := c.input.Value()
	wrappedInput := wordwrap.String(inputText, availableWidth)
	b.WriteString(wrappedInput)
	if wrappedInput != "" && !strings.HasSuffix(wrappedInput, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(components.Styles.MetadataStyle.Render("Press Enter to submit or ESC to cancel"))

	return b.String()
}

// IsActive returns whether feedback mode is currently active
func (c *FeedbackInputComponent) IsActive() bool {
	return c.inFeedbackMode
}
