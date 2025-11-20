package queries

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/transformers"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

// LoadDocumentData loads document data for a specific document.
// Returns document entity transformed into view model ready for presentation.
//
// Pre-loads:
// - Document entity by ID
//
// Eliminates N+1 queries by loading all required data upfront.
func LoadDocumentData(
	ctx context.Context,
	repo repositories.DocumentRepository,
	documentID string,
) (*viewmodels.DocumentViewModel, error) {
	// Fetch document
	doc, err := repo.FindDocumentByID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// Transform to view model
	vm := transformers.TransformToDocumentViewModel(doc)

	return vm, nil
}
