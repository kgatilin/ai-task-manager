package pluginsdk

import "context"

// Command represents a CLI command provided by a plugin.
// Commands are invoked as: dw <plugin-name> <command-name> [args]
//
// Example: dw claude-code init --force
type Command interface {
	// GetName returns the command name (e.g., "init", "log")
	GetName() string

	// GetDescription returns a brief description of what the command does
	GetDescription() string

	// GetUsage returns usage instructions (e.g., "init [--force]")
	GetUsage() string

	// Execute runs the command with provided arguments
	// The context provides access to plugin services and I/O
	Execute(ctx context.Context, args []string) error
}
