package app

import (
	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

// PluginContext provides access to app-layer services for plugins
// This lives in app layer (not domain) to avoid domain depending on app
type PluginContext struct {
	// EventRepo provides access to logged events
	EventRepo domain.EventRepository

	// AnalysisRepo provides access to session analyses
	AnalysisRepo domain.AnalysisRepository

	// Config is the project's configuration
	Config *domain.Config

	// CWD is the current working directory
	CWD string

	// DBPath is the path to the database
	DBPath string
}
