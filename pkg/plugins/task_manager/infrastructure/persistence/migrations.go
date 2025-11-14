package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

const (
	// SchemaVersion is the current database schema version
	SchemaVersion = 6
)

// SQL table creation statements
const (
	createRoadmapsTable = `
CREATE TABLE IF NOT EXISTS roadmaps (
    id TEXT PRIMARY KEY,
    vision TEXT NOT NULL,
    success_criteria TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
`

	createTracksTable = `
CREATE TABLE IF NOT EXISTS tracks (
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
`

	createTrackDependenciesTable = `
CREATE TABLE IF NOT EXISTS track_dependencies (
    track_id TEXT NOT NULL,
    depends_on_id TEXT NOT NULL,
    PRIMARY KEY (track_id, depends_on_id),
    FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE,
    FOREIGN KEY (depends_on_id) REFERENCES tracks(id) ON DELETE CASCADE
)
`

	createTasksTable = `
CREATE TABLE IF NOT EXISTS tasks (
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
`

	createIterationsTable = `
CREATE TABLE IF NOT EXISTS iterations (
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
`

	createIterationTasksTable = `
CREATE TABLE IF NOT EXISTS iteration_tasks (
    iteration_number INTEGER NOT NULL,
    task_id TEXT NOT NULL,
    PRIMARY KEY (iteration_number, task_id),
    FOREIGN KEY (iteration_number) REFERENCES iterations(number) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
)
`

	createProjectMetadataTable = `
CREATE TABLE IF NOT EXISTS project_metadata (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
)
`

	createAcceptanceCriteriaTable = `
CREATE TABLE IF NOT EXISTS acceptance_criteria (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    description TEXT NOT NULL,
    verification_type TEXT NOT NULL,
    status TEXT NOT NULL,
    notes TEXT,
    testing_instructions TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
)
`

	// Indexes for common queries
	createTracksRoadmapIDIndex = `
CREATE INDEX IF NOT EXISTS idx_tracks_roadmap_id ON tracks(roadmap_id)
`

	createTracksStatusIndex = `
CREATE INDEX IF NOT EXISTS idx_tracks_status ON tracks(status)
`

	createTracksRankIndex = `
CREATE INDEX IF NOT EXISTS idx_tracks_rank ON tracks(rank)
`

	createTasksTrackIDIndex = `
CREATE INDEX IF NOT EXISTS idx_tasks_track_id ON tasks(track_id)
`

	createTasksStatusIndex = `
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)
`

	createTasksRankIndex = `
CREATE INDEX IF NOT EXISTS idx_tasks_rank ON tasks(rank)
`

	createIterationsStatusIndex = `
CREATE INDEX IF NOT EXISTS idx_iterations_status ON iterations(status)
`

	createIterationsRankIndex = `
CREATE INDEX IF NOT EXISTS idx_iterations_rank ON iterations(rank)
`

	createIterationTasksIterationIndex = `
CREATE INDEX IF NOT EXISTS idx_iteration_tasks_iteration ON iteration_tasks(iteration_number)
`

	createIterationTasksTaskIndex = `
CREATE INDEX IF NOT EXISTS idx_iteration_tasks_task ON iteration_tasks(task_id)
`

	createAcceptanceCriteriaTaskIDIndex = `
CREATE INDEX IF NOT EXISTS idx_ac_task_id ON acceptance_criteria(task_id)
`

	createAcceptanceCriteriaStatusIndex = `
CREATE INDEX IF NOT EXISTS idx_ac_status ON acceptance_criteria(status)
`

	createADRsTable = `
CREATE TABLE IF NOT EXISTS adrs (
    id TEXT PRIMARY KEY,
    track_id TEXT NOT NULL,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    context TEXT NOT NULL,
    decision TEXT NOT NULL,
    consequences TEXT NOT NULL,
    alternatives TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    superseded_by TEXT,
    FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE,
    FOREIGN KEY(superseded_by) REFERENCES adrs(id) ON DELETE SET NULL
)
`

	createADRsTrackIDIndex = `
CREATE INDEX IF NOT EXISTS idx_adrs_track_id ON adrs(track_id)
`

	createADRsStatusIndex = `
CREATE INDEX IF NOT EXISTS idx_adrs_status ON adrs(status)
`

	createDocumentsTable = `
CREATE TABLE IF NOT EXISTS documents (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL CHECK(length(title) > 0 AND length(title) <= 200),
    type TEXT NOT NULL CHECK(type IN ('adr', 'plan', 'retrospective', 'other')),
    status TEXT NOT NULL CHECK(status IN ('draft', 'published', 'archived')),
    content TEXT NOT NULL,
    track_id TEXT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    iteration_number INTEGER NULL REFERENCES iterations(number) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata TEXT DEFAULT '{}',
    CHECK(NOT (track_id IS NOT NULL AND iteration_number IS NOT NULL))
)
`

	createDocumentsTrackIDIndex = `
CREATE INDEX IF NOT EXISTS idx_documents_track_id ON documents(track_id)
`

	createDocumentsIterationNumberIndex = `
CREATE INDEX IF NOT EXISTS idx_documents_iteration_number ON documents(iteration_number)
`

	createDocumentsTypeIndex = `
CREATE INDEX IF NOT EXISTS idx_documents_type ON documents(type)
`
)

