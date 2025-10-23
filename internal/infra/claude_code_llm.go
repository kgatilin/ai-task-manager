package infra

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

// ClaudeCodeLLM implements the domain.LLM interface using the Claude CLI tool
// It wraps the claude command-line tool for AI interactions
type ClaudeCodeLLM struct {
	logger *Logger
	config *domain.Config
}

// NewClaudeCodeLLM creates a new Claude Code LLM instance with default config
func NewClaudeCodeLLM(logger *Logger) *ClaudeCodeLLM {
	if logger == nil {
		logger = NewDefaultLogger()
	}
	return &ClaudeCodeLLM{
		logger: logger,
		config: domain.DefaultConfig(),
	}
}

// NewClaudeCodeLLMWithConfig creates a new Claude Code LLM instance with custom config
func NewClaudeCodeLLMWithConfig(logger *Logger, config *domain.Config) *ClaudeCodeLLM {
	if logger == nil {
		logger = NewDefaultLogger()
	}
	if config == nil {
		config = domain.DefaultConfig()
	}
	return &ClaudeCodeLLM{
		logger: logger,
		config: config,
	}
}

// Query executes an LLM query using the Claude CLI tool
// Streams output to stderr in real-time for progress visibility
// The prompt parameter is treated as the user prompt unless options specify system prompt mode
func (l *ClaudeCodeLLM) Query(ctx context.Context, prompt string, options *domain.LLMOptions) (string, error) {
	// Use provided options or create defaults
	if options == nil {
		options = &domain.LLMOptions{}
	}

	// Build command arguments
	args := []string{}

	// Apply model from options or config
	model := options.Model
	if model == "" {
		model = l.config.Analysis.Model
	}
	if model != "" {
		args = append(args, "--model", model)
	}

	var userPrompt string

	// Apply system prompt mode from options or config
	systemPromptMode := options.SystemPromptMode
	if systemPromptMode == "" {
		systemPromptMode = l.config.Analysis.ClaudeOptions.SystemPromptMode
	}

	if systemPromptMode == "replace" {
		args = append(args, "--system-prompt", prompt)
		// When using system prompt, we need a user prompt too
		userPrompt = "Analyze the session data provided in the system prompt."
	} else if systemPromptMode == "append" {
		args = append(args, "--append-system-prompt", prompt)
		userPrompt = "Analyze the session data."
	} else {
		// No system prompt mode, use prompt directly
		userPrompt = prompt
	}

	// Apply allowed tools from options or config
	allowedTools := options.AllowedTools
	if len(allowedTools) == 0 {
		allowedTools = l.config.Analysis.ClaudeOptions.AllowedTools
	}
	if len(allowedTools) > 0 {
		args = append(args, "--allowed-tools", strings.Join(allowedTools, ","))
	}

	// Add the user prompt last
	args = append(args, userPrompt)

	if l.logger != nil {
		l.logger.Debug("Executing: claude %s", strings.Join(args, " "))
	}
	cmd := exec.CommandContext(ctx, "claude", args...)

	var stdout, stderr bytes.Buffer

	// Use MultiWriter to stream to both the buffer and os.Stderr for real-time feedback
	cmd.Stdout = io.MultiWriter(&stdout, os.Stderr)
	cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)

	if l.logger != nil {
		l.logger.Debug("Running Claude CLI command...")
	}
	if err := cmd.Run(); err != nil {
		if l.logger != nil {
			l.logger.Error("Claude CLI command failed: %v", err)
		}
		return "", fmt.Errorf("claude command failed: %w, stderr: %s", err, stderr.String())
	}
	if l.logger != nil {
		l.logger.Debug("Claude CLI command completed successfully")
	}

	return strings.TrimSpace(stdout.String()), nil
}

// EstimateTokens provides a rough estimate of tokens needed for the given prompt
// Uses the simple heuristic: ~4 characters per token (common approximation for Claude models)
func (l *ClaudeCodeLLM) EstimateTokens(prompt string) int {
	// Estimate tokens: ~4 characters per token (conservative estimate)
	return len(prompt) / 4
}

// GetModel returns the currently configured model
func (l *ClaudeCodeLLM) GetModel() string {
	return l.config.Analysis.Model
}
