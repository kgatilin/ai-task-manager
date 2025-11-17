package components_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/components"
)

func TestNewQuitKey(t *testing.T) {
	k := components.NewQuitKey()
	verifyKeyHelp(t, k, "q", "quit")
}

func TestNewBackKey(t *testing.T) {
	k := components.NewBackKey()
	verifyKeyHelp(t, k, "esc", "back")
}

func TestNewHelpKey(t *testing.T) {
	k := components.NewHelpKey()
	verifyKeyHelp(t, k, "?", "help")
}

func TestNewUpKey(t *testing.T) {
	k := components.NewUpKey()
	verifyKeyHelp(t, k, "↑/k", "move up")
}

func TestNewDownKey(t *testing.T) {
	k := components.NewDownKey()
	verifyKeyHelp(t, k, "↓/j", "move down")
}

func TestNewEnterKey(t *testing.T) {
	k := components.NewEnterKey()
	verifyKeyHelp(t, k, "enter", "select")
}

// verifyKeyHelp is a helper to verify key binding help text
func verifyKeyHelp(t *testing.T, k key.Binding, expectedKey, expectedDesc string) {
	help := k.Help()
	if help.Key != expectedKey {
		t.Errorf("expected key %q, got %q", expectedKey, help.Key)
	}
	if help.Desc != expectedDesc {
		t.Errorf("expected desc %q, got %q", expectedDesc, help.Desc)
	}
}

// TestKeyBindingsConsistency verifies that factory functions produce consistent results
func TestKeyBindingsConsistency(t *testing.T) {
	// Call factories multiple times and verify consistent results
	testCases := []struct {
		name        string
		factory     func() key.Binding
		expectedKey string
		expectedDesc string
	}{
		{"quit", components.NewQuitKey, "q", "quit"},
		{"back", components.NewBackKey, "esc", "back"},
		{"help", components.NewHelpKey, "?", "help"},
		{"up", components.NewUpKey, "↑/k", "move up"},
		{"down", components.NewDownKey, "↓/j", "move down"},
		{"enter", components.NewEnterKey, "enter", "select"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call factory multiple times and verify consistency
			for i := 0; i < 3; i++ {
				k := tc.factory()
				help := k.Help()
				if help.Key != tc.expectedKey || help.Desc != tc.expectedDesc {
					t.Errorf("iteration %d: expected (%q, %q), got (%q, %q)",
						i, tc.expectedKey, tc.expectedDesc, help.Key, help.Desc)
				}
			}
		})
	}
}

// TestKeyBindingsEnabled verifies that factory-created bindings are enabled by default
func TestKeyBindingsEnabled(t *testing.T) {
	factories := []func() key.Binding{
		components.NewQuitKey,
		components.NewBackKey,
		components.NewHelpKey,
		components.NewUpKey,
		components.NewDownKey,
		components.NewEnterKey,
	}

	for i, factory := range factories {
		k := factory()
		if !k.Enabled() {
			t.Errorf("factory %d: expected key binding to be enabled", i)
		}
	}
}
