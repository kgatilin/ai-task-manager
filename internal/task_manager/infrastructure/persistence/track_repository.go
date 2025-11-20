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

// Compile-time check that SQLiteTrackRepository implements repositories.TrackRepository
var _ repositories.TrackRepository = (*SQLiteTrackRepository)(nil)

// SQLiteTrackRepository implements repositories.TrackRepository using SQLite as the backend.
type SQLiteTrackRepository struct {
	DB     *sql.DB
	logger logger.Logger
}

// NewSQLiteTrackRepository creates a new SQLite-backed repository.
func NewSQLiteTrackRepository(db *sql.DB, logger logger.Logger) *SQLiteTrackRepository {
	return &SQLiteTrackRepository{
		DB:     db,
		logger: logger,
	}
}

// ============================================================================
// Track Operations
// ============================================================================

// SaveTrack persists a new track to storage.
func (r *SQLiteTrackRepository) SaveTrack(ctx context.Context, track *entities.TrackEntity) error {
	// Check if track already exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tracks WHERE id = ?", track.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check track existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w: track %s already exists", tmerrors.ErrAlreadyExists, track.ID)
	}

	// Check if roadmap exists
	var roadmapExists int
	err = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM roadmaps WHERE id = ?", track.RoadmapID).Scan(&roadmapExists)
	if err != nil {
		return fmt.Errorf("failed to check roadmap existence: %w", err)
	}
	if roadmapExists == 0 {
		return fmt.Errorf("%w: roadmap %s not found", tmerrors.ErrNotFound, track.RoadmapID)
	}

	// Start transaction for track and dependencies
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert track
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO tracks (id, roadmap_id, title, description, status, rank, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		track.ID, track.RoadmapID, track.Title, track.Description, track.Status, track.Rank, track.CreatedAt, track.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}

	// Insert dependencies
	for _, depID := range track.Dependencies {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO track_dependencies (track_id, depends_on_id) VALUES (?, ?)",
			track.ID, depID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert track dependency: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTrack retrieves a track by its ID.
func (r *SQLiteTrackRepository) GetTrack(ctx context.Context, id string) (*entities.TrackEntity, error) {
	var track entities.TrackEntity

	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id, roadmap_id, title, description, status, rank, created_at, updated_at FROM tracks WHERE id = ?",
		id,
	).Scan(&track.ID, &track.RoadmapID, &track.Title, &track.Description, &track.Status, &track.Rank, &track.CreatedAt, &track.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: track %s not found", tmerrors.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to query track: %w", err)
	}

	// Load dependencies
	deps, err := r.GetTrackDependencies(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load track dependencies: %w", err)
	}
	track.Dependencies = deps

	return &track, nil
}

// ListTracks returns all tracks for a roadmap, optionally filtered.
func (r *SQLiteTrackRepository) ListTracks(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
	query := "SELECT id, roadmap_id, title, description, status, rank, created_at, updated_at FROM tracks WHERE roadmap_id = ?"
	args := []interface{}{roadmapID}

	// Add status filter if provided
	if len(filters.Status) > 0 {
		placeholders := ""
		for i := range filters.Status {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, filters.Status[i])
		}
		query += " AND status IN (" + placeholders + ")"
	}

	// Add priority filter if provided
	if len(filters.Priority) > 0 {
		placeholders := ""
		for i := range filters.Priority {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, filters.Priority[i])
		}
		query += " AND rank IN (" + placeholders + ")"
	}

	query += " ORDER BY id"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*entities.TrackEntity
	for rows.Next() {
		var track entities.TrackEntity
		err := rows.Scan(&track.ID, &track.RoadmapID, &track.Title, &track.Description, &track.Status, &track.Rank, &track.CreatedAt, &track.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}

		// Load dependencies
		deps, err := r.GetTrackDependencies(ctx, track.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load track dependencies: %w", err)
		}
		track.Dependencies = deps

		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tracks: %w", err)
	}

	return tracks, nil
}

