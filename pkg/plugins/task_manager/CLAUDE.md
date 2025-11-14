# Package: task_manager

**Path**: `pkg/plugins/task_manager`

**Role**: Hierarchical roadmap management (Roadmap → Track → Task → Iteration → ADR → AcceptanceCriteria)

**Architecture**: Clean architecture with DDD (Domain → Application → Infrastructure → Presentation)

**Database**: SQLite per-project (`.darwinflow/projects/<name>/roadmap.db`)

---

## Package Structure

```
pkg/plugins/task_manager/
├── domain/                          # Business logic (zero external dependencies)
│   ├── entities/                    # 7 aggregates
│   │   ├── roadmap_entity.go        # Root aggregate (vision, success criteria)
│   │   ├── track_entity.go          # Work streams (status, priority, dependencies)
│   │   ├── task_entity.go           # Atomic work units (todo/in-progress/done)
│   │   ├── iteration_entity.go      # Time-boxed groupings (planned/current/complete)
│   │   ├── adr_entity.go            # Architecture decisions (proposed/accepted/rejected)
│   │   ├── acceptance_criteria_entity.go  # Task verification (not-started/verified/failed)
│   │   ├── value_objects.go         # Status/Priority enums, ID types
│   │   └── *_entity_test.go
│   ├── services/                    # Domain services (stateless business logic)
│   │   ├── validation_service.go    # ID validation, format checks
│   │   ├── dependency_service.go    # Circular dependency detection (DFS algorithm)
│   │   ├── iteration_service.go     # Iteration lifecycle validation
│   │   └── *_service_test.go
│   ├── events/                      # Domain events
│   │   └── events.go                # 20+ event types (created/updated/status_changed)
│   ├── repositories/                # Repository interfaces (6 focused interfaces)
│   │   ├── roadmap_repository.go    # Roadmap CRUD
│   │   ├── track_repository.go      # Track CRUD + dependency management
│   │   ├── task_repository.go       # Task CRUD + iteration membership
│   │   ├── iteration_repository.go  # Iteration CRUD + task relationships
│   │   ├── adr_repository.go        # ADR CRUD + track association
│   │   └── acceptance_criteria_repository.go  # AC CRUD + verification status
│   └── CLAUDE.md                    # Domain layer guidance
│
├── application/                     # Use cases and orchestration
│   ├── track_service.go             # Track operations (CRUD + dependencies)
│   ├── task_service.go              # Task operations (CRUD + move)
│   ├── iteration_service.go         # Iteration operations (CRUD + lifecycle)
│   ├── adr_service.go               # ADR operations (CRUD + status transitions)
│   ├── ac_service.go                # AC operations (CRUD + verification)
│   ├── dto/                         # Data Transfer Objects
│   │   ├── track_dto.go             # TrackDTO, CreateTrackInput, UpdateTrackInput
│   │   ├── task_dto.go              # TaskDTO, CreateTaskInput, UpdateTaskInput
│   │   ├── iteration_dto.go         # IterationDTO, CreateIterationInput, etc.
│   │   ├── adr_dto.go               # ADRDTO, CreateADRInput, etc.
│   │   ├── ac_dto.go                # AcceptanceCriteriaDTO, CreateACInput, etc.
│   │   └── helpers.go               # Entity→DTO, DTO→Entity conversions
│   ├── mocks/                       # Generated mocks (mockery)
│   │   ├── mock_track_repository.go
│   │   ├── mock_task_repository.go
│   │   └── ... (6 repository mocks)
│   ├── *_service_test.go            # Service tests (126 tests, 82.1% coverage)
│   └── CLAUDE.md                    # Application layer guidance
│
├── infrastructure/                  # Technical implementations
│   └── persistence/                 # Database persistence
│       ├── roadmap_repository.go    # SQLite implementation
│       ├── track_repository.go      # SQLite implementation + dependency queries
│       ├── task_repository.go       # SQLite implementation + iteration joins
│       ├── iteration_repository.go  # SQLite implementation + task relationships
│       ├── adr_repository.go        # SQLite implementation + track foreign key
│       ├── acceptance_criteria_repository.go  # SQLite + verification status queries
│       ├── migrations.go            # Schema migrations (8 tables)
│       ├── event_emitting_repository.go  # Decorator for event emission
│       ├── repository_composite.go  # Composite pattern (legacy compatibility)
│       ├── *_repository_test.go     # Integration tests with real SQLite
│       └── CLAUDE.md                # Infrastructure layer guidance
│
├── presentation/                    # User interface layer
│   └── cli/                         # CLI command adapters
│       ├── track_adapters.go        # 7 track commands (create/list/show/update/delete/add-dep/remove-dep)
│       ├── task_adapters.go         # 7 task commands (create/list/show/update/delete/move/validate)
│       ├── iteration_adapters.go    # 10 iteration commands (create/list/show/current/update/start/complete/add-task/remove-task/delete)
│       ├── adr_adapters.go          # 7 ADR commands (create/list/show/update/supersede/deprecate/check)
│       ├── ac_adapters.go           # 9 AC commands (add/list/list-iteration/show/update/verify/fail/failed/delete)
│       ├── project_adapters.go      # 5 project commands (create/list/switch/show/delete)
│       ├── roadmap_adapters.go      # 3 roadmap commands (init/show/update)
│       └── CLAUDE.md                # Presentation layer guidance
│
├── e2e_test/                        # End-to-end tests
│   ├── e2e_test.go                  # Base suite (binary build, project setup)
│   ├── project_test.go              # Project command tests
│   ├── track_test.go                # Track command tests
│   ├── task_test.go                 # Task command tests
│   ├── iteration_test.go            # Iteration command tests
│   ├── adr_test.go                  # ADR command tests
│   ├── ac_test.go                   # Acceptance criteria tests
│   ├── workflow_test.go             # Complete workflow integration tests
│   └── CLAUDE.md                    # E2E test patterns and best practices
│
├── plugin.go                        # Plugin registration + dependency injection
└── CLAUDE.md                        # This file
```

