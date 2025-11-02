package task_manager

import (
	"database/sql"
	"fmt"
)

// migrateV3ToV4 migrates database from schema version 3 (priority TEXT) to version 4 (rank INTEGER)
func migrateV3ToV4(db *sql.DB) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if tracks table has priority column (version 3)
	var hasPriority bool
	rows, err := tx.Query("PRAGMA table_info(tracks)")
	if err != nil {
		return fmt.Errorf("failed to check tracks table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return fmt.Errorf("failed to scan column info: %w", err)
		}
		if name == "priority" {
			hasPriority = true
			break
		}
	}
	rows.Close()

	if !hasPriority {
		// Already migrated or new database
		return tx.Commit()
	}

	fmt.Println("Migrating database from schema v3 to v4 (priority -> rank)...")

	// MIGRATE TRACKS TABLE
	// 1. Create new tracks table with rank
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS tracks_new (
			id TEXT PRIMARY KEY,
			roadmap_id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL,
			rank INTEGER NOT NULL DEFAULT 500,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY(roadmap_id) REFERENCES roadmaps(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create new tracks table: %w", err)
	}

	// 2. Copy data with priority -> rank conversion
	_, err = tx.Exec(`
		INSERT INTO tracks_new (id, roadmap_id, title, description, status, rank, created_at, updated_at)
		SELECT
			id,
			roadmap_id,
			title,
			description,
			status,
			CASE priority
				WHEN 'critical' THEN 100
				WHEN 'high' THEN 200
				WHEN 'medium' THEN 300
				WHEN 'low' THEN 400
				ELSE 500
			END as rank,
			created_at,
			updated_at
		FROM tracks
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate tracks data: %w", err)
	}

	// 3. Drop old table and rename new one
	if _, err = tx.Exec("DROP TABLE tracks"); err != nil {
		return fmt.Errorf("failed to drop old tracks table: %w", err)
	}
	if _, err = tx.Exec("ALTER TABLE tracks_new RENAME TO tracks"); err != nil {
		return fmt.Errorf("failed to rename new tracks table: %w", err)
	}

	// MIGRATE TASKS TABLE
	// Check if tasks table has priority column
	var hasTaskPriority bool
	rows, err = tx.Query("PRAGMA table_info(tasks)")
	if err != nil {
		return fmt.Errorf("failed to check tasks table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return fmt.Errorf("failed to scan task column info: %w", err)
		}
		if name == "priority" {
			hasTaskPriority = true
			break
		}
	}
	rows.Close()

	if hasTaskPriority {
		// 1. Create new tasks table with rank
		_, err = tx.Exec(`
			CREATE TABLE IF NOT EXISTS tasks_new (
				id TEXT PRIMARY KEY,
				track_id TEXT NOT NULL,
				title TEXT NOT NULL,
				description TEXT,
				status TEXT NOT NULL,
				rank INTEGER NOT NULL DEFAULT 500,
				branch TEXT,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL,
				FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create new tasks table: %w", err)
		}

		// 2. Copy data with priority -> rank conversion
		_, err = tx.Exec(`
			INSERT INTO tasks_new (id, track_id, title, description, status, rank, branch, created_at, updated_at)
			SELECT
				id,
				track_id,
				title,
				description,
				status,
				CASE priority
					WHEN 'critical' THEN 100
					WHEN 'high' THEN 200
					WHEN 'medium' THEN 300
					WHEN 'low' THEN 400
					ELSE 500
				END as rank,
				branch,
				created_at,
				updated_at
			FROM tasks
		`)
		if err != nil {
			return fmt.Errorf("failed to migrate tasks data: %w", err)
		}

		// 3. Drop old table and rename new one
		if _, err = tx.Exec("DROP TABLE tasks"); err != nil {
			return fmt.Errorf("failed to drop old tasks table: %w", err)
		}
		if _, err = tx.Exec("ALTER TABLE tasks_new RENAME TO tasks"); err != nil {
			return fmt.Errorf("failed to rename new tasks table: %w", err)
		}
	}

	// MIGRATE ITERATIONS TABLE
	// Check if iterations table has rank column
	var hasIterRank bool
	rows, err = tx.Query("PRAGMA table_info(iterations)")
	if err != nil {
		return fmt.Errorf("failed to check iterations table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return fmt.Errorf("failed to scan iteration column info: %w", err)
		}
		if name == "rank" {
			hasIterRank = true
			break
		}
	}
	rows.Close()

	if !hasIterRank {
		// 1. Create new iterations table with rank
		_, err = tx.Exec(`
			CREATE TABLE IF NOT EXISTS iterations_new (
				number INTEGER PRIMARY KEY,
				name TEXT NOT NULL,
				goal TEXT,
				status TEXT NOT NULL,
				rank INTEGER NOT NULL DEFAULT 500,
				deliverable TEXT,
				started_at TIMESTAMP,
				completed_at TIMESTAMP,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create new iterations table: %w", err)
		}

		// 2. Copy data with rank based on iteration number (number * 100)
		_, err = tx.Exec(`
			INSERT INTO iterations_new (number, name, goal, status, rank, deliverable, started_at, completed_at, created_at, updated_at)
			SELECT
				number,
				name,
				goal,
				status,
				number * 100 as rank,
				deliverable,
				started_at,
				completed_at,
				created_at,
				updated_at
			FROM iterations
		`)
		if err != nil {
			return fmt.Errorf("failed to migrate iterations data: %w", err)
		}

		// 3. Drop old table and rename new one
		if _, err = tx.Exec("DROP TABLE iterations"); err != nil {
			return fmt.Errorf("failed to drop old iterations table: %w", err)
		}
		if _, err = tx.Exec("ALTER TABLE iterations_new RENAME TO iterations"); err != nil {
			return fmt.Errorf("failed to rename new iterations table: %w", err)
		}
	}

	// Recreate indexes
	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_tracks_roadmap_id ON tracks(roadmap_id)")
	if err != nil {
		return fmt.Errorf("failed to create tracks roadmap_id index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_tracks_status ON tracks(status)")
	if err != nil {
		return fmt.Errorf("failed to create tracks status index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_tracks_rank ON tracks(rank)")
	if err != nil {
		return fmt.Errorf("failed to create tracks rank index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_track_id ON tasks(track_id)")
	if err != nil {
		return fmt.Errorf("failed to create tasks track_id index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)")
	if err != nil {
		return fmt.Errorf("failed to create tasks status index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_rank ON tasks(rank)")
	if err != nil {
		return fmt.Errorf("failed to create tasks rank index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_iterations_status ON iterations(status)")
	if err != nil {
		return fmt.Errorf("failed to create iterations status index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_iterations_rank ON iterations(rank)")
	if err != nil {
		return fmt.Errorf("failed to create iterations rank index: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	fmt.Println("✓ Migration to schema v4 complete!")
	return nil
}
