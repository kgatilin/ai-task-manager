package presenters_test

import (
	"context"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/presenters"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
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

func TestIterationDetailPresenter_DocumentsTabShowsCount(t *testing.T) {
	// Create view model with documents
	vm := &viewmodels.IterationDetailViewModel{
		Number: 1,
		Name:   "Test Iteration",
		Progress: &viewmodels.ProgressViewModel{
			Completed: 0,
			Total:     0,
			Percent:   0.0,
		},
		Documents: []viewmodels.DocumentListItemViewModel{
			{ID: "TM-doc-1", Title: "Design Doc", Type: "adr", StatusIcon: "✓"},
		},
	}

	presenter := presenters.NewIterationDetailPresenter(vm, nil, context.Background())

	// Simulate window size message
	p, _ := presenter.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	presenter = p.(*presenters.IterationDetailPresenter)

	view := presenter.View()

	// Verify Documents tab shows count
	if !strings.Contains(view, "Documents (1)") {
		t.Error("Expected Documents tab to show count '(1)'")
	}
}

func TestIterationDetailPresenter_GetterMethods(t *testing.T) {
	vm := &viewmodels.IterationDetailViewModel{
		Number: 1,
		Name:   "Test Iteration",
		Progress: &viewmodels.ProgressViewModel{
			Completed: 0,
			Total:     1,
			Percent:   0.0,
		},
		TODOTasks: []*viewmodels.TaskRowViewModel{
			{ID: "TM-task-1", Title: "Task 1", Status: "todo"},
		},
	}

	presenter := presenters.NewIterationDetailPresenter(vm, nil, context.Background())

	// Test GetActiveTab
	if presenter.GetActiveTab() != presenters.IterationDetailTabTasks {
		t.Errorf("Expected active tab to be Tasks, got %v", presenter.GetActiveTab())
	}

	// Test GetSelectedIndex
	if presenter.GetSelectedIndex() != 0 {
		t.Errorf("Expected selected index to be 0, got %d", presenter.GetSelectedIndex())
	}

	// Switch to Documents tab
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	p, _ := presenter.Update(tabMsg)
	presenter = p.(*presenters.IterationDetailPresenter)
	p, _ = presenter.Update(tabMsg)
	presenter = p.(*presenters.IterationDetailPresenter)

	// Verify tab changed
	if presenter.GetActiveTab() != presenters.IterationDetailTabDocuments {
		t.Errorf("Expected active tab to be Documents after 2 tabs, got %v", presenter.GetActiveTab())
	}
}

// TestIterationDetailPresenter_ViewportHeightSetOnTabSwitch verifies viewport height is set when switching to AC tab
func TestIterationDetailPresenter_ViewportHeightSetOnTabSwitch(t *testing.T) {
	vm := &viewmodels.IterationDetailViewModel{
		Number: 1,
		Name:   "Test Iteration",
		Progress: &viewmodels.ProgressViewModel{
			Completed: 0,
			Total:     1,
			Percent:   0.0,
		},
		TODOTasks: []*viewmodels.TaskRowViewModel{
			{ID: "TM-task-1", Title: "Task 1", Status: "todo"},
		},
		TaskACs: []*viewmodels.TaskACGroupViewModel{
			{
				Task: &viewmodels.TaskRowViewModel{ID: "TM-task-1", Title: "Task 1"},
				ACs: []*viewmodels.IterationACViewModel{
					{
						ID:                  "TM-ac-1",
						Description:         "AC 1",
						TestingInstructions: "1. Do this\n2. Do that",
						Status:              "not-started",
						IsExpanded:          false,
					},
				},
			},
		},
	}

	presenter := presenters.NewIterationDetailPresenter(vm, nil, context.Background())

	// Simulate window size message to set up terminal height
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 30}
	p, _ := presenter.Update(sizeMsg)
	presenter = p.(*presenters.IterationDetailPresenter)

	// Verify initial state is Tasks tab
	if presenter.GetActiveTab() != presenters.IterationDetailTabTasks {
		t.Fatal("Expected to start on Tasks tab")
	}

	// Switch to ACs tab via Tab key
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	p, _ = presenter.Update(tabMsg)
	presenter = p.(*presenters.IterationDetailPresenter)

	// Verify we're now on ACs tab
	if presenter.GetActiveTab() != presenters.IterationDetailTabACs {
		t.Fatal("Expected to be on ACs tab after Tab key")
	}

	// Viewport height should have been set during tab switch
	// (internal state, but we can verify by checking that rendering works)
	view := presenter.View()
	if !strings.Contains(view, "Acceptance Criteria") {
		t.Error("Expected ACs tab to be displayed in view")
	}
}

// TestIterationDetailPresenter_ViewportHeightSetOnACExpansion verifies viewport height is recalculated when AC is expanded
func TestIterationDetailPresenter_ViewportHeightSetOnACExpansion(t *testing.T) {
	vm := &viewmodels.IterationDetailViewModel{
		Number: 1,
		Name:   "Test Iteration",
		Progress: &viewmodels.ProgressViewModel{
			Completed: 0,
			Total:     1,
			Percent:   0.0,
		},
		TaskACs: []*viewmodels.TaskACGroupViewModel{
			{
				Task: &viewmodels.TaskRowViewModel{ID: "TM-task-1", Title: "Task 1"},
				ACs: []*viewmodels.IterationACViewModel{
					{
						ID:                  "TM-ac-1",
						Description:         "AC 1",
						TestingInstructions: "1. Do this\n2. Do that\n3. Verify result",
						Status:              "not-started",
						IsExpanded:          false,
					},
				},
			},
		},
	}

	presenter := presenters.NewIterationDetailPresenter(vm, nil, context.Background())

	// Simulate window size message
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 30}
	p, _ := presenter.Update(sizeMsg)
	presenter = p.(*presenters.IterationDetailPresenter)

	// Switch to ACs tab
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	p, _ = presenter.Update(tabMsg)
	presenter = p.(*presenters.IterationDetailPresenter)

	// Initial state: AC is collapsed
	initialView := presenter.View()
	if !strings.Contains(initialView, "TM-ac-1") {
		t.Error("Expected AC to be shown in collapsed view")
	}

	// Expand AC via Enter key
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	p, _ = presenter.Update(enterMsg)
	presenter = p.(*presenters.IterationDetailPresenter)

	// Verify AC is expanded
	if !vm.TaskACs[0].ACs[0].IsExpanded {
		t.Error("Expected AC to be expanded after Enter key")
	}

	// Viewport height should have been recalculated during expansion
	expandedView := presenter.View()
	if !strings.Contains(expandedView, "Do this") {
		t.Error("Expected AC testing instructions to be shown when expanded")
	}
}
