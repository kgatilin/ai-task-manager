package task_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// FileWatcher watches the tasks directory for changes and emits events.
// It uses fsnotify to watch for file system changes.
type FileWatcher struct {
	mu           sync.Mutex
	logger       pluginsdk.Logger
	tasksDir     string
	watcher      *fsnotify.Watcher
	ctx          context.Context
	cancel       context.CancelFunc
	eventChan    chan<- pluginsdk.Event
	isRunning    bool
	done         chan struct{}
	trackedFiles map[string]*TaskEntity // Track known file states
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(logger pluginsdk.Logger, tasksDir string) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &FileWatcher{
		logger:       logger,
		tasksDir:     tasksDir,
		watcher:      watcher,
		trackedFiles: make(map[string]*TaskEntity),
		done:         make(chan struct{}),
	}, nil
}

// Start begins watching the tasks directory
func (fw *FileWatcher) Start(ctx context.Context, eventChan chan<- pluginsdk.Event) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.isRunning {
		return fmt.Errorf("file watcher is already running")
	}

	// Create tasks directory if it doesn't exist
	if err := os.MkdirAll(fw.tasksDir, 0755); err != nil {
		return fmt.Errorf("failed to create tasks directory: %w", err)
	}

	// Create a cancellable context
	fw.ctx, fw.cancel = context.WithCancel(ctx)
	fw.eventChan = eventChan

	// Watch the tasks directory
	if err := fw.watcher.Add(fw.tasksDir); err != nil {
		fw.watcher.Close()
		return fmt.Errorf("failed to watch directory: %w", err)
	}

	// Load existing tasks
	if err := fw.loadExistingTasks(); err != nil {
		fw.logger.Warn("failed to load existing tasks: %v", err)
	}

	fw.isRunning = true

	// Start the watch loop in a background goroutine
	go fw.watchLoop()

	fw.logger.Info("file watcher started for directory", "path", fw.tasksDir)
	return nil
}

// Stop stops watching the tasks directory
func (fw *FileWatcher) Stop() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if !fw.isRunning {
		return nil
	}

	fw.isRunning = false

	// Cancel the context to stop the watch loop
	if fw.cancel != nil {
		fw.cancel()
	}

	// Wait for the watch loop to finish
	select {
	case <-fw.done:
		fw.logger.Info("file watcher stopped")
	case <-time.After(5 * time.Second):
		fw.logger.Warn("file watcher did not stop gracefully, forcing close")
	}

	// Close the watcher
	if fw.watcher != nil {
		fw.watcher.Close()
	}

	return nil
}

// watchLoop runs in a background goroutine and processes file system events
func (fw *FileWatcher) watchLoop() {
	defer close(fw.done)

	for {
		select {
		case <-fw.ctx.Done():
			return

		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}
			fw.handleFileEvent(event)

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			fw.logger.Error("file watcher error", "error", err)
		}
	}
}

// handleFileEvent processes a single file system event
func (fw *FileWatcher) handleFileEvent(event fsnotify.Event) {
	// Only process JSON files
	if filepath.Ext(event.Name) != ".json" {
		return
	}

	taskID := fileNameToTaskID(event.Name)
	if taskID == "" {
		return
	}

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		fw.handleTaskCreated(event.Name, taskID)

	case event.Op&fsnotify.Write == fsnotify.Write:
		fw.handleTaskUpdated(event.Name, taskID)

	case event.Op&fsnotify.Remove == fsnotify.Remove:
		fw.handleTaskDeleted(taskID)
	}
}

// handleTaskCreated processes a task creation event
func (fw *FileWatcher) handleTaskCreated(filePath string, taskID string) {
	task, err := fw.loadTaskFromFile(filePath)
	if err != nil {
		fw.logger.Warn("failed to load task after creation", "path", filePath, "error", err)
		return
	}

	fw.mu.Lock()
	fw.trackedFiles[filePath] = task
	fw.mu.Unlock()

	// Emit event
	event := pluginsdk.Event{
		Type:      EventTaskCreated,
		Source:    PluginSourceName,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"id":    task.ID,
			"title": task.Title,
			"status": task.Status,
		},
		Version: "1.0",
	}

	select {
	case fw.eventChan <- event:
		fw.logger.Debug("emitted task.created event", "task_id", task.ID)
	case <-fw.ctx.Done():
		return
	}
}

// handleTaskUpdated processes a task update event
func (fw *FileWatcher) handleTaskUpdated(filePath string, taskID string) {
	task, err := fw.loadTaskFromFile(filePath)
	if err != nil {
		fw.logger.Warn("failed to load task after update", "path", filePath, "error", err)
		return
	}

	fw.mu.Lock()
	oldTask := fw.trackedFiles[filePath]
	fw.trackedFiles[filePath] = task
	fw.mu.Unlock()

	// Only emit event if status changed
	if oldTask == nil || oldTask.Status != task.Status {
		event := pluginsdk.Event{
			Type:      EventTaskUpdated,
			Source:    PluginSourceName,
			Timestamp: time.Now(),
			Payload: map[string]interface{}{
				"id":        task.ID,
				"title":     task.Title,
				"status":    task.Status,
				"old_status": "",
			},
			Version: "1.0",
		}

		if oldTask != nil {
			event.Payload["old_status"] = oldTask.Status
		}

		select {
		case fw.eventChan <- event:
			fw.logger.Debug("emitted task.updated event", "task_id", task.ID)
		case <-fw.ctx.Done():
			return
		}
	}
}

// handleTaskDeleted processes a task deletion event
func (fw *FileWatcher) handleTaskDeleted(taskID string) {
	fw.mu.Lock()
	delete(fw.trackedFiles, taskIDToFileName(fw.tasksDir, taskID))
	fw.mu.Unlock()

	event := pluginsdk.Event{
		Type:      EventTaskDeleted,
		Source:    PluginSourceName,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"id": taskID,
		},
		Version: "1.0",
	}

	select {
	case fw.eventChan <- event:
		fw.logger.Debug("emitted task.deleted event", "task_id", taskID)
	case <-fw.ctx.Done():
		return
	}
}

// loadExistingTasks loads all existing task files
func (fw *FileWatcher) loadExistingTasks() error {
	entries, err := os.ReadDir(fw.tasksDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(fw.tasksDir, entry.Name())
		task, err := fw.loadTaskFromFile(filePath)
		if err != nil {
			fw.logger.Warn("failed to load existing task", "path", filePath, "error", err)
			continue
		}

		fw.mu.Lock()
		fw.trackedFiles[filePath] = task
		fw.mu.Unlock()
	}

	return nil
}

// loadTaskFromFile loads a task from a JSON file
func (fw *FileWatcher) loadTaskFromFile(filePath string) (*TaskEntity, error) {
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

// fileNameToTaskID extracts task ID from file name
func fileNameToTaskID(filePath string) string {
	fileName := filepath.Base(filePath)
	if len(fileName) > 5 { // Remove .json extension
		return fileName[:len(fileName)-5]
	}
	return ""
}

// taskIDToFileName creates file path from task ID
func taskIDToFileName(tasksDir, taskID string) string {
	return filepath.Join(tasksDir, taskID+".json")
}
