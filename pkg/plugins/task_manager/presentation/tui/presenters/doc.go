// Package presenters implements the MVP (Model-View-Presenter) pattern for the TUI.
//
// Presenters own view state, handle view logic, and prepare data for rendering.
// Each presenter corresponds to one view and implements the Presenter interface.
//
// Architecture Rules:
// - One presenter per view (loading, error, dashboard, iteration_detail, task_detail)
// - Presenters call queries (never repositories directly)
// - Presenters work with ViewModels (never entities)
// - Each presenter ~200-350 lines max
// - KeyMaps defined separately for each presenter
//
// Presenter Interface:
// - Init() tea.Cmd - Returns command to load data
// - Update(msg tea.Msg) (Presenter, tea.Cmd) - Handles messages
// - View() string - Renders ViewModel using Bubbles components
package presenters
