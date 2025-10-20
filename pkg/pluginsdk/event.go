package pluginsdk

import "time"

// Event is the standard event structure for plugin-emitted events.
// Events are emitted by plugins and stored in the framework's event database.
type Event struct {
	// Type is the event type identifier (e.g., "tool.invoked", "task.created").
	// Use dot notation to namespace events: "<domain>.<action>".
	Type string

	// Source is the name of the plugin that emitted this event
	Source string

	// Timestamp is when the event occurred
	Timestamp time.Time

	// Payload contains the event-specific data.
	// Structure depends on the event type.
	Payload map[string]interface{}

	// Metadata contains additional context about the event.
	// Common fields: session_id, user_id, environment, etc.
	Metadata map[string]string

	// Version is the schema version for this event (for future schema evolution).
	// Default value: "1.0"
	// Used to handle backward compatibility when event schemas change.
	Version string
}
