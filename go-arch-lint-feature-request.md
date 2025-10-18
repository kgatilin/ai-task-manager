# Feature Request: Enforce Architectural Rules on Test Files

## Summary

Add support for enforcing architectural rules on test files (`*_test.go`) to prevent architectural violations from being introduced through test code.

## Problem Statement

Currently, `go-arch-lint` only validates architectural rules for production code, ignoring test files. This creates a blind spot where architectural violations can be introduced through tests, leading to:

1. **Inconsistent architecture** - Tests bypass architectural boundaries that production code must respect
2. **Implementation coupling** - Tests become tightly coupled to implementation details (database schemas, SQL syntax, internal packages)
3. **Fragile tests** - Tests break when implementation details change, even when interfaces remain stable
4. **Bad examples** - Developers copy patterns from tests into production code, spreading violations

### Real-World Example

In a DDD project using the `ddd` preset, the architecture requires:
```
cmd → internal/app → internal/infra → internal/domain
```

**Production code** (correctly follows architecture):
```go
// cmd/dw/logs.go
package main

import (
    "github.com/myproject/internal/app"
    "github.com/myproject/internal/infra"
)

func handleLogs(args []string) {
    repo, _ := infra.NewSQLiteEventRepository(dbPath)
    service := app.NewLogsService(repo, repo)
    // Uses app layer service - architecturally correct ✓
}
```

**Test code** (violates architecture, but not caught):
```go
// cmd/dw/logs_test.go
package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"  // Direct dependency on DB driver
)

func TestLogs(t *testing.T) {
    db, _ := sql.Open("sqlite3", dbPath)  // Bypasses repository layer
    db.Exec("INSERT INTO events ...")     // Direct SQL - violates DDD ✗

    // Test should use:
    // repo.Save(ctx, event) instead
}
```

**Current behavior**: Linter passes ✓ (incorrectly - violation not detected)

**Desired behavior**: Linter should fail with:
```
[ERROR] Test Architectural Violation
  File: cmd/dw/logs_test.go
  Issue: Test file imports 'github.com/mattn/go-sqlite3' which violates layer rules
  Rule: cmd layer should not directly import database drivers
  Fix: Use infra layer repositories (internal/infra) for database operations
```

## Current Workarounds

### 1. Manual Code Review
- Error-prone
- Time-consuming
- Inconsistently applied

### 2. Separate Test Package Suffix
```go
package main_integration_test  // Different package = different rules?
```
- Confusing
- Still doesn't enforce architecture

### 3. Custom Linter Integration
- Requires maintaining separate tooling
- Duplicates go-arch-lint logic

## Proposed Solution

### Configuration Option 1: Simple Toggle

```yaml
# .goarchlint
lint_test_files: true  # Default: false for backward compatibility

structure:
  required_directories:
    cmd: "Application entry points"
    internal/app: "Application services"
    internal/infra: "Infrastructure implementations"
    internal/domain: "Core business logic"

rules:
  directories_import:
    cmd: [internal/app, internal/infra]
    internal/app: [internal/domain]
    internal/infra: [internal/domain]
    internal/domain: []

  # When lint_test_files: true, these rules apply to *_test.go files too
```

### Configuration Option 2: Granular Control

```yaml
# .goarchlint
test_rules:
  # Apply same rules as production code
  enforce_architecture: true

  # Optional: Allow specific exemptions for test infrastructure
  exemptions:
    packages:
      - "testing"           # Standard library testing tools
      - "github.com/stretchr/testify"  # Common test frameworks

    # Allow specific violations only in test files
    allow_direct_imports:
      "cmd/*_test.go":
        - "database/sql"    # If you REALLY need it (not recommended)

  # Require tests to use proper layers
  require_test_helpers: true  # Enforce using test builders/factories
```

### Configuration Option 3: Test Type Differentiation

```yaml
# .goarchlint
test_rules:
  default:
    enforce_architecture: true

  # Different rules for different test types
  by_suffix:
    "*_integration_test.go":
      allow_direct_sql: true
      allow_direct_http: true
      rationale: "Integration tests may need direct access"

    "*_test.go":
      enforce_architecture: true
      rationale: "Unit tests must respect architectural boundaries"

    "*_e2e_test.go":
      exempt_from_rules: true
      rationale: "E2E tests can use any dependencies"
```

## Rationale & Benefits

### 1. **Consistency**
Tests are code too. If architecture matters in production, it matters in tests.

### 2. **Better Test Quality**
- Tests using proper interfaces are more maintainable
- Tests validate actual user-facing behavior, not implementation details
- Refactoring becomes safer

### 3. **Educational Value**
- Developers learn correct patterns from test examples
- Prevents "test-only" anti-patterns from spreading

### 4. **Catch Issues Early**
Many projects have clean production code but messy test code. This prevents that divergence.

### 5. **Gradual Adoption**
With `lint_test_files: false` as default, existing projects aren't broken. Teams can opt-in when ready.

## Example Use Cases

### Use Case 1: DDD Architecture
**Before** (violation not caught):
```go
// cmd/service/handler_test.go
func TestHandler(t *testing.T) {
    db, _ := sql.Open("postgres", connStr)  // ✗ Violates DDD
    db.Exec("INSERT INTO users ...")
}
```

