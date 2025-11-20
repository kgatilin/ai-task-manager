# Task Manager (tm)

**A hierarchical roadmap and task management CLI tool built with Clean Architecture**

Task Manager is a developer-focused CLI tool for managing complex projects through a structured hierarchy: Roadmap → Tracks → Tasks → Iterations. Built with Domain-Driven Design principles and SQLite storage, it provides powerful command-line workflows and an interactive terminal UI for planning, tracking, and executing software development work.

## Features

- **Hierarchical Project Structure**: Organize work from vision (Roadmap) down to concrete tasks
  - **Roadmap**: Project vision and success criteria
  - **Tracks**: Major work streams with dependencies and priorities
  - **Tasks**: Atomic work units with status tracking (todo/in-progress/done)
  - **Iterations**: Time-boxed groupings for sprint planning
- **Acceptance Criteria**: Task verification with detailed testing instructions
- **Architecture Decision Records (ADRs)**: Document architectural choices and their rationale
- **Document Management**: Plans, retrospectives, and other project documentation
- **Multi-Project Support**: Isolated project databases for separate roadmaps
- **Interactive TUI**: Keyboard-driven terminal interface for browsing and managing work
- **Clean Architecture**: Strict separation of concerns with DDD principles
- **SQLite Storage**: Fast, file-based persistence per project

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/kgatilin/ai-task-manager.git
cd ai-task-manager

# Build the CLI
make build

# Optional: Install to GOPATH/bin
make install
```

### Basic Usage

```bash
# Create a project
./tm project create myproject

# Initialize a roadmap
./tm roadmap init \
  --vision "Build a scalable event-driven architecture" \
  --success-criteria "Support 10+ microservices, 99.9% uptime"

# Create a track (major work area)
./tm track create \
  --title "Core Infrastructure" \
  --description "Foundation services and frameworks" \
  --priority high

# Create a task
./tm task create \
  --track TM-track-1 \
  --title "Implement event bus" \
  --priority high

# Create an iteration
./tm iteration create \
  --name "Sprint 1" \
  --goal "Complete core infrastructure" \
  --deliverable "Event bus and message broker"

# Add tasks to iteration
./tm iteration add-task 1 TM-task-1

# Launch interactive UI
./tm ui
```

## Core Concepts

### Project Hierarchy

```
Roadmap (vision + success criteria)
  └── Track (work stream)
        └── Task (concrete work)
              └── Acceptance Criteria (verification steps)
```

**Iterations** group tasks across tracks for time-boxed execution (sprints).

### Entity Types

- **Roadmap**: Root entity defining project vision and success criteria (one per project)
- **Track**: Major work stream with status, priority, and dependencies
- **Task**: Atomic work unit with title, description, status (todo/in-progress/done)
- **Iteration**: Time-boxed grouping with goal, deliverable, and task membership
- **Acceptance Criteria**: Task verification with description and testing instructions
- **Document**: ADRs, plans, retrospectives attached to tracks or iterations

### Multi-Project Isolation

Each project has its own SQLite database (`.tm/projects/<name>/roadmap.db`), providing complete isolation between roadmaps. Use `tm project switch` to change the active project, or use `--project <name>` flag to override on any command.

## Architecture

Task Manager follows Clean Architecture with Domain-Driven Design:

```
cmd/tm → internal/task_manager → [domain, application, infrastructure, presentation]
```

### Layers

- **Domain** (`internal/task_manager/domain/`): Pure business logic, zero dependencies
  - Entities: 7 aggregates (Roadmap, Track, Task, Iteration, ADR, AcceptanceCriteria, Document)
  - Services: Domain services (validation, dependency detection, lifecycle management)
  - Repositories: Interfaces only (implemented in infrastructure)
- **Application** (`internal/task_manager/application/`): Use cases and orchestration
  - Services: High-level operations coordinating domain logic
  - DTOs: Data transfer objects for external communication
  - Mocks: Generated mocks for testing (126 tests, 82.1% coverage)
- **Infrastructure** (`internal/task_manager/infrastructure/`): Technical implementations
  - Persistence: SQLite repository implementations with migrations
  - Event emission: Domain event publishing decorator
- **Presentation** (`internal/task_manager/presentation/`): User interfaces
  - CLI: ~48 Cobra commands for all operations
  - TUI: Interactive terminal UI (Bubble Tea framework)

**Dependency Rule**: Dependencies flow inward only. Domain has zero external dependencies.

## Commands

### Project Management

```bash
# Create a new project
tm project create <name>