**Total**: ~48 CLI commands, 7 aggregates, 8 database tables, 126 application tests

---

## Domain Model (Quick Reference)

**Roadmap** (Root Aggregate)
- Fields: ID, Vision, SuccessCriteria
- Purpose: Single root for all tracks per project
- Commands: `roadmap init/show/update`

**Track** (Work Stream)
- Fields: ID, Title, Description, Status (not-started/in-progress/complete/blocked/waiting), Priority (critical/high/medium/low), Rank
- Purpose: Major work areas with dependency management
- Key: Can depend on other tracks (circular dependency prevention via DFS)
- Commands: `track create/list/show/update/delete/add-dependency/remove-dependency`

**Task** (Atomic Work)
- Fields: ID, TrackID, Title, Description, Status (todo/in-progress/done), Rank, Branch
- Purpose: Concrete work items within tracks
- Key: Can belong to iterations, has acceptance criteria
- Commands: `task create/list/show/update/delete/move/validate`

**Iteration** (Time-Boxed Grouping)
- Fields: Number (auto-increment), Name, Goal, Deliverable, Status (planned/current/complete)
- Purpose: Group tasks from multiple tracks for time-boxed delivery
- Key: Only one "current" iteration at a time
- Commands: `iteration create/list/show/current/update/start/complete/add-task/remove-task/delete`

**ADR** (Architecture Decision Record)
- Fields: ID, TrackID, Title, Context, Decision, Consequences, Alternatives, Status (proposed/accepted/rejected/superseded/deprecated)
- Purpose: Document architectural decisions at track level
- Key: Immutable once accepted (create new ADR to change)
- Commands: `adr create/list/show/update/supersede/deprecate/check`

**AcceptanceCriteria** (Task Verification)
- Fields: ID, TaskID, Description, TestingInstructions, Status (not-started/pending-review/verified/failed), Feedback
- Purpose: Define "done" for tasks with verification steps
- Key: Must verify all ACs before task completion
- Commands: `ac add/list/list-iteration/show/update/verify/fail/failed/delete`

