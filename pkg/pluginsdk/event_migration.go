package pluginsdk

// eventMigrationHandler defines the interface for event schema migration.
// This allows future schema evolution while maintaining backward compatibility.
type eventMigrationHandler struct {
	// supportedVersions maps schema version strings to their migration functions
	supportedVersions map[string]func(*Event) *Event
}

// migrateEvent migrates an event from its current version to the latest schema version.
// This is a no-op for schema version "1.0" and serves as the foundation for future migrations.
//
// Usage:
// - Version "1.0": No migration needed (current version)
// - Future versions: Migration functions would be added here to handle schema changes
//
// Example migration function (for future use):
//   func migrateFromV1ToV2(event *Event) *Event {
//       // Transform event structure from v1 to v2
//       return event
//   }
func migrateEvent(event *Event) *Event {
	if event == nil {
		return nil
	}

	// Ensure version is set
	if event.Version == "" {
		event.Version = "1.0"
	}

	// For version "1.0", no migration is needed
	// Future versions would apply their respective migrations here
	switch event.Version {
	case "1.0":
		// Current version - no migration needed
		return event
	default:
		// Unknown version - return as-is with default version set
		event.Version = "1.0"
		return event
	}
}
