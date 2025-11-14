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

// Compile-time check that SQLiteAcceptanceCriteriaRepository implements repositories.AcceptanceCriteriaRepository
var _ repositories.AcceptanceCriteriaRepository = (*SQLiteAcceptanceCriteriaRepository)(nil)

// SQLiteAcceptanceCriteriaRepository implements repositories.AcceptanceCriteriaRepository using SQLite as the backend.
type SQLiteAcceptanceCriteriaRepository struct {
	DB     *sql.DB
	logger pluginsdk.Logger
}

// NewSQLiteAcceptanceCriteriaRepository creates a new SQLite-backed repository.
func NewSQLiteAcceptanceCriteriaRepository(db *sql.DB, logger pluginsdk.Logger) *SQLiteAcceptanceCriteriaRepository {
	return &SQLiteAcceptanceCriteriaRepository{
		DB:     db,
		logger: logger,
	}
}

// ============================================================================
// Acceptance Criteria Operations
// ============================================================================

// SaveAC persists a new acceptance criterion to storage.
func (r *SQLiteAcceptanceCriteriaRepository) SaveAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error {
	// Check if AC already exists
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM acceptance_criteria WHERE id = ?", ac.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check AC existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w: AC %s already exists", pluginsdk.ErrAlreadyExists, ac.ID)
	}

	// Verify task exists
	var taskExists int
	err = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tasks WHERE id = ?", ac.TaskID).Scan(&taskExists)
	if err != nil {
		return fmt.Errorf("failed to verify task: %w", err)
	}
	if taskExists == 0 {
		return fmt.Errorf("%w: task %s not found", pluginsdk.ErrNotFound, ac.TaskID)
	}

	_, err = r.DB.ExecContext(
		ctx,
		"INSERT INTO acceptance_criteria (id, task_id, description, verification_type, status, notes, testing_instructions, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		ac.ID, ac.TaskID, ac.Description, string(ac.VerificationType), string(ac.Status), ac.Notes, ac.TestingInstructions, ac.CreatedAt, ac.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert AC: %w", err)
	}

	return nil
}

// GetAC retrieves an acceptance criterion by its ID.
func (r *SQLiteAcceptanceCriteriaRepository) GetAC(ctx context.Context, id string) (*entities.AcceptanceCriteriaEntity, error) {
	var ac entities.AcceptanceCriteriaEntity

	var testingInstructions sql.NullString
	err := r.DB.QueryRowContext(
		ctx,
		"SELECT id, task_id, description, verification_type, status, notes, testing_instructions, created_at, updated_at FROM acceptance_criteria WHERE id = ?",
		id,
	).Scan(&ac.ID, &ac.TaskID, &ac.Description, (*string)(&ac.VerificationType), (*string)(&ac.Status), &ac.Notes, &testingInstructions, &ac.CreatedAt, &ac.UpdatedAt)

	if testingInstructions.Valid {
		ac.TestingInstructions = testingInstructions.String
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: AC %s not found", pluginsdk.ErrNotFound, id)
		}
		return nil, fmt.Errorf("failed to query AC: %w", err)
	}

	return &ac, nil
}

// ListAC returns all acceptance criteria for a task.
func (r *SQLiteAcceptanceCriteriaRepository) ListAC(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT id, task_id, description, verification_type, status, notes, testing_instructions, created_at, updated_at FROM acceptance_criteria WHERE task_id = ? ORDER BY created_at ASC",
		taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query ACs: %w", err)
	}
	defer rows.Close()

	var acs []*entities.AcceptanceCriteriaEntity
	for rows.Next() {
		var ac entities.AcceptanceCriteriaEntity
		var testingInstructions sql.NullString
		err := rows.Scan(&ac.ID, &ac.TaskID, &ac.Description, (*string)(&ac.VerificationType), (*string)(&ac.Status), &ac.Notes, &testingInstructions, &ac.CreatedAt, &ac.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan AC: %w", err)
		}
		if testingInstructions.Valid {
			ac.TestingInstructions = testingInstructions.String
		}
		acs = append(acs, &ac)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ACs: %w", err)
	}

	return acs, nil
}

