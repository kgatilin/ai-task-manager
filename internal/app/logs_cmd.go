package app

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// LogsServiceInterface defines the interface for logs operations
type LogsServiceInterface interface {
	ListRecentLogs(ctx context.Context, limit, sessionLimit int, sessionID string, ordered bool) ([]*LogRecord, error)
	ExecuteRawQuery(ctx context.Context, query string) (*pluginsdk.QueryResult, error)
}

// LogsCommandHandler handles the logs command presentation logic
type LogsCommandHandler struct {
	service LogsServiceInterface
	out     io.Writer
}

// NewLogsCommandHandler creates a new logs command handler
func NewLogsCommandHandler(service LogsServiceInterface, out io.Writer) *LogsCommandHandler {
	return &LogsCommandHandler{
		service: service,
		out:     out,
	}
}

// ListLogs displays logs based on the provided options
func (h *LogsCommandHandler) ListLogs(ctx context.Context, limit, sessionLimit int, sessionID string, ordered bool, format string) error {
	records, err := h.service.ListRecentLogs(ctx, limit, sessionLimit, sessionID, ordered)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Fprintln(h.out, "No logs found.")
		fmt.Fprintln(h.out, "Run 'dw init' or a plugin's init command to initialize logging.")
		return nil
	}

	// Handle CSV format
	if format == "csv" {
		return FormatLogsAsCSV(h.out, records)
	}

	// Handle Markdown format
	if format == "markdown" {
		return FormatLogsAsMarkdown(h.out, records)
	}

	// Validate format
	if format != "text" && format != "" {
		return fmt.Errorf("invalid format '%s'. Valid formats: text, csv, markdown", format)
	}

	// Display logs in text format
	if sessionID != "" {
		fmt.Fprintf(h.out, "Showing %d logs for session %s:\n\n", len(records), sessionID)
	} else if sessionLimit > 0 {
		fmt.Fprintf(h.out, "Showing %d logs from %d most recent sessions:\n\n", len(records), sessionLimit)
	} else {
		fmt.Fprintf(h.out, "Showing %d most recent logs:\n\n", len(records))
	}

	for i, record := range records {
		fmt.Fprint(h.out, FormatLogRecord(i, record))
	}

	return nil
}

// ExecuteRawQuery executes a raw SQL query and displays the results
func (h *LogsCommandHandler) ExecuteRawQuery(ctx context.Context, query string) error {
	result, err := h.service.ExecuteRawQuery(ctx, query)
	if err != nil {
		return err
	}

	// Print column headers
	for i, col := range result.Columns {
		if i > 0 {
			fmt.Fprint(h.out, " | ")
		}
		fmt.Fprint(h.out, col)
	}
	fmt.Fprintln(h.out)
	fmt.Fprintln(h.out, repeatString("-", 80))

	// Print rows
	for _, row := range result.Rows {
		for i, val := range row {
			if i > 0 {
				fmt.Fprint(h.out, " | ")
			}
			fmt.Fprint(h.out, FormatQueryValue(val))
		}
		fmt.Fprintln(h.out)
	}

	fmt.Fprintln(h.out)
	fmt.Fprintf(h.out, "(%d rows)\n", len(result.Rows))
	return nil
}

func repeatString(s string, count int) string {
	return strings.Repeat(s, count)
}
