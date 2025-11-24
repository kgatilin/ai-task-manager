# Adding New Tasks

Detailed guide for creating tasks with good acceptance criteria.

## Creating Tasks

### Basic Task Creation

```bash
tm task create \
  --track TM-track-X \
  --title "User-facing feature name" \
  --description "Detailed description of what needs to be done" \
  --priority high|medium|low \
  --rank 100
```

**Required fields:**
- `--track`: Which track (work stream) this task belongs to
- `--title`: Short, descriptive name

**Optional fields:**
- `--description`: Detailed explanation
- `--priority`: high, medium (default), low
- `--rank`: Ordering within track (default: 0)

### Task Granularity

**Good task (user-facing feature):**
- Title: "Add task comment commands"
- Has 3-5 acceptance criteria
- Delivers observable value
- Can be demonstrated independently

**Too granular (implementation detail):**
- Title: "Create TaskComment entity" ❌
- Title: "Add database migration" ❌
- Title: "Define CommentRepository interface" ❌

**Fix:** Merge these into one task: "Add task comment commands"

## Adding Acceptance Criteria

### The Two-Field Structure

**CRITICAL:** Always use separate fields:

```bash
tm ac add TM-task-X \
  --description "WHAT must be verified" \
  --testing-instructions "HOW to verify it"
```

### Description Field (WHAT)

**Good descriptions:**
- "Users can list tasks filtered by status"
- "CLI shows task acceptance criteria in detail view"
- "Task comments persist across sessions"
- "Domain layer has 90%+ test coverage"

**Bad descriptions:**
- "TaskRepository.FindByStatus implements filtering" ❌
- "Add WHERE clause to SQL query" ❌
- "Implement GetComments method" ❌

**Rule:** Describe **what users can verify**, not how code implements it.

### Testing Instructions Field (HOW)

**Format:** Numbered steps with exact commands and expected outcomes.

**Good testing instructions:**
```
1. Run: tm task list --status todo
2. Verify: Output shows only tasks with status 'todo'
3. Run: tm task list --status done
4. Verify: Output shows only tasks with status 'done'
5. Run: tm task list
6. Verify: Output shows all tasks regardless of status
```

**Bad testing instructions:**
```
Test the filtering functionality
```

**Rule:** Be specific. Include exact commands and expected results.

## Complete Example

### User-Facing Feature Task

```bash
# Create the task
tm task create \
  --track TM-track-3 \
  --title "Add task comment commands" \
  --priority high \
  --rank 100

# Add acceptance criteria
tm ac add TM-task-15 \
  --description "Users can add comments to tasks" \
  --testing-instructions "1. Run: tm task comment add TM-task-1 --text 'Test comment'
2. Verify: Command succeeds with confirmation message
3. Run: tm task show TM-task-1
4. Verify: Comment appears in task details"

tm ac add TM-task-15 \
  --description "Users can list all comments on a task" \
  --testing-instructions "1. Add two comments: tm task comment add TM-task-1 --text 'First'
2. Add another: tm task comment add TM-task-1 --text 'Second'
3. Run: tm task comment list TM-task-1
4. Verify: Both comments appear with timestamps"

tm ac add TM-task-15 \
  --description "Comments persist across CLI sessions" \
  --testing-instructions "1. Add comment: tm task comment add TM-task-1 --text 'Persisted'
2. Exit the terminal
3. Open new terminal
4. Run: tm task comment list TM-task-1
5. Verify: Comment still exists"

tm ac add TM-task-15 \
  --description "CLI validates comment text is not empty" \
  --testing-instructions "1. Run: tm task comment add TM-task-1 --text ''
2. Verify: Command fails with error: 'Comment text cannot be empty'
3. Run: tm task comment add TM-task-1
4. Verify: Command fails with error: 'Comment text is required'"
```

## Good vs Bad Examples

### Example 1: Test Coverage

**❌ Bad:**
```bash
tm ac add TM-task-X \
  --description "Domain layer tests are implemented with good coverage"
```

**✅ Good:**
```bash
tm ac add TM-task-X \
  --description "Domain layer has 90%+ test coverage with all tests passing" \
  --testing-instructions "1. Run: cd internal/task_manager/domain
2. Run: go test ./... -coverprofile=coverage.out
3. Run: go tool cover -func=coverage.out | grep total
4. Verify: Total coverage >= 90%
5. Run: go test ./... -v
6. Verify: All tests pass with zero failures"
```

### Example 2: CLI Command

**❌ Bad:**
```bash
tm ac add TM-task-X \
  --description "TaskListCommand implements status filtering in Execute method"
```

**✅ Good:**
```bash
tm ac add TM-task-X \
  --description "Users can filter task list by status" \
  --testing-instructions "1. Run: tm task list --status todo
2. Verify: Shows only todo tasks
3. Run: tm task list --status in-progress
4. Verify: Shows only in-progress tasks
5. Run: tm task list --status done
6. Verify: Shows only done tasks"
```

### Example 3: Data Persistence

**❌ Bad:**
```bash
tm ac add TM-task-X \
  --description "TaskRepository.Save correctly inserts into SQLite database"
```

**✅ Good:**
```bash
tm ac add TM-task-X \
  --description "Task data persists across CLI restarts" \
  --testing-instructions "1. Create task: tm task create --track TM-track-1 --title 'Test'
2. Note the task ID (e.g., TM-task-5)
3. Exit terminal completely
4. Open new terminal
5. Run: tm task show TM-task-5
6. Verify: Task data is still present"
```

## Tips for Writing Good AC

1. **Think like a user** - What can they observe and verify?
2. **Be specific** - Use exact commands in testing instructions
3. **Test boundaries** - Include validation and error cases
4. **One AC per behavior** - Don't combine multiple things
5. **Make it runnable** - Anyone should be able to copy-paste and verify
6. **Avoid implementation** - Don't mention classes, methods, repositories
7. **Focus on outcomes** - What happens, not how it happens

## Common Mistakes

**Mixing implementation and behavior:**
```bash
# ❌ Bad - mentions implementation
--description "Repository implements caching with TTL"

# ✅ Good - focuses on behavior
--description "CLI responds to queries under 100ms after initial load"
```

**Vague testing instructions:**
```bash
# ❌ Bad - too vague
--testing-instructions "Test that it works"

# ✅ Good - specific steps
--testing-instructions "1. Run: tm task list
2. Verify: Response time < 100ms
3. Run: tm task list again
4. Verify: Response time < 10ms (cached)"
```

**Too granular:**
```bash
# ❌ Bad - implementation detail
tm task create --title "Add TaskComment entity"
tm task create --title "Implement CommentRepository"
tm task create --title "Add comment CLI commands"

# ✅ Good - complete feature
tm task create --title "Add task comment commands"
# Then add AC for: add comment, list comments, persistence, validation
```
