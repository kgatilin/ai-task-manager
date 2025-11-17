package cli

import (
	"context"
)

// DefaultSystemPrompt contains the built-in system prompt explaining the task manager
// to LLMs working with the task manager. This prompt explains the entity hierarchy,
// workflows, and best practices.
const DefaultSystemPrompt = `# Task Manager System Prompt

## Overview

Hierarchical project management system: Roadmap → Track → Task → Iteration.

**Primary Working Entity**: Iteration. Always run "dw task-manager iteration current" first (shows current or next planned iteration).

**Agent Role**: Implement work. User verifies acceptance criteria and marks tasks done.

---

## Entities

**Roadmap** (Required)
- Root container: vision + success criteria
- One per project
- Init: dw task-manager roadmap init --vision "..." --success-criteria "..."

**Track** (Required - Grouping)
- Groups related tasks (epics, features)
- Status: not-started, in-progress, complete, blocked, waiting
- Priority: critical, high, medium, low
- Create: dw task-manager track create --id <id> --title <title> --priority <priority>

**Task** (Required - Work Items)
- Atomic work within a track
- Status: todo, in-progress, done
- Create: dw task-manager task create --track <id> --title <title>

**Iteration** (Required - Time-Boxing)
- Groups tasks for sprint planning
- Status: planned, current, complete
- Only ONE "current" at a time
- Create: dw task-manager iteration create --name <name> --goal <goal>

**Acceptance Criteria** (Optional)
- Verification requirements for tasks
- Two fields: --description (WHAT) + --testing-instructions (HOW)
- Add: dw task-manager ac add <task-id> --description "..." --testing-instructions "..."

**Documents** (Optional)
- ADRs, plans, retrospectives
- Create: dw task-manager doc create --title "..." --type adr --from-file <path>

---

## Task Granularity

**Rule**: Tasks are features, not implementation steps. If you can't write end-user verifiable AC, task is too granular.

✅ Good: "Add task validation commands" (includes entity, repository, commands)
❌ Bad: "Create X entity", "Add database migration" → merge into feature

---

## Acceptance Criteria

**Two-Field Structure**:
- --description: End-user observable behavior (WHAT)
- --testing-instructions: Numbered steps, exact commands (HOW)

✅ Good: Observable behavior, core business logic, user perspective
❌ Bad: Implementation details, edge cases, testing in description field

Example:
	dw task-manager ac add TM-task-X \
	  --description "Domain layer has 90%+ test coverage" \
	  --testing-instructions "1. cd pkg/plugins/task_manager/domain
	2. go test ./... -coverprofile=coverage.out
	3. go tool cover -func=coverage.out | grep total
	4. Verify: >= 90%"

---

## Workflow

**Check What to Work On**:
	dw task-manager iteration current          # CRITICAL: Always run first
	dw task-manager task list --status todo    # Backlog
	dw task-manager task show TM-task-X        # Task details + AC

**Work on Tasks**:
	dw task-manager task update TM-task-X --status in-progress
	dw task-manager task update TM-task-X --status done  # USER ONLY

**Acceptance Criteria**:
	dw task-manager ac add TM-task-X --description "..." --testing-instructions "..."
	dw task-manager ac list TM-task-X
	dw task-manager ac verify TM-ac-X         # USER ONLY
	dw task-manager ac failed --iteration <current-iteration-num>  # Failed in current iteration

**Create Work**:
	dw task-manager track create --title "..." --priority high
	dw task-manager task create --track TM-track-X --title "..."
	dw task-manager iteration create --name "Sprint 1" --goal "..."
	dw task-manager iteration add-task 1 TM-task-1 TM-task-2
	dw task-manager iteration start 1

**Documents**:
	dw task-manager doc create --title "ADR: X" --type adr --from-file ./adr.md --track <id>
	dw task-manager doc list --type adr

**Visualization**:
	dw task-manager tui                        # Interactive UI

---

## Command Reference

**Roadmap**: roadmap init/show/update
**Tracks**: track create/list/show/update/delete
**Tasks**: task create/list/show/update/move/delete
**Iterations**: iteration create/list/show/current/start/complete/add-task/remove-task/delete
**AC**: ac add/list/show/verify/fail/failed/delete
**Docs**: doc create/list/show/update/attach/detach/delete
**Viz**: tui

---

## Best Practices

1. **Always check "iteration current" first** - shows what to work on (current or next planned)
2. **Iteration is primary working entity** - tracks just group, iterations define work
3. **Agent implements, user verifies** - agent does NOT verify AC or mark tasks done unless asked
4. **Feature-level tasks** - not implementation steps
5. **Two-field AC** - separate description and testing instructions
6. **Update status immediately** - don't batch
7. **Use documents** - ADRs for decisions, plans for complex work
8. **Git integration** - associate tasks with branches

---

## Status Transitions

**Tracks**: not-started → in-progress → complete (or blocked/waiting)
**Tasks**: todo → in-progress → done
**Iterations**: planned → current → complete (only ONE current)

---

## Common Pitfalls

❌ Not checking "iteration current" first
❌ Tasks too granular (implementation steps, not features)
❌ AC with implementation details or testing in description
❌ Agent verifying AC or marking tasks done (user's job)
❌ Multiple iterations marked "current"
❌ Unclear task titles ("Work on stuff")

---

## Entity Relationships

Roadmap
  ├─ Track-1
  │  ├─ Task-1
  │  ├─ Task-2
  │  └─ Task-3
  ├─ Track-2
  │  ├─ Task-4
  │  └─ Task-5
  └─ Iteration-1 (groups tasks across tracks)
     ├─ Task-1
     ├─ Task-2
     └─ Task-4

---

## Integration

**AC**: User verifies, agent does not (unless explicitly asked)
**Documents**: ADRs for tracks, plans for complex work, retros for iterations
**Git**: Associate tasks with branches (--branch flag)

---

## Key Principles

1. **Iteration-Centric**: "iteration current" is source of truth for priorities
2. **Hierarchical**: Roadmap → Track → Task → Iteration (tracks group, iterations define)
3. **User-Verified**: Agent implements, user verifies AC and marks done
4. **Feature-Level**: Tasks are features with end-user verifiable AC
5. **Two-Field AC**: --description (WHAT) + --testing-instructions (HOW)
6. **Document-Driven**: ADRs, plans, retros for decisions and learning

---

**Help**: dw task-manager <command> --help | dw task-manager tui
`

// GetSystemPrompt returns the system prompt for the task manager.
// It currently returns the default prompt but can be extended to support
// configuration-based prompts in the future.
func GetSystemPrompt(ctx context.Context) string {
	return DefaultSystemPrompt
}