// UpdateAC updates an existing acceptance criterion.
func (r *SQLiteAcceptanceCriteriaRepository) UpdateAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error {
	result, err := r.DB.ExecContext(
		ctx,
		"UPDATE acceptance_criteria SET task_id = ?, description = ?, verification_type = ?, status = ?, notes = ?, testing_instructions = ?, updated_at = ? WHERE id = ?",
		ac.TaskID, ac.Description, string(ac.VerificationType), string(ac.Status), ac.Notes, ac.TestingInstructions, ac.UpdatedAt, ac.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update AC: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: AC %s not found", pluginsdk.ErrNotFound, ac.ID)
	}

	return nil
}

// DeleteAC removes an acceptance criterion from storage.
func (r *SQLiteAcceptanceCriteriaRepository) DeleteAC(ctx context.Context, id string) error {
	result, err := r.DB.ExecContext(ctx, "DELETE FROM acceptance_criteria WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete AC: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: AC %s not found", pluginsdk.ErrNotFound, id)
	}

	return nil
}

// ListACByTask is an alias for ListAC for consistency with other repositories.
func (r *SQLiteAcceptanceCriteriaRepository) ListACByTask(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
	return r.ListAC(ctx, taskID)
}

// ListACByIteration returns all acceptance criteria for all tasks in an iteration.
func (r *SQLiteAcceptanceCriteriaRepository) ListACByIteration(ctx context.Context, iterationNum int) ([]*entities.AcceptanceCriteriaEntity, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT ac.id, ac.task_id, ac.description, ac.verification_type, ac.status, ac.notes, ac.testing_instructions, ac.created_at, ac.updated_at
		 FROM acceptance_criteria ac
		 JOIN tasks t ON ac.task_id = t.id
		 JOIN iteration_tasks it ON t.id = it.task_id
		 WHERE it.iteration_number = ?
		 ORDER BY ac.created_at ASC`,
		iterationNum,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query ACs by iteration: %w", err)
	}
	defer rows.Close()

	var acs []*entities.AcceptanceCriteriaEntity
	for rows.Next() {
		var ac entities.AcceptanceCriteriaEntity
		var testingInstructions sql.NullString
		err := rows.Scan(&ac.ID, &ac.TaskID, &ac.Description, (*string)(&ac.VerificationType), (*string)(&ac.Status), &ac.Notes, &testingInstructions, &ac.CreatedAt, &ac.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan AC: %w", err)
		}
		if testingInstructions.Valid {
			ac.TestingInstructions = testingInstructions.String
		}
		acs = append(acs, &ac)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ACs: %w", err)
	}

	return acs, nil
}

// ListFailedAC returns all acceptance criteria with status "failed".
func (r *SQLiteAcceptanceCriteriaRepository) ListFailedAC(ctx context.Context, filters entities.ACFilters) ([]*entities.AcceptanceCriteriaEntity, error) {
	query := `SELECT ac.id, ac.task_id, ac.description, ac.verification_type, ac.status, ac.notes, ac.testing_instructions, ac.created_at, ac.updated_at
		      FROM acceptance_criteria ac`

	var joins []string
	var conditions []string
	var args []interface{}

	// Base condition: status = failed
	conditions = append(conditions, "ac.status = ?")
	args = append(args, string(entities.ACStatusFailed))

	// Add iteration filter
	if filters.IterationNum != nil {
		joins = append(joins, "JOIN iteration_tasks it ON ac.task_id = it.task_id")
		conditions = append(conditions, "it.iteration_number = ?")
		args = append(args, *filters.IterationNum)
	}

	// Add track filter
	if filters.TrackID != "" {
		joins = append(joins, "JOIN tasks t ON ac.task_id = t.id")
		conditions = append(conditions, "t.track_id = ?")
		args = append(args, filters.TrackID)
	}

	// Add task filter
	if filters.TaskID != "" {
		conditions = append(conditions, "ac.task_id = ?")
		args = append(args, filters.TaskID)
	}

	// Build final query
	for _, join := range joins {
		query += " " + join
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY ac.created_at ASC"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query failed ACs: %w", err)
	}
	defer rows.Close()

	var acs []*entities.AcceptanceCriteriaEntity
	for rows.Next() {
		var ac entities.AcceptanceCriteriaEntity
		var testingInstructions sql.NullString
		err := rows.Scan(&ac.ID, &ac.TaskID, &ac.Description, (*string)(&ac.VerificationType), (*string)(&ac.Status), &ac.Notes, &testingInstructions, &ac.CreatedAt, &ac.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan AC: %w", err)
		}
		if testingInstructions.Valid {
			ac.TestingInstructions = testingInstructions.String
		}
		acs = append(acs, &ac)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ACs: %w", err)
	}

	return acs, nil
}
