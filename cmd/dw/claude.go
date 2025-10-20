package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kgatilin/darwinflow-pub/internal/app"
	"github.com/kgatilin/darwinflow-pub/internal/infra"
)

func handleClaudeCommand(args []string) {
	if len(args) < 1 {
		printClaudeUsage()
		os.Exit(1)
	}

	subcommand := args[0]

	switch subcommand {
	case "init":
		handleInit(args[1:])
	case "log":
		handleLog(args[1:])
	case "auto-summary":
		handleAutoSummary(args[1:])
	case "auto-summary-exec":
		handleAutoSummaryExec(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown claude subcommand: %s\n\n", subcommand)
		printClaudeUsage()
		os.Exit(1)
	}
}

func printClaudeUsage() {
	fmt.Println("Usage: dw claude <subcommand>")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  init              Initialize Claude Code logging infrastructure")
	fmt.Println("  log <event-type>  Log a Claude Code event (reads JSON from stdin)")
	fmt.Println("  auto-summary      Auto-trigger session summary (called by SessionEnd hook)")
	fmt.Println("  auto-summary-exec Internal: Execute summary in background (do not call directly)")
	fmt.Println()
}

// handleAutoSummary handles auto-triggered session summaries on SessionEnd
// This is called by the SessionEnd hook and spawns a background process to do the actual work
// Returns immediately so Claude Code is not blocked
func handleAutoSummary(args []string) {
	// Read hook input from stdin
	stdinData, err := io.ReadAll(os.Stdin)
	if err != nil {
		// Fail silently - don't disrupt Claude Code
		return
	}

	// Try to parse as hook input to extract session ID
	hookInput, err := infra.ParseHookInput(io.NopCloser(bytes.NewReader(stdinData)))
	if err != nil {
		// Not valid hook input, fail silently
		return
	}

	// Get session ID
	sessionID := hookInput.SessionID
	if sessionID == "" {
		// No session ID, can't analyze
		return
	}

	// Load config to check if auto-summary is enabled
	logger := infra.NewDefaultLogger()
	configLoader := infra.NewConfigLoader(logger)
	config, err := configLoader.LoadConfig("")
	if err != nil {
		// Config load failed, fail silently
		return
	}

	// Check if auto-summary is enabled
	if !config.Analysis.AutoSummaryEnabled {
		// Auto-summary disabled, silently exit
		return
	}

	// Spawn detached background process to execute the summary
	// This allows the hook to return immediately while analysis runs in background
	if err := spawnBackgroundSummary(sessionID); err != nil {
		// Fail silently - don't disrupt Claude Code
		return
	}

	// Return immediately - background process will handle the analysis
}

// spawnBackgroundSummary spawns a detached background process to execute the summary
func spawnBackgroundSummary(sessionID string) error {
	// Get the path to the current executable
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// Create command: dw claude auto-summary-exec <session-id>
	cmd := exec.Command(executable, "claude", "auto-summary-exec", sessionID)

	// Detach from parent process
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	// Start the process without waiting for it to complete
	if err := cmd.Start(); err != nil {
		return err
	}

	// Don't wait for the process - let it run in background
	// The process will continue even after the parent exits
	return nil
}

// handleAutoSummaryExec executes the actual summary analysis
// This is called by the background process spawned by handleAutoSummary
func handleAutoSummaryExec(args []string) {
	if len(args) < 1 {
		// No session ID provided
		return
	}

	sessionID := args[0]

	// Load config
	logger := infra.NewDefaultLogger()
	configLoader := infra.NewConfigLoader(logger)
	config, err := configLoader.LoadConfig("")
	if err != nil {
		// Config load failed, exit silently
		return
	}

	// Get the prompt name from config
	promptName := config.Analysis.AutoSummaryPrompt
	if promptName == "" {
		promptName = "session_summary"
	}

	// Create repository and services
	repo, err := infra.NewSQLiteEventRepository(app.DefaultDBPath)
	if err != nil {
		return
	}
	defer repo.Close()

	logsService := app.NewLogsService(repo, repo)
	llmExecutor := app.NewClaudeCLIExecutorWithConfig(logger, config)
	analysisService := app.NewAnalysisService(repo, repo, logsService, llmExecutor, logger, config)

	// Execute the analysis
	_, _ = analysisService.AnalyzeSessionWithPrompt(context.Background(), sessionID, promptName)
	// Ignore errors - this is best-effort background analysis
}

func handleInit(args []string) {
	dbPath := app.DefaultDBPath

	fmt.Println("Initializing Claude Code logging for DarwinFlow...")
	fmt.Println()

	// Ensure database directory exists before creating repository
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating database directory: %v\n", err)
		os.Exit(1)
	}

	// Create infrastructure dependencies
	repository, err := infra.NewSQLiteEventRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating repository: %v\n", err)
		os.Exit(1)
	}
	defer repository.Close()

	hookConfigManager, err := infra.NewHookConfigManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating hook config manager: %v\n", err)
		os.Exit(1)
	}

	// Create application service
	setupService := app.NewSetupService(repository, hookConfigManager)

	// Initialize logging infrastructure
	ctx := context.Background()
	if err := setupService.Initialize(ctx, dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Created logging database:", dbPath)
	fmt.Println("✓ Added hooks to Claude Code settings:", setupService.GetSettingsPath())
	fmt.Println()
	fmt.Println("DarwinFlow logging is now active for all Claude Code sessions.")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Restart Claude Code to activate the hooks")
	fmt.Println("  2. Events will be automatically logged to", dbPath)
}

