package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
	"github.com/spf13/cobra"
)

// NewUICommand creates a Cobra command for launching the interactive TUI
func NewUICommand(
	repo domain.RoadmapRepository,
	logger logger.Logger,
) *cobra.Command {
	return &cobra.Command{
		Use:   "ui",
		Short: "Launch interactive TUI",
		Long: `Launch the interactive terminal user interface (TUI) for the task manager.

The interactive TUI provides a graphical terminal interface for managing
roadmaps, tracks, tasks, and iterations.

Navigation:
  j/k or ↑/↓     Navigate between items
  Enter          View details / drill down
  esc            Go back to previous view
  r              Refresh data
  q              Quit`,
		Example: `  tm ui`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI(cmd.Context(), repo, logger)
		},
	}
}

// runTUI executes the TUI application
func runTUI(
	ctx context.Context,
	repo domain.RoadmapRepository,
	logger logger.Logger,
) error {
	// Create the TUI app model
	appModel := NewAppModelNew(ctx, repo, logger)

	// Start the Bubble Tea program
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
