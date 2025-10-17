package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/kgatilin/darwinflow-pub/pkg/claude"
)

type logsOptions struct {
	limit int
	query string
	help  bool
}

func parseLogsFlags(args []string) (*logsOptions, error) {
	fs := flag.NewFlagSet("logs", flag.ContinueOnError)
	opts := &logsOptions{}

	fs.IntVar(&opts.limit, "limit", 20, "Number of most recent logs to display")
	fs.StringVar(&opts.query, "query", "", "Arbitrary SQL query to execute")
	fs.BoolVar(&opts.help, "help", false, "Show help and database schema")

	fs.Usage = printLogsUsage

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return opts, nil
}

func handleLogs(args []string) {
	opts, err := parseLogsFlags(args)
	if err != nil {
		os.Exit(1)
	}

	// Show help if requested
	if opts.help {
		printLogsHelp()
		return
	}

	dbPath := claude.DefaultDBPath

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Database not found at %s\n", dbPath)
		fmt.Fprintf(os.Stderr, "Run 'dw claude init' to initialize logging.\n")
		os.Exit(1)
	}

	// Handle arbitrary SQL query
	if opts.query != "" {
		if err := executeRawQuery(dbPath, opts.query); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle standard log listing
	if err := listLogs(dbPath, opts.limit); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printLogsUsage() {
	fmt.Println("Usage: dw logs [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --limit N       Number of most recent logs to display (default: 20)")
	fmt.Println("  --query SQL     Execute an arbitrary SQL query")
	fmt.Println("  --help          Show help and database schema")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  dw logs                                    # Show 20 most recent logs")
	fmt.Println("  dw logs --limit 50                         # Show 50 most recent logs")
	fmt.Println("  dw logs --query \"SELECT * FROM events\"     # Run custom SQL query")
	fmt.Println()
}

func printLogsHelp() {
	fmt.Println("DarwinFlow Logs - Database Schema")
	fmt.Println()
	fmt.Println("DATABASE STRUCTURE:")
	fmt.Println()
	fmt.Println("Table: events")
	fmt.Println("  Columns:")
	fmt.Println("    - id           TEXT PRIMARY KEY    Unique event identifier (UUID)")
	fmt.Println("    - timestamp    INTEGER NOT NULL    Unix timestamp in milliseconds")
	fmt.Println("    - event_type   TEXT NOT NULL       Event type (e.g., 'tool.invoked', 'chat.message.user')")
	fmt.Println("    - payload      TEXT NOT NULL       JSON payload with event-specific data")
	fmt.Println("    - content      TEXT NOT NULL       Normalized searchable content")
	fmt.Println()
	fmt.Println("  Indexes:")
	fmt.Println("    - idx_events_timestamp          ON events(timestamp)")
	fmt.Println("    - idx_events_type               ON events(event_type)")
	fmt.Println("    - idx_events_timestamp_type     ON events(timestamp, event_type)")
	fmt.Println()
	fmt.Println("FTS5 Virtual Table: events_fts (if available)")
	fmt.Println("  Full-text search on content field")
	fmt.Println()
	fmt.Println("COMMON EVENT TYPES:")
	fmt.Println("  - tool.invoked              Tool was invoked (Read, Write, Bash, etc.)")
	fmt.Println("  - tool.result               Tool execution completed")
	fmt.Println("  - chat.message.user         User sent a message")
	fmt.Println("  - chat.message.assistant    Assistant sent a message")
	fmt.Println("  - chat.started              Chat session started")
	fmt.Println("  - file.read                 File was read")
	fmt.Println("  - file.written              File was written")
	fmt.Println("  - context.changed           Context changed")
	fmt.Println("  - error                     Error occurred")
	fmt.Println()
	fmt.Println("QUERY EXAMPLES:")
	fmt.Println("  # Count events by type")
	fmt.Println("  dw logs --query \"SELECT event_type, COUNT(*) as count FROM events GROUP BY event_type ORDER BY count DESC\"")
	fmt.Println()
	fmt.Println("  # Find all tool invocations in last hour")
	fmt.Printf("  dw logs --query \"SELECT * FROM events WHERE event_type = 'tool.invoked' AND timestamp > strftime('%%s', 'now', '-1 hour') * 1000\"\n")
	fmt.Println()
	fmt.Println("  # Search content for specific text")
	fmt.Printf("  dw logs --query \"SELECT * FROM events WHERE content LIKE '%%sqlite%%' LIMIT 10\"\n")
	fmt.Println()
	fmt.Println("Database location:", claude.DefaultDBPath)
	fmt.Println()
}

type logRecord struct {
	ID        string
	Timestamp int64
	EventType string
	Payload   []byte
	Content   string
}

func queryLogs(dbPath string, limit int) ([]logRecord, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	query := "SELECT id, timestamp, event_type, payload, content FROM events ORDER BY timestamp DESC LIMIT ?"
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var records []logRecord

	for rows.Next() {
		var r logRecord
		var payloadStr string

		if err := rows.Scan(&r.ID, &r.Timestamp, &r.EventType, &payloadStr, &r.Content); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		r.Payload = []byte(payloadStr)
		records = append(records, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

func formatLogRecord(i int, record logRecord) string {
	var output string

	timestamp := time.UnixMilli(record.Timestamp)
	output += fmt.Sprintf("[%d] %s\n", i+1, timestamp.Format("2006-01-02 15:04:05.000"))
	output += fmt.Sprintf("    Event: %s\n", record.EventType)
	output += fmt.Sprintf("    ID: %s\n", record.ID)

	// Pretty print JSON payload
	var payload interface{}
	if err := json.Unmarshal(record.Payload, &payload); err == nil {
		prettyPayload, _ := json.MarshalIndent(payload, "    ", "  ")
		output += fmt.Sprintf("    Payload: %s\n", string(prettyPayload))
	} else {
		output += fmt.Sprintf("    Payload: %s\n", string(record.Payload))
	}

	if record.Content != "" {
		// Truncate content if too long
		content := record.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		output += fmt.Sprintf("    Content: %s\n", content)
	}

	output += "\n"
	return output
}

func listLogs(dbPath string, limit int) error {
	records, err := queryLogs(dbPath, limit)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Println("No logs found.")
		fmt.Println("Run 'dw claude init' to initialize logging, then use Claude Code to generate events.")
		return nil
	}

	// Display logs
	fmt.Printf("Showing %d most recent logs:\n\n", len(records))

	for i, record := range records {
		fmt.Print(formatLogRecord(i, record))
	}

	return nil
}

func formatQueryValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "NULL"
	case []byte:
		// Try to parse as JSON for pretty printing
		var jsonObj interface{}
		if err := json.Unmarshal(v, &jsonObj); err == nil {
			jsonBytes, _ := json.Marshal(jsonObj)
			str := string(jsonBytes)
			if len(str) > 100 {
				str = str[:100] + "..."
			}
			return str
		}
		str := string(v)
		if len(str) > 100 {
			str = str[:100] + "..."
		}
		return str
	case string:
		if len(v) > 100 {
			return v[:100] + "..."
		}
		return v
	case int64:
		// Check if it might be a timestamp (13 digits for milliseconds)
		if v > 1000000000000 && v < 9999999999999 {
			t := time.UnixMilli(v)
			return fmt.Sprintf("%d (%s)", v, t.Format("2006-01-02 15:04:05"))
		}
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func executeRawQuery(dbPath string, query string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	// Print column headers
	for i, col := range columns {
		if i > 0 {
			fmt.Print(" | ")
		}
		fmt.Print(col)
	}
	fmt.Println()
	fmt.Println(repeatString("-", 80))

	// Prepare scan destinations
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Print rows
	rowCount := 0
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		for i, val := range values {
			if i > 0 {
				fmt.Print(" | ")
			}
			fmt.Print(formatQueryValue(val))
		}
		fmt.Println()
		rowCount++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	fmt.Println()
	fmt.Printf("(%d rows)\n", rowCount)
	return nil
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
