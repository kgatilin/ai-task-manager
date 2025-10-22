# Package: task_manager

**Path**: `pkg/plugins/task_manager`

**Role**: Example plugin demonstrating real-time event streaming with fsnotify

---

## Quick Reference

- **Files**: 6
- **Exports**: 10+
- **Dependencies**: `pkg/pluginsdk`, `github.com/fsnotify/fsnotify`
- **Layer**: Plugin implementation (example)
- **Type**: External plugin (reference implementation)

---

## Purpose

This is an **example plugin** created as Task 4.3 to demonstrate:
- Real-time event streaming using fsnotify
- File-based task storage (JSON)
- Full plugin capability implementation (IEntityProvider, ICommandProvider, IEventEmitter)
- Proper test coverage and linter compliance

## Quick Start

### Using the Plugin

To integrate this plugin into DarwinFlow:

```go
// In cmd/dw/plugin_registration.go
import "github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"

func registerExternalPlugins(registry *internal.PluginRegistry) {
    plugin, err := task_manager.NewTaskManagerPlugin(logger, workingDir)
    if err != nil {
        log.Fatalf("failed to create task-manager plugin: %v", err)
    }
    registry.RegisterPlugin(plugin)
}
```

### Using the CLI

```bash
# Initialize task directory
dw task-manager init

# Create a new task
dw task-manager create "Implement feature" --description "Add new feature" --priority high

# List all tasks
dw task-manager list

# List tasks with specific status
dw task-manager list --status todo

# Update a task status
dw task-manager update task-123456 --status done

# Update multiple fields
dw task-manager update task-123456 --status in-progress --priority high
```

---

## Architecture

### Plugin Implementation (`plugin.go`)

**TaskManagerPlugin** implements:
- `pluginsdk.Plugin` - Base plugin interface
- `pluginsdk.IEntityProvider` - Query tasks by ID or filters
- `pluginsdk.ICommandProvider` - CLI commands (init, create, list, update)
- `pluginsdk.IEventEmitter` - Real-time file change events

**Key Methods**:
- `GetInfo()` - Plugin metadata (name, version, description)
- `GetCapabilities()` - Lists implemented capabilities
- `GetEntityTypes()` - Returns "task" entity type
- `Query(ctx, query)` - Query tasks with filters and pagination
- `GetEntity(ctx, id)` - Get task by ID
- `UpdateEntity(ctx, id, fields)` - Update task fields
- `GetCommands()` - Returns CLI commands
- `StartEventStream(ctx, eventChan)` - Begin file watching
- `StopEventStream()` - Stop file watching

### Entity (`task_entity.go`)

**TaskEntity** implements:
- `pluginsdk.IExtensible` - Required entity interface
- `pluginsdk.ITrackable` - Status and progress tracking

**Fields**:
- `ID` - Unique task identifier
- `Title` - Task title
- `Description` - Task description
- `Status` - "todo", "in-progress", "done"
- `Priority` - "low", "medium", "high"
- `CreatedAt` - Creation timestamp
- `UpdatedAt` - Last update timestamp

**Status → Progress Mapping**:
- "todo" → 0.0
- "in-progress" → 0.5
- "done" → 1.0

### Commands (`commands.go`)

**4 Commands**:

1. **InitCommand** (`init`)
   - Usage: `dw task-manager init`
   - Creates `.darwinflow/tasks/` directory

2. **CreateCommand** (`create`)
   - Usage: `dw task-manager create <title> [--description <desc>] [--priority <priority>]`
   - Creates new task JSON file with unique ID
   - Triggers `task.created` event

3. **ListCommand** (`list`)
   - Usage: `dw task-manager list [--status <status>]`
   - Lists all tasks or filtered by status
   - Displays table: ID, Title, Status, Priority

4. **UpdateCommand** (`update`)
   - Usage: `dw task-manager update <id> --status <status> [--title] [--priority]`
   - Updates task fields
   - Triggers `task.updated` event

### File Watcher (`watcher.go`)

**FileWatcher** monitors `.darwinflow/tasks/` directory for changes:

**Lifecycle**:
1. `NewFileWatcher()` - Create watcher with fsnotify
2. `Start(ctx, eventChan)` - Watch directory and emit events
3. Background goroutine processes file system events
4. `Stop()` - Gracefully stop watching

**Events Emitted**:
- `task.created` - New task file created
- `task.updated` - Task file modified (status change)
- `task.deleted` - Task file deleted

**Event Payload**:
```json
{
  "type": "task.created",
  "source": "task-manager",
  "timestamp": "2025-10-22T10:00:00Z",
  "payload": {
    "id": "task-123",
    "title": "Example task",
    "status": "todo"
  }
}
```

**Implementation Details**:
- Uses fsnotify for cross-platform file watching
- Tracks file state to detect actual changes vs temporary files
- Respects context cancellation for clean shutdown
- Thread-safe with mutex protection
- Graceful timeout on shutdown (5 seconds)

### Event Types (`events.go`)

**Constants**:
- `EventTaskCreated = "task.created"`
- `EventTaskUpdated = "task.updated"`
- `EventTaskDeleted = "task.deleted"`
- `PluginSourceName = "task-manager"`