# List all projects (* marks active)
tm project list

# Switch active project
tm project switch <name>

# Show current project
tm project show

# Delete a project
tm project delete <name> --force
```

### Roadmap Commands

```bash
# Initialize roadmap
tm roadmap init \
  --vision "Your project vision" \
  --success-criteria "Measurable success criteria"

# Show roadmap
tm roadmap show

# Update roadmap
tm roadmap update \
  --vision "Updated vision" \
  --success-criteria "New criteria"
```

### Track Commands (Work Streams)

```bash
# Create track
tm track create \
  --title "Track Title" \
  --description "Description" \
  --priority high|medium|low \
  --rank 100

# List tracks
tm track list
tm track list --status in-progress --priority high

# Show track details
tm track show TM-track-1

# Update track
tm track update TM-track-1 \
  --status planning|in-progress|done|blocked \
  --priority critical|high|medium|low

# Manage dependencies
tm track add-dependency TM-track-2 TM-track-1    # track-2 depends on track-1
tm track remove-dependency TM-track-2 TM-track-1

# Delete track
tm track delete TM-track-1 --force
```

### Task Commands (Work Items)

```bash
# Create task
tm task create \
  --track TM-track-1 \
  --title "Task Title" \
  --description "Description" \
  --priority high|medium|low \
  --rank 100

# List tasks
tm task list
tm task list --track TM-track-1 --status todo

# Show task details
tm task show TM-task-1

# Update task
tm task update TM-task-1 \
  --status todo|in-progress|done \
  --priority high|medium|low \
  --branch feat/my-feature

# Move task to different track
tm task move TM-task-1 --track TM-track-2

# Delete task
tm task delete TM-task-1 --force
```

### Iteration Commands (Sprints)

```bash
# Create iteration
tm iteration create \
  --name "Sprint 1" \
  --goal "Sprint goal" \
  --deliverable "Expected deliverable"

# List iterations
tm iteration list

# Show iteration details
tm iteration show 1

# Show current iteration (or next planned)
tm iteration current

# Add/remove tasks
tm iteration add-task 1 TM-task-1 TM-task-2
tm iteration remove-task 1 TM-task-1

# Start iteration (mark as current)
tm iteration start 1

# Complete iteration
tm iteration complete 1

# Delete iteration
tm iteration delete 1 --force
```

### Acceptance Criteria Commands

```bash
# Add acceptance criterion
tm ac add TM-task-1 \
  --description "What must be verified" \
  --testing-instructions "1. Run command X\n2. Verify Y\n3. Check Z"

# List acceptance criteria
tm ac list TM-task-1
tm ac list-iteration 1               # All ACs in iteration

# Show AC details
tm ac show TM-ac-1

# Update AC
tm ac update TM-ac-1 \
  --description "Updated description" \
  --testing-instructions "Updated instructions"

# Verify AC
tm ac verify TM-ac-1

# Mark AC as failed
tm ac fail TM-ac-1 --feedback "Error message or reason"

# List failed ACs
tm ac failed                          # All failed
tm ac failed --iteration 1            # Failed in iteration 1
tm ac failed --task TM-task-1         # Failed for task

# Delete AC
tm ac delete TM-ac-1 --force
```

### Document Commands (ADRs, Plans, etc.)

```bash
# Create document
tm doc create \
  --title "ADR: Use Event Sourcing" \
  --type adr|plan|retrospective \
  --content "Document content..." \
  --track TM-track-1

# Or from file
tm doc create \
  --title "ADR: Use Event Sourcing" \
  --type adr \
  --from-file ./docs/adr-001.md \
  --track TM-track-1

# List documents
tm doc list
tm doc list --type adr

# Show document
tm doc show TM-doc-1

# Update document
tm doc update TM-doc-1 --from-file ./updated.md

# Attach to track or iteration
tm doc attach TM-doc-1 --track TM-track-2
tm doc attach TM-doc-1 --iteration 3

# Detach document
tm doc detach TM-doc-1

# Delete document
tm doc delete TM-doc-1 --force
```

### ADR Commands (Architecture Decision Records)

```bash
# Create ADR
tm adr create \
  --title "Use microservices architecture" \
  --context "Need to scale independently" \
  --decision "Adopt microservices with event bus" \
  --consequences "Increased complexity, better scalability" \
  --track TM-track-1

