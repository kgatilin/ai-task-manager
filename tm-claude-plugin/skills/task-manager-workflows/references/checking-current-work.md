# Checking Current Work

Detailed guide for understanding what to work on.

## Priority Order

Always follow this sequence:

### 1. Check Current Iteration

```bash
tm iteration current
```

**What it shows:**
- Current active iteration (status: `in-progress`)
- OR next planned iteration if none is active (status: `planning`)
- Iteration number, name, goal, and deliverables

**Why start here:** The iteration is your primary focus. Everything else is context.

### 2. View Iteration Details

```bash
tm iteration show <num>
```

**What it shows:**
- Goal and deliverables
- All tasks in the iteration
- Task statuses (todo/in-progress/done)
- Task priorities

**Use this to:** Understand the full scope of the iteration and what remains.

### 3. View Specific Task

```bash
tm task show TM-task-X
```

**What it shows:**
- Task title and description
- Current status
- Priority
- Associated track
- All acceptance criteria with testing instructions
- Git branch (if set)

**Use this to:** Understand exactly what needs to be done and how to verify completion.

### 4. View Track Context (Optional)

```bash
tm track show TM-track-X
```

**What it shows:**
- Track title and description
- All tasks in the track
- Track status and priority
- Associated documents (ADRs, plans)
- Dependencies

**Use this to:** Understand the broader context and architectural decisions.

## Common Checks

### What Am I Working On?

```bash
tm task list --status in-progress
```

Shows all tasks currently marked as in-progress.

### What's in My Backlog?

```bash
tm task list --status todo
```

Shows all tasks waiting to be started.

### What Have I Completed?

```bash
tm task list --status done
```

Shows all completed tasks.

### What's Blocked?

```bash
tm ac failed --iteration <num>
```

Shows acceptance criteria that have failed verification, indicating blockers.

### Track-Specific Work

```bash
tm task list --track TM-track-X
```

Shows all tasks in a specific track (work stream).

## Daily Standup Pattern

**Morning routine:**
```bash
# 1. What iteration am I on?
tm iteration current

# 2. What tasks are in it?
tm iteration show <num>

# 3. What am I currently working on?
tm task list --status in-progress

# 4. Any blockers?
tm ac failed --iteration <num>
```

## Project Overview Pattern

**Quick status check:**
```bash
# Current focus
tm iteration current

# Active work
tm task list --status in-progress

# Blockers
tm ac failed --iteration <num>
```

This gives you a complete picture in three commands.

## Tips

- **Always start with `tm iteration current`** - it's your north star
- **Check AC before starting work** - use `tm task show TM-task-X` to see what needs verification
- **Update status as you work** - don't batch status updates
- **Check failed AC regularly** - run `tm ac failed --iteration <num>` to catch issues early
