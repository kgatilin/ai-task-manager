package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/repositories"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// Compile-time check that SQLiteAggregateRepository implements repositories.AggregateRepository
var _ repositories.AggregateRepository = (*SQLiteAggregateRepository)(nil)

// SQLiteAggregateRepository implements repositories.AggregateRepository using SQLite as the backend.
type SQLiteAggregateRepository struct {
	DB     *sql.DB
	logger pluginsdk.Logger
}

// NewSQLiteAggregateRepository creates a new SQLite-backed repository.
func NewSQLiteAggregateRepository(db *sql.DB, logger pluginsdk.Logger) *SQLiteAggregateRepository {
	return &SQLiteAggregateRepository{
		DB:     db,
		logger: logger,
	}
}

// ============================================================================
// Aggregate Queries
// ============================================================================

// GetRoadmapWithTracks retrieves a roadmap with all its tracks.
func (r *SQLiteAggregateRepository) GetRoadmapWithTracks(ctx context.Context, roadmapID string) (*entities.RoadmapEntity, error) {
	// Get roadmap
	var roadmap entities.RoadmapEntity
	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id, vision, success_criteria, created_at, updated_at FROM roadmaps WHERE id = ?",
		roadmapID,
	).Scan(&roadmap.ID, &roadmap.Vision, &roadmap.SuccessCriteria, &roadmap.CreatedAt, &roadmap.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: roadmap %s not found", pluginsdk.ErrNotFound, roadmapID)
		}
		return nil, fmt.Errorf("failed to query roadmap: %w", err)
	}

	// Load all tracks for this roadmap
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, roadmap_id, title, description, status, rank, created_at, updated_at FROM tracks WHERE roadmap_id = ? ORDER BY id",
		roadmapID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load tracks: %w", err)
	}
	defer rows.Close()

	// Note: RoadmapEntity doesn't have a Tracks field, so we just verify loading works
	for rows.Next() {
		var (
			id, roadmapID, title, desc, status string
			rank                               int
			createdAt, updatedAt               string
		)
		if err := rows.Scan(&id, &roadmapID, &title, &desc, &status, &rank, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tracks: %w", err)
	}

	return &roadmap, nil
}

// ============================================================================
// Project Metadata Operations
// ============================================================================

// GetProjectMetadata retrieves a metadata value by key.
func (r *SQLiteAggregateRepository) GetProjectMetadata(ctx context.Context, key string) (string, error) {
	var value string
	err := r.DB.QueryRowContext(ctx, "SELECT value FROM project_metadata WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%w: metadata key %s not found", pluginsdk.ErrNotFound, key)
		}
		return "", fmt.Errorf("failed to query metadata: %w", err)
	}
	return value, nil
}

// SetProjectMetadata sets a metadata value by key.
func (r *SQLiteAggregateRepository) SetProjectMetadata(ctx context.Context, key, value string) error {
	_, err := r.DB.ExecContext(
		ctx,
		"INSERT OR REPLACE INTO project_metadata (key, value) VALUES (?, ?)",
		key, value,
	)
	if err != nil {
		return fmt.Errorf("failed to set metadata: %w", err)
	}
	return nil
}

// GetProjectCode retrieves the project code (e.g., "DW" for darwinflow).
// Returns "DW" as default if not set.
func (r *SQLiteAggregateRepository) GetProjectCode(ctx context.Context) string {
	code, err := r.GetProjectMetadata(ctx, "project_code")
	if err != nil {
		// Return default if not set
		return "DW"
	}
	return code
}

// GetNextSequenceNumber retrieves the next sequence number for an entity type.
// Entity types: "task", "track", "iter", "ac", "adr"
func (r *SQLiteAggregateRepository) GetNextSequenceNumber(ctx context.Context, entityType string) (int, error) {
	var maxNum int
	var query string

	switch entityType {
	case "task":
		// Parse existing task IDs to find max number
		query = "SELECT id FROM tasks"
	case "track":
		// Parse existing track IDs to find max number
		query = "SELECT id FROM tracks"
	case "iter":
		// For iterations, use the number column directly
		err := r.DB.QueryRowContext(ctx, "SELECT COALESCE(MAX(number), 0) FROM iterations").Scan(&maxNum)
		if err != nil {
			return 0, fmt.Errorf("failed to get max iteration number: %w", err)
		}
		return maxNum + 1, nil
	case "ac":
		// Parse existing AC IDs to find max number
		query = "SELECT id FROM acceptance_criteria"
	case "adr":
		// Parse existing ADR IDs to find max number
		query = "SELECT id FROM adrs"
	default:
		return 0, fmt.Errorf("%w: invalid entity type: %s", pluginsdk.ErrInvalidArgument, entityType)
	}

	// For tasks and tracks, we need to parse IDs
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to query %s IDs: %w", entityType, err)
	}
	defer rows.Close()

	maxNum = 0
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("failed to scan ID: %w", err)
		}

		// Parse the numeric part from IDs like "DW-task-123" or "DW-track-5"
		// Format: {CODE}-{entity}-{number}
		// Split by "-" and parse the last part
		parts := strings.Split(id, "-")
		if len(parts) >= 3 {
			var num int
			_, err := fmt.Sscanf(parts[len(parts)-1], "%d", &num)
			if err == nil && num > maxNum {
				maxNum = num
			}
		}
	}

	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating IDs: %w", err)
	}

	return maxNum + 1, nil
}
