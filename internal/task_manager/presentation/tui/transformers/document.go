package transformers

import (
	"strings"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

// TransformDocumentsToListItems converts a slice of DocumentEntity to DocumentListItemViewModel.
// Pure function - processes all entities and returns list items.
func TransformDocumentsToListItems(entities []*entities.DocumentEntity) []viewmodels.DocumentListItemViewModel {
	if len(entities) == 0 {
		return []viewmodels.DocumentListItemViewModel{}
	}

	items := make([]viewmodels.DocumentListItemViewModel, 0, len(entities))
	for _, entity := range entities {
		if entity == nil {
			continue
		}

		item := viewmodels.DocumentListItemViewModel{
			ID:         entity.ID,
			Title:      entity.Title,
			Type:       formatDocumentType(entity.Type),
			StatusIcon: getDocumentStatusIcon(entity.Status),
		}
		items = append(items, item)
	}

	return items
}

// formatDocumentType converts DocumentType value object to display label.
// Examples: "adr" -> "ADR", "plan" -> "Plan", "retrospective" -> "Retrospective"
func formatDocumentType(docType entities.DocumentType) string {
	switch docType {
	case entities.DocumentTypeADR:
		return "ADR"
	case entities.DocumentTypePlan:
		return "Plan"
	case entities.DocumentTypeRetrospective:
		return "Retrospective"
	case entities.DocumentTypeOther:
		return "Other"
	default:
		return strings.ToUpper(string(docType))
	}
}

// getDocumentStatusIcon returns visual icon for document status.
// ✓ for published/approved documents, ○ for draft, ✗ for archived.
func getDocumentStatusIcon(status entities.DocumentStatus) string {
	switch status {
	case entities.DocumentStatusPublished:
		return "✓"
	case entities.DocumentStatusDraft:
		return "○"
	case entities.DocumentStatusArchived:
		return "✗"
	default:
		return "?"
	}
}

// TransformToDocumentViewModel transforms a document entity to a document view model
func TransformToDocumentViewModel(doc *entities.DocumentEntity) *viewmodels.DocumentViewModel {
	vm := viewmodels.NewDocumentViewModel(
		doc.ID,
		doc.Title,
		string(doc.Type),
		string(doc.Status),
		doc.Content,
		doc.TrackID,
		doc.IterationNumber,
		formatTime(doc.CreatedAt),
		formatTime(doc.UpdatedAt),
	)

	// Pre-compute display fields
	vm.StatusLabel = GetDocumentStatusLabel(string(doc.Status))
	vm.StatusColor = GetDocumentStatusColor(string(doc.Status))
	vm.TypeLabel = GetDocumentTypeLabel(string(doc.Type))
	vm.Icon = GetDocumentStatusIcon(string(doc.Status))

	return vm
}

// GetDocumentStatusLabel returns human-readable status label
func GetDocumentStatusLabel(status string) string {
	switch status {
	case string(entities.DocumentStatusDraft):
		return "Draft"
	case string(entities.DocumentStatusPublished):
		return "Published"
	case string(entities.DocumentStatusArchived):
		return "Archived"
	default:
		return status
	}
}

// GetDocumentStatusColor returns lipgloss color name for status
func GetDocumentStatusColor(status string) string {
	switch status {
	case string(entities.DocumentStatusDraft):
		return "11" // Yellow for draft
	case string(entities.DocumentStatusPublished):
		return "10" // Green for published
	case string(entities.DocumentStatusArchived):
		return "240" // Gray for archived
	default:
		return "240"
	}
}

// GetDocumentTypeLabel returns human-readable type label
func GetDocumentTypeLabel(docType string) string {
	switch docType {
	case string(entities.DocumentTypeADR):
		return "ADR"
	case string(entities.DocumentTypePlan):
		return "Plan"
	case string(entities.DocumentTypeRetrospective):
		return "Retrospective"
	case string(entities.DocumentTypeOther):
		return "Other"
	default:
		return docType
	}
}

// GetDocumentStatusIcon returns status icon
func GetDocumentStatusIcon(status string) string {
	switch status {
	case string(entities.DocumentStatusDraft):
		return "📝" // Draft icon
	case string(entities.DocumentStatusPublished):
		return "✓" // Check mark for published
	case string(entities.DocumentStatusArchived):
		return "📦" // Archive icon
	default:
		return "📄"
	}
}

// formatTime formats a time.Time to a human-readable string
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}
