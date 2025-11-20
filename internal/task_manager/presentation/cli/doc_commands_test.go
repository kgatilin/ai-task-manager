package cli_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Test Cobra Command Structure
// ============================================================================

// TestNewDocCommands verifies that NewDocCommands returns a valid Cobra command group
func TestNewDocCommands_Structure(t *testing.T) {
	docCommands := cli.NewDocCommands(nil)

	assert.NotNil(t, docCommands, "NewDocCommands should return a command group")
	assert.Equal(t, "doc", docCommands.Name(), "command name should be 'doc'")
	assert.NotEmpty(t, docCommands.Short, "command should have short description")
	assert.NotEmpty(t, docCommands.Long, "command should have long description")
}

// TestDocCommands_AllSubcommands verifies all 7 subcommands are present
func TestDocCommands_AllSubcommands(t *testing.T) {
	docCommands := cli.NewDocCommands(nil)

	expectedSubcommands := []string{
		"create",
		"list",
		"show",
		"update",
		"attach",
		"detach",
		"delete",
	}

	commandNames := make(map[string]bool)
	for _, cmd := range docCommands.Commands() {
		commandNames[cmd.Name()] = true
	}

	for _, name := range expectedSubcommands {
		assert.True(t, commandNames[name], "command '%s' should exist", name)
	}
}

// ============================================================================
// Test Flags and Command Structure
// ============================================================================

// TestDocCommands_HasRequiredSubcommands verifies all required commands exist and are properly configured
func TestDocCommands_HasRequiredSubcommands(t *testing.T) {
	docCommands := cli.NewDocCommands(nil)

	subcommands := make(map[string]*cobra.Command)
	for _, cmd := range docCommands.Commands() {
		subcommands[cmd.Name()] = cmd
	}

	// Verify all subcommands exist
	assert.NotNil(t, subcommands["create"], "create command should exist")
	assert.NotNil(t, subcommands["list"], "list command should exist")
	assert.NotNil(t, subcommands["show"], "show command should exist")
	assert.NotNil(t, subcommands["update"], "update command should exist")
	assert.NotNil(t, subcommands["attach"], "attach command should exist")
	assert.NotNil(t, subcommands["detach"], "detach command should exist")
	assert.NotNil(t, subcommands["delete"], "delete command should exist")

	// Verify create has required flags
	create := subcommands["create"]
	assert.NotNil(t, create.Flag("title"), "create should have --title flag")
	assert.NotNil(t, create.Flag("type"), "create should have --type flag")

	// Verify list has filter flags
	list := subcommands["list"]
	assert.NotNil(t, list.Flag("track"), "list should have --track flag")
	assert.NotNil(t, list.Flag("iteration"), "list should have --iteration flag")
	assert.NotNil(t, list.Flag("type"), "list should have --type flag")

	// Verify show requires positional argument
	show := subcommands["show"]
	assert.NotNil(t, show.Args, "show should have args validation")

	// Verify update has update flags
	update := subcommands["update"]
	assert.NotNil(t, update.Flag("content"), "update should have --content flag")
	assert.NotNil(t, update.Flag("from-file"), "update should have --from-file flag")
	assert.NotNil(t, update.Flag("status"), "update should have --status flag")
	assert.NotNil(t, update.Flag("detach"), "update should have --detach flag")

	// Verify attach requires positional argument and has flags
	attach := subcommands["attach"]
	assert.NotNil(t, attach.Args, "attach should have args validation")
	assert.NotNil(t, attach.Flag("track"), "attach should have --track flag")
	assert.NotNil(t, attach.Flag("iteration"), "attach should have --iteration flag")

	// Verify detach requires positional argument
	detach := subcommands["detach"]
	assert.NotNil(t, detach.Args, "detach should have args validation")

	// Verify delete requires positional argument and has force flag
	delete := subcommands["delete"]
	assert.NotNil(t, delete.Args, "delete should have args validation")
	assert.NotNil(t, delete.Flag("force"), "delete should have --force flag")
}

// ============================================================================
// Test Aliases
// ============================================================================

// TestDocCommands_HasAliases verifies the doc command group has aliases
func TestDocCommands_HasAliases(t *testing.T) {
	docCommands := cli.NewDocCommands(nil)

	assert.NotEmpty(t, docCommands.Aliases, "doc command should have aliases")
	assert.Contains(t, docCommands.Aliases, "d", "doc command should have 'd' alias")
}
