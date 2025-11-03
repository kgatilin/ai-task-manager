package task_manager

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	// FormatTable is the default table format
	FormatTable OutputFormat = "table"
	// FormatLLM is the structured text format optimized for LLM consumption
	FormatLLM OutputFormat = "llm"
	// FormatJSON is the JSON format
	FormatJSON OutputFormat = "json"
)

// ParseOutputFormat parses a format string and returns the OutputFormat
func ParseOutputFormat(format string) (OutputFormat, error) {
	switch strings.ToLower(format) {
	case "", "table":
		return FormatTable, nil
	case "llm":
		return FormatLLM, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("invalid format: %s (must be table, llm, or json)", format)
	}
}

// OutputFormatter provides methods for formatting output in different formats
type OutputFormatter struct {
	writer io.Writer
	format OutputFormat
}

// NewOutputFormatter creates a new output formatter
func NewOutputFormatter(writer io.Writer, format OutputFormat) *OutputFormatter {
	return &OutputFormatter{
		writer: writer,
		format: format,
	}
}

// OutputJSON marshals data to JSON and writes it
func (f *OutputFormatter) OutputJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintf(f.writer, "%s\n", jsonData)
	return err
}

// ContextualHints provides contextual next-step suggestions after command output
type ContextualHints struct {
	hints []string
}

// NewContextualHints creates a new contextual hints builder
func NewContextualHints() *ContextualHints {
	return &ContextualHints{
		hints: make([]string, 0),
	}
}

// Add adds a hint to the list
func (h *ContextualHints) Add(hint string) {
	h.hints = append(h.hints, hint)
}

// Output writes the hints to the writer
func (h *ContextualHints) Output(writer io.Writer) {
	if len(h.hints) == 0 {
		return
	}
	fmt.Fprintf(writer, "\nNext steps:\n")
	for _, hint := range h.hints {
		fmt.Fprintf(writer, "  %s\n", hint)
	}
}
