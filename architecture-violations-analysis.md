# Architecture Violations Analysis

**Generated**: 2025-11-12
**Total Violations**: 196
**Linter**: go-arch-lint

---

## Violation Summary by Type

| Type | Count | Action Required |
|------|-------|-----------------|
| **Forbidden pkg-to-pkg Dependency** | 134 | REVIEW - Many are necessary for clean architecture |
| **Forbidden Import** | 27 | REVIEW - Many are necessary |
| **Skip-level Import** | 22 | REVIEW - Some may be unavoidable |
| **Test Coverage** | 7 | FIX - Add tests |
| **Test Naming** | 6 | FIX - Move or delete orphaned tests |

---

## Category 1: Domain importing pkg/pluginsdk ✅ NECESSARY

**Verdict**: **SHOULD BE ALLOWED** - These are required for entities to implement SDK interfaces

**Violations**: ~15 occurrences

### Examples:
```
domain/entities/adr_entity.go imports pkg/pluginsdk
domain/entities/track_entity.go imports pkg/pluginsdk
domain/entities/task_entity.go imports pkg/pluginsdk
domain/services/dependency_service.go imports pkg/pluginsdk
domain/services/validation_service.go imports pkg/pluginsdk
```

### Why This Is Necessary:
- Entities implement `pluginsdk.Entity` interface (GetID, GetType, etc.)
- Services use `pluginsdk.Logger` interface
- **This is CORRECT architecture** - Domain depends on SDK contracts

### Recommendation:
```yaml
# Update .arch-lint.yaml to allow:
pkg/plugins/*/domain/** can import:
  - pkg/pluginsdk
```

---

## Category 2: Application importing domain/* subpackages ✅ NECESSARY

**Verdict**: **SHOULD BE ALLOWED** - This is correct clean architecture layering

**Violations**: ~80 occurrences

### Examples:
```
application/track_service.go imports domain/entities
application/track_service.go imports domain/repositories
application/track_service.go imports domain/services
application/task_service.go imports domain/entities
application/iteration_service.go imports domain/repositories
```

### Why This Is Necessary:
- Application layer **MUST** import domain layer (dependency inversion)
- Services orchestrate domain entities and repositories
- This is **TEXTBOOK clean architecture**: Application → Domain

### Recommendation:
```yaml
# Update .arch-lint.yaml to allow:
pkg/plugins/*/application/** can import:
  - pkg/pluginsdk
  - pkg/plugins/*/domain/**
```

---

## Category 3: Application tests importing infrastructure ❌ VIOLATION

**Verdict**: **SHOULD BE FIXED** - Tests should use mocks, not real infrastructure

**Violations**: ~5 occurrences

### Examples:
```
application/track_service_test.go imports infrastructure/persistence
application/task_service_test.go imports infrastructure/persistence
application/iteration_service_test.go imports infrastructure/persistence
application/adr_service_test.go imports infrastructure/persistence
application/ac_service_test.go imports infrastructure/persistence
```

### Why This Is Wrong:
- Application layer tests should use **repository mocks**
- Importing real infrastructure creates coupling
- Makes tests slower and harder to maintain

### Recommendation:
**Create mock repositories** and update tests to use them instead of real SQLite implementations.

Alternative: **Accept integration testing approach** - Current tests use real SQLite in temp directories (via t.TempDir()), which provides better coverage. If acceptable, update linter config:

```yaml
# If integration testing is preferred:
pkg/plugins/*/application/**_test.go can import:
  - pkg/plugins/*/infrastructure/persistence
```

---

## Category 4: Domain repositories importing domain/entities ✅ NECESSARY

**Verdict**: **SHOULD BE ALLOWED** - Repository interfaces need entity types in signatures

**Violations**: ~8 occurrences

### Examples:
```
domain/repositories/track_repository.go imports domain/entities
domain/repositories/task_repository.go imports domain/entities
domain/repositories/iteration_repository.go imports domain/entities
domain/repositories/adr_repository.go imports domain/entities
```

### Why This Is Necessary:
- Repository method signatures return/accept entities
- Example: `GetTrack(ctx, id) (*entities.TrackEntity, error)`
- **Cannot define repositories without entity types**