**Project** (Multi-Project Support)
- Purpose: Isolated SQLite databases per project (`.darwinflow/projects/<name>/roadmap.db`)
- Commands: `project create/list/switch/show/delete`

---

## Architecture Decisions (Why This Structure)

### 1. Unified Services (NOT CQRS)

**Pattern**: One service per aggregate (`TrackService`, `TaskService`, etc.) handling all operations.

**Why NOT CQRS**:
- Operations require orchestration (dependencies, validation, events)
- Read models would duplicate domain entities
- Queries need same business rules as commands
- Added complexity without benefit for this domain

**Key**: Service methods return DTOs, never domain entities (isolation between layers).

### 2. Repository Interface Segregation (6 Interfaces)

**Pattern**: One repository interface per aggregate (ISP compliance).

**Why NOT monolithic**:
- Consumers depend only on methods they use
- Clear ownership per aggregate
- Independent testing (mock only needed repositories)
- Prevents "god repository" anti-pattern

**Implementation**: 6 files in `infrastructure/persistence/*_repository.go`

### 3. Event Emission via Decorator

**Pattern**: `EventEmittingRepository` wraps base repositories.

```go
// plugin.go wiring
baseRepo := persistence.NewTrackRepository(db)
eventRepo := persistence.NewEventEmittingRepository(baseRepo, eventBus, "track")
trackService := application.NewTrackService(eventRepo, depService)
```

**Why decorator**:
- Base repositories stay pure (no event bus dependency)
- Event emission happens AFTER successful persistence
- Single responsibility (persistence vs notification)
- Application services unaware of event emission

### 4. DTO Conversion at Service Boundary

**Pattern**: Domain entities stay in domain/application. DTOs cross to presentation.

**Why**:
- Presentation doesn't import domain entities directly
- Service can change entity structure without breaking CLI
- DTOs are serialization-safe (no pointers, simplified types)

**Conversion**: `application/dto/helpers.go` contains all Entity↔DTO conversions.

### 5. Domain Service vs Entity Method

**Entity method** (default):
- Single entity validation/behavior
- Examples: `Track.Validate()`, `Task.CanComplete()`

**Domain service** (when needed):
- Multi-entity coordination
- Complex algorithms (DFS for circular dependencies)
- Stateless business logic
- Examples: `DependencyService.CheckCircular()`, `IterationService.CanStart()`

**Rule**: Start with entity method. Extract to domain service if needs multiple aggregates or complex algorithm.

### 6. Mock Placement: `application/mocks/`

**Location**: `application/mocks/`, NOT `domain/repositories/mocks/`

**Why**:
- Mocks are test infrastructure (not domain)
- Application tests are primary consumers
- Keeps domain/ free of test utilities
- `mockery` generates into application/mocks/ by convention

---

## Dependency Flow (Critical Rules)

```
Presentation (CLI)
    ↓ imports: application services + DTOs, domain types
Application (Services + DTOs)
    ↓ imports: domain interfaces + entities
Domain (Entities + Interfaces + Services)
    ↑ implemented by
Infrastructure (Repositories + Migrations)
```

**Layer Rules**:
1. **Domain**: Imports NOTHING (zero external dependencies)
2. **Application**: Imports domain ONLY
3. **Infrastructure**: Implements domain interfaces (dependency inversion)
4. **Presentation**: Uses application services (thin adapters, no business logic)

---

## Decision Guide: Where Does X Go?

### Adding Field to Entity

**Example**: Add `ArchivedAt *time.Time` to Track

1. **Domain entity** (`domain/entities/track_entity.go`): Add field + `Archive()` method
2. **DTO** (`application/dto/track_dto.go`): Add `ArchivedAt *time.Time` to `TrackDTO`
3. **Conversion** (`application/dto/helpers.go`): Update `ToTrackDTO()` and `FromTrackDTO()`
4. **Migration** (`infrastructure/persistence/migrations.go`): `ALTER TABLE tracks ADD COLUMN archived_at TIMESTAMP`
5. **Repository** (`infrastructure/persistence/track_repository.go`): Update `Save()` to persist field
6. **Service** (`application/track_service.go`): Add `ArchiveTrack(id string) error` method
7. **CLI** (`presentation/cli/track_adapters.go`): Add `track archive <id>` command
8. **E2E test** (`e2e_test/track_test.go`): Test archive command

