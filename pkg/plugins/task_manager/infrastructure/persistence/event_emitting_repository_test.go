package persistence_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/events"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/infrastructure/persistence"
)

// ============================================================================
// Mock Event Bus
// ============================================================================

// mockEventBus captures emitted events for verification
type mockEventBus struct {
	events []pluginsdk.BusEvent
}

func newMockEventBus() *mockEventBus {
	return &mockEventBus{
		events: []pluginsdk.BusEvent{},
	}
}

func (m *mockEventBus) Publish(ctx context.Context, event pluginsdk.BusEvent) error {
	m.events = append(m.events, event)
	return nil
}

func (m *mockEventBus) Subscribe(filter pluginsdk.EventFilter, handler pluginsdk.EventHandler) (string, error) {
	return "sub-id", nil
}

func (m *mockEventBus) Unsubscribe(subscriptionID string) error {
	return nil
}

// getEventByType returns the first event with the given type
func (m *mockEventBus) getEventByType(eventType string) *pluginsdk.BusEvent {
	for _, e := range m.events {
		if e.Type == eventType {
			return &e
		}
	}
	return nil
}

// getEventCount returns the number of events emitted
func (m *mockEventBus) getEventCount() int {
	return len(m.events)
}

// reset clears all captured events
func (m *mockEventBus) reset() {
	m.events = []pluginsdk.BusEvent{}
}

// setupEventEmittingRepo creates a repository wrapped with event emission
func setupEventEmittingRepo(t *testing.T, db *sql.DB) (*persistence.EventEmittingRepository, *mockEventBus) {
	logger := createTestLogger()

	// Create the underlying repository (contains all 6 repo implementations)
	baseRepo := persistence.NewSQLiteRoadmapRepository(db, logger)

	// Create mock event bus
	mockBus := newMockEventBus()

	// Wrap with event-emitting decorator
	eventRepo := persistence.NewEventEmittingRepository(baseRepo, mockBus, logger)

	return eventRepo, mockBus
}

// ============================================================================
// Roadmap Event Tests
// ============================================================================

func TestEventEmittingRepository_SaveRoadmap_EmitsCreatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Create and save roadmap
	roadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Test Vision",
		"Test Criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)

	err := repo.SaveRoadmap(ctx, roadmap)
	if err != nil {
		t.Fatalf("SaveRoadmap failed: %v", err)
	}

	// Verify event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventRoadmapCreated)
	if event == nil {
		t.Fatal("expected RoadmapCreated event, but none found")
	}

	// Verify event type and source
	if event.Type != events.EventRoadmapCreated {
		t.Errorf("expected event type %s, got %s", events.EventRoadmapCreated, event.Type)
	}
	if event.Source != events.PluginSourceName {
		t.Errorf("expected source %s, got %s", events.PluginSourceName, event.Source)
	}
}

func TestEventEmittingRepository_UpdateRoadmap_EmitsUpdatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Create roadmap first
	roadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Test Vision",
		"Test Criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	repo.SaveRoadmap(ctx, roadmap)

	// Reset mock to clear creation event
	mockBus.reset()

	// Update roadmap
	roadmap.Vision = "Updated Vision"
	roadmap.UpdatedAt = time.Now().UTC()

	err := repo.UpdateRoadmap(ctx, roadmap)
	if err != nil {
		t.Fatalf("UpdateRoadmap failed: %v", err)
	}

	// Verify update event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventRoadmapUpdated)
	if event == nil {
		t.Fatal("expected RoadmapUpdated event, but none found")
	}
}

// ============================================================================
// Track Event Tests
// ============================================================================

