# Dependency Graph

## cmd/dw/claude.go
depends on:
  - local:pkg/claude
    - DefaultDBPath
    - InitializeLogging
    - LogFromStdin
    - NewSettingsManager

## cmd/dw/main.go
depends on:

## internal/events/event.go
depends on:
  - external:github.com/google/uuid
    - New

## internal/hooks/config.go
depends on: (none)

## internal/storage/storage.go
depends on:

## pkg/claude/cli.go
depends on:
  - local:internal/events
    - ChatMessageAssistant
    - ChatMessageUser
    - ChatStarted
    - ContextChanged
    - Error
    - FileRead
    - FileWritten
    - ToolInvoked
    - ToolResult
    - Type
  - local:internal/hooks
    - HookInput

## pkg/claude/context.go
depends on:

## pkg/claude/logger.go
depends on:
  - local:internal/events
    - ChatMessageAssistant
    - ChatMessageUser
    - ChatPayload
    - ChatStarted
    - NewEvent
    - ToolInvoked
    - ToolPayload
    - ToolResult
    - Type
  - local:internal/hooks
    - HookInput
  - local:internal/storage
    - Record
    - Store

## pkg/claude/settings.go
depends on:
  - local:internal/hooks
    - DefaultConfig
    - HookConfig
    - HookMatcher
    - MergeConfig

## pkg/claude/sqlite.go
depends on:
  - local:internal/storage
    - Filter
    - Record
  - external:github.com/mattn/go-sqlite3

## pkg/claude/transcript.go
depends on:

