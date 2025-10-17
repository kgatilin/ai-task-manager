package claude

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kgatilin/darwinflow-pub/internal/storage"
)

// SQLiteStore implements storage.Store using SQLite
type SQLiteStore struct {
	db   *sql.DB
	path string
}

// NewSQLiteStore creates a new SQLite-backed store
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	return &SQLiteStore{
		db:   db,
		path: dbPath,
	}, nil
}

// Init initializes the database schema
func (s *SQLiteStore) Init(ctx context.Context) error {
	// Create base schema
	baseSchema := `
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			timestamp INTEGER NOT NULL,
			event_type TEXT NOT NULL,
			payload TEXT NOT NULL,
			content TEXT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
		CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
		CREATE INDEX IF NOT EXISTS idx_events_timestamp_type ON events(timestamp, event_type);
	`

	_, err := s.db.ExecContext(ctx, baseSchema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Try to create FTS5 virtual table (optional, may not be available)
	ftsSchema := `
		CREATE VIRTUAL TABLE IF NOT EXISTS events_fts USING fts5(
			content,
			content=events,
			content_rowid=rowid
		);

		CREATE TRIGGER IF NOT EXISTS events_fts_insert AFTER INSERT ON events BEGIN
			INSERT INTO events_fts(rowid, content) VALUES (new.rowid, new.content);
		END;

		CREATE TRIGGER IF NOT EXISTS events_fts_delete AFTER DELETE ON events BEGIN
			DELETE FROM events_fts WHERE rowid = old.rowid;
		END;

		CREATE TRIGGER IF NOT EXISTS events_fts_update AFTER UPDATE ON events BEGIN
			DELETE FROM events_fts WHERE rowid = old.rowid;
			INSERT INTO events_fts(rowid, content) VALUES (new.rowid, new.content);
		END;
	`

	// Attempt FTS5, but don't fail if unavailable
	_, _ = s.db.ExecContext(ctx, ftsSchema)

	return nil
}

// Store saves an event record
func (s *SQLiteStore) Store(ctx context.Context, record storage.Record) error {
	query := `
		INSERT INTO events (id, timestamp, event_type, payload, content)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		record.ID,
		record.Timestamp,
		record.EventType,
		string(record.Payload),
		record.Content,
	)

	if err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	return nil
}

// Query retrieves events based on filter criteria
func (s *SQLiteStore) Query(ctx context.Context, filter storage.Filter) ([]storage.Record, error) {
	var conditions []string
	var args []interface{}

	// Build WHERE clause
	if filter.StartTime != nil {
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, filter.StartTime.UnixMilli())
	}

	if filter.EndTime != nil {
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, filter.EndTime.UnixMilli())
	}

	if len(filter.EventTypes) > 0 {
		placeholders := make([]string, len(filter.EventTypes))
		for i, et := range filter.EventTypes {
			placeholders[i] = "?"
			args = append(args, et)
		}
		conditions = append(conditions, fmt.Sprintf("event_type IN (%s)", strings.Join(placeholders, ",")))
	}

	// Build query
	query := "SELECT id, timestamp, event_type, payload, content FROM events"

	if filter.SearchText != "" {
		// Try FTS search first, fall back to LIKE if FTS not available
		ftsQuery := `
			SELECT e.id, e.timestamp, e.event_type, e.payload, e.content
			FROM events e
			JOIN events_fts fts ON fts.rowid = e.rowid
			WHERE fts.content MATCH ?
		`
		ftsArgs := append([]interface{}{filter.SearchText}, args...)

		if len(conditions) > 0 {
			ftsQuery += " AND " + strings.Join(conditions, " AND ")
		}

		// Try FTS query
		_, err := s.db.QueryContext(ctx, ftsQuery+" LIMIT 1", ftsArgs...)
		if err == nil {
			// FTS is available
			query = ftsQuery
			args = ftsArgs
		} else {
			// Fall back to LIKE search
			conditions = append([]string{"content LIKE ?"}, conditions...)
			args = append([]interface{}{"%" + filter.SearchText + "%"}, args...)
			if len(conditions) > 0 {
				query += " WHERE " + strings.Join(conditions, " AND ")
			}
		}
	} else if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY timestamp DESC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			query += " OFFSET ?"
			args = append(args, filter.Offset)
		}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var records []storage.Record
	for rows.Next() {
		var r storage.Record
		var payloadStr string

		if err := rows.Scan(&r.ID, &r.Timestamp, &r.EventType, &payloadStr, &r.Content); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		r.Payload = []byte(payloadStr)
		records = append(records, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
