package domain

import "context"

// LLMOptions contains configuration options for LLM queries
type LLMOptions struct {
	Model           string
	MaxTokens       int
	Temperature     float64
	SystemPrompt    string
	SystemPromptMode string // "replace" or "append"
	AllowedTools    []string
}

// LLM defines the interface for language model interactions
// Implementations can use Claude CLI, Anthropic API, OpenAI, or other LLM providers
type LLM interface {
	// Query executes an LLM request with the given prompt and options
	// Returns the LLM's response as a string
	Query(ctx context.Context, prompt string, options *LLMOptions) (string, error)

	// EstimateTokens provides a rough estimate of tokens needed for the given prompt
	// Uses a simple heuristic (chars/4) for consistency across implementations
	EstimateTokens(prompt string) int

	// GetModel returns the current model being used by this LLM
	GetModel() string
}
