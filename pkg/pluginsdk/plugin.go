package pluginsdk

import (
	"context"
	"io"
)

// Plugin is the base interface that ALL plugins must implement.
// It provides basic plugin metadata and declares what capabilities the plugin supports.
type Plugin interface {
	// GetInfo returns basic metadata about the plugin
	GetInfo() PluginInfo

	// GetCapabilities returns a list of capability interface names this plugin implements.
	// Examples: "IEntityProvider", "ICommandProvider", "IEventEmitter"
	// The framework uses this to route requests to the appropriate plugin methods.
	GetCapabilities() []string
}

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
	// Name is the unique identifier for the plugin (e.g., "claude-code", "task-manager")
	Name string

	// Version is the semantic version of the plugin (e.g., "1.0.0")
	Version string

	// Description is a human-readable description of what the plugin does
	Description string

	// IsCore indicates whether this is a built-in plugin shipped with DarwinFlow.
	// Core plugins are loaded automatically, while external plugins are discovered.
	IsCore bool
}

// PluginContext is the runtime context provided to plugins by the framework.
// It provides access to logging, working directory, and event emission.
type PluginContext interface {
	// GetLogger returns a logger for the plugin to use
	GetLogger() Logger

	// GetWorkingDir returns the current working directory of the DarwinFlow project
	GetWorkingDir() string

	// EmitEvent sends an event to the framework's event store.
	// This is the primary way plugins communicate events to the framework.
	EmitEvent(ctx context.Context, event Event) error
}

// CommandContext extends PluginContext with I/O streams for command execution.
// It is provided to commands when they are executed via the CLI.
type CommandContext interface {
	PluginContext

	// GetStdout returns the output stream for the command.
	// Commands should write their output here.
	GetStdout() io.Writer

	// GetStdin returns the input stream for the command.
	// Commands can read user input from here.
	GetStdin() io.Reader
}

// Logger is the interface for plugin logging.
// The framework provides an implementation that plugins use to log messages.
type Logger interface {
	// Debug logs a debug-level message.
	// Debug messages are only shown when debug logging is enabled.
	Debug(msg string, keysAndValues ...interface{})

	// Info logs an info-level message.
	// Info messages are shown during normal operation.
	Info(msg string, keysAndValues ...interface{})

	// Warn logs a warning-level message.
	// Warnings indicate potential issues that don't prevent operation.
	Warn(msg string, keysAndValues ...interface{})

	// Error logs an error-level message.
	// Errors indicate failures that may affect functionality.
	Error(msg string, keysAndValues ...interface{})
}
