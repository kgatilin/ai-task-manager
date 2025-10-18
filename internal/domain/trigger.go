package domain

// TriggerType represents the Claude Code hook event types that can trigger event logging
// These are domain concepts representing the points in the interaction where events occur
type TriggerType string

const (
	TriggerPreToolUse       TriggerType = "PreToolUse"
	TriggerPostToolUse      TriggerType = "PostToolUse"
	TriggerNotification     TriggerType = "Notification"
	TriggerUserPromptSubmit TriggerType = "UserPromptSubmit"
	TriggerStop             TriggerType = "Stop"
	TriggerSubagentStop     TriggerType = "SubagentStop"
	TriggerPreCompact       TriggerType = "PreCompact"
	TriggerSessionStart     TriggerType = "SessionStart"
	TriggerSessionEnd       TriggerType = "SessionEnd"
)
