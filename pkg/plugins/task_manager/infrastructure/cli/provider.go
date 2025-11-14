package cli

import (
	"database/sql"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
)

// PluginProvider provides access to plugin infrastructure for CLI commands.
// This interface allows infrastructure commands to access plugin resources
// without tight coupling to the plugin implementation.
type PluginProvider interface {
	// GetWorkingDir returns the working directory
	GetWorkingDir() string

	// GetLogger returns the plugin logger
	GetLogger() pluginsdk.Logger

	// GetActiveProject returns the active project name
	GetActiveProject() (string, error)

	// SetActiveProject sets the active project
	SetActiveProject(projectName string) error

	// GetProjectDatabase returns a database connection for the specified project
	GetProjectDatabase(projectName string) (*sql.DB, error)

	// GetRepositoryForProject returns a repository for the specified project
	// Returns the repository, a cleanup function, and an error
	GetRepositoryForProject(projectName string) (domain.RoadmapRepository, func(), error)
}
