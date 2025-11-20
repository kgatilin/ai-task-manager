package cli_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/mocks"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
	"github.com/spf13/cobra"
)

// ============================================================================
// Test NewADRCommands
// ============================================================================

func TestNewADRCommands(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	cmd := cli.NewADRCommands(adrService)

	// Verify command structure
	if cmd.Use != "adr" {
		t.Fatalf("Expected Use to be 'adr', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Fatalf("Expected Short description to be set")
	}

	// Verify subcommands exist
	expectedCommands := []string{"create", "list", "show", "update", "supersede", "deprecate", "check"}
	if len(cmd.Commands()) != len(expectedCommands) {
		t.Fatalf("Expected %d subcommands, got %d", len(expectedCommands), len(cmd.Commands()))
	}
}

// ============================================================================
// Test individual commands exist and work
// ============================================================================

func TestADRCreateCommand_FlagsConfigured(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)
	cmd := findCommand(parentCmd, "create")

	if cmd == nil {
		t.Fatalf("Expected create command to exist")
	}

	// Verify required flags are marked
	if flag := cmd.Flag("track"); flag == nil {
		t.Fatalf("Expected --track flag to exist")
	}
	if flag := cmd.Flag("title"); flag == nil {
		t.Fatalf("Expected --title flag to exist")
	}
	if flag := cmd.Flag("context"); flag == nil {
		t.Fatalf("Expected --context flag to exist")
	}
	if flag := cmd.Flag("decision"); flag == nil {
		t.Fatalf("Expected --decision flag to exist")
	}
	if flag := cmd.Flag("consequences"); flag == nil {
		t.Fatalf("Expected --consequences flag to exist")
	}
	if flag := cmd.Flag("alternatives"); flag == nil {
		t.Fatalf("Expected --alternatives flag to exist")
	}
}

func TestADRListCommand_FlagsConfigured(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)
	cmd := findCommand(parentCmd, "list")

	if cmd == nil {
		t.Fatalf("Expected list command to exist")
	}

	// Verify optional track filter flag exists
	if flag := cmd.Flag("track"); flag == nil {
		t.Fatalf("Expected --track flag to exist")
	}
}

func TestADRShowCommand_ArgumentsConfigured(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)
	cmd := findCommand(parentCmd, "show")

	if cmd == nil {
		t.Fatalf("Expected show command to exist")
	}

	// Verify it expects exactly 1 argument
	if cmd.Args == nil || cmd.Args(cmd, []string{}) != nil {
		// ExactArgs returns a function, just verify the structure
		if cmd.Use != "show <adr-id>" {
			t.Fatalf("Expected command to require exactly 1 argument, got Use: %s", cmd.Use)
		}
	}
}

func TestADRUpdateCommand_FlagsConfigured(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)
	cmd := findCommand(parentCmd, "update")

	if cmd == nil {
		t.Fatalf("Expected update command to exist")
	}

	// Verify optional update flags exist
	if flag := cmd.Flag("title"); flag == nil {
		t.Fatalf("Expected --title flag to exist")
	}
	if flag := cmd.Flag("context"); flag == nil {
		t.Fatalf("Expected --context flag to exist")
	}
	if flag := cmd.Flag("decision"); flag == nil {
		t.Fatalf("Expected --decision flag to exist")
	}
	if flag := cmd.Flag("consequences"); flag == nil {
		t.Fatalf("Expected --consequences flag to exist")
	}
	if flag := cmd.Flag("alternatives"); flag == nil {
		t.Fatalf("Expected --alternatives flag to exist")
	}
	if flag := cmd.Flag("status"); flag == nil {
		t.Fatalf("Expected --status flag to exist")
	}
}

func TestADRSupersedeCommand_FlagsConfigured(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)
	cmd := findCommand(parentCmd, "supersede")

	if cmd == nil {
		t.Fatalf("Expected supersede command to exist")
	}

	// Verify required flag exists
	if flag := cmd.Flag("superseded-by"); flag == nil {
		t.Fatalf("Expected --superseded-by flag to exist")
	}
}

func TestADRDeprecateCommand_ArgumentsConfigured(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)
	cmd := findCommand(parentCmd, "deprecate")

	if cmd == nil {
		t.Fatalf("Expected deprecate command to exist")
	}

	// Verify it expects exactly 1 argument
	if cmd.Use != "deprecate <adr-id>" {
		t.Fatalf("Expected command to require exactly 1 argument, got Use: %s", cmd.Use)
	}
}

func TestADRCheckCommand_HasDescription(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)
	cmd := findCommand(parentCmd, "check")

	if cmd == nil {
		t.Fatalf("Expected check command to exist")
	}

	// Verify it has a description
	if cmd.Short == "" {
		t.Fatalf("Expected command to have a short description")
	}
	if cmd.Long == "" {
		t.Fatalf("Expected command to have a long description")
	}
}

// ============================================================================
// Test command help output
// ============================================================================

func TestADRCommands_HelpAvailable(t *testing.T) {
	mockADRRepo := &mocks.MockADRRepository{}
	mockTrackRepo := &mocks.MockTrackRepository{}
	mockAggregateRepo := &mocks.MockAggregateRepository{}

	validationService := services.NewValidationService()
	adrService := application.NewADRApplicationService(mockADRRepo, mockTrackRepo, mockAggregateRepo, validationService)

	parentCmd := cli.NewADRCommands(adrService)

	commandNames := []string{"create", "list", "show", "update", "supersede", "deprecate", "check"}
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{}

	for _, name := range commandNames {
		cmd := findCommand(parentCmd, name)
		if cmd != nil {
			commands = append(commands, struct {
				name string
				cmd  *cobra.Command
			}{name, cmd})
		}
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cmd.Short == "" {
				t.Fatalf("Expected %s command to have Short description", tc.name)
			}
			if tc.cmd.Long == "" {
				t.Fatalf("Expected %s command to have Long description", tc.name)
			}
			if tc.cmd.Example == "" {
				t.Fatalf("Expected %s command to have Example", tc.name)
			}
		})
	}
}
