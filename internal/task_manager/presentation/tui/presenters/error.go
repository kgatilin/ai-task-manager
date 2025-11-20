package presenters

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/components"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
	"github.com/muesli/reflow/wordwrap"
)

// ErrorPresenter presents an error view with message and help text
type ErrorPresenter struct {
	viewModel *viewmodels.ErrorViewModel
	help      components.Help
	quitKey   key.Binding
	backKey   key.Binding
	width     int
}

// NewErrorPresenter creates a new error presenter
func NewErrorPresenter(vm *viewmodels.ErrorViewModel) *ErrorPresenter {
	return &ErrorPresenter{
		viewModel: vm,
		help:      components.NewHelp(),
		quitKey:   components.NewQuitKey(),
		backKey:   components.NewBackKey(),
		width:     80, // default width
	}
}

func (p *ErrorPresenter) Init() tea.Cmd {
	return nil
}

func (p *ErrorPresenter) Update(msg tea.Msg) (Presenter, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.quitKey):
			return p, tea.Quit
		case key.Matches(msg, p.backKey):
			if p.viewModel.CanGoBack {
				return p, func() tea.Msg { return BackMsgNew{} }
			}
		}
	}
	return p, nil
}

func (p *ErrorPresenter) View() string {
	var b strings.Builder

	// Calculate available width (leave margin)
	availableWidth := p.width - 4
	if availableWidth < 40 {
		availableWidth = 40
	}

	b.WriteString("\n")
	b.WriteString(components.Styles.ErrorTitleStyle.Render("Error"))
	b.WriteString("\n\n")

	// Wrap error message
	wrappedMessage := wordwrap.String(p.viewModel.ErrorMessage, availableWidth)
	b.WriteString(components.Styles.ErrorMessageStyle.Render(wrappedMessage))
	b.WriteString("\n")

	if p.viewModel.Details != "" {
		b.WriteString("\n")
		wrappedDetails := wordwrap.String(p.viewModel.Details, availableWidth)
		b.WriteString(components.Styles.ErrorDetailsStyle.Render(wrappedDetails))
		b.WriteString("\n")
	}

	if p.viewModel.RetryAction != "" {
		b.WriteString("\n")
		wrappedRetry := wordwrap.String("Suggestion: "+p.viewModel.RetryAction, availableWidth)
		b.WriteString(components.Styles.ErrorDetailsStyle.Render(wrappedRetry))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	keys := []key.Binding{}
	if p.viewModel.CanGoBack {
		keys = append(keys, p.backKey)
	}
	keys = append(keys, p.quitKey)
	b.WriteString(p.help.ShortHelpView(keys))
	b.WriteString("\n")

	return b.String()
}
