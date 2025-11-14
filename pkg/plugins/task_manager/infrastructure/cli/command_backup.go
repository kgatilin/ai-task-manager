package cli

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ============================================================================
// BackupCommand creates a backup of the task-manager database
// ============================================================================

type BackupCommand struct {
	Provider PluginProvider
	project  string
}

func (c *BackupCommand) GetName() string {
	return "backup"
}

func (c *BackupCommand) GetDescription() string {
	return "Create a backup of the task-manager database"
}

func (c *BackupCommand) GetUsage() string {
	return "dw task-manager backup [--project <name>]"
}

func (c *BackupCommand) GetHelp() string {
	return `Creates a timestamped backup of the task-manager database.

The backup is stored in ~/.darwinflow/task-manager/backups/<project-name>/
with a timestamp-based filename (e.g., task-manager-2025-11-03-143022.db).

Before creating the backup, the database integrity is validated using SQLite's
PRAGMA integrity_check. If the database is corrupted, the backup will fail.

Auto-cleanup keeps only the 10 most recent backups. Older backups are
automatically deleted.

Flags:
  --project <name>  Override active project (optional)

Examples:
  # Backup active project
  dw task-manager backup

  # Backup specific project
  dw task-manager backup --project production

Notes:
  - Backups are stored in ~/.darwinflow/task-manager/backups/<project-name>/
  - Only the 10 most recent backups are kept (older ones auto-deleted)
  - Database integrity is validated before backup creation
  - Filename format: task-manager-YYYY-MM-DD-HHMMSS.db`
}

func (c *BackupCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		if args[i] == "--project" && i+1 < len(args) {
			c.project = args[i+1]
			i++
		}
	}

	// Determine project name
	projectName := c.project
	if projectName == "" {
		var err error
		projectName, err = c.Provider.GetActiveProject()
		if err != nil {
			return fmt.Errorf("failed to get active project: %w", err)
		}
	}

	// Get database path
	dbPath := c.getDatabasePath(projectName)
	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("database not found for project '%s' (run 'dw task-manager roadmap init' first)", projectName)
		}
		return fmt.Errorf("failed to access database: %w", err)
	}

	// Validate database integrity
	if err := c.validateDatabaseIntegrity(dbPath); err != nil {
		return fmt.Errorf("database integrity check failed: %w", err)
	}

	// Create backup directory
	backupDir := c.getBackupDirectory(projectName)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02-150405")
	backupFilename := fmt.Sprintf("task-manager-%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupFilename)

	// Copy database to backup location
	if err := c.copyFile(dbPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Auto-cleanup old backups (keep only 10 most recent)
	if err := c.cleanupOldBackups(backupDir, 10); err != nil {
		// Log warning but don't fail the backup
		fmt.Fprintf(cmdCtx.GetStdout(), "Warning: failed to cleanup old backups: %v\n", err)
	}

	fmt.Fprintf(cmdCtx.GetStdout(), "Backup created successfully: %s\n", backupPath)
	fmt.Fprintf(cmdCtx.GetStdout(), "Project: %s\n", projectName)
	fmt.Fprintf(cmdCtx.GetStdout(), "Timestamp: %s\n", timestamp)

	return nil
}

// getDatabasePath returns the path to the database for the given project
func (c *BackupCommand) getDatabasePath(projectName string) string {
	return filepath.Join(c.Provider.GetWorkingDir(), ".darwinflow", "projects", projectName, "roadmap.db")
}

// getBackupDirectory returns the backup directory for the given project
func (c *BackupCommand) getBackupDirectory(projectName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to working directory if home dir unavailable
		return filepath.Join(c.Provider.GetWorkingDir(), ".darwinflow", "task-manager", "backups", projectName)
	}
	return filepath.Join(homeDir, ".darwinflow", "task-manager", "backups", projectName)
}

// validateDatabaseIntegrity validates the database using SQLite's PRAGMA integrity_check
func (c *BackupCommand) validateDatabaseIntegrity(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var result string
	err = db.QueryRow("PRAGMA integrity_check").Scan(&result)
	if err != nil {
		return fmt.Errorf("failed to run integrity check: %w", err)
	}

	if result != "ok" {
		return fmt.Errorf("integrity check failed: %s", result)
	}

	return nil
}

