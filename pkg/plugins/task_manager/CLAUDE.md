# Package: task_manager

**Path**: `pkg/plugins/task_manager`

**Role**: Hierarchical roadmap management plugin (Roadmap → Track → Task → Iteration → ADR → AcceptanceCriteria)

---

## ⚠️ REFACTORING IN PROGRESS ⚠️

**Status**: Undergoing clean architecture refactoring (Iteration #27)

**Active Iteration**: #27 - TM Refactoring: Foundation Layers (Domain/Application/Infrastructure)

**Related Tasks**:
- TM-task-126: Phase 1 - Domain Layer Extraction (in-progress)
- TM-task-127: Phase 2 - Infrastructure Layer (in-progress)
- TM-task-128: Phase 3 - Application Layer with Unified Services (COMPLETE ✅)
- TM-task-136: Repository Interface Segregation (in-progress)
- TM-task-137: Update CLI Adapters to Use Unified Services (todo - BLOCKING)

**Current State**: Application layer complete with unified services. CLI adapters need update to use new services.

**Target Architecture**: Clean architecture with DDD principles (Domain → Application → Infrastructure → Presentation)

**Documentation References**:
- **Main Architecture**: `/workspace/CLAUDE.md` - DarwinFlow architecture guide
- **Domain Layer**: `domain/CLAUDE.md` - Domain entities, services, events, repository interfaces
- **Application Layer**: `application/CLAUDE.md` - Unified application services (TM-task-128 ✅)
- **Infrastructure Layer**: `infrastructure/CLAUDE.md` - Repository implementations, migrations
- **Presentation Layer**: `presentation/CLAUDE.md` - CLI adapters (work in progress)
- **Plugin SDK**: `pkg/pluginsdk/CLAUDE.md` - SDK documentation

---

## Overview

The task-manager plugin provides comprehensive project/product roadmap management with:
- **Multi-project support** - Separate isolated roadmaps (e.g., "production" vs "test")
- **Hierarchical structure** - Roadmap → Track → Task → Iteration → ADR → AcceptanceCriteria
- **SQLite database storage** - Per-project databases with full schema management
- **Full CLI commands** - Comprehensive commands across all entities
- **Event bus integration** - Cross-plugin communication with event types
- **Interactive TUI** - Visualization and management with project context
- **Event sourcing** - Complete audit trail for all changes
- **Clean Architecture** - DDD with separated layers (Domain, Application, Infrastructure, Presentation)

---

## Target Architecture (Clean Architecture + DDD)

### Layer Structure

```
pkg/plugins/task_manager/
├── domain/                          # Pure business logic (zero external dependencies)
│   ├── entities/                    # Domain entities (7 aggregates)
│   │   ├── roadmap_entity.go
│   │   ├── track_entity.go
│   │   ├── task_entity.go
│   │   ├── iteration_entity.go
│   │   ├── adr_entity.go
│   │   ├── acceptance_criteria_entity.go
│   │   └── *_entity_test.go
│   ├── services/                    # Domain services (business rules)
│   │   ├── validation_service.go    # ID validation, format validation
│   │   ├── dependency_service.go    # Circular dependency detection (DFS)
│   │   ├── iteration_service.go     # Iteration lifecycle validation
│   │   └── *_service_test.go
│   ├── events/                      # Domain events
│   │   └── events.go
│   ├── repositories/                # Repository interfaces (6 focused interfaces)
│   │   ├── roadmap_repository.go    # Roadmap CRUD
│   │   ├── track_repository.go      # Track CRUD + dependencies
│   │   ├── task_repository.go       # Task CRUD + iteration management
│   │   ├── iteration_repository.go  # Iteration CRUD + lifecycle
│   │   ├── adr_repository.go        # ADR CRUD + track relationship
│   │   └── acceptance_criteria_repository.go  # AC CRUD + verification
│   └── CLAUDE.md                    # Domain layer architectural guidance
│
├── application/                     # Application services (use cases) ✅ COMPLETE
│   ├── track_service.go             # All track operations (CRUD + dependencies)
│   ├── task_service.go              # All task operations (CRUD + move)
│   ├── iteration_service.go         # All iteration operations (CRUD + lifecycle)
│   ├── adr_service.go               # All ADR operations (CRUD + status)
│   ├── ac_service.go                # All AC operations (CRUD + verification)
│   ├── dto/                         # Data Transfer Objects
│   │   ├── track_dto.go
│   │   ├── task_dto.go
│   │   ├── iteration_dto.go
│   │   ├── adr_dto.go
│   │   ├── ac_dto.go
│   │   └── helpers.go
│   ├── *_service_test.go            # Service tests (126 tests, 82.1% coverage)
│   └── CLAUDE.md                    # Application layer architectural guidance
│
├── infrastructure/                  # Infrastructure implementations
│   └── persistence/                 # Repository implementations (6 focused repos)
│       ├── roadmap_repository.go
│       ├── track_repository.go
│       ├── task_repository.go
│       ├── iteration_repository.go
│       ├── adr_repository.go
│       ├── acceptance_criteria_repository.go
│       ├── migrations.go            # Database migrations
│       ├── event_emitting_repository.go  # Event emission decorator
│       └── CLAUDE.md                # Infrastructure layer guidance
│
├── presentation/                    # Presentation layer (CLI adapters)
│   └── cli/                         # CLI framework adapters (⚠️ NEEDS UPDATE)
│       ├── track_adapters.go        # Bridge CLI → application services
│       ├── task_adapters.go
│       ├── iteration_adapters.go
│       ├── adr_adapters.go
│       ├── ac_adapters.go
│       └── CLAUDE.md                # Presentation layer guidance
│
└── plugin.go                        # Plugin registration and wiring
```

### Dependency Flow (Clean Architecture)

```
Presentation (CLI adapters)
    ↓ depends on
Application (Services + DTOs)
    ↓ depends on
Domain (Entities + Interfaces + Services)
    ↑ implemented by
Infrastructure (Repository implementations)
```

**Key Principle**: Dependencies point inward. Domain layer has zero external dependencies.

---

## Domain Model

### Aggregates (7 total)

**Roadmap (Root Aggregate)**
- Single active roadmap per project
- Contains vision and success criteria
- Parent to all tracks

**Track (Major Work Area)**
- Represents work streams (e.g., "Framework Core", "Plugin System")
- Has status (not-started, in-progress, complete, blocked, waiting)
- Has priority (critical, high, medium, low)
- Can depend on other tracks (with circular dependency prevention)
- Contains multiple tasks
- Associated with ADRs (Architecture Decision Records)

**Task (Concrete Work Item)**
- Belongs to a track
- Has status (todo, in-progress, done)
- Can have git branch association
- Atomic unit of work
- Can be grouped into iterations
- Has acceptance criteria

**Iteration (Time-Boxed Grouping)**
- Groups tasks from multiple tracks
- Has status (planned, current, complete)
- Only one can be "current" at a time
- Auto-incrementing iteration numbers
- Deliverable-oriented (goal, deliverable description)

**ADR (Architecture Decision Record)**
- Documents architectural decisions for tracks
- Has status (proposed, accepted, rejected, superseded, deprecated)
- Links to specific track
- Can supersede other ADRs
- Immutable once accepted (create new ADR to change)

**AcceptanceCriteria (Task Verification)**
- Defines what "done" means for a task
- Has description and testing instructions
- Has status (not-started, pending-review, verified, failed)
- Can be verified or failed with feedback
- Links to specific task

---

## Current Implementation Status

### ✅ Complete (TM-task-128)
- **Application Layer**: 5 unified services with 82.1% test coverage
- **DTOs**: Input/output types for all services
- **Tests**: 126 comprehensive tests (all passing)

### 🚧 In Progress
- **Domain Layer** (TM-task-126): Entities, services, repository interfaces
- **Infrastructure Layer** (TM-task-127): Repository implementations, migrations
- **Repository Segregation** (TM-task-136): Split monolithic repository into 6 focused interfaces

### ⚠️ Blocking (TM-task-137)
- **CLI Adapters**: Need update to use new application services
- **Issue**: Adapters still import deleted CQRS packages (application/commands, application/queries)
- **Impact**: Full build fails, CLI commands won't work until adapters updated

### Migration Status

**Deleted (TM-task-128)**:
- Old CQRS pattern (application/commands/, application/queries/)
- 11 command/query handler files

**Created (TM-task-128)**:
- 5 unified application services (Track, Task, Iteration, ADR, AC)
- 6 DTO files (one per aggregate + helpers)
- 10 comprehensive test files

**Next Steps (TM-task-137)**:
- Update 5 CLI adapter files to use application services
- Remove imports to deleted packages
- Full build/test suite will pass after completion

---

## Multi-Project Architecture

**Project Isolation:**
- Each project has its own SQLite database in `.darwinflow/projects/<project-name>/roadmap.db`
- Active project tracked in `.darwinflow/active-project.txt`
- Complete data isolation between projects
- Auto-migration from legacy single-database structure

**Use Cases:**
- Separate "production" and "test" roadmaps
- Multiple product roadmaps in one workspace
- Experimentation without affecting real data

**Commands:**
- All entity commands support `--project <name>` flag to override active project
- 5 dedicated project management commands (create, list, switch, show, delete)

---

## Database Schema

**8 Tables** (post-refactoring):
- `roadmaps` - Roadmap entities (id, vision, success_criteria)
- `tracks` - Track entities (id, roadmap_id, title, description, status, priority, rank)
- `track_dependencies` - Track dependency relationships (track_id, depends_on_id)
- `tasks` - Task entities (id, track_id, title, description, status, rank, branch)
- `iterations` - Iteration entities (number, roadmap_id, name, goal, status, deliverable)
- `iteration_tasks` - Iteration-task relationships (iteration_number, task_id)
- `adrs` - Architecture Decision Records (id, track_id, title, context, decision, status)
- `acceptance_criteria` - Acceptance criteria (id, task_id, description, testing_instructions, status)

All tables have:
- Primary keys and foreign keys
- Proper indexes on frequently queried columns
- Created_at and updated_at timestamps
- Referential integrity constraints

---

## Commands Overview

**Note**: Commands currently use old implementation. After TM-task-137 completion, they will use new application services.

### Project Commands (5 commands)

```bash
dw task-manager project create <name>
dw task-manager project list
dw task-manager project switch <name>
dw task-manager project show
dw task-manager project delete <name> [--force]
```

### Roadmap Commands (3 commands)

```bash
dw task-manager roadmap init --vision "..." --success-criteria "..."
dw task-manager roadmap show
dw task-manager roadmap update [--vision "..."] [--success-criteria "..."]
```

### Track Commands (7 commands)

```bash
dw task-manager track create --id <id> --title <title> [--description] [--rank]
dw task-manager track list [--status <status>]
dw task-manager track show <track-id>
dw task-manager track update <track-id> [--title] [--description] [--status] [--rank]
dw task-manager track delete <track-id> [--force]
dw task-manager track add-dependency <track-id> <depends-on>
dw task-manager track remove-dependency <track-id> <depends-on>
```

### Task Commands (7 commands)

```bash
dw task-manager task create --track <track-id> --title <title> [--description] [--rank]
dw task-manager task list [--track <track-id>] [--status <status>]
dw task-manager task show <task-id>
dw task-manager task update <task-id> [--title] [--description] [--status] [--rank] [--branch]
dw task-manager task delete <task-id> [--force]
dw task-manager task move <task-id> --track <new-track-id>
dw task-manager task validate <task-id>  # Validate acceptance criteria
```

### Iteration Commands (10 commands)

```bash
dw task-manager iteration create --name <name> --goal <goal> --deliverable <deliverable>
dw task-manager iteration list
dw task-manager iteration show <iteration-number> [--full]
dw task-manager iteration current
dw task-manager iteration update <number> [--name] [--goal] [--deliverable]
dw task-manager iteration start <iteration-number>
dw task-manager iteration complete <iteration-number>
dw task-manager iteration add-task <iteration> <task-id> [<task-id>...]
dw task-manager iteration remove-task <iteration> <task-id> [<task-id>...]
dw task-manager iteration delete <iteration-number> [--force]
```

### ADR Commands (7 commands)

```bash
dw task-manager adr create <track-id> --title <title> --context <context> --decision <decision>
dw task-manager adr list [--track <track-id>] [--status <status>]
dw task-manager adr show <adr-id>
dw task-manager adr update <adr-id> [--title] [--context] [--decision] [--consequences] [--alternatives]
dw task-manager adr supersede <adr-id> --superseded-by <new-adr-id>
dw task-manager adr deprecate <adr-id>
dw task-manager adr check <track-id>  # Check if track has required ADR
```

### Acceptance Criteria Commands (9 commands)

```bash
dw task-manager ac add <task-id> --description <desc> --testing-instructions <instructions>
dw task-manager ac list <task-id>
dw task-manager ac list-iteration <iteration-number>
dw task-manager ac show <ac-id>
dw task-manager ac update <ac-id> [--description] [--testing-instructions]
dw task-manager ac verify <ac-id>
dw task-manager ac fail <ac-id> --feedback <feedback>
dw task-manager ac failed [--task <task-id>] [--iteration <iteration-number>]
dw task-manager ac delete <ac-id> [--force]
```

### TUI Command

```bash
dw task-manager tui
```

**Total Commands**: ~48 commands across all entities

---

## Event Bus Integration

The plugin emits events for all CRUD operations:

**Roadmap Events:**
- `task-manager.roadmap.created`
- `task-manager.roadmap.updated`

**Track Events:**
- `task-manager.track.created`
- `task-manager.track.updated`
- `task-manager.track.status_changed`
- `task-manager.track.completed`
- `task-manager.track.blocked`

**Task Events:**
- `task-manager.task.created`
- `task-manager.task.updated`
- `task-manager.task.status_changed`
- `task-manager.task.completed`
- `task-manager.task.moved`

**Iteration Events:**
- `task-manager.iteration.created`
- `task-manager.iteration.updated`
- `task-manager.iteration.started`
- `task-manager.iteration.completed`

**ADR Events:**
- `task-manager.adr.created`
- `task-manager.adr.updated`
- `task-manager.adr.superseded`
- `task-manager.adr.deprecated`

**Acceptance Criteria Events:**
- `task-manager.ac.created`
- `task-manager.ac.verified`
- `task-manager.ac.failed`

Other plugins can subscribe to these events for notifications, automation, etc.

---

## Testing

**Application Layer Coverage**: 82.1% (126 tests) ✅

**Test Organization** (post-refactoring):
- Domain entity tests: Validation, state transitions, business rules
- Domain service tests: Circular dependencies, lifecycle validation
- Application service tests: 126 tests (all CRUD + business operations)
- Infrastructure repository tests: SQLite integration, migrations
- CLI adapter tests: Command execution, flag parsing

**Running Tests**:

```bash
# Application layer (new unified services)
go test ./pkg/plugins/task_manager/application/... -v -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
# Expected: 82.1% coverage

# Domain layer
go test ./pkg/plugins/task_manager/domain/... -v

# Infrastructure layer
go test ./pkg/plugins/task_manager/infrastructure/... -v

# All tests (will fail until TM-task-137 complete)
go test ./pkg/plugins/task_manager/... -v
```

---

## ⚠️ CRITICAL TESTING RULES ⚠️

### 🚫 NEVER USE REAL DIRECTORIES IN TESTS 🚫

**ABSOLUTE RULE: ALL tests MUST use `t.TempDir()` for file operations.**

**❌ NEVER DO THIS:**
```go
func TestBadExample(t *testing.T) {
    homeDir, _ := os.UserHomeDir()
    testDir := filepath.Join(homeDir, ".darwinflow")  // ❌ BAD!
    // Writing to real user directories in tests
}
```

**✅ ALWAYS DO THIS:**
```go
func TestGoodExample(t *testing.T) {
    tmpDir := t.TempDir()  // ✅ GOOD! Auto-cleanup after test
    testDir := filepath.Join(tmpDir, ".darwinflow")
    // All file operations in isolated temp directory
}
```

### Why This Matters

1. **Test Isolation**: Tests must not interfere with each other or real data
2. **CI/CD**: Tests run in clean environments without side effects
3. **Reproducibility**: Tests must produce identical results on every run
4. **Safety**: Never risk corrupting real user data during testing
5. **Cleanup**: `t.TempDir()` automatically removes test files after completion

---

## Plugin Architecture

### Plugin Interface Implementation

**TaskManagerPlugin** implements:
- `pluginsdk.Plugin` - Base plugin interface
- `pluginsdk.IEntityProvider` - Query roadmaps, tracks, tasks, iterations, ADRs, ACs
- `pluginsdk.ICommandProvider` - All CLI commands

**Key Methods:**
- `GetInfo()` - Plugin metadata (name, version, description)
- `GetCapabilities()` - Lists implemented capabilities
- `GetEntityTypes()` - Returns entity types (roadmap, track, task, iteration, adr, ac)
- `Query(ctx, query)` - Query entities with filters and pagination
- `GetEntity(ctx, id)` - Get entity by ID
- `UpdateEntity(ctx, id, fields)` - Update entity fields
- `GetCommands()` - Returns all CLI commands

---

## Key Design Decisions

1. **Clean Architecture**: Separated concerns (Domain → Application → Infrastructure → Presentation)
2. **Unified Services**: One service per aggregate handling all operations (NOT CQRS)
3. **Repository Segregation**: 6 focused repository interfaces (ISP compliance)
4. **Event Sourcing**: All changes emit events for audit trail and cross-plugin notifications
5. **SQLite Persistence**: Reliable local storage without external dependencies
6. **TUI Integration**: Bubble Tea framework for rich terminal user experience
7. **Track Dependencies**: Enables workflow management and blocking detection
8. **Iteration Grouping**: Time-boxed work organizing across tracks
9. **ADR Pattern**: Documenting architectural decisions at track level
10. **Acceptance Criteria**: Explicit verification requirements for tasks

---

## Refactoring Progress

**Iteration #27**: TM Refactoring: Foundation Layers (Domain/Application/Infrastructure)

**Goal**: Extract and establish the three foundation layers of clean architecture

**Deliverable**: Complete domain/, application/, and infrastructure/ packages with full test coverage

**Status**: Application layer complete (82.1% coverage). CLI adapters need update (TM-task-137).

**Timeline**:
- Phase 1 (TM-task-126): Domain Layer Extraction - in-progress
- Phase 2 (TM-task-127): Infrastructure Layer - in-progress
- Phase 3 (TM-task-128): Application Layer - COMPLETE ✅
- Phase 4 (TM-task-136): Repository Segregation - in-progress
- Phase 5 (TM-task-137): CLI Adapter Update - todo (BLOCKING)

**Next Steps**:
1. Complete TM-task-137 (update CLI adapters to use application services)
2. Verify full build passes
3. Complete remaining phases (TM-task-126, TM-task-127, TM-task-136)
4. Close iteration #27

---

## References

- **Main Architecture**: `/workspace/CLAUDE.md` - DarwinFlow architecture guide
- **Domain Layer**: `domain/CLAUDE.md` - Domain layer architectural guidance
- **Application Layer**: `application/CLAUDE.md` - Application layer architectural guidance (TM-task-128)
- **Infrastructure Layer**: `infrastructure/CLAUDE.md` - Infrastructure layer guidance
- **Presentation Layer**: `presentation/CLAUDE.md` - Presentation layer guidance
- **Plugin SDK**: `pkg/pluginsdk/CLAUDE.md` - SDK documentation
- **Iteration #27**: See `dw task-manager iteration show 27 --full` for refactoring details

---

*Last Updated: 2025-11-11 (Iteration #27 - Application Layer Complete)*
