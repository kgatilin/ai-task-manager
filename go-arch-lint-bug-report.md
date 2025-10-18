# go-arch-lint Bug Report

## Summary
`go-arch-lint` reports violations for imports that are explicitly allowed in the configuration. The error message itself is contradictory, stating that a package violates a rule while simultaneously showing that the rule permits the import.

## Configuration
File: `.goarchlint`
```yaml
rules:
    directories_import:
        cmd:
            - internal/app
            - internal/infra
        cmd/dw:
            - internal/app
            - internal/infra
        internal/app:
            - internal/domain
        internal/domain: []
        internal/infra:
            - internal/domain
```

## Actual Imports
File: `cmd/dw/claude.go`
```go
import (
    "github.com/kgatilin/darwinflow-pub/internal/app"
    "github.com/kgatilin/darwinflow-pub/internal/infra"
)
```

File: `cmd/dw/logs.go`
```go
import (
    "github.com/kgatilin/darwinflow-pub/internal/app"
)
```

## Expected Behavior
These imports should be allowed since:
1. The configuration explicitly allows `cmd` to import from `internal/app` and `internal/infra`
2. The configuration also explicitly allows `cmd/dw` to import from `internal/app` and `internal/infra`
3. The actual imports match the allowed list exactly

## Actual Behavior
```
[ERROR] Forbidden Import
  File: cmd/dw/claude.go
  Issue: cmd/dw imports internal/app
  Rule: cmd can only import from: [internal/app internal/infra]
  Fix: Restructure dependencies according to allowed imports

[ERROR] Forbidden Import
  File: cmd/dw/claude.go
  Issue: cmd/dw imports internal/infra
  Rule: cmd can only import from: [internal/app internal/infra]
  Fix: Restructure dependencies according to allowed imports

[ERROR] Forbidden Import
  File: cmd/dw/logs.go
  Issue: cmd/dw imports internal/app
  Rule: cmd can only import from: [internal/app internal/infra]
  Fix: Restructure dependencies according to allowed imports
```

## Analysis
The error message is contradictory:
- **Issue**: "cmd/dw imports internal/app"
- **Rule**: "cmd can only import from: [internal/app internal/infra]"

The linter correctly identifies that the rule allows imports from `[internal/app internal/infra]`, but then incorrectly reports importing `internal/app` as a violation.

## Possible Root Causes
1. **Directory matching issue**: The linter may not correctly match subdirectories like `cmd/dw` to parent rules like `cmd`
2. **Logic inversion bug**: The validation logic may have an inverted condition (e.g., checking if import is IN allowed list vs NOT IN allowed list)
3. **Rule precedence**: If both `cmd` and `cmd/dw` rules exist, there may be confusion about which rule applies

## Reproduction
1. Create a Go project with the module path `github.com/kgatilin/darwinflow-pub`
2. Create the directory structure:
   - `cmd/dw/` (with Go files)
   - `internal/app/` (with Go files)
   - `internal/infra/` (with Go files)
   - `internal/domain/` (with Go files)
3. Add imports from `cmd/dw/*.go` to `internal/app` and `internal/infra`
4. Create `.goarchlint` with the configuration shown above
5. Run `go-arch-lint .`

## Environment
- Go version: 1.25.1
- Module: github.com/kgatilin/darwinflow-pub
- go-arch-lint: (version unknown, please check with `go-arch-lint --version`)

## Workaround Attempted
Tried adding explicit `cmd/dw` rule in addition to `cmd` rule, but same error persists.
