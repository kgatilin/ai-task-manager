package cli

import (
	"context"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
)

// PluginProvider is an interface for accessing plugin functionality
// This allows the TUI command to work with the plugin without importing the root package
type PluginProvider interface {
	GetRepositoryForProject(project string) (domain.RoadmapRepository, func(), error)
	GetLogger() pluginsdk.Logger
	GetActiveProject() (string, error)
}

// TUICommand launches the interactive TUI for browsing roadmaps and tracks
type TUICommand struct {
	Plugin  PluginProvider
	project string
}

// GetName returns the command name
func (c *TUICommand) GetName() string {
	return "tui"
}

// GetDescription returns the command description
func (c *TUICommand) GetDescription() string {
	return "Launch interactive TUI to browse roadmaps and tracks"
}

// GetHelp returns the command help text
func (c *TUICommand) GetHelp() string {
	return `Usage: dw task-manager tui

Launch an interactive terminal user interface to browse and manage roadmaps, tracks, and tasks.

Navigation:
  j/k or ↑/↓     Navigate between items
  Enter          View details / drill down
  esc            Go back to previous view
  r              Refresh data
  q              Quit

Views:
  Roadmap List   Browse all tracks in the active roadmap
  Track Detail   View tasks within a selected track
`
}

// GetUsage returns the command usage
func (c *TUICommand) GetUsage() string {
	return "tui"
}

// Execute runs the TUI command
func (c *TUICommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
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

	// Create the TUI model with project name
	appModel := NewAppModelWithProject(ctx, repo, c.Plugin.GetLogger(), projectName)

	// Start the Bubble Tea program
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}

// TUIScaffold provides utility functions for TUI-related operations
type TUIScaffold struct{}

// RunProgram is a testable wrapper around tea.NewProgram
// This allows us to mock the Bubble Tea program in tests
var RunProgram = func(model tea.Model) error {
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// StartTUI starts the task manager TUI with the given context and output writer
// This function is exposed for external callers (e.g., other plugins, CLI wrappers)
func StartTUI(ctx context.Context, repository domain.RoadmapRepository, logger pluginsdk.Logger, output io.Writer) error {
	if repository == nil {
		return fmt.Errorf("repository is required")
	}

	appModel := NewAppModel(ctx, repository, logger)

	p := tea.NewProgram(appModel, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(output, "TUI error: %v\n", err)
		return err
	}

	return nil
}
