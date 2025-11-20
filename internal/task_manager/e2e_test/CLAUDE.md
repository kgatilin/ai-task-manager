# E2E Test Suite - Task Manager Standalone Binary

## Overview

End-to-end tests validate the complete CLI command workflow by building the standalone `tm` binary from source and executing real commands. Tests ensure CLI commands work correctly with database persistence, command parsing, and cross-entity operations.

**Binary**: Standalone `tm` CLI
**Commands**: All tests use `tm <command>` directly (e.g., `tm project create`, `tm track create`, `tm iteration start`)
**Working Directory**: Tests use `TM_WORKING_DIR` environment variable pointing to `.tm/`

---

## Architecture

### Test Main Function

**CRITICAL**: Binary built once in `TestMain()` before all tests:

```go
func TestMain(m *testing.M) {
    // Build binary: ../../../../tm-e2e-test
    projectRoot := "../../../../"
    tmBinaryPath = filepath.Join(projectRoot, tmBinaryName)
    buildCmd := exec.Command("go", "build", "-o", tmBinaryPath, "./cmd/tm")
    // Run tests, then cleanup
}
```

**Why**: Tests the actual codebase, not system-installed binary. Build failures caught before tests run.

### Suite Types

**Base Suite (`E2ETestSuite`):**
- Creates unique project per suite (`e2e-test-{timestamp}`)
- Isolated working directory in `/tmp`
- Initializes roadmap in `SetupSuite()`
- All tests within suite share same project/database

**Specialized Suites:**

1. **Regular Suites** (TrackTestSuite, TaskTestSuite, etc.):
   - Tests create their own data
   - Work regardless of DB state
   - Run in shared project

2. **Workflow Suites** (WorkflowTestSuite, IterationWorkflowTestSuite):
   - Test complete workflows
   - Require clean DB state
   - Separate suite instance = separate project = isolated DB

---

## Test Types

### Individual Command Tests
Test single commands in isolation:
- Create/Update/Delete/Show/List operations
- Error cases (missing flags, invalid IDs)
- Validation (duplicate checks, constraints)

**Example**: `TestTrackCreate`, `TestTaskUpdate`, `TestIterationDelete`

### Workflow Tests
Test complete multi-command scenarios:
- Cross-entity workflows (Track → ADR → Task → AC → Iteration)
- Lifecycle tests (create → update → verify → complete → delete)
- State-dependent operations (numbering, empty state checks)

**Example**: `TestCompleteTaskLifecycle`, `TestIterationCompleteWorkflow`

### State-Dependent Tests
Tests that require clean database:
- Empty list verification
- Sequential numbering (iterations 1, 2, 3)
- MAX+1 numbering after deletion (TM-ac-522)

**Solution**: Separate suite instance (gets own project/database)

---

## Mandatory Requirements

### 1. Temporary Directories

**ALWAYS use temp directories for test isolation:**

```go
func (s *E2ETestSuite) SetupSuite() {
    tmpDir := os.TempDir()
    s.testWorkingDir = filepath.Join(tmpDir, "tm-e2e-wd-"+s.projectName)
    os.MkdirAll(s.testWorkingDir, 0755)
}
```

**Why**:
- Tests run in isolation
- No conflicts with real user data
- CI/CD compatibility
- Automatic cleanup

### 2. Environment Variable

**Set `TM_WORKING_DIR` for all commands:**

```go
func (s *E2ETestSuite) run(args ...string) (string, error) {
    cmd := exec.Command(tmBinaryPath, args...)
    cmd.Env = append(os.Environ(), "TM_WORKING_DIR="+s.testWorkingDir)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

**Why**: Ensures consistent working directory for `.tm/active-project.txt` and project databases.

### 3. Unique Project Names

**Generate unique project per suite:**

```go
s.projectName = "e2e-test-" + strconv.FormatInt(time.Now().UnixNano(), 36)
```

**Why**: Parallel test execution without conflicts.

---

## Helper Methods

### parseID(output, entityPattern)

**Extracts entity IDs from command output:**

```go
trackID := s.parseID(output, "-track-")  // Finds "EET-track-1"
taskID := s.parseID(output, "-task-")    // Finds "EET-task-5"
adrID := s.parseID(output, "-adr-")      // Finds "EET-adr-2"
acID := s.parseID(output, "-ac-")        // Finds "EET-ac-10"
```

**IMPORTANT**: Use entity pattern (e.g., `"-track-"`), NOT output text (e.g., `"Created track:"`).

**Why**: Output format changes between commands, but ID format is consistent.

### parseIterationNumber(output)

**Extracts iteration number from create output:**

```go
iterNumber := s.parseIterationNumber(output)  // Finds "3" from "Number: 3"
```

**Why**: Iteration numbers auto-increment; cannot hardcode.

---

## Best Practices

### Do:
- ✅ Extract IDs dynamically using `parseID()`
- ✅ Use `t.TempDir()` or suite-level temp directories
- ✅ Create separate suites for state-dependent tests
- ✅ Test error cases (missing flags, invalid values)
- ✅ Verify command output contains expected data
- ✅ Use `requireSuccess()` / `requireError()` helpers

### Don't:
- ❌ Hardcode entity IDs or iteration numbers
- ❌ Write to real user directories (`~/.tm`)
- ❌ Use system-installed `tm` binary from PATH
- ❌ Assume test execution order within suite
- ❌ Share mutable state between tests
- ❌ Use `parseID()` with output text patterns

---

## Known Limitations

### Parallel Operations
**Parallel entity creation tests are SKIPPED** (`TestParallelSafety`).

**Issue**: Race conditions in database operations when creating entities concurrently.

**Status**: Not currently supported. Tests marked with `s.T().Skip()`.

### Suite Isolation
**Tests within a suite run sequentially**, not in parallel.

**Why**: testify/suite shares single instance across tests. For parallel execution, use regular Go tests with `t.Parallel()`.

---

## Test Execution

**Run all E2E tests:**
```bash
go test ./internal/task_manager/e2e_test/... -v
```

**Run specific suite:**
```bash
go test ./internal/task_manager/e2e_test/... -v -run TestIterationSuite
```

**Run single test:**
```bash
go test ./internal/task_manager/e2e_test/... -v -run TestIterationSuite/TestIterationCreate
```

**Expected behavior:**
- `tm` binary builds once before all tests
- Each suite gets unique project + working directory
- Tests pass regardless of order
- Skipped tests show as `SKIP` (not failures)

---

## Adding New Tests

### For Individual Commands:
1. Add test to appropriate existing suite (TrackTestSuite, TaskTestSuite, etc.)
2. Create entity, extract ID with `parseID()`
3. Test command, verify output
4. Tests are automatically isolated (suite handles setup)

### For Workflows:
1. Add to `WorkflowTestSuite` if needs existing data
2. Create NEW suite if needs clean database
3. Document state requirements in test comments
4. Use descriptive step comments (Step 1, Step 2, etc.)

### For State-Dependent Features:
1. Create separate suite instance (e.g., `MyFeatureWorkflowTestSuite`)
2. Override `SetupSuite()` if needed
3. Use `parseID()` / `parseIterationNumber()` for dynamic extraction
4. Verify both success and empty states

---

## Command Examples

**All tests use `tm` directly**:

```bash
# Project commands
tm project create "My Project"
tm project list
tm project switch "My Project"

# Roadmap commands
tm roadmap init --vision "..." --success-criteria "..."
tm roadmap show

# Track commands
tm track create --title "..." --description "..."
tm track list
tm track show TM-track-1

# Task commands
tm task create --track TM-track-1 --title "..."
tm task list
tm task update TM-task-1 --status in-progress

# Iteration commands
tm iteration create --name "..." --goal "..."
tm iteration start 1
tm iteration add-task 1 TM-task-1

# ADR commands
tm adr create --track TM-track-1 --title "..." --context "..." --decision "..."
tm adr list

# Acceptance Criteria commands
tm ac add TM-task-1 --description "..." --testing-instructions "..."
tm ac list TM-task-1
tm ac verify TM-ac-1
```

---

**Last Updated**: 2025-11-20 (Iteration #37 - Standalone `tm` binary)
