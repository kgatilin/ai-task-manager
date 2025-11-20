# Package: task_manager

**Path**: `internal/task_manager`

**Role**: Core Task Manager application - hierarchical roadmap management (Roadmap → Track → Task → Iteration → Document → AcceptanceCriteria)

**Architecture**: Clean Architecture with Domain-Driven Design (Domain → Application → Infrastructure → Presentation)

**Database**: SQLite per-project (`~/.tm/projects/<name>/roadmap.db`)

---

## Package Structure

```
internal/task_manager/
├── domain/                          # Business logic (zero external dependencies)
│   ├── entities/                    # 7 core aggregates
│   │   ├── roadmap_entity.go        # Root aggregate (vision, success criteria)
│   │   ├── track_entity.go          # Work streams (status, priority, dependencies)
│   │   ├── task_entity.go           # Atomic work units (todo/in-progress/done)
│   │   ├── iteration_entity.go      # Time-boxed groupings (planned/active/complete)
│   │   ├── document_entity.go       # Documentation (ADR, plans, retrospectives)
│   │   ├── acceptance_criteria_entity.go  # Task verification (not-started/verified/failed)
│   │   ├── value_objects.go         # Status/Priority enums, ID types
│   │   └── *_entity_test.go
│   ├── services/                    # Domain services (stateless business logic)
│   │   ├── validation_service.go    # ID validation, format checks
│   │   ├── dependency_service.go    # Circular dependency detection (DFS algorithm)
│   │   ├── iteration_service.go     # Iteration lifecycle validation
│   │   └── *_service_test.go
│   ├── events/                      # Domain events (currently unused)
│   │   └── events.go                # Event type definitions for future use
│   ├── repositories/                # Repository interfaces (7 focused interfaces)
│   │   ├── roadmap_repository.go    # Roadmap CRUD
│   │   ├── track_repository.go      # Track CRUD + dependency management
│   │   ├── task_repository.go       # Task CRUD + iteration membership
│   │   ├── iteration_repository.go  # Iteration CRUD + task relationships
│   │   ├── document_repository.go   # Document CRUD + track/iteration association
│   │   ├── acceptance_criteria_repository.go  # AC CRUD + verification status
│   │   └── project_repository.go    # Project management (multi-project support)
│
├── application/                     # Use cases and orchestration
│   ├── track_service.go             # Track operations (CRUD + dependencies)
│   ├── task_service.go              # Task operations (CRUD + move)
│   ├── iteration_service.go         # Iteration operations (CRUD + lifecycle)
│   ├── document_service.go          # Document operations (CRUD + status transitions)
│   ├── ac_service.go                # AC operations (CRUD + verification)
│   ├── project_service.go           # Project operations (CRUD + switching)
│   ├── dto/                         # Data Transfer Objects
│   │   ├── track_dto.go             # TrackDTO, CreateTrackInput, UpdateTrackInput
│   │   ├── task_dto.go              # TaskDTO, CreateTaskInput, UpdateTaskInput
│   │   ├── iteration_dto.go         # IterationDTO, CreateIterationInput, etc.
│   │   ├── document_dto.go          # DocumentDTO, CreateDocumentInput, etc.
│   │   ├── ac_dto.go                # AcceptanceCriteriaDTO, CreateACInput, etc.
│   │   ├── project_dto.go           # ProjectDTO, CreateProjectInput, etc.
│   │   └── helpers.go               # Entity→DTO, DTO→Entity conversions
│   ├── mocks/                       # Generated mocks (mockery)
│   │   ├── mock_track_repository.go
│   │   ├── mock_task_repository.go
│   │   └── ... (7 repository mocks)
│   └── *_service_test.go            # Service tests (high coverage, blackbox)
│
├── infrastructure/                  # Technical implementations
│   └── persistence/                 # Database persistence
│       ├── roadmap_repository.go    # SQLite implementation
│       ├── track_repository.go      # SQLite implementation + dependency queries
│       ├── task_repository.go       # SQLite implementation + iteration joins
│       ├── iteration_repository.go  # SQLite implementation + task relationships
│       ├── document_repository.go   # SQLite implementation + track/iteration foreign keys
│       ├── acceptance_criteria_repository.go  # SQLite + verification status queries
│       ├── project_repository.go    # SQLite + active project tracking
│       ├── migrations.go            # Schema migrations (9 tables)
│       └── *_repository_test.go     # Integration tests with real SQLite
│
├── presentation/                    # User interface layer
│   └── cli/                         # CLI command adapters
│       ├── track_commands.go        # 7 track commands (create/list/show/update/delete/add-dep/remove-dep)
│       ├── task_commands.go         # 6 task commands (create/list/show/update/delete/move)
│       ├── iteration_commands.go    # 10 iteration commands (create/list/show/current/update/start/complete/add-task/remove-task/delete)
│       ├── document_commands.go     # 7 document commands (create/list/show/update/attach/detach/delete)
│       ├── ac_commands.go           # 9 AC commands (add/list/list-iteration/show/update/verify/fail/failed/delete)
│       ├── project_commands.go      # 5 project commands (create/list/switch/show/delete)
│       └── roadmap_commands.go      # 3 roadmap commands (init/show/update)
│
├── e2e_test/                        # End-to-end tests
│   ├── e2e_test.go                  # Base suite (binary build, project setup)
│   ├── project_test.go              # Project command tests
│   ├── track_test.go                # Track command tests
│   ├── task_test.go                 # Task command tests
│   ├── iteration_test.go            # Iteration command tests
│   ├── document_test.go             # Document command tests
│   ├── ac_test.go                   # Acceptance criteria tests
│   ├── workflow_test.go             # Complete workflow integration tests
│   └── CLAUDE.md                    # E2E test patterns and best practices
│
└── CLAUDE.md                        # This file
```

