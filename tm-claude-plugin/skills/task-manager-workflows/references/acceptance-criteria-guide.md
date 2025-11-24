# Acceptance Criteria Guide

Comprehensive guide for writing excellent acceptance criteria.

## The Golden Rule

**Acceptance criteria describe WHAT users can verify, not HOW code implements it.**

## Two-Field Structure

Always use both fields correctly:

```bash
tm ac add TM-task-X \
  --description "WHAT must be verified" \
  --testing-instructions "HOW to verify it"
```

### Description Field (WHAT)

**Purpose:** Describe the observable behavior or outcome.

**Good descriptions:**
- Focus on user-facing functionality
- Describe observable behavior
- Avoid implementation details
- Written from user perspective

**Examples:**
- ✅ "Users can list tasks filtered by status"
- ✅ "CLI displays help text for all commands"
- ✅ "Task data persists across CLI restarts"
- ✅ "Domain layer has 90%+ test coverage"

**Bad descriptions:**
- ❌ "TaskRepository.FindByStatus implements filtering"
- ❌ "Add WHERE clause to SQL query"
- ❌ "Implement help command handler"
- ❌ "Mock repository in tests"

### Testing Instructions Field (HOW)

**Purpose:** Provide exact steps to verify the acceptance criterion.

**Format:**
```
1. Run: <exact command>
2. Verify: <expected outcome>
3. Run: <another command>
4. Verify: <expected outcome>
```

**Rules:**
- Numbered steps
- Exact commands (copy-paste ready)
- Clear expected outcomes
- Anyone can execute without extra knowledge

## Complete Examples

### Example 1: CLI Feature

```bash
tm ac add TM-task-X \
  --description "Users can filter task list by status" \
  --testing-instructions "1. Create test tasks: tm task create --track TM-track-1 --title 'Todo task'
2. Update one: tm task update TM-task-1 --status in-progress
3. Update another: tm task update TM-task-2 --status done
4. Run: tm task list --status todo
5. Verify: Shows only todo tasks
6. Run: tm task list --status in-progress
7. Verify: Shows only in-progress tasks
8. Run: tm task list --status done
9. Verify: Shows only done tasks"
```

### Example 2: Data Persistence

```bash
tm ac add TM-task-X \
  --description "Task comments persist across CLI sessions" \
  --testing-instructions "1. Add comment: tm task comment add TM-task-1 --text 'Test comment'
2. Note the comment ID
3. Exit terminal completely (close window)
4. Open new terminal window
5. Run: tm task comment list TM-task-1
6. Verify: Comment with same ID and text is present"
```

### Example 3: Validation

```bash
tm ac add TM-task-X \
  --description "CLI validates task title is not empty" \
  --testing-instructions "1. Run: tm task create --track TM-track-1 --title ''
2. Verify: Command fails with error 'Task title cannot be empty'
3. Run: tm task create --track TM-track-1
4. Verify: Command fails with error 'Task title is required'
5. Run: tm task create --track TM-track-1 --title 'Valid'
6. Verify: Command succeeds and task is created"
```

### Example 4: Test Coverage

```bash
tm ac add TM-task-X \
  --description "Domain layer has 90%+ test coverage with all tests passing" \
  --testing-instructions "1. Navigate: cd internal/task_manager/domain
2. Run: go test ./... -coverprofile=coverage.out
3. Run: go tool cover -func=coverage.out | grep total
4. Verify: Total coverage shows >= 90%
5. Run: go test ./... -v
6. Verify: All tests pass, zero failures
7. Verify: No skipped tests"
```

### Example 5: Help Documentation

```bash
tm ac add TM-task-X \
  --description "All CLI commands display helpful usage information" \
  --testing-instructions "1. Run: tm --help
2. Verify: Shows list of all commands
3. Run: tm task --help
4. Verify: Shows task subcommands
5. Run: tm task create --help
6. Verify: Shows required flags and examples
7. Run: tm iteration --help
8. Verify: Shows iteration subcommands"
```

## Good vs Bad Patterns

### Pattern: CLI Output

**❌ Bad:**
```bash
--description "ListCommand formats output properly"
```

**✅ Good:**
```bash
--description "Task list displays ID, title, and status in table format"
--testing-instructions "1. Run: tm task list
2. Verify: Output shows table headers: ID | Title | Status
3. Verify: Each row shows task ID (e.g., TM-task-1)
4. Verify: Each row shows task title
5. Verify: Each row shows task status (todo/in-progress/done)"
```

### Pattern: Error Handling

