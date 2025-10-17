package main

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kgatilin/darwinflow-pub/pkg/claude"
)

func TestRepeatString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		count  int
		want   string
	}{
		{
			name:  "repeat dash 5 times",
			s:     "-",
			count: 5,
			want:  "-----",
		},
		{
			name:  "repeat empty string",
			s:     "",
			count: 10,
			want:  "",
		},
		{
			name:  "repeat zero times",
			s:     "x",
			count: 0,
			want:  "",
		},
		{
			name:  "repeat multi-char string",
			s:     "ab",
			count: 3,
			want:  "ababab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repeatString(tt.s, tt.count)
			if got != tt.want {
				t.Errorf("repeatString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListLogs_EmptyDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create empty database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// listLogs should handle empty database gracefully
	// We can't easily test the output without refactoring, but we can
	// verify the database queries work correctly
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	query := "SELECT id, timestamp, event_type, payload, content FROM events ORDER BY timestamp DESC LIMIT ?"
	rows, err := db.Query(query, 10)
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 rows in empty database, got %d", count)
	}
}

func TestListLogs_WithData(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database with test data
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Insert test records
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	testRecords := []struct {
		id        string
		timestamp int64
		eventType string
		payload   string
		content   string
	}{
		{
			id:        "test-1",
			timestamp: time.Now().UnixMilli(),
			eventType: "chat.message.user",
			payload:   `{"message":"hello"}`,
			content:   "hello",
		},
		{
			id:        "test-2",
			timestamp: time.Now().UnixMilli() + 1000,
			eventType: "tool.invoked",
			payload:   `{"tool":"Read"}`,
			content:   "Read tool",
		},
		{
			id:        "test-3",
			timestamp: time.Now().UnixMilli() + 2000,
			eventType: "tool.result",
			payload:   `{"result":"success"}`,
			content:   "success",
		},
	}

	for _, r := range testRecords {
		_, err := db.Exec(
			"INSERT INTO events (id, timestamp, event_type, payload, content) VALUES (?, ?, ?, ?, ?)",
			r.id, r.timestamp, r.eventType, r.payload, r.content,
		)
		if err != nil {
			t.Fatalf("Insert failed: %v", err)
		}
	}

	// Test query with limit
	query := "SELECT id, timestamp, event_type, payload, content FROM events ORDER BY timestamp DESC LIMIT ?"
	rows, err := db.Query(query, 2)
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var id string
		var timestamp int64
		var eventType string
		var payload string
		var content string

		if err := rows.Scan(&id, &timestamp, &eventType, &payload, &content); err != nil {
			t.Errorf("Scan failed: %v", err)
		}
	}

	if count != 2 {
		t.Errorf("Expected 2 rows with limit=2, got %d", count)
	}

	// Test query without limit (should get all)
	rows2, err := db.Query("SELECT id FROM events")
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	defer rows2.Close()

	count2 := 0
	for rows2.Next() {
		count2++
	}

	if count2 != 3 {
		t.Errorf("Expected 3 total rows, got %d", count2)
	}
}

func TestExecuteRawQuery_SelectQuery(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database with test data
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Insert test data
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO events (id, timestamp, event_type, payload, content) VALUES (?, ?, ?, ?, ?)",
		"test-1", time.Now().UnixMilli(), "chat.started", `{"session":"123"}`, "session started",
	)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Test that query executes successfully
	rows, err := db.Query("SELECT event_type, COUNT(*) as count FROM events GROUP BY event_type")
	if err != nil {
		t.Errorf("Aggregate query failed: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var eventType string
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			t.Errorf("Scan failed: %v", err)
		}
		if eventType == "chat.started" && count == 1 {
			found = true
		}
	}

	if !found {
		t.Error("Expected to find chat.started event with count=1")
	}
}

func TestExecuteRawQuery_InvalidQuery(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Test invalid query
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Query("SELECT * FROM nonexistent_table")
	if err == nil {
		t.Error("Expected error for invalid query, got nil")
	}
}

