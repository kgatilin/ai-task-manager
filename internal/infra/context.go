package infra

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ContextDetector determines the current context from environment and filesystem
type ContextDetector struct{}

// NewContextDetector creates a new context detector
func NewContextDetector() *ContextDetector {
	return &ContextDetector{}
}

// DetectContext determines the current context
func (d *ContextDetector) DetectContext() string {
	// Priority 1: Environment variable
	if ctx := os.Getenv("DW_CONTEXT"); ctx != "" {
		return ctx
	}

	// Priority 2: Parse from current working directory
	cwd, err := os.Getwd()
	if err == nil {
		if ctx := d.parseContextFromPath(cwd); ctx != "" {
			return ctx
		}
	}

	// Default to unknown
	return "unknown"
}

// parseContextFromPath extracts context information from a file path
func (d *ContextDetector) parseContextFromPath(path string) string {
	// Look for .darwinflow directory
	current := path
	for {
		dwPath := filepath.Join(current, ".darwinflow")
		if info, err := os.Stat(dwPath); err == nil && info.IsDir() {
			// Found .darwinflow directory, use the parent directory name as project
			projectName := filepath.Base(current)
			return "project/" + projectName
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached root
			break
		}
		current = parent
	}

	// Fallback: use the last directory component
	projectName := filepath.Base(path)
	if projectName != "" && projectName != "." && projectName != "/" {
		return "project/" + projectName
	}

	return ""
}

// NormalizeContent creates a human-readable text representation for full-text search
func NormalizeContent(eventType, payload string) string {
	// Parse the payload JSON to extract key information
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		// If we can't parse, fall back to simple combination
		return strings.TrimSpace(eventType + ": " + payload)
	}

	// Build human-readable content based on event type
	switch eventType {
	case "tool.invoked":
		tool := getString(data, "tool")
		// Parameters can be an object or string
		if params, ok := data["parameters"]; ok && params != nil {
			// Format parameters as JSON for readability
			paramsJSON, _ := json.Marshal(params)
			paramsStr := string(paramsJSON)
			// Truncate if too long
			if len(paramsStr) > 500 {
				paramsStr = paramsStr[:500] + "..."
			}
			return fmt.Sprintf("Tool: %s\nParameters: %s", tool, paramsStr)
		}
		return fmt.Sprintf("Tool: %s", tool)

	case "tool.result":
		tool := getString(data, "tool")
		result := getString(data, "result")
		if result != "" {
			return fmt.Sprintf("Tool: %s\nResult: %s", tool, result)
		}
		return fmt.Sprintf("Tool: %s completed", tool)

	case "chat.message.user", "chat.message.assistant":
		message := getString(data, "message")
		return message

	case "file.read", "file.written":
		filePath := getString(data, "file_path")
		changes := getString(data, "changes")
		if changes != "" {
			return fmt.Sprintf("File: %s\nChanges: %s", filePath, changes)
		}
		return fmt.Sprintf("File: %s", filePath)

	default:
		// For unknown event types, create a readable summary
		var parts []string
		for key, value := range data {
			if key == "context" {
				continue // Skip context field to reduce noise
			}
			if str := fmt.Sprintf("%v", value); str != "" && str != "<nil>" {
				parts = append(parts, fmt.Sprintf("%s: %v", key, value))
			}
		}
		return strings.Join(parts, "\n")
	}
}

// getString safely extracts a string value from a map
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		// Convert other types to string
		return fmt.Sprintf("%v", val)
	}
	return ""
}
