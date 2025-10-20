package pluginsdk

import "io"

// PluginContext provides services to plugins without exposing internal types.
// This is the primary way plugins access system capabilities.
type PluginContext interface {
	// GetLogger returns a logger for the plugin
	GetLogger() Logger

	// GetDBPath returns the database path
	GetDBPath() string

	// GetWorkingDir returns the current working directory
	GetWorkingDir() string
}

// CommandContext provides context for command execution.
// This extends PluginContext with I/O capabilities for interactive commands.
type CommandContext interface {
	PluginContext

	// GetOutput returns the output writer for command output
	GetOutput() io.Writer

	// GetInput returns the input reader for command input
	GetInput() io.Reader
}

// ToolContext provides context for tool execution.
// Tools are project-scoped operations (dw project <tool>)
type ToolContext interface {
	PluginContext

	// GetOutput returns the output writer for tool output
	GetOutput() io.Writer
}
