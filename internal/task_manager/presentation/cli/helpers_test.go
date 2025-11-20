package cli_test

import "github.com/spf13/cobra"

// findCommand finds a subcommand by name within a parent command
func findCommand(parent *cobra.Command, name string) *cobra.Command {
	for _, cmd := range parent.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