func handleLog(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: event type required")
		fmt.Fprintln(os.Stderr, "Usage: dw claude log <event-type>")
		os.Exit(1)
	}

	eventTypeStr := args[0]

	// Get max param length from environment or use default
	maxParamLength := 30
	if envVal := os.Getenv("DW_MAX_PARAM_LENGTH"); envVal != "" {
		if parsed, err := fmt.Sscanf(envVal, "%d", &maxParamLength); err == nil && parsed == 1 {
			// Successfully parsed
		}
	}

	// Silently execute (errors shouldn't disrupt Claude Code)
	if err := logFromStdin(eventTypeStr, maxParamLength); err != nil {
		// Silently fail - don't disrupt Claude Code
		return
	}
}

func logFromStdin(eventTypeStr string, maxParamLength int) error {
	// Read hook input from stdin
	stdinData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	// Try to parse as hook input
	hookInput, err := infra.ParseHookInput(io.NopCloser(bytes.NewReader(stdinData)))
	if err != nil {
		// Not valid hook input, fail silently
		return nil
	}

	// Map event type string to domain event type
	eventMapper := &app.EventMapper{}
	eventType := eventMapper.MapEventType(eventTypeStr)

	// Create infrastructure dependencies
	repository, err := infra.NewSQLiteEventRepository(app.DefaultDBPath)
	if err != nil {
		return err
	}
	defer repository.Close()

	transcriptParser := infra.NewTranscriptParser()
	contextDetector := infra.NewContextDetector()

	// Create application service
	loggerService := app.NewLoggerService(
		repository,
		transcriptParser,
		contextDetector,
		infra.NormalizeContent,
	)
	defer loggerService.Close()

	// Convert infra.HookInput to app.HookInputData
	hookInputData := app.HookInputData{
		SessionID:      hookInput.SessionID,
		TranscriptPath: hookInput.TranscriptPath,
		CWD:            hookInput.CWD,
		PermissionMode: hookInput.PermissionMode,
		HookEventName:  hookInput.HookEventName,
		ToolName:       hookInput.ToolName,
		ToolInput:      hookInput.ToolInput,
		ToolOutput:     hookInput.ToolOutput,
		Error:          hookInput.Error,
		UserMessage:    hookInput.UserMessage,
		Prompt:         hookInput.Prompt,
	}

	// Log event
	ctx := context.Background()
	return loggerService.LogFromHookInput(ctx, hookInputData, eventType, maxParamLength)
}
