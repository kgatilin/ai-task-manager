package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/mocks"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
)

// ============================================================================
// InitRoadmap Tests
// ============================================================================

func TestRoadmapApplicationService_InitRoadmap_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	input := dto.CreateRoadmapDTO{
		Vision:          "Build extensible framework",
		SuccessCriteria: "Support 10 plugins",
	}

	roadmap, err := service.InitRoadmap(ctx, input)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if roadmap == nil {
		t.Fatal("Expected roadmap to be returned")
	}
	if roadmap.Vision != input.Vision {
		t.Errorf("Expected vision %q, got %q", input.Vision, roadmap.Vision)
	}
	if roadmap.SuccessCriteria != input.SuccessCriteria {
		t.Errorf("Expected success criteria %q, got %q", input.SuccessCriteria, roadmap.SuccessCriteria)
	}
	if roadmap.ID == "" {
		t.Error("Expected roadmap ID to be generated")
	}
}

func TestRoadmapApplicationService_InitRoadmap_EmptyVision(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	input := dto.CreateRoadmapDTO{
		Vision:          "",
		SuccessCriteria: "Support 10 plugins",
	}

	_, err := service.InitRoadmap(ctx, input)

	// Verify
	if err == nil {
		t.Fatal("Expected error for empty vision")
	}
	if !errors.Is(err, tmerrors.ErrInvalidArgument) {
		t.Errorf("Expected ErrInvalidArgument, got %v", err)
	}
}

func TestRoadmapApplicationService_InitRoadmap_EmptySuccessCriteria(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	input := dto.CreateRoadmapDTO{
		Vision:          "Build extensible framework",
		SuccessCriteria: "",
	}

	_, err := service.InitRoadmap(ctx, input)

	// Verify
	if err == nil {
		t.Fatal("Expected error for empty success criteria")
	}
	if !errors.Is(err, tmerrors.ErrInvalidArgument) {
		t.Errorf("Expected ErrInvalidArgument, got %v", err)
	}
}

func TestRoadmapApplicationService_InitRoadmap_RoadmapAlreadyExists(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create existing roadmap
	existingRoadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Existing vision",
		"Existing criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, existingRoadmap)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	input := dto.CreateRoadmapDTO{
		Vision:          "Build extensible framework",
		SuccessCriteria: "Support 10 plugins",
	}

	_, err := service.InitRoadmap(ctx, input)

	// Verify
	if err == nil {
		t.Fatal("Expected error when roadmap already exists")
	}
	if !errors.Is(err, tmerrors.ErrAlreadyExists) && err.Error() == "" {
		t.Errorf("Expected error about existing roadmap, got %v", err)
	}
}

// ============================================================================
// GetRoadmap Tests
// ============================================================================

func TestRoadmapApplicationService_GetRoadmap_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create existing roadmap
	existingRoadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Existing vision",
		"Existing criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, existingRoadmap)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	roadmap, err := service.GetRoadmap(ctx)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if roadmap == nil {
		t.Fatal("Expected roadmap to be returned")
	}
	if roadmap.ID != existingRoadmap.ID {
		t.Errorf("Expected roadmap ID %q, got %q", existingRoadmap.ID, roadmap.ID)
	}
}

func TestRoadmapApplicationService_GetRoadmap_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	_, err := service.GetRoadmap(ctx)

	// Verify
	if err == nil {
		t.Fatal("Expected error when roadmap not found")
	}
	if !errors.Is(err, tmerrors.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// ============================================================================
// UpdateRoadmap Tests
// ============================================================================

func TestRoadmapApplicationService_UpdateRoadmap_VisionOnly(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create existing roadmap
	existingRoadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Old vision",
		"Old criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, existingRoadmap)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	newVision := "New vision"
	input := dto.UpdateRoadmapDTO{
		Vision: &newVision,
	}

	roadmap, err := service.UpdateRoadmap(ctx, input)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if roadmap.Vision != newVision {
		t.Errorf("Expected vision %q, got %q", newVision, roadmap.Vision)
	}
	if roadmap.SuccessCriteria != existingRoadmap.SuccessCriteria {
		t.Errorf("Expected success criteria to remain unchanged")
	}
}

func TestRoadmapApplicationService_UpdateRoadmap_SuccessCriteriaOnly(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create existing roadmap
	existingRoadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Old vision",
		"Old criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, existingRoadmap)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	newCriteria := "New success criteria"
	input := dto.UpdateRoadmapDTO{
		SuccessCriteria: &newCriteria,
	}

	roadmap, err := service.UpdateRoadmap(ctx, input)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if roadmap.SuccessCriteria != newCriteria {
		t.Errorf("Expected success criteria %q, got %q", newCriteria, roadmap.SuccessCriteria)
	}
	if roadmap.Vision != existingRoadmap.Vision {
		t.Errorf("Expected vision to remain unchanged")
	}
}

func TestRoadmapApplicationService_UpdateRoadmap_BothFields(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create existing roadmap
	existingRoadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Old vision",
		"Old criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, existingRoadmap)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	newVision := "New vision"
	newCriteria := "New success criteria"
	input := dto.UpdateRoadmapDTO{
		Vision:          &newVision,
		SuccessCriteria: &newCriteria,
	}

	roadmap, err := service.UpdateRoadmap(ctx, input)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if roadmap.Vision != newVision {
		t.Errorf("Expected vision %q, got %q", newVision, roadmap.Vision)
	}
	if roadmap.SuccessCriteria != newCriteria {
		t.Errorf("Expected success criteria %q, got %q", newCriteria, roadmap.SuccessCriteria)
	}
}

func TestRoadmapApplicationService_UpdateRoadmap_NoFields(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	input := dto.UpdateRoadmapDTO{}

	_, err := service.UpdateRoadmap(ctx, input)

	// Verify
	if err == nil {
		t.Fatal("Expected error when no fields provided")
	}
}

func TestRoadmapApplicationService_UpdateRoadmap_EmptyVision(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create existing roadmap
	existingRoadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Old vision",
		"Old criteria",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, existingRoadmap)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	emptyVision := ""
	input := dto.UpdateRoadmapDTO{
		Vision: &emptyVision,
	}

	_, err := service.UpdateRoadmap(ctx, input)

	// Verify
	if err == nil {
		t.Fatal("Expected error for empty vision")
	}
	if !errors.Is(err, tmerrors.ErrInvalidArgument) {
		t.Errorf("Expected ErrInvalidArgument, got %v", err)
	}
}

