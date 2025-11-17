Document: TM-doc-1763133562553464460
  Title:       TUI Folder Structure Refactoring Specification
  Type:        plan
  Status:      draft
  Created:     2025-11-14 15:19:22
  Updated:     2025-11-14 16:19:27
  Attachment:  Iteration 28

Content:
# Iteration 28 - TUI Folder Structure Refactoring Specification

**Purpose**: Refactor iteration 28 TUI implementation to match planned architecture with proper folder organization

**Reference**: `.agent/details/refactoring/final/IMPLEMENTATION-INDEX.md` Phase 4-7

---

## Current State (Iteration 28 Implementation)

**Location**: `pkg/plugins/task_manager/presentation/cli/`

**Files** (8 files, 3,567 total lines):
```
presentation/cli/
├── command_tui_new.go           (89 lines)   - Command registration
├── tui_new_app.go               (257 lines)  - Root Bubble Tea app
├── tui_new_presenters.go        (1,242 lines) - ALL 5 presenters (monolithic)
├── tui_new_queries.go           (111 lines)  - ALL 3 queries (monolithic)
├── tui_new_transformers.go      (216 lines)  - ALL transformers (monolithic)
├── tui_new_types.go             (244 lines)  - ALL ViewModels (monolithic)
├── tui_new_types_test.go        (742 lines)  - ViewModel tests
└── tui_new_transformers_test.go (755 lines)  - Transformer tests
```

**Problems**:
1. ❌ Wrong location (`cli/` instead of `tui/`)
2. ❌ Monolithic files (all presenters in 1 file, all ViewModels in 1 file)
3. ❌ No components layer (code duplication)
4. ❌ Mixed concerns (ViewModels grouped by implementation, not by feature)

---

## Target State (Planned Architecture)

**Location**: `pkg/plugins/task_manager/presentation/tui/`

**Files** (38 files, organized by layer and feature):

### Layer 1: ViewModels (11 files)
```
presentation/tui/viewmodels/
├── common.go                    # Shared types (if any)
├── loading_vm.go                # LoadingViewModel
├── error_vm.go                  # ErrorViewModel
├── roadmap_list_vm.go           # RoadmapListViewModel + nested (IterationCardVM, TrackCardVM, BacklogTaskVM)
├── iteration_detail_vm.go       # IterationDetailViewModel + nested (ProgressVM, TaskRowVM, IterationACVM)
├── task_detail_vm.go            # TaskDetailViewModel + nested (ACDetailVM, TrackInfoVM, IterationMembershipVM)
├── loading_vm_test.go
├── error_vm_test.go
├── roadmap_list_vm_test.go
├── iteration_detail_vm_test.go
└── task_detail_vm_test.go
```

**Organization Principle**: One file per view type (RoadmapList, IterationDetail, TaskDetail)

**ViewModels to Split**:
- From `tui_new_types.go` (244 lines) → 6 files (~40-50 lines each)

### Layer 2: Transformers (7 files)
```
presentation/tui/transformers/
├── formatting_helpers.go        # RenderStatusBadge(), RenderProgressBar(), FormatDate()
├── roadmap_transformer.go       # TransformToRoadmapListViewModel()
├── iteration_transformer.go     # TransformToIterationDetailViewModel()
├── task_transformer.go          # TransformToTaskDetailViewModel()
├── formatting_helpers_test.go
├── roadmap_transformer_test.go
├── iteration_transformer_test.go
└── task_transformer_test.go
```

**Organization Principle**: One transformer per entity type (Roadmap, Iteration, Task)

**Transformers to Split**:
- From `tui_new_transformers.go` (216 lines) → 4 files (~50-70 lines each)

### Layer 3: Queries (4 files)
```
presentation/tui/queries/
├── roadmap_queries.go           # GetRoadmapViewData() → RoadmapListViewModel
├── iteration_queries.go         # GetIterationDetailViewData() → IterationDetailViewModel
├── task_queries.go              # GetTaskDetailViewData() → TaskDetailViewModel
└── queries_test.go              # Mock repository tests
```