**After** (properly caught):
```
[ERROR] Architecture Violation in Test
  File: cmd/service/handler_test.go:5
  Issue: 'cmd' layer directly imports 'database/sql'
  Expected: Use internal/infra repository layer

  Fix:
    repo := infra.NewUserRepository(connStr)
    repo.Save(ctx, user)
```

### Use Case 2: Shared External Imports
**Before** (violation not caught):
```go
// pkg/api/handler_test.go
import "github.com/lib/pq"  // PostgreSQL driver

// internal/infra/postgres.go
import "github.com/lib/pq"  // Same driver

// Linter doesn't complain about pkg importing same driver as infra
```

**After** (properly caught):
```
[ERROR] Shared External Import in Test
  Package: github.com/lib/pq
  Imported by:
    - pkg/api/handler_test.go (layer: pkg)
    - internal/infra/postgres.go (layer: internal/infra)

  Rule: External packages should be owned by a single layer
  Fix: pkg tests should use mocks or infra layer repositories
```

### Use Case 3: Layer Jumping in Tests
**Before**:
```go
// cmd/cli/command_test.go
import "myproject/internal/domain"  // ✗ cmd → domain (skips app layer)

func TestCommand(t *testing.T) {
    user := domain.NewUser(...)  // Direct domain access
}
```

**After**:
```
[ERROR] Layer Jumping in Test
  File: cmd/cli/command_test.go
  Issue: 'cmd' imports 'internal/domain' directly
  Expected: cmd → internal/app → internal/domain

  Fix: Use app layer services:
    service := app.NewUserService(...)
    service.CreateUser(...)
```

## Implementation Approach

### Phase 1: Basic Support
1. Add `lint_test_files` boolean flag
2. When enabled, apply existing rules to `*_test.go` files
3. Treat test files as belonging to the same layer as their package

### Phase 2: Test-Specific Rules
1. Add `test_rules` configuration section
2. Support exemptions for common test packages
3. Generate test-specific error messages

### Phase 3: Advanced Features
1. Support different rules by test type (unit vs integration)
2. Add test helper detection (builders, factories)
3. Suggest correct patterns in error messages

## Edge Cases to Consider

### 1. Test Packages with `_test` suffix
```go
package mypackage_test  // Black-box testing
```
**Handling**: Should follow rules of the package being tested, not its own layer

### 2. Shared Test Utilities
```go
// internal/testutil/helpers.go
```
**Handling**: Should have its own layer rules or be exempt

### 3. Table-Driven Tests with Inline Functions
```go
tests := []struct{
    setup func() *sql.DB  // Direct SQL reference in test table
}
```
**Handling**: Lint the types used in test structures too

### 4. `//go:build integration` Tags
```go
//go:build integration

package mypackage_test
```
**Handling**: Could allow different rules based on build tags

## Migration Path

### For Existing Projects

```yaml
# .goarchlint - Start conservatively
test_rules:
  enforce_architecture: false  # Don't break existing tests
  warn_only: true              # Show warnings, don't fail CI

  # Gradually tighten
  progressive_mode:
    max_violations: 10         # Fail if violations exceed threshold
    auto_decrease: true        # Decrease threshold over time
```

### For New Projects

```yaml
# .goarchlint - Strict from day one
lint_test_files: true
test_rules:
  enforce_architecture: true
  no_exemptions: true
```

## Alternatives Considered

### Alternative 1: Separate Tool
**Pros**: Focused scope, independent development
**Cons**: Duplicates go-arch-lint logic, extra dependency, inconsistent with main tool

### Alternative 2: Custom Lint Comments
```go
//arch:lint:disable
db, _ := sql.Open(...)
```
**Pros**: Fine-grained control
**Cons**: Easy to abuse, hard to audit, requires discipline

### Alternative 3: Build Tag Exemptions
```go
//go:build !archlint
```
**Pros**: Leverages existing Go tooling
**Cons**: Confusing, hard to discover, requires separate test runs

## Related Work

- **golangci-lint** - Has file-level exemptions but no architectural awareness
- **ArchUnit (Java)** - Tests architecture rules in tests themselves
- **NDepend (.NET)** - Enforces rules on test assemblies

## Questions for Maintainers

1. Would you prefer a simple toggle (`lint_test_files: true`) or more granular control?
2. Should test files be treated as part of their package's layer, or have separate layer classification?
3. What should be the default behavior (`false` for backward compatibility)?
4. Should we support progressive/gradual migration modes?

## Conclusion

Adding test file linting would close a significant gap in go-arch-lint's coverage. While tests require some special handling, they should generally follow the same architectural principles as production code. This feature would help teams maintain consistent, maintainable codebases where tests serve as good examples rather than architectural anti-patterns.

---

**Proposed Configuration (Recommendation)**:
```yaml
# .goarchlint
version: "1.0"

# Simple and clear
lint_test_files: true  # Default: false

# Optional exemptions for common cases
test_exemptions:
  packages:
    - "testing"
    - "github.com/stretchr/testify"

  # Allow per-layer exemptions if needed
  per_layer:
    "cmd":
      warn_only: true  # Start with warnings
```

This strikes a balance between strictness and pragmatism, allowing teams to adopt gradually while providing clear architectural guidance.
