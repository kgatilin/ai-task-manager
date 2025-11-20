package cli_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/mocks"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
)

// ============================================================================
// NewRoadmapCommands Tests
// ============================================================================

func TestNewRoadmapCommands(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	cmd := cli.NewRoadmapCommands(roadmapService)

	// Verify command structure
	if cmd.Use != "roadmap" {
		t.Errorf("Expected command name 'roadmap', got %q", cmd.Use)
	}

	// Verify subcommands exist
	subcommands := map[string]bool{}
	for _, sub := range cmd.Commands() {
		subcommands[sub.Use] = true
	}
	if !subcommands["init"] {
		t.Error("init subcommand should exist")
	}
	if !subcommands["show"] {
		t.Error("show subcommand should exist")
	}
	if !subcommands["update"] {
		t.Error("update subcommand should exist")
	}
}

// ============================================================================
// roadmap init command Tests
// ============================================================================

func TestRoadmapInitCommand_Success(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "init")

	if cmd == nil {
		t.Fatalf("Expected init command to exist")
	}

	cmd.SetContext(context.Background())

	// Capture output
	outputBuffer := &bytes.Buffer{}
	cmd.SetOut(outputBuffer)

	// Parse flags
	args := []string{
		"--vision", "Build a framework",
		"--success-criteria", "Support 10 plugins",
	}
	err := cmd.ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Execute
	err = cmd.RunE(cmd, []string{})

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := outputBuffer.String()
	if !contains(output, "Roadmap created successfully") {
		t.Error("Output should contain 'Roadmap created successfully'")
	}
	if !contains(output, "Vision") {
		t.Error("Output should contain 'Vision'")
	}
	if !contains(output, "Build a framework") {
		t.Error("Output should contain 'Build a framework'")
	}
}

func TestRoadmapInitCommand_MissingVision(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "init")

	if cmd == nil {
		t.Fatalf("Expected init command to exist")
	}

	cmd.SetContext(context.Background())

	// Parse flags without --vision
	args := []string{
		"--success-criteria", "Support 10 plugins",
	}
	err := cmd.ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Execute
	err = cmd.RunE(cmd, []string{})

	// Verify
	if err == nil {
		t.Error("Expected error for missing --vision")
	}
	if !contains(err.Error(), "--vision is required") {
		t.Errorf("Expected '--vision is required' error, got %v", err)
	}
}

func TestRoadmapInitCommand_MissingSuccessCriteria(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "init")

	if cmd == nil {
		t.Fatalf("Expected init command to exist")
	}

	cmd.SetContext(context.Background())

	// Parse flags without --success-criteria
	args := []string{
		"--vision", "Build a framework",
	}
	err := cmd.ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Execute
	err = cmd.RunE(cmd, []string{})

	// Verify
	if err == nil {
		t.Error("Expected error for missing --success-criteria")
	}
	if !contains(err.Error(), "--success-criteria is required") {
		t.Errorf("Expected '--success-criteria is required' error, got %v", err)
	}
}

// ============================================================================
// roadmap show command Tests
// ============================================================================

func TestRoadmapShowCommand_Success(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create test roadmap
	now := time.Now().UTC()
	roadmap, err := entities.NewRoadmapEntity(
		"roadmap-1",
		"Build a framework",
		"Support 10 plugins",
		now,
		now,
	)
	if err != nil {
		t.Fatalf("Failed to create test roadmap: %v", err)
	}

	// Setup mock to return our roadmap
	mockRoadmapRepo.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "show")

	if cmd == nil {
		t.Fatalf("Expected show command to exist")
	}

	cmd.SetContext(context.Background())

	// Capture output
	outputBuffer := &bytes.Buffer{}
	cmd.SetOut(outputBuffer)

	// Execute
	err = cmd.RunE(cmd, []string{})

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := outputBuffer.String()
	if !contains(output, "Roadmap:") {
		t.Error("Output should contain 'Roadmap:'")
	}
	if !contains(output, "roadmap-1") {
		t.Error("Output should contain roadmap ID")
	}
	if !contains(output, "Build a framework") {
		t.Error("Output should contain vision")
	}
}

// ============================================================================
// roadmap update command Tests
// ============================================================================

func TestRoadmapUpdateCommand_UpdateVisionSuccess(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create test roadmap
	now := time.Now().UTC()
	roadmap, err := entities.NewRoadmapEntity(
		"roadmap-1",
		"Old vision",
		"Support 10 plugins",
		now,
		now,
	)
	if err != nil {
		t.Fatalf("Failed to create test roadmap: %v", err)
	}

	// Setup mock to return our roadmap
	mockRoadmapRepo.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockRoadmapRepo.UpdateRoadmapFunc = func(ctx context.Context, r *entities.RoadmapEntity) error {
		return nil
	}

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "update")

	if cmd == nil {
		t.Fatalf("Expected update command to exist")
	}

	cmd.SetContext(context.Background())

	// Capture output
	outputBuffer := &bytes.Buffer{}
	cmd.SetOut(outputBuffer)

	// Parse flags
	args := []string{
		"--vision", "New vision",
	}
	err = cmd.ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Execute
	err = cmd.RunE(cmd, []string{})

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := outputBuffer.String()
	if !contains(output, "Roadmap updated successfully") {
		t.Error("Output should contain 'Roadmap updated successfully'")
	}
}

