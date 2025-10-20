package domain

import (
	"context"
	"errors"
)

// Common errors
var (
	ErrNotFound = errors.New("entity not found")
)

// Plugin represents a plugin that provides entities to the system.
// Plugins can be built-in (core) or external (loaded from project directory).
type Plugin interface {
	// GetInfo returns metadata about this plugin
	GetInfo() PluginInfo

	// GetEntityTypes returns the types of entities this plugin provides
	GetEntityTypes() []EntityTypeInfo

	// Query returns entities matching the given query
	Query(ctx context.Context, query EntityQuery) ([]IExtensible, error)

	// GetEntity retrieves a single entity by ID
	GetEntity(ctx context.Context, entityID string) (IExtensible, error)

	// UpdateEntity updates an entity's fields
	UpdateEntity(ctx context.Context, entityID string, fields map[string]interface{}) (IExtensible, error)
}

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
	// Name is the unique identifier for this plugin (e.g., "claude-code")
	Name string

	// Version is the plugin version
	Version string

	// Description is a human-readable description
	Description string

	// IsCore indicates if this is a built-in core plugin
	IsCore bool
}

// EntityTypeInfo describes an entity type provided by a plugin
type EntityTypeInfo struct {
	// Type is the entity type name (e.g., "session", "task")
	Type string

	// DisplayName is the human-readable name (e.g., "Claude Session")
	DisplayName string

	// DisplayNamePlural is the plural form (e.g., "Claude Sessions")
	DisplayNamePlural string

	// Capabilities lists the capability interfaces this entity type supports
	Capabilities []string

	// Icon is an optional emoji or symbol for UI display
	Icon string
}

// EntityQuery specifies criteria for querying entities
type EntityQuery struct {
	// EntityType filters by entity type (empty = all types from plugin)
	EntityType string

	// Capabilities filters to entities supporting ALL listed capabilities
	Capabilities []string

	// Filters are field-level filters (field name -> expected value)
	Filters map[string]interface{}

	// Limit restricts the number of results (0 = no limit)
	Limit int

	// Offset for pagination
	Offset int

	// OrderBy specifies sort field (empty = plugin default order)
	OrderBy string

	// OrderDesc reverses sort order
	OrderDesc bool
}

// Command represents a CLI command provided by a plugin (e.g., dw <plugin-name> <command>)
type Command interface {
	// GetName returns the command name (e.g., "init", "log")
	GetName() string

	// GetDescription returns a brief description of what the command does
	GetDescription() string

	// GetUsage returns usage instructions (e.g., "init [--force]")
	GetUsage() string

	// Execute runs the command with provided arguments
	Execute(ctx context.Context, args []string) error
}

// Tool represents a project-scoped tool provided by a plugin (e.g., dw project <tool-name>)
type Tool interface {
	// GetName returns the tool's command name (used as: dw project <name>)
	GetName() string

	// GetDescription returns a brief description of what the tool does
	GetDescription() string

	// GetUsage returns usage instructions (e.g., "analyze [--format=json]")
	GetUsage() string

	// Execute runs the tool with provided arguments
	// Note: Tools receive PluginContext from app layer (not defined in domain)
	Execute(ctx context.Context, args []string) error
}

// ICommandProvider is a capability that plugins can implement to provide CLI commands
type ICommandProvider interface {
	// GetCommands returns the commands provided by this plugin
	GetCommands() []Command
}

// IToolProvider is a capability that plugins can implement to provide project tools
type IToolProvider interface {
	// GetTools returns the tools provided by this plugin
	GetTools() []Tool
}
