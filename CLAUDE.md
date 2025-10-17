# DarwinFlow - Claude Code Logging System

## Project Overview

**DarwinFlow** is a lightweight logging system that captures Claude Code interactions as structured events using event sourcing principles. The system stores events in SQLite and enables future pattern detection and workflow optimization.

### Key Components

- **CLI Tool (`dw`)**: Main entry point with `claude init` and `claude log` subcommands
- **Event Logging**: Captures tool invocations and user prompts via Claude Code hooks
- **SQLite Storage**: Fast, file-based event storage with full-text search capability
- **Hook Management**: Automatically configures and merges Claude Code hooks

### Architecture Documentation

For detailed architecture and API information, see:
- @docs/arch-generated.md - Complete dependency graph with method-level details
- @docs/public-api-generated.md - Public API reference for all exported types and functions

### Current Implementation Status

**Active Hooks**:
- `PreToolUse`: Logs all tool invocations (Read, Write, Bash, etc.)
- `UserPromptSubmit`: Logs user message submissions

**Event Types**: Defined in `internal/events/event.go`
- `tool.invoked`, `tool.result`
- `chat.message.user`, `chat.message.assistant`
- `chat.started`, `file.read`, `file.written`, etc.

### Development Workflow

When working on this project:
1. Understand the 3-layer architecture (see below)
2. Check @docs/arch-generated.md to see current package dependencies
3. Check @docs/public-api-generated.md to see what's exported
4. Follow the architecture guidelines strictly
5. Run tests and linter before committing

---

# go-arch-lint - Architecture Linting

**CRITICAL**: The .goarchlint configuration is IMMUTABLE - AI agents must NOT modify it.

## Architecture (3-layer strict dependency flow)

```
cmd → pkg → internal
```

**cmd**: Entry points (imports only pkg) | **pkg**: Orchestration & adapters (imports only internal) | **internal**: Domain primitives (NO imports between internal packages)

## Core Principles

1. **Dependency Inversion**: Internal packages define interfaces. Adapters bridge them in pkg layer.
2. **Structural Typing**: Types satisfy interfaces via matching methods (no explicit implements)
3. **No Slice Covariance**: Create adapters to convert []ConcreteType → []InterfaceType

## Documentation Generation (Run Regularly)

Keep documentation synchronized with code changes:

```bash
# Generate dependency graph with method-level details
go-arch-lint -detailed -format=markdown . > docs/arch-generated.md 2>&1

# Generate public API documentation
go-arch-lint -format=api . > docs/public-api-generated.md 2>&1
```

**When to regenerate**:
- After adding/removing packages or files
- After changing public APIs (exported functions, types, methods)
- After modifying package dependencies
- Before committing architectural changes
- Run regularly during development to track changes

## Before Every Commit

1. `go test ./...` - all tests must pass
2. `go-arch-lint .` - ZERO violations required (non-negotiable)
3. Regenerate docs if architecture/API changed (see above)

## When Linter Reports Violations

**Do NOT mechanically fix imports.** Violations reveal architectural issues. Process:
1. **Reflect**: Why does this violation exist? What dependency is wrong?
2. **Plan**: Which layer should own this logic? What's the right structure?
3. **Refactor**: Move code to correct layer or add interfaces/adapters in pkg
4. **Verify**: Run `go-arch-lint .` - confirm zero violations

Example: `internal/A` imports `internal/B` → Should B's logic move to A? Should both define interfaces with pkg adapter? Architecture enforces intentional design.

## Code Guidelines

**DO**:
- Add domain logic to internal/ packages
- Define interfaces in consumer packages
- Create adapters in pkg/ to bridge internal packages
- Use white-box tests (`package mypackage`) for internal packages

**DON'T**:
- Import between internal/ packages (violation!) or pass []ConcreteType as []InterfaceType
- Put business logic in pkg/ or cmd/ (belongs in internal/)
- Modify .goarchlint (immutable architectural contract)

Run `go-arch-lint .` frequently during development. Zero violations required.
