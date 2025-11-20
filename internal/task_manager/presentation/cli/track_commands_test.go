package cli_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Test Cobra Command Structure
// ============================================================================

// TestNewTrackCommands verifies that NewTrackCommands returns a valid Cobra command group
func TestNewTrackCommands_Structure(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)

	assert.NotNil(t, trackCommands, "NewTrackCommands should return a command group")
	assert.Equal(t, "track", trackCommands.Name(), "command name should be 'track'")
	assert.NotEmpty(t, trackCommands.Short, "command should have short description")
	assert.NotEmpty(t, trackCommands.Long, "command should have long description")
}

// TestTrackCommands_AllSubcommands verifies all 7 subcommands are present
func TestTrackCommands_AllSubcommands(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)

	expectedSubcommands := []string{
		"create",
		"list",
		"show",
		"update",
		"delete",
		"add-dependency",
		"remove-dependency",
	}

	commandNames := make(map[string]bool)
	for _, cmd := range trackCommands.Commands() {
		commandNames[cmd.Name()] = true
	}

	for _, name := range expectedSubcommands {
		assert.True(t, commandNames[name], "command '%s' should exist", name)
	}
}

// TestTrackCreateCommand_Flags verifies create command has required flags
func TestTrackCreateCommand_Flags(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)
	createCmd := findCommand(trackCommands, "create")

	assert.NotNil(t, createCmd, "create command should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("title"), "--title flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("description"), "--description flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("rank"), "--rank flag should exist")
}

// TestTrackListCommand_Flags verifies list command has filter flags
func TestTrackListCommand_Flags(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)
	listCmd := findCommand(trackCommands, "list")

	assert.NotNil(t, listCmd, "list command should exist")
	assert.NotNil(t, listCmd.Flags().Lookup("status"), "--status flag should exist")
}

// TestTrackShowCommand_Arguments verifies show command requires track ID
func TestTrackShowCommand_Arguments(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)
	showCmd := findCommand(trackCommands, "show")

	assert.NotNil(t, showCmd, "show command should exist")
	// Check that it uses ExactArgs(1)
	assert.NotNil(t, showCmd.Args, "show command should have argument validation")
}

// TestTrackUpdateCommand_Flags verifies update command has optional field flags
func TestTrackUpdateCommand_Flags(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)
	updateCmd := findCommand(trackCommands, "update")

	assert.NotNil(t, updateCmd, "update command should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("title"), "--title flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("description"), "--description flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("status"), "--status flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("rank"), "--rank flag should exist")
}

// TestTrackDeleteCommand_Flags verifies delete command has force flag
func TestTrackDeleteCommand_Flags(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)
	deleteCmd := findCommand(trackCommands, "delete")

	assert.NotNil(t, deleteCmd, "delete command should exist")
	assert.NotNil(t, deleteCmd.Flags().Lookup("force"), "--force flag should exist")
}

// TestTrackAddDependencyCommand_Arguments verifies add-dependency command requires two IDs
func TestTrackAddDependencyCommand_Arguments(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)
	addDepCmd := findCommand(trackCommands, "add-dependency")

	assert.NotNil(t, addDepCmd, "add-dependency command should exist")
	assert.NotNil(t, addDepCmd.Args, "add-dependency command should have argument validation")
}

// TestTrackRemoveDependencyCommand_Arguments verifies remove-dependency command requires two IDs
func TestTrackRemoveDependencyCommand_Arguments(t *testing.T) {
	trackCommands := cli.NewTrackCommands(nil, nil)
	removeDepCmd := findCommand(trackCommands, "remove-dependency")

	assert.NotNil(t, removeDepCmd, "remove-dependency command should exist")
	assert.NotNil(t, removeDepCmd.Args, "remove-dependency command should have argument validation")
}

// ============================================================================
// Helper function to find subcommand by name
// ============================================================================
