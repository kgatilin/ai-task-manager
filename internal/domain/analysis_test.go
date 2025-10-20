package domain_test

import (
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

func TestNewSessionAnalysis(t *testing.T) {
	sessionID := "test-session-123"
	analysisResult := "This is the analysis result from LLM"
	modelUsed := "claude-sonnet-4"
	promptUsed := "analyze this session"

	analysis := domain.NewSessionAnalysis(sessionID, analysisResult, modelUsed, promptUsed)

	// Verify required fields
	if analysis.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
	if analysis.SessionID != sessionID {
		t.Errorf("Expected SessionID = %q, got %q", sessionID, analysis.SessionID)
	}
	if analysis.AnalysisResult != analysisResult {
		t.Errorf("Expected AnalysisResult = %q, got %q", analysisResult, analysis.AnalysisResult)
	}
	if analysis.ModelUsed != modelUsed {
		t.Errorf("Expected ModelUsed = %q, got %q", modelUsed, analysis.ModelUsed)
	}
	if analysis.PromptUsed != promptUsed {
		t.Errorf("Expected PromptUsed = %q, got %q", promptUsed, analysis.PromptUsed)
	}

	// Verify defaults for backward compatibility
	if analysis.AnalysisType != "tool_analysis" {
		t.Errorf("Expected AnalysisType = %q, got %q", "tool_analysis", analysis.AnalysisType)
	}
	if analysis.PromptName != "analysis" {
		t.Errorf("Expected PromptName = %q, got %q", "analysis", analysis.PromptName)
	}

	// Verify timestamp is recent
	if time.Since(analysis.AnalyzedAt) > time.Second {
		t.Errorf("Expected recent timestamp, got %v", analysis.AnalyzedAt)
	}

	// Verify PatternsSummary is empty by default
	if analysis.PatternsSummary != "" {
		t.Errorf("Expected empty PatternsSummary, got %q", analysis.PatternsSummary)
	}
}

func TestNewSessionAnalysisWithType(t *testing.T) {
	tests := []struct {
		name           string
		sessionID      string
		analysisResult string
		modelUsed      string
		promptUsed     string
		analysisType   string
		promptName     string
	}{
		{
			name:           "creates session summary analysis",
			sessionID:      "session-1",
			analysisResult: "Session summary text",
			modelUsed:      "claude-opus-4",
			promptUsed:     "summarize this session",
			analysisType:   "session_summary",
			promptName:     "session_summary",
		},
		{
			name:           "creates tool analysis",
			sessionID:      "session-2",
			analysisResult: "Tool analysis text",
			modelUsed:      "claude-sonnet-4",
			promptUsed:     "analyze tool usage",
			analysisType:   "tool_analysis",
			promptName:     "tool_analysis",
		},
		{
			name:           "creates custom analysis type",
			sessionID:      "session-3",
			analysisResult: "Custom analysis",
			modelUsed:      "claude-haiku-4",
			promptUsed:     "custom prompt",
			analysisType:   "custom_type",
			promptName:     "custom_prompt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := domain.NewSessionAnalysisWithType(
				tt.sessionID,
				tt.analysisResult,
				tt.modelUsed,
				tt.promptUsed,
				tt.analysisType,
				tt.promptName,
			)

			if analysis.ID == "" {
				t.Error("Expected ID to be generated, got empty string")
			}
			if analysis.SessionID != tt.sessionID {
				t.Errorf("Expected SessionID = %q, got %q", tt.sessionID, analysis.SessionID)
			}
			if analysis.AnalysisResult != tt.analysisResult {
				t.Errorf("Expected AnalysisResult = %q, got %q", tt.analysisResult, analysis.AnalysisResult)
			}
			if analysis.ModelUsed != tt.modelUsed {
				t.Errorf("Expected ModelUsed = %q, got %q", tt.modelUsed, analysis.ModelUsed)
			}
			if analysis.PromptUsed != tt.promptUsed {
				t.Errorf("Expected PromptUsed = %q, got %q", tt.promptUsed, analysis.PromptUsed)
			}
			if analysis.AnalysisType != tt.analysisType {
				t.Errorf("Expected AnalysisType = %q, got %q", tt.analysisType, analysis.AnalysisType)
			}
			if analysis.PromptName != tt.promptName {
				t.Errorf("Expected PromptName = %q, got %q", tt.promptName, analysis.PromptName)
			}

			// Verify timestamp is recent
			if time.Since(analysis.AnalyzedAt) > time.Second {
				t.Errorf("Expected recent timestamp, got %v", analysis.AnalyzedAt)
			}
		})
	}
}