**Total**: ~47 CLI commands, 7 aggregates, 9 database tables

---

## Domain Model (Quick Reference)

**Roadmap** (Root Aggregate)
- Fields: ID, Vision, SuccessCriteria
- Purpose: Single root for all tracks per project
- Commands: `tm roadmap init/show/update`

**Track** (Work Stream)
- Fields: ID, Title, Description, Status (not-started/in-progress/complete/blocked/waiting), Priority (critical/high/medium/low), Rank
- Purpose: Major work areas with dependency management
- Key: Can depend on other tracks (circular dependency prevention via DFS)
- Commands: `tm track create/list/show/update/delete/add-dependency/remove-dependency`

**Task** (Atomic Work)
- Fields: ID, TrackID, Title, Description, Status (todo/in-progress/done), Rank, Branch
- Purpose: Concrete work items within tracks
- Key: Can belong to iterations, has acceptance criteria
- Commands: `tm task create/list/show/update/delete/move`

**Iteration** (Time-Boxed Grouping)
- Fields: Number (auto-increment), Name, Goal, Deliverable, Status (planned/active/complete), LockedAt
- Purpose: Group tasks from multiple tracks for time-boxed delivery
- Key: Only one "active" iteration at a time, locking prevents changes during execution
- Commands: `tm iteration create/list/show/current/update/start/complete/add-task/remove-task/delete`

**Document** (Documentation)
- Fields: ID, Title, Type (adr/plan/retrospective/guide/other), Content (markdown), Status, TrackID (optional), IterationNumber (optional)
- Purpose: Capture ADRs, planning docs, retrospectives, and other documentation
- Key: Can be attached to tracks or iterations, supports markdown content
- Commands: `tm doc create/list/show/update/attach/detach/delete`

**AcceptanceCriteria** (Task Verification)
- Fields: ID, TaskID, Description, TestingInstructions, Status (not-started/pending-review/verified/failed), Feedback
- Purpose: Define "done" for tasks with verification steps
- Key: Must verify all ACs before task completion
- Commands: `tm ac add/list/list-iteration/show/update/verify/fail/failed/delete`

**Project** (Multi-Project Support)
- Purpose: Isolated SQLite databases per project (`~/.tm/projects/<name>/roadmap.db`)
- Commands: `tm project create/list/switch/show/delete`

---

## Architecture Decisions (Why This Structure)

### 1. Clean Architecture + Domain-Driven Design

**Pattern**: The entire application follows Clean Architecture layering with DDD principles.

**Why Clean Architecture**:
- **Independence**: Business logic (domain) has zero external dependencies
- **Testability**: Each layer can be tested in isolation
- **Flexibility**: Can swap infrastructure (SQLite → Postgres) without touching domain
- **Clarity**: Clear separation of concerns, easy to navigate and understand

