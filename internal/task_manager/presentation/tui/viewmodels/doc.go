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
// - One file per view type (loading, error, dashboard, iteration_detail, task_detail)
// - Nested ViewModels in same file as parent (e.g., IterationCardViewModel in dashboard.go)
package viewmodels
