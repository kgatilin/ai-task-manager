package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

// NoOpLogger is a logger that does nothing (for backward compatibility)
type NoOpLogger struct{}

func (l *NoOpLogger) Debug(format string, args ...interface{}) {}
func (l *NoOpLogger) Info(format string, args ...interface{})  {}
func (l *NoOpLogger) Warn(format string, args ...interface{})  {}
func (l *NoOpLogger) Error(format string, args ...interface{}) {}

// LLMExecutor defines the interface for executing LLM queries
type LLMExecutor interface {
	// Execute runs an LLM query with the given prompt and returns the response
	Execute(ctx context.Context, prompt string) (string, error)
}

// Logger interface for dependency injection
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// AnalysisService handles session analysis operations
type AnalysisService struct {
	eventRepo    domain.EventRepository
	analysisRepo domain.AnalysisRepository
	logsService  *LogsService
	llmExecutor  LLMExecutor
	logger       Logger
	config       *domain.Config
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(
	eventRepo domain.EventRepository,
	analysisRepo domain.AnalysisRepository,
	logsService *LogsService,
	llmExecutor LLMExecutor,
	logger Logger,
	config *domain.Config,
) *AnalysisService {
	if config == nil {
		config = domain.DefaultConfig()
	}
	return &AnalysisService{
		eventRepo:    eventRepo,
		analysisRepo: analysisRepo,
		logsService:  logsService,
		llmExecutor:  llmExecutor,
		logger:       logger,
		config:       config,
	}
}

// AnalyzeSession analyzes a specific session with the default analysis prompt
// This is kept for backward compatibility - uses "tool_analysis" prompt
func (s *AnalysisService) AnalyzeSession(ctx context.Context, sessionID string) (*domain.SessionAnalysis, error) {
	return s.AnalyzeSessionWithPrompt(ctx, sessionID, "tool_analysis")
}

// AnalyzeSessionWithPrompt analyzes a specific session with a named prompt from config
func (s *AnalysisService) AnalyzeSessionWithPrompt(ctx context.Context, sessionID, promptName string) (*domain.SessionAnalysis, error) {
	// Get session logs
	s.logger.Debug("Fetching logs for session %s", sessionID)
	logs, err := s.logsService.ListRecentLogs(ctx, 0, 0, sessionID, true)
	if err != nil {
		s.logger.Error("Failed to get session logs: %v", err)
		return nil, fmt.Errorf("failed to get session logs: %w", err)
	}

	if len(logs) == 0 {
		s.logger.Warn("No logs found for session %s", sessionID)
		return nil, fmt.Errorf("no logs found for session %s", sessionID)
	}
	s.logger.Debug("Found %d log records for session %s", len(logs), sessionID)

	// Format logs as markdown
	s.logger.Debug("Formatting logs as markdown")
	var buf bytes.Buffer
	if err := FormatLogsAsMarkdown(&buf, logs); err != nil {
		s.logger.Error("Failed to format logs: %v", err)
		return nil, fmt.Errorf("failed to format logs: %w", err)
	}

	// Get analysis prompt from config
	promptTemplate, exists := s.config.Prompts[promptName]
	if !exists || promptTemplate == "" {
		s.logger.Warn("Prompt %s not found in config, using default tool_analysis", promptName)
		promptTemplate = domain.DefaultToolAnalysisPrompt
		promptName = "tool_analysis"
	}

	prompt := promptTemplate + buf.String()
	s.logger.Debug("Generated prompt with %d characters (%d KB)", len(prompt), len(prompt)/1024)

	// Execute LLM analysis
	s.logger.Info("Invoking Claude CLI for %s analysis...", promptName)
	analysisResult, err := s.llmExecutor.Execute(ctx, prompt)
	if err != nil {
		s.logger.Error("Failed to execute LLM analysis: %v", err)
		return nil, fmt.Errorf("failed to execute LLM analysis: %w", err)
	}
	s.logger.Debug("Claude CLI returned %d characters", len(analysisResult))

	// Create and save analysis with type
	s.logger.Debug("Saving analysis to database")
	analysis := domain.NewSessionAnalysisWithType(
		sessionID,
		analysisResult,
		s.config.Analysis.Model,
		promptTemplate,
		promptName, // analysis type matches prompt name
		promptName,
	)

	if err := s.analysisRepo.SaveAnalysis(ctx, analysis); err != nil {
		s.logger.Error("Failed to save analysis: %v", err)
		return nil, fmt.Errorf("failed to save analysis: %w", err)
	}

	s.logger.Info("Analysis completed successfully")
	return analysis, nil
}

// AnalyzeMultipleSessions analyzes multiple sessions with a specific prompt
// Returns a map of sessionID -> analysis, and any errors encountered
func (s *AnalysisService) AnalyzeMultipleSessions(ctx context.Context, sessionIDs []string, promptName string) (map[string]*domain.SessionAnalysis, []error) {
	results := make(map[string]*domain.SessionAnalysis)
	var errors []error

	for _, sessionID := range sessionIDs {
		analysis, err := s.AnalyzeSessionWithPrompt(ctx, sessionID, promptName)
		if err != nil {
			errors = append(errors, fmt.Errorf("session %s: %w", sessionID, err))
			continue
		}
		results[sessionID] = analysis
	}

	return results, errors
}

// GetLastSession returns the ID of the most recent session
func (s *AnalysisService) GetLastSession(ctx context.Context) (string, error) {
	logs, err := s.logsService.ListRecentLogs(ctx, 1, 0, "", false)
	if err != nil {
		return "", fmt.Errorf("failed to get last session: %w", err)
	}

	if len(logs) == 0 {
		return "", fmt.Errorf("no sessions found")
	}

	return logs[0].SessionID, nil
}

// GetUnanalyzedSessions returns all session IDs that haven't been analyzed
func (s *AnalysisService) GetUnanalyzedSessions(ctx context.Context) ([]string, error) {
	return s.analysisRepo.GetUnanalyzedSessionIDs(ctx)
}

// GetAnalysis retrieves the analysis for a session
func (s *AnalysisService) GetAnalysis(ctx context.Context, sessionID string) (*domain.SessionAnalysis, error) {
	return s.analysisRepo.GetAnalysisBySessionID(ctx, sessionID)
}

// GetAllSessionIDs retrieves all session IDs, ordered by most recent first
// If limit > 0, returns only the latest N sessions
func (s *AnalysisService) GetAllSessionIDs(ctx context.Context, limit int) ([]string, error) {
	return s.analysisRepo.GetAllSessionIDs(ctx, limit)
}

// ClaudeCLIExecutor implements LLMExecutor using the claude CLI tool
type ClaudeCLIExecutor struct {
	logger Logger
	config *domain.Config
}

// NewClaudeCLIExecutor creates a new Claude CLI executor
func NewClaudeCLIExecutor(logger Logger) *ClaudeCLIExecutor {
	if logger == nil {
		logger = &NoOpLogger{}
	}
	return &ClaudeCLIExecutor{
		logger: logger,
		config: domain.DefaultConfig(),
	}
}

// NewClaudeCLIExecutorWithConfig creates a new Claude CLI executor with custom config
func NewClaudeCLIExecutorWithConfig(logger Logger, config *domain.Config) *ClaudeCLIExecutor {
	if logger == nil {
		logger = &NoOpLogger{}
	}
	if config == nil {
		config = domain.DefaultConfig()
	}
	return &ClaudeCLIExecutor{
		logger: logger,
		config: config,
	}
}

// Execute runs claude -p with the given prompt
// Streams output to stderr in real-time for progress visibility
// The prompt parameter is treated as the user prompt unless config specifies system prompt mode
func (e *ClaudeCLIExecutor) Execute(ctx context.Context, prompt string) (string, error) {
	return e.ExecuteWithOptions(ctx, prompt, nil)
}

// ExecuteWithOptions runs claude with custom options
// options can override config settings (model, tokenLimit, etc.)
func (e *ClaudeCLIExecutor) ExecuteWithOptions(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	// Build command arguments
	args := []string{"-p"}

	// Apply model from config or options
	model := e.config.Analysis.Model
	if opt, ok := options["model"].(string); ok && opt != "" {
		model = opt
	}
	if model != "" {
		args = append(args, "--model", model)
	}

	// Apply system prompt mode from config
	if e.config.Analysis.ClaudeOptions.SystemPromptMode == "replace" {
		args = append(args, "--system-prompt", prompt)
		// When using system prompt, we need a user prompt too
		// Use empty prompt to just get the system prompt to work
		prompt = "Analyze the session data provided in the system prompt."
	} else if e.config.Analysis.ClaudeOptions.SystemPromptMode == "append" {
		args = append(args, "--append-system-prompt", prompt)
		prompt = "Analyze the session data."
	}

	// Apply allowed tools from config
	if len(e.config.Analysis.ClaudeOptions.AllowedTools) > 0 {
		args = append(args, "--allowed-tools", strings.Join(e.config.Analysis.ClaudeOptions.AllowedTools, " "))
	} else {
		// Empty allowed tools = no tools
		args = append(args, "--allowed-tools", "")
	}

	// Add the user prompt last
	args = append(args, prompt)

	e.logger.Debug("Executing: claude %s", strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "claude", args...)

	var stdout, stderr bytes.Buffer

	// Use MultiWriter to stream to both the buffer and os.Stderr for real-time feedback
	cmd.Stdout = io.MultiWriter(&stdout, os.Stderr)
	cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)

	e.logger.Debug("Running Claude CLI command...")
	if err := cmd.Run(); err != nil {
		e.logger.Error("Claude CLI command failed: %v", err)
		return "", fmt.Errorf("claude command failed: %w, stderr: %s", err, stderr.String())
	}
	e.logger.Debug("Claude CLI command completed successfully")

	return strings.TrimSpace(stdout.String()), nil
}
