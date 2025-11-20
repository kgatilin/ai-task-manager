package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/persistence"
)

// Helper to create a test database
func createTestDB(t *testing.T) *sql.DB {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Initialize schema
	if err := persistence.InitSchema(db); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return db
}

// Helper to create a test logger
func createTestLogger() logger.Logger {
	return &testLogger{}
}

type testLogger struct{}

func (l *testLogger) Debug(msg string, fields ...interface{})        {}
func (l *testLogger) Info(msg string, fields ...interface{})         {}
func (l *testLogger) Warn(msg string, fields ...interface{})         {}
func (l *testLogger) Error(msg string, fields ...interface{})        {}
func (l *testLogger) WithFields(fields ...interface{}) logger.Logger { return l }
func (l *testLogger) GetLevel() logger.Level                         { return 1 }
func (l *testLogger) SetLevel(level logger.Level)                    {}

// ============================================================================
// Roadmap Tests
// ============================================================================

func TestSaveAndGetRoadmap(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	ctx := context.Background()

	// Create a roadmap
	roadmap, err := entities.NewRoadmapEntity(
		"roadmap-1",
		"Build the best system",
		"Deliver on time and quality",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("failed to create roadmap entity: %v", err)
	}

	// Save roadmap
	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	// Get roadmap
	retrieved, err := repo.GetRoadmap(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("failed to get roadmap: %v", err)
	}

	if retrieved.ID != roadmap.ID {
		t.Errorf("expected roadmap ID %s, got %s", roadmap.ID, retrieved.ID)
	}
	if retrieved.Vision != roadmap.Vision {
		t.Errorf("expected vision %s, got %s", roadmap.Vision, retrieved.Vision)
	}
}

func TestSaveRoadmapDuplicate(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	ctx := context.Background()

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())

	if err := repo.SaveRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to save roadmap: %v", err)
	}

	// Try to save duplicate
	if err := repo.SaveRoadmap(ctx, roadmap); err == nil {
		t.Error("expected error when saving duplicate roadmap")
	} else if !errors.Is(err, tmerrors.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got: %v", err)
	}
}

func TestGetRoadmapNotFound(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	ctx := context.Background()

	_, err := repo.GetRoadmap(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent roadmap")
	} else if !errors.Is(err, tmerrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestUpdateRoadmap(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	ctx := context.Background()

	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	// Update roadmap
	roadmap.Vision = "new vision"
	roadmap.UpdatedAt = time.Now().UTC()

	if err := repo.UpdateRoadmap(ctx, roadmap); err != nil {
		t.Fatalf("failed to update roadmap: %v", err)
	}

	// Verify update
	retrieved, _ := repo.GetRoadmap(ctx, "roadmap-1")
	if retrieved.Vision != "new vision" {
		t.Errorf("expected vision to be updated, got %s", retrieved.Vision)
	}
}

func TestGetActiveRoadmap(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	ctx := context.Background()

	now := time.Now().UTC()

	// Create first roadmap
	roadmap1, _ := entities.NewRoadmapEntity("roadmap-1", "vision1", "criteria1", now, now)
	repo.SaveRoadmap(ctx, roadmap1)

	time.Sleep(10 * time.Millisecond)

	// Create second roadmap (more recent)
	roadmap2, _ := entities.NewRoadmapEntity("roadmap-2", "vision2", "criteria2", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap2)

	// Get active roadmap should return the most recent one
	active, err := repo.GetActiveRoadmap(ctx)
	if err != nil {
		t.Fatalf("failed to get active roadmap: %v", err)
	}

	if active.ID != "roadmap-2" {
		t.Errorf("expected roadmap-2, got %s", active.ID)
	}
}
