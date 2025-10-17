# Public API

## claude

### Types

- **Logger**
  - Methods:
    - (*Logger) Close() error
    - (*Logger) LogEvent(events.Type, interface{}) error
    - (*Logger) LogFromHookInput(hooks.HookInput, events.Type, int) error
- **SQLiteStore**
  - Methods:
    - (*SQLiteStore) Close() error
    - (*SQLiteStore) Init(context.Context) error
    - (*SQLiteStore) Query(context.Context, storage.Filter) ([]storage.Record, error)
    - (*SQLiteStore) Store(context.Context, storage.Record) error
- *Settings*
  - Properties:
    - Hooks map[string][]hooks.HookMatcher
    - Other map[string]interface{}
- **SettingsManager**
  - Methods:
    - (*SettingsManager) AddDarwinFlowHooks() error
    - (*SettingsManager) GetSettingsPath() string
    - (*SettingsManager) ReadSettings() (*Settings, error)
    - (*SettingsManager) WriteSettings(*Settings) error
- *TranscriptEntry*
  - Properties:
    - Role string
    - Content string
    - Type string
    - Name string
    - Parameters map[string]interface{}
    - Input map[string]interface{}

### Package Functions

- DetectContext() string
- ExtractLastAssistantMessage(string) (string, error)
- ExtractLastToolUse(string, int) (string, string, error)
- ExtractLastUserMessage(string) (string, error)
- InitializeLogging(string) error
- LogFromStdin(string, int) error
- NewLogger(string) (*Logger, error)
- NewSQLiteStore(string) (*SQLiteStore, error)
- NewSettingsManager() (*SettingsManager, error)
- NormalizeContent(string) string
- ParseTranscript(string) ([]TranscriptEntry, error)

### Constants

- DefaultDBPath

## events

### Types

- *ChatPayload*
  - Properties:
    - Message string
    - Context string
- *ContextPayload*
  - Properties:
    - Context string
    - Description string
- *ErrorPayload*
  - Properties:
    - Error string
    - StackTrace string
    - Context string
- *Event*
  - Properties:
    - ID string
    - Timestamp int64
    - Event Type
    - Payload json.RawMessage
    - Content string
- *FilePayload*
  - Properties:
    - FilePath string
    - Changes string
    - DurationMs int64
    - Context string
- *ToolPayload*
  - Properties:
    - Tool string
    - Parameters string
    - Result string
    - DurationMs int64
    - Context string
- *Type*

### Package Functions

- NewEvent(Type, interface{}, string) (*Event, error)

### Constants

- ChatMessageAssistant
- ChatMessageUser
- ChatStarted
- ContextChanged
- Error
- FileRead
- FileWritten
- ToolInvoked
- ToolResult

## hooks

### Types

- *EventType*
- *HookAction*
  - Properties:
    - Type string
    - Command string
    - Timeout int
- *HookConfig*
  - Properties:
    - Hooks map[string][]HookMatcher
- *HookInput*
  - Properties:
    - SessionID string
    - TranscriptPath string
    - CWD string
    - HookEventName string
- *HookMatcher*
  - Properties:
    - Matcher string
    - Hooks []HookAction

### Package Functions

- DefaultConfig() HookConfig
- MergeConfig(HookConfig) HookConfig

### Constants

- Notification
- PostToolUse
- PreCompact
- PreToolUse
- SessionEnd
- SessionStart
- Stop
- SubagentStop
- UserPromptSubmit

## storage

### Types

- *Filter*
  - Properties:
    - StartTime *time.Time
    - EndTime *time.Time
    - EventTypes []string
    - Context string
    - SearchText string
    - Limit int
    - Offset int
- *Record*
  - Properties:
    - ID string
    - Timestamp int64
    - EventType string
    - Payload []byte
    - Content string
- *Store*