**Why Domain-Driven Design**:
- **Ubiquitous Language**: Entities and operations match user mental model
- **Aggregates**: Natural boundaries (Track, Task, Iteration) reflect business concepts
- **Domain Services**: Complex business logic (dependency cycles, iteration lifecycle) stays pure
- **Repository Pattern**: Clean abstraction over persistence

**Key Principle**: Business rules live in domain layer. Everything else is infrastructure.

### 2. Unified Services (NOT CQRS)

**Pattern**: One service per aggregate (`TrackService`, `TaskService`, etc.) handling all operations.

**Why NOT CQRS**:
- Operations require orchestration (dependencies, validation, lifecycle management)
- Read models would duplicate domain entities (no performance benefit for this scale)
- Queries need same business rules as commands (task status, iteration locking)
- Added complexity without benefit for this domain size

**Key**: Service methods return DTOs, never domain entities (isolation between layers).

### 3. Repository Interface Segregation (7 Interfaces)

**Pattern**: One repository interface per aggregate (ISP compliance).

**Why NOT monolithic**:
- Consumers depend only on methods they use
- Clear ownership per aggregate
- Independent testing (mock only needed repositories)
- Prevents "god repository" anti-pattern
- Easy to understand what operations are available per entity

**Implementation**: 7 files in `infrastructure/persistence/*_repository.go`

### 4. No Event Bus

**Why**:
- No external event consumers in current design
- Event emission would be infrastructure overhead without immediate benefit
- Domain events still defined in `domain/events/` for future use
- Simplifies dependency injection

**Future**: Can re-add event emission if needed (async workflows, audit logs, webhooks).

### 5. DTO Conversion at Service Boundary

**Pattern**: Domain entities stay in domain/application. DTOs cross to presentation.

**Why**:
- Presentation doesn't import domain entities directly
- Service can change entity structure without breaking CLI
- DTOs are serialization-safe (no pointers, simplified types)
- Clear contract between application and presentation layers

**Conversion**: `application/dto/helpers.go` contains all Entity↔DTO conversions.

### 6. Domain Service vs Entity Method

**Entity method** (default):
- Single entity validation/behavior
- Examples: `Track.Validate()`, `Task.CanComplete()`
- First choice for new behavior

**Domain service** (when needed):
- Multi-entity coordination
- Complex algorithms (DFS for circular dependencies)
- Stateless business logic that doesn't belong to one entity
- Examples: `DependencyService.CheckCircular()`, `IterationService.CanStart()`

**Rule**: Start with entity method. Extract to domain service if needs multiple aggregates or complex algorithm.

### 7. Mock Placement: `application/mocks/`

**Location**: `application/mocks/`, NOT `domain/repositories/mocks/`

**Why**:
- Mocks are test infrastructure (not domain concepts)
- Application tests are primary consumers
- Keeps domain/ free of test utilities
- `mockery` generates into application/mocks/ by convention

**Note**: Domain layer tests NEVER use mocks (domain has no external dependencies).

---

## Dependency Flow (Critical Rules)

```
Presentation (CLI)
    ↓ imports: application services + DTOs, domain types
Application (Services + DTOs)
    ↓ imports: domain interfaces + entities + services
Domain (Entities + Interfaces + Services)
    ↑ implemented by
Infrastructure (Repositories + Migrations)
```

**Layer Rules**:
1. **Domain**: Imports NOTHING (zero external dependencies, pure business logic)
2. **Application**: Imports domain ONLY (orchestration, no infrastructure)
3. **Infrastructure**: Implements domain interfaces (dependency inversion principle)
4. **Presentation**: Uses application services (thin adapters, no business logic)

**Critical**: These rules are enforced by `go-arch-lint`. Violations break the build.

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
7. **CLI** (`presentation/cli/track_commands.go`): Add `track archive <id>` command
8. **E2E test** (`e2e_test/track_test.go`): Test archive command

### Adding Validation

**Entity-level validation** (`domain/entities/`):
- Field constraints (required, length, format)
- Single-entity invariants
- Example: `Track.Validate()` checks title not empty, max 200 chars
- **When**: Validation only needs data from one entity