**❌ Bad:**
```bash
--description "Repository returns proper error type"
```

**✅ Good:**
```bash
--description "CLI shows clear error when task not found"
--testing-instructions "1. Run: tm task show TM-task-999
2. Verify: Shows error 'Task TM-task-999 not found'
3. Verify: Exit code is non-zero
4. Verify: Error message is red/highlighted"
```

### Pattern: Performance

**❌ Bad:**
```bash
--description "Query is optimized with index"
```

**✅ Good:**
```bash
--description "Task list loads in under 100ms for 1000 tasks"
--testing-instructions "1. Create 1000 tasks using script
2. Clear any caches
3. Run: time tm task list
4. Verify: Execution time < 100ms
5. Run again: time tm task list
6. Verify: Subsequent calls also < 100ms"
```

### Pattern: Integration

**❌ Bad:**
```bash
--description "Service calls repository correctly"
```

**✅ Good:**
```bash
--description "Creating a task automatically adds it to track"
--testing-instructions "1. Run: tm task create --track TM-track-1 --title 'Test'
2. Note the task ID
3. Run: tm track show TM-track-1
4. Verify: New task appears in track's task list
5. Run: tm task show <task-id>
6. Verify: Task shows track: TM-track-1"
```

## Writing Tips

### Start with "Users can..."

Good starters for description:
- "Users can [action]"
- "CLI displays [information]"
- "System persists [data]"
- "Command validates [input]"

**Examples:**
- "Users can create tasks with optional description"
- "CLI displays acceptance criteria in task detail view"
- "System persists task status across restarts"
- "Command validates iteration number is positive"

### Use Exact Commands

Testing instructions must have exact, copy-paste commands:

**❌ Vague:**
```
1. Create a task
2. Check it's there
```

**✅ Exact:**
```
1. Run: tm task create --track TM-track-1 --title 'Test Task'
2. Note the task ID from output (e.g., TM-task-5)
3. Run: tm task show TM-task-5
4. Verify: Title shows 'Test Task'
```

### Include Verification Steps

Every testing instruction needs clear verification:

**❌ Missing verification:**
```
1. Run: tm task list
```

**✅ With verification:**
```
1. Run: tm task list
2. Verify: Output contains table with columns: ID, Title, Status
3. Verify: All tasks are listed
4. Verify: No error messages appear
```

### Test Edge Cases

Good AC includes boundary and error conditions:

```bash
# Normal case
tm ac add TM-task-X \
  --description "Users can create tasks with title" \
  --testing-instructions "..."

# Edge case: empty title
tm ac add TM-task-X \
  --description "CLI rejects empty task titles" \
  --testing-instructions "..."

# Edge case: very long title
tm ac add TM-task-X \
  --description "CLI handles task titles up to 200 characters" \
  --testing-instructions "..."
```

## Common Mistakes

### Mistake 1: Implementation in Description

**❌ Wrong:**
```bash
--description "TaskRepository.Save inserts row into SQLite"
```

**✅ Right:**
```bash
--description "Created tasks persist after CLI restart"
```

### Mistake 2: No Clear Verification

**❌ Wrong:**
```bash
--testing-instructions "1. Run the command
2. Check it works"
```

**✅ Right:**
```bash
--testing-instructions "1. Run: tm task list
2. Verify: Output shows table with columns ID, Title, Status
3. Verify: All existing tasks appear in the list"
```

### Mistake 3: Too Many Things in One AC

**❌ Wrong:**
```bash
--description "Users can create, list, update, and delete tasks"
```

**✅ Right (separate AC):**
```bash
# AC 1
--description "Users can create tasks"

# AC 2
--description "Users can list all tasks"

# AC 3
--description "Users can update task status"

# AC 4
--description "Users can delete tasks"
```

### Mistake 4: Testing Instructions in Description

**❌ Wrong:**
```bash
--description "Users can list tasks. Run tm task list and verify output."
```

**✅ Right:**
```bash
--description "Users can list all tasks"
--testing-instructions "1. Run: tm task list
2. Verify: Shows all tasks in table format"
```

## Checklist for Good AC

Before submitting an AC, verify:

- [ ] Description focuses on user-observable behavior
- [ ] No implementation details in description
- [ ] Testing instructions are numbered steps
- [ ] Every command is exact and copy-paste ready
- [ ] Every step has clear verification
- [ ] Anyone can execute without code knowledge
- [ ] Covers the happy path
- [ ] Includes error cases if relevant
- [ ] One acceptance criterion per behavior
- [ ] Both fields are filled out
