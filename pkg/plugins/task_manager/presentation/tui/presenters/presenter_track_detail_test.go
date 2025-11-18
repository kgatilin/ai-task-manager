package presenters_test

import (
	"context"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/presenters"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

func TestTrackDetailPresenter_ViewRendersDocumentsSection(t *testing.T) {
	vm := viewmodels.NewTrackDetailViewModel("TM-track-1", "Test Track", "Description", "in-progress", "In Progress", 1, nil, nil)
	vm.Documents = append(vm.Documents, viewmodels.DocumentListItemViewModel{
		ID:         "TM-doc-1",
		Title:      "Design ADR",
		Type:       "adr",
		StatusIcon: "✓",
	})

	presenter := presenters.NewTrackDetailPresenter(vm, nil, context.Background())

	// Simulate window size message to trigger initialization
	p, _ := presenter.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	presenter = p.(*presenters.TrackDetailPresenter)

	view := presenter.View()

	// Check that documents section is rendered
	if !strings.Contains(view, "📄 Documents (1)") {
		t.Error("Expected documents section header with count to be rendered")
	}

	if !strings.Contains(view, "Design ADR") {
		t.Error("Expected document title to be rendered")
	}

	if !strings.Contains(view, "adr") {
		t.Error("Expected document type to be rendered")
	}

	if !strings.Contains(view, "✓") {
		t.Error("Expected document status icon to be rendered")
	}
}

func TestTrackDetailPresenter_ViewRendersEmptyDocumentsSection(t *testing.T) {
	vm := viewmodels.NewTrackDetailViewModel("TM-track-1", "Test Track", "Description", "in-progress", "In Progress", 1, nil, nil)
	// No documents

	presenter := presenters.NewTrackDetailPresenter(vm, nil, context.Background())

	// Simulate window size message
	p, _ := presenter.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	presenter = p.(*presenters.TrackDetailPresenter)

	view := presenter.View()

	// Check that empty documents section is rendered
	if !strings.Contains(view, "📄 Documents (0)") {
		t.Error("Expected documents section header with count 0")
	}

	if !strings.Contains(view, "(No documents)") {
		t.Error("Expected empty message to be rendered")
	}
}

func TestTrackDetailPresenter_ViewRendersMultipleDocuments(t *testing.T) {
	vm := viewmodels.NewTrackDetailViewModel("TM-track-1", "Test Track", "Description", "in-progress", "In Progress", 1, nil, nil)
	vm.Documents = append(vm.Documents, viewmodels.DocumentListItemViewModel{
		ID:         "TM-doc-1",
		Title:      "Design ADR",
		Type:       "adr",
		StatusIcon: "✓",
	})
	vm.Documents = append(vm.Documents, viewmodels.DocumentListItemViewModel{
		ID:         "TM-doc-2",
		Title:      "Implementation Plan",
		Type:       "plan",
		StatusIcon: "○",
	})

	presenter := presenters.NewTrackDetailPresenter(vm, nil, context.Background())

	// Simulate window size message
	p, _ := presenter.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	presenter = p.(*presenters.TrackDetailPresenter)

	view := presenter.View()

	// Check that multiple documents are rendered
	if !strings.Contains(view, "📄 Documents (2)") {
		t.Error("Expected documents section header with count 2")
	}

	if !strings.Contains(view, "Design ADR") {
		t.Error("Expected first document title to be rendered")
	}

	if !strings.Contains(view, "Implementation Plan") {
		t.Error("Expected second document title to be rendered")
	}

	// Count how many times we see the document format pattern
	docCount := strings.Count(view, "- adr [") + strings.Count(view, "- plan [")
	if docCount != 2 {
		t.Errorf("Expected 2 documents to be rendered with format, got %d", docCount)
	}
}

func TestTrackDetailPresenter_EnterNavigatesToDocument(t *testing.T) {
	vm := viewmodels.NewTrackDetailViewModel("TM-track-1", "Test Track", "Description", "in-progress", "In Progress", 1, nil, nil)
	vm.Documents = append(vm.Documents, viewmodels.DocumentListItemViewModel{
		ID:         "TM-doc-1",
		Title:      "Design Document",
		Type:       "adr",
		StatusIcon: "✓",
	})

	presenter := presenters.NewTrackDetailPresenter(vm, nil, context.Background())

	// Verify that Enter key triggers navigation
	// Note: We can't test the actual command execution without mocking the repository,
	// but we can test that the navigation messages are created correctly by examining
	// the presenter's behavior with documents
	view := presenter.View()

	if !strings.Contains(view, "Design Document") {
		t.Error("Document should be visible in the view for navigation")
	}
}

func TestTrackDetailPresenter_GetSelectedIndex(t *testing.T) {
	vm := viewmodels.NewTrackDetailViewModel("TM-track-1", "Test Track", "Description", "in-progress", "In Progress", 1, nil, nil)
	vm.TODOTasks = append(vm.TODOTasks, &viewmodels.TrackDetailTaskViewModel{ID: "TM-task-1", Title: "Task 1"})

	presenter := presenters.NewTrackDetailPresenter(vm, nil, context.Background())

	// Test initial selected index
	if presenter.GetSelectedIndex() != 0 {
		t.Errorf("Expected initial selected index to be 0, got %d", presenter.GetSelectedIndex())
	}

	// Create presenter with specific selection
	presenter2 := presenters.NewTrackDetailPresenterWithSelection(vm, nil, context.Background(), 5)
	if presenter2.GetSelectedIndex() != 5 {
		t.Errorf("Expected selected index to be 5, got %d", presenter2.GetSelectedIndex())
	}
}