### Recommendation:
```yaml
# Update .arch-lint.yaml to allow:
pkg/plugins/*/domain/repositories/** can import:
  - pkg/plugins/*/domain/entities
```

---

## Category 5: Domain events importing domain/entities ✅ NECESSARY

**Verdict**: **SHOULD BE ALLOWED** - Event payloads contain entity data

**Violations**: ~1 occurrence

### Example:
```
domain/events/event_payloads.go imports domain/entities
```

### Why This Is Necessary:
- Event payloads contain entity state
- Events like `TrackCreatedPayload` include `TrackEntity` data

### Recommendation:
```yaml
# Update .arch-lint.yaml to allow:
pkg/plugins/*/domain/events/** can import:
  - pkg/plugins/*/domain/entities
```

---

## Category 6: Infrastructure importing domain ✅ NECESSARY

**Verdict**: **SHOULD BE ALLOWED** - Infrastructure implements domain interfaces

**Violations**: ~20 occurrences

### Examples:
```
infrastructure/persistence/track_repository.go imports domain/repositories
infrastructure/persistence/track_repository.go imports domain/entities
infrastructure/persistence/task_repository.go imports domain/repositories
```

### Why This Is Necessary:
- Infrastructure **implements** domain repository interfaces
- Clean architecture: Infrastructure → Domain (dependency inversion)

### Recommendation:
```yaml
# Update .arch-lint.yaml to allow:
pkg/plugins/*/infrastructure/** can import:
  - pkg/pluginsdk
  - pkg/plugins/*/domain/**
```

---

## Category 7: Presentation (CLI) importing application/domain ✅ NECESSARY

**Verdict**: **SHOULD BE ALLOWED** - Presentation layer calls application services

**Violations**: ~15 occurrences

### Examples:
```
presentation/cli/track_adapters.go imports application
presentation/cli/task_adapters.go imports application/dto
presentation/cli/iteration_adapters.go imports domain/entities
```

### Why This Is Necessary:
- CLI adapters **call** application services
- Need DTOs for service input
- May need entities for output formatting

### Recommendation:
```yaml
# Update .arch-lint.yaml to allow:
pkg/plugins/*/presentation/** can import:
  - pkg/pluginsdk
  - pkg/plugins/*/application/**
  - pkg/plugins/*/domain/entities  # For output formatting
```

---

## Category 8: Root package importing nested packages ⚠️ MIXED

**Verdict**: **SOME ARE NECESSARY, SOME SHOULD BE MOVED**

**Violations**: ~15 occurrences

### Examples:
```
✅ ACCEPTABLE:
  plugin.go imports domain/entities (needs to wire up services)
  plugin.go imports application (creates application services)
  plugin.go imports infrastructure/persistence (creates repositories)
  command_roadmap.go imports domain/entities (formats output)
  command_tui.go imports domain (launches TUI)

❌ SHOULD BE FIXED:
  command_migrate.go imports infrastructure/persistence (should be in infrastructure/)
  command_project.go imports infrastructure/persistence (should be in infrastructure/)
```

