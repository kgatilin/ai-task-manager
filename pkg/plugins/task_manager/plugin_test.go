package task_manager_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// MockLogger is a simple logger for testing
type MockLogger struct {
	messages []string
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	m.messages = append(m.messages, "DEBUG: "+msg)
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.messages = append(m.messages, "INFO: "+msg)
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.messages = append(m.messages, "WARN: "+msg)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.messages = append(m.messages, "ERROR: "+msg)
}

// TestNewTaskManagerPlugin tests plugin creation
func TestNewTaskManagerPlugin(t *testing.T) {
	dir := t.TempDir()
	logger := &MockLogger{}

	plugin, err := task_manager.NewTaskManagerPlugin(logger, dir, nil)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	if plugin == nil {
		t.Error("plugin should not be nil")
	}

	info := plugin.GetInfo()
	if info.Name != "task-manager" {
		t.Errorf("expected plugin name 'task-manager', got %q", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", info.Version)
	}
}

// TestGetCapabilities tests that plugin reports correct capabilities
func TestGetCapabilities(t *testing.T) {
	dir := t.TempDir()
	logger := &MockLogger{}

	plugin, err := task_manager.NewTaskManagerPlugin(logger, dir, nil)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	capabilities := plugin.GetCapabilities()
	expected := []string{"IEntityProvider", "ICommandProvider", "IEventEmitter"}

	if len(capabilities) != len(expected) {
		t.Errorf("expected %d capabilities, got %d", len(expected), len(capabilities))
	}

	capMap := make(map[string]bool)
	for _, cap := range capabilities {
		capMap[cap] = true
	}

	for _, exp := range expected {
		if !capMap[exp] {
			t.Errorf("missing capability: %s", exp)
		}
	}
}

// TestGetEntityTypes tests entity type info
func TestGetEntityTypes(t *testing.T) {
	dir := t.TempDir()
	logger := &MockLogger{}

	plugin, err := task_manager.NewTaskManagerPlugin(logger, dir, nil)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	types := plugin.GetEntityTypes()
	if len(types) != 1 {
		t.Errorf("expected 1 entity type, got %d", len(types))
	}

	if types[0].Type != "task" {
		t.Errorf("expected entity type 'task', got %q", types[0].Type)
	}
}

// TestTaskEntityGetters tests TaskEntity field getters
func TestTaskEntityGetters(t *testing.T) {
	now := time.Now().UTC()
	entity, _ := entities.NewTaskEntity(
		"task-123",
		"track-test",
		"Test Task",
		"Test Description",
		"todo",
		200,
		"",
		now,
		now,
	)

	if entity.GetID() != "task-123" {
		t.Errorf("expected ID 'task-123', got %q", entity.GetID())
	}

	if entity.GetType() != "task" {
		t.Errorf("expected type 'task', got %q", entity.GetType())
	}

	if entity.GetStatus() != "todo" {
		t.Errorf("expected status 'todo', got %q", entity.GetStatus())
	}

	fields := entity.GetAllFields()
	if fields["title"] != "Test Task" {
		t.Errorf("expected title 'Test Task', got %q", fields["title"])
	}
}

// TestTaskEntityProgress tests progress calculation
func TestTaskEntityProgress(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		status   string
		expected float64
	}{
		{"todo", 0.0},
		{"in-progress", 0.5},
		{"done", 1.0},
	}

	for _, test := range tests {
		entity, _ := entities.NewTaskEntity("id", "track-test", "title", "desc", test.status, 300, "", now, now)
		progress := entity.GetProgress()
		if progress != test.expected {
			t.Errorf("for status %q, expected progress %.1f, got %.1f", test.status, test.expected, progress)
		}
	}
}

// TestQueryTasks tests the Query method
// TestCreateCommandExecution tests the create command (REMOVED - obsolete file-based command)
// The create command has been removed in favor of database-based CLI commands.

// TestQueryTasks tests the Query method (REMOVED - obsolete file-based storage)
// Query method is deprecated in favor of database-based CLI commands.

// TestListCommand tests the list command (REMOVED - obsolete file-based command)
// The list command now uses database repository, not file-based storage.

// TestUpdateCommand tests the update command (REMOVED - obsolete file-based command)
// The update command now uses database repository, not file-based storage.

// TestInitCommand tests the init command (REMOVED - obsolete file-based command)
// The init command has been removed in favor of database initialization.

// TestUpdateEntity tests the UpdateEntity method (REMOVED - obsolete file-based storage)
// UpdateEntity method is deprecated in favor of database-based CLI commands.

// TestEventStreamStartStop tests event stream start and stop
func TestEventStreamStartStop(t *testing.T) {
	dir := t.TempDir()
	logger := &MockLogger{}

	plugin, err := task_manager.NewTaskManagerPlugin(logger, dir, nil)
	if err != nil {
		t.Fatalf("failed to create plugin: %v", err)
	}

	// Create event channel
	eventChan := make(chan pluginsdk.Event, 100)
	defer close(eventChan)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start event stream
	err = plugin.StartEventStream(ctx, eventChan)
	if err != nil {
		t.Fatalf("failed to start event stream: %v", err)
	}

	// Stop event stream
	err = plugin.StopEventStream()
	if err != nil {
		t.Fatalf("failed to stop event stream: %v", err)
	}
}

// MockCommandContext implements pluginsdk.CommandContext
type MockCommandContext struct {
	workingDir string
	stdout     *bytes.Buffer
	stdin      *bytes.Buffer
}

func (m *MockCommandContext) GetLogger() pluginsdk.Logger {
	return &MockLogger{}
}

func (m *MockCommandContext) GetWorkingDir() string {
	return m.workingDir
}

func (m *MockCommandContext) EmitEvent(ctx context.Context, event pluginsdk.Event) error {
	return nil
}

func (m *MockCommandContext) GetStdout() io.Writer {
	return m.stdout
}

func (m *MockCommandContext) GetStdin() io.Reader {
	return m.stdin
}

// TestCreateCommand is a helper for testing
type TestCreateCommand struct {
}

func (t *TestCreateCommand) GetName() string {
	return "create"
}

func (t *TestCreateCommand) GetDescription() string {
	return "Create task"
}

func (t *TestCreateCommand) GetUsage() string {
	return "create <title>"
}

func (t *TestCreateCommand) GetHelp() string {
	return ""
}

func (t *TestCreateCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Just create the task directory for testing
	dir := cmdCtx.GetWorkingDir()
	tasksDir := filepath.Join(dir, ".darwinflow", "tasks")
	return os.MkdirAll(tasksDir, 0755)
}
