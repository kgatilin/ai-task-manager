package cli_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Test Cobra Command Structure
// ============================================================================

// TestNewACCommands verifies that NewACCommands returns a valid Cobra command group
func TestNewACCommands_Structure(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)

	assert.NotNil(t, acCommands, "NewACCommands should return a command group")
	assert.Equal(t, "ac", acCommands.Name(), "command name should be 'ac'")
	assert.NotEmpty(t, acCommands.Short, "command should have short description")
	assert.NotEmpty(t, acCommands.Long, "command should have long description")
}

// TestACCommands_AllSubcommands verifies all 9 subcommands are present
func TestACCommands_AllSubcommands(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)

	expectedSubcommands := []string{
		"add",
		"list",
		"list-iteration",
		"show",
		"update",
		"verify",
		"fail",
		"failed",
		"delete",
	}

	commandNames := make(map[string]bool)
	for _, cmd := range acCommands.Commands() {
		commandNames[cmd.Name()] = true
	}

	for _, name := range expectedSubcommands {
		assert.True(t, commandNames[name], "command '%s' should exist", name)
	}
}

// TestACAddCommand_Flags verifies add command has required flags
func TestACAddCommand_Flags(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	addCmd := findCommand(acCommands, "add")

	assert.NotNil(t, addCmd, "add command should exist")
	assert.NotNil(t, addCmd.Flags().Lookup("description"), "--description flag should exist")
	assert.NotNil(t, addCmd.Flags().Lookup("testing-instructions"), "--testing-instructions flag should exist")
}

// TestACListCommand_Arguments verifies list command requires task ID
func TestACListCommand_Arguments(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	listCmd := findCommand(acCommands, "list")

	assert.NotNil(t, listCmd, "list command should exist")
	assert.NotNil(t, listCmd.Args, "list command should have argument validation")
}

// TestACShowCommand_Arguments verifies show command requires AC ID
func TestACShowCommand_Arguments(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	showCmd := findCommand(acCommands, "show")

	assert.NotNil(t, showCmd, "show command should exist")
	assert.NotNil(t, showCmd.Args, "show command should have argument validation")
}

// TestACUpdateCommand_Flags verifies update command has optional field flags
func TestACUpdateCommand_Flags(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	updateCmd := findCommand(acCommands, "update")

	assert.NotNil(t, updateCmd, "update command should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("description"), "--description flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("testing-instructions"), "--testing-instructions flag should exist")
}

// TestACVerifyCommand_Arguments verifies verify command requires AC ID
func TestACVerifyCommand_Arguments(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	verifyCmd := findCommand(acCommands, "verify")

	assert.NotNil(t, verifyCmd, "verify command should exist")
	assert.NotNil(t, verifyCmd.Args, "verify command should have argument validation")
}

// TestACFailCommand_Flags verifies fail command has required feedback flag
func TestACFailCommand_Flags(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	failCmd := findCommand(acCommands, "fail")

	assert.NotNil(t, failCmd, "fail command should exist")
	assert.NotNil(t, failCmd.Flags().Lookup("feedback"), "--feedback flag should exist")
}

// TestACFailedCommand_Flags verifies failed command has optional filter flags
func TestACFailedCommand_Flags(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	failedCmd := findCommand(acCommands, "failed")

	assert.NotNil(t, failedCmd, "failed command should exist")
	assert.NotNil(t, failedCmd.Flags().Lookup("iteration"), "--iteration flag should exist")
	assert.NotNil(t, failedCmd.Flags().Lookup("track"), "--track flag should exist")
	assert.NotNil(t, failedCmd.Flags().Lookup("task"), "--task flag should exist")
}

// TestACDeleteCommand_Flags verifies delete command has force flag
func TestACDeleteCommand_Flags(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	deleteCmd := findCommand(acCommands, "delete")

	assert.NotNil(t, deleteCmd, "delete command should exist")
	assert.NotNil(t, deleteCmd.Flags().Lookup("force"), "--force flag should exist")
}

// TestACListIterationCommand_Arguments verifies list-iteration command requires iteration number
func TestACListIterationCommand_Arguments(t *testing.T) {
	acCommands := cli.NewACCommands(nil, nil)
	listIterCmd := findCommand(acCommands, "list-iteration")

	assert.NotNil(t, listIterCmd, "list-iteration command should exist")
	assert.NotNil(t, listIterCmd.Args, "list-iteration command should have argument validation")
}

// ============================================================================
// Helper function
// ============================================================================

// findCommand is a helper function to find a subcommand by name
