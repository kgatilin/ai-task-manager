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
// - One transformer file per entity type (dashboard, iteration_detail, task_detail)
//
// Organization:
// - dashboard.go: Roadmap/iteration/track transformations + filtering
// - iteration_detail.go: Iteration detail + task grouping transformations
// - task_detail.go: Task detail + AC transformations
package transformers