**Organization Principle**: One query file per view's data loading

**Queries to Split**:
- From `tui_new_queries.go` (111 lines) → 3 files (~35-40 lines each)

### Layer 4: Presenters (6 files)
```
presentation/tui/presenters/
├── presenter.go                 # Base Presenter interface
├── loading.go                   # LoadingPresenter
├── error.go                     # ErrorPresenter
├── roadmap_list.go              # RoadmapListPresenter (~250 lines)
├── iteration_detail.go          # IterationDetailPresenter (~350 lines)
├── task_detail.go               # TaskDetailPresenter (~250 lines)
└── keymaps.go                   # All KeyMap definitions
```

**Organization Principle**: One presenter per view (matches ViewModels 1:1)

**Presenters to Split**:
- From `tui_new_presenters.go` (1,242 lines) → 6 files (~200-350 lines each)

### Layer 5: Components (Future - Not in this task)
```
presentation/tui/components/
├── status_badge.go              # Reusable status badge rendering
├── task_list.go                 # Reusable task list component
├── progress_bar.go              # Iteration progress bars
├── box.go                       # Bordered boxes
└── keymaps.go                   # Shared key bindings
```

**Note**: Components extraction is deferred to future iteration (Phase 7 in IMPLEMENTATION-INDEX.md)

### Layer 6: Application (2 files)
```
presentation/tui/
├── app.go                       # Root Bubble Tea app (from tui_new_app.go)
└── command.go                   # TUI command registration (from command_tui_new.go)
```

---

## File-by-File Migration Plan

### Step 1: Create Directory Structure
```bash
mkdir -p pkg/plugins/task_manager/presentation/tui/{viewmodels,transformers,queries,presenters}
```

### Step 2: Move and Split ViewModels

**From**: `presentation/cli/tui_new_types.go` (244 lines)

**To**: 6 files in `presentation/tui/viewmodels/`

| Source Lines | Target File | ViewModels Included | Lines |
|--------------|-------------|---------------------|-------|
| 1-36 | `loading_vm.go` | LoadingViewModel, ErrorViewModel | ~40 |
| 37-100 | `roadmap_list_vm.go` | RoadmapListViewModel, IterationCardVM, TrackCardVM, BacklogTaskVM | ~70 |
| 101-160 | `iteration_detail_vm.go` | ProgressViewModel, TaskRowViewModel, IterationACViewModel, IterationDetailViewModel | ~70 |
| 161-244 | `task_detail_vm.go` | TaskDetailViewModel, ACDetailViewModel, TrackInfoViewModel, IterationMembershipViewModel | ~90 |

**Tests**: Split `tui_new_types_test.go` accordingly

### Step 3: Move and Split Transformers

**From**: `presentation/cli/tui_new_transformers.go` (216 lines)

**To**: 4 files in `presentation/tui/transformers/`

| Source Lines | Target File | Functions Included | Lines |
|--------------|-------------|-------------------|-------|
| 1-50 | `formatting_helpers.go` | FilterActiveIterations(), FilterActiveTracks(), FilterBacklogTasks() | ~50 |
| 51-120 | `roadmap_transformer.go` | TransformToRoadmapListViewModel() | ~70 |
| 121-170 | `iteration_transformer.go` | TransformToIterationDetailViewModel() | ~50 |
| 171-216 | `task_transformer.go` | TransformToTaskDetailViewModel() | ~50 |

**Tests**: Split `tui_new_transformers_test.go` accordingly

### Step 4: Move and Split Queries

**From**: `presentation/cli/tui_new_queries.go` (111 lines)

**To**: 3 files in `presentation/tui/queries/`

| Source Lines | Target File | Functions Included | Lines |
|--------------|-------------|-------------------|-------|
| 1-20 | `queries.go` (common) | QueryService struct definition | ~20 |
| 21-50 | `roadmap_queries.go` | LoadRoadmapListData() | ~35 |
| 51-80 | `iteration_queries.go` | LoadIterationDetailData() | ~35 |
| 81-111 | `task_queries.go` | LoadTaskDetailData() | ~35 |

