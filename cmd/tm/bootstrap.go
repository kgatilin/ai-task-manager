package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
	infralogger "github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/logger"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/infrastructure/persistence"
)

// App contains all dependencies for the tm binary
type App struct {
	Logger           logger.Logger
	ConfigPath       string
	WorkingDir       string
	ActiveProject    string
	RepositoryCommon *persistence.SQLiteRepositoryComposite

	// Domain services (stateless)
	ValidationService      *services.ValidationService
	DomainIterationService *services.IterationService

	// Application services
	TrackService     *application.TrackApplicationService
	TaskService      *application.TaskApplicationService
	IterationService *application.IterationApplicationService
	ADRService       *application.ADRApplicationService
	ACService        *application.ACApplicationService
	RoadmapService   *application.RoadmapApplicationService
	DocumentService  *application.DocumentApplicationService
	ProjectService   *application.ProjectApplicationService
}

// BootstrapApp initializes the application
func BootstrapApp() (*App, error) {
	// Determine config path
	configPath := GetConfigPath()

	// Create simple logger
	logger := infralogger.NewStandardLogger(logger.LevelInfo)

	// Resolve working directory
	workingDir := persistence.ResolveWorkingDirectory()

	// Get active project
	activeProject, err := persistence.GetActiveProject(workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get active project: %w", err)
	}

	// Open database for active project
	_, db, err := persistence.OpenProjectDatabase(workingDir, activeProject)
	if err != nil {
		return nil, fmt.Errorf("failed to open project database: %w", err)
	}

	// Create repository composite (provides all 6 focused repositories)
	// Repository composite owns the DB connection and handles cleanup
	repoComposite := persistence.NewSQLiteRepositoryComposite(db, logger)

	// Create domain services (stateless)
	validationService := services.NewValidationService()
	domainIterationService := services.NewIterationService()

	// Create application services with injected dependencies
	trackService := application.NewTrackApplicationService(
		repoComposite.Track,
		repoComposite.Roadmap,
		repoComposite.Aggregate,
		validationService,
	)

	taskService := application.NewTaskApplicationService(
		repoComposite.Task,
		repoComposite.Track,
		repoComposite.Aggregate,
		repoComposite.AC,
		validationService,
	)

	iterationAppService := application.NewIterationApplicationService(
		repoComposite.Iteration,
		repoComposite.Task,
		repoComposite.Aggregate,
		domainIterationService,
		validationService,
	)

	adrService := application.NewADRApplicationService(
		repoComposite.ADR,
		repoComposite.Track,
		repoComposite.Aggregate,
		validationService,
	)

	acService := application.NewACApplicationService(
		repoComposite.AC,
		repoComposite.Task,
		repoComposite.Aggregate,
		validationService,
	)

	roadmapService := application.NewRoadmapApplicationService(
		repoComposite.Roadmap,
		repoComposite.Track,
		repoComposite.Task,
		repoComposite.Iteration,
		validationService,
	)

	documentService := application.NewDocumentApplicationService(
		repoComposite.Document,
		repoComposite.Track,
		repoComposite.Iteration,
	)

	// Create project management repository and service
	projectMgmtRepo := persistence.NewFileSystemProjectManagementRepository(workingDir)
	projectService := application.NewProjectService(
		projectMgmtRepo,
		validationService,
	)

	// Create app instance with all dependencies
	app := &App{
		Logger:                 logger,
		ConfigPath:             configPath,
		WorkingDir:             workingDir,
		ActiveProject:          activeProject,
		RepositoryCommon:       repoComposite,
		ValidationService:      validationService,
		DomainIterationService: domainIterationService,
		TrackService:           trackService,
		TaskService:            taskService,
		IterationService:       iterationAppService,
		ADRService:             adrService,
		ACService:              acService,
		RoadmapService:         roadmapService,
		DocumentService:        documentService,
		ProjectService:         projectService,
	}

	return app, nil
}

// Close closes database connections and cleanup
func (a *App) Close() error {
	if a.RepositoryCommon != nil {
		return a.RepositoryCommon.Close()
	}
	return nil
}

// getConfigPath returns the config file path (~/.tm/config.yaml)
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./.tm/config.yaml"
	}
	return filepath.Join(homeDir, ".tm", "config.yaml")
}
