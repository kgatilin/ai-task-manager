package cli_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Test Cobra Command Structure
// ============================================================================

// TestNewIterationCommands verifies that NewIterationCommands returns a valid Cobra command group
func TestNewIterationCommands_Structure(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)

	assert.NotNil(t, iterationCommands, "NewIterationCommands should return a command group")
	assert.Equal(t, "iteration", iterationCommands.Name(), "command name should be 'iteration'")
	assert.NotEmpty(t, iterationCommands.Short, "command should have short description")
	assert.NotEmpty(t, iterationCommands.Long, "command should have long description")
}

// TestIterationCommands_AllSubcommands verifies all 10 subcommands are present
func TestIterationCommands_AllSubcommands(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)

	expectedSubcommands := []string{
		"create",
		"list",
		"show",
		"current",
		"start",
		"complete",
		"add-task",
		"remove-task",
		"delete",
		"update",
	}

	commandNames := make(map[string]bool)
	for _, cmd := range iterationCommands.Commands() {
		commandNames[cmd.Name()] = true
	}

	for _, name := range expectedSubcommands {
		assert.True(t, commandNames[name], "command '%s' should exist", name)
	}
}

// TestIterationCreateCommand_Flags verifies create command has required flags
func TestIterationCreateCommand_Flags(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	createCmd := findCommand(iterationCommands, "create")

	assert.NotNil(t, createCmd, "create command should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("name"), "--name flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("goal"), "--goal flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("deliverable"), "--deliverable flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("rank"), "--rank flag should exist")
}

// TestIterationListCommand_Structure verifies list command exists
func TestIterationListCommand_Structure(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	listCmd := findCommand(iterationCommands, "list")

	assert.NotNil(t, listCmd, "list command should exist")
	assert.NotEmpty(t, listCmd.Short, "list command should have description")
}

// TestIterationShowCommand_Arguments verifies show command requires iteration number
func TestIterationShowCommand_Arguments(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	showCmd := findCommand(iterationCommands, "show")

	assert.NotNil(t, showCmd, "show command should exist")
	assert.NotNil(t, showCmd.Args, "show command should have argument validation")
	assert.NotNil(t, showCmd.Flags().Lookup("full"), "--full flag should exist")
}

// TestIterationCurrentCommand_Structure verifies current command exists
func TestIterationCurrentCommand_Structure(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	currentCmd := findCommand(iterationCommands, "current")

	assert.NotNil(t, currentCmd, "current command should exist")
	assert.NotEmpty(t, currentCmd.Short, "current command should have description")
	assert.NotNil(t, currentCmd.Flags().Lookup("full"), "--full flag should exist")
}

// TestIterationStartCommand_Arguments verifies start command requires iteration number
func TestIterationStartCommand_Arguments(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	startCmd := findCommand(iterationCommands, "start")

	assert.NotNil(t, startCmd, "start command should exist")
	assert.NotNil(t, startCmd.Args, "start command should have argument validation")
}

// TestIterationCompleteCommand_Arguments verifies complete command requires iteration number
func TestIterationCompleteCommand_Arguments(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	completeCmd := findCommand(iterationCommands, "complete")

	assert.NotNil(t, completeCmd, "complete command should exist")
	assert.NotNil(t, completeCmd.Args, "complete command should have argument validation")
}

// TestIterationAddTaskCommand_Arguments verifies add-task command requires iteration and tasks
func TestIterationAddTaskCommand_Arguments(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	addTaskCmd := findCommand(iterationCommands, "add-task")

	assert.NotNil(t, addTaskCmd, "add-task command should exist")
	assert.NotNil(t, addTaskCmd.Args, "add-task command should have argument validation")
}

// TestIterationRemoveTaskCommand_Arguments verifies remove-task command requires iteration and tasks
func TestIterationRemoveTaskCommand_Arguments(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	removeTaskCmd := findCommand(iterationCommands, "remove-task")

	assert.NotNil(t, removeTaskCmd, "remove-task command should exist")
	assert.NotNil(t, removeTaskCmd.Args, "remove-task command should have argument validation")
}

// TestIterationDeleteCommand_Arguments verifies delete command requires iteration number
func TestIterationDeleteCommand_Arguments(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	deleteCmd := findCommand(iterationCommands, "delete")

	assert.NotNil(t, deleteCmd, "delete command should exist")
	assert.NotNil(t, deleteCmd.Args, "delete command should have argument validation")
}

// TestIterationUpdateCommand_Flags verifies update command has optional field flags
func TestIterationUpdateCommand_Flags(t *testing.T) {
	iterationCommands := cli.NewIterationCommands(nil, nil, nil)
	updateCmd := findCommand(iterationCommands, "update")

	assert.NotNil(t, updateCmd, "update command should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("name"), "--name flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("goal"), "--goal flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("deliverable"), "--deliverable flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("rank"), "--rank flag should exist")
}

// ============================================================================
// Helper function to find subcommand by name
// ============================================================================
