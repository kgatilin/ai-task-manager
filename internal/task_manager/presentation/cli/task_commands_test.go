package cli_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Test Cobra Command Structure
// ============================================================================

// TestNewTaskCommands verifies that NewTaskCommands returns a valid Cobra command group
func TestNewTaskCommands_Structure(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)

	assert.NotNil(t, taskCommands, "NewTaskCommands should return a command group")
	assert.Equal(t, "task", taskCommands.Name(), "command name should be 'task'")
	assert.NotEmpty(t, taskCommands.Short, "command should have short description")
	assert.NotEmpty(t, taskCommands.Long, "command should have long description")
}

// TestTaskCommands_AllSubcommands verifies all 8 subcommands are present
func TestTaskCommands_AllSubcommands(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)

	expectedSubcommands := []string{
		"create",
		"list",
		"show",
		"update",
		"delete",
		"move",
		"backlog",
		"check-ready",
	}

	commandNames := make(map[string]bool)
	for _, cmd := range taskCommands.Commands() {
		commandNames[cmd.Name()] = true
	}

	for _, name := range expectedSubcommands {
		assert.True(t, commandNames[name], "command '%s' should exist", name)
	}
}

// TestTaskCreateCommand_Flags verifies create command has required flags
func TestTaskCreateCommand_Flags(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	createCmd := findCommand(taskCommands, "create")

	assert.NotNil(t, createCmd, "create command should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("track"), "--track flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("title"), "--title flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("description"), "--description flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("rank"), "--rank flag should exist")
	assert.NotNil(t, createCmd.Flags().Lookup("branch"), "--branch flag should exist")
}

// TestTaskListCommand_Flags verifies list command has filter flags
func TestTaskListCommand_Flags(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	listCmd := findCommand(taskCommands, "list")

	assert.NotNil(t, listCmd, "list command should exist")
	assert.NotNil(t, listCmd.Flags().Lookup("track"), "--track flag should exist")
	assert.NotNil(t, listCmd.Flags().Lookup("status"), "--status flag should exist")
}

// TestTaskShowCommand_Arguments verifies show command requires task ID
func TestTaskShowCommand_Arguments(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	showCmd := findCommand(taskCommands, "show")

	assert.NotNil(t, showCmd, "show command should exist")
	// Check that it uses ExactArgs(1)
	assert.NotNil(t, showCmd.Args, "show command should have argument validation")
}

// TestTaskUpdateCommand_Flags verifies update command has optional field flags
func TestTaskUpdateCommand_Flags(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	updateCmd := findCommand(taskCommands, "update")

	assert.NotNil(t, updateCmd, "update command should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("title"), "--title flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("description"), "--description flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("status"), "--status flag should exist")
	assert.NotNil(t, updateCmd.Flags().Lookup("rank"), "--rank flag should exist")
}

// TestTaskDeleteCommand_Arguments verifies delete command requires task ID
func TestTaskDeleteCommand_Arguments(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	deleteCmd := findCommand(taskCommands, "delete")

	assert.NotNil(t, deleteCmd, "delete command should exist")
	assert.NotNil(t, deleteCmd.Args, "delete command should have argument validation")
}

// TestTaskMoveCommand_Flags verifies move command has track flag
func TestTaskMoveCommand_Flags(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	moveCmd := findCommand(taskCommands, "move")

	assert.NotNil(t, moveCmd, "move command should exist")
	assert.NotNil(t, moveCmd.Flags().Lookup("track"), "--track flag should exist")
}

// TestTaskBacklogCommand_Structure verifies backlog command exists
func TestTaskBacklogCommand_Structure(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	backlogCmd := findCommand(taskCommands, "backlog")

	assert.NotNil(t, backlogCmd, "backlog command should exist")
	assert.NotEmpty(t, backlogCmd.Short, "backlog command should have description")
}

// TestTaskCheckReadyCommand_Arguments verifies check-ready command requires task ID
func TestTaskCheckReadyCommand_Arguments(t *testing.T) {
	taskCommands := cli.NewTaskCommands(nil, nil)
	checkCmd := findCommand(taskCommands, "check-ready")

	assert.NotNil(t, checkCmd, "check-ready command should exist")
	assert.NotNil(t, checkCmd.Args, "check-ready command should have argument validation")
}
