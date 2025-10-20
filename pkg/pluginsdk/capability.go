package pluginsdk

import "context"

// IEntityProvider is a plugin capability for providing queryable entities.
// Plugins that implement this can be queried for entities via the framework's registry.
type IEntityProvider interface {
	Plugin

	// GetEntityTypes returns metadata about all entity types this plugin provides
	GetEntityTypes() []EntityTypeInfo

	// Query returns entities matching the given query criteria
	Query(ctx context.Context, query EntityQuery) ([]IExtensible, error)

	// GetEntity retrieves a specific entity by ID
	GetEntity(ctx context.Context, entityID string) (IExtensible, error)
}

// IEntityUpdater is a plugin capability for supporting entity updates.
// It extends IEntityProvider with the ability to modify entities.
type IEntityUpdater interface {
	IEntityProvider

	// UpdateEntity modifies an entity's fields and returns the updated entity.
	// The fields map contains field names as keys and new values.
	UpdateEntity(ctx context.Context, entityID string, fields map[string]interface{}) (IExtensible, error)
}

// ICommandProvider is a plugin capability for providing CLI commands.
// Plugins that implement this can register commands accessible via `dw project <command>`.
type ICommandProvider interface {
	Plugin

	// GetCommands returns all commands provided by this plugin
	GetCommands() []Command
}

// IEventEmitter is a plugin capability for emitting real-time events.
// Plugins that implement this can stream events to the framework's event store.
type IEventEmitter interface {
	Plugin
	// Event emission is handled via PluginContext.EmitEvent() for built-in plugins
	// or stdout JSON streams for subprocess plugins.
}

// EntityTypeInfo describes an entity type provided by a plugin
type EntityTypeInfo struct {
	// Type is the unique identifier for this entity type (e.g., "session", "task")
	Type string

	// DisplayName is the human-readable singular name (e.g., "Claude Session", "Task")
	DisplayName string

	// DisplayNamePlural is the human-readable plural name (e.g., "Claude Sessions", "Tasks")
	DisplayNamePlural string

	// Capabilities is a list of entity capability interfaces this type implements
	// Examples: ["IExtensible", "ITrackable"]
	Capabilities []string

	// Icon is an optional emoji or symbol representing this entity type.
	// Used in UI displays.
	Icon string

	// Description is a human-readable description of this entity type
	Description string
}

// EntityQuery represents a query for entities from a plugin.
// Plugins receive this query and return matching entities.
type EntityQuery struct {
	// EntityType is the type of entities to query (e.g., "session", "task")
	EntityType string

	// Filters contains query filters as key-value pairs.
	// The structure and supported filters depend on the plugin and entity type.
	// Common filters: "status", "created_after", "tag", etc.
	Filters map[string]interface{}

	// Limit is the maximum number of entities to return.
	// 0 means no limit.
	Limit int

	// Offset is the number of entities to skip (for pagination)
	Offset int

	// SortBy specifies the field to sort results by.
	// Empty string means no specific sorting (plugin default).
	SortBy string

	// SortDesc indicates whether to sort in descending order.
	// False means ascending order.
	SortDesc bool
}