func TestSessionAnalysis_AllFields(t *testing.T) {
	// Test that all fields can be set and retrieved
	analysis := &domain.SessionAnalysis{
		ID:              "test-id",
		SessionID:       "test-session",
		AnalyzedAt:      time.Now(),
		AnalysisResult:  "full analysis text",
		ModelUsed:       "claude-sonnet-4",
		PromptUsed:      "test prompt",
		PatternsSummary: "summary of patterns",
		AnalysisType:    "tool_analysis",
		PromptName:      "custom_prompt",
	}

	// Verify all fields are accessible
	if analysis.ID != "test-id" {
		t.Errorf("ID mismatch")
	}
	if analysis.SessionID != "test-session" {
		t.Errorf("SessionID mismatch")
	}
	if analysis.AnalysisResult != "full analysis text" {
		t.Errorf("AnalysisResult mismatch")
	}
	if analysis.ModelUsed != "claude-sonnet-4" {
		t.Errorf("ModelUsed mismatch")
	}
	if analysis.PromptUsed != "test prompt" {
		t.Errorf("PromptUsed mismatch")
	}
	if analysis.PatternsSummary != "summary of patterns" {
		t.Errorf("PatternsSummary mismatch")
	}
	if analysis.AnalysisType != "tool_analysis" {
		t.Errorf("AnalysisType mismatch")
	}
	if analysis.PromptName != "custom_prompt" {
		t.Errorf("PromptName mismatch")
	}
}

func TestToolSuggestion_Fields(t *testing.T) {
	// Test that ToolSuggestion struct can be created and fields accessed
	suggestion := domain.ToolSuggestion{
		Name:        "CodeAnalyzer",
		Description: "Analyzes code patterns",
		Rationale:   "Would speed up analysis by 50%",
		Examples:    []string{"example 1", "example 2"},
	}

	if suggestion.Name != "CodeAnalyzer" {
		t.Errorf("Name mismatch")
	}
	if suggestion.Description != "Analyzes code patterns" {
		t.Errorf("Description mismatch")
	}
	if suggestion.Rationale != "Would speed up analysis by 50%" {
		t.Errorf("Rationale mismatch")
	}
	if len(suggestion.Examples) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(suggestion.Examples))
	}
	if suggestion.Examples[0] != "example 1" {
		t.Errorf("Examples[0] mismatch")
	}
}

func TestNewSessionAnalysis_UniqueIDs(t *testing.T) {
	// Verify that multiple calls generate unique IDs
	analysis1 := domain.NewSessionAnalysis("session-1", "result1", "model1", "prompt1")
	analysis2 := domain.NewSessionAnalysis("session-1", "result2", "model2", "prompt2")

	if analysis1.ID == analysis2.ID {
		t.Error("Expected unique IDs for different analyses, got same ID")
	}
}

func TestNewSessionAnalysisWithType_UniqueIDs(t *testing.T) {
	// Verify that multiple calls generate unique IDs
	analysis1 := domain.NewSessionAnalysisWithType("session-1", "result1", "model1", "prompt1", "type1", "name1")
	analysis2 := domain.NewSessionAnalysisWithType("session-1", "result2", "model2", "prompt2", "type2", "name2")

	if analysis1.ID == analysis2.ID {
		t.Error("Expected unique IDs for different analyses, got same ID")
	}
}

func TestSessionAnalysis_EmptyFields(t *testing.T) {
	// Test that analysis can be created with empty optional fields
	analysis := domain.NewSessionAnalysisWithType("session-1", "", "", "", "", "")

	if analysis.ID == "" {
		t.Error("ID should be generated even with empty fields")
	}
	if analysis.SessionID != "session-1" {
		t.Error("SessionID should be preserved")
	}
	if analysis.AnalysisResult != "" {
		t.Error("Empty AnalysisResult should be preserved")
	}
}
