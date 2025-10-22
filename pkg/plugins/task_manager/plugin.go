package task_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// Ensure plugin implements required SDK interfaces
var (
	_ pluginsdk.Plugin           = (*TaskManagerPlugin)(nil)
	_ pluginsdk.IEntityProvider  = (*TaskManagerPlugin)(nil)
	_ pluginsdk.ICommandProvider = (*TaskManagerPlugin)(nil)
	_ pluginsdk.IEventEmitter    = (*TaskManagerPlugin)(nil)
)

// TaskManagerPlugin provides task management with real-time file watching.
// It implements Plugin, IEntityProvider, ICommandProvider, and IEventEmitter interfaces.
type TaskManagerPlugin struct {
	logger      pluginsdk.Logger
	workingDir  string
	tasksDir    string
	fileWatcher *FileWatcher
}

// NewTaskManagerPlugin creates a new task manager plugin
func NewTaskManagerPlugin(logger pluginsdk.Logger, workingDir string) (*TaskManagerPlugin, error) {
	tasksDir := filepath.Join(workingDir, ".darwinflow", "tasks")

	fileWatcher, err := NewFileWatcher(logger, tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	return &TaskManagerPlugin{
		logger:      logger,
		workingDir:  workingDir,
		tasksDir:    tasksDir,
		fileWatcher: fileWatcher,
	}, nil
}

// GetInfo returns metadata about this plugin (SDK interface)
func (p *TaskManagerPlugin) GetInfo() pluginsdk.PluginInfo {
	return pluginsdk.PluginInfo{
		Name:        "task-manager",
		Version:     "1.0.0",
		Description: "Task tracking with real-time file watching",
		IsCore:      false,
	}
}

// GetCapabilities returns the capability interfaces this plugin implements (SDK interface)
func (p *TaskManagerPlugin) GetCapabilities() []string {
	return []string{"IEntityProvider", "ICommandProvider", "IEventEmitter"}
}

// GetEntityTypes returns the entity types this plugin provides (SDK interface)
func (p *TaskManagerPlugin) GetEntityTypes() []pluginsdk.EntityTypeInfo {
	return []pluginsdk.EntityTypeInfo{
		{
			Type:              "task",
			DisplayName:       "Task",
			DisplayNamePlural: "Tasks",
			Capabilities:      []string{"IExtensible", "ITrackable"},
			Icon:              "✓",
			Description:       "Task with status tracking",
		},
	}
}

// Query returns entities matching the given query (SDK interface)
func (p *TaskManagerPlugin) Query(ctx context.Context, query pluginsdk.EntityQuery) ([]pluginsdk.IExtensible, error) {
	// Ensure tasks directory exists
	if err := os.MkdirAll(p.tasksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create tasks directory: %w", err)
	}

	// Read all task files
	entries, err := os.ReadDir(p.tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %w", err)
	}

	entities := make([]pluginsdk.IExtensible, 0)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(p.tasksDir, entry.Name())
		task, err := p.loadTaskFromFile(filePath)
		if err != nil {
			p.logger.Warn("failed to load task", "path", filePath, "error", err)
			continue
		}

		// Apply filters if specified
		if !p.matchesFilters(task, query.Filters) {
			continue
		}

		entities = append(entities, task)
	}

	// Apply offset and limit
	if query.Offset > 0 {
		if query.Offset >= len(entities) {
			return []pluginsdk.IExtensible{}, nil
		}
		entities = entities[query.Offset:]
	}

	if query.Limit > 0 && len(entities) > query.Limit {
		entities = entities[:query.Limit]
	}

	return entities, nil
}

// GetEntity retrieves a single entity by ID (SDK interface)
func (p *TaskManagerPlugin) GetEntity(ctx context.Context, entityID string) (pluginsdk.IExtensible, error) {
	filePath := filepath.Join(p.tasksDir, entityID+".json")
	task, err := p.loadTaskFromFile(filePath)
	if err != nil {
		return nil, pluginsdk.ErrNotFound
	}
	return task, nil
}

// UpdateEntity updates an entity's fields (SDK interface)
func (p *TaskManagerPlugin) UpdateEntity(ctx context.Context, entityID string, fields map[string]interface{}) (pluginsdk.IExtensible, error) {
	filePath := filepath.Join(p.tasksDir, entityID+".json")
	task, err := p.loadTaskFromFile(filePath)
	if err != nil {
		return nil, pluginsdk.ErrNotFound
	}

	// Update fields
	if title, ok := fields["title"]; ok {
		if titleStr, ok := title.(string); ok {
			task.Title = titleStr
		}
	}
	if description, ok := fields["description"]; ok {
		if descStr, ok := description.(string); ok {
			task.Description = descStr
		}
	}
	if status, ok := fields["status"]; ok {
		if statusStr, ok := status.(string); ok {
			task.Status = statusStr
		}
	}
	if priority, ok := fields["priority"]; ok {
		if priorityStr, ok := priority.(string); ok {
			task.Priority = priorityStr
		}
	}

	// Save updated task
	if err := p.saveTaskToFile(filePath, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

// GetCommands returns the CLI commands provided by this plugin (SDK interface)
func (p *TaskManagerPlugin) GetCommands() []pluginsdk.Command {
	return []pluginsdk.Command{
		&InitCommand{plugin: p},
		&CreateCommand{plugin: p},
		&ListCommand{plugin: p},
		&UpdateCommand{plugin: p},
	}
}

// StartEventStream begins streaming events to the provided channel (SDK interface)
func (p *TaskManagerPlugin) StartEventStream(ctx context.Context, eventChan chan<- pluginsdk.Event) error {
	p.logger.Info("starting event stream for task-manager plugin")
	return p.fileWatcher.Start(ctx, eventChan)
}

// StopEventStream stops the event stream (SDK interface)
func (p *TaskManagerPlugin) StopEventStream() error {
	p.logger.Info("stopping event stream for task-manager plugin")
	return p.fileWatcher.Stop()
}

// Helper methods

// loadTaskFromFile loads a task from a JSON file
func (p *TaskManagerPlugin) loadTaskFromFile(filePath string) (*TaskEntity, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var task TaskEntity
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// saveTaskToFile saves a task to a JSON file
func (p *TaskManagerPlugin) saveTaskToFile(filePath string, task *TaskEntity) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// matchesFilters checks if an entity matches the given filters
func (p *TaskManagerPlugin) matchesFilters(entity pluginsdk.IExtensible, filters map[string]interface{}) bool {
	if len(filters) == 0 {
		return true
	}

	for key, expectedValue := range filters {
		actualValue := entity.GetField(key)
		if actualValue != expectedValue {
			return false
		}
	}

	return true
}