// InitSchema initializes the database schema with all required tables and indexes.
// It's safe to call multiple times (uses IF NOT EXISTS).
func InitSchema(db *sql.DB) error {
	// First create project_metadata table if it doesn't exist
	if _, err := db.Exec(createProjectMetadataTable); err != nil {
		return fmt.Errorf("failed to create project_metadata table: %w", err)
	}

	// Check if we need to migrate from version 3 to version 4 (priority -> rank)
	var currentVersion int
	err := db.QueryRow("SELECT CAST(value AS INTEGER) FROM project_metadata WHERE key = 'schema_version'").Scan(&currentVersion)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check schema version: %w", err)
	}

	// If no version found, check if we have old tables with priority column
	if err == sql.ErrNoRows || currentVersion == 0 {
		// Check if tracks table exists and has priority column
		var hasPriorityColumn bool
		rows, err := db.Query("PRAGMA table_info(tracks)")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var cid int
				var name, typ string
				var notnull, pk int
				var dfltValue sql.NullString
				if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err == nil {
					if name == "priority" {
						hasPriorityColumn = true
						break
					}
				}
			}
			rows.Close()
		}

		// If we found priority column, it's a v3 database that needs migration
		if hasPriorityColumn {
			currentVersion = 3
		}
	}

	// If we have version 3, run migration
	if currentVersion == 3 {
		if err := migrateV3ToV4(db); err != nil {
			return fmt.Errorf("failed to migrate from v3 to v4: %w", err)
		}
		currentVersion = 4
	}

	// If we have version 4, run migration
	if currentVersion == 4 {
		if err := migrateV4ToV5(db); err != nil {
			return fmt.Errorf("failed to migrate from v4 to v5: %w", err)
		}
	}

	statements := []string{
		createRoadmapsTable,
		createTracksTable,
		createTrackDependenciesTable,
		createTasksTable,
		createIterationsTable,
		createIterationTasksTable,
		createProjectMetadataTable,
		createAcceptanceCriteriaTable,
		createADRsTable,
		createDocumentsTable,
		createTracksRoadmapIDIndex,
		createTracksStatusIndex,
		createTracksRankIndex,
		createTasksTrackIDIndex,
		createTasksStatusIndex,
		createTasksRankIndex,
		createIterationsStatusIndex,
		createIterationsRankIndex,
		createIterationTasksIterationIndex,
		createIterationTasksTaskIndex,
		createAcceptanceCriteriaTaskIDIndex,
		createAcceptanceCriteriaStatusIndex,
		createADRsTrackIDIndex,
		createADRsStatusIndex,
		createDocumentsTrackIDIndex,
		createDocumentsIterationNumberIndex,
		createDocumentsTypeIndex,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	// Update schema version
	_, err = db.Exec("INSERT OR REPLACE INTO project_metadata (key, value) VALUES ('schema_version', ?)", SchemaVersion)
	if err != nil {
		return fmt.Errorf("failed to update schema version: %w", err)
	}

	return nil
}

