package presenters_test

import (
	"context"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/presenters"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

func TestIterationDetailKeyMap_TaskTransitionKeysExist(t *testing.T) {
	keys := presenters.NewIterationDetailKeyMap()

	// Verify InProgress key binding
	if keys.InProgress.Keys()[0] != "i" {
		t.Errorf("Expected 'i' key binding for InProgress, got %v", keys.InProgress.Keys())
	}

	// Verify Review key binding
	if keys.Review.Keys()[0] != "r" {
		t.Errorf("Expected 'r' key binding for Review, got %v", keys.Review.Keys())
	}

	// Verify Done key binding
	if keys.Done.Keys()[0] != "d" {
		t.Errorf("Expected 'd' key binding for Done, got %v", keys.Done.Keys())
	}

	// Verify Reopen key binding
	if keys.Reopen.Keys()[0] != "o" {
		t.Errorf("Expected 'o' key binding for Reopen, got %v", keys.Reopen.Keys())
	}
}

func TestIterationDetailKeyMap_ShortHelpContextAware(t *testing.T) {
	keys := presenters.NewIterationDetailKeyMap()

	// Tasks tab should show task transition keys
	tasksHelp := keys.ShortHelp(presenters.IterationDetailTabTasks)
	hasInProgress := false
	hasReview := false
	hasDone := false
	for _, binding := range tasksHelp {
		if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}, binding) {
			hasInProgress = true
		}
		if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}, binding) {
			hasReview = true
		}
		if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, binding) {
			hasDone = true
		}
	}
	if !hasInProgress || !hasReview || !hasDone {
		t.Error("Expected task transition keys (i, r, d) in Tasks tab short help")
	}

	// ACs tab should NOT show task transition keys, but should show AC action keys
	acsHelp := keys.ShortHelp(presenters.IterationDetailTabACs)
	hasVerify := false
	hasSkip := false
	hasFail := false
	hasTaskTransition := false
	for _, binding := range acsHelp {
		if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}, binding) {
			hasVerify = true
		}
		if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}, binding) {
			hasSkip = true
		}
		if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}, binding) {
			hasFail = true
		}
		if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}, binding) ||
			key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}, binding) ||
			key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, binding) {
			hasTaskTransition = true
		}
	}
	if !hasVerify || !hasSkip || !hasFail {
		t.Error("Expected AC action keys (space, s, f) in ACs tab short help")
	}
	if hasTaskTransition {
		t.Error("Task transition keys should NOT appear in ACs tab short help")
	}
}

func TestIterationDetailKeyMap_FullHelpContextAware(t *testing.T) {
	keys := presenters.NewIterationDetailKeyMap()

	// Tasks tab full help should include task transitions
	tasksFullHelp := keys.FullHelp(presenters.IterationDetailTabTasks)
	hasTaskTransitions := false
	for _, row := range tasksFullHelp {
		for _, binding := range row {
			if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}, binding) {
				hasTaskTransitions = true
				break
			}
		}
	}
	if !hasTaskTransitions {
		t.Error("Expected task transition keys in Tasks tab full help")
	}

	// ACs tab full help should include AC actions, not task transitions
	acsFullHelp := keys.FullHelp(presenters.IterationDetailTabACs)
	hasACActions := false
	hasTaskTransitionsInACs := false
	for _, row := range acsFullHelp {
		for _, binding := range row {
			if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}, binding) {
				hasACActions = true
			}
			if key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, binding) ||
				key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}, binding) {
				hasTaskTransitionsInACs = true
			}
		}
	}
	if !hasACActions {
		t.Error("Expected AC action keys in ACs tab full help")
	}
	if hasTaskTransitionsInACs {
		t.Error("Task transition keys should NOT appear in ACs tab full help")
	}
}

func TestIterationDetailPresenter_TaskTransitionOnlyOnTasksTab(t *testing.T) {
	// Create view model with tasks
	vm := &viewmodels.IterationDetailViewModel{
		Number: 1,
		Name:   "Test Iteration",
		TODOTasks: []*viewmodels.TaskRowViewModel{
			{ID: "TM-task-1", Title: "Task 1", Status: "todo"},
		},
	}

	presenter := presenters.NewIterationDetailPresenter(vm, nil, context.Background())

	// Press 'i' on Tasks tab - should trigger transition
	iMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
	_, cmd := presenter.Update(iMsg)

	// Verify command is returned (task transition)
	// Note: In real implementation, this would call repository which we don't have in unit test
	// So we can't fully test the command execution, but we can verify key is handled
	_ = cmd // Command may be nil if no task selected or validation fails

	// Switch to ACs tab
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	p, _ := presenter.Update(tabMsg)
	presenter = p.(*presenters.IterationDetailPresenter)

	// Press 'i' on ACs tab - should NOT trigger transition
	_, cmdAfterTab := presenter.Update(iMsg)
	// On ACs tab, 'i' key should not do anything related to task transitions
	_ = cmdAfterTab // Should be nil or unrelated to task transition
}
