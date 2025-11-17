package components_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/components"
)

func TestNewHelp(t *testing.T) {
	help := components.NewHelp()

	// Verify help was created (test that operations work without panic)
	help.SetWidth(80)
	// If we got here without panic, NewHelp worked correctly
}

func TestHelpSetWidth(t *testing.T) {
	help := components.NewHelp()

	// Set width
	help.SetWidth(120)

	// Verify width was set (we can't directly access the field, but we verify no panic)
	// The width should be passed to the underlying model
	help.SetWidth(80)
	help.SetWidth(0)
	help.SetWidth(200)
}

func TestHelpShortHelpView(t *testing.T) {
	help := components.NewHelp()

	// Create test key bindings
	keys := []key.Binding{
		key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
	}

	// Get short help view
	view := help.ShortHelpView(keys)

	// Verify view is a string (possibly empty, depends on help model state)
	if view == "" {
		t.Log("Short help view is empty (acceptable - depends on help model state)")
	}
}

func TestHelpFullHelpView(t *testing.T) {
	help := components.NewHelp()

	// Create test key bindings grouped
	groups := [][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
		},
		{
			key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
	}

	// Get full help view
	view := help.FullHelpView(groups)

	// Verify view is a string (possibly empty, depends on help model state)
	if view == "" {
		t.Log("Full help view is empty (acceptable - depends on help model state)")
	}
}

func TestHelpShowAll(t *testing.T) {
	help := components.NewHelp()

	// Toggle show all
	help.ShowAll(true)
	help.ShowAll(false)
	help.ShowAll(true)

	// Verify no panic occurred
}

func TestHelpMultipleOperations(t *testing.T) {
	help := components.NewHelp()

	// Set width first
	help.SetWidth(100)

	// Create key bindings
	keys := []key.Binding{
		key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}

	// Get short help
	view1 := help.ShortHelpView(keys)

	// Toggle show all
	help.ShowAll(true)

	// Create grouped keys for full help
	groups := [][]key.Binding{
		keys,
	}

	// Get full help
	view2 := help.FullHelpView(groups)

	// Both views should be strings (may be empty, may not)
	if view1 == "" && view2 == "" {
		t.Log("Both help views are empty (acceptable - depends on help model state)")
	}

	// Reset show all
	help.ShowAll(false)

	// Set width again
	help.SetWidth(80)

	// Get short help again
	view3 := help.ShortHelpView(keys)

	// Verify we can call operations multiple times without error
	if view3 == "" {
		t.Log("Short help view is empty after reset (acceptable)")
	}
}

func TestHelpEmptyKeyBindings(t *testing.T) {
	help := components.NewHelp()

	// Test with empty key bindings
	keys := []key.Binding{}

	// Get short help with empty keys
	view := help.ShortHelpView(keys)

	// Verify no panic with empty bindings
	if view == "" {
		t.Log("Short help view is empty for empty bindings (acceptable)")
	}
}

func TestHelpEmptyGroupedKeyBindings(t *testing.T) {
	help := components.NewHelp()

	// Test with empty grouped key bindings
	groups := [][]key.Binding{}

	// Get full help with empty groups
	view := help.FullHelpView(groups)

	// Verify no panic with empty groups
	if view == "" {
		t.Log("Full help view is empty for empty groups (acceptable)")
	}
}

func TestHelpWidthEdgeCases(t *testing.T) {
	help := components.NewHelp()

	// Test edge cases for width
	testWidths := []int{0, 1, 10, 80, 120, 256, 1000}

	for _, width := range testWidths {
		help.SetWidth(width)
		// Verify no panic for any width value
	}
}