func TestExecuteRawQuery_TimestampFormatting(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database with test data
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Insert test data with known timestamp
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	knownTime := time.Date(2025, 10, 17, 12, 30, 45, 0, time.UTC)
	timestamp := knownTime.UnixMilli()

	_, err = db.Exec(
		"INSERT INTO events (id, timestamp, event_type, payload, content) VALUES (?, ?, ?, ?, ?)",
		"test-ts", timestamp, "test.event", `{}`, "test",
	)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Query and verify timestamp can be read
	rows, err := db.Query("SELECT timestamp FROM events WHERE id = ?", "test-ts")
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var ts int64
		if err := rows.Scan(&ts); err != nil {
			t.Errorf("Scan failed: %v", err)
		}

		if ts != timestamp {
			t.Errorf("Expected timestamp %d, got %d", timestamp, ts)
		}

		// Verify we can convert it back to time
		retrievedTime := time.UnixMilli(ts)
		if !retrievedTime.Equal(knownTime) {
			t.Errorf("Expected time %v, got %v", knownTime, retrievedTime)
		}
	} else {
		t.Error("Expected to find row with timestamp")
	}
}

func TestDatabaseSchema(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Verify schema
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	// Check that events table exists with correct columns
	rows, err := db.Query("PRAGMA table_info(events)")
	if err != nil {
		t.Fatalf("PRAGMA table_info failed: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id":         false,
		"timestamp":  false,
		"event_type": false,
		"payload":    false,
		"content":    false,
	}

	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dfltValue interface{}
		var pk int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			t.Errorf("Scan failed: %v", err)
		}

		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
		}
	}

	// Verify all expected columns were found
	for col, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column %q not found in events table", col)
		}
	}

	// Check that indexes exist
	rows2, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='events'")
	if err != nil {
		t.Fatalf("Query for indexes failed: %v", err)
	}
	defer rows2.Close()

	indexCount := 0
	for rows2.Next() {
		var name string
		if err := rows2.Scan(&name); err != nil {
			t.Errorf("Scan failed: %v", err)
		}
		indexCount++
	}

	// Should have at least a few indexes (excluding auto-created ones)
	if indexCount < 3 {
		t.Errorf("Expected at least 3 indexes, got %d", indexCount)
	}
}

func TestParseLogsFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    logsOptions
		wantErr bool
	}{
		{
			name: "default flags",
			args: []string{},
			want: logsOptions{limit: 20, query: "", help: false},
		},
		{
			name: "custom limit",
			args: []string{"--limit", "50"},
			want: logsOptions{limit: 50, query: "", help: false},
		},
		{
			name: "with query",
			args: []string{"--query", "SELECT * FROM events"},
			want: logsOptions{limit: 20, query: "SELECT * FROM events", help: false},
		},
		{
			name: "help flag",
			args: []string{"--help"},
			want: logsOptions{limit: 20, query: "", help: true},
		},
		{
			name: "multiple flags",
			args: []string{"--limit", "100", "--query", "SELECT COUNT(*) FROM events"},
			want: logsOptions{limit: 100, query: "SELECT COUNT(*) FROM events", help: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogsFlags(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLogsFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.limit != tt.want.limit {
					t.Errorf("limit = %d, want %d", got.limit, tt.want.limit)
				}
				if got.query != tt.want.query {
					t.Errorf("query = %q, want %q", got.query, tt.want.query)
				}
				if got.help != tt.want.help {
					t.Errorf("help = %v, want %v", got.help, tt.want.help)
				}
			}
		})
	}
}

func TestQueryLogs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create and initialize database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Insert test data
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	for i := 0; i < 5; i++ {
		_, err := db.Exec(
			"INSERT INTO events (id, timestamp, event_type, payload, content) VALUES (?, ?, ?, ?, ?)",
			fmt.Sprintf("test-%d", i),
			time.Now().UnixMilli()+int64(i*1000),
			"test.event",
			`{"test":"data"}`,
			"test content",
		)
		if err != nil {
			t.Fatalf("Insert failed: %v", err)
		}
	}

	// Test querying logs
	records, err := queryLogs(dbPath, 3)
	if err != nil {
		t.Errorf("queryLogs failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// Verify ordering (should be DESC by timestamp, so newest first)
	if len(records) >= 2 {
		if records[0].Timestamp < records[1].Timestamp {
			t.Error("Records not ordered by timestamp DESC")
		}
	}
}

func TestQueryLogs_EmptyDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create empty database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	records, err := queryLogs(dbPath, 10)
	if err != nil {
		t.Errorf("queryLogs failed: %v", err)
	}

	if len(records) != 0 {
		t.Errorf("Expected 0 records from empty DB, got %d", len(records))
	}
}

