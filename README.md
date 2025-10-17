# DarwinFlow

**Capture, store, and analyze your Claude Code interactions**

DarwinFlow is a lightweight logging system that automatically captures all Claude Code interactions as structured events. Built with event sourcing principles, it enables pattern detection, workflow optimization, and deep insights into your AI-assisted development sessions.

## Features

- **Automatic Logging**: Captures all Claude Code events via hooks (tool invocations, user prompts, etc.)
- **Event Sourcing**: Immutable event log enabling replay and analysis
- **SQLite Storage**: Fast, file-based storage with full-text search
- **Zero Performance Impact**: Non-blocking, concurrent-safe logging
- **Context-Aware**: Automatically detects project context from environment
- **Clean Architecture**: Strict 3-layer design (`cmd → pkg → internal`)

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/kgatilin/darwinflow-pub.git
cd darwinflow-pub

# Build the CLI
go build -o dw ./cmd/dw

# Install to your PATH (optional)
go install ./cmd/dw
```

### Initialize Logging

```bash
# Set up Claude Code logging infrastructure
dw claude init
```

This will:
- Create the SQLite database at `~/.darwinflow/logs/events.db`
- Add hooks to your Claude Code settings (typically `~/.claude/settings.json`)
- Configure automatic event capture for PreToolUse and UserPromptSubmit hooks

### Start Using Claude Code

After running `dw claude init`, restart Claude Code. All your interactions will now be automatically logged!

## Architecture

DarwinFlow follows a strict 3-layer architecture enforced by [go-arch-lint](https://github.com/fdaines/go-arch-lint):

```
cmd → pkg → internal
```

- **cmd/dw**: CLI entry points (`dw claude init`, `dw claude log`)
- **pkg/claude**: Orchestration layer (settings management, logging coordination)
- **internal/**: Domain primitives (events, hooks config, storage interfaces)

### Key Components

- **Events** (`internal/events`): Event types and payload definitions
- **Hooks** (`internal/hooks`): Claude Code hook configuration and merging logic
- **Storage** (`internal/storage`): Storage interface definitions
- **Logger** (`pkg/claude`): Event logging and database interaction
- **Settings Manager** (`pkg/claude`): Claude Code settings file management

## Usage

### Commands

```bash
# Initialize logging infrastructure
dw claude init

# Log an event (typically called by hooks)
dw claude log <event-type>
```

### Event Types

Currently captured events:

- `tool.invoked` - Claude Code tool invocation (Read, Write, Bash, etc.)
- `chat.message.user` - User prompt submission

### Environment Variables

- `DW_CONTEXT` - Set the current context (e.g., `project/myapp`)
- `DW_MAX_PARAM_LENGTH` - Maximum parameter length for logging (default: 30)

## Development

### Prerequisites

- Go 1.25.1 or later
- [go-arch-lint](https://github.com/fdaines/go-arch-lint) for architecture validation

### Building

```bash
# Build the CLI
make

# Run tests
make test
```

### Architecture Compliance

Before committing, ensure:

1. All tests pass: `go test ./...`
2. Zero architecture violations: `go-arch-lint .`
3. Documentation is up-to-date (see [CLAUDE.md](./CLAUDE.md))

### Generated Documentation

Architecture and API documentation is generated automatically:

```bash
# Generate dependency graph
go-arch-lint -detailed -format=markdown . > docs/arch-generated.md

# Generate public API reference
go-arch-lint -format=api . > docs/public-api-generated.md
```

## Project Structure

```
darwinflow-pub/
├── cmd/dw/              # CLI entry points
│   ├── main.go          # Main command router
│   └── claude.go        # Claude subcommand handlers
├── pkg/claude/          # Orchestration & adapters
│   ├── logger.go        # Event logging
│   ├── settings.go      # Settings file management
│   ├── sqlite.go        # SQLite storage adapter
│   ├── transcript.go    # Transcript parsing
│   └── cli.go           # CLI helper functions
├── internal/            # Domain primitives
│   ├── events/          # Event definitions
│   ├── hooks/           # Hook configuration
│   └── storage/         # Storage interfaces
├── docs/                # Generated documentation
│   ├── arch-generated.md      # Dependency graph
│   └── public-api-generated.md # Public API reference
├── CLAUDE.md            # AI agent instructions
└── README.md            # This file
```

## Roadmap

### V1 (Current)
- ✅ Basic event capture (PreToolUse, UserPromptSubmit)
- ✅ SQLite storage with full-text search
- ✅ Hook management and merging

### V2 (Planned)
- Vector embeddings for semantic search
- Pattern detection across sessions
- Enhanced context extraction

### V3 (Future)
- Workflow optimization suggestions
- Self-modifying commands based on patterns
- Advanced analytics and insights

## Contributing

Contributions are welcome! Please ensure:

1. Code follows the 3-layer architecture
2. All tests pass (`go test ./...`)
3. Architecture linter passes (`go-arch-lint .`)
4. Documentation is updated for API/architecture changes

## License

MIT License - See LICENSE file for details

## Acknowledgments

Built to enhance [Claude Code](https://www.anthropic.com/claude/code) workflows and enable AI-assisted development insights.