### Adding Validation

**Entity-level validation** (`domain/entities/`):
- Field constraints (required, length, format)
- Single-entity invariants
- Example: `Track.Validate()` checks title not empty, max 200 chars

**Domain service validation** (`domain/services/`):
- Cross-entity constraints
- Complex business rules
- Example: `DependencyService.CheckCircular()`, `IterationService.CanStart()`

**Application service** (orchestration, NOT validation):
- Calls domain validation methods
- Never implements validation logic itself
- Returns validation errors from domain

### Adding CLI Command

**New command**: `dw task-manager track export <id> --format json`

1. **Application service** (`application/track_service.go`): Add `ExportTrack(id, format string)` if orchestration needed
2. **CLI adapter** (`presentation/cli/track_adapters.go`): Add command, parse flags, call service, format output
3. **E2E test** (`e2e_test/track_test.go`): Test command with various flags

**Rule**: CLI adapters are thin. Business logic goes in application/domain.

### Adding Query

**New query**: "Find all tracks blocked by track X"

1. **Repository interface** (`domain/repositories/track_repository.go`): Add `FindBlockedBy(trackID string) ([]*Track, error)`
2. **Repository implementation** (`infrastructure/persistence/track_repository.go`): Implement with SQL JOIN on `track_dependencies`
3. **Application service** (`application/track_service.go`): Add method if orchestration needed, or call repo directly
4. **CLI adapter** (`presentation/cli/track_adapters.go`): Add command or flag to existing command
5. **Tests**: Infrastructure test for SQL, application test for service orchestration

---

## Testing Strategy

### Domain Layer (`domain/*_test.go`)
- Pure unit tests (no mocks, no external dependencies)
- Test entity validation, state transitions, business rules
- Test domain services (DFS algorithm, lifecycle validation)
- **Never mock**: Domain has no external dependencies

### Application Layer (`application/*_service_test.go`)
- Mock repository interfaces (use `application/mocks/`)
- Verify orchestration (call order, multiple repos)
- Verify DTO conversion
- **Mock**: Repository interfaces
- **Don't mock**: Domain services (pure functions, no external dependencies)
- **Current coverage**: 82.1% (126 tests)

**Example pattern**:
```go
mockRepo := mocks.NewMockTrackRepository(t)
mockRepo.EXPECT().FindByID("TM-track-1").Return(track, nil)
mockRepo.EXPECT().Save(mock.Anything).Return(nil)

service := application.NewTrackService(mockRepo, depService)
dto, err := service.UpdateTrack("TM-track-1", updateInput)
```

### Infrastructure Layer (`infrastructure/persistence/*_test.go`)
- Real SQLite database in `t.TempDir()`
- Test migrations, CRUD, queries, constraints, referential integrity
- **Never mock**: Database (defeats purpose of integration test)