### Step 5: Move and Split Presenters

**From**: `presentation/cli/tui_new_presenters.go` (1,242 lines)

**To**: 6 files in `presentation/tui/presenters/`

| Source Lines | Target File | Content | Lines |
|--------------|-------------|---------|-------|
| 1-30 | `presenter.go` | Base Presenter interface | ~30 |
| 31-100 | `loading.go` | LoadingPresenter + LoadingKeyMap | ~70 |
| 101-200 | `error.go` | ErrorPresenter + ErrorKeyMap | ~100 |
| 201-600 | `roadmap_list.go` | RoadmapListPresenter + RoadmapListKeyMap | ~400 |
| 601-950 | `iteration_detail.go` | IterationDetailPresenter + IterationDetailKeyMap | ~350 |
| 951-1242 | `task_detail.go` | TaskDetailPresenter + TaskDetailKeyMap | ~292 |

### Step 6: Move Application Files

**From**: `presentation/cli/`

**To**: `presentation/tui/`

| Source File | Target File | Changes |
|-------------|-------------|---------|
| `tui_new_app.go` | `app.go` | Update imports to use `presentation/tui/` packages |
| `command_tui_new.go` | `command.go` | Update imports to use `presentation/tui/` packages |

### Step 7: Update All Imports

**Files to Update**:
1. `presentation/tui/app.go` - Import from `presentation/tui/presenters/`, `presentation/tui/queries/`
2. `presentation/tui/presenters/*.go` - Import from `presentation/tui/viewmodels/`, `presentation/tui/queries/`
3. `presentation/tui/queries/*.go` - Import from `presentation/tui/transformers/`
4. `presentation/tui/transformers/*.go` - Import from `presentation/tui/viewmodels/`, `domain/entities/`
5. `plugin.go` - Update command registration import

**Import Pattern Changes**:
```go
// OLD
import "github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/cli"

// NEW
import "github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui"
import "github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/presenters"
import "github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
```

### Step 8: Delete Old Files

**After verification**:
```bash
rm pkg/plugins/task_manager/presentation/cli/tui_new_*.go
rm pkg/plugins/task_manager/presentation/cli/command_tui_new.go
```

---

## Package Documentation

### presentation/tui/viewmodels/

**Package comment** (add to `common.go` or first file):
```go
// Package viewmodels provides view-specific data structures for the Task Manager TUI.
//
// ViewModels are pure data structures with ZERO business logic. They represent
// pre-computed, display-ready data optimized for rendering in Bubble Tea views.
//
// Architecture Rules:
// - ViewModels have ZERO imports (only stdlib allowed)
// - All transformations done in transformers/ layer (not lazily)
// - Immutable after creation (read-only)
// - Flattened/denormalized for rendering efficiency
//
// Organization:
// - One file per view type (RoadmapList, IterationDetail, TaskDetail)
// - Nested ViewModels in same file as parent (e.g., IterationCardVM in roadmap_list_vm.go)
```

### presentation/tui/transformers/

**Package comment**:
```go
// Package transformers provides pure functions to transform domain entities into ViewModels.
//
// Transformers are responsible for:
// - Entity → ViewModel conversion
// - Pre-computing display fields (status badges, progress bars, formatted dates)
// - Flattening nested structures for rendering efficiency
// - Filtering and grouping data for views
//
// Architecture Rules:
// - Pure functions (no side effects, no repository calls)
// - Imports: domain/entities, presentation/tui/viewmodels
// - One transformer file per entity type (Roadmap, Iteration, Task)
// - formatting_helpers.go contains shared utility functions
```

### presentation/tui/queries/

**Package comment**:
```go
// Package queries provides view-optimized data loading for the TUI.
//
// Query services load domain entities from repositories and transform them
// into ViewModels using transformers. They orchestrate data loading for
// specific views (e.g., RoadmapList, IterationDetail).
//
// Architecture Rules:
// - One query file per view's data loading needs
// - Eliminates N+1 queries by pre-loading related data
// - Returns ViewModels (not entities) to presenters
// - Imports: domain/repositories, transformers, viewmodels
```