func TestFormatLogRecord(t *testing.T) {
	record := logRecord{
		ID:        "test-123",
		Timestamp: time.Date(2025, 10, 17, 12, 30, 45, 0, time.UTC).UnixMilli(),
		EventType: "chat.message.user",
		Payload:   []byte(`{"message":"hello"}`),
		Content:   "hello",
	}

	output := formatLogRecord(0, record)

	// Verify output contains expected elements
	if !contains(output, "test-123") {
		t.Error("Output should contain ID")
	}
	if !contains(output, "chat.message.user") {
		t.Error("Output should contain event type")
	}
	if !contains(output, "hello") {
		t.Error("Output should contain content")
	}
}

func TestFormatLogRecord_LongContent(t *testing.T) {
	longContent := string(make([]byte, 300))
	for i := range longContent {
		longContent = longContent[:i] + "x"
	}

	record := logRecord{
		ID:        "test-456",
		Timestamp: time.Now().UnixMilli(),
		EventType: "test.event",
		Payload:   []byte(`{}`),
		Content:   longContent,
	}

	output := formatLogRecord(0, record)

	// Content should be truncated
	if !contains(output, "...") {
		t.Error("Long content should be truncated with ...")
	}
}

func TestFormatQueryValue(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "nil value",
			input: nil,
			want:  "NULL",
		},
		{
			name:  "string value",
			input: "test",
			want:  "test",
		},
		{
			name:  "long string",
			input: string(make([]byte, 150)),
			want:  "...",
		},
		{
			name:  "int64 regular",
			input: int64(42),
			want:  "42",
		},
		{
			name:  "int64 timestamp",
			input: int64(1697548245000),
			want:  "1697548245000 (2023-10-17",
		},
		{
			name:  "json bytes",
			input: []byte(`{"key":"value"}`),
			want:  `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatQueryValue(tt.input)
			if !contains(got, tt.want) {
				t.Errorf("formatQueryValue() = %q, should contain %q", got, tt.want)
			}
		})
	}
}

func TestListLogs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create and initialize database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Insert test data
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO events (id, timestamp, event_type, payload, content) VALUES (?, ?, ?, ?, ?)",
		"test-1", time.Now().UnixMilli(), "test.event", `{"test":"data"}`, "test content",
	)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Test listLogs - should not error
	err = listLogs(dbPath, 10)
	if err != nil {
		t.Errorf("listLogs failed: %v", err)
	}
}

func TestListLogs_EmptyDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create empty database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Test listLogs with empty database - should not error
	err = listLogs(dbPath, 10)
	if err != nil {
		t.Errorf("listLogs with empty DB failed: %v", err)
	}
}

func TestListLogs_InvalidPath(t *testing.T) {
	// Test with non-existent database
	err := listLogs("/nonexistent/path/db.sqlite", 10)
	if err == nil {
		t.Error("Expected error for non-existent database, got nil")
	}
}

func TestExecuteRawQuery_Success(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create and initialize database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Test executeRawQuery with valid query
	err = executeRawQuery(dbPath, "SELECT COUNT(*) FROM events")
	if err != nil {
		t.Errorf("executeRawQuery failed: %v", err)
	}
}

func TestExecuteRawQuery_InvalidSQL(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Test with invalid SQL
	err = executeRawQuery(dbPath, "INVALID SQL QUERY")
	if err == nil {
		t.Error("Expected error for invalid SQL, got nil")
	}
}

func TestExecuteRawQuery_WithData(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create and initialize database
	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	store.Close()

	// Insert test data
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO events (id, timestamp, event_type, payload, content) VALUES (?, ?, ?, ?, ?)",
		"test-1", time.Now().UnixMilli(), "test.event", `{"test":"data"}`, "test content",
	)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Test executeRawQuery with data
	err = executeRawQuery(dbPath, "SELECT * FROM events WHERE event_type = 'test.event'")
	if err != nil {
		t.Errorf("executeRawQuery with data failed: %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