### Presentation Layer (`presentation/cli/*_test.go`)
- Mock application services
- Test flag parsing, error handling, output formatting
- **Never test**: Business logic (that's in application/domain)

### E2E Tests (`e2e_test/*_test.go`)
- Build binary from source: `go build -o /tmp/dw-e2e-test ./cmd/dw`
- Execute real commands, verify output
- Test complete workflows (create → update → delete)
- Test cross-entity operations (track → task → iteration → AC)
- **See**: `e2e_test/CLAUDE.md` for detailed patterns

---

## Common Anti-Patterns

### ❌ Domain Importing Infrastructure
```go
// domain/entities/track_entity.go
import "pkg/plugins/task_manager/infrastructure/persistence" // WRONG!
```
**Fix**: Define interface in `domain/repositories/`, implement in `infrastructure/persistence/`.

### ❌ Business Logic in Presentation
```go
// presentation/cli/track_adapters.go
if track.Status == "complete" && len(dependents) > 0 {
    return errors.New("cannot complete track with active dependents") // WRONG!
}
```
**Fix**: Move to `domain/services/dependency_service.go` or entity method.

### ❌ Application Service Returning Domain Entity
```go
// application/track_service.go
func (s *TrackService) GetTrack(id string) (*entities.Track, error) // WRONG!
```
**Fix**: Return DTO: `func (s *TrackService) GetTrack(id string) (*dto.TrackDTO, error)`

### ❌ CLI Adapter Calling Repository Directly
```go
// presentation/cli/track_adapters.go
track, err := trackRepo.FindByID(id) // WRONG!
```
**Fix**: Call application service: `trackDTO, err := trackService.GetTrack(id)`

### ❌ Mocking Domain Services
```go
// application/track_service_test.go
mockDepService := mocks.NewMockDependencyService(t) // WRONG!
```
**Fix**: Use real domain service (stateless, no external dependencies).

### ❌ SQL in Application Layer
```go
// application/track_service.go
rows, err := db.Query("SELECT * FROM tracks WHERE status = ?", status) // WRONG!
```
**Fix**: Add method to repository interface, implement in infrastructure.

---

## Dependency Injection (Plugin Wiring)

**Pattern** (`plugin.go`):
1. Create infrastructure (DB connection, repositories)
2. Wrap repositories with `EventEmittingRepository` (decorator)
3. Create domain services (stateless, pure functions)
4. Inject wrapped repositories + domain services into application services
5. Inject application services into CLI adapters
6. Register commands with framework

**Key**: Dependencies injected via constructors (NOT global variables or singletons).

**Example**:
```go
// Infrastructure
db := persistence.OpenDB(dbPath)
baseTrackRepo := persistence.NewTrackRepository(db)

// Decorator
eventTrackRepo := persistence.NewEventEmittingRepository(baseTrackRepo, eventBus, "track")

// Domain services
validationService := services.NewValidationService()
dependencyService := services.NewDependencyService(baseTrackRepo)

// Application services
trackService := application.NewTrackService(eventTrackRepo, dependencyService)

// Presentation
trackCommands := cli.NewTrackCommands(trackService)
```

---

## Multi-Project Architecture

**Isolation**: Each project → own SQLite DB (`.darwinflow/projects/<name>/roadmap.db`)

**Active project**: Tracked in `.darwinflow/active-project.txt`

**Commands**: All commands support `--project <name>` flag (overrides active project)

**Migration**: Auto-migrates from legacy single-database structure

**Use cases**:
- Separate "production" and "test" roadmaps
- Multiple product roadmaps in one workspace
- Experimentation without affecting real data

---

## Event Bus Integration

**20+ event types** (`domain/events/events.go`):
- Roadmap: `created`, `updated`
- Track: `created`, `updated`, `status_changed`, `completed`, `blocked`
- Task: `created`, `updated`, `status_changed`, `completed`, `moved`
- Iteration: `created`, `updated`, `started`, `completed`
- ADR: `created`, `updated`, `superseded`, `deprecated`
- AC: `created`, `verified`, `failed`

**Emission**: Via `EventEmittingRepository` decorator (after successful persistence)

**Subscription**: Other plugins can subscribe for notifications, automation, analytics

---

## Key References

- **Workflow**: `/workspace/CLAUDE.md` "Task Manager - Core Workflow" - How to use plugin (commands, best practices)
- **E2E Tests**: `e2e_test/CLAUDE.md` - E2E test patterns, best practices, examples
- **Domain Layer**: `domain/CLAUDE.md` - Domain-specific guidance
- **Application Layer**: `application/CLAUDE.md` - Application service patterns
- **Infrastructure Layer**: `infrastructure/CLAUDE.md` - Repository implementation guidance
- **Presentation Layer**: `presentation/CLAUDE.md` - CLI adapter patterns
- **SDK**: `pkg/pluginsdk/CLAUDE.md` - Plugin SDK documentation
- **Framework**: `/workspace/CLAUDE.md` - DarwinFlow architecture

---

**Last Updated**: 2025-11-14 (Clean architecture with DDD)
