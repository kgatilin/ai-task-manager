package claude_code

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// Ensure plugin implements SDK ICommandProvider
var _ pluginsdk.ICommandProvider = (*ClaudeCodePlugin)(nil)

// GetCommands returns the CLI commands provided by this plugin (SDK interface)
func (p *ClaudeCodePlugin) GetCommands() []pluginsdk.Command {
	return []pluginsdk.Command{
		&InitCommand{plugin: p},
		&EmitEventCommand{plugin: p},
		&LogCommand{plugin: p},
		&AutoSummaryCommand{plugin: p},
		&AutoSummaryExecCommand{plugin: p},
		&SessionSummaryCommand{plugin: p},
	}
}

// InitCommand initializes Claude Code logging infrastructure
type InitCommand struct {
	plugin *ClaudeCodePlugin
}

func (c *InitCommand) GetName() string {
	return "init"
}

func (c *InitCommand) GetDescription() string {
	return "Initialize Claude Code logging infrastructure"
}

func (c *InitCommand) GetUsage() string {
	return "init"
}

func (c *InitCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	if c.plugin.handler == nil {
		return fmt.Errorf("handler not initialized")
	}

	// Use the default database path
	return c.plugin.handler.Init(ctx, c.plugin.dbPath)
}

// EmitEventCommand emits an event via the plugin SDK context
// This command reads a structured event from stdin and emits it through the plugin context.
// All errors are logged but never propagated - this ensures hook execution is never disrupted.
//
// Input format (JSON from stdin):
//
//	{
//	  "type": "tool.invoked",
//	  "source": "claude-code",
//	  "timestamp": "2025-10-20T10:30:00Z",
//	  "payload": { "tool": "Read", "parameters": {...} },
//	  "metadata": { "session_id": "abc123", "cwd": "/workspace" },
//	  "version": "1.0"
//	}
//
// Required fields: type, source, metadata.session_id
// The command validates input and emits to the framework's event store.
type EmitEventCommand struct {
	plugin *ClaudeCodePlugin
}

func NewEmitEventCommand(plugin *ClaudeCodePlugin) *EmitEventCommand {
	return &EmitEventCommand{
		plugin: plugin,
	}
}

func (c *EmitEventCommand) GetName() string {
	return "emit-event"
}

func (c *EmitEventCommand) GetDescription() string {
	return "Emit an event via plugin context (reads JSON from stdin)"
}

func (c *EmitEventCommand) GetUsage() string {
	return "emit-event"
}

func (c *EmitEventCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Add timeout to prevent infinite hangs
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Safely recover from panics
	defer func() {
		if r := recover(); r != nil {
			c.plugin.logger.Error("emit-event: panic recovered: %v", r)
		}
	}()

	// Read stdin
	stdinData, err := io.ReadAll(cmdCtx.GetStdin())
	if err != nil {
		c.plugin.logger.Debug("emit-event: failed to read stdin: %v", err)
		return nil // Silently fail - don't disrupt Claude Code
	}

	if len(stdinData) == 0 {
		c.plugin.logger.Debug("emit-event: empty stdin")
		return nil // Silently fail - don't disrupt Claude Code
	}

	// Parse JSON into pluginsdk.Event
	var event pluginsdk.Event
	if err := json.Unmarshal(stdinData, &event); err != nil {
		c.plugin.logger.Debug("emit-event: invalid JSON: %v", err)
		return nil // Silently fail - don't disrupt Claude Code
	}

	// Validate required fields
	if event.Type == "" {
		c.plugin.logger.Debug("emit-event: missing required field: type")
		return nil
	}

	if event.Source == "" {
		c.plugin.logger.Debug("emit-event: missing required field: source")
		return nil
	}

	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}

	sessionID, ok := event.Metadata["session_id"]
	if !ok || sessionID == "" {
		c.plugin.logger.Debug("emit-event: missing required field: metadata.session_id")
		return nil
	}

	// Set default timestamp if missing
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Set default version if missing
	if event.Version == "" {
		event.Version = "1.0"
	}

	// Initialize empty payload if nil
	if event.Payload == nil {
		event.Payload = make(map[string]interface{})
	}

	// Emit event via plugin context (silently fail if DB error)
	if err := cmdCtx.EmitEvent(ctxWithTimeout, event); err != nil {
		c.plugin.logger.Debug("emit-event: failed to emit event: %v", err)
		return nil // Silently fail - don't disrupt Claude Code
	}

	return nil
}

// LogCommand logs a Claude Code event from hook input
// DEPRECATED: Use EmitEventCommand instead (will be removed in v2.0)
type LogCommand struct {
	plugin *ClaudeCodePlugin
}

func (c *LogCommand) GetName() string {
	return "log"
}

func (c *LogCommand) GetDescription() string {
	return "Log a Claude Code event (reads JSON from stdin)"
}

func (c *LogCommand) GetUsage() string {
	return "log <event-type>"
}

