package storage

import (
	"context"
	"time"
)

// Store defines the interface for persisting events
type Store interface {
	// Init initializes the storage (creates tables, indexes, etc.)
	Init(ctx context.Context) error

	// Store saves an event record
	Store(ctx context.Context, record Record) error

	// Query retrieves events based on filter criteria
	Query(ctx context.Context, filter Filter) ([]Record, error)

	// Close closes the storage connection
	Close() error
}

// Record represents a storable event record
type Record struct {
	ID        string
	Timestamp int64  // Unix timestamp in milliseconds
	EventType string
	Payload   []byte // JSON payload
	Content   string // Normalized content for search
}

// Filter defines query parameters for retrieving events
type Filter struct {
	// Time range
	StartTime *time.Time
	EndTime   *time.Time

	// Event type filtering
	EventTypes []string

	// Context filtering
	Context string

	// Full-text search
	SearchText string

	// Pagination
	Limit  int
	Offset int
}
