package presenters

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/components"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// LoadingPresenter presents a loading view with a spinner and message
type LoadingPresenter struct {
	viewModel *viewmodels.LoadingViewModel
	spinner   components.Spinner
}

// NewLoadingPresenter creates a new loading presenter
func NewLoadingPresenter(vm *viewmodels.LoadingViewModel) *LoadingPresenter {
	return &LoadingPresenter{
		viewModel: vm,
		spinner:   components.NewSpinner(),
	}
}

func (p *LoadingPresenter) Init() tea.Cmd {
	if p.viewModel.ShowSpinner {
		return p.spinner.Tick()
	}
	return nil
}

func (p *LoadingPresenter) Update(msg tea.Msg) (Presenter, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return p, tea.Quit
		}
	case interface{}: // Handles all Bubble Tea messages including spinner.TickMsg
		if p.viewModel.ShowSpinner {
			return p, p.spinner.Update(msg)
		}
	}
	return p, nil
}

func (p *LoadingPresenter) View() string {
	var output string
	if p.viewModel.ShowSpinner {
		output = p.spinner.View() + " "
	}
	output += components.Styles.LoadingStyle.Render(p.viewModel.Message)
	return "\n" + output + "\n\n"
}