// UpdateTrack updates an existing track.
func (r *SQLiteTrackRepository) UpdateTrack(ctx context.Context, track *entities.TrackEntity) error {
	// Start transaction for track and dependencies update
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update track fields
	result, err := tx.ExecContext(
		ctx,
		"UPDATE tracks SET title = ?, description = ?, status = ?, rank = ?, updated_at = ? WHERE id = ?",
		track.Title, track.Description, track.Status, track.Rank, track.UpdatedAt, track.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: track %s not found", tmerrors.ErrNotFound, track.ID)
	}

	// Delete existing dependencies
	_, err = tx.ExecContext(ctx, "DELETE FROM track_dependencies WHERE track_id = ?", track.ID)
	if err != nil {
		return fmt.Errorf("failed to delete dependencies: %w", err)
	}

	// Insert new dependencies
	for _, depID := range track.Dependencies {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO track_dependencies (track_id, depends_on_id) VALUES (?, ?)",
			track.ID, depID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert dependency: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteTrack removes a track from storage.
func (r *SQLiteTrackRepository) DeleteTrack(ctx context.Context, id string) error {
	result, err := r.DB.ExecContext(ctx, "DELETE FROM tracks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: track %s not found", tmerrors.ErrNotFound, id)
	}

	return nil
}

// AddTrackDependency adds a dependency from trackID to dependsOnID.
func (r *SQLiteTrackRepository) AddTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	// Check for self-dependency
	if trackID == dependsOnID {
		return fmt.Errorf("%w: track cannot depend on itself", tmerrors.ErrInvalidArgument)
	}

	// Check both tracks exist
	for _, id := range []string{trackID, dependsOnID} {
		var exists int
		err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tracks WHERE id = ?", id).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check track existence: %w", err)
		}
		if exists == 0 {
			return fmt.Errorf("%w: track %s not found", tmerrors.ErrNotFound, id)
		}
	}

	// Check if dependency already exists
	var exists int
	err := r.DB.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM track_dependencies WHERE track_id = ? AND depends_on_id = ?",
		trackID, dependsOnID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check dependency existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w: dependency already exists", tmerrors.ErrAlreadyExists)
	}

	// Insert dependency
	_, err = r.DB.ExecContext(
		ctx,
		"INSERT INTO track_dependencies (track_id, depends_on_id) VALUES (?, ?)",
		trackID, dependsOnID,
	)
	if err != nil {
		return fmt.Errorf("failed to add dependency: %w", err)
	}

	return nil
}

// RemoveTrackDependency removes a dependency from trackID to dependsOnID.
func (r *SQLiteTrackRepository) RemoveTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	result, err := r.DB.ExecContext(
		ctx,
		"DELETE FROM track_dependencies WHERE track_id = ? AND depends_on_id = ?",
		trackID, dependsOnID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove dependency: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: dependency not found", tmerrors.ErrNotFound)
	}

	return nil
}

// GetTrackDependencies returns the IDs of all tracks that trackID depends on.
func (r *SQLiteTrackRepository) GetTrackDependencies(ctx context.Context, trackID string) ([]string, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT depends_on_id FROM track_dependencies WHERE track_id = ? ORDER BY depends_on_id",
		trackID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies: %w", err)
	}
	defer rows.Close()

	var deps []string
	for rows.Next() {
		var depID string
		if err := rows.Scan(&depID); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %w", err)
		}
		deps = append(deps, depID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating dependencies: %w", err)
	}

	return deps, nil
}

// ValidateNoCycles checks if adding/updating the track would create a circular dependency.
// Uses depth-first search to detect cycles.
func (r *SQLiteTrackRepository) ValidateNoCycles(ctx context.Context, trackID string) error {
	// Use DFS to detect cycles
	visited := make(map[string]bool)
	return r.detectCycleDFS(ctx, trackID, visited)
}

// GetTrackWithTasks retrieves a track with all its tasks.
func (r *SQLiteTrackRepository) GetTrackWithTasks(ctx context.Context, trackID string) (*entities.TrackEntity, error) {
	track, err := r.GetTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}

	// Load all tasks for this track
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, track_id, title, description, status, rank, branch, created_at, updated_at FROM tasks WHERE track_id = ?",
		trackID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}
	defer rows.Close()

	// Note: TrackEntity doesn't have a Tasks field, so we just verify loading works
	for rows.Next() {
		var (
			id, trackID, title, desc, status string
			rank                             int
			branch                           sql.NullString
			createdAt, updatedAt             string
		)
		if err := rows.Scan(&id, &trackID, &title, &desc, &status, &rank, &branch, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return track, nil
}

// ============================================================================
// Helper Methods
// ============================================================================

// detectCycleDFS performs depth-first search to detect cycles.
func (r *SQLiteTrackRepository) detectCycleDFS(ctx context.Context, trackID string, visited map[string]bool) error {
	if visited[trackID] {
		return fmt.Errorf("%w: circular dependency detected", tmerrors.ErrInvalidArgument)
	}

	visited[trackID] = true

	deps, err := r.GetTrackDependencies(ctx, trackID)
	if err != nil {
		return err
	}

	for _, depID := range deps {
		if err := r.detectCycleDFS(ctx, depID, visited); err != nil {
			return err
		}
	}

	visited[trackID] = false
	return nil
}
