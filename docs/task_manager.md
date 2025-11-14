# Package: task_manager_test

**Path**: `pkg/plugins/task_manager`

## Overview

- **Files**: 29
- **Exports**: 244

## Dependencies

This package imports:

**Local packages**:
- `pkg/plugins/task_manager/application`
- `pkg/plugins/task_manager/infrastructure/persistence`
- `pkg/plugins/task_manager/domain/services`
- `pkg/plugins/task_manager`
- `pkg/plugins/task_manager/domain/entities`
- `pkg/pluginsdk`
- `pkg/plugins/task_manager/presentation/cli`
- `pkg/plugins/task_manager/domain`

**External packages**:
- `path/filepath`
- `gopkg.in/yaml.v3`
- `bytes`
- `database/sql`
- `io`
- `github.com/charmbracelet/glamour`
- `context`
- `sort`
- `github.com/charmbracelet/lipgloss`
- `time`
- `fmt`
- `os`
- `encoding/json`
- `github.com/mattn/go-sqlite3`
- `github.com/fsnotify/fsnotify`
- `strings`
- `regexp`
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/bubbles/textarea`
- `github.com/charmbracelet/bubbles/viewport`
- `sync`
- `testing`
- `errors`

## Exported API

### Types

#### ACsLoadedMsg

```go
ACsLoadedMsg
```

**Properties**:

- ACs []*entities.AcceptanceCriteriaEntity
- Error error

#### ADRConfig

```go
ADRConfig
```

**Properties**:

- Required bool
- EnforceOnTaskCompletion bool

#### ADRsLoadedMsg

```go
ADRsLoadedMsg
```

**Properties**:

- ADRs []*entities.ADREntity
- Error error

#### AppModel

```go
AppModel
```

**Methods**:

- `(*AppModel) GetACStartLines() map[string]int`
- `(*AppModel) GetCurrentView() ViewMode`
- `(*AppModel) GetIterationDetailViewportYOffset() int`
- `(*AppModel) GetSelectedIterationACIdx() int`
- `(*AppModel) GetSelectedIterationIdx() int`
- `(*AppModel) GetSelectedIterationTaskIdx() int`
- `(*AppModel) GetTaskStartLines() map[string]int`
- `(*AppModel) Init() tea.Cmd`
- `(*AppModel) RenderProgressBar(float64, int) string`
- `(*AppModel) SetACs([]*entities.AcceptanceCriteriaEntity)`
- `(*AppModel) SetCurrentIteration(*entities.IterationEntity)`
- `(*AppModel) SetCurrentTrack(*entities.TrackEntity)`
- `(*AppModel) SetCurrentView(ViewMode)`
- `(*AppModel) SetDimensions(int)`
- `(*AppModel) SetError(error)`
- `(*AppModel) SetIterationDetailFocusAC(bool)`
- `(*AppModel) SetIterationTasks([]*entities.TaskEntity)`
- `(*AppModel) SetIterations([]*entities.IterationEntity)`
- `(*AppModel) SetRoadmap(*entities.RoadmapEntity)`
- `(*AppModel) SetSelectedACIdx(int)`
- `(*AppModel) SetSelectedIterationIdx(int)`
- `(*AppModel) SetTasks([]*entities.TaskEntity)`
- `(*AppModel) SetTracks([]*entities.TrackEntity)`
- `(*AppModel) Update(tea.Msg) (tea.Model, tea.Cmd)`
- `(*AppModel) View() string`

#### BackMsg

```go
BackMsg
```

#### BackupCommand

```go
BackupCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*BackupCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*BackupCommand) GetDescription() string`
- `(*BackupCommand) GetHelp() string`
- `(*BackupCommand) GetName() string`
- `(*BackupCommand) GetUsage() string`

#### BackupListCommand

```go
BackupListCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*BackupListCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*BackupListCommand) GetDescription() string`
- `(*BackupListCommand) GetHelp() string`
- `(*BackupListCommand) GetName() string`
- `(*BackupListCommand) GetUsage() string`

#### Config

```go
Config
```

**Properties**:

- ADR ADRConfig

#### ContextualHints

```go
ContextualHints
```

**Methods**:

- `(*ContextualHints) Add(string)`
- `(*ContextualHints) Output(io.Writer)`

#### CreateCommand

```go
CreateCommand
```

**Methods**:

- `(*CreateCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*CreateCommand) GetDescription() string`
- `(*CreateCommand) GetHelp() string`
- `(*CreateCommand) GetName() string`
- `(*CreateCommand) GetUsage() string`

#### ErrorMsg

```go
ErrorMsg
```

**Properties**:

- Error error

#### FileWatcher

```go
FileWatcher
```

**Methods**:

- `(*FileWatcher) Start(context.Context, chan pluginsdk.Event) error`
- `(*FileWatcher) Stop() error`

#### FullRoadmapDataLoadedMsg

```go
FullRoadmapDataLoadedMsg
```

**Properties**:

- IterationTasks map[int][]*entities.TaskEntity
- TrackTasks map[string][]*entities.TaskEntity
- BacklogTasks []*entities.TaskEntity
- Error error

#### Hotkey

```go
Hotkey
```

**Properties**:

- Keys []string
- Description string

#### HotkeyGroup

```go
HotkeyGroup
```

**Properties**:

- Name string
- Hotkeys []Hotkey

#### InitCommand

```go
InitCommand
```

**Methods**:

- `(*InitCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*InitCommand) GetDescription() string`
- `(*InitCommand) GetHelp() string`
- `(*InitCommand) GetName() string`
- `(*InitCommand) GetUsage() string`

#### ItemSelectionType

```go
ItemSelectionType
```

#### IterationDetailLoadedMsg

```go
IterationDetailLoadedMsg
```

**Properties**:

- Iteration *entities.IterationEntity
- Tasks []*entities.TaskEntity
- ACs []*entities.AcceptanceCriteriaEntity
- Error error

#### IterationsLoadedMsg

```go
IterationsLoadedMsg
```

**Properties**:

- Iterations []*entities.IterationEntity
- Error error

#### ListCommand

```go
ListCommand
```

**Methods**:

- `(*ListCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*ListCommand) GetDescription() string`
- `(*ListCommand) GetHelp() string`
- `(*ListCommand) GetName() string`
- `(*ListCommand) GetUsage() string`

#### MigrateIDsCommand

```go
MigrateIDsCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*MigrateIDsCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*MigrateIDsCommand) GetDescription() string`
- `(*MigrateIDsCommand) GetHelp() string`
- `(*MigrateIDsCommand) GetName() string`
- `(*MigrateIDsCommand) GetUsage() string`

#### OutputFormat

```go
OutputFormat
```

#### OutputFormatter

```go
OutputFormatter
```

**Methods**:

- `(*OutputFormatter) OutputJSON(interface{}) error`

#### ProjectCreateCommand

```go
ProjectCreateCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*ProjectCreateCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*ProjectCreateCommand) GetDescription() string`
- `(*ProjectCreateCommand) GetHelp() string`
- `(*ProjectCreateCommand) GetName() string`
- `(*ProjectCreateCommand) GetUsage() string`

#### ProjectDeleteCommand

```go
ProjectDeleteCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*ProjectDeleteCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*ProjectDeleteCommand) GetDescription() string`
- `(*ProjectDeleteCommand) GetHelp() string`
- `(*ProjectDeleteCommand) GetName() string`
- `(*ProjectDeleteCommand) GetUsage() string`

#### ProjectListCommand

```go
ProjectListCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*ProjectListCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*ProjectListCommand) GetDescription() string`
- `(*ProjectListCommand) GetHelp() string`
- `(*ProjectListCommand) GetName() string`
- `(*ProjectListCommand) GetUsage() string`

#### ProjectShowCommand

```go
ProjectShowCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*ProjectShowCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*ProjectShowCommand) GetDescription() string`
- `(*ProjectShowCommand) GetHelp() string`
- `(*ProjectShowCommand) GetName() string`
- `(*ProjectShowCommand) GetUsage() string`

#### ProjectSwitchCommand

```go
ProjectSwitchCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*ProjectSwitchCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*ProjectSwitchCommand) GetDescription() string`
- `(*ProjectSwitchCommand) GetHelp() string`
- `(*ProjectSwitchCommand) GetName() string`
- `(*ProjectSwitchCommand) GetUsage() string`

#### PromptCommand

```go
PromptCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*PromptCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*PromptCommand) GetDescription() string`
- `(*PromptCommand) GetHelp() string`
- `(*PromptCommand) GetName() string`
- `(*PromptCommand) GetUsage() string`

#### RestoreCommand

```go
RestoreCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*RestoreCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*RestoreCommand) GetDescription() string`
- `(*RestoreCommand) GetHelp() string`
- `(*RestoreCommand) GetName() string`
- `(*RestoreCommand) GetUsage() string`

#### RoadmapFullCommand

```go
RoadmapFullCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*RoadmapFullCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*RoadmapFullCommand) GetDescription() string`
- `(*RoadmapFullCommand) GetHelp() string`
- `(*RoadmapFullCommand) GetName() string`
- `(*RoadmapFullCommand) GetUsage() string`

#### RoadmapInitCommand

```go
RoadmapInitCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*RoadmapInitCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*RoadmapInitCommand) GetDescription() string`
- `(*RoadmapInitCommand) GetHelp() string`
- `(*RoadmapInitCommand) GetName() string`
- `(*RoadmapInitCommand) GetUsage() string`

#### RoadmapLoadedMsg

```go
RoadmapLoadedMsg
```

**Properties**:

- Roadmap *entities.RoadmapEntity
- Tracks []*entities.TrackEntity
- Iterations []*entities.IterationEntity
- BacklogTasks []*entities.TaskEntity
- Error error

#### RoadmapShowCommand

```go
RoadmapShowCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*RoadmapShowCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*RoadmapShowCommand) GetDescription() string`
- `(*RoadmapShowCommand) GetHelp() string`
- `(*RoadmapShowCommand) GetName() string`
- `(*RoadmapShowCommand) GetUsage() string`

#### RoadmapUpdateCommand

```go
RoadmapUpdateCommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*RoadmapUpdateCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*RoadmapUpdateCommand) GetDescription() string`
- `(*RoadmapUpdateCommand) GetHelp() string`
- `(*RoadmapUpdateCommand) GetName() string`
- `(*RoadmapUpdateCommand) GetUsage() string`

#### TUICommand

```go
TUICommand
```

**Properties**:

- Plugin *TaskManagerPlugin

**Methods**:

- `(*TUICommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*TUICommand) GetDescription() string`
- `(*TUICommand) GetHelp() string`
- `(*TUICommand) GetName() string`
- `(*TUICommand) GetUsage() string`

#### TUIScaffold

```go
TUIScaffold
```

#### TaskDetailLoadedMsg

```go
TaskDetailLoadedMsg
```

**Properties**:

- Task *entities.TaskEntity
- Error error

#### TaskManagerPlugin

```go
TaskManagerPlugin
```

**Methods**:

- `(*TaskManagerPlugin) GetCapabilities() []string`
- `(*TaskManagerPlugin) GetCommands() []pluginsdk.Command`
- `(*TaskManagerPlugin) GetConfig() *Config`
- `(*TaskManagerPlugin) GetEntity(context.Context, string) (pluginsdk.IExtensible, error)`
- `(*TaskManagerPlugin) GetEntityTypes() []pluginsdk.EntityTypeInfo`
- `(*TaskManagerPlugin) GetInfo() pluginsdk.PluginInfo`
- `(*TaskManagerPlugin) GetRepository() domain.RoadmapRepository`
- `(*TaskManagerPlugin) GetRepositoryForProject(string) (domain.RoadmapRepository, func(...), error)`
- `(*TaskManagerPlugin) Query(context.Context, pluginsdk.EntityQuery) ([]pluginsdk.IExtensible, error)`
- `(*TaskManagerPlugin) StartEventStream(context.Context, chan pluginsdk.Event) error`
- `(*TaskManagerPlugin) StopEventStream() error`
- `(*TaskManagerPlugin) UpdateEntity(context.Context, string, map[string]interface{}) (pluginsdk.IExtensible, error)`

#### TrackDetailLoadedMsg

```go
TrackDetailLoadedMsg
```

**Properties**:

- Track *entities.TrackEntity
- Tasks []*entities.TaskEntity
- Error error

#### UpdateCommand

```go
UpdateCommand
```

**Methods**:

- `(*UpdateCommand) Execute(context.Context, pluginsdk.CommandContext, []string) error`
- `(*UpdateCommand) GetDescription() string`
- `(*UpdateCommand) GetHelp() string`
- `(*UpdateCommand) GetName() string`
- `(*UpdateCommand) GetUsage() string`

#### ViewMode

```go
ViewMode
```

### Functions

- `DefaultConfig() *Config`
- `GetCurrentTime() time.Time`
- `GetSystemPrompt(context.Context) string`
- `InitSchema(*sql.DB) error`
- `LoadConfig(string) (*Config, error)`
- `MigrateFromFileStorage(*sql.DB, string) error`
- `NewAppModel(context.Context, domain.RoadmapRepository, pluginsdk.Logger) *AppModel`
- `NewAppModelWithProject(context.Context, domain.RoadmapRepository, pluginsdk.Logger, string) *AppModel`
- `NewContextualHints() *ContextualHints`
- `NewFileWatcher(pluginsdk.Logger, string) (*FileWatcher, error)`
- `NewOutputFormatter(io.Writer, OutputFormat) *OutputFormatter`
- `NewTaskManagerPlugin(pluginsdk.Logger, string, interface{}) (*TaskManagerPlugin, error)`
- `NewTaskManagerPluginWithDatabase(pluginsdk.Logger, string, *sql.DB, interface{}) (*TaskManagerPlugin, error)`
- `ParseOutputFormat(string) (OutputFormat, error)`
- `SaveConfig(string, *Config) error`
- `StartTUI(context.Context, domain.RoadmapRepository, pluginsdk.Logger, io.Writer) error`
- `WrapText(string, int) string`

### Constants

- `DefaultSystemPrompt`
- `EventACAutomaticallyVerified`
- `EventACCreated`
- `EventACDeleted`
- `EventACFailed`
- `EventACPendingReview`
- `EventACUpdated`
- `EventACVerified`
- `EventADRCreated`
- `EventADRDeprecated`
- `EventADRSuperseded`
- `EventADRUpdated`
- `EventIterationCompleted`
- `EventIterationCreated`
- `EventIterationStarted`
- `EventIterationUpdated`
- `EventRoadmapCreated`
- `EventRoadmapUpdated`
- `EventTaskCompleted`
- `EventTaskCreated`
- `EventTaskDeleted`
- `EventTaskStatusChanged`
- `EventTaskUpdated`
- `EventTrackBlocked`
- `EventTrackCompleted`
- `EventTrackCreated`
- `EventTrackStatusChanged`
- `EventTrackUpdated`
- `FormatJSON`
- `FormatLLM`
- `FormatTable`
- `PluginSourceName`
- `SchemaVersion`
- `SelectBacklog`
- `SelectIterations`
- `SelectTracks`
- `ViewACDetail`
- `ViewACFailInput`
- `ViewACList`
- `ViewADRList`
- `ViewError`
- `ViewIterationDetail`
- `ViewIterationList`
- `ViewLoading`
- `ViewRoadmapList`
- `ViewTaskDetail`
- `ViewTrackDetail`

### Variables

- `RunProgram`

## Files

- `pkg/plugins/task_manager/acceptance_criteria_entity_test.go`
- `pkg/plugins/task_manager/adr_entity_test.go`
- `pkg/plugins/task_manager/command_backup.go`
- `pkg/plugins/task_manager/command_migrate.go`
- `pkg/plugins/task_manager/command_project.go`
- `pkg/plugins/task_manager/command_prompt.go`
- `pkg/plugins/task_manager/command_roadmap.go`
- `pkg/plugins/task_manager/command_tui.go`
- `pkg/plugins/task_manager/commands.go`
- `pkg/plugins/task_manager/config.go`
- `pkg/plugins/task_manager/config_test.go`
- `pkg/plugins/task_manager/events.go`
- `pkg/plugins/task_manager/iteration_entity_test.go`
- `pkg/plugins/task_manager/migration_v3_to_v4.go`
- `pkg/plugins/task_manager/migration_v4_to_v5.go`
- `pkg/plugins/task_manager/migration_v4_to_v5_test.go`
- `pkg/plugins/task_manager/output_formatter.go`
- `pkg/plugins/task_manager/output_formatter_test.go`
- `pkg/plugins/task_manager/plugin.go`
- `pkg/plugins/task_manager/plugin_test.go`
- `pkg/plugins/task_manager/prompt.go`
- `pkg/plugins/task_manager/prompt_test.go`
- `pkg/plugins/task_manager/roadmap_entity_test.go`
- `pkg/plugins/task_manager/schema.go`
- `pkg/plugins/task_manager/sqlite_repository_test.go`
- `pkg/plugins/task_manager/track_entity_test.go`
- `pkg/plugins/task_manager/tui_models.go`
- `pkg/plugins/task_manager/tui_models_test.go`
- `pkg/plugins/task_manager/watcher.go`

---

*Generated by `go-arch-lint -format=package`*