func TestEventEmittingRepository_SaveTrack_EmitsCreatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Create roadmap first
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)
	mockBus.reset()

	// Create and save track
	track, _ := entities.NewTrackEntity(
		"track-1",
		"roadmap-1",
		"Test Track",
		"Description",
		"not-started",
		100,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)

	err := repo.SaveTrack(ctx, track)
	if err != nil {
		t.Fatalf("SaveTrack failed: %v", err)
	}

	// Verify event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventTrackCreated)
	if event == nil {
		t.Fatal("expected TrackCreated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateTrack_EmitsUpdatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity(
		"track-1",
		"roadmap-1",
		"Test Track",
		"Description",
		"not-started",
		100,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	repo.SaveTrack(ctx, track)
	mockBus.reset()

	// Update track (no status change)
	track.Title = "Updated Track"
	track.UpdatedAt = time.Now().UTC()

	err := repo.UpdateTrack(ctx, track)
	if err != nil {
		t.Fatalf("UpdateTrack failed: %v", err)
	}

	// Verify only update event was emitted (no status change)
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventTrackUpdated)
	if event == nil {
		t.Fatal("expected TrackUpdated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateTrack_EmitsStatusChangedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity(
		"track-1",
		"roadmap-1",
		"Test Track",
		"Description",
		"not-started",
		100,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	repo.SaveTrack(ctx, track)
	mockBus.reset()

	// Update track status
	track.Status = "in-progress"
	track.UpdatedAt = time.Now().UTC()

	err := repo.UpdateTrack(ctx, track)
	if err != nil {
		t.Fatalf("UpdateTrack failed: %v", err)
	}

	// Verify both update and status_changed events were emitted
	if mockBus.getEventCount() != 2 {
		t.Errorf("expected 2 events, got %d", mockBus.getEventCount())
	}

	updateEvent := mockBus.getEventByType(events.EventTrackUpdated)
	if updateEvent == nil {
		t.Error("expected TrackUpdated event, but none found")
	}

	statusEvent := mockBus.getEventByType(events.EventTrackStatusChanged)
	if statusEvent == nil {
		t.Error("expected TrackStatusChanged event, but none found")
	}
}

func TestEventEmittingRepository_UpdateTrack_EmitsCompletedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity(
		"track-1",
		"roadmap-1",
		"Test Track",
		"Description",
		"not-started",
		100,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	repo.SaveTrack(ctx, track)
	mockBus.reset()

	// Update track to complete status
	track.Status = "complete"
	track.UpdatedAt = time.Now().UTC()

	err := repo.UpdateTrack(ctx, track)
	if err != nil {
		t.Fatalf("UpdateTrack failed: %v", err)
	}

	// Verify updated, status_changed, and completed events were emitted
	if mockBus.getEventCount() != 3 {
		t.Errorf("expected 3 events, got %d", mockBus.getEventCount())
	}

	completedEvent := mockBus.getEventByType(events.EventTrackCompleted)
	if completedEvent == nil {
		t.Error("expected TrackCompleted event, but none found")
	}
}

func TestEventEmittingRepository_UpdateTrack_EmitsBlockedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity(
		"track-1",
		"roadmap-1",
		"Test Track",
		"Description",
		"not-started",
		100,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	repo.SaveTrack(ctx, track)
	mockBus.reset()

	// Update track to blocked status
	track.Status = "blocked"
	track.UpdatedAt = time.Now().UTC()

	err := repo.UpdateTrack(ctx, track)
	if err != nil {
		t.Fatalf("UpdateTrack failed: %v", err)
	}

	// Verify updated, status_changed, and blocked events were emitted
	if mockBus.getEventCount() != 3 {
		t.Errorf("expected 3 events, got %d", mockBus.getEventCount())
	}

	blockedEvent := mockBus.getEventByType(events.EventTrackBlocked)
	if blockedEvent == nil {
		t.Error("expected TrackBlocked event, but none found")
	}
}

// ============================================================================
// Task Event Tests
// ============================================================================

func TestEventEmittingRepository_SaveTask_EmitsCreatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)
	mockBus.reset()

	// Create and save task
	task, _ := entities.NewTaskEntity(
		"task-1",
		"track-1",
		"Test Task",
		"Description",
		"todo",
		100,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)

	err := repo.SaveTask(ctx, task)
	if err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	// Verify event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventTaskCreated)
	if event == nil {
		t.Fatal("expected TaskCreated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateTask_EmitsUpdatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, and task
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)
	mockBus.reset()

	// Update task (no status change)
	task.Title = "Updated Task"
	task.UpdatedAt = time.Now().UTC()

	err := repo.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	// Verify only update event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventTaskUpdated)
	if event == nil {
		t.Fatal("expected TaskUpdated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateTask_EmitsStatusChangedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, and task
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)
	mockBus.reset()

	// Update task status
	task.Status = "in-progress"
	task.UpdatedAt = time.Now().UTC()

	err := repo.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	// Verify both update and status_changed events were emitted
	if mockBus.getEventCount() != 2 {
		t.Errorf("expected 2 events, got %d", mockBus.getEventCount())
	}

	updateEvent := mockBus.getEventByType(events.EventTaskUpdated)
	if updateEvent == nil {
		t.Error("expected TaskUpdated event, but none found")
	}

	statusEvent := mockBus.getEventByType(events.EventTaskStatusChanged)
	if statusEvent == nil {
		t.Error("expected TaskStatusChanged event, but none found")
	}
}

func TestEventEmittingRepository_UpdateTask_EmitsCompletedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, and task
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)
	mockBus.reset()

	// Update task to done status
	task.Status = "done"
	task.UpdatedAt = time.Now().UTC()

	err := repo.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	// Verify updated, status_changed, and completed events were emitted
	if mockBus.getEventCount() != 3 {
		t.Errorf("expected 3 events, got %d", mockBus.getEventCount())
	}

	completedEvent := mockBus.getEventByType(events.EventTaskCompleted)
	if completedEvent == nil {
		t.Error("expected TaskCompleted event, but none found")
	}
}

