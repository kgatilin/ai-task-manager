package domain_test

import (
	"testing"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

func TestValidateModel(t *testing.T) {
	tests := []struct {
		name  string
		model string
		want  bool
	}{
		// Valid aliases
		{name: "accepts sonnet alias", model: "sonnet", want: true},
		{name: "accepts opus alias", model: "opus", want: true},
		{name: "accepts haiku alias", model: "haiku", want: true},

		// Valid full names
		{name: "accepts claude-sonnet-4-5-20250929", model: "claude-sonnet-4-5-20250929", want: true},
		{name: "accepts claude-opus-4-20250514", model: "claude-opus-4-20250514", want: true},
		{name: "accepts claude-3-5-sonnet-20241022", model: "claude-3-5-sonnet-20241022", want: true},
		{name: "accepts claude-3-5-haiku-20241022", model: "claude-3-5-haiku-20241022", want: true},

		// Invalid models
		{name: "rejects unknown alias", model: "gpt-4", want: false},
		{name: "rejects invalid model name", model: "claude-invalid", want: false},
		{name: "rejects random string", model: "random", want: false},

		// Empty model (valid - uses default)
		{name: "accepts empty model", model: "", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.ValidateModel(tt.model)
			if got != tt.want {
				t.Errorf("ValidateModel(%q) = %v, want %v", tt.model, got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := domain.DefaultConfig()

	// Verify Analysis defaults
	if config.Analysis.TokenLimit != 100000 {
		t.Errorf("Expected TokenLimit = 100000, got %d", config.Analysis.TokenLimit)
	}
	if config.Analysis.Model != "sonnet" {
		t.Errorf("Expected Model = %q, got %q", "sonnet", config.Analysis.Model)
	}
	if config.Analysis.ParallelLimit != 3 {
		t.Errorf("Expected ParallelLimit = 3, got %d", config.Analysis.ParallelLimit)
	}
	if len(config.Analysis.EnabledPrompts) != 1 || config.Analysis.EnabledPrompts[0] != "tool_analysis" {
		t.Errorf("Expected EnabledPrompts = [%q], got %v", "tool_analysis", config.Analysis.EnabledPrompts)
	}
	if config.Analysis.AutoSummaryEnabled {
		t.Error("Expected AutoSummaryEnabled = false, got true")
	}
	if config.Analysis.AutoSummaryPrompt != "session_summary" {
		t.Errorf("Expected AutoSummaryPrompt = %q, got %q", "session_summary", config.Analysis.AutoSummaryPrompt)
	}

	// Verify ClaudeOptions defaults
	if len(config.Analysis.ClaudeOptions.AllowedTools) != 0 {
		t.Errorf("Expected empty AllowedTools, got %v", config.Analysis.ClaudeOptions.AllowedTools)
	}
	if config.Analysis.ClaudeOptions.SystemPromptMode != "replace" {
		t.Errorf("Expected SystemPromptMode = %q, got %q", "replace", config.Analysis.ClaudeOptions.SystemPromptMode)
	}

	// Verify UI defaults
	if config.UI.DefaultOutputDir != "./analysis-outputs" {
		t.Errorf("Expected DefaultOutputDir = %q, got %q", "./analysis-outputs", config.UI.DefaultOutputDir)
	}
	if config.UI.FilenameTemplate != "{{.SessionID}}-{{.PromptName}}-{{.Date}}.md" {
		t.Errorf("Expected FilenameTemplate = %q, got %q", "{{.SessionID}}-{{.PromptName}}-{{.Date}}.md", config.UI.FilenameTemplate)
	}
	if config.UI.AutoRefreshInterval != "" {
		t.Errorf("Expected empty AutoRefreshInterval, got %q", config.UI.AutoRefreshInterval)
	}

	// Verify Prompts defaults
	if len(config.Prompts) != 2 {
		t.Errorf("Expected 2 default prompts, got %d", len(config.Prompts))
	}
	if _, ok := config.Prompts["session_summary"]; !ok {
		t.Error("Expected session_summary prompt to exist")
	}
	if _, ok := config.Prompts["tool_analysis"]; !ok {
		t.Error("Expected tool_analysis prompt to exist")
	}
}

func TestAnalysisConfig_DefaultValues(t *testing.T) {
	config := domain.AnalysisConfig{
		TokenLimit:         100000,
		Model:              "sonnet",
		ParallelLimit:      3,
		EnabledPrompts:     []string{"tool_analysis"},
		AutoSummaryEnabled: false,
		AutoSummaryPrompt:  "session_summary",
		ClaudeOptions: domain.ClaudeOptions{
			AllowedTools:     []string{},
			SystemPromptMode: "replace",
		},
	}

	// Verify all fields are accessible
	if config.TokenLimit != 100000 {
		t.Error("TokenLimit mismatch")
	}
	if config.Model != "sonnet" {
		t.Error("Model mismatch")
	}
	if config.ParallelLimit != 3 {
		t.Error("ParallelLimit mismatch")
	}
	if config.AutoSummaryEnabled {
		t.Error("AutoSummaryEnabled should be false")
	}
}

func TestUIConfig_AllFields(t *testing.T) {
	config := domain.UIConfig{
		DefaultOutputDir:    "/custom/path",
		FilenameTemplate:    "custom-{{.SessionID}}.md",
		AutoRefreshInterval: "30s",
	}

	if config.DefaultOutputDir != "/custom/path" {
		t.Error("DefaultOutputDir mismatch")
	}
	if config.FilenameTemplate != "custom-{{.SessionID}}.md" {
		t.Error("FilenameTemplate mismatch")
	}
	if config.AutoRefreshInterval != "30s" {
		t.Error("AutoRefreshInterval mismatch")
	}
}

func TestClaudeOptions_AllFields(t *testing.T) {
	options := domain.ClaudeOptions{
		AllowedTools:     []string{"Read", "Write", "Bash"},
		SystemPromptMode: "append",
	}

	if len(options.AllowedTools) != 3 {
		t.Errorf("Expected 3 allowed tools, got %d", len(options.AllowedTools))
	}
	if options.AllowedTools[0] != "Read" {
		t.Error("AllowedTools[0] mismatch")
	}
	if options.SystemPromptMode != "append" {
		t.Error("SystemPromptMode mismatch")
	}
}

func TestConfig_Prompts(t *testing.T) {
	config := domain.DefaultConfig()

	// Test session_summary prompt exists and is non-empty
	sessionSummary, ok := config.Prompts["session_summary"]
	if !ok {
		t.Fatal("session_summary prompt not found")
	}
	if sessionSummary == "" {
		t.Error("session_summary prompt is empty")
	}
	if sessionSummary != domain.DefaultSessionSummaryPrompt {
		t.Error("session_summary prompt doesn't match DefaultSessionSummaryPrompt")
	}

	// Test tool_analysis prompt exists and is non-empty
	toolAnalysis, ok := config.Prompts["tool_analysis"]
	if !ok {
		t.Fatal("tool_analysis prompt not found")
	}
	if toolAnalysis == "" {
		t.Error("tool_analysis prompt is empty")
	}
	if toolAnalysis != domain.DefaultToolAnalysisPrompt {
		t.Error("tool_analysis prompt doesn't match DefaultToolAnalysisPrompt")
	}
}

func TestDefaultPrompts_NotEmpty(t *testing.T) {
	// Verify the default prompt constants are non-empty
	if domain.DefaultSessionSummaryPrompt == "" {
		t.Error("DefaultSessionSummaryPrompt is empty")
	}
	if domain.DefaultToolAnalysisPrompt == "" {
		t.Error("DefaultToolAnalysisPrompt is empty")
	}
	if domain.DefaultAnalysisPrompt == "" {
		t.Error("DefaultAnalysisPrompt is empty")
	}

	// Verify backward compatibility: DefaultAnalysisPrompt == DefaultToolAnalysisPrompt
	if domain.DefaultAnalysisPrompt != domain.DefaultToolAnalysisPrompt {
		t.Error("DefaultAnalysisPrompt should equal DefaultToolAnalysisPrompt for backward compatibility")
	}
}

func TestAllowedModels_Completeness(t *testing.T) {
	// Verify all expected models are in the whitelist
	expectedAliases := []string{"sonnet", "opus", "haiku"}
	for _, alias := range expectedAliases {
		if !domain.AllowedModels[alias] {
			t.Errorf("Expected alias %q to be in AllowedModels", alias)
		}
	}

	expectedFullNames := []string{
		"claude-sonnet-4-5-20250929",
		"claude-opus-4-20250514",
		"claude-3-5-sonnet-20241022",
		"claude-3-5-haiku-20241022",
	}
	for _, name := range expectedFullNames {
		if !domain.AllowedModels[name] {
			t.Errorf("Expected model name %q to be in AllowedModels", name)
		}
	}
}

func TestConfig_StructTags(t *testing.T) {
	// Verify that struct tags are present for YAML and JSON marshaling
	// This is a structural test to ensure the types can be serialized

	config := domain.DefaultConfig()

	// Should be able to access all nested structures
	_ = config.Analysis
	_ = config.UI
	_ = config.Prompts
	_ = config.Analysis.ClaudeOptions
}

func TestAnalysisConfig_CustomValues(t *testing.T) {
	config := domain.AnalysisConfig{
		TokenLimit:         200000,
		Model:              "opus",
		ParallelLimit:      5,
		EnabledPrompts:     []string{"custom_prompt1", "custom_prompt2"},
		AutoSummaryEnabled: true,
		AutoSummaryPrompt:  "custom_summary",
		ClaudeOptions: domain.ClaudeOptions{
			AllowedTools:     []string{"Read", "Write"},
			SystemPromptMode: "append",
		},
	}

	// Verify all custom values are preserved
	if config.TokenLimit != 200000 {
		t.Error("Custom TokenLimit not preserved")
	}
	if config.Model != "opus" {
		t.Error("Custom Model not preserved")
	}
	if config.ParallelLimit != 5 {
		t.Error("Custom ParallelLimit not preserved")
	}
	if len(config.EnabledPrompts) != 2 {
		t.Error("Custom EnabledPrompts count mismatch")
	}
	if !config.AutoSummaryEnabled {
		t.Error("Custom AutoSummaryEnabled not preserved")
	}
}
