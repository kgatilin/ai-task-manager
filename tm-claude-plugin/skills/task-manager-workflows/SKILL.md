---
name: task-manager-workflows
description: Expert guidance for Task Manager CLI (tm). Helps navigate iterations, tasks, acceptance criteria, and documents. Knows when to use tm iteration current, tm task list, tm ac add, tm iteration validate, and other tm commands. Covers planning workflows, execution patterns, verification processes, and acceptance criteria best practices for software development workflows.
version: 1.0.0
allowed-tools: Bash(tm:*), Bash(./tm:*)
---

# Task Manager Workflows

Task Manager (tm) is a CLI tool for organizing software development. It structures work into **Roadmap → Tracks → Tasks → Iterations** with acceptance criteria for verification.

## Quick Start: The One Command You Need

**Always start here:**
```bash
tm iteration current        # Shows what to work on right now
tm iteration current --full # Detailed view with tasks and acceptance criteria
```

This shows your active iteration or the next planned one. Everything else flows from this.

## Core Workflows

### 1. Checking Current Work

**See what's happening:**
```bash
tm iteration current              # Current iteration
tm iteration show <num>           # Iteration details
tm task list --status in-progress # Active tasks
```

**See details:** [references/checking-current-work.md](references/checking-current-work.md)

### 2. Working on an Iteration

**Task lifecycle:** `todo` → `in-progress` → `done`

```bash
tm task show TM-task-X                     # View task and acceptance criteria
tm task update TM-task-X --status in-progress  # Start work
# ... do the work ...
tm iteration validate <num>                # Check if work is complete
tm ac failed --iteration <num>             # See any failures
```

**Important:** Update status as you work. User verifies AC and marks tasks done.

**See details:** [references/working-on-iteration.md](references/working-on-iteration.md)

### 3. Adding New Tasks

**Create tasks with acceptance criteria:**
```bash
tm task create --track TM-track-X --title "..." --rank 100
tm ac add TM-task-X \
  --description "End-user observable behavior" \
  --testing-instructions "1. Step one
2. Step two
3. Verify: Expected outcome"
```

**Key principle:** Acceptance criteria describe **what users can verify**, not how it's implemented.

**See details:** [references/adding-new-tasks.md](references/adding-new-tasks.md)

### 4. Preparing New Iteration

**Plan and start iteration:**
```bash
# 1. Create iteration
tm iteration create --name "..." --goal "..." --deliverable "..."

# 2. Add tasks
tm iteration add-task <num> TM-task-1 TM-task-2 TM-task-3

# 3. Start iteration
tm iteration start <num>
```

**Optional:** Create documents (ADRs, plans) for architectural decisions.

**See details:** [references/preparing-new-iteration.md](references/preparing-new-iteration.md)

## Acceptance Criteria Best Practices

**Good AC (user-verifiable):**
- ✅ Describes observable behavior
- ✅ Can be tested by end user
- ✅ Focuses on outcomes, not implementation
- ✅ Separate `--description` (WHAT) from `--testing-instructions` (HOW)

**Bad AC (implementation details):**
- ❌ Mentions repositories, services, internal methods
- ❌ Describes code structure
- ❌ Things only developers care about

**Example - Good:**
```bash
tm ac add TM-task-X \
  --description "Users can list tasks filtered by status" \
  --testing-instructions "1. Run: tm task list --status todo
2. Verify: Shows only todo tasks
3. Run: tm task list --status done
4. Verify: Shows only done tasks"
```

**Example - Bad:**
```bash
tm ac add TM-task-X \
  --description "TaskRepository.FindByStatus implements proper filtering with SQL WHERE clause"
```

**See details:** [references/acceptance-criteria-guide.md](references/acceptance-criteria-guide.md)

## Task Granularity

**Core principle:** If you can't write user-verifiable AC, the task is too granular.

**Good task:**
- ✅ Complete user-facing feature
- ✅ Has 3-5 end-user verifiable AC
- ✅ Delivers observable value

**Too granular:**
- ❌ "Create X entity" → Merge into feature that uses it
- ❌ "Add database migration" → Part of feature implementation
- ❌ "Define X interface" → Merge into service implementation

**Example:** Instead of separate tasks for "Create TaskComment entity", "Add CommentRepository", "Add comment CLI", combine into one: "Add task comment commands" with AC covering the full user-facing feature.

## Common Patterns

### Daily Workflow
```bash
tm iteration current              # What iteration?
tm iteration show <num>           # What's in it?
tm task list --status in-progress # What am I working on?
```

### Validation
```bash
tm iteration validate <num>       # Check all AC
tm ac failed --iteration <num>    # See failures
# Fix issues
tm iteration validate <num>       # Re-validate
```

### Quick Status
```bash
tm iteration current
tm task list --status in-progress
tm ac failed --iteration <num>
```

## Agent vs User Responsibilities

**Agent (Claude):**
- Implement work defined in tasks
- Add acceptance criteria during planning
- Update task status to `in-progress` when starting
- Create ADRs and iteration plans

**User:**
- Verify acceptance criteria (run testing instructions)
- Mark AC as verified or failed
- Mark tasks as `done` after verification
- Complete iterations
- Make strategic decisions

**Important:** Agent does NOT verify AC or mark tasks done unless explicitly asked.

## Troubleshooting

**"My AC doesn't make sense"**
→ Rewrite to focus on user-observable behavior. Ask: "What can a user verify without looking at code?"

**"Task is too big"**
→ Break into multiple user-facing tasks (each with 3-5 AC)

**"Don't know what to work on"**
→ Run `tm iteration current` then `tm iteration show <num>`

**"AC keeps failing"**
→ This is normal! Fix issues and re-run `tm iteration validate <num>`

## Key Commands Summary

```bash
# Always start here
tm iteration current

# Task management
tm task list [--status todo|in-progress|done]
tm task show TM-task-X
tm task update TM-task-X --status in-progress

# Acceptance criteria
tm ac add TM-task-X --description "..." --testing-instructions "..."
tm ac list TM-task-X
tm ac list-iteration <num> --with-testing  # See all AC in iteration
tm ac failed --iteration <num>

# Iteration workflow
tm iteration show <num>
tm iteration validate <num>

# Documents
tm doc create --title "..." --type adr --from-file ./file.md --track TM-track-X
```

## Related Resources

- **Command Reference**: See [references/command-reference.md](references/command-reference.md)
- **Detailed Workflows**: See [references/](references/) for in-depth guides
- **Command Help**: `tm --help` and `tm <command> --help`
- **GitHub**: https://github.com/kgatilin/ai-task-manager

**Remember:** Focus on end-user verifiable behavior, not implementation details!
