package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgatilin/darwinflow-pub/internal/app"
	"github.com/kgatilin/darwinflow-pub/internal/domain"
	"github.com/kgatilin/darwinflow-pub/internal/infra"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/claude_code"
)

// AppServices contains all app-layer services needed by commands
type AppServices struct {
	LogsService     *app.LogsService
	AnalysisService *app.AnalysisService
	SetupService    *app.SetupService
	ConfigLoader    app.ConfigLoader
	Logger          app.Logger
	Config          *domain.Config
	DBPath          string
}

// InitializeApp creates all infrastructure and app services
func InitializeApp(dbPath, configPath string, debugMode bool) (*AppServices, error) {
	// 1. Create logger
	var logger *infra.Logger
	if debugMode {
		logger = infra.NewDebugLogger()
	} else {
		logger = infra.NewDefaultLogger()
	}

	// 2. Ensure database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// 3. Create repository
	repo, err := infra.NewSQLiteEventRepository(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// 4. Load config
	configLoader := infra.NewConfigLoader(logger)
	config, err := configLoader.LoadConfig(configPath)
	if err != nil {
		// Non-fatal - use default config
		logger.Warn("Failed to load config, using defaults: %v", err)
		config = domain.DefaultConfig()
	}

	// 5. Create app services
	logsService := app.NewLogsService(repo, repo)
	llmExecutor := app.NewClaudeCLIExecutorWithConfig(logger, config)
	analysisService := app.NewAnalysisService(repo, repo, logsService, llmExecutor, logger, config)

	// 6. Create setup service (for init command)
	hookConfigManager, err := infra.NewHookConfigManager()
	if err != nil {
		logger.Warn("Failed to create hook config manager: %v", err)
		// Non-fatal - setupService will be nil
	}
	var setupService *app.SetupService
	if hookConfigManager != nil {
		setupService = app.NewSetupService(repo, hookConfigManager)
	}

	return &AppServices{
		LogsService:     logsService,
		AnalysisService: analysisService,
		SetupService:    setupService,
		ConfigLoader:    configLoader,
		Logger:          logger,
		Config:          config,
		DBPath:          dbPath,
	}, nil
}

// RegisterPlugins registers all built-in plugins
func RegisterPlugins(registry *app.PluginRegistry, services *AppServices) error {
	// Create claude command handler for plugin commands
	var handler *app.ClaudeCommandHandler

	// Only create handler if we have all required services
	if services.SetupService != nil {
		transcriptParser := infra.NewTranscriptParser()
		contextDetector := infra.NewContextDetector()
		hookInputParser := infra.NewHookInputParserAdapter()
		eventMapper := &app.EventMapper{}

		loggerService := app.NewLoggerService(
			// Repository will be accessed through services
			nil, // Will be set up per-command as needed
			transcriptParser,
			contextDetector,
			infra.NormalizeContent,
		)

		handler = app.NewClaudeCommandHandler(
			services.SetupService,
			loggerService,
			services.AnalysisService,
			hookInputParser,
			eventMapper,
			services.ConfigLoader,
			services.Logger,
			os.Stdout,
		)
	}

	// Register claude_code plugin
	claudePlugin := claude_code.NewClaudeCodePlugin(
		services.AnalysisService,
		services.LogsService,
		services.Logger,
		services.SetupService,
		handler,
		services.DBPath,
	)

	if err := registry.RegisterPlugin(claudePlugin); err != nil {
		return fmt.Errorf("failed to register claude-code plugin: %w", err)
	}

	// Future: Load external plugins from ~/.darwinflow/plugins/

	return nil
}