func TestRoadmapUpdateCommand_UpdateSuccessCriteria(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create test roadmap
	now := time.Now().UTC()
	roadmap, err := entities.NewRoadmapEntity(
		"roadmap-1",
		"Build a framework",
		"Old criteria",
		now,
		now,
	)
	if err != nil {
		t.Fatalf("Failed to create test roadmap: %v", err)
	}

	// Setup mocks
	mockRoadmapRepo.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockRoadmapRepo.UpdateRoadmapFunc = func(ctx context.Context, r *entities.RoadmapEntity) error {
		return nil
	}

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "update")

	if cmd == nil {
		t.Fatalf("Expected update command to exist")
	}

	cmd.SetContext(context.Background())

	// Capture output
	outputBuffer := &bytes.Buffer{}
	cmd.SetOut(outputBuffer)

	// Parse flags
	args := []string{
		"--success-criteria", "New criteria",
	}
	err = cmd.ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Execute
	err = cmd.RunE(cmd, []string{})

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := outputBuffer.String()
	if !contains(output, "Roadmap updated successfully") {
		t.Error("Output should contain 'Roadmap updated successfully'")
	}
}

func TestRoadmapUpdateCommand_UpdateBothFields(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	// Create test roadmap
	now := time.Now().UTC()
	roadmap, err := entities.NewRoadmapEntity(
		"roadmap-1",
		"Old vision",
		"Old criteria",
		now,
		now,
	)
	if err != nil {
		t.Fatalf("Failed to create test roadmap: %v", err)
	}

	// Setup mocks
	mockRoadmapRepo.GetActiveRoadmapFunc = func(ctx context.Context) (*entities.RoadmapEntity, error) {
		return roadmap, nil
	}

	mockRoadmapRepo.UpdateRoadmapFunc = func(ctx context.Context, r *entities.RoadmapEntity) error {
		return nil
	}

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "update")

	if cmd == nil {
		t.Fatalf("Expected update command to exist")
	}

	cmd.SetContext(context.Background())

	// Capture output
	outputBuffer := &bytes.Buffer{}
	cmd.SetOut(outputBuffer)

	// Parse flags
	args := []string{
		"--vision", "New vision",
		"--success-criteria", "New criteria",
	}
	err = cmd.ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Execute
	err = cmd.RunE(cmd, []string{})

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := outputBuffer.String()
	if !contains(output, "Roadmap updated successfully") {
		t.Error("Output should contain 'Roadmap updated successfully'")
	}
}

func TestRoadmapUpdateCommand_NoFieldsSet(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "update")

	if cmd == nil {
		t.Fatalf("Expected update command to exist")
	}

	cmd.SetContext(context.Background())

	// Execute with no flags
	err := cmd.RunE(cmd, []string{})

	// Verify
	if err == nil {
		t.Error("Expected error when no fields are specified")
	}
	if !contains(err.Error(), "at least one field must be specified to update") {
		t.Errorf("Expected error about at least one field, got %v", err)
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestRoadmapCommands_IntegrationFlow(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	rootCmd := cli.NewRoadmapCommands(roadmapService)

	// Verify root command metadata
	if rootCmd.Use != "roadmap" {
		t.Errorf("Expected command name 'roadmap', got %q", rootCmd.Use)
	}

	// Verify all subcommands are accessible
	initCmd, _, err := rootCmd.Find([]string{"init"})
	if err != nil {
		t.Fatalf("Failed to find init command: %v", err)
	}
	if initCmd == nil {
		t.Error("init command should not be nil")
	}

	showCmd, _, err := rootCmd.Find([]string{"show"})
	if err != nil {
		t.Fatalf("Failed to find show command: %v", err)
	}
	if showCmd == nil {
		t.Error("show command should not be nil")
	}

	updateCmd, _, err := rootCmd.Find([]string{"update"})
	if err != nil {
		t.Fatalf("Failed to find update command: %v", err)
	}
	if updateCmd == nil {
		t.Error("update command should not be nil")
	}
}

func TestRoadmapInitCommand_FlagValidation(t *testing.T) {
	mockRoadmapRepo := mocks.NewMockRoadmapRepository()
	mockTrackRepo := mocks.NewMockTrackRepository()
	mockTaskRepo := mocks.NewMockTaskRepository()
	mockIterationRepo := mocks.NewMockIterationRepository()
	validationSvc := services.NewValidationService()

	roadmapService := application.NewRoadmapApplicationService(
		mockRoadmapRepo,
		mockTrackRepo,
		mockTaskRepo,
		mockIterationRepo,
		validationSvc,
	)

	parentCmd := cli.NewRoadmapCommands(roadmapService)
	cmd := findCommand(parentCmd, "init")

	if cmd == nil {
		t.Fatalf("Expected init command to exist")
	}

	cmd.SetContext(context.Background())

	// Verify flags exist
	if cmd.Flags().Lookup("vision") == nil {
		t.Error("--vision flag should exist")
	}
	if cmd.Flags().Lookup("success-criteria") == nil {
		t.Error("--success-criteria flag should exist")
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