**Domain service validation** (`domain/services/`):
- Cross-entity constraints
- Complex business rules requiring multiple entities
- Stateless algorithms (no entity state mutation)
- Example: `DependencyService.CheckCircular()`, `IterationService.CanStart()`
- **When**: Validation needs data from multiple entities or complex graph algorithms

**Application service** (orchestration, NOT validation):
- Calls domain validation methods
- Never implements validation logic itself
- Returns validation errors from domain
- **Role**: Coordinate calls to domain, don't duplicate domain logic

**Anti-pattern**: Application service implementing business rules (belongs in domain).

### Adding CLI Command

**New command**: `tm track export <id> --format json`

1. **Application service** (`application/track_service.go`): Add `ExportTrack(id, format string)` if orchestration needed
2. **CLI command** (`presentation/cli/track_commands.go`): Add command, parse flags, call service, format output
3. **E2E test** (`e2e_test/track_test.go`): Test command with various flags

**Rule**: CLI commands are thin adapters. Business logic goes in application/domain.

**What goes in CLI**:
- Flag parsing (`--format json`)
- Calling application service
- Output formatting (tables, JSON, etc.)
- Error message presentation

**What goes in application/domain**:
- Business logic (what fields to export)
- Validation (is format valid?)
- Data retrieval and transformation

### Adding Query

**New query**: "Find all tracks blocked by track X"

1. **Repository interface** (`domain/repositories/track_repository.go`): Add `FindBlockedBy(trackID string) ([]*Track, error)`
2. **Repository implementation** (`infrastructure/persistence/track_repository.go`): Implement with SQL JOIN on `track_dependencies`
3. **Application service** (`application/track_service.go`): Add method if orchestration needed, or call repo directly
4. **CLI command** (`presentation/cli/track_commands.go`): Add command or flag to existing command
5. **Tests**: Infrastructure test for SQL correctness, application test for service orchestration if needed

**When to add to repository**:
- Query is reusable across multiple commands
- Query is complex (joins, aggregations, filtering)
- Query is a domain concept ("blocked by" is a dependency relationship)

**When to query directly in service**:
- One-off query for specific use case
- Simple query (single table, no joins)
- BUT: Still prefer repository for testability

### Adding New Aggregate

**Example**: Add "Milestone" aggregate (group of iterations)

1. **Domain entity** (`domain/entities/milestone_entity.go`): Define entity, value objects, methods
2. **Repository interface** (`domain/repositories/milestone_repository.go`): Define CRUD + queries
3. **DTO** (`application/dto/milestone_dto.go`): Define DTOs and conversion helpers
4. **Application service** (`application/milestone_service.go`): Define use cases
5. **Migration** (`infrastructure/persistence/migrations.go`): Add `milestones` table
6. **Repository impl** (`infrastructure/persistence/milestone_repository.go`): Implement interface
7. **CLI commands** (`presentation/cli/milestone_commands.go`): Create/list/show/update/delete
8. **Tests**: Domain tests, application tests (mocked repo), infrastructure tests (real DB), E2E tests
9. **Wire up** (`cmd/tm/main.go`): Create repository, service, commands in dependency injection

**Consider**:
- Is this truly a new aggregate, or should it be a field on existing aggregate?
- Does it have its own lifecycle independent of other aggregates?
- Does it have clear boundaries and invariants?

---

## Testing Strategy

### Domain Layer (`domain/*_test.go`)

**What to test**:
- Pure unit tests (no mocks, no external dependencies)
- Entity validation, state transitions, business rules
- Domain services (DFS algorithm, lifecycle validation, dependency checks)

**How to test**:
- Package: `package entities_test` (blackbox testing)
- Create entities directly, call methods, assert results
- **Never mock**: Domain has no external dependencies (nothing to mock)

**Example**:
```go
func TestTrack_Validate(t *testing.T) {
    track := entities.NewTrack("", "TM-track-1", "Test Track", "Description", ...)
    err := track.Validate()
    assert.NoError(t, err)
}
```

### Application Layer (`application/*_service_test.go`)

**What to test**:
- Service orchestration (correct repository calls, correct order)
- DTO conversion (entity → DTO, DTO → entity)
- Error handling and validation coordination
- Multi-repository transactions

