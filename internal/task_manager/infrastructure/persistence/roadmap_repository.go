package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
)

// Compile-time check that SQLiteRoadmapRepository implements repositories.RoadmapRepository
var _ repositories.RoadmapRepository = (*SQLiteRoadmapRepository)(nil)

// SQLiteRoadmapRepository implements repositories.RoadmapRepository using SQLite as the backend.
type SQLiteRoadmapRepository struct {
	DB     *sql.DB
	logger logger.Logger
}

// NewSQLiteRoadmapOnlyRepository creates a new SQLite-backed roadmap-focused repository.
// This is internal to the persistence layer and not exported for general use.
// Use NewSQLiteRepositoryComposite for the full repository interface.
func NewSQLiteRoadmapOnlyRepository(db *sql.DB, logger logger.Logger) *SQLiteRoadmapRepository {
	return &SQLiteRoadmapRepository{
		DB:     db,
		logger: logger,
	}
}

// ============================================================================
// Roadmap Operations
// ============================================================================

// SaveRoadmap persists a new roadmap to storage.
func (r *SQLiteRoadmapRepository) SaveRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error {
	// Check if roadmap already exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM roadmaps WHERE id = ?", roadmap.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check roadmap existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w: roadmap %s already exists", tmerrors.ErrAlreadyExists, roadmap.ID)
	}

	_, err = r.DB.ExecContext(
		ctx,
		"INSERT INTO roadmaps (id, vision, success_criteria, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		roadmap.ID, roadmap.Vision, roadmap.SuccessCriteria, roadmap.CreatedAt, roadmap.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert roadmap: %w", err)
	}

	return nil
}

// GetRoadmap retrieves a roadmap by its ID.
func (r *SQLiteRoadmapRepository) GetRoadmap(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
	var roadmap entities.RoadmapEntity

	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id, vision, success_criteria, created_at, updated_at FROM roadmaps WHERE id = ?",
		id,
	).Scan(&roadmap.ID, &roadmap.Vision, &roadmap.SuccessCriteria, &roadmap.CreatedAt, &roadmap.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: roadmap %s not found", tmerrors.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to query roadmap: %w", err)
	}

	return &roadmap, nil
}

// GetActiveRoadmap retrieves the most recently created roadmap.
func (r *SQLiteRoadmapRepository) GetActiveRoadmap(ctx context.Context) (*entities.RoadmapEntity, error) {
	var roadmap entities.RoadmapEntity

	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id, vision, success_criteria, created_at, updated_at FROM roadmaps ORDER BY created_at DESC LIMIT 1",
	).Scan(&roadmap.ID, &roadmap.Vision, &roadmap.SuccessCriteria, &roadmap.CreatedAt, &roadmap.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: no active roadmap found", tmerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to query active roadmap: %w", err)
	}

	return &roadmap, nil
}

// UpdateRoadmap updates an existing roadmap.
func (r *SQLiteRoadmapRepository) UpdateRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error {
	result, err := r.DB.ExecContext(
		ctx,
		"UPDATE roadmaps SET vision = ?, success_criteria = ?, updated_at = ? WHERE id = ?",
		roadmap.Vision, roadmap.SuccessCriteria, roadmap.UpdatedAt, roadmap.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update roadmap: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: roadmap %s not found", tmerrors.ErrNotFound, roadmap.ID)
	}

	return nil
}