// ============================================================================
// Iteration Event Tests
// ============================================================================

func TestEventEmittingRepository_SaveIteration_EmitsCreatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)
	mockBus.reset()

	// Create and save iteration
	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(
		1,
		"Iteration 1",
		"Goal",
		"Deliverable",
		[]string{},
		"planned",
		100,
		time.Time{},
		time.Time{},
		now,
		now,
	)

	err := repo.SaveIteration(ctx, iteration)
	if err != nil {
		t.Fatalf("SaveIteration failed: %v", err)
	}

	// Verify event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventIterationCreated)
	if event == nil {
		t.Fatal("expected IterationCreated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateIteration_EmitsUpdatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and iteration
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iter 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	repo.SaveIteration(ctx, iteration)
	mockBus.reset()

	// Update iteration
	iteration.Name = "Updated Iteration"
	iteration.UpdatedAt = time.Now().UTC()

	err := repo.UpdateIteration(ctx, iteration)
	if err != nil {
		t.Fatalf("UpdateIteration failed: %v", err)
	}

	// Verify update event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventIterationUpdated)
	if event == nil {
		t.Fatal("expected IterationUpdated event, but none found")
	}
}

func TestEventEmittingRepository_StartIteration_EmitsStartedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and iteration
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iter 1", "Goal", "Deliverable", []string{}, "planned", 100, time.Time{}, time.Time{}, now, now)
	repo.SaveIteration(ctx, iteration)
	mockBus.reset()

	// Start iteration
	err := repo.StartIteration(ctx, 1)
	if err != nil {
		t.Fatalf("StartIteration failed: %v", err)
	}

	// Verify started event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventIterationStarted)
	if event == nil {
		t.Fatal("expected IterationStarted event, but none found")
	}
}

func TestEventEmittingRepository_CompleteIteration_EmitsCompletedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and iteration (must be current to complete)
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	now := time.Now().UTC()
	iteration, _ := entities.NewIterationEntity(1, "Iter 1", "Goal", "Deliverable", []string{}, "current", 100, time.Time{}, time.Time{}, now, now)
	repo.SaveIteration(ctx, iteration)
	mockBus.reset()

	// Complete iteration
	err := repo.CompleteIteration(ctx, 1)
	if err != nil {
		t.Fatalf("CompleteIteration failed: %v", err)
	}

	// Verify completed event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventIterationCompleted)
	if event == nil {
		t.Fatal("expected IterationCompleted event, but none found")
	}
}

// ============================================================================
// ADR Event Tests
// ============================================================================