# List ADRs
tm adr list
tm adr list --status accepted --track TM-track-1

# Show ADR
tm adr show TM-adr-1

# Update ADR
tm adr update TM-adr-1 \
  --status proposed|accepted|rejected|superseded|deprecated

# Supersede ADR (with reason)
tm adr supersede TM-adr-1 TM-adr-2

# Deprecate ADR
tm adr deprecate TM-adr-1 --reason "No longer relevant"

# Check for superseded/deprecated ADRs
tm adr check

# Delete ADR
tm adr delete TM-adr-1 --force
```

### Interactive TUI

```bash
# Launch terminal UI
tm ui
```

**Navigation:**
- `j/k` or `↑/↓` - Move up/down
- `Enter` - Select/drill down
- `i` - Switch to iteration view
- `r` - Refresh data
- `Esc` - Go back
- `q` - Quit

**Features:**
- Roadmap overview with tracks and tasks
- Track details with nested task lists
- Iteration planning and progress
- Dependency visualization
- Status and priority filtering

## Development

### Prerequisites

- Go 1.25.1 or later
- SQLite3 (included via `mattn/go-sqlite3`)

### Building

```bash
# Build binary
make build          # Creates ./tm

# Run tests
make test           # Run all tests

# Install to GOPATH/bin
make install
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Test Structure:**
- **Unit tests**: `*_test.go` in each package (black-box pattern)
- **Integration tests**: Infrastructure layer with real SQLite
- **E2E tests**: `internal/task_manager/e2e_test/` - full CLI command validation

**Coverage Target**: 70-80% (current: domain 90%+, application 82.1%, infrastructure 52.3%)

### Architecture Validation

Task Manager uses strict dependency rules enforced by linting:

```bash
# Validate architecture compliance
go-arch-lint .

# Regenerate architecture documentation
go-arch-lint docs
```

**Dependency Rules:**
- Domain → NOTHING (pure business logic)
- Application → Domain
- Infrastructure → Domain
- Presentation → Application + Infrastructure

## Project Structure

```
ai-task-manager/
├── cmd/tm/                              # CLI entry point
│   ├── main.go                          # Bootstrap and command registration
│   └── ...                              # Command handlers
├── internal/task_manager/               # Task manager implementation
│   ├── domain/                          # Business logic (zero dependencies)
│   │   ├── entities/                    # Roadmap, Track, Task, Iteration, etc.
│   │   ├── services/                    # Domain services (validation, dependencies)
│   │   ├── events/                      # Domain events
│   │   └── repositories/                # Repository interfaces
│   ├── application/                     # Use cases and orchestration
│   │   ├── *_service.go                 # Application services
│   │   ├── dto/                         # Data transfer objects
│   │   └── mocks/                       # Generated mocks (mockery)
│   ├── infrastructure/                  # Technical implementations
│   │   └── persistence/                 # SQLite repositories + migrations
│   ├── presentation/                    # User interfaces
│   │   └── cli/                         # Cobra command adapters
│   ├── e2e_test/                        # End-to-end tests
│   └── plugin.go                        # Dependency injection
├── docs/                                # Documentation
├── Makefile                             # Build automation
├── go.mod                               # Go module definition
├── CLAUDE.md                            # AI agent instructions
└── README.md                            # This file
```

## Database

Each project uses a SQLite database stored at `.tm/projects/<name>/roadmap.db`.

**Schema (8 tables):**
- `roadmaps` - One per project (vision, success criteria)
- `tracks` - Work streams with status/priority
- `track_dependencies` - Track dependencies (junction table)
- `tasks` - Work items with status
- `iterations` - Time-boxed groupings
- `iteration_tasks` - Iteration membership (junction table)
- `adrs` - Architecture decision records
- `acceptance_criteria` - Task verification criteria
- `documents` - Plans, retrospectives, etc.

Migrations run automatically on first access to ensure schema is up-to-date.

## Key Dependencies

- **Cobra**: CLI framework for command structure
- **Bubble Tea**: Terminal UI framework for interactive TUI
- **SQLite3**: Embedded database for persistence
- **Glamour**: Markdown rendering for terminal
- **Lipgloss**: Terminal styling and layouts

## Contributing

Contributions are welcome! Please ensure:

1. Code follows Clean Architecture principles
2. All tests pass (`make test`)
3. Domain layer has zero external dependencies
4. Coverage remains above 70% for new code
5. Documentation is updated for new features

## License

MIT License - See LICENSE file for details