// copyFile copies a file from src to dst
func (c *BackupCommand) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Sync to ensure data is written to disk
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// cleanupOldBackups removes old backups, keeping only the most recent 'keep' backups
func (c *BackupCommand) cleanupOldBackups(backupDir string, keep int) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	// Filter only backup files
	var backupFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "task-manager-") && strings.HasSuffix(entry.Name(), ".db") {
			backupFiles = append(backupFiles, entry)
		}
	}

	// If we have fewer backups than the keep limit, nothing to delete
	if len(backupFiles) <= keep {
		return nil
	}

	// Sort by modification time (newest first)
	sort.Slice(backupFiles, func(i, j int) bool {
		infoI, _ := backupFiles[i].Info()
		infoJ, _ := backupFiles[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Delete old backups beyond the keep limit
	for i := keep; i < len(backupFiles); i++ {
		filePath := filepath.Join(backupDir, backupFiles[i].Name())
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to remove old backup %s: %w", backupFiles[i].Name(), err)
		}
	}

	return nil
}

// ============================================================================
// RestoreCommand restores a backup of the task-manager database
// ============================================================================

type RestoreCommand struct {
	Provider     PluginProvider
	project      string
	backupFile   string
	skipValidate bool
}

func (c *RestoreCommand) GetName() string {
	return "restore"
}

func (c *RestoreCommand) GetDescription() string {
	return "Restore a backup of the task-manager database"
}

func (c *RestoreCommand) GetUsage() string {
	return "dw task-manager restore <backup-file> [--project <name>] [--skip-validate]"
}

func (c *RestoreCommand) GetHelp() string {
	return `Restores a backup of the task-manager database.

Before restoring, the backup file is validated using SQLite's PRAGMA
integrity_check. If the backup is corrupted, the restore will fail.

A safety backup of the current database is created in the project folder
(next to the existing database) before the restore operation. The safety
backup is named: roadmap-safety-backup-YYYY-MM-DD-HHMMSS.db

Arguments:
  <backup-file>  Path to the backup file (can be relative or absolute)

Flags:
  --project <name>     Override active project (optional)
  --skip-validate      Skip integrity validation (not recommended)

Examples:
  # Restore from backup file
  dw task-manager restore ~/.darwinflow/task-manager/backups/default/task-manager-2025-11-03-143022.db

  # Restore specific project
  dw task-manager restore task-manager-2025-11-03-143022.db --project production

  # Restore without validation (not recommended)
  dw task-manager restore backup.db --skip-validate

Notes:
  - A safety backup is created before restore (in project folder)
  - Backup integrity is validated before restore (unless --skip-validate)
  - Current database is completely replaced
  - Safety backup filename: roadmap-safety-backup-YYYY-MM-DD-HHMMSS.db`
}

func (c *RestoreCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse arguments and flags
	if len(args) == 0 {
		return fmt.Errorf("backup file path is required")
	}
	c.backupFile = args[0]

	for i := 1; i < len(args); i++ {
		if args[i] == "--project" && i+1 < len(args) {
			c.project = args[i+1]
			i++
		} else if args[i] == "--skip-validate" {
			c.skipValidate = true
		}
	}

	// Determine project name
	projectName := c.project
	if projectName == "" {
		var err error
		projectName, err = c.Provider.GetActiveProject()
		if err != nil {
			return fmt.Errorf("failed to get active project: %w", err)
		}
	}

	// Verify backup file exists
	if _, err := os.Stat(c.backupFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("backup file not found: %s", c.backupFile)
		}
		return fmt.Errorf("failed to access backup file: %w", err)
	}

	// Validate backup file integrity (unless skipped)
	if !c.skipValidate {
		if err := c.validateDatabaseIntegrity(c.backupFile); err != nil {
			return fmt.Errorf("backup file integrity check failed: %w", err)
		}
	}

	// Get current database path
	dbPath := c.getDatabasePath(projectName)

	// Create safety backup of current database (if it exists)
	if _, err := os.Stat(dbPath); err == nil {
		timestamp := time.Now().Format("2006-01-02-150405")
		safetyBackupName := fmt.Sprintf("roadmap-safety-backup-%s.db", timestamp)
		safetyBackupPath := filepath.Join(filepath.Dir(dbPath), safetyBackupName)

		if err := c.copyFile(dbPath, safetyBackupPath); err != nil {
			return fmt.Errorf("failed to create safety backup: %w", err)
		}

		fmt.Fprintf(cmdCtx.GetStdout(), "Safety backup created: %s\n", safetyBackupPath)
	}

	// Copy backup file to database location
	if err := c.copyFile(c.backupFile, dbPath); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	fmt.Fprintf(cmdCtx.GetStdout(), "Database restored successfully from: %s\n", c.backupFile)
	fmt.Fprintf(cmdCtx.GetStdout(), "Project: %s\n", projectName)
	fmt.Fprintf(cmdCtx.GetStdout(), "Database: %s\n", dbPath)

	return nil
}

