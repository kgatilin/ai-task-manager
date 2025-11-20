package cli

import "github.com/spf13/cobra"

// FindCommandByName searches for a subcommand by name within a parent command.
// Returns nil if no matching subcommand is found.
func FindCommandByName(parent *cobra.Command, name string) *cobra.Command {
	for _, cmd := range parent.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
