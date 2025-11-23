//go:build !headless
// +build !headless

package main

import (
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui"
	"github.com/spf13/cobra"
)

// registerTUICommand registers the interactive TUI command.
// This implementation is used for the full build (default).
// When built with -tags headless, root_headless.go provides a stub implementation instead.
func registerTUICommand(rootCmd *cobra.Command, app *App) {
	rootCmd.AddCommand(tui.NewUICommand(app.RepositoryCommon, app.Logger))
}
