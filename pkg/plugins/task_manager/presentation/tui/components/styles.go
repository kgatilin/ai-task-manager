package components

import "github.com/charmbracelet/lipgloss"

// Icons defines all icons used in the TUI
// This is the single source of truth for icon constants
var Icons = struct {
	// Iteration status icons
	IterationPlanned  string
	IterationCurrent  string
	IterationComplete string

	// Task status icons
	TaskTodo       string
	TaskInProgress string
	TaskReview     string
	TaskDone       string
	TaskCancelled  string

	// Track status icons
	TrackNotStarted string
	TrackInProgress string
	TrackComplete   string
	TrackBlocked    string
	TrackWaiting    string

	// AC status icons (same as entity StatusIndicator)
	ACNotStarted            string
	ACVerified              string
	ACAutomaticallyVerified string
	ACPendingReview         string
	ACFailed                string
	ACSkipped               string
}{
	// Iteration status icons
	IterationPlanned:  "📋",
	IterationCurrent:  "▶",
	IterationComplete: "✓",

	// Task status icons
	TaskTodo:       "○",
	TaskInProgress: "◐",
	TaskReview:     "◑",
	TaskDone:       "●",
	TaskCancelled:  "⊗",

	// Track status icons
	TrackNotStarted: "○",
	TrackInProgress: "◐",
	TrackComplete:   "●",
	TrackBlocked:    "⊠",
	TrackWaiting:    "⏸",

	// AC status icons (matches entities.AcceptanceCriteriaEntity.StatusIndicator())
	ACNotStarted:            "○",
	ACVerified:              "✓",
	ACAutomaticallyVerified: "✓",
	ACPendingReview:         "⏸",
	ACFailed:                "✗",
	ACSkipped:               "⊘",
}

// ColorScheme defines all colors used in the TUI.
// This is the single source of truth for the color palette.
var ColorScheme = struct {
	Accent       string // "205" - Primary accent (magenta/pink)
	ErrorTitle   string // "196" - Error red
	ErrorMessage string // "203" - Error message pink
	Muted        string // "240" - Muted gray for metadata
	SectionTitle string // "cyan" - Section headers
	Success      string // "green" - Progress/success indicators
	Warning      string // "yellow" - Warning/pending states
	Info         string // "blue" - Informational states
	Failed       string // "160" - Failed/error states (dark red)
	Skipped      string // "gray" - Skipped/disabled states
	Current      string // "magenta" - Current/active iteration highlight
}{
	Accent:       "205",
	ErrorTitle:   "196",
	ErrorMessage: "203",
	Muted:        "240",
	SectionTitle: "cyan",
	Success:      "green",
	Warning:      "yellow",
	Info:         "blue",
	Failed:       "160",
	Skipped:      "gray",
	Current:      "magenta",
}

// Styles contains all pre-defined lipgloss styles used across the TUI.
// This is the single source of truth for styling.
// Each style is initialized once at package load time.
var Styles = struct {
	// General styles
	TitleStyle    lipgloss.Style // Bold + accent color
	SectionStyle  lipgloss.Style // Bold + cyan
	MetadataStyle lipgloss.Style // Muted gray
	SelectedStyle lipgloss.Style // Accent color

	// Error view styles
	ErrorTitleStyle   lipgloss.Style // Bold + error red
	ErrorMessageStyle lipgloss.Style // Error message pink
	ErrorDetailsStyle lipgloss.Style // Muted + italic

	// Loading view styles
	LoadingStyle lipgloss.Style // Bold + accent

	// Task/AC styles
	ProgressStyle lipgloss.Style // Green for progress
	TestingStyle  lipgloss.Style // Muted + italic for testing instructions

	// Tab styles
	TabStyle       lipgloss.Style // Bold for inactive tab
	ActiveTabStyle lipgloss.Style // Bold + underline + accent for active tab

	// Component accent style
	AccentStyle lipgloss.Style // Accent color (for spinner, etc.)

	// Status-specific styles
	StatusPlannedStyle      lipgloss.Style // Planned iteration (info blue)
	StatusCurrentStyle      lipgloss.Style // Current iteration (bold magenta)
	StatusCompleteStyle     lipgloss.Style // Complete iteration (green)
	StatusTodoStyle         lipgloss.Style // Todo task (info blue)
	StatusInProgressStyle   lipgloss.Style // In-progress task (warning yellow)
	StatusReviewStyle       lipgloss.Style // Review task (warning yellow)
	StatusDoneStyle         lipgloss.Style // Done task (success green)
	StatusNotStartedStyle   lipgloss.Style // Not started track (muted gray)
	StatusBlockedStyle      lipgloss.Style // Blocked track (failed red)
	StatusWaitingStyle      lipgloss.Style // Waiting track (warning yellow)
	ACVerifiedStyle         lipgloss.Style // Verified AC (success green)
	ACFailedStyle           lipgloss.Style // Failed AC (failed red + bold)
	ACPendingStyle          lipgloss.Style // Pending AC (warning yellow)
	ACSkippedStyle          lipgloss.Style // Skipped AC (skipped gray)
}{
	TitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorScheme.Accent)),

	SectionStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorScheme.SectionTitle)),

	MetadataStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Muted)),

	SelectedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Accent)),

	ErrorTitleStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorScheme.ErrorTitle)),

	ErrorMessageStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.ErrorMessage)),

	ErrorDetailsStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Muted)).
		Italic(true),

	LoadingStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorScheme.Accent)),

	ProgressStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Success)),

	TestingStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Muted)).
		Italic(true),

	TabStyle: lipgloss.NewStyle().
		Bold(true),

	ActiveTabStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorScheme.Accent)).
		Underline(true),

	AccentStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Accent)),

	// Status-specific styles
	StatusPlannedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Info)),

	StatusCurrentStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorScheme.Current)),

	StatusCompleteStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Success)),

	StatusTodoStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Info)),

	StatusInProgressStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Warning)),

	StatusReviewStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Warning)),

	StatusDoneStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Success)),

	StatusNotStartedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Muted)),

	StatusBlockedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Failed)),

	StatusWaitingStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Warning)),

	ACVerifiedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Success)),

	ACFailedStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorScheme.Failed)),

	ACPendingStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Warning)),

	ACSkippedStyle: lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorScheme.Skipped)),
}
