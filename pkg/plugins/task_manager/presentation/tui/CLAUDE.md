# Package: presentation/tui

**Path**: `pkg/plugins/task_manager/presentation/tui`

**Role**: Terminal User Interface using MVP pattern with strict layer separation

**Framework**: Bubble Tea + Bubbles + Lipgloss

---

## Architecture: MVP with 4 Layers

```
ViewModels (Pure Data)           # Zero imports, immutable, display-ready
    ↑ created by
Transformers (Entity→ViewModel)  # Pure functions, pre-compute display data
    ↑ called by
Queries (Fetch + Transform)      # Orchestrate loading, eliminate N+1
    ↑ called by
Presenters (UI Rendering)        # Bubble Tea components, keyboard handling
    ↑ composed by
App (State Machine)              # Navigation, view transitions
```

**Key Principle**: Data flows down (dependencies), events flow up (messages). Each layer has ZERO knowledge of layers above it.

---

## Layer Responsibilities

### ViewModels (`viewmodels/*.go`)
- **Pure data structures** with ZERO logic
- Only stdlib imports (no domain, no tea package)
- All fields public, pre-computed (no lazy evaluation)
- Immutable after construction
- **See**: `viewmodels/doc.go` for architecture rules

### Transformers (`transformers/*.go`)
- **Pure functions**: Entity → ViewModel
- Import: domain/entities, viewmodels (NEVER presenters)
- Pre-compute ALL display data (status labels, icons, timestamps, progress bars)
- No side effects (no repository calls, no logging)
- **See**: `transformers/doc.go` for patterns

### Queries (`queries/*.go`)
- **Orchestrate data loading** for specific views
- Fetch all data at once (eliminate N+1 queries)
- Transform entities → ViewModels via transformers
- Return complete ViewModel ready for rendering
- **See**: `queries/doc.go` for query patterns

### Presenters (`presenters/*.go`)
- **Bubble Tea components** implementing Presenter interface
- Model = ViewModel + UI state (selectedIndex, expandedItems, activeTab)
- Init(): Load data via query
- Update(msg): Handle input → update UI state OR reload data
- View(): Render ViewModel using components (NEVER transform data)
- **See**: `presenters/doc.go`, `presenters/base.go` for interface

### App (`app.go`)
- **Root Bubble Tea model** (state machine)
- Manages view transitions via ViewStateNew enum
- Tracks navigation state (previousView, currentIterationNumber, currentTaskID)
- Delegates Update/View to active presenter
- Handles global keys (q=quit, esc=back)

---

## Component Reuse

### Bubbles Components
- **Spinner** (`components/spinner.go`): Wrapper with consistent styling
- **Help** (`components/help.go`): Wrapper with centralized key bindings
- **KeyMap** (`components/keybindings.go`): Centralized key definitions, progressive help

### Custom Components
- **ScrollHelper** (`components/scroll_helper.go`): Auto-scroll for long lists (keep selected item in view)
- **Styles** (`components/styles.go`): Centralized lipgloss styles (single source of truth)

**Rule**: All presenters use `components.Styles.*` for styling. Never create lipgloss styles in presenters.

**See**: `components/doc.go` for color scheme and style patterns

---

## Key Patterns

### Initialization
- Set default width/height (80x24) to prevent 0-width wrapping
- Request `tea.WindowSize()` from Init()
- Handle `tea.WindowSizeMsg` in Update()

### Text Wrapping
- Use `wordwrap.String()` + `indent.String()` from github.com/muesli/reflow
- `indent.String()` adds spaces after EVERY newline (not just first)

### Progressive Help
- ShortHelp: 3-4 essential keys (default)
- FullHelp: All keys (toggle with '?')
- Context-aware: Different keys for different tabs/views

### Error Handling
- Presenters return `presenters.ErrorMsg{Err: err}`
- App shows error view
- Track previousView before error (enable escape to navigate back)

### Data Reloading
- After mutations: reload data via query
- Preserve selection: pass selectedIndex in loaded message
- Restore selection after reload

### Message Passing
- Define custom messages in `presenters/messages.go`
- Presenters return messages for navigation (DrillIntoIterationMsg, BackMsg)
- App handles messages and decides next view

**See**: Existing presenter files for implementation examples

---

## Navigation Flow

```
Dashboard (ViewRoadmapListNew)
  ├─ Enter on iteration → IterationDetail
  │   ├─ Enter on task → TaskDetail
  │   │   └─ Esc → back to IterationDetail
  │   └─ Esc → back to Dashboard
  └─ Esc → Exit
```

**Error Recovery**: Escape from error view returns to previous view (not quit)

**State Tracking**: App tracks currentIterationNumber, currentTaskID for navigation context

---

## Testing Strategy

- **ViewModels**: Pure unit tests (100% coverage target)
- **Transformers**: Entity → ViewModel tests (100% coverage target)
- **Queries**: Mock repositories (70-80% coverage)
- **Presenters**: Integration tests (focus on key handling, state transitions)

**See**: `*_test.go` files in each package

---

## Anti-Patterns

### ❌ Computing Display Data in Presenters
- ViewModels must have pre-computed fields (StatusLabel, Icon, Progress)
- Presenters render strings, never transform data

### ❌ Transformers Calling Repositories
- Transformers are pure functions (no side effects)
- Queries call repositories, pass entities to transformers

### ❌ ViewModels with Methods
- ViewModels are data structures (no logic)
- Pre-compute in transformers, store in fields

### ❌ Presenters Importing Domain Entities
- Presenters work with ViewModels only
- Never import domain/entities in presenters

### ❌ Creating Styles in Presenters
- Use `components.Styles.*` for all styling
- Never create lipgloss styles inline

---

## Key Decisions

### Why MVP over MVVM?
- Simpler: Explicit UI state management (no two-way binding)
- Testable: Deterministic message passing
- Bubble Tea-friendly: Update(msg) → (model, cmd) maps naturally

### Why 4 Layers?
- Eliminates N+1 queries (queries fetch all data at once)
- Thin presenters (one query call, complete ViewModel)
- Easy to optimize (add caching, batching to queries layer)

### Why Pre-compute Display Data?
- Pure rendering in presenters (no conditional logic)
- Testable transformers (pure functions)
- Performance: Compute once, render many times

---

## Common Tasks

### Adding New View
1. Create ViewModel (`viewmodels/new_view.go`)
2. Create Transformer (`transformers/new_view.go`)
3. Create Query (`queries/new_view.go`)
4. Create Presenter (`presenters/new_view.go`)
5. Add ViewState enum to `app.go`
6. Handle navigation messages in `app.go`

### Adding New Action
1. Define message type (`presenters/messages.go`)
2. Handle key in presenter Update()
3. Update KeyMap (`components/keybindings.go`)
4. Update help text in presenter View()

---

## References

- **Bubble Tea**: https://github.com/charmbracelet/bubbletea
- **Bubbles**: https://github.com/charmbracelet/bubbles (spinner, help, key)
- **Lipgloss**: https://github.com/charmbracelet/lipgloss (styling)
- **Reflow**: https://github.com/muesli/reflow (wordwrap, indent)
- **Package docs**: Read `*/doc.go` in each subdirectory for detailed patterns
- **Examples**: Read existing presenter files for implementation patterns
- **Parent architecture**: `../CLAUDE.md` (task_manager plugin structure)

---

**Last Updated**: 2025-11-17 (Iteration 28 - MVP pattern with layered architecture)
