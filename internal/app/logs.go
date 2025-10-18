package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

// LogRecord represents a formatted log entry for display
type LogRecord struct {
	ID        string
	Timestamp time.Time
	EventType string
	Payload   json.RawMessage
	Content   string
}

// LogsService provides methods for querying and displaying logs
type LogsService struct {
	repo          domain.EventRepository
	rawExecutor   domain.RawQueryExecutor
}

// NewLogsService creates a new logs service
func NewLogsService(repo domain.EventRepository, rawExecutor domain.RawQueryExecutor) *LogsService {
	return &LogsService{
		repo:        repo,
		rawExecutor: rawExecutor,
	}
}

// ListRecentLogs retrieves the most recent N logs
func (s *LogsService) ListRecentLogs(ctx context.Context, limit int) ([]*LogRecord, error) {
	query := domain.EventQuery{
		Limit: limit,
	}

	events, err := s.repo.FindByQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}

	records := make([]*LogRecord, len(events))
	for i, event := range events {
		payloadBytes, err := event.MarshalPayload()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}

		records[i] = &LogRecord{
			ID:        event.ID,
			Timestamp: event.Timestamp,
			EventType: string(event.Type),
			Payload:   payloadBytes,
			Content:   event.Content,
		}
	}

	return records, nil
}

// ExecuteRawQuery executes an arbitrary SQL query
func (s *LogsService) ExecuteRawQuery(ctx context.Context, query string) (*domain.QueryResult, error) {
	return s.rawExecutor.ExecuteRawQuery(ctx, query)
}

// FormatLogRecord formats a single log record for display
func FormatLogRecord(index int, record *LogRecord) string {
	var output string

	output += fmt.Sprintf("[%d] %s\n", index+1, record.Timestamp.Format("2006-01-02 15:04:05.000"))
	output += fmt.Sprintf("    Event: %s\n", record.EventType)
	output += fmt.Sprintf("    ID: %s\n", record.ID)

	// Pretty print JSON payload with nested JSON expansion
	var payload interface{}
	if err := json.Unmarshal(record.Payload, &payload); err == nil {
		// Expand nested JSON strings in the payload
		expandedPayload := expandNestedJSON(payload)
		prettyPayload, _ := json.MarshalIndent(expandedPayload, "    ", "  ")
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

// expandNestedJSON recursively expands JSON strings within a data structure
func expandNestedJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// Recursively process map values
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = expandNestedJSON(value)
		}
		return result

	case []interface{}:
		// Recursively process array elements
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = expandNestedJSON(value)
		}
		return result

	case string:
		// Try to parse string as JSON
		if len(v) > 0 && (v[0] == '{' || v[0] == '[') {
			var parsed interface{}
			if err := json.Unmarshal([]byte(v), &parsed); err == nil {
				// Successfully parsed, recursively expand
				return expandNestedJSON(parsed)
			}
		}
		// Not JSON or parsing failed, return as-is
		return v

	default:
		// Return other types as-is
		return v
	}
}

// FormatQueryValue formats a value from a raw query result for display
func FormatQueryValue(val interface{}) string {
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
