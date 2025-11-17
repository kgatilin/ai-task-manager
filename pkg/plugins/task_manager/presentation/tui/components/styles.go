package components

import "github.com/charmbracelet/lipgloss"

// ColorScheme defines all colors used in the TUI.
// This is the single source of truth for the color palette.
var ColorScheme = struct {
	Accent       string // "205" - Primary accent (magenta/pink)
	ErrorTitle   string // "196" - Error red
	ErrorMessage string // "203" - Error message pink
	Muted        string // "240" - Muted gray for metadata
	SectionTitle string // "cyan" - Section headers
	Success      string // "green" - Progress/success indicators
}{
	Accent:       "205",
	ErrorTitle:   "196",
	ErrorMessage: "203",
	Muted:        "240",
	SectionTitle: "cyan",
	Success:      "green",
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
}
