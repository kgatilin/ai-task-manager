package infra

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// TranscriptEntry represents a line in the Claude Code transcript (JSONL format)
type TranscriptEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`

	// Tool use specific fields
	Type       string                 `json:"type,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Input      map[string]interface{} `json:"input,omitempty"`
}

// TranscriptParser parses Claude Code transcript files
type TranscriptParser struct{}

// NewTranscriptParser creates a new transcript parser
func NewTranscriptParser() *TranscriptParser {
	return &TranscriptParser{}
}

// Parse reads the transcript file and returns all entries
func (p *TranscriptParser) Parse(path string) ([]TranscriptEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open transcript: %w", err)
	}
	defer file.Close()

	var entries []TranscriptEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var entry TranscriptEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Skip malformed lines
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading transcript: %w", err)
	}

	return entries, nil
}

// ExtractLastToolUse extracts tool name and parameters from the last tool use in transcript
func (p *TranscriptParser) ExtractLastToolUse(transcriptPath string, maxParamLength int) (string, string, error) {
	entries, err := p.Parse(transcriptPath)
	if err != nil {
		return "", "", err
	}

	// Find the last tool use entry
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		if entry.Type == "tool_use" || entry.Name != "" {
			toolName := entry.Name
			if toolName == "" {
				toolName = "unknown"
			}

			// Serialize parameters
			var params map[string]interface{}
			if entry.Parameters != nil {
				params = entry.Parameters
			} else if entry.Input != nil {
				params = entry.Input
			}

			paramsJSON, _ := json.Marshal(params)
			paramsStr := string(paramsJSON)

			// Trim if too long
			if len(paramsStr) > maxParamLength {
				paramsStr = paramsStr[:maxParamLength] + "..."
			}

			return toolName, paramsStr, nil
		}
	}

	return "", "", fmt.Errorf("no tool use found in transcript")
}

// ExtractLastUserMessage extracts the last user message from transcript
func (p *TranscriptParser) ExtractLastUserMessage(transcriptPath string) (string, error) {
	entries, err := p.Parse(transcriptPath)
	if err != nil {
		return "", err
	}

	// Find the last user message
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if entry.Role == "user" && entry.Content != "" {
			return entry.Content, nil
		}
	}

	return "", fmt.Errorf("no user message found in transcript")
}

// ExtractLastAssistantMessage extracts the last assistant message from transcript
func (p *TranscriptParser) ExtractLastAssistantMessage(transcriptPath string) (string, error) {
	entries, err := p.Parse(transcriptPath)
	if err != nil {
		return "", err
	}

	// Find the last assistant message
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if entry.Role == "assistant" && entry.Content != "" {
			// Truncate if very long
			content := entry.Content
			if len(content) > 1000 {
				content = content[:1000] + "... (truncated)"
			}
			return content, nil
		}
	}

	return "", fmt.Errorf("no assistant message found in transcript")
}
