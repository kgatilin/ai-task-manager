# Claude Code Logging Requirements Specification

## Overview

Build a lightweight logging system for Claude Code interactions within the `dw` tool. The system captures all Claude Code events as structured logs using event sourcing principles, stores them in SQLite with vector search capabilities, and enables future pattern detection and workflow optimization.

## Goals

1. **Comprehensive instrumentation**: Capture all Claude Code interactions automatically
2. **Event sourcing foundation**: Treat all interactions as immutable events with structured payloads
3. **Semantic search capability**: Enable future pattern detection through vector embeddings
4. **Lightweight and fast**: Logging should not impact Claude Code performance
5. **Self-modifying potential**: Data structure supports future automation and optimization

## Components

### 1. CLI Commands

#### `dw claude init`
- **Purpose**: Set up Claude Code logging infrastructure
- **Behavior**:
  - Create `.darwinflow/logs/` directory structure if it doesn't exist
  - Initialize SQLite database (`events.db`) with required schema
  - Check if Claude Code hooks file exists
  - If hooks file doesn't exist:
    - Generate new hooks file with all event handlers
    - Configure handlers to call `dw claude log`
  - If hooks file exists:
    - Parse existing hooks
    - Add missing event handlers for logging
    - Preserve existing custom hooks
    - Update handlers to include logging calls
  - Output status: which hooks were added/updated
  - Show hooks file location and next steps

#### `dw claude log`
- **Purpose**: Log a Claude Code event (called by hooks)
- **Behavior**:
  - Accept event type and payload as command arguments
  - Extract current context (project/entity) from environment or cwd
  - Generate normalized content representation from payload
  - Insert event into SQLite database
  - Return immediately
  - Handle logging failures gracefully (never block Claude Code)

### 2. Event Types

All events captured from Claude Code:

- `chat.started` - New conversation initiated
- `chat.message.user` - User message sent
- `chat.message.assistant` - Assistant response received
- `tool.invoked` - Tool call made (filesystem, web_search, etc.)
- `tool.result` - Tool execution completed
- `file.read` - File accessed by Claude
- `file.written` - File modified by Claude
- `context.changed` - User switched project/entity context
- `error` - Error occurred during interaction

### 3. Log Structure

Each log entry must contain:

- **id**: Unique identifier (UUID)
- **timestamp**: Unix timestamp in milliseconds
- **event**: Event type from defined set
- **payload**: JSON object with event-specific data
  - message (for chat events)
  - tool (for tool events)
  - file_path (for file events)
  - context (project/entity identifier)
  - duration_ms (for timed events)
  - error (for error events)
  - Extensible for future event types
- **content**: Normalized text representation for future analysis

### 4. Storage Layer

#### Technology Choice
- **SQLite** for structured event storage
- Single database file: `.darwinflow/logs/events.db`

#### Schema Requirements

**Events Table**:
- Primary key: event id
- Indexed fields: timestamp, event type, context, file_path
- Full-text search on content field (for text-based queries)
- JSON payload stored as text blob

#### Query Patterns to Support
- Temporal filtering: events in date range
- Event type filtering: all file reads, all tool calls
- Context filtering: all events for specific project
- Full-text search: find events by content keywords
- Combined queries: text search within context + time range

### 5. Content Normalization Strategy

The `content` field is used for embedding generation. Normalization rules per event type:

- **chat.message.user**: Full message text
- **chat.message.assistant**: Full response text (or truncated if >8k tokens)
- **tool.invoked**: Tool name + parameters summary
- **tool.result**: Tool name + result summary (not full result)
- **file.read**: File path + first 50 lines or 2k characters
- **file.written**: File path + change description (not full diff)
- **context.changed**: New context identifier + description
- **error**: Error message + stack trace summary

### 6. Context Detection

The logger must determine current context automatically:

**Priority order**:
1. Environment variable `DW_CONTEXT` (explicitly set by dw tool)
2. Parse from current working directory (detect project from path)
3. Check for context markers in recent events
4. Default to "unknown" if unable to determine

Context format: `{type}/{identifier}` (e.g., `project/darwinflow`, `entity/user-auth-epic`)

### 7. Claude Code Hooks File Management

**Location**: Determined by Claude Code configuration (typically `~/.config/claude/hooks.ts`)

**Hook Functions to Implement**:
- `onChatStarted()` - Log when new chat begins
- `onUserMessage(message)` - Log user input
- `onAssistantMessage(message)` - Log Claude response
- `onToolInvoked(tool, params)` - Log tool calls
- `onToolResult(tool, result)` - Log tool results
- `onFileRead(path)` - Log file access
- `onFileWritten(path, changes)` - Log file modifications
- `onError(error)` - Log errors

**Hook Implementation Pattern**:
Each hook should:
1. Extract relevant data from hook parameters
2. Call `dw claude log <event-type> <json-payload>`
3. Never throw errors (wrap in try-catch)
4. Return immediately (no async waiting)

**New Hooks File Creation**:
When no hooks file exists:
- Generate complete hooks file with all logging handlers
- Use TypeScript template with proper types
- Include inline documentation for each hook
- Add commented examples for custom extensions

**Existing Hooks File Update**:
When hooks file already exists:
- Parse the file to detect existing hook functions
- For each hook:
  - If hook doesn't exist: add new hook function
  - If hook exists but has no logging: append logging call at end of function
  - If hook exists with logging: update logging call to match current format
  - Preserve all existing custom logic
- Maintain code formatting and style
- Create backup of original file before modification
- Report what was changed in user-friendly format

**Hook Detection Strategy**:
- Parse TypeScript AST to find exported functions
- Match function names against expected hook names
- Detect if `dw claude log` call already exists in function body
- Identify function structure (sync vs async, return type)

**Safety Guarantees**:
- Always create `.bak` backup before modifying
- Validate syntax after modification (rollback on failure)
- Never remove existing code, only append
- Warn if hook structure is non-standard (can't safely update)
- Show diff of proposed changes before applying (with --dry-run flag)

### 8. Performance Requirements

**Logging Performance**:
- `dw claude log` must not block or slow down Claude Code
- No blocking on network calls
- Handle concurrent writes from multiple Claude Code instances

**Storage Efficiency**:
- Database size should not impact logging speed

### 9. Error Handling

**Logging Failures**:
- Never crash or block Claude Code
- Log to stderr if database unavailable

## Design Decisions

### Why Event Sourcing?
- Immutable log enables replay and analysis
- Supports future "what if" scenario testing
- No data loss from aggregation too early

### Why SQLite?
- Zero-configuration database
- File-based (fits `.darwinflow` directory model)
- Excellent performance for read-heavy + append-only workloads
- Simple full-text search support built-in
- Future extension: sqlite-vec for semantic search

## Success Criteria

1. Claude Code interactions logged without blocking or slowing down the tool
2. `dw claude init` creates new hooks file or safely updates existing one
3. System automatically captures context from `dw` tool
4. Events stored in SQLite with full-text search capability

## Out of Scope for V1

- Vector embeddings and semantic search (future extension)
- Pattern detection and workflow optimization (future extension)

## Future Extensions

### Vector Embeddings (V2)
- Add `dw claude analyze` command for batch embedding generation
- Separate embeddings table linked to events via event_id
- Support for semantic search and similarity queries

### Pattern Detection & Optimization (V3)
- Analyze event sequences to identify workflows
- Self-modifying: system adds commands based on learned patterns
- Optimization suggestions panel
- Cache generated summaries
- Learn which file access patterns are common