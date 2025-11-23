# Task Manager (tm)

## CRITICAL: Path Convention

**NEVER use absolute paths** in documentation, commands, or code references.

- ✅ **Always use relative paths from project root**: `internal/task_manager/domain/task.go`
- ✅ **Project root is your frame of reference**: All paths start from repository root
- ❌ **Never use `/workspace/`**: Different environments use different absolute paths
- ❌ **Never use `/Users/...` or `/home/...`**: Machine-specific paths break portability

**Why**: Project root may be `/workspace/` in container, `/Users/name/projects/task-manager/` on Mac, `/home/user/tm/` on Linux, etc. Only relative paths are portable.

**Examples**:
- Good: `internal/task_manager/CLAUDE.md`
- Good: `.agent/iteration-38-plan-2025-11-22.md`
- Good: `cmd/tm/main.go`
- Bad: `/workspace/internal/task_manager/CLAUDE.md`
- Bad: `/Users/kgatilin/PersonalProjects/darwinflow-pub/internal/...`

---

## Project Overview

**Task Manager** is a standalone CLI tool for managing software development workflows using Clean Architecture and Domain-Driven Design principles. It provides roadmap tracking, task management, iterations, and acceptance criteria verification.

### Core Features

- **Roadmap Management**: Organize work into tracks, tasks, and iterations
- **Acceptance Criteria**: Define and verify task completion criteria
- **Iteration Workflow**: Plan and execute work in focused iterations
- **Document Management**: ADRs, plans, and retrospectives attached to tracks/iterations
- **SQLite Storage**: Per-project database in `.taskmanager/data.db`
- **Clean Architecture**: Strict layer separation enforced by go-arch-lint
- **Cobra CLI**: Robust command-line interface with subcommands

---

## Task Manager - Core Workflow

**Understanding What to Work On:**

```bash
# CRITICAL: Always check this first (shows current or next planned iteration)
tm iteration current

# View all tasks (with status filtering)
tm task list                    # All tasks
tm task list --status todo      # Backlog
tm task list --status in-progress  # Active work

# View track details and tasks
tm track show TM-track-X

# View task details (including acceptance criteria)
tm task show TM-task-X
```

**Working on Tasks:**

```bash
# Start a task (todo → in-progress)
tm task update TM-task-X --status in-progress

# Complete a task (USER ONLY - agent does NOT mark done unless asked)
tm task update TM-task-X --status done

# Return to backlog (in-progress → todo)
tm task update TM-task-X --status todo
```

**Acceptance Criteria:**

```bash
# Add acceptance criterion to task
tm ac add TM-task-X --description "..."

# List task acceptance criteria
tm ac list TM-task-X

# Mark as verified (USER ONLY - agent does NOT verify unless asked)
tm ac verify TM-ac-X

# Mark as failed with feedback (USER ONLY)
tm ac fail TM-ac-X --feedback "..."

# List failed ACs in current iteration (most useful)
tm ac failed --iteration <current-iteration-num>

# Or filter by task/track
tm ac failed --task TM-task-X
tm ac failed --track TM-track-X
```

**Creating New Work:**

```bash
# 1. Create a new track
tm track create --title "..." --description "..." --rank 100

# 2. (Optional) Create ADR document for the track
tm doc create \
  --title "ADR: ..." \
  --type adr \
  --content "# Context\n...\n\n# Decision\n...\n\n# Consequences\n..." \
  --track <track-id>

# Or from file
tm doc create \
  --title "ADR: ..." \
  --type adr \
  --from-file ./docs/adr.md \
  --track <track-id>

# 3. Create tasks in the track with acceptance criteria
tm task create --track TM-track-X --title "..." --rank 100
tm ac add TM-task-X --description "..."

# 4. Create iteration and add tasks
tm iteration create --name "..." --goal "..." --deliverable "..."
tm iteration add-task <iter-num> TM-task-1 TM-task-2

# 5. Start working on iteration
tm iteration start <iter-num>
```

**Document Management Commands:**