func TestEventEmittingRepository_SaveADR_EmitsCreatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap and track
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)
	mockBus.reset()

	// Create and save ADR
	now := time.Now().UTC()
	adr, err := entities.NewADREntity(
		"adr-1",
		"track-1",
		"Test ADR",
		"proposed",
		"Context",
		"Decision",
		"Consequences",
		"Alternatives",
		now,
		now,
		nil,
	)
	if err != nil {
		t.Fatalf("NewADREntity failed: %v", err)
	}

	err = repo.SaveADR(ctx, adr)
	if err != nil {
		t.Fatalf("SaveADR failed: %v", err)
	}

	// Verify event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventADRCreated)
	if event == nil {
		t.Fatal("expected ADRCreated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateADR_EmitsUpdatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, and ADR
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	now := time.Now().UTC()
	adr, err := entities.NewADREntity("adr-1", "track-1", "ADR", "proposed", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	if err != nil {
		t.Fatalf("NewADREntity failed: %v", err)
	}
	repo.SaveADR(ctx, adr)
	mockBus.reset()

	// Update ADR
	adr.Title = "Updated ADR"
	adr.UpdatedAt = time.Now().UTC()

	err = repo.UpdateADR(ctx, adr)
	if err != nil {
		t.Fatalf("UpdateADR failed: %v", err)
	}

	// Verify update event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventADRUpdated)
	if event == nil {
		t.Fatal("expected ADRUpdated event, but none found")
	}
}

func TestEventEmittingRepository_SupersedeADR_EmitsSupersededEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, and two ADRs
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("adr-1", "track-1", "ADR 1", "accepted", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	repo.SaveADR(ctx, adr1)

	adr2, _ := entities.NewADREntity("adr-2", "track-1", "ADR 2", "proposed", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	repo.SaveADR(ctx, adr2)
	mockBus.reset()

	// Supersede ADR 1 with ADR 2
	err := repo.SupersedeADR(ctx, "adr-1", "adr-2")
	if err != nil {
		t.Fatalf("SupersedeADR failed: %v", err)
	}

	// Verify superseded event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventADRSuperseded)
	if event == nil {
		t.Fatal("expected ADRSuperseded event, but none found")
	}
}

func TestEventEmittingRepository_DeprecateADR_EmitsDeprecatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, and ADR
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	now := time.Now().UTC()
	adr, _ := entities.NewADREntity("adr-1", "track-1", "ADR", "accepted", "Context", "Decision", "Consequences", "Alternatives", now, now, nil)
	repo.SaveADR(ctx, adr)
	mockBus.reset()

	// Deprecate ADR
	err := repo.DeprecateADR(ctx, "adr-1")
	if err != nil {
		t.Fatalf("DeprecateADR failed: %v", err)
	}

	// Verify deprecated event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventADRDeprecated)
	if event == nil {
		t.Fatal("expected ADRDeprecated event, but none found")
	}
}

// ============================================================================
// Acceptance Criteria Event Tests
// ============================================================================

func TestEventEmittingRepository_SaveAC_EmitsCreatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, and task
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)
	mockBus.reset()

	// Create and save AC
	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity(
		"ac-1",
		"task-1",
		"Test AC",
		entities.VerificationTypeManual,
		"Testing instructions",
		now,
		now,
	)

	err := repo.SaveAC(ctx, ac)
	if err != nil {
		t.Fatalf("SaveAC failed: %v", err)
	}

	// Verify event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventACCreated)
	if event == nil {
		t.Fatal("expected ACCreated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateAC_EmitsUpdatedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, task, and AC
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "Instructions", now, now)
	repo.SaveAC(ctx, ac)
	mockBus.reset()

	// Update AC (no status change)
	ac.Description = "Updated AC"
	ac.UpdatedAt = time.Now().UTC()

	err := repo.UpdateAC(ctx, ac)
	if err != nil {
		t.Fatalf("UpdateAC failed: %v", err)
	}

	// Verify update event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventACUpdated)
	if event == nil {
		t.Fatal("expected ACUpdated event, but none found")
	}
}

func TestEventEmittingRepository_UpdateAC_EmitsVerifiedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, task, and AC
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "Instructions", now, now)
	repo.SaveAC(ctx, ac)
	mockBus.reset()

	// Update AC to verified status
	ac.Status = entities.ACStatusVerified
	ac.UpdatedAt = time.Now().UTC()

	err := repo.UpdateAC(ctx, ac)
	if err != nil {
		t.Fatalf("UpdateAC failed: %v", err)
	}

	// Verify verified event was emitted (not just update)
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventACVerified)
	if event == nil {
		t.Fatal("expected ACVerified event, but none found")
	}
}

