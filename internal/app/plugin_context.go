package app

import (
	"io"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ProjectContext provides access to app-layer services for plugin tools.
// This is the internal context used by the tool registry (not SDK).
// Renamed from PluginContext to avoid confusion with SDK PluginContext.
type ProjectContext struct {
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

// pluginContextAdapter adapts internal services to SDK PluginContext interface.
// This allows plugins to access system capabilities without depending on internal types.
type pluginContextAdapter struct {
	logger     Logger
	dbPath     string
	workingDir string
}

// NewPluginContext creates a new plugin context adapter
func NewPluginContext(logger Logger, dbPath, workingDir string) pluginsdk.PluginContext {
	return &pluginContextAdapter{
		logger:     logger,
		dbPath:     dbPath,
		workingDir: workingDir,
	}
}

func (p *pluginContextAdapter) GetLogger() pluginsdk.Logger {
	return &loggerAdapter{inner: p.logger}
}

func (p *pluginContextAdapter) GetDBPath() string {
	return p.dbPath
}

func (p *pluginContextAdapter) GetWorkingDir() string {
	return p.workingDir
}

// loggerAdapter adapts app.Logger to pluginsdk.Logger
type loggerAdapter struct {
	inner Logger
}

func (l *loggerAdapter) Debug(format string, args ...interface{}) {
	l.inner.Debug(format, args...)
}

func (l *loggerAdapter) Info(format string, args ...interface{}) {
	l.inner.Info(format, args...)
}

func (l *loggerAdapter) Warn(format string, args ...interface{}) {
	l.inner.Warn(format, args...)
}

func (l *loggerAdapter) Error(format string, args ...interface{}) {
	l.inner.Error(format, args...)
}

// commandContextAdapter adapts internal services to SDK CommandContext interface
type commandContextAdapter struct {
	pluginContextAdapter
	output io.Writer
	input  io.Reader
}

// NewCommandContext creates a new command context adapter
func NewCommandContext(logger Logger, dbPath, workingDir string, output io.Writer, input io.Reader) pluginsdk.CommandContext {
	return &commandContextAdapter{
		pluginContextAdapter: pluginContextAdapter{
			logger:     logger,
			dbPath:     dbPath,
			workingDir: workingDir,
		},
		output: output,
		input:  input,
	}
}

func (c *commandContextAdapter) GetOutput() io.Writer {
	return c.output
}

func (c *commandContextAdapter) GetInput() io.Reader {
	return c.input
}

// toolContextAdapter adapts internal services to SDK ToolContext interface
type toolContextAdapter struct {
	pluginContextAdapter
	output io.Writer
}

// NewToolContext creates a new tool context adapter
func NewToolContext(logger Logger, dbPath, workingDir string, output io.Writer) pluginsdk.ToolContext {
	return &toolContextAdapter{
		pluginContextAdapter: pluginContextAdapter{
			logger:     logger,
			dbPath:     dbPath,
			workingDir: workingDir,
		},
		output: output,
	}
}

func (t *toolContextAdapter) GetOutput() io.Writer {
	return t.output
}
