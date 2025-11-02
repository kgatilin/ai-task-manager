package task_manager

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ============================================================================
// MigrateIDsCommand migrates old-format IDs to new human-readable format
// ============================================================================

type MigrateIDsCommand struct {
	Plugin  *TaskManagerPlugin
	project string
	dryRun  bool
}

func (c *MigrateIDsCommand) GetName() string {
	return "migrate-ids"
}

func (c *MigrateIDsCommand) GetDescription() string {
	return "Migrate old-format IDs to new human-readable format"
}

func (c *MigrateIDsCommand) GetUsage() string {
	return "dw task-manager migrate-ids [--project <name>] [--dry-run]"
}

func (c *MigrateIDsCommand) GetHelp() string {
	return `Migrates existing track and task IDs from old timestamp-based format
to new human-readable format (e.g., track-1730419200000 → DW-track-1).

This command is safe to run multiple times - it will skip IDs that are
already in the new format.

Flags:
  --project <name>  Project to migrate (default: active project)
  --dry-run         Show what would be migrated without making changes

Examples:
  # Migrate active project
  dw task-manager migrate-ids

  # Preview migration for a specific project
  dw task-manager migrate-ids --project production --dry-run

Notes:
  - Migration runs in a transaction - either all changes succeed or none do
  - Old IDs are preserved temporarily during migration for safety
  - All foreign key relationships are updated automatically
  - Iteration IDs don't need migration (already integers)`
}

func (c *MigrateIDsCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--dry-run":
			c.dryRun = true
		}
	}

	// Get repository for project
	repo, cleanup, err := c.Plugin.getRepositoryForProject(c.project)
	if err != nil {
		return err
	}
	defer cleanup()

	// Get project code
	projectCode := repo.GetProjectCode(ctx)
	if projectCode == "" {
		return fmt.Errorf("project code not set - cannot generate new IDs")
	}

	// Get underlying database for migration
	sqliteRepo, ok := repo.(*SQLiteRoadmapRepository)
	if !ok {
		// If wrapped in EventEmittingRepository, unwrap it
		if eventRepo, ok := repo.(*EventEmittingRepository); ok {
			if sqliteBaseRepo, ok := eventRepo.repo.(*SQLiteRoadmapRepository); ok {
				sqliteRepo = sqliteBaseRepo
			} else {
				return fmt.Errorf("unsupported repository type for migration")
			}
		} else {
			return fmt.Errorf("unsupported repository type for migration")
		}
	}

	fmt.Fprintf(cmdCtx.GetStdout(), "Analyzing %s project for migration...\n\n", c.project)

	// Check what needs migrating
	trackCount, taskCount, err := c.analyzeData(ctx, sqliteRepo)
	if err != nil {
		return fmt.Errorf("failed to analyze data: %w", err)
	}

	if trackCount == 0 && taskCount == 0 {
		fmt.Fprintf(cmdCtx.GetStdout(), "✓ No IDs need migration - all IDs are already in new format\n")
		return nil
	}

	fmt.Fprintf(cmdCtx.GetStdout(), "Found:\n")
	fmt.Fprintf(cmdCtx.GetStdout(), "  - %d track(s) to migrate\n", trackCount)
	fmt.Fprintf(cmdCtx.GetStdout(), "  - %d task(s) to migrate\n", taskCount)
	fmt.Fprintf(cmdCtx.GetStdout(), "\n")

	if c.dryRun {
		fmt.Fprintf(cmdCtx.GetStdout(), "Dry run - no changes made\n")
		return nil
	}

	// Perform migration
	if err := c.performMigration(ctx, sqliteRepo, projectCode); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Fprintf(cmdCtx.GetStdout(), "✓ Migration completed successfully!\n")
	fmt.Fprintf(cmdCtx.GetStdout(), "  - %d track IDs updated\n", trackCount)
	fmt.Fprintf(cmdCtx.GetStdout(), "  - %d task IDs updated\n", taskCount)

	return nil
}

// analyzeData checks how many entities need migration
func (c *MigrateIDsCommand) analyzeData(ctx context.Context, repo *SQLiteRoadmapRepository) (trackCount, taskCount int, err error) {
	// Pattern for old timestamp-based IDs
	oldPattern := regexp.MustCompile(`^(track|task)-\d{13,}$`)

	// Count tracks with old format
	rows, err := repo.db.QueryContext(ctx, "SELECT id FROM tracks")
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return 0, 0, err
		}
		if oldPattern.MatchString(id) {
			trackCount++
		}
	}

	// Count tasks with old format
	rows, err = repo.db.QueryContext(ctx, "SELECT id FROM tasks")
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return 0, 0, err
		}
		if oldPattern.MatchString(id) {
			taskCount++
		}
	}

	return trackCount, taskCount, nil
}