```bash
# Create document (ADR, plan, retrospective, etc.)
tm doc create \
  --title "..." \
  --type adr \
  --from-file ./docs/adr.md \
  --track TM-track-X

# Or create inline
tm doc create \
  --title "..." \
  --type plan \
  --content "# Planning doc..."

# List documents (filter by type)
tm doc list
tm doc list --type adr

# Show document
tm doc show TM-doc-X

# Update document
tm doc update TM-doc-X --from-file ./updated.md

# Attach to track or iteration
tm doc attach TM-doc-X --track TM-track-Y
tm doc attach TM-doc-X --iteration 5

# Detach document
tm doc detach TM-doc-X

# Delete document
tm doc delete TM-doc-X [--force]
```

**Priority Guidance**:
- **CRITICAL**: Always run `tm iteration current` first to see what to work on
- Shows current active iteration OR next planned iteration if none active
- Iteration is the primary working entity (tracks are just grouping)

**Best Practices**:
- Update task status as you work (don't batch updates)
- **Agent responsibility**: Implement work defined in tasks
- **User responsibility**: Verify all acceptance criteria and mark tasks "done"
- Agent does NOT verify AC or mark tasks done unless explicitly asked
- Use `tm iteration current` to stay focused
- Use documents (ADRs, plans, retrospectives) for architecture decisions and planning

### Writing Good Acceptance Criteria

**Core Principle**: Acceptance criteria must describe **end-user verifiable functionality** that focuses on **core business logic**, not implementation details or edge cases.

**Command Structure**:
```bash
tm ac add <task-id> \
  --description "What must be verified (end-user observable)" \
  --testing-instructions "Step-by-step instructions to verify"
```

**CRITICAL**: Use separate fields:
- `--description`: The acceptance criterion itself (WHAT needs to be verified)
- `--testing-instructions`: Step-by-step instructions (HOW to verify it)

**Good AC Characteristics**:
- ✅ Describes WHAT the user can verify, not HOW it's implemented
- ✅ Focuses on observable behavior and outcomes
- ✅ Can be tested/verified by an end user
- ✅ Addresses core business logic
- ✅ Written from user perspective
- ✅ Testing instructions in separate field with numbered steps

**Bad AC Characteristics**:
- ❌ Implementation details (repositories, services, internal methods)
- ❌ Edge cases and technical minutiae
- ❌ Things only developers care about
- ❌ Internal code structure or architecture
- ❌ Testing instructions mixed into description field

**Examples**:

Good AC with proper separation:
```bash
tm ac add TM-task-X \
  --description "Domain layer has 90%+ test coverage with all tests passing" \
  --testing-instructions "1. Run: cd internal/task_manager/domain
2. Run: go test ./... -coverprofile=coverage.out
3. Run: go tool cover -func=coverage.out | grep total
4. Verify: total coverage >= 90%
5. Run: go test ./... -v
6. Verify: All tests pass with zero failures"
```

Bad AC (everything in description):
```bash
tm ac add TM-task-X \
  --description "Domain layer has 90%+ test coverage

Testing instructions:
1. Run: go test ./...
2. Verify coverage >= 90%"
```

**Testing Instructions Best Practices**:
- Start each step with a number
- Use exact commands (copy-paste ready)
- Include verification steps ("Verify: X should show Y")
- Make it reproducible by anyone
- Focus on observable outcomes, not internal state

### Task Granularity

**Core Principle**: If you cannot write end-user verifiable acceptance criteria for a task, the task is likely **too granular** and should be merged into a larger, user-facing task.

**Good Task Granularity**:
- ✅ Represents a complete user-facing feature or capability
- ✅ Has at least 3-5 end-user verifiable acceptance criteria
- ✅ Delivers observable value to the end user
- ✅ Can be demonstrated and tested independently

**Too Granular (merge into larger task)**:
- ❌ "Create X entity" - implementation detail, merge into command that uses it
- ❌ "Add database migration" - implementation detail, happens as part of feature
- ❌ "Define X interface" - implementation detail, merge into service that implements it
- ❌ "Add X field to entity" - implementation detail, merge into feature that uses it

**Examples**:

Good Tasks:
- ✅ "Add task validation and comment commands" (includes entity creation, repository, commands)
- ✅ "Implement iteration locking workflow" (includes status fields, entities, commands)
- ✅ "Show iteration membership in task details" (includes repository method, CLI, TUI)

Too Granular (should be merged):
- ❌ "Create TaskComment entity" → Merge into "Add task validation commands"
- ❌ "Database migration for task planning" → Merge into "Add task planning commands"
- ❌ "Add iteration status fields" → Merge into "Add iteration locking commands"

**Guidelines**:
- Tasks should represent **features**, not implementation steps
- Implementation details (entities, migrations, repositories) are part of feature delivery
- Each task should answer: "What can the user now do that they couldn't before?"
- If the answer is "nothing visible", the task is too granular

---

## Architecture Overview

**Task Manager** follows Clean Architecture principles with strict layer separation:

```
Presentation → Application → Domain ← Infrastructure
     ↓              ↓           ↑            ↑
   (CLI)      (Use Cases)  (Entities)   (Database)
```

### Why Clean Architecture?

1. **Testability**: Core business logic isolated and easily testable
2. **Independence**: Business rules don't depend on frameworks, UI, or database
3. **Flexibility**: Easy to swap implementations (CLI → Web, SQLite → Postgres)
4. **Maintainability**: Clear boundaries reduce coupling and increase cohesion

### Why Cobra?

**Cobra** was chosen for the CLI framework because:
- Industry-standard (used by kubectl, docker, hugo)
- Excellent command organization and subcommand support
- Built-in help generation
- Flag parsing with validation
- Minimal boilerplate
- Active maintenance and community support

### Why SQLite Per-Project?

**Per-project SQLite databases** (`.taskmanager/data.db`) because:
- Zero server setup or configuration
- Fast local queries (no network latency)
- Easy backup (just copy the file)
- Git-friendly (can be committed or gitignored)
- Each project has isolated task data
- No shared state between projects

---

## Package Structure

All code lives under `internal/task_manager/` following Clean Architecture layers:

**Domain Layer** (`internal/task_manager/domain/`):
- Pure business logic with zero external dependencies
- Entities: Track, Task, Iteration, AcceptanceCriteria, Document
- Repository interfaces (contracts only)
- Domain services and aggregates
- Value objects and domain events

**Application Layer** (`internal/task_manager/application/`):
- Use cases and application services
- Orchestrates domain logic and infrastructure
- Transaction coordination
- Business workflow implementation
- Maps domain entities to DTOs (if needed)

**Infrastructure Layer** (`internal/task_manager/infrastructure/`):
- Database implementations (SQLite repositories)
- File system operations
- External service integrations
- Repository implementations for domain interfaces

**Presentation Layer** (`internal/task_manager/presentation/`):
- CLI commands (Cobra)
- Command handlers
- Input validation and parsing
- Output formatting
- User interaction

**Entry Point** (`cmd/tm/`):
- Main entry point
- Dependency injection and wiring
- Configuration loading
- Application bootstrap

**E2E Tests** (`internal/task_manager/e2e_test/`):
- End-to-end integration tests
- Full workflow validation
- Real database and file system tests

### Layer Dependencies (Enforced by go-arch-lint)

```
Presentation → Application → Domain
Infrastructure → Domain

NEVER:
Domain → Application
Domain → Infrastructure
Domain → Presentation
```

**Key Rules**:
- Domain layer imports NOTHING from other layers
- Application imports Domain only
- Infrastructure imports Domain only (implements interfaces)
- Presentation imports Application and Domain (orchestrates)
- Entry point imports all layers (wires dependencies)

---

## Architecture Quick Reference

### Dependency Rules

- **internal/task_manager/domain**: Imports NOTHING from other task_manager packages
- **internal/task_manager/application**: May import `domain` only
- **internal/task_manager/infrastructure**: May import `domain` only
- **internal/task_manager/presentation**: May import `domain`, `application`
- **cmd/tm**: May import all layers for dependency injection

**Key Principle**: Dependencies flow inward toward domain. Domain is the most stable layer.

### Core Principles

1. **Dependency Inversion**: Define interfaces in domain, implement in infrastructure
2. **Separation of Concerns**: Each layer has a single, well-defined responsibility
3. **Domain-Centric**: Business logic is isolated and protected
4. **Repository Pattern**: One repository per aggregate root
5. **Use Cases**: Application services represent user actions/workflows
6. **Clean Boundaries**: No leaking of infrastructure details into domain

### Repository Pattern

Each aggregate root has a repository:
- **TrackRepository**: Manages Track aggregates
- **TaskRepository**: Manages Task aggregates
- **IterationRepository**: Manages Iteration aggregates
- **AcceptanceCriteriaRepository**: Manages AcceptanceCriteria
- **DocumentRepository**: Manages Document entities

**Pattern**:
```go
// Domain defines interface
type TaskRepository interface {
    Save(task *Task) error
    FindByID(id string) (*Task, error)
    // ... other methods
}

// Infrastructure implements
type sqliteTaskRepository struct {
    db *sql.DB
}

func (r *sqliteTaskRepository) Save(task *Task) error {
    // SQLite implementation
}

// Application receives injected repository
type TaskService struct {
    repo domain.TaskRepository
}
```

---

## Development Workflow

**Note**: When the user refers to "workflow", they mean these CLAUDE.md instructions.

### Building and Running

```bash
# Build binary to ./tm
make build

# Install to GOPATH/bin
make install

# Run tests
make test

# Or directly
go build -o tm ./cmd/tm
go test ./...
```

### Working on Features

1. Understand which layer the change belongs to:
   - New business rule? → Domain
   - New workflow? → Application
   - New database query? → Infrastructure
   - New command? → Presentation

2. Read relevant package `CLAUDE.md` for layer-specific guidance:
   - `internal/task_manager/CLAUDE.md` - Overall architecture
   - `internal/task_manager/domain/CLAUDE.md` - Domain patterns
   - `internal/task_manager/application/CLAUDE.md` - Use case patterns
   - etc.

3. Follow Clean Architecture rules:
   - Domain never imports from other layers
   - Define interfaces in domain, implement in infrastructure
   - Use dependency injection (constructor injection)

4. Write tests for new functionality (target 70-80% coverage)

5. Update documentation when adding features

6. Run tests and linter before committing

7. Commit after each logical task/iteration

### Large Tasks - Use Task Tool Delegation

For substantial refactorings or multi-package features:

1. **Decompose** into layer-sized chunks (Domain → Application → Infrastructure → Presentation)
2. **Delegate** each chunk sequentially using Task tool
3. **Review** sub-agent reports between chunks
4. **Verify** all tests/linter pass after completion

**Final Checklist** (use TodoWrite):
- [ ] Run `go test ./...` - all tests pass
- [ ] Run `go-arch-lint .` - zero violations
- [ ] Update README.md (if commands/features changed)
- [ ] Update CLAUDE.md (if workflow/architecture changed)
- [ ] Commit with concise message

### Reporting Completed Work

**CRITICAL**: When reporting work completion to the user, follow these guidelines:

**DO**:
- ✅ Highlight deviations from the original plan
- ✅ Emphasize questions that require user decision
- ✅ Call out issues that need user attention
- ✅ Note any blockers or unexpected challenges
- ✅ Use clear, separate blocks for user action items

**DON'T**:
- ❌ Provide detailed summaries of what was implemented (user knows the tasks)
- ❌ Include standard "verify acceptance criteria" instructions (goes without saying)
- ❌ Repeat task descriptions or requirements back to user
- ❌ List every file changed or every test added (user can see git commit)
- ❌ Explain what was supposed to happen (user defined the tasks)

**Report Format**:

```markdown
# Implementation Complete: [Task/Iteration Name]

## Status
[One line: Complete / Complete with deviations / Blocked]

## Deviations from Plan
[Only if there were deviations - explain what and why]

## Questions for User
[Only if decisions needed - clear, actionable questions]

## Issues Requiring Attention
[Only if there are blockers or problems]

## Commit
[Commit hash and one-line summary]
```

**Example - Good Report**:
```markdown
# Implementation Complete: Iteration #27 TODO Tasks

## Status
Complete with minor deviations

## Deviations from Plan
1. Mocks placed in `application/mocks/` instead of `domain/repositories/mocks/` (AC-482 specifies domain layer)
2. Infrastructure coverage at 52.3% vs 60% target (7.7% gap)

## Questions for User
1. Should mocks stay in application/ or move to domain/repositories/? (affects linter violations)
2. Is 52.3% infrastructure coverage acceptable, or should I add more tests?

## Commit
8b103f8 feat: complete iteration #27 TODO tasks - test architecture refactoring
```

**Example - Bad Report**:
```markdown
# Implementation Complete: Iteration #27 TODO Tasks

## Summary
Successfully implemented all 4 TODO tasks...
[3 paragraphs explaining what was done]

## Implementation Phases
Phase 1: Created mocks...
Phase 2: Refactored tests...
[Detailed breakdown of every step]

## Files Modified
- Created: application/mocks/ (6 files)
- Modified: 10 files
[Long list of every file]

## Next Steps - USER ACTION REQUIRED
1. Review All Acceptance Criteria
2. Verify Each AC
3. Close Tasks
[Standard AC verification instructions]
```

---

## Before Every Commit

1. `go test ./...` - all tests must pass
2. `go-arch-lint .` - **ZERO violations required** (non-negotiable)
3. Update README.md / CLAUDE.md if functionality changed

---

## When Linter Reports Violations

**Don't mechanically fix imports.** Violations reveal architectural issues.

**Process**:
1. **Reflect**: Why does this violation exist? Wrong dependency?
2. **Plan**: Which layer should own this? Right structure?
3. **Refactor**: Move code to correct layer
4. **Verify**: `go-arch-lint .` → zero violations

**Common Violations**:

❌ **Domain imports Application/Infrastructure/Presentation**
→ Fix: Domain must be pure. Move code to appropriate layer or use interfaces.

❌ **Application imports Infrastructure**
→ Fix: Application should depend on domain interfaces only. Infrastructure implements interfaces, gets injected.

❌ **Circular dependencies between layers**
→ Fix: Review layer responsibilities. Use dependency inversion.

✅ **Application imports Domain** - OK (expected)
✅ **Infrastructure imports Domain** - OK (implements interfaces)
✅ **Presentation imports Application + Domain** - OK (orchestrates)

**Example - Application needs database**:
- ✅ `domain/repositories.go` defines `TaskRepository` interface
- ✅ `infrastructure/sqlite_task_repository.go` implements `TaskRepository`
- ✅ `application/task_service.go` receives injected `domain.TaskRepository`
- ✅ `cmd/tm/main.go` wires concrete implementation

---

## Testing

**Coverage Target**: 70-80%

**Package Naming**: `package pkgname_test` (black-box testing preferred)

**File Naming**: `*_test.go` in same directory

**Test Naming**:
- `TestFunctionName` or `TestType_Method`
- Examples: `TestNewTaskService`, `TestTask_Validate`

**Running Tests**:
```bash
go test ./...                               # All tests
go test -cover ./...                        # With coverage
go test -coverprofile=coverage.out ./...    # Coverage report
go tool cover -html=coverage.out            # View in browser

# Layer-specific tests
go test ./internal/task_manager/domain/...
go test ./internal/task_manager/application/...
go test ./internal/task_manager/infrastructure/...
```

**Best Practices**:
- Each test is independent (no shared state)
- Use `t.TempDir()` for file/database operations
- Use `defer` for cleanup
- Test public API only (black-box) when possible
- `t.Fatalf()` for setup failures, `t.Errorf()` for assertions
- Mock repositories for application layer tests
- Use real database for infrastructure tests (with temp DB)

**Test Organization**:
- **Unit tests**: Domain and application layer (mock dependencies)
- **Integration tests**: Infrastructure layer (real database)
- **E2E tests**: `e2e_test/` package (full workflow validation)

---

## Documentation

**CRITICAL**: Documentation must be updated when functionality changes.

### Documentation Types

**README.md** (user-facing):
- New commands or flags
- New features
- Changed behavior
- Installation instructions
- Quick start guide

**CLAUDE.md** (this file - development workflow):
- Workflow changes
- Architecture changes
- New patterns or conventions
- Layer responsibilities
- Decision rationale (Why Cobra? Why SQLite?)

**Package CLAUDE.md** (package-level architectural guidance):
- What belongs in this package vs elsewhere
- Layer-specific patterns and rules
- Testing strategies
- Examples of proper layer usage

### Documentation Checklist

- [ ] Code implemented and tested
- [ ] README.md updated (if user-facing changes)
- [ ] CLAUDE.md updated (if workflow/architecture changes)
- [ ] Package CLAUDE.md updated (if package responsibilities changed)
- [ ] All tests pass
- [ ] Linter passes (zero violations)

---

## Key References

- **Main Architecture**: `internal/task_manager/CLAUDE.md` - Task Manager architecture overview
- **Package Documentation**: `<package>/CLAUDE.md` - Layer-specific architectural guidance
- **Linter**: `go-arch-lint .` - Validate architecture compliance
- **Clean Architecture**: [The Clean Architecture (Uncle Bob)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- **DDD**: [Domain-Driven Design Reference](https://www.domainlanguage.com/ddd/reference/)

---

## Common Patterns

### Adding a New Command

1. **Presentation**: Create command file in `presentation/cli/`
   - Define Cobra command
   - Parse flags and validate input
   - Call application service

2. **Application**: Create/update service in `application/`
   - Implement use case
   - Orchestrate domain logic
   - Handle transactions

3. **Domain**: Update entities/repositories if needed
   - Add domain logic to entities
   - Define new repository methods in interfaces

4. **Infrastructure**: Implement repository methods
   - Add SQL queries
   - Implement data mapping

5. **Wire**: Update `cmd/tm/main.go`
   - Register new command
   - Inject dependencies

### Adding a New Entity

1. **Domain**: Define entity in `domain/`
   - Create struct with business logic
   - Add validation methods
   - Define repository interface

2. **Infrastructure**: Implement repository
   - Create SQL migrations
   - Implement CRUD operations
   - Handle transactions

3. **Application**: Create service if needed
   - Implement use cases for entity
   - Coordinate with other services

4. **Presentation**: Add commands
   - Create/update/delete commands
   - List/show commands
   - Output formatting

### Repository Pattern Example

```go
// 1. Domain defines interface (domain/repositories.go)
type TaskRepository interface {
    Save(task *Task) error
    FindByID(id string) (*Task, error)
    FindAll() ([]*Task, error)
}

// 2. Infrastructure implements (infrastructure/sqlite_task_repository.go)
type sqliteTaskRepository struct {
    db *sql.DB
}

func NewSQLiteTaskRepository(db *sql.DB) domain.TaskRepository {
    return &sqliteTaskRepository{db: db}
}

func (r *sqliteTaskRepository) Save(task *domain.Task) error {
    // SQL implementation
}

// 3. Application receives injection (application/task_service.go)
type TaskService struct {
    taskRepo domain.TaskRepository
}

func NewTaskService(taskRepo domain.TaskRepository) *TaskService {
    return &TaskService{taskRepo: taskRepo}
}

// 4. Main wires dependencies (cmd/tm/main.go)
db := setupDatabase()
taskRepo := infrastructure.NewSQLiteTaskRepository(db)
taskService := application.NewTaskService(taskRepo)
```

---

**Remember**:
- Domain is the heart - protect it
- Dependencies flow inward
- Use dependency injection
- Test each layer independently
- Zero linter violations before commit
