package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/infrastructure/cli"
)

// PluginProvider is an alias for the infrastructure provider interface
type PluginProvider = cli.PluginProvider

// TUINewCommand launches the new MVP TUI for task manager
type TUINewCommand struct {
	Plugin  PluginProvider
	project string
}

func (c *TUINewCommand) GetName() string {
	return "tui-new"
}

func (c *TUINewCommand) GetDescription() string {
	return "Launch new MVP TUI (Dashboard → Iteration → Task navigation)"
}

func (c *TUINewCommand) GetHelp() string {
	return `Usage: dw task-manager tui-new [--project <name>]

Launch the new MVP terminal user interface with core navigation flow:
- Dashboard: View all iterations
- Iteration Details: View tasks in an iteration
- Task Details: View acceptance criteria for a task

Navigation:
  j/k or ↑/↓     Navigate between items
  Enter          View details / drill down
  esc            Go back to previous view
  r              Refresh data
  q              Quit

Flags:
  --project <name>    Use specific project (overrides active project)
`
}

func (c *TUINewCommand) GetUsage() string {
	return "tui-new [--project <name>]"
}

func (c *TUINewCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		}
	}

	// Get repository for project
	repo, cleanup, err := c.Plugin.GetRepositoryForProject(c.project)
	if err != nil {
		return err
	}
	defer cleanup()

	// Determine project name for display
	projectName := c.project
	if projectName == "" {
		projectName, err = c.Plugin.GetActiveProject()
		if err != nil {
			projectName = "default"
		}
	}

	// Create the TUI app model
	appModel := NewAppModelNew(ctx, repo, c.Plugin.GetLogger(), projectName)

	// Start the Bubble Tea program
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
