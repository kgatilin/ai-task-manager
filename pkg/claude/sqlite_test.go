package claude_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/internal/storage"
	"github.com/kgatilin/darwinflow-pub/pkg/claude"
)

func TestNewSQLiteStore(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	defer store.Close()
}

func TestSQLiteStore_Init(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	err = store.Init(ctx)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}

	// Init should be idempotent
	err = store.Init(ctx)
	if err != nil {
		t.Errorf("Second Init failed: %v", err)
	}
}

func TestSQLiteStore_Store(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	record := storage.Record{
		ID:        "test-id-123",
		Timestamp: time.Now().UnixMilli(),
		EventType: "chat.message.user",
		Payload:   []byte(`{"message": "test"}`),
		Content:   "test message",
	}

	err = store.Store(ctx, record)
	if err != nil {
		t.Errorf("Store failed: %v", err)
	}
}

func TestSQLiteStore_Query(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := claude.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	if err := store.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Store some test records
	records := []storage.Record{
		{
			ID:        "test-1",
			Timestamp: time.Now().UnixMilli(),
			EventType: "chat.message.user",
			Payload:   []byte(`{"message": "hello"}`),
			Content:   "hello",
		},
		{
			ID:        "test-2",
			Timestamp: time.Now().UnixMilli() + 1000,
			EventType: "tool.invoked",
			Payload:   []byte(`{"tool": "Read"}`),
			Content:   "Read tool",
		},
	}

	for _, r := range records {
		if err := store.Store(ctx, r); err != nil {
			t.Fatalf("Store failed: %v", err)
		}
	}

	// Query all records
	filter := storage.Filter{
		Limit: 10,
	}

	results, err := store.Query(ctx, filter)
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 records, got %d", len(results))
	}

	// Query by event type
	filter = storage.Filter{
		EventTypes: []string{"chat.message.user"},
		Limit:      10,
	}

	results, err = store.Query(ctx, filter)
	if err != nil {
		t.Errorf("Query by event type failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 record, got %d", len(results))
	}

	if results[0].EventType != "chat.message.user" {
		t.Errorf("Expected event type 'chat.message.user', got '%s'", results[0].EventType)
	}
}