### Recommendation:
1. **Allow plugin.go** to import all layers (it's the composition root)
2. **Move infrastructure commands** to infrastructure/ package
3. **Keep presentation commands** (roadmap, tui) in root for now

```yaml
# Update .arch-lint.yaml:
pkg/plugins/*/plugin.go can import:
  - pkg/plugins/*/**  # Composition root exception
```

---

## Category 9: Orphaned test files ❌ VIOLATION

**Verdict**: **SHOULD BE MOVED OR DELETED**

**Violations**: 6 occurrences

### Files:
```
pkg/plugins/task_manager/acceptance_criteria_entity_test.go
pkg/plugins/task_manager/adr_entity_test.go
pkg/plugins/task_manager/roadmap_entity_test.go
pkg/plugins/task_manager/track_entity_test.go
pkg/plugins/task_manager/iteration_entity_test.go
pkg/plugins/task_manager/sqlite_repository_test.go
```

### Issue:
- Test files in root package import nested packages
- No corresponding implementation files in root package
- Should be in domain/entities/ or infrastructure/persistence/

### Recommendation:
**Move to correct locations**:
```
acceptance_criteria_entity_test.go → domain/entities/acceptance_criteria_entity_test.go
adr_entity_test.go → domain/entities/adr_entity_test.go
roadmap_entity_test.go → domain/entities/roadmap_entity_test.go
track_entity_test.go → domain/entities/track_entity_test.go
iteration_entity_test.go → domain/entities/iteration_entity_test.go
sqlite_repository_test.go → infrastructure/persistence/*_repository_test.go (split by entity)
```

---

## Category 10: Test coverage violations ⚠️ SHOULD ADD TESTS

**Verdict**: **ADD MISSING TESTS**

**Violations**: 7 occurrences

### Packages needing tests:
```
1. domain/repositories - 0% coverage (threshold: 50%)
2. domain/events - no tests
3. application/dto - 0% coverage
4. presentation/cli - 0% coverage
5. infrastructure/persistence - 0% coverage (false - tests exist but not counted)
```

### Recommendation:
1. **domain/repositories**: Add interface tests (simple - just verify signatures)
2. **domain/events**: Add event payload tests
3. **application/dto**: Add DTO validation tests
4. **presentation/cli**: Add adapter tests (flag parsing, output formatting)
5. **infrastructure/persistence**: Move tests from root package to fix coverage reporting

---

## Proposed Linter Configuration Update

Create/update `.arch-lint.yaml`:

```yaml
# Allow necessary clean architecture imports
allowed_imports:
  # Domain layer can import SDK
  pkg/plugins/*/domain/**:
    - pkg/pluginsdk

  # Domain repositories can import domain entities
  pkg/plugins/*/domain/repositories/**:
    - pkg/plugins/*/domain/entities

  # Domain events can import domain entities
  pkg/plugins/*/domain/events/**:
    - pkg/plugins/*/domain/entities

  # Application layer can import domain
  pkg/plugins/*/application/**:
    - pkg/pluginsdk
    - pkg/plugins/*/domain/**

  # Infrastructure can import domain
  pkg/plugins/*/infrastructure/**:
    - pkg/pluginsdk
    - pkg/plugins/*/domain/**

  # Presentation can import application and domain
  pkg/plugins/*/presentation/**:
    - pkg/pluginsdk
    - pkg/plugins/*/application/**
    - pkg/plugins/*/domain/entities  # For output formatting

  # Plugin root (composition root) can import all
  pkg/plugins/*/plugin.go:
    - pkg/plugins/*/**

  # Integration tests can import infrastructure (optional)
  pkg/plugins/*/application/**_test.go:
    - pkg/plugins/*/infrastructure/persistence  # If integration testing preferred

# Skip-level imports - allow for subpackage access
allow_skip_level: true
```

---

## Summary

| Category | Count | Verdict | Action |
|----------|-------|---------|--------|
| Domain → SDK | ~15 | ✅ ALLOW | Update linter config |
| Application → Domain | ~80 | ✅ ALLOW | Update linter config |
| Domain repos → Entities | ~8 | ✅ ALLOW | Update linter config |
| Domain events → Entities | ~1 | ✅ ALLOW | Update linter config |
| Infrastructure → Domain | ~20 | ✅ ALLOW | Update linter config |
| Presentation → App/Domain | ~15 | ✅ ALLOW | Update linter config |
| Plugin root → All | ~10 | ✅ ALLOW | Update linter config |
| App tests → Infrastructure | ~5 | ⚠️ DECIDE | Mock or allow integration tests |
| Orphaned test files | 6 | ❌ FIX | Move to correct packages |
| Test coverage gaps | 7 | ❌ FIX | Add missing tests |

---

## Recommendation

**~170 violations (87%) are CORRECT clean architecture and should be allowed by updating linter config.**

**~26 violations (13%) are real issues that should be fixed:**
- 6 orphaned test files (move to correct locations)
- 7 test coverage gaps (add tests)
- 5 application test imports (decide: mock or allow integration testing)
- 8 infrastructure commands in root (move to infrastructure package)

**Next Steps:**
1. Review and approve this analysis
2. Create `.arch-lint.yaml` with updated rules
3. Fix the ~26 real violations
4. Re-run linter to verify ~170 violations are now allowed