// MigrateFromFileStorage migrates existing task JSON files to the database.
// It creates a "legacy-tasks" track if needed and imports all tasks from the file storage directory.
func MigrateFromFileStorage(db *sql.DB, tasksDir string) error {
	// Check if tasks directory exists
	if _, err := os.Stat(tasksDir); os.IsNotExist(err) {
		// No existing tasks to migrate
		return nil
	}

	// First, check if there are any tasks already in the database
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing tasks: %w", err)
	}

	// If database already has tasks, skip migration
	if count > 0 {
		return nil
	}

	// Read task files from directory
	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		return fmt.Errorf("failed to read tasks directory: %w", err)
	}

	// Check if there are any task files
	taskFiles := []os.DirEntry{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			taskFiles = append(taskFiles, entry)
		}
	}

	// If no task files, nothing to migrate
	if len(taskFiles) == 0 {
		return nil
	}

	// Create a default roadmap for legacy tasks
	legacyRoadmapID := "legacy-roadmap"
	legacyTrackID := "track-legacy-tasks"

	// Check if legacy roadmap exists
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM roadmaps WHERE id = ?", legacyRoadmapID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for legacy roadmap: %w", err)
	}

	// Only create if doesn't exist
	if exists == 0 {
		roadmap, err := entities.NewRoadmapEntity(
			legacyRoadmapID,
			"Legacy Tasks from File Storage",
			"Migrate existing tasks to database",
			GetCurrentTime(),
			GetCurrentTime(),
		)
		if err != nil {
			return fmt.Errorf("failed to create legacy roadmap: %w", err)
		}

		_, err = db.Exec(
			"INSERT INTO roadmaps (id, vision, success_criteria, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			roadmap.ID, roadmap.Vision, roadmap.SuccessCriteria, roadmap.CreatedAt, roadmap.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert legacy roadmap: %w", err)
		}

		// Create a legacy track
		track, err := entities.NewTrackEntity(
			legacyTrackID,
			legacyRoadmapID,
			"Legacy Tasks",
			"Tasks migrated from file-based storage",
			"not-started",
			300, // low priority = 300 rank
			[]string{},
			GetCurrentTime(),
			GetCurrentTime(),
		)
		if err != nil {
			return fmt.Errorf("failed to create legacy track: %w", err)
		}

		_, err = db.Exec(
			"INSERT INTO tracks (id, roadmap_id, title, description, status, rank, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			track.ID, track.RoadmapID, track.Title, track.Description, track.Status, track.Rank, track.CreatedAt, track.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert legacy track: %w", err)
		}
	}

	// Migrate task files
	migratedCount := 0
	for _, entry := range taskFiles {
		taskPath := filepath.Join(tasksDir, entry.Name())
		data, err := os.ReadFile(taskPath)
		if err != nil {
			// Log error but continue with next file
			fmt.Printf("Warning: failed to read task file %s: %v\n", entry.Name(), err)
			continue
		}

		// Unmarshal JSON
		var oldTask entities.TaskEntity
		if err := json.Unmarshal(data, &oldTask); err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to parse task file %s: %v\n", entry.Name(), err)
			continue
		}

		// Insert into database (force legacy track assignment)
		_, err = db.Exec(
			"INSERT INTO tasks (id, track_id, title, description, status, rank, branch, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			oldTask.ID, legacyTrackID, oldTask.Title, oldTask.Description, oldTask.Status, oldTask.Rank, oldTask.Branch, oldTask.CreatedAt, oldTask.UpdatedAt,
		)
		if err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to migrate task %s: %v\n", oldTask.ID, err)
			continue
		}

		migratedCount++
	}

	if migratedCount > 0 {
		fmt.Printf("Migrated %d tasks to database\n", migratedCount)
	}

	return nil
}

// GetCurrentTime returns the current time in UTC.
// This is a helper function for consistent timestamp handling.
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

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

// migrateV4ToV5 migrates database from schema version 4 to version 5
// Adds testing_instructions column to acceptance_criteria table
func migrateV4ToV5(db *sql.DB) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if acceptance_criteria table has testing_instructions column (version 5)
	var hasTestingInstructions bool
	rows, err := tx.Query("PRAGMA table_info(acceptance_criteria)")
	if err != nil {
		return fmt.Errorf("failed to check acceptance_criteria table: %w", err)
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
		if name == "testing_instructions" {
			hasTestingInstructions = true
			break
		}
	}
	rows.Close()

	if hasTestingInstructions {
		// Already migrated or new database
		return tx.Commit()
	}

	fmt.Println("Migrating database from schema v4 to v5 (adding testing_instructions column)...")

	// MIGRATE ACCEPTANCE_CRITERIA TABLE
	// 1. Create new acceptance_criteria table with testing_instructions
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS acceptance_criteria_new (
			id TEXT PRIMARY KEY,
			task_id TEXT NOT NULL,
			description TEXT NOT NULL,
			verification_type TEXT NOT NULL,
			status TEXT NOT NULL,
			notes TEXT,
			testing_instructions TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create new acceptance_criteria table: %w", err)
	}

	// 2. Copy data from old table to new table
	_, err = tx.Exec(`
		INSERT INTO acceptance_criteria_new (id, task_id, description, verification_type, status, notes, testing_instructions, created_at, updated_at)
		SELECT
			id,
			task_id,
			description,
			verification_type,
			status,
			notes,
			'' as testing_instructions,
			created_at,
			updated_at
		FROM acceptance_criteria
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate acceptance_criteria data: %w", err)
	}

	// 3. Drop old table and rename new one
	if _, err = tx.Exec("DROP TABLE acceptance_criteria"); err != nil {
		return fmt.Errorf("failed to drop old acceptance_criteria table: %w", err)
	}
	if _, err = tx.Exec("ALTER TABLE acceptance_criteria_new RENAME TO acceptance_criteria"); err != nil {
		return fmt.Errorf("failed to rename new acceptance_criteria table: %w", err)
	}

	// 4. Recreate indexes
	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_ac_task_id ON acceptance_criteria(task_id)")
	if err != nil {
		return fmt.Errorf("failed to create ac_task_id index: %w", err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_ac_status ON acceptance_criteria(status)")
	if err != nil {
		return fmt.Errorf("failed to create ac_status index: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	fmt.Println("✓ Migration to schema v5 complete!")
	return nil
}