func (c *LogCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("event type required")
	}

	if c.plugin.handler == nil {
		return fmt.Errorf("handler not initialized")
	}

	eventTypeStr := args[0]

	// Get max param length from environment or use default
	maxParamLength := 30
	if envVal := os.Getenv("DW_MAX_PARAM_LENGTH"); envVal != "" {
		if parsed, err := fmt.Sscanf(envVal, "%d", &maxParamLength); err == nil && parsed == 1 {
			// Successfully parsed
		}
	}

	// Read stdin data from command context
	stdinData, err := io.ReadAll(cmdCtx.GetStdin())
	if err != nil {
		// Silently fail - don't disrupt Claude Code
		return nil
	}

	// Execute (silently - errors shouldn't disrupt Claude Code)
	_ = c.plugin.handler.Log(ctx, eventTypeStr, stdinData, maxParamLength)
	return nil
}

// AutoSummaryCommand handles auto-triggered session summaries on SessionEnd
type AutoSummaryCommand struct {
	plugin *ClaudeCodePlugin
}

func (c *AutoSummaryCommand) GetName() string {
	return "auto-summary"
}

func (c *AutoSummaryCommand) GetDescription() string {
	return "Auto-trigger session summary (called by SessionEnd hook)"
}

func (c *AutoSummaryCommand) GetUsage() string {
	return "auto-summary"
}

func (c *AutoSummaryCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	if c.plugin.handler == nil {
		return fmt.Errorf("handler not initialized")
	}

	// Read stdin data from command context
	stdinData, err := io.ReadAll(cmdCtx.GetStdin())
	if err != nil {
		// Silently fail - don't disrupt Claude Code
		return nil
	}

	// Execute (silently - errors shouldn't disrupt Claude Code)
	_ = c.plugin.handler.AutoSummary(ctx, stdinData)
	return nil
}

// AutoSummaryExecCommand executes the actual summary analysis in background
type AutoSummaryExecCommand struct {
	plugin *ClaudeCodePlugin
}

func (c *AutoSummaryExecCommand) GetName() string {
	return "auto-summary-exec"
}

func (c *AutoSummaryExecCommand) GetDescription() string {
	return "Internal: Execute summary in background (do not call directly)"
}

func (c *AutoSummaryExecCommand) GetUsage() string {
	return "auto-summary-exec <session-id>"
}

func (c *AutoSummaryExecCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	if len(args) < 1 {
		// No session ID provided
		return nil
	}

	if c.plugin.handler == nil {
		return fmt.Errorf("handler not initialized")
	}

	sessionID := args[0]

	// Execute (silently - errors shouldn't disrupt background analysis)
	_ = c.plugin.handler.AutoSummaryExec(ctx, sessionID)
	return nil
}

// SessionSummaryCommand provides a quick summary of a session
type SessionSummaryCommand struct {
	plugin *ClaudeCodePlugin
}

func (c *SessionSummaryCommand) GetName() string {
	return "session-summary"
}

func (c *SessionSummaryCommand) GetDescription() string {
	return "Display a summary of a Claude Code session"
}

func (c *SessionSummaryCommand) GetUsage() string {
	return "session-summary --session-id <id> | --last"
}

func (c *SessionSummaryCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	sessionID := ""
	last := false

	// Simple flag parsing
	for i := 0; i < len(args); i++ {
		if args[i] == "--session-id" && i+1 < len(args) {
			sessionID = args[i+1]
			i++
		} else if args[i] == "--last" {
			last = true
		}
	}

	// Determine which session to summarize
	var targetSessionID string
	if last {
		lastID, err := c.plugin.analysisService.GetLastSession(ctx)
		if err != nil {
			return fmt.Errorf("failed to get last session: %w", err)
		}
		targetSessionID = lastID
	} else if sessionID != "" {
		targetSessionID = sessionID
	} else {
		return fmt.Errorf("must specify either --session-id or --last")
	}

	// Get session entity
	entity, err := c.plugin.buildSessionEntity(ctx, targetSessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Get output writer from command context
	out := cmdCtx.GetStdout()

	// Display summary
	fmt.Fprintln(out, "Session Summary")
	fmt.Fprintln(out, "===============")
	fmt.Fprintf(out, "Session ID: %s\n", entity.GetID())
	fmt.Fprintf(out, "Event Count: %d\n", entity.GetField("event_count"))
	fmt.Fprintf(out, "First Event: %v\n", entity.GetField("first_event"))
	fmt.Fprintf(out, "Last Event: %v\n", entity.GetField("last_event"))
	fmt.Fprintf(out, "Token Count: ~%d\n", entity.GetField("token_count"))
	fmt.Fprintf(out, "Status: %s\n", entity.GetStatus())

	// Display analyses if available
	analyses := entity.GetAnalyses()
	if len(analyses) > 0 {
		fmt.Fprintf(out, "\nAnalyses: %d\n", len(analyses))
		for i, analysis := range analyses {
			fmt.Fprintf(out, "  [%d] %s (%s)\n", i+1, analysis.PromptName, analysis.ModelUsed)
			if analysis.PatternsSummary != "" {
				fmt.Fprintf(out, "      Summary: %s\n", analysis.PatternsSummary)
			}
		}
	} else {
		fmt.Fprintln(out, "\nNo analyses available")
	}

	return nil
}
