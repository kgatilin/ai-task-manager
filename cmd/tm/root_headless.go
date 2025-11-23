//go:build headless
// +build headless

package main

import (
	"errors"

	"github.com/spf13/cobra"
)

// registerTUICommand registers a stub TUI command for headless builds.
// This implementation is used when built with -tags headless.
// When built normally (without headless tag), root_full.go provides the real TUI implementation instead.
func registerTUICommand(rootCmd *cobra.Command, app *App) {
	// Register stub command that returns clear error message
	uiCmd := &cobra.Command{
		Use:   "ui",
		Short: "Interactive terminal UI (not available in headless build)",
		Long:  "The interactive TUI is not available in this headless build. Use CLI commands instead.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("TUI not available in headless build - use CLI commands")
		},
	}
	rootCmd.AddCommand(uiCmd)
}
