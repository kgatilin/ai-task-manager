// Package components provides reusable TUI components with centralized styling.
//
// This package wraps Bubble Tea components (spinner, help, key bindings) with
// project-specific styling and theming. It provides:
//   - Centralized color scheme and styles
//   - Consistent component behavior across views
//   - Abstraction layer over external TUI framework
//
// Architecture Rules:
//   - Components are thin wrappers around Bubble Tea
//   - All styling defined in styles.go (single source of truth)
//   - No business logic (only UI rendering)
//   - Imports: Bubble Tea components, lipgloss (no domain layer)
//
// Usage:
//   - Presenters import this package instead of bubbles/* directly
//   - Use Styles.* for all lipgloss styling
//   - Use component constructors (NewSpinner, NewHelp) for consistent setup
//
// Color Scheme:
//   - Accent (205): Primary magenta/pink - titles, selected items
//   - ErrorTitle (196): Error red - error titles
//   - ErrorMessage (203): Error message pink - error messages
//   - Muted (240): Muted gray - metadata, details, testing instructions
//   - SectionTitle (cyan): Section headers and dividers
//   - Success (green): Progress indicators, success states
//
// Style Patterns:
//   - TitleStyle: Bold + accent (page titles, major headings)
//   - SectionStyle: Bold + cyan (section headers like "Active Iterations")
//   - MetadataStyle: Muted gray (timestamps, IDs, secondary info)
//   - SelectedStyle: Accent color (highlighted items)
//   - ErrorTitleStyle: Bold + error red (error headers)
//   - ErrorMessageStyle: Error pink (main error messages)
//   - ErrorDetailsStyle: Muted + italic (error context)
//   - LoadingStyle: Bold + accent (loading messages)
//   - ProgressStyle: Green (progress bars, success indicators)
//   - TestingStyle: Muted + italic (testing instructions, AC details)
//   - AccentStyle: Accent color (component accents, spinners)
package components