---

## Storage Format

### File Structure

```
.darwinflow/
  tasks/
    task-1729604400000000000.json
    task-1729604401000000000.json
    ...
```

### Task JSON File

```json
{
  "id": "task-1729604400000000000",
  "title": "Implement feature",
  "description": "Add new feature to system",
  "status": "in-progress",
  "priority": "high",
  "created_at": "2025-10-22T10:00:00Z",
  "updated_at": "2025-10-22T10:15:30Z"
}
```

---

## Testing

### Test Coverage: 59.1%

**Tests**:
- `TestNewTaskManagerPlugin` - Plugin creation
- `TestGetCapabilities` - Capability reporting
- `TestGetEntityTypes` - Entity type metadata
- `TestTaskEntityGetters` - Entity field access
- `TestTaskEntityProgress` - Progress calculation
- `TestQueryTasks` - Entity querying
- `TestCreateCommandExecution` - Create command
- `TestListCommand` - List command with filters
- `TestUpdateCommand` - Update command
- `TestInitCommand` - Init command
- `TestEventStreamStartStop` - Event stream lifecycle
- `TestEventEmissionOnTaskCreation` - Event emission
- `TestUpdateEntity` - Entity updates

### Running Tests

```bash
# Run all tests
go test ./pkg/plugins/task_manager -v

# Check coverage
go test -cover ./pkg/plugins/task_manager

# Generate coverage report
go test -coverprofile=coverage.out ./pkg/plugins/task_manager
go tool cover -html=coverage.out
```

---

## Architecture Principles

### What's Here

✅ **Plugin implementation** - TaskManagerPlugin
✅ **Entity type** - TaskEntity with IExtensible and ITrackable
✅ **Commands** - Init, create, list, update
✅ **Event streaming** - File watcher with fsnotify
✅ **Event types** - Task lifecycle events
✅ **Comprehensive tests** - 59% coverage, 13 test cases
✅ **Only pluginsdk imports** - No internal/* dependencies

### Key Patterns

1. **Plugin Structure**: Minimal, focused on task management
2. **File-Based Storage**: JSON files for simplicity and event-driven design
3. **Real-Time Events**: fsnotify integration for immediate notifications
4. **Clean Separation**: Commands, entity, watcher, and events in separate files
5. **Error Handling**: Graceful degradation and timeout management

### Plugin Rules

- ✅ Implements all required SDK interfaces correctly
- ✅ Only imports `pkg/pluginsdk` (no `internal/*`)
- ✅ Fully self-contained and extractable
- ✅ Thread-safe operations (watcher with mutex)
- ✅ Respects context cancellation
- ✅ Graceful error handling

---

## Performance Characteristics

### Event Streaming
- **Latency**: File changes detected within 100ms (system-dependent)
- **Memory**: Minimal overhead; tracks file state for de-duplication
- **CPU**: Negligible when idle, scales with file system activity

### File I/O
- **Create**: O(1) - Direct JSON write
- **Query**: O(n) - Reads all task files from directory
- **Update**: O(1) - Direct JSON rewrite
- **Delete**: O(1) - File removal

### Scalability
- Works well for hundreds of tasks
- File system limits apply (typically 1000s of files)
- For larger datasets, consider database backend

---

## Future Enhancements

Possible extensions:

1. **Database Backend** - Replace JSON files with SQLite
2. **Task Dependencies** - Add parent/child relationships
3. **Recurring Tasks** - Implement schedules (cron-like)
4. **Tags/Categories** - Organize tasks
5. **Search** - Full-text search across tasks
6. **Analytics** - Track task completion metrics
7. **Notifications** - Alert on status changes
8. **Sync** - Multi-device synchronization

---

## Troubleshooting

### File Watcher Not Detecting Changes

**Cause**: Some file systems don't send fsnotify events reliably (especially over network mounts).

**Solution**: Test with local storage or configure longer polling intervals.

### High CPU Usage

**Cause**: Too many file system events (rapid file creation/deletion).

**Solution**: Add debouncing or batch events in watcher.

### Tasks Not Persisting

**Cause**: Permissions issue on `.darwinflow/tasks/` directory.

**Solution**: Verify directory is writable: `chmod 755 .darwinflow/tasks`

---

## Files

- `plugin.go` - TaskManagerPlugin implementation
- `task_entity.go` - TaskEntity with IExtensible and ITrackable
- `commands.go` - Init, Create, List, Update commands
- `events.go` - Event type constants
- `watcher.go` - File watcher with fsnotify
- `plugin_test.go` - Comprehensive tests (59% coverage)
- `CLAUDE.md` - This file

---

## Example: Integration Guide

To use this plugin as a template for your own plugins:

1. **Copy structure**: Use this plugin's file organization
2. **Replace entity**: Implement your entity type (inherits IExtensible)
3. **Implement commands**: Copy command pattern from commands.go
4. **Add storage**: Extend SaveTask/LoadTask for your backend
5. **Add events**: Emit events for important state changes
6. **Write tests**: Achieve 70%+ coverage minimum

---

*Created as Task 4.3: Example plugin demonstrating real-time event streaming with fsnotify*
