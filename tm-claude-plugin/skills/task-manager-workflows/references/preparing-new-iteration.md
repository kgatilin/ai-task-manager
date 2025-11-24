# Preparing New Iteration

Detailed guide for planning and starting iterations.

## Iteration Planning Process

### 1. Create the Iteration

```bash
tm iteration create \
  --name "Iteration 5: Core Features" \
  --goal "Implement core task management features" \
  --deliverable "Users can create, list, update, and delete tasks"
```

**Required fields:**
- `--name`: Clear, descriptive name (include iteration number)
- `--goal`: What you're trying to achieve
- `--deliverable`: What users get at the end

**Tips:**
- **Name**: Include iteration number for easy reference
- **Goal**: Keep it concise (1 sentence)
- **Deliverable**: Describe from user perspective

### 2. Select Tasks

**Review available tasks:**
```bash
tm task list --status todo
```

**Check task details:**
```bash
tm task show TM-task-X
```

Read the description and acceptance criteria to understand the work involved.

**Consider:**
- Task priority (high priority tasks first)
- Task dependencies (complete blocking tasks first)
- Iteration capacity (don't overload)
- Goal alignment (tasks should support iteration goal)

### 3. Add Tasks to Iteration

```bash
tm iteration add-task <num> TM-task-1 TM-task-2 TM-task-3
```

**Add multiple tasks at once:**
You can list multiple task IDs separated by spaces.

**Review iteration scope:**
```bash
tm iteration show <num>
```

Check that the iteration is reasonable in scope (typically 5-10 tasks).

### 4. Create Supporting Documents (Optional)

**Create iteration plan:**
```bash
tm doc create \
  --title "Iteration 5 Plan" \
  --type plan \
  --from-file ./iteration-5-plan.md \
  --iteration 5
```

**Create ADR for architectural decisions:**
```bash
tm doc create \
  --title "ADR: Use SQLite for Persistence" \
  --type adr \
  --from-file ./adr-001-sqlite.md \
  --track TM-track-1
```

### 5. Start the Iteration

```bash
tm iteration start <num>
```

This marks the iteration as `in-progress` and makes it the active iteration.

**Verify it started:**
```bash
tm iteration current
```

Should show your iteration as active.

## Complete Example

### Planning Iteration 5

```bash
# 1. Create iteration
tm iteration create \
  --name "Iteration 5: Task Comments" \
  --goal "Add commenting functionality to tasks" \
  --deliverable "Users can add, view, and manage comments on tasks"

# 2. Review available tasks
tm task list --status todo

# Output shows:
# TM-task-15: Add task comment commands
# TM-task-16: Show comments in task detail view
# TM-task-17: Add comment search functionality
# TM-task-18: Export comments to markdown

# 3. Check task details
tm task show TM-task-15
tm task show TM-task-16

# 4. Add tasks to iteration (select tasks that fit the goal)
tm iteration add-task 5 TM-task-15 TM-task-16

# 5. Review iteration
tm iteration show 5

# 6. Create iteration plan (optional)
tm doc create \
  --title "Iteration 5 Plan" \
  --type plan \
  --content "# Iteration 5: Task Comments

## Goal
Add commenting functionality to tasks

## Deliverable
Users can add, view, and manage comments on tasks

## Tasks
1. TM-task-15: Add task comment commands
2. TM-task-16: Show comments in task detail view

## Success Criteria
- All acceptance criteria verified
- Documentation updated
- Tests passing" \
  --iteration 5

# 7. Start the iteration
tm iteration start 5

# 8. Verify
tm iteration current
```

## Iteration Sizing

### Good Iteration Size

**5-10 tasks** is typically reasonable:
- Small enough to complete in focused work period
- Large enough to deliver meaningful value
- Easy to track and manage

**Each task has 3-5 acceptance criteria:**
- Total: 15-50 AC per iteration
- Provides clear verification points
- Reasonable validation scope

### Too Large

**20+ tasks** is usually too much:
- Harder to complete
- Difficult to track
- Loses focus

**Fix:** Break into multiple iterations

### Too Small

**1-2 tasks** is usually too small:
- Doesn't deliver enough value
- Overhead of iteration management
- Consider task granularity

**Fix:** Add more tasks or merge tasks

## Iteration Lifecycle

```
planning → in-progress → completed
```

### Planning Status

**What it means:**
- Iteration exists but hasn't started
- Can add/remove tasks
- Can modify goal/deliverable

**Commands:**
```bash
tm iteration create ...    # Creates in planning status
tm iteration add-task ...  # Modify task membership
```

### In-Progress Status

**What it means:**
- Iteration is active
- This is your current focus
- Tasks are being worked on

**Commands:**
```bash
tm iteration start <num>   # Move to in-progress
tm task update ... --status in-progress  # Work on tasks
tm iteration validate <num>  # Check progress
```

### Completed Status

**What it means:**
- All work is done
- All AC verified
- Deliverable achieved

**Commands:**
```bash
tm iteration complete <num>  # Mark as completed (USER ONLY)
```

## Planning Patterns

### Bottom-Up Planning

```bash
# 1. Create tasks first
tm task create --track TM-track-1 --title "Feature A"
tm task create --track TM-track-1 --title "Feature B"
tm task create --track TM-track-2 --title "Feature C"

# 2. Add acceptance criteria
tm ac add TM-task-X --description "..." --testing-instructions "..."

# 3. Create iteration and add tasks
tm iteration create --name "..." --goal "..." --deliverable "..."
tm iteration add-task <num> TM-task-X TM-task-Y TM-task-Z
```

**Use when:** You have existing backlog of well-defined tasks.

### Top-Down Planning

```bash
# 1. Create iteration with goal
tm iteration create \
  --name "Iteration 5" \
  --goal "Add user authentication" \
  --deliverable "Users can register, login, and logout"

# 2. Create tasks to achieve goal
tm task create --track TM-track-1 --title "Add registration"
tm task create --track TM-track-1 --title "Add login"
tm task create --track TM-track-1 --title "Add logout"

# 3. Add tasks to iteration
tm iteration add-task 5 TM-task-X TM-task-Y TM-task-Z

# 4. Add acceptance criteria to each task
tm ac add TM-task-X --description "..." --testing-instructions "..."
```

**Use when:** Starting fresh with a clear iteration goal.

### Hybrid Planning

```bash
# 1. Create iteration
tm iteration create --name "..." --goal "..." --deliverable "..."

# 2. Add existing tasks that fit
tm iteration add-task <num> TM-task-1 TM-task-2

# 3. Create new tasks as needed
tm task create --track TM-track-X --title "..."
tm iteration add-task <num> TM-task-X
```

**Use when:** You have some existing tasks and need to create more.

## Tips

- **Clear goal**: Iteration goal should be specific and achievable
- **User-focused deliverable**: Describe what users get, not what you build
- **Right-sized**: 5-10 tasks is usually good
- **Cohesive**: Tasks should relate to the iteration goal
- **Track docs**: Use ADRs for architectural decisions
- **Plan document**: Optional but helpful for complex iterations
- **Start explicitly**: Run `tm iteration start <num>` to make it active
