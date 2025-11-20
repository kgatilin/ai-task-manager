package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	tmerrors "github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/logger"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/repositories"
)

// Compile-time check that SQLiteADRRepository implements repositories.ADRRepository
var _ repositories.ADRRepository = (*SQLiteADRRepository)(nil)

// SQLiteADRRepository implements repositories.ADRRepository using SQLite as the backend.
type SQLiteADRRepository struct {
	DB     *sql.DB
	logger logger.Logger
}

// NewSQLiteADRRepository creates a new SQLite-backed repository.
func NewSQLiteADRRepository(db *sql.DB, logger logger.Logger) *SQLiteADRRepository {
	return &SQLiteADRRepository{
		DB:     db,
		logger: logger,
	}
}

// ============================================================================
// ADR Operations
// ============================================================================

// SaveADR persists a new ADR to storage.
func (r *SQLiteADRRepository) SaveADR(ctx context.Context, adr *entities.ADREntity) error {
	// Check if ADR already exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM adrs WHERE id = ?", adr.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check ADR existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w: ADR %s already exists", tmerrors.ErrAlreadyExists, adr.ID)
	}

	// Check if track exists
	err = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tracks WHERE id = ?", adr.TrackID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check track existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("%w: track %s does not exist", tmerrors.ErrNotFound, adr.TrackID)
	}

	_, err = r.DB.ExecContext(
		ctx,
		"INSERT INTO adrs (id, track_id, title, status, context, decision, consequences, alternatives, created_at, updated_at, superseded_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		adr.ID, adr.TrackID, adr.Title, adr.Status, adr.Context, adr.Decision, adr.Consequences, adr.Alternatives, adr.CreatedAt, adr.UpdatedAt, adr.SupersededBy,
	)
	if err != nil {
		return fmt.Errorf("failed to insert ADR: %w", err)
	}

	return nil
}

// GetADR retrieves an ADR by its ID.
func (r *SQLiteADRRepository) GetADR(ctx context.Context, id string) (*entities.ADREntity, error) {
	row := r.DB.QueryRowContext(
		ctx,
		"SELECT id, track_id, title, status, context, decision, consequences, alternatives, created_at, updated_at, superseded_by FROM adrs WHERE id = ?",
		id,
	)

	var adr entities.ADREntity
	var supersededBy sql.NullString
	err := row.Scan(
		&adr.ID, &adr.TrackID, &adr.Title, &adr.Status, &adr.Context, &adr.Decision, &adr.Consequences, &adr.Alternatives, &adr.CreatedAt, &adr.UpdatedAt, &supersededBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: ADR %s not found", tmerrors.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to query ADR: %w", err)
	}

	if supersededBy.Valid {
		adr.SupersededBy = &supersededBy.String
	}

	return &adr, nil
}

// ListADRs returns all ADRs, optionally filtered by track.
func (r *SQLiteADRRepository) ListADRs(ctx context.Context, trackID *string) ([]*entities.ADREntity, error) {
	query := "SELECT id, track_id, title, status, context, decision, consequences, alternatives, created_at, updated_at, superseded_by FROM adrs"
	var args []interface{}

	if trackID != nil {
		query += " WHERE track_id = ?"
		args = append(args, *trackID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query ADRs: %w", err)
	}
	defer rows.Close()

	var adrs []*entities.ADREntity
	for rows.Next() {
		var adr entities.ADREntity
		var supersededBy sql.NullString
		err := rows.Scan(
			&adr.ID, &adr.TrackID, &adr.Title, &adr.Status, &adr.Context, &adr.Decision, &adr.Consequences, &adr.Alternatives, &adr.CreatedAt, &adr.UpdatedAt, &supersededBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ADR: %w", err)
		}

		if supersededBy.Valid {
			adr.SupersededBy = &supersededBy.String
		}

		adrs = append(adrs, &adr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ADRs: %w", err)
	}

	return adrs, nil
}

// UpdateADR updates an existing ADR.
func (r *SQLiteADRRepository) UpdateADR(ctx context.Context, adr *entities.ADREntity) error {
	// Check if ADR exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM adrs WHERE id = ?", adr.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check ADR existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("%w: ADR %s not found", tmerrors.ErrNotFound, adr.ID)
	}

	_, err = r.DB.ExecContext(
		ctx,
		"UPDATE adrs SET title = ?, status = ?, context = ?, decision = ?, consequences = ?, alternatives = ?, updated_at = ?, superseded_by = ? WHERE id = ?",
		adr.Title, adr.Status, adr.Context, adr.Decision, adr.Consequences, adr.Alternatives, adr.UpdatedAt, adr.SupersededBy, adr.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update ADR: %w", err)
	}

	return nil
}

// SupersedeADR marks an ADR as superseded by another ADR.
func (r *SQLiteADRRepository) SupersedeADR(ctx context.Context, adrID, supersededByID string) error {
	// Check if both ADRs exist
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM adrs WHERE id = ?", adrID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check ADR existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("%w: ADR %s not found", tmerrors.ErrNotFound, adrID)
	}

	err = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM adrs WHERE id = ?", supersededByID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check superseding ADR existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("%w: ADR %s not found", tmerrors.ErrNotFound, supersededByID)
	}

	now := time.Now().UTC()
	_, err = r.DB.ExecContext(
		ctx,
		"UPDATE adrs SET status = ?, superseded_by = ?, updated_at = ? WHERE id = ?",
		string(entities.ADRStatusSuperseded), supersededByID, now, adrID,
	)
	if err != nil {
		return fmt.Errorf("failed to supersede ADR: %w", err)
	}

	return nil
}

// DeprecateADR marks an ADR as deprecated.
func (r *SQLiteADRRepository) DeprecateADR(ctx context.Context, adrID string) error {
	// Check if ADR exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM adrs WHERE id = ?", adrID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check ADR existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("%w: ADR %s not found", tmerrors.ErrNotFound, adrID)
	}

	now := time.Now().UTC()
	_, err = r.DB.ExecContext(
		ctx,
		"UPDATE adrs SET status = ?, updated_at = ? WHERE id = ?",
		string(entities.ADRStatusDeprecated), now, adrID,
	)
	if err != nil {
		return fmt.Errorf("failed to deprecate ADR: %w", err)
	}

	return nil
}

// GetADRsByTrack returns all ADRs for a specific track.
func (r *SQLiteADRRepository) GetADRsByTrack(ctx context.Context, trackID string) ([]*entities.ADREntity, error) {
	return r.ListADRs(ctx, &trackID)
}
