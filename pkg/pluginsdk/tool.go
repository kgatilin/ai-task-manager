package pluginsdk

import "context"

// Tool represents a project-scoped tool provided by a plugin.
// Tools are invoked as: dw project <tool-name> [args]
//
// Tools receive a ToolContext which provides access to:
// - Logger for output
// - Database path for persistence
// - Working directory for file operations
// - Output writer for formatted results
//
// Example: dw project session-summary --last
type Tool interface {
	// GetName returns the tool's command name (used as: dw project <name>)
	GetName() string

	// GetDescription returns a brief description of what the tool does
	GetDescription() string

	// GetUsage returns usage instructions (e.g., "analyze [--format=json]")
	GetUsage() string

	// Execute runs the tool with provided arguments
	// The context provides access to plugin services and output
	Execute(ctx context.Context, args []string) error
}
