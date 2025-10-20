package pluginsdk

// TriggerType represents system hook event types that can trigger plugins
// These are the public SDK trigger types that plugins use to declare what events they respond to
type TriggerType string

const (
	// TriggerBeforeToolUse fires before any tool execution
	// Used by: Claude Code (PreToolUse hook)
	TriggerBeforeToolUse TriggerType = "trigger.tool.before"

	// TriggerAfterToolUse fires after tool execution
	TriggerAfterToolUse TriggerType = "trigger.tool.after"

	// TriggerUserInput fires when user provides input
	// Used by: Claude Code (UserPromptSubmit hook)
	TriggerUserInput TriggerType = "trigger.user.input"

	// TriggerSessionStart fires when session begins
	TriggerSessionStart TriggerType = "trigger.session.start"

	// TriggerSessionEnd fires when session ends
	// Used by: Claude Code (SessionEnd hook)
	TriggerSessionEnd TriggerType = "trigger.session.end"
)

// HookConfiguration describes a single hook provided by a plugin
type HookConfiguration struct {
	// TriggerType is the event type that triggers this hook
	// Examples: "trigger.tool.before", "trigger.user.input", "trigger.session.end"
	TriggerType string

	// Name is a human-readable name for the hook
	// Examples: "PreToolUse", "UserPromptSubmit"
	Name string

	// Description explains what this hook does
	Description string

	// Command is the CLI command to execute when the hook triggers
	// Example: "dw claude-code emit-event"
	Command string

	// Timeout is the maximum seconds this hook should take (0 = no timeout)
	Timeout int
}
