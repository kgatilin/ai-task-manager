// Package queries provides view-optimized data loading for the TUI.
//
// Query services load domain entities from repositories and transform them
// into ViewModels using transformers. They orchestrate data loading for
// specific views (e.g., RoadmapList, IterationDetail, TaskDetail).
//
// Architecture Rules:
// - One query file per view's data loading needs
// - Eliminates N+1 queries by pre-loading related data
// - Returns ViewModels (not entities) to presenters
// - Imports: domain/repositories, presentation/tui/transformers, presentation/tui/viewmodels
//
// Organization:
// - dashboard.go: Loads data for dashboard/roadmap list view
// - iteration_detail.go: Loads data for iteration detail view
// - task_detail.go: Loads data for task detail view
package queries