// performMigration executes the migration in a transaction
func (c *MigrateIDsCommand) performMigration(ctx context.Context, repo *SQLiteRoadmapRepository, projectCode string) error {
	// Start transaction
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Migrate tracks
	if err := c.migrateTracks(ctx, tx, projectCode); err != nil {
		return fmt.Errorf("failed to migrate tracks: %w", err)
	}

	// Migrate tasks
	if err := c.migrateTasks(ctx, tx, projectCode); err != nil {
		return fmt.Errorf("failed to migrate tasks: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// migrateTracks migrates track IDs
func (c *MigrateIDsCommand) migrateTracks(ctx context.Context, tx *sql.Tx, projectCode string) error {
	// Pattern for old timestamp-based IDs
	oldPattern := regexp.MustCompile(`^track-\d{13,}$`)

	// Get all tracks with old IDs, ordered by created_at to preserve creation order
	rows, err := tx.QueryContext(ctx, "SELECT id, created_at FROM tracks ORDER BY created_at ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	type trackToMigrate struct {
		oldID     string
		createdAt time.Time
	}
	var tracksToMigrate []trackToMigrate

	for rows.Next() {
		var id string
		var createdAt time.Time
		if err := rows.Scan(&id, &createdAt); err != nil {
			return err
		}
		if oldPattern.MatchString(id) {
			tracksToMigrate = append(tracksToMigrate, trackToMigrate{
				oldID:     id,
				createdAt: createdAt,
			})
		}
	}

	// Get current sequence number for tracks
	var nextNum int
	err = tx.QueryRowContext(ctx,
		"SELECT COALESCE(MAX(CAST(SUBSTR(id, LENGTH(?) + 8) AS INTEGER)), 0) + 1 FROM tracks WHERE id LIKE ? || '-track-%'",
		projectCode, projectCode).Scan(&nextNum)
	if err != nil {
		nextNum = 1 // Start from 1 if no existing tracks
	}

	// Migrate each track
	for _, track := range tracksToMigrate {
		newID := fmt.Sprintf("%s-track-%d", projectCode, nextNum)
		nextNum++

		// Update track ID
		_, err := tx.ExecContext(ctx,
			"UPDATE tracks SET id = ? WHERE id = ?",
			newID, track.oldID)
		if err != nil {
			return fmt.Errorf("failed to update track %s: %w", track.oldID, err)
		}

		// Update track dependencies
		_, err = tx.ExecContext(ctx,
			"UPDATE track_dependencies SET track_id = ? WHERE track_id = ?",
			newID, track.oldID)
		if err != nil {
			return fmt.Errorf("failed to update track dependencies for %s: %w", track.oldID, err)
		}

		_, err = tx.ExecContext(ctx,
			"UPDATE track_dependencies SET depends_on_id = ? WHERE depends_on_id = ?",
			newID, track.oldID)
		if err != nil {
			return fmt.Errorf("failed to update track dependency references for %s: %w", track.oldID, err)
		}

		// Update tasks that reference this track
		_, err = tx.ExecContext(ctx,
			"UPDATE tasks SET track_id = ? WHERE track_id = ?",
			newID, track.oldID)
		if err != nil {
			return fmt.Errorf("failed to update tasks for track %s: %w", track.oldID, err)
		}
	}

	// Update sequence number
	if len(tracksToMigrate) > 0 {
		_, err = tx.ExecContext(ctx,
			"INSERT OR REPLACE INTO sequences (name, value) VALUES (?, ?)",
			"track", nextNum)
		if err != nil {
			return fmt.Errorf("failed to update track sequence: %w", err)
		}
	}

	return nil
}

// migrateTasks migrates task IDs
func (c *MigrateIDsCommand) migrateTasks(ctx context.Context, tx *sql.Tx, projectCode string) error {
	// Pattern for old timestamp-based IDs
	oldPattern := regexp.MustCompile(`^task-\d{13,}$`)

	// Get all tasks with old IDs, ordered by created_at to preserve creation order
	rows, err := tx.QueryContext(ctx, "SELECT id, created_at FROM tasks ORDER BY created_at ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	type taskToMigrate struct {
		oldID     string
		createdAt time.Time
	}
	var tasksToMigrate []taskToMigrate

	for rows.Next() {
		var id string
		var createdAt time.Time
		if err := rows.Scan(&id, &createdAt); err != nil {
			return err
		}
		if oldPattern.MatchString(id) {
			tasksToMigrate = append(tasksToMigrate, taskToMigrate{
				oldID:     id,
				createdAt: createdAt,
			})
		}
	}

	// Get current sequence number for tasks
	var nextNum int
	err = tx.QueryRowContext(ctx,
		"SELECT COALESCE(MAX(CAST(SUBSTR(id, LENGTH(?) + 7) AS INTEGER)), 0) + 1 FROM tasks WHERE id LIKE ? || '-task-%'",
		projectCode, projectCode).Scan(&nextNum)
	if err != nil {
		nextNum = 1 // Start from 1 if no existing tasks
	}

	// Migrate each task
	for _, task := range tasksToMigrate {
		newID := fmt.Sprintf("%s-task-%d", projectCode, nextNum)
		nextNum++

		// Update task ID
		_, err := tx.ExecContext(ctx,
			"UPDATE tasks SET id = ? WHERE id = ?",
			newID, task.oldID)
		if err != nil {
			return fmt.Errorf("failed to update task %s: %w", task.oldID, err)
		}

		// Update iteration_tasks references
		_, err = tx.ExecContext(ctx,
			"UPDATE iteration_tasks SET task_id = ? WHERE task_id = ?",
			newID, task.oldID)
		if err != nil {
			return fmt.Errorf("failed to update iteration_tasks for task %s: %w", task.oldID, err)
		}
	}

	// Update sequence number
	if len(tasksToMigrate) > 0 {
		_, err = tx.ExecContext(ctx,
			"INSERT OR REPLACE INTO sequences (name, value) VALUES (?, ?)",
			"task", nextNum)
		if err != nil {
			return fmt.Errorf("failed to update task sequence: %w", err)
		}
	}

	return nil
}