func TestEventEmittingRepository_UpdateAC_EmitsFailedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, task, and AC
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "Instructions", now, now)
	repo.SaveAC(ctx, ac)
	mockBus.reset()

	// Update AC to failed status
	ac.Status = entities.ACStatusFailed
	ac.Notes = "Test failed"
	ac.UpdatedAt = time.Now().UTC()

	err := repo.UpdateAC(ctx, ac)
	if err != nil {
		t.Fatalf("UpdateAC failed: %v", err)
	}

	// Verify failed event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventACFailed)
	if event == nil {
		t.Fatal("expected ACFailed event, but none found")
	}
}

func TestEventEmittingRepository_DeleteAC_EmitsDeletedEvent(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap, track, task, and AC
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	track, _ := entities.NewTrackEntity("track-1", "roadmap-1", "Track", "", "not-started", 100, []string{}, time.Now().UTC(), time.Now().UTC())
	repo.SaveTrack(ctx, track)

	task, _ := entities.NewTaskEntity("task-1", "track-1", "Task", "", "todo", 100, "", time.Now().UTC(), time.Now().UTC())
	repo.SaveTask(ctx, task)

	now := time.Now().UTC()
	ac := entities.NewAcceptanceCriteriaEntity("ac-1", "task-1", "AC", entities.VerificationTypeManual, "Instructions", now, now)
	repo.SaveAC(ctx, ac)
	mockBus.reset()

	// Delete AC
	err := repo.DeleteAC(ctx, "ac-1")
	if err != nil {
		t.Fatalf("DeleteAC failed: %v", err)
	}

	// Verify deleted event was emitted
	if mockBus.getEventCount() != 1 {
		t.Errorf("expected 1 event, got %d", mockBus.getEventCount())
	}

	event := mockBus.getEventByType(events.EventACDeleted)
	if event == nil {
		t.Fatal("expected ACDeleted event, but none found")
	}
}

// ============================================================================
// Read-Only Operations (No Events)
// ============================================================================

func TestEventEmittingRepository_ReadOperations_DoNotEmitEvents(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, mockBus := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Setup: Create roadmap
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)
	mockBus.reset()

	// Perform read operations
	repo.GetRoadmap(ctx, "roadmap-1")
	repo.GetActiveRoadmap(ctx)
	repo.ListTracks(ctx, "roadmap-1", entities.TrackFilters{})
	repo.ListIterations(ctx)

	// Verify no events were emitted
	if mockBus.getEventCount() != 0 {
		t.Errorf("expected 0 events for read operations, got %d", mockBus.getEventCount())
	}
}

// ============================================================================
// Nil Event Bus Tests
// ============================================================================

func TestEventEmittingRepository_NilEventBus_DoesNotPanic(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	logger := createTestLogger()
	baseRepo := persistence.NewSQLiteRoadmapRepository(db, logger)

	// Create event-emitting repo with nil event bus
	eventRepo := persistence.NewEventEmittingRepository(baseRepo, nil, logger)

	ctx := context.Background()

	// Operations should work without panicking
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())

	err := eventRepo.SaveRoadmap(ctx, roadmap)
	if err != nil {
		t.Fatalf("SaveRoadmap failed with nil event bus: %v", err)
	}

	// Verify roadmap was saved
	retrieved, err := eventRepo.GetRoadmap(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("GetRoadmap failed: %v", err)
	}
	if retrieved.ID != "roadmap-1" {
		t.Errorf("expected roadmap ID roadmap-1, got %s", retrieved.ID)
	}
}

// ============================================================================
// Underlying Repository Method Tests
// ============================================================================

func TestEventEmittingRepository_DelegatesToUnderlyingRepository(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo, _ := setupEventEmittingRepo(t, db)
	ctx := context.Background()

	// Test that CRUD operations actually persist data
	roadmap, _ := entities.NewRoadmapEntity("roadmap-1", "vision", "criteria", time.Now().UTC(), time.Now().UTC())
	repo.SaveRoadmap(ctx, roadmap)

	// Verify data was persisted by retrieving it
	retrieved, err := repo.GetRoadmap(ctx, "roadmap-1")
	if err != nil {
		t.Fatalf("GetRoadmap failed: %v", err)
	}

	if retrieved.Vision != "vision" {
		t.Errorf("expected vision 'vision', got '%s'", retrieved.Vision)
	}
}
