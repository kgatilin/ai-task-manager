package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/infrastructure/persistence"
)

// ============================================================================
// ADR Tests
// ============================================================================

func TestSaveAndGetADR(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	adrRepo := persistence.NewSQLiteADRRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create ADR
	adr, err := entities.NewADREntity(
		"adr-1",
		"track-1",
		"Use SQLite for storage",
		"proposed",
		"We need a database",
		"Use SQLite",
		"Good performance",
		"PostgreSQL, MySQL",
		time.Now().UTC(),
		time.Now().UTC(),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create ADR: %v", err)
	}

	// Save ADR
	if err := adrRepo.SaveADR(ctx, adr); err != nil {
		t.Fatalf("failed to save ADR: %v", err)
	}

	// Get ADR
	retrieved, err := adrRepo.GetADR(ctx, "adr-1")
	if err != nil {
		t.Fatalf("failed to get ADR: %v", err)
	}

	if retrieved.ID != adr.ID || retrieved.Title != adr.Title {
		t.Errorf("ADR mismatch")
	}
	if retrieved.Status != "proposed" {
		t.Errorf("expected status proposed, got %s", retrieved.Status)
	}
}

func TestListADRs(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	adrRepo := persistence.NewSQLiteADRRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create multiple ADRs
	for i := 1; i <= 3; i++ {
		id := "adr-" + string(rune(48+i))
		adr, _ := entities.NewADREntity(id, "track-1", "ADR "+string(rune(48+i)), "proposed", "context", "decision", "consequences", "", time.Now().UTC(), time.Now().UTC(), nil)
		adrRepo.SaveADR(ctx, adr)
	}

	// List ADRs
	trackID := "track-1"
	adrs, err := adrRepo.ListADRs(ctx, &trackID)
	if err != nil {
		t.Fatalf("failed to list ADRs: %v", err)
	}

	if len(adrs) != 3 {
		t.Errorf("expected 3 ADRs, got %d", len(adrs))
	}
}

func TestListAllADRs(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	adrRepo := persistence.NewSQLiteADRRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	// Create ADRs with different statuses
	adr1, _ := entities.NewADREntity("adr-1", "track-1", "ADR 1", "proposed", "context", "decision", "consequences", "", time.Now().UTC(), time.Now().UTC(), nil)
	adr2, _ := entities.NewADREntity("adr-2", "track-1", "ADR 2", "accepted", "context", "decision", "consequences", "", time.Now().UTC(), time.Now().UTC(), nil)

	adrRepo.SaveADR(ctx, adr1)
	adrRepo.SaveADR(ctx, adr2)

	// List all ADRs (nil trackID)
	adrs, err := adrRepo.ListADRs(ctx, nil)
	if err != nil {
		t.Fatalf("failed to list ADRs: %v", err)
	}

	if len(adrs) < 2 {
		t.Errorf("expected at least 2 ADRs, got %d", len(adrs))
	}
}

func TestUpdateADR(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	adrRepo := persistence.NewSQLiteADRRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	trackRepo.SaveTrack(ctx, track)

	adr, _ := entities.NewADREntity("adr-1", "track-1", "Old Title", "proposed", "context", "decision", "consequences", "", time.Now().UTC(), time.Now().UTC(), nil)
	adrRepo.SaveADR(ctx, adr)

	// Update ADR
	adr.Title = "New Title"
	adr.Status = "accepted"
	adr.UpdatedAt = time.Now().UTC()

	if err := adrRepo.UpdateADR(ctx, adr); err != nil {
		t.Fatalf("failed to update ADR: %v", err)
	}

	// Verify update
	retrieved, _ := adrRepo.GetADR(ctx, "adr-1")
	if retrieved.Title != "New Title" {
		t.Errorf("expected title to be updated, got %s", retrieved.Title)
	}
	if retrieved.Status != "accepted" {
		t.Errorf("expected status accepted, got %s", retrieved.Status)
	}
}

func TestGetADRsByTrack(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	roadmapRepo := persistence.NewSQLiteRoadmapRepository(db, createTestLogger())
	trackRepo := persistence.NewSQLiteTrackRepository(db, createTestLogger())
	adrRepo := persistence.NewSQLiteADRRepository(db, createTestLogger())
	ctx := context.Background()

	// Setup
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	roadmapRepo.SaveRoadmap(ctx, roadmap)

	track1, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track 1", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())
	track2, _ := entities.NewTrackEntity("track-2", "roadmap-1", "Track 2", "", "not-started", 200, []string{}, time.Now().UTC(), time.Now().UTC())

	trackRepo.SaveTrack(ctx, track1)
	trackRepo.SaveTrack(ctx, track2)

	// Create ADRs for different tracks
	adr1, _ := entities.NewADREntity("adr-1", "track-1", "ADR 1", "proposed", "context", "decision", "consequences", "", time.Now().UTC(), time.Now().UTC(), nil)
	adr2, _ := entities.NewADREntity("adr-2", "track-1", "ADR 2", "proposed", "context", "decision", "consequences", "", time.Now().UTC(), time.Now().UTC(), nil)
	adr3, _ := entities.NewADREntity("adr-3", "track-2", "ADR 3", "proposed", "context", "decision", "consequences", "", time.Now().UTC(), time.Now().UTC(), nil)

	adrRepo.SaveADR(ctx, adr1)
	adrRepo.SaveADR(ctx, adr2)
	adrRepo.SaveADR(ctx, adr3)

	// Get ADRs for track-1
	adrs, err := adrRepo.GetADRsByTrack(ctx, "track-1")
	if err != nil {
		t.Fatalf("failed to get ADRs by track: %v", err)
	}

	if len(adrs) != 2 {
		t.Errorf("expected 2 ADRs for track-1, got %d", len(adrs))
	}
}

