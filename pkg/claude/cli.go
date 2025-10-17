package claude

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/kgatilin/darwinflow-pub/internal/events"
	"github.com/kgatilin/darwinflow-pub/internal/hooks"
)

// LogFromStdin reads hook input from stdin and logs the appropriate event
func LogFromStdin(eventTypeStr string, maxParamLength int) error {
	// Read hook input from stdin
	stdinData, err := io.ReadAll(os.Stdin)
	if err != nil {
		// Silently fail - don't disrupt Claude Code
		return nil
	}

	// Parse as hook input if it's JSON
	var hookInput hooks.HookInput
	var payload interface{}

	if len(stdinData) > 0 {
		if err := json.Unmarshal(stdinData, &hookInput); err != nil {
			// If not valid hook input, try as generic payload
			if err := json.Unmarshal(stdinData, &payload); err != nil {
				// Not valid JSON, use raw string
				payload = map[string]interface{}{
					"raw": string(stdinData),
				}
			}
		}
	}

	// Map event type string to events.Type
	eventType := mapEventType(eventTypeStr)

	// Create logger
	logger, err := NewLogger(DefaultDBPath)
	if err != nil {
		// Silently fail - don't disrupt Claude Code
		return nil
	}
	defer logger.Close()

	// Log event
	if hookInput.SessionID != "" {
		// We have valid hook input
		if err := logger.LogFromHookInput(hookInput, eventType, maxParamLength); err != nil {
			// Silently fail
			return nil
		}
	} else if payload != nil {
		// We have generic payload
		if err := logger.LogEvent(eventType, payload); err != nil {
			// Silently fail
			return nil
		}
	} else {
		// No input, log minimal event
		minimalPayload := map[string]interface{}{
			"timestamp": hookInput.SessionID,
		}
		if err := logger.LogEvent(eventType, minimalPayload); err != nil {
			// Silently fail
			return nil
		}
	}

	return nil
}

// mapEventType maps string event types to events.Type
func mapEventType(eventTypeStr string) events.Type {
	// Normalize the string
	normalized := strings.ToLower(strings.ReplaceAll(eventTypeStr, "_", "."))

	switch normalized {
	case "chat.started":
		return events.ChatStarted
	case "chat.ended", "chat.end":
		return events.ChatStarted // Reuse for now
	case "chat.message.user", "user.message":
		return events.ChatMessageUser
	case "chat.message.assistant", "assistant.message":
		return events.ChatMessageAssistant
	case "tool.invoked", "tool.invoke":
		return events.ToolInvoked
	case "tool.result":
		return events.ToolResult
	case "file.read":
		return events.FileRead
	case "file.written", "file.write":
		return events.FileWritten
	case "context.changed", "context.change":
		return events.ContextChanged
	case "error":
		return events.Error
	default:
		// Default to generic event
		return events.Type(normalized)
	}
}