### presentation/tui/presenters/

**Package comment**:
```go
// Package presenters implements the MVP (Model-View-Presenter) pattern for the TUI.
//
// Presenters own view state, handle view logic, and prepare data for rendering.
// Each presenter corresponds to one view and implements the Presenter interface.
//
// Architecture Rules:
// - One presenter per view (RoadmapList, IterationDetail, TaskDetail)
// - Presenters call queries (never repositories directly)
// - Presenters work with ViewModels (never entities)
// - Each presenter ~200-350 lines max
// - KeyMaps defined in keymaps.go (shared)
//
// Presenter Interface:
// - Init() tea.Cmd - Returns command to load data
// - Update(msg tea.Msg) (Presenter, tea.Cmd) - Handles messages
// - View() string - Renders ViewModel using Bubbles components
```

---

## Verification Checklist

After refactoring:

### Directory Structure
- [ ] `presentation/tui/viewmodels/` exists with 6+ files
- [ ] `presentation/tui/transformers/` exists with 4 files
- [ ] `presentation/tui/queries/` exists with 3 files
- [ ] `presentation/tui/presenters/` exists with 6 files
- [ ] `presentation/tui/app.go` exists
- [ ] `presentation/tui/command.go` exists
- [ ] Old `presentation/cli/tui_new_*.go` files deleted

### Tests
- [ ] All tests pass: `go test ./pkg/plugins/task_manager/presentation/tui/... -v`
- [ ] Test coverage maintained (100% ViewModels, 100% Transformers)
- [ ] Tests organized alongside code (viewmodels/*_test.go, transformers/*_test.go)

### Linter
- [ ] Zero violations: `go-arch-lint .`
- [ ] Check: `presentation/tui/viewmodels/` has ZERO imports (only stdlib)
- [ ] Check: `presentation/tui/transformers/` imports viewmodels + entities
- [ ] Check: `presentation/tui/presenters/` imports viewmodels + queries
- [ ] Check: No imports of `presentation/cli/` anywhere

### Functionality
- [ ] `dw task-manager tui-new` launches successfully
- [ ] All navigation works (Dashboard → Iteration → Task → back)
- [ ] AC actions work (space verify, s skip, f fail)
- [ ] Iteration reordering works (shift+up/down)
- [ ] No regressions (all features from iteration 28 work)

### Documentation
- [ ] Package comments added to each package (viewmodels, transformers, queries, presenters)
- [ ] README updated if needed (command location unchanged: `dw task-manager tui-new`)

---

## Success Criteria

1. ✅ All 8 monolithic files split into 38 focused files
2. ✅ Correct location: `presentation/tui/` (not `cli/`)
3. ✅ Organized by layer: viewmodels/ transformers/ queries/ presenters/
4. ✅ Organized by feature: One file per view/entity type
5. ✅ All tests pass with same coverage (100% ViewModels, 100% Transformers)
6. ✅ Zero linter violations
7. ✅ Zero regressions (all iteration 28 features work)
8. ✅ Package documentation complete

---

## Out of Scope (Future Iterations)

The following are NOT included in this refactoring task:

1. **Components extraction** - Reusable components (status_badge, task_list, etc.) deferred to separate iteration
2. **New views** - No new views added (ViewTrackDetail, ViewADRList, etc.)
3. **Bubbles component migration** - Using list.Model, table.Model, etc. for existing views (future optimization)
4. **Query optimization** - N+1 query elimination (requires repository changes)

These will be addressed in subsequent iterations per IMPLEMENTATION-INDEX.md phases 7-9.

---

## References

- **Architecture Plan**: `.agent/details/refactoring/final/IMPLEMENTATION-INDEX.md`
- **UI Architecture**: `.agent/details/refactoring/final/ui-architecture.md`
- **Domain-UI Integration**: `.agent/details/refactoring/final/domain-ui-integration.md`
- **Proposed Structure**: `.agent/details/refactoring/final/proposed-structure.md`

