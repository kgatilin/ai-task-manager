package tui_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/internal/app/tui"
	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

func TestNewAnalysisViewerModel(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:     "test-session",
		ViewType:   "session",
		PromptUsed: "test_prompt",
		ModelUsed:  "test-model",
		Result:     "# Test Analysis\n\nThis is a test analysis result.",
		Timestamp:  time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Model should be created successfully
	// Note: NewAnalysisViewerModel returns a value type, not a pointer,
	// so we can't check for nil. The function always returns a valid model.
	_ = model
}

func TestAnalysisViewerModel_Init(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:      "test-session",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		Result: "Test result",
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)
	cmd := model.Init()

	if cmd != nil {
		t.Error("Init() should return nil command")
	}
}

func TestAnalysisViewerModel_UpdateWindowSize(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:      "test-session",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		Result: "Test result",
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("WindowSizeMsg should return nil command")
	}

	_, ok := updatedModel.(tui.AnalysisViewerModel)
	if !ok {
		t.Error("Update should return AnalysisViewerModel")
	}
}

func TestAnalysisViewerModel_UpdateEsc(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:      "test-session-abc",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		Result: "Test result",
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Initialize viewport with window size first
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Send esc key
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Esc key should return a command")
	}

	// Execute command to verify it's BackToDetailMsg
	if cmd != nil {
		result := cmd()
		if _, ok := result.(tui.BackToDetailMsg); !ok {
			t.Error("Expected BackToDetailMsg from esc key")
		}
	}
}

func TestAnalysisViewerModel_UpdateScrolling(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:      "test-session",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		Result: "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Initialize viewport
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 10})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Test down arrow scrolling
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, cmd := model.Update(msg)

	if cmd == nil {
		// Scrolling might not return a command, that's ok
	}

	_, ok := updatedModel.(tui.AnalysisViewerModel)
	if !ok {
		t.Error("Update should return AnalysisViewerModel after scrolling")
	}
}

func TestAnalysisViewerModel_ViewNotReady(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:      "test-session",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		Result: "Test result",
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// View before initialization should return initializing message
	view := model.View()

	if view == "" {
		t.Error("View() should return non-empty string")
	}
}

func TestAnalysisViewerModel_ViewReady(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:      "test-session-xyz",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		ModelUsed:      "test-model",
		Result: "# Analysis\n\nTest content",
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Initialize with window size
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	model = updatedModel.(tui.AnalysisViewerModel)

	// View after initialization should show content
	view := model.View()

	if view == "" {
		t.Error("View() should return non-empty string after initialization")
	}
}

func TestAnalysisViewerModel_UpdateScrollKeys(t *testing.T) {
	analysis := &domain.Analysis{
		ViewID:      "test-session",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		ModelUsed:      "test-model",
		Result: "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8",
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Initialize viewport
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 10})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Test up arrow
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Model is a value type, always valid after Update

	// Test page up
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Model is a value type, always valid after Update

	// Test page down
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Model is a value type, always valid after Update
	_ = model
}

func TestAnalysisViewerModel_FooterScrollPercent(t *testing.T) {
	longAnalysis := ""
	for i := 0; i < 500; i++ {
		longAnalysis += "Analysis line " + string(rune(i)) + "\n"
	}

	analysis := &domain.Analysis{
		ViewID:      "test-session",
		Result: longAnalysis,
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Initialize
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 15})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Scroll through content
	for i := 0; i < 20; i++ {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = updatedModel.(tui.AnalysisViewerModel)

		view := model.View()
		if view == "" {
			t.Error("View should show scroll percentage")
		}
	}
}

func TestAnalysisViewerModel_RenderContent_Branches(t *testing.T) {
	// Test with very long content to trigger different rendering
	longResult := ""
	for i := 0; i < 1000; i++ {
		longResult += "# Heading\n\nParagraph text here.\n\n"
	}

	analysis := &domain.Analysis{
		ViewID:      "test-session-long-analysis",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		ModelUsed:      "test-model",
		Result: longResult,
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Initialize with different widths to test word wrap
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 50})
	model = updatedModel.(tui.AnalysisViewerModel)

	view1 := model.View()
	if view1 == "" {
		t.Error("View should return non-empty string with long analysis")
	}

	// Test with very small width
	updatedModel2, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 50})
	model = updatedModel2.(tui.AnalysisViewerModel)

	view2 := model.View()
	if view2 == "" {
		t.Error("View should return non-empty string with small width")
	}
}

func TestAnalysisViewerModel_RenderError(t *testing.T) {
	// Test with empty analysis result to check error handling
	analysis := &domain.Analysis{
		ViewID:      "test-session",
		ViewType:   "session",
		PromptUsed:     "test_prompt",
		ModelUsed:      "test-model",
		Result: "", // Empty result
		Timestamp:     time.Now(),
	}

	model := tui.NewAnalysisViewerModel(analysis)

	// Initialize with window size
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	model = updatedModel.(tui.AnalysisViewerModel)

	// Call View
	view := model.View()

	if view == "" {
		t.Error("View should return non-empty string even with empty analysis")
	}
}
