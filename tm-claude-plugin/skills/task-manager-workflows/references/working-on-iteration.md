# Working on an Iteration

Detailed guide for executing work within an iteration.

## Task Lifecycle

Tasks move through three statuses:

```
todo → in-progress → done
```

### Starting Work on a Task

**1. Review the task:**
```bash
tm task show TM-task-X
```

Read the description and acceptance criteria carefully. Understand what needs to be verified.

**2. Mark as in-progress:**
```bash
tm task update TM-task-X --status in-progress
```

Do this **immediately** when you start work. Don't wait.

**3. Do the work:**
Implement the feature, fix the bug, or complete the task according to its description.

**4. Check acceptance criteria:**
```bash
tm ac list TM-task-X
```

Review what needs to be verified. Run the testing instructions yourself to check your work.

### During Work

**Update status immediately:**
Don't batch status updates. When you start a task, mark it in-progress right away.

**Check AC as you work:**
Refer to acceptance criteria frequently to ensure you're building the right thing.

**Ask for clarification:**
If AC is unclear, ask the user to clarify before proceeding.

## Validation

### Validating Your Work

**Run full iteration validation:**
```bash
tm iteration validate <num>
```

This checks ALL acceptance criteria in the iteration and reports pass/fail status.

**See failures:**
```bash
tm ac failed --iteration <num>
```

Shows all failed AC grouped by task, with feedback explaining why they failed.

**Task-specific validation:**
```bash
tm ac list TM-task-X
tm ac failed --task TM-task-X
```

Check just one task's acceptance criteria.

### Handling Failed AC

**1. Review the failure:**
```bash
tm ac failed --task TM-task-X
```

Read the feedback to understand what went wrong.

**2. Fix the issue:**
Implement the necessary changes to satisfy the AC.

**3. Re-validate:**
```bash
tm iteration validate <num>
```

Run validation again to check if the issue is resolved.

**4. Repeat until all pass:**
This is normal! AC verification is an iterative process.

## Completing Tasks

### Agent Responsibility

**What agents do:**
- Implement the work
- Update status to `in-progress` when starting
- Self-check AC by running testing instructions
- Fix any issues found during validation

**What agents do NOT do:**
- Mark AC as verified (that's the user's job)
- Mark tasks as `done` (unless explicitly asked)
- Complete iterations (user decision)

### User Responsibility

**What users do:**
- Run testing instructions from AC
- Mark AC as verified: `tm ac verify TM-ac-X`
- Mark AC as failed: `tm ac fail TM-ac-X --feedback "..."`
- Mark tasks as done: `tm task update TM-task-X --status done`
- Complete iterations: `tm iteration complete <num>`

### Marking Task as Done

**Only the user marks tasks done:**
```bash
tm task update TM-task-X --status done
```

Agents do NOT do this unless explicitly asked.

**Before marking done:**
1. All acceptance criteria should be verified
2. Run `tm ac failed --task TM-task-X` to ensure nothing failed
3. User has tested the work

## Patterns

### Standard Work Pattern

```bash
# Morning: Check what to work on
tm iteration current
tm iteration show <num>

# Start a task
tm task show TM-task-5
tm task update TM-task-5 --status in-progress

# During work: Check AC
tm ac list TM-task-5

# After work: Validate
tm iteration validate <num>
tm ac failed --iteration <num>

# Fix any issues and re-validate
tm iteration validate <num>
```

### Multiple Tasks Pattern

```bash
# See all in-progress work
tm task list --status in-progress

# Switch between tasks
tm task show TM-task-5
# Work on task 5
tm task show TM-task-7
# Work on task 7
```

**Note:** Try to focus on one task at a time, but multiple in-progress tasks is okay.

### Blocked Task Pattern

```bash
# Can't complete a task?
# Move it back to todo
tm task update TM-task-X --status todo

# Document the blocker as AC failure
tm ac fail TM-ac-X --feedback "Blocked because: [reason]"

# Work on a different task
tm task list --status todo
tm task update TM-task-Y --status in-progress
```

## Tips

- **Update status immediately** when you start work
- **Check AC frequently** to ensure you're on track
- **Run validation often** - don't wait until the end
- **Expect failures** - validation is iterative, not one-shot
- **Don't mark tasks done yourself** - that's the user's job
- **Focus on AC** - they define what "done" means
