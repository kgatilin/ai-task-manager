package presenters

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/components"
)

// getStatusStyle returns the appropriate style for a status based on its color name
func getStatusStyle(colorName string) lipgloss.Style {
	switch colorName {
	case "info":
		return components.Styles.StatusTodoStyle
	case "warning":
		return components.Styles.StatusInProgressStyle
	case "success":
		return components.Styles.StatusCompleteStyle
	case "current":
		return components.Styles.StatusCurrentStyle
	case "planned":
		return components.Styles.StatusPlannedStyle
	case "error":
		return components.Styles.StatusBlockedStyle
	case "failed":
		return components.Styles.StatusBlockedStyle
	case "muted":
		return components.Styles.MetadataStyle
	default:
		return lipgloss.NewStyle()
	}
}
