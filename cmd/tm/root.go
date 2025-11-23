package main

import (
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

// NewRootCmd creates the root Cobra command for the task manager.
func NewRootCmd(app *App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "tm",
		Short:   "Task Manager - Manage tasks, iterations, and roadmaps",
		Long:    "tm is a command-line tool for managing tasks, iterations, tracks, and roadmaps",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show help if no subcommand is provided
			return cmd.Help()
		},
	}

	// Set up cobra configuration
	rootCmd.SilenceUsage = true

	// Add global flags
	rootCmd.PersistentFlags().String("project", "", "Specify project name")

	// Add special commands (version, ui, completion, prompt)
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewCompletionCommand())
	rootCmd.AddCommand(cli.NewPromptCommand(cli.GetSystemPrompt))

	// Add application commands from the Cobra command groups
	if app != nil {
		// Register TUI command (implementation varies by build tag)
		registerTUICommand(rootCmd, app)

		// Add project commands
		rootCmd.AddCommand(cli.NewProjectCommands(app.ProjectService))

		// Add task commands from the Cobra command group
		rootCmd.AddCommand(cli.NewTaskCommands(app.TaskService, app.ACService))

		// Add iteration commands from the Cobra command group
		rootCmd.AddCommand(cli.NewIterationCommands(app.IterationService, app.DocumentService, app.ACService))

		// Add AC commands from the Cobra command group
		rootCmd.AddCommand(cli.NewACCommands(app.ACService, app.TaskService))

		// Add track commands from the Cobra command group
		rootCmd.AddCommand(cli.NewTrackCommands(app.TrackService, app.DocumentService))

		// Add ADR commands from the Cobra command group
		rootCmd.AddCommand(cli.NewADRCommands(app.ADRService))

		// Add roadmap commands from the Cobra command group
		rootCmd.AddCommand(cli.NewRoadmapCommands(app.RoadmapService))

		// Add document commands from the Cobra command group
		rootCmd.AddCommand(cli.NewDocCommands(app.DocumentService))
	}

	return rootCmd
}
