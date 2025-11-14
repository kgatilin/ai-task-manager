package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/dto"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/mocks"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/services"
)

// setupADRTestService creates a test service with mock repositories
func setupADRTestService(t *testing.T) (*application.ADRApplicationService, context.Context, *mocks.MockADRRepository, *mocks.MockTrackRepository, *mocks.MockAggregateRepository) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}
	validationService := services.NewValidationService()

	service := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)
	ctx := context.Background()

	return service, ctx, mockADRRepo, mockTrackRepo, mockAggregateRepo
}

// createTestTrackForADRMock creates a test track entity for mock configuration
func createTestTrackForADRMock(t *testing.T, trackID string) *entities.TrackEntity {
	now := time.Now().UTC()
	roadmapID := "roadmap-" + trackID
	track, err := entities.NewTrackEntity(trackID, roadmapID, "Test Track", "Test description", "not-started", 100, []string{}, now, now)
	if err != nil {
		t.Fatalf("failed to create test track: %v", err)
	}
	return track
}

// TestADRService_CreateADR_Success tests successful ADR creation
func TestADRService_CreateADR_Success(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		if id == track.ID {
			return track, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	mockADRRepo.SaveADRFunc = func(ctx context.Context, adr *entities.ADREntity) error {
		return nil
	}

	input := dto.CreateADRDTO{
		
		TrackID:      track.ID,
		Title:        "Test ADR",
		Context:      "Test context",
		Decision:     "Test decision",
		Consequences: "Test consequences",
		Alternatives: "Test alternatives",
		Status:       "proposed",
	}

	adr, err := service.CreateADR(ctx, input)
	if err != nil {
		t.Fatalf("CreateADR() failed: %v", err)
	}

	if adr.ID == "" {
		t.Error("adr.ID should not be empty (auto-generated)")
	}
	if adr.Title != input.Title {
		t.Errorf("adr.Title = %q, want %q", adr.Title, input.Title)
	}
	if adr.Status != input.Status {
		t.Errorf("adr.Status = %q, want %q", adr.Status, input.Status)
	}
}

// TestADRService_CreateADR_InvalidID tests ADR creation with invalid ID
// NOTE: This test is now obsolete because CreateADRDTO no longer has an ID field.
// The service auto-generates IDs internally, so there's no "invalid ID" scenario for create operations.
// Keeping this test as a stub for documentation purposes.
func TestADRService_CreateADR_InvalidID(t *testing.T) {
	t.Skip("Test obsolete: CreateADRDTO no longer accepts ID field (service auto-generates)")
}

// TestADRService_CreateADR_EmptyTitle tests ADR creation with empty title
func TestADRService_CreateADR_EmptyTitle(t *testing.T) {
	service, ctx, _, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	input := dto.CreateADRDTO{
		
		TrackID:      track.ID,
		Title:        "", // Empty title
		Context:      "Test context",
		Decision:     "Test decision",
		Consequences: "Test consequences",
	}

	_, err := service.CreateADR(ctx, input)
	if err == nil {
		t.Fatal("CreateADR() should fail with empty title")
	}
}

// TestADRService_CreateADR_TrackNotFound tests ADR creation with non-existent track
func TestADRService_CreateADR_TrackNotFound(t *testing.T) {
	service, ctx, _, mockTrackRepo, _ := setupADRTestService(t)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	input := dto.CreateADRDTO{
		
		TrackID:      "nonexistent",
		Title:        "Test ADR",
		Context:      "Test context",
		Decision:     "Test decision",
		Consequences: "Test consequences",
	}

	_, err := service.CreateADR(ctx, input)
	if err == nil {
		t.Fatal("CreateADR() should fail with non-existent track")
	}
}

// TestADRService_CreateADR_DefaultStatus tests ADR creation with default status
func TestADRService_CreateADR_DefaultStatus(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockADRRepo.SaveADRFunc = func(ctx context.Context, adr *entities.ADREntity) error {
		return nil
	}

	input := dto.CreateADRDTO{
		
		TrackID:      track.ID,
		Title:        "Test ADR",
		Context:      "Test context",
		Decision:     "Test decision",
		Consequences: "Test consequences",
		Status:       "", // Empty status should default to proposed
	}

	adr, err := service.CreateADR(ctx, input)
	if err != nil {
		t.Fatalf("CreateADR() failed: %v", err)
	}

	if adr.Status != "proposed" {
		t.Errorf("adr.Status = %q, want %q", adr.Status, "proposed")
	}
}

// TestADRService_UpdateADR_Success tests successful ADR update
func TestADRService_UpdateADR_Success(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	existingADR, _ := entities.NewADREntity("TM-adr-1", track.ID, "Original Title", "proposed", "Original context", "Original decision", "Original consequences", "", now, now, nil)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		if id == existingADR.ID {
			return existingADR, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	mockADRRepo.UpdateADRFunc = func(ctx context.Context, adr *entities.ADREntity) error {
		return nil
	}

	// Update ADR
	newTitle := "Updated Title"
	newContext := "Updated context"
	updateInput := dto.UpdateADRDTO{
		ID:      existingADR.ID, // MUST set ID for update operations
		Title:   &newTitle,
		Context: &newContext,
	}

	adr, err := service.UpdateADR(ctx, updateInput)
	if err != nil {
		t.Fatalf("UpdateADR() failed: %v", err)
	}

	if adr.Title != newTitle {
		t.Errorf("adr.Title = %q, want %q", adr.Title, newTitle)
	}
	if adr.Context != newContext {
		t.Errorf("adr.Context = %q, want %q", adr.Context, newContext)
	}
}

// TestADRService_UpdateADR_NotFound tests updating non-existent ADR
func TestADRService_UpdateADR_NotFound(t *testing.T) {
	service, ctx, mockADRRepo, _, _ := setupADRTestService(t)

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	newTitle := "Updated Title"
	updateInput := dto.UpdateADRDTO{
		ID:    "nonexistent",
		Title: &newTitle,
	}

	_, err := service.UpdateADR(ctx, updateInput)
	if err == nil {
		t.Fatal("UpdateADR() should fail for non-existent ADR")
	}
}

// TestADRService_UpdateADR_PartialUpdate tests partial ADR update
func TestADRService_UpdateADR_PartialUpdate(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	existingADR, _ := entities.NewADREntity("TM-adr-1", track.ID, "Original Title", "proposed", "Original context", "Original decision", "Original consequences", "", now, now, nil)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		if id == existingADR.ID {
			return existingADR, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	mockADRRepo.UpdateADRFunc = func(ctx context.Context, adr *entities.ADREntity) error {
		return nil
	}

	// Update only title
	newTitle := "Updated Title"
	updateInput := dto.UpdateADRDTO{
		ID:    existingADR.ID, // MUST set ID for update operations
		Title: &newTitle,
	}

	adr, err := service.UpdateADR(ctx, updateInput)
	if err != nil {
		t.Fatalf("UpdateADR() failed: %v", err)
	}

	if adr.Title != newTitle {
		t.Errorf("adr.Title = %q, want %q", adr.Title, newTitle)
	}
	// Other fields should remain unchanged
	if adr.Context != existingADR.Context {
		t.Errorf("adr.Context changed: got %q, want %q", adr.Context, existingADR.Context)
	}
	if adr.Decision != existingADR.Decision {
		t.Errorf("adr.Decision changed: got %q, want %q", adr.Decision, existingADR.Decision)
	}
}

// TestADRService_SupersedeADR_Success tests successful ADR supersession
func TestADRService_SupersedeADR_Success(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("TM-adr-1", track.ID, "Original ADR", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)
	adr2, _ := entities.NewADREntity("TM-adr-2", track.ID, "New ADR", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		if id == adr1.ID {
			return adr1, nil
		}
		if id == adr2.ID {
			return adr2, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	mockADRRepo.UpdateADRFunc = func(ctx context.Context, adr *entities.ADREntity) error {
		adr1 = adr // Update reference
		return nil
	}

	// Supersede ADR-1 with ADR-2
	err := service.SupersedeADR(ctx, "TM-adr-1", "TM-adr-2")
	if err != nil {
		t.Fatalf("SupersedeADR() failed: %v", err)
	}

	// Verify status changed
	adr, err := service.GetADR(ctx, "TM-adr-1")
	if err != nil {
		t.Fatalf("GetADR() failed: %v", err)
	}

	if adr.Status != "superseded" {
		t.Errorf("adr.Status = %q, want %q", adr.Status, "superseded")
	}
	if adr.SupersededBy == nil || *adr.SupersededBy != "TM-adr-2" {
		t.Errorf("adr.SupersededBy = %v, want %q", adr.SupersededBy, "TM-adr-2")
	}
}

// TestADRService_SupersedeADR_ADRNotFound tests superseding non-existent ADR
func TestADRService_SupersedeADR_ADRNotFound(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("TM-adr-1", track.ID, "Test ADR", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		if id == adr1.ID {
			return adr1, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	// Try to supersede non-existent ADR
	err := service.SupersedeADR(ctx, "nonexistent", "TM-adr-1")
	if err == nil {
		t.Fatal("SupersedeADR() should fail with non-existent ADR")
	}
}

// TestADRService_SupersedeADR_SupersededByNotFound tests superseding with non-existent superseding ADR
func TestADRService_SupersedeADR_SupersededByNotFound(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("TM-adr-1", track.ID, "Test ADR", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		if id == adr1.ID {
			return adr1, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	// Try to supersede with non-existent superseding ADR
	err := service.SupersedeADR(ctx, "TM-adr-1", "nonexistent")
	if err == nil {
		t.Fatal("SupersedeADR() should fail with non-existent superseding ADR")
	}
}

// TestADRService_DeprecateADR_Success tests successful ADR deprecation
func TestADRService_DeprecateADR_Success(t *testing.T) {
	service, ctx, mockADRRepo, mockTrackRepo, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	existingADR, _ := entities.NewADREntity("TM-adr-1", track.ID, "Test ADR", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockTrackRepo.GetTrackFunc = func(ctx context.Context, id string) (*entities.TrackEntity, error) {
		return track, nil
	}

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		if id == existingADR.ID {
			return existingADR, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	mockADRRepo.UpdateADRFunc = func(ctx context.Context, adr *entities.ADREntity) error {
		existingADR = adr // Update reference
		return nil
	}

	// Deprecate ADR
	err := service.DeprecateADR(ctx, "TM-adr-1")
	if err != nil {
		t.Fatalf("DeprecateADR() failed: %v", err)
	}

	// Verify status changed
	adr, err := service.GetADR(ctx, "TM-adr-1")
	if err != nil {
		t.Fatalf("GetADR() failed: %v", err)
	}

	if adr.Status != "deprecated" {
		t.Errorf("adr.Status = %q, want %q", adr.Status, "deprecated")
	}
}

// TestADRService_DeprecateADR_NotFound tests deprecating non-existent ADR
func TestADRService_DeprecateADR_NotFound(t *testing.T) {
	service, ctx, mockADRRepo, _, _ := setupADRTestService(t)

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	err := service.DeprecateADR(ctx, "nonexistent")
	if err == nil {
		t.Fatal("DeprecateADR() should fail for non-existent ADR")
	}
}

// TestADRService_GetADR_Success tests successful ADR retrieval
func TestADRService_GetADR_Success(t *testing.T) {
	service, ctx, mockADRRepo, _, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	existingADR, _ := entities.NewADREntity("TM-adr-1", track.ID, "Test ADR", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		if id == existingADR.ID {
			return existingADR, nil
		}
		return nil, pluginsdk.ErrNotFound
	}

	// Get ADR
	adr, err := service.GetADR(ctx, "TM-adr-1")
	if err != nil {
		t.Fatalf("GetADR() failed: %v", err)
	}

	if adr.ID != existingADR.ID {
		t.Errorf("adr.ID = %q, want %q", adr.ID, existingADR.ID)
	}
	if adr.Title != existingADR.Title {
		t.Errorf("adr.Title = %q, want %q", adr.Title, existingADR.Title)
	}
}

// TestADRService_GetADR_NotFound tests retrieving non-existent ADR
func TestADRService_GetADR_NotFound(t *testing.T) {
	service, ctx, mockADRRepo, _, _ := setupADRTestService(t)

	mockADRRepo.GetADRFunc = func(ctx context.Context, id string) (*entities.ADREntity, error) {
		return nil, pluginsdk.ErrNotFound
	}

	_, err := service.GetADR(ctx, "nonexistent")
	if err == nil {
		t.Fatal("GetADR() should fail for non-existent ADR")
	}
}

// TestADRService_ListADRs_Success tests successful ADR listing
func TestADRService_ListADRs_Success(t *testing.T) {
	service, ctx, mockADRRepo, _, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("TM-adr-1", track.ID, "Test ADR 1", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)
	adr2, _ := entities.NewADREntity("TM-adr-2", track.ID, "Test ADR 2", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)
	adr3, _ := entities.NewADREntity("TM-adr-3", track.ID, "Test ADR 3", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockADRRepo.ListADRsFunc = func(ctx context.Context, trackID *string) ([]*entities.ADREntity, error) {
		return []*entities.ADREntity{adr1, adr2, adr3}, nil
	}

	// List all ADRs
	adrs, err := service.ListADRs(ctx, nil)
	if err != nil {
		t.Fatalf("ListADRs() failed: %v", err)
	}

	if len(adrs) != 3 {
		t.Fatalf("ListADRs() returned %d ADRs, want 3", len(adrs))
	}
}

// TestADRService_ListADRs_WithTrackFilter tests listing ADRs with track filter
func TestADRService_ListADRs_WithTrackFilter(t *testing.T) {
	service, ctx, mockADRRepo, _, _ := setupADRTestService(t)
	track1 := createTestTrackForADRMock(t, "TM-track-1")
	track2 := createTestTrackForADRMock(t, "TM-track-2")

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("TM-adr-1", track1.ID, "Test ADR 1", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)
	adr2, _ := entities.NewADREntity("TM-adr-2", track2.ID, "Test ADR 2", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockADRRepo.ListADRsFunc = func(ctx context.Context, trackID *string) ([]*entities.ADREntity, error) {
		if trackID != nil && *trackID == track1.ID {
			return []*entities.ADREntity{adr1}, nil
		}
		if trackID != nil && *trackID == track2.ID {
			return []*entities.ADREntity{adr2}, nil
		}
		return []*entities.ADREntity{adr1, adr2}, nil
	}

	// List ADRs for track 1
	adrs, err := service.ListADRs(ctx, &track1.ID)
	if err != nil {
		t.Fatalf("ListADRs() failed: %v", err)
	}

	if len(adrs) != 1 {
		t.Fatalf("ListADRs() returned %d ADRs, want 1", len(adrs))
	}
	if adrs[0].TrackID != track1.ID {
		t.Errorf("adrs[0].TrackID = %q, want %q", adrs[0].TrackID, track1.ID)
	}
}

// TestADRService_GetADRsByTrack_Success tests retrieving ADRs by track
func TestADRService_GetADRsByTrack_Success(t *testing.T) {
	service, ctx, mockADRRepo, _, _ := setupADRTestService(t)
	track := createTestTrackForADRMock(t, "TM-track-1")

	now := time.Now().UTC()
	adr1, _ := entities.NewADREntity("TM-adr-1", track.ID, "Test ADR", "proposed", "Test context", "Test decision", "Test consequences", "", now, now, nil)

	mockADRRepo.GetADRsByTrackFunc = func(ctx context.Context, trackID string) ([]*entities.ADREntity, error) {
		if trackID == track.ID {
			return []*entities.ADREntity{adr1}, nil
		}
		return []*entities.ADREntity{}, nil
	}

	// Get ADRs by track
	adrs, err := service.GetADRsByTrack(ctx, track.ID)
	if err != nil {
		t.Fatalf("GetADRsByTrack() failed: %v", err)
	}

	if len(adrs) != 1 {
		t.Fatalf("GetADRsByTrack() returned %d ADRs, want 1", len(adrs))
	}
	if adrs[0].TrackID != track.ID {
		t.Errorf("adrs[0].TrackID = %q, want %q", adrs[0].TrackID, track.ID)
	}
}
