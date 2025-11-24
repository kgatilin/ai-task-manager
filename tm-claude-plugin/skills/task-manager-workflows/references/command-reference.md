# Complete Task Manager Command Reference

**Version:** 1.0.0

## Global Flags

All commands support:
- `--project <name>` - Override active project for this command
- `--help` - Show command help

## Project Management

Commands for managing multiple isolated project databases.

**Aliases:** `project`, `proj`, `p`

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

## Roadmap Commands

```bash
# Initialize roadmap (one per project)
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

## Track Commands (Work Streams)

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

## Task Commands

Commands for creating, updating, and managing tasks within tracks.

**Aliases:** `task`, `t`

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
tm task list --status in-progress
tm task list --status done

# List backlog (all TODO tasks)
tm task backlog

# Show task details (includes acceptance criteria)
tm task show TM-task-1

# Check if all AC are verified
tm task check-ready TM-task-1

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

## Iteration Commands

Commands for creating, updating, and managing iterations.

**Aliases:** `iteration`, `iter`

```bash
# Create iteration
tm iteration create \
  --name "Sprint 1" \
  --goal "Sprint goal" \
  --deliverable "Expected deliverable"

# List iterations
tm iteration list

# Show iteration details
tm iteration show <num>

# Show current iteration (or next planned if none active)
tm iteration current
tm iteration current --full   # Show full details with timestamps and tasks with AC

# Add/remove tasks
tm iteration add-task <num> TM-task-1 TM-task-2
tm iteration remove-task <num> TM-task-1

# Start iteration (mark as in-progress)
tm iteration start <num>

# Update iteration
tm iteration update <num> \
  --name "New name" \
  --goal "New goal" \
  --deliverable "New deliverable"

# Complete iteration
tm iteration complete <num>

# Delete iteration
tm iteration delete <num> --force
```

## Acceptance Criteria Commands

Commands for creating, verifying, and managing acceptance criteria for tasks.

**Aliases:** `ac`, `criterion`, `criteria`

```bash
# Add acceptance criterion
tm ac add TM-task-1 \
  --description "What must be verified" \
  --testing-instructions "1. Run command X
2. Verify Y
3. Check Z"

# List acceptance criteria
tm ac list TM-task-1                           # For specific task
tm ac list-iteration <num>                     # All ACs in iteration
tm ac list-iteration <num> --with-testing      # Include testing instructions

# Show AC details
tm ac show TM-ac-1

# Update AC
tm ac update TM-ac-1 \
  --description "Updated description" \
  --testing-instructions "Updated instructions"

# Verify AC (mark as passed)
tm ac verify TM-ac-1

# Mark AC as failed
tm ac fail TM-ac-1 --feedback "Error message or reason"

# Mark AC as skipped
tm ac skip TM-ac-1

# List failed ACs
tm ac failed                          # All failed
tm ac failed --iteration <num>        # Failed in iteration
tm ac failed --task TM-task-1         # Failed for specific task
tm ac failed --track TM-track-1       # Failed for specific track

# Delete AC
tm ac delete TM-ac-1 --force
```

## Document Commands (ADRs, Plans, etc.)

```bash
# Create document from file
tm doc create \
  --title "ADR: Use Event Sourcing" \
  --type adr|plan|retrospective \
  --from-file ./docs/adr-001.md \
  --track TM-track-1

# Create document with inline content
tm doc create \
  --title "ADR: Use Event Sourcing" \
  --type adr \
  --content "# Context
...

# Decision
..." \
  --track TM-track-1

# List documents
tm doc list
tm doc list --type adr
tm doc list --type plan

# Show document
tm doc show TM-doc-1

# Update document
tm doc update TM-doc-1 --from-file ./updated.md
tm doc update TM-doc-1 --content "Updated content..."

# Attach to track or iteration
tm doc attach TM-doc-1 --track TM-track-2
tm doc attach TM-doc-1 --iteration 3

# Detach document (remove association)
tm doc detach TM-doc-1

# Delete document
tm doc delete TM-doc-1 --force
```

## ADR Commands (Architecture Decision Records)

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

# Update ADR status
tm adr update TM-adr-1 \
  --status proposed|accepted|rejected|superseded|deprecated

# Supersede ADR (mark as replaced by another ADR)
tm adr supersede TM-adr-1 TM-adr-2

# Deprecate ADR
tm adr deprecate TM-adr-1 --reason "No longer relevant"

# Check for superseded/deprecated ADRs
tm adr check

# Delete ADR
tm adr delete TM-adr-1 --force
```

## Interactive TUI

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

## Common Command Patterns

### Quick Status Check
```bash
tm iteration current                          # Current iteration
tm task list --status in-progress             # Active work
tm ac failed --iteration <num>                # Blockers
```

### Start New Work
```bash
tm task show TM-task-X                        # Review task
tm task update TM-task-X --status in-progress # Start work
```

### Validate Work
```bash
tm iteration validate <num>                   # Check all AC
tm ac failed --iteration <num>                # See failures
```

### Create Complete Feature
```bash
# 1. Create track
tm track create --title "..." --description "..." --rank 100

# 2. Create tasks
tm task create --track TM-track-X --title "..." --rank 100

# 3. Add acceptance criteria
tm ac add TM-task-X --description "..." --testing-instructions "..."

# 4. Create iteration
tm iteration create --name "..." --goal "..." --deliverable "..."

# 5. Add tasks to iteration
tm iteration add-task <num> TM-task-X

# 6. Start iteration
tm iteration start <num>
```

## Flag Reference

### Common Command Aliases

- `project`, `proj`, `p` → Project commands
- `task`, `t` → Task commands
- `iteration`, `iter` → Iteration commands
- `ac`, `criterion`, `criteria` → Acceptance criteria commands

### Status Values

**Track Status:**
- `planning` - Initial planning phase
- `in-progress` - Active work
- `done` - Completed
- `blocked` - Waiting on dependencies

**Task Status:**
- `todo` - Not started (backlog)
- `in-progress` - Currently being worked on
- `done` - Completed

**Iteration Status:**
- `planning` - Being planned
- `in-progress` - Active iteration
- `completed` - Finished

**AC Status:**
- `pending` - Not yet verified
- `verified` - Passed verification
- `failed` - Failed verification (includes feedback)

### Priority Values

- `critical` - Urgent and important (tracks only)
- `high` - Important
- `medium` - Normal priority (default)
- `low` - Nice to have

### Document Types

- `adr` - Architecture Decision Record
- `plan` - Planning document (iteration plans, feature plans)
- `retrospective` - Retrospective document (after iteration completion)

## Database Location

Each project stores its data in:
```
.tm/projects/<project-name>/roadmap.db
```

This allows complete isolation between projects and easy backup/versioning.