// getDatabasePath returns the path to the database for the given project
func (c *RestoreCommand) getDatabasePath(projectName string) string {
	return filepath.Join(c.Provider.GetWorkingDir(), ".darwinflow", "projects", projectName, "roadmap.db")
}

// validateDatabaseIntegrity validates the database using SQLite's PRAGMA integrity_check
func (c *RestoreCommand) validateDatabaseIntegrity(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var result string
	err = db.QueryRow("PRAGMA integrity_check").Scan(&result)
	if err != nil {
		return fmt.Errorf("failed to run integrity check: %w", err)
	}

	if result != "ok" {
		return fmt.Errorf("integrity check failed: %s", result)
	}

	return nil
}

// copyFile copies a file from src to dst
func (c *RestoreCommand) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Sync to ensure data is written to disk
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// ============================================================================
// BackupListCommand lists all available backups
// ============================================================================

type BackupListCommand struct {
	Provider PluginProvider
	project  string
}

func (c *BackupListCommand) GetName() string {
	return "backup list"
}

func (c *BackupListCommand) GetDescription() string {
	return "List all available backups"
}

func (c *BackupListCommand) GetUsage() string {
	return "dw task-manager backup list [--project <name>]"
}

func (c *BackupListCommand) GetHelp() string {
	return `Lists all available backups for the active or specified project.

Displays backup filename, size, and creation time. Backups are sorted
by creation time (newest first).

Flags:
  --project <name>  Override active project (optional)

Examples:
  # List backups for active project
  dw task-manager backup list

  # List backups for specific project
  dw task-manager backup list --project production

Output format:
  Backups for project: default
  Location: ~/.darwinflow/task-manager/backups/default/

  task-manager-2025-11-03-143022.db    1.2 MB    2025-11-03 14:30:22
  task-manager-2025-11-03-120000.db    1.1 MB    2025-11-03 12:00:00
  task-manager-2025-11-02-180000.db    1.0 MB    2025-11-02 18:00:00

  Total: 3 backups`
}

func (c *BackupListCommand) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		if args[i] == "--project" && i+1 < len(args) {
			c.project = args[i+1]
			i++
		}
	}

	// Determine project name
	projectName := c.project
	if projectName == "" {
		var err error
		projectName, err = c.Provider.GetActiveProject()
		if err != nil {
			return fmt.Errorf("failed to get active project: %w", err)
		}
	}

	// Get backup directory
	backupDir := c.getBackupDirectory(projectName)

	// Check if backup directory exists
	if _, err := os.Stat(backupDir); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(cmdCtx.GetStdout(), "No backups found for project: %s\n", projectName)
			fmt.Fprintf(cmdCtx.GetStdout(), "Location: %s\n", backupDir)
			return nil
		}
		return fmt.Errorf("failed to access backup directory: %w", err)
	}

	// Read backup files
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	// Filter and collect backup files with metadata
	type backupInfo struct {
		name    string
		size    int64
		modTime time.Time
	}
	var backups []backupInfo

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), "task-manager-") || !strings.HasSuffix(entry.Name(), ".db") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, backupInfo{
			name:    entry.Name(),
			size:    info.Size(),
			modTime: info.ModTime(),
		})
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].modTime.After(backups[j].modTime)
	})

	// Display results
	fmt.Fprintf(cmdCtx.GetStdout(), "Backups for project: %s\n", projectName)
	fmt.Fprintf(cmdCtx.GetStdout(), "Location: %s\n\n", backupDir)

	if len(backups) == 0 {
		fmt.Fprintf(cmdCtx.GetStdout(), "No backups found.\n")
		return nil
	}

	// Display backup list
	for _, backup := range backups {
		sizeStr := formatSize(backup.size)
		timeStr := backup.modTime.Format("2006-01-02 15:04:05")
		fmt.Fprintf(cmdCtx.GetStdout(), "%-40s  %8s  %s\n", backup.name, sizeStr, timeStr)
	}

	fmt.Fprintf(cmdCtx.GetStdout(), "\nTotal: %d backups\n", len(backups))

	return nil
}

// getBackupDirectory returns the backup directory for the given project
func (c *BackupListCommand) getBackupDirectory(projectName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to working directory if home dir unavailable
		return filepath.Join(c.Provider.GetWorkingDir(), ".darwinflow", "task-manager", "backups", projectName)
	}
	return filepath.Join(homeDir, ".darwinflow", "task-manager", "backups", projectName)
}

// formatSize formats a size in bytes to a human-readable string
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
