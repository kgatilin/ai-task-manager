package pluginsdk

import "time"

// IExtensible is the REQUIRED base capability that ALL entities must implement.
// It provides core identity and introspection for any entity in the system.
//
// Required methods:
//   - GetID: Unique identifier for the entity
//   - GetType: Entity type name (e.g., "session", "task", "roadmap")
//   - GetCapabilities: List of capability interfaces this entity implements
//   - GetField: Access individual fields by name
//   - GetAllFields: Get all entity fields as a map
//
// Currently used by: TUI entity list/detail views, PluginRegistry routing
type IExtensible interface {
	// GetID returns the unique identifier for this entity
	GetID() string

	// GetType returns the entity type (e.g., "session", "task", "roadmap")
	GetType() string

	// GetCapabilities returns list of capability names this entity supports
	// (e.g., ["IExtensible", "IHasContext", "ITrackable"])
	GetCapabilities() []string

	// GetField retrieves a named field value
	// Returns nil if field doesn't exist
	GetField(name string) interface{}

	// GetAllFields returns all fields as a map
	GetAllFields() map[string]interface{}
}

// IHasContext is an OPTIONAL capability for entities with related contextual data.
// Entities implement this when they have associated files, activity records, or metadata.
//
// Use cases:
//   - Sessions with tool invocations and file references
//   - Tasks with associated code changes
//   - Projects with related artifacts
//
// Currently used by: [none - available for future use]
type IHasContext interface {
	IExtensible

	// GetContext returns contextual information about this entity
	GetContext() *EntityContext
}

// EntityContext contains contextual information about an entity
type EntityContext struct {
	// RelatedEntities are other entities connected to this one
	// Key is entity type, value is list of entity IDs
	RelatedEntities map[string][]string

	// LinkedFiles are file paths referenced by this entity
	LinkedFiles []string

	// RecentActivity is a log of recent actions on this entity
	RecentActivity []ActivityRecord

	// Metadata for any additional context
	Metadata map[string]interface{}
}

// ActivityRecord represents a single activity event related to an entity
type ActivityRecord struct {
	// Timestamp is when the activity occurred
	Timestamp time.Time `json:"timestamp"`

	// Type is the kind of activity (e.g., "created", "updated", "analyzed")
	Type string `json:"type"`

	// Description is a human-readable description of the activity
	Description string `json:"description"`

	// Actor is who/what performed the activity (user, system, etc.)
	Actor string `json:"actor"`
}

// ITrackable is an OPTIONAL capability for entities with status and progress tracking.
// Entities implement this when they have lifecycle states and completion metrics.
//
// Use cases:
//   - Sessions that can be active/completed/analyzed
//   - Tasks with todo/in-progress/done states
//   - Workflows with multi-step progress
//
// Currently used by: [none - available for future use]
type ITrackable interface {
	IExtensible

	// GetStatus returns the current status (e.g., "active", "completed", "blocked")
	GetStatus() string

	// GetProgress returns completion progress as a value between 0.0 and 1.0
	GetProgress() float64

	// IsBlocked returns true if the entity is blocked from progressing
	IsBlocked() bool

	// GetBlockReason returns the reason for blocking, or empty string if not blocked
	GetBlockReason() string
}

// ISchedulable is an OPTIONAL capability for entities with time-based scheduling.
// Entities implement this when they have start dates, due dates, or deadlines.
//
// Use cases:
//   - Tasks with deadlines
//   - Scheduled workflows
//   - Time-boxed sessions
//
// Currently used by: [none - available for future use]
type ISchedulable interface {
	IExtensible

	// GetStartDate returns when the entity should/did start
	GetStartDate() *time.Time

	// GetDueDate returns when the entity should be completed
	GetDueDate() *time.Time

	// IsOverdue returns true if past due date and not complete
	IsOverdue() bool
}

// IRelatable is an OPTIONAL capability for entities with explicit relationships.
// Entities implement this when they reference other entities.
//
// Use cases:
//   - Tasks related to sessions
//   - Projects containing tasks
//   - Hierarchical workflows
//
// Currently used by: [none - available for future use]
type IRelatable interface {
	IExtensible

	// GetRelated returns IDs of related entities of the specified type
	GetRelated(entityType string) []string

	// GetAllRelations returns all relationships grouped by entity type
	GetAllRelations() map[string][]string
}
