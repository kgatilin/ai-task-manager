package hooks

// HookConfig represents the hooks configuration for Claude Code
type HookConfig struct {
	Hooks map[string][]HookMatcher `json:"hooks"`
}

// HookMatcher represents a hook matcher and its associated hooks
type HookMatcher struct {
	Matcher string       `json:"matcher,omitempty"`
	Hooks   []HookAction `json:"hooks"`
}

// HookAction represents a single hook action
type HookAction struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
}

// HookInput represents the standardized input passed to hooks via stdin
type HookInput struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	HookEventName  string `json:"hook_event_name"`
}

// EventType represents the Claude Code hook event types
type EventType string

const (
	PreToolUse      EventType = "PreToolUse"
	PostToolUse     EventType = "PostToolUse"
	Notification    EventType = "Notification"
	UserPromptSubmit EventType = "UserPromptSubmit"
	Stop            EventType = "Stop"
	SubagentStop    EventType = "SubagentStop"
	PreCompact      EventType = "PreCompact"
	SessionStart    EventType = "SessionStart"
	SessionEnd      EventType = "SessionEnd"
)

// DefaultConfig returns the default hooks configuration for DarwinFlow logging
func DefaultConfig() HookConfig {
	return HookConfig{
		Hooks: map[string][]HookMatcher{
			string(PreToolUse): {
				{
					Matcher: "*", // Match all tools
					Hooks: []HookAction{
						{
							Type:    "command",
							Command: "dw claude log tool.invoked",
							Timeout: 5,
						},
					},
				},
			},
			string(UserPromptSubmit): {
				{
					Hooks: []HookAction{
						{
							Type:    "command",
							Command: "dw claude log chat.message.user",
							Timeout: 5,
						},
					},
				},
			},
		},
	}
}

// MergeConfig merges new hooks into existing configuration
func MergeConfig(existing, new HookConfig) HookConfig {
	merged := HookConfig{
		Hooks: make(map[string][]HookMatcher),
	}

	// Copy existing hooks
	for event, matchers := range existing.Hooks {
		merged.Hooks[event] = append([]HookMatcher{}, matchers...)
	}

	// Add new hooks
	for event, newMatchers := range new.Hooks {
		if existing, ok := merged.Hooks[event]; ok {
			// Event exists, merge matchers
			merged.Hooks[event] = mergeMatchers(existing, newMatchers)
		} else {
			// New event, add all matchers
			merged.Hooks[event] = newMatchers
		}
	}

	return merged
}

// mergeMatchers merges new matchers into existing ones
func mergeMatchers(existing, new []HookMatcher) []HookMatcher {
	result := append([]HookMatcher{}, existing...)

	for _, newMatcher := range new {
		found := false
		for i, existingMatcher := range result {
			if existingMatcher.Matcher == newMatcher.Matcher {
				// Matcher exists, merge hooks
				result[i].Hooks = mergeHooks(existingMatcher.Hooks, newMatcher.Hooks)
				found = true
				break
			}
		}
		if !found {
			// New matcher, add it
			result = append(result, newMatcher)
		}
	}

	return result
}

// mergeHooks merges new hooks into existing ones, avoiding duplicates
func mergeHooks(existing, new []HookAction) []HookAction {
	result := append([]HookAction{}, existing...)

	for _, newHook := range new {
		duplicate := false
		for _, existingHook := range result {
			if existingHook.Command == newHook.Command {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, newHook)
		}
	}

	return result
}