func TestRoadmapApplicationService_UpdateRoadmap_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	newVision := "New vision"
	input := dto.UpdateRoadmapDTO{
		Vision: &newVision,
	}

	_, err := service.UpdateRoadmap(ctx, input)

	// Verify
	if err == nil {
		t.Fatal("Expected error when roadmap not found")
	}
	if !errors.Is(err, tmerrors.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// ============================================================================
// GetFullOverview Tests
// ============================================================================

func TestRoadmapApplicationService_GetFullOverview_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create roadmap
	roadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Build extensible framework",
		"Support 10 plugins",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, roadmap)

	// Create track
	track, _ := entities.NewTrackEntity(
		"TM-track-1",
		"roadmap-1",
		"Plugin System",
		"Implement plugin architecture",
		"in-progress",
		100,
		[]string{},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockTrackRepo.SaveTrack(ctx, track)

	// Create task
	task, _ := entities.NewTaskEntity(
		"TM-task-1",
		"TM-track-1",
		"Create plugin SDK",
		"Define plugin interface",
		"todo",
		100,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockTaskRepo.SaveTask(ctx, task)

	// Create iteration
	now := time.Now().UTC()
	zero := time.Time{}
	iteration, _ := entities.NewIterationEntity(
		1,
		"Sprint 1",
		"Complete plugin SDK",
		"Plugin SDK ready for use",
		[]string{"TM-task-1"},
		"planned",
		100,
		zero, // startedAt
		zero, // completedAt
		now,  // createdAt
		now,  // updatedAt
	)
	mockIterationRepo.SaveIteration(ctx, iteration)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	options := dto.RoadmapOverviewOptions{
		Verbose:  false,
		Sections: nil,
	}

	overview, err := service.GetFullOverview(ctx, options)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if overview == nil {
		t.Fatal("Expected overview to be returned")
	}
	if overview.Roadmap == nil {
		t.Error("Expected roadmap in overview")
	}
	if len(overview.Tracks) != 1 {
		t.Errorf("Expected 1 track, got %d", len(overview.Tracks))
	}
	if len(overview.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(overview.Tasks))
	}
	if len(overview.Iterations) != 1 {
		t.Errorf("Expected 1 iteration, got %d", len(overview.Iterations))
	}
}

func TestRoadmapApplicationService_GetFullOverview_RoadmapNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	options := dto.RoadmapOverviewOptions{
		Verbose:  false,
		Sections: nil,
	}

	_, err := service.GetFullOverview(ctx, options)

	// Verify
	if err == nil {
		t.Fatal("Expected error when roadmap not found")
	}
	if !errors.Is(err, tmerrors.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestRoadmapApplicationService_GetFullOverview_EmptyData(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create roadmap only (no tracks, tasks, iterations)
	roadmap, _ := entities.NewRoadmapEntity(
		"roadmap-1",
		"Build extensible framework",
		"Support 10 plugins",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	mockRoadmapRepo.SaveRoadmap(ctx, roadmap)

	service := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	// Test
	options := dto.RoadmapOverviewOptions{
		Verbose:  false,
		Sections: nil,
	}

	overview, err := service.GetFullOverview(ctx, options)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if overview == nil {
		t.Fatal("Expected overview to be returned")
	}
	if len(overview.Tracks) != 0 {
		t.Errorf("Expected 0 tracks, got %d", len(overview.Tracks))
	}
	if len(overview.Tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(overview.Tasks))
	}
	if len(overview.Iterations) != 0 {
		t.Errorf("Expected 0 iterations, got %d", len(overview.Iterations))
	}
}