**How to test**:
- Package: `package application_test` (blackbox testing)
- Mock repository interfaces (use `application/mocks/`)
- Verify orchestration logic, not business logic (that's in domain)

**Mock**: Repository interfaces only
**Don't mock**: Domain services (pure functions, use real implementations)

**Example**:
```go
mockRepo := mocks.NewMockTrackRepository(t)
mockRepo.EXPECT().FindByID("TM-track-1").Return(track, nil)
mockRepo.EXPECT().Save(mock.Anything).Return(nil)

service := application.NewTrackService(mockRepo, depService)
dto, err := service.UpdateTrack("TM-track-1", updateInput)
```

### Infrastructure Layer (`infrastructure/persistence/*_test.go`)

**What to test**:
- Real SQLite database operations (no mocks)
- Migrations, CRUD, complex queries, constraints
- Referential integrity, indexes, query performance
- SQL correctness, edge cases (NULLs, empty strings, etc.)

**How to test**:
- Package: `package persistence_test` (blackbox testing)
- Use `t.TempDir()` for isolated SQLite database
- Run migrations, execute operations, verify results
- **Never mock**: Database (defeats purpose of integration test)

**Example**:
```go
db := setupTestDB(t) // t.TempDir() + migrations
repo := persistence.NewTrackRepository(db)

track := entities.NewTrack(...)
err := repo.Save(track)
assert.NoError(t, err)

found, err := repo.FindByID(track.ID())
assert.Equal(t, track.Title(), found.Title())
```

### Presentation Layer (`presentation/cli/*_test.go`)

**What to test**:
- Flag parsing, error handling, output formatting
- Argument validation (user input errors)
- Service integration (correct service method called)

**How to test**:
- Mock application services
- Test input parsing and output formatting
- **Never test**: Business logic (that's in application/domain)

**Note**: Currently minimal presentation tests (E2E tests cover CLI comprehensively).

### E2E Tests (`e2e_test/*_test.go`)

**What to test**:
- Build binary from source: `go build -o /tmp/tm-e2e-test ./cmd/tm`
- Execute real commands, verify output (stdout, stderr, exit codes)
- Test complete workflows (create → update → delete)
- Test cross-entity operations (track → task → iteration → AC)
- Test multi-project scenarios
- Test error cases (invalid input, missing data, constraint violations)

**How to test**:
- Use `exec.Command()` to run binary in isolated environment
- Use `t.TempDir()` for isolated home directory
- Parse command output (JSON, tables, error messages)
- Chain commands to test workflows

**Coverage**: E2E tests provide comprehensive integration testing across all layers.

**See**: `e2e_test/CLAUDE.md` for detailed patterns, best practices, and examples.

---

## Common Anti-Patterns

### ❌ Domain Importing Infrastructure
```go
// domain/entities/track_entity.go
import "internal/task_manager/infrastructure/persistence" // WRONG!
```
**Fix**: Define interface in `domain/repositories/`, implement in `infrastructure/persistence/`.

**Why wrong**: Domain must be independent of infrastructure (core Clean Architecture principle).

### ❌ Business Logic in Presentation
```go
// presentation/cli/track_commands.go
if track.Status == "complete" && len(dependents) > 0 {
    return errors.New("cannot complete track with active dependents") // WRONG!
}
```
**Fix**: Move to `domain/services/dependency_service.go` or entity method.

**Why wrong**: Presentation is for I/O, not business rules. Rules belong in domain.

### ❌ Application Service Returning Domain Entity
```go
// application/track_service.go
func (s *TrackService) GetTrack(id string) (*entities.Track, error) // WRONG!
```
**Fix**: Return DTO: `func (s *TrackService) GetTrack(id string) (*dto.TrackDTO, error)`

**Why wrong**: Breaks layer isolation. Presentation should never see domain entities.

### ❌ CLI Command Calling Repository Directly
```go
// presentation/cli/track_commands.go
track, err := trackRepo.FindByID(id) // WRONG!
```
**Fix**: Call application service: `trackDTO, err := trackService.GetTrack(id)`

**Why wrong**: Bypasses application orchestration. Presentation → Application → Domain → Infrastructure.

### ❌ Mocking Domain Services
```go
// application/track_service_test.go
mockDepService := mocks.NewMockDependencyService(t) // WRONG!
```
**Fix**: Use real domain service (stateless, no external dependencies).

**Why wrong**: Domain services are pure functions. No reason to mock them.

### ❌ SQL in Application Layer
```go
// application/track_service.go
rows, err := db.Query("SELECT * FROM tracks WHERE status = ?", status) // WRONG!
```
**Fix**: Add method to repository interface, implement in infrastructure.

**Why wrong**: Application shouldn't know about SQL. That's infrastructure concern.

### ❌ Business Logic in DTO Conversion
```go
// application/dto/helpers.go
func ToTrackDTO(track *entities.Track) *TrackDTO {
    if track.Status() == "complete" && len(track.Dependencies()) > 0 {
        // complex business logic here... // WRONG!
    }
}
```
**Fix**: Keep conversions simple. Business logic belongs in domain/application service.

**Why wrong**: DTOs are data structures. Logic belongs in domain or service methods.

### ❌ Repository Implementing Business Logic
```go
// infrastructure/persistence/track_repository.go
func (r *TrackRepository) Save(track *entities.Track) error {
    if track.Status() == "complete" {
        // validate dependents, update related tracks, etc. // WRONG!
    }
}
```
**Fix**: Keep repository to persistence only. Business logic in domain service or entity.

**Why wrong**: Infrastructure shouldn't contain business rules. That's domain responsibility.

---

## Dependency Injection (Application Wiring)

**Pattern** (`cmd/tm/main.go`):
1. Create infrastructure (DB connection, file system access)
2. Create repositories (implement domain interfaces)
3. Create domain services (stateless, pure functions)
4. Inject repositories + domain services into application services
5. Inject application services into CLI commands
6. Wire up commands to CLI framework

**Key**: Dependencies injected via constructors (NOT global variables or singletons).

**Example** (simplified):
```go
// Infrastructure
db := persistence.OpenDB(dbPath)
trackRepo := persistence.NewTrackRepository(db)
taskRepo := persistence.NewTaskRepository(db)

// Domain services
validationService := services.NewValidationService()
dependencyService := services.NewDependencyService(trackRepo)

// Application services
trackService := application.NewTrackService(trackRepo, dependencyService)
taskService := application.NewTaskService(taskRepo, trackRepo)

// Presentation
rootCmd.AddCommand(
    cli.NewTrackCommands(trackService),
    cli.NewTaskCommands(taskService),
)
```

**Benefits**:
- Easy to test (inject mocks)
- Easy to swap implementations (SQLite → Postgres)
- Clear dependency graph (no hidden dependencies)
- No global state (testable, thread-safe)

---

## Multi-Project Architecture

**Isolation**: Each project → own SQLite DB (`~/.tm/projects/<name>/roadmap.db`)

**Active project**: Tracked in `~/.tm/active-project.txt`

**Commands**: All commands support `--project <name>` flag (overrides active project)

**Use cases**:
- Separate "production" and "test" roadmaps
- Multiple product roadmaps in one workspace
- Experimentation without affecting real data
- Team collaboration (different team members, different projects)

**Database per project**:
- Complete isolation (no cross-project contamination)
- Easy backup/restore (copy single DB file)
- Easy sharing (send DB file to teammate)
- Independent migrations (can test schema changes in test project)

---

## Key Principles (Summary)

1. **Clean Architecture**: Business logic (domain) has zero external dependencies
2. **Dependency Inversion**: Domain defines interfaces, infrastructure implements them
3. **Single Responsibility**: Each layer has one reason to change
4. **Interface Segregation**: Small, focused repository interfaces (not monolithic)
5. **DTO Boundary**: Domain entities never cross to presentation layer
6. **Domain Services**: Multi-entity coordination, complex algorithms, stateless logic
7. **No Mocking Domain**: Domain services are pure functions (use real implementations)
8. **Infrastructure Tests**: Use real database (no mocks in integration tests)
9. **E2E Coverage**: Comprehensive workflow testing via binary execution
10. **Dependency Injection**: Constructor injection (no globals, no singletons)

---

## Key References

- **Workflow**: Root `CLAUDE.md` "Task Manager - Core Workflow" - How to use the application (commands, best practices)
- **Architecture Index**: `docs/arch-index.md` - Full package structure and dependencies
- **E2E Tests**: `e2e_test/CLAUDE.md` - E2E test patterns, best practices, examples
- **Entry Point**: `cmd/tm/main.go` - Application bootstrap and dependency injection

---

**Last Updated**: 2025-11-20 (Standalone application, Clean Architecture + DDD)
