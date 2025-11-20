package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates a Cobra command for showing version information
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long: `Display the version of the task manager CLI.

Shows the current version of the tm binary.`,
		Example: `  tm version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "tm - Task Manager CLI v%s\n", version)
		},
	}
}
