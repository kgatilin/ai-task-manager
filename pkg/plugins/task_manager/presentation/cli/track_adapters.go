package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/dto"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ============================================================================
// TrackCreateCommandAdapter - Adapts CLI to CreateTrackCommand use case
// ============================================================================

// TrackCreateCommandAdapter adapts track create CLI command to application use case
type TrackCreateCommandAdapter struct {
	TrackService *application.TrackApplicationService

	// CLI flags (parsed from args)
	project     string
	title       string
	description string
	rank        int
}

func (c *TrackCreateCommandAdapter) GetName() string {
	return "track create"
}

func (c *TrackCreateCommandAdapter) GetDescription() string {
	return "Create a new track"
}

func (c *TrackCreateCommandAdapter) GetUsage() string {
	return "dw task-manager track create --title <title> [--description <desc>] [--rank <rank>]"
}

func (c *TrackCreateCommandAdapter) GetHelp() string {
	return `Creates a new track in the active roadmap.

A track represents a major work area with multiple tasks and iterations.
All tracks must belong to an active roadmap - create one first with
'dw task-manager roadmap init'.

Flags:
  --title <title>          Track title (required)
  --description <desc>     Track description (optional)
  --rank <rank>            Track rank (optional, default: 500)
                          Range: 1-1000 (lower = higher priority)
  --project <name>         Project name (optional, uses active project if not specified)

Examples:
  # Create a basic track
  dw task-manager track create --title "Plugin System"

  # Create with custom rank
  dw task-manager track create \
    --title "Plugin System" \
    --description "Implement extensible plugin architecture" \
    --rank 100

Notes:
  - Track ID is auto-generated in format: <PROJECT_CODE>-track-<number> (e.g., DW-track-1)
  - An active roadmap must exist (create with 'dw task-manager roadmap init')
  - Initial status is automatically set to 'not-started'
  - No dependencies are added initially (use track add-dependency)
  - Rank determines ordering: lower values appear first (1=highest, 1000=lowest)`
}

func (c *TrackCreateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	c.rank = 500 // default
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--title":
			if i+1 < len(args) {
				c.title = args[i+1]
				i++
			}
		case "--description":
			if i+1 < len(args) {
				c.description = args[i+1]
				i++
			}
		case "--rank":
			if i+1 < len(args) {
				var err error
				c.rank, err = strconv.Atoi(args[i+1])
				if err != nil || c.rank < 1 || c.rank > 1000 {
					return fmt.Errorf("invalid rank: must be between 1 and 1000")
				}
				i++
			}
		}
	}

	// Validate required flags
	if c.title == "" {
		return fmt.Errorf("--title is required")
	}

	// Get active roadmap ID
	// Note: The service requires roadmap ID to verify roadmap exists
	roadmap, err := c.TrackService.GetActiveRoadmap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active roadmap: %w (create one with 'dw task-manager roadmap init')", err)
	}

	// Create DTO with roadmap ID
	input := dto.CreateTrackDTO{
		RoadmapID:   roadmap.ID,
		Title:       c.title,
		Description: c.description,
		Status:      "not-started",
		Rank:        c.rank,
	}

	// Execute via application service
	track, err := c.TrackService.CreateTrack(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Track created successfully\n")
	fmt.Fprintf(out, "  ID:          %s\n", track.ID)
	fmt.Fprintf(out, "  Title:       %s\n", track.Title)
	fmt.Fprintf(out, "  Status:      %s\n", track.Status)
	fmt.Fprintf(out, "  Rank:        %d\n", track.Rank)
	if track.Description != "" {
		fmt.Fprintf(out, "  Description: %s\n", track.Description)
	}

	return nil
}

// ============================================================================
// TrackUpdateCommandAdapter - Adapts CLI to UpdateTrackCommand use case
// ============================================================================

// TrackUpdateCommandAdapter adapts track update CLI command to application use case
type TrackUpdateCommandAdapter struct {
	TrackService *application.TrackApplicationService

	// CLI flags (parsed from args)
	project     string
	trackID     string
	title       *string
	description *string
	status      *string
	rank        *int
}

func (c *TrackUpdateCommandAdapter) GetName() string {
	return "track update"
}

func (c *TrackUpdateCommandAdapter) GetDescription() string {
	return "Update an existing track"
}

func (c *TrackUpdateCommandAdapter) GetUsage() string {
	return "dw task-manager track update <track-id> [--title <title>] [--description <desc>] [--status <status>] [--rank <rank>]"
}

func (c *TrackUpdateCommandAdapter) GetHelp() string {
	return `Updates an existing track's fields.

At least one field must be specified to update.

Flags:
  --title <title>          New track title
  --description <desc>     New track description
  --status <status>        New track status (not-started, in-progress, complete, blocked, waiting)
  --rank <rank>            New track rank (1-1000, lower = higher priority)
  --project <name>         Project name (optional, uses active project if not specified)

Examples:
  # Update track title
  dw task-manager track update TM-track-1 --title "New Title"

  # Update multiple fields
  dw task-manager track update TM-track-1 \
    --title "Plugin Architecture" \
    --status in-progress \
    --rank 100

  # Mark track as complete
  dw task-manager track update TM-track-1 --status complete`
}

func (c *TrackUpdateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse track ID (first positional argument)
	if len(args) == 0 {
		return fmt.Errorf("track ID is required")
	}
	c.trackID = args[0]
	args = args[1:]

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--title":
			if i+1 < len(args) {
				val := args[i+1]
				c.title = &val
				i++
			}
		case "--description":
			if i+1 < len(args) {
				val := args[i+1]
				c.description = &val
				i++
			}
		case "--status":
			if i+1 < len(args) {
				val := args[i+1]
				c.status = &val
				i++
			}
		case "--rank":
			if i+1 < len(args) {
				rankVal, err := strconv.Atoi(args[i+1])
				if err != nil || rankVal < 1 || rankVal > 1000 {
					return fmt.Errorf("invalid rank: must be between 1 and 1000")
				}
				c.rank = &rankVal
				i++
			}
		}
	}

	// Validate at least one field is provided
	if c.title == nil && c.description == nil && c.status == nil && c.rank == nil {
		return fmt.Errorf("at least one field must be specified to update")
	}

	// Create DTO
	input := dto.UpdateTrackDTO{
		ID:          c.trackID,
		Title:       c.title,
		Description: c.description,
		Status:      c.status,
		Rank:        c.rank,
	}

	// Execute via application service
	track, err := c.TrackService.UpdateTrack(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Track updated successfully\n")
	fmt.Fprintf(out, "  ID:          %s\n", track.ID)
	fmt.Fprintf(out, "  Title:       %s\n", track.Title)
	fmt.Fprintf(out, "  Status:      %s\n", track.Status)
	fmt.Fprintf(out, "  Rank:        %d\n", track.Rank)
	if track.Description != "" {
		fmt.Fprintf(out, "  Description: %s\n", track.Description)
	}

	return nil
}

// ============================================================================
// TrackListCommandAdapter - Lists tracks with optional filtering
// ============================================================================

type TrackListCommandAdapter struct {
	TrackService *application.TrackApplicationService

	// CLI flags
	project string
	status  string
}

func (c *TrackListCommandAdapter) GetName() string {
	return "track list"
}

func (c *TrackListCommandAdapter) GetDescription() string {
	return "List tracks"
}

func (c *TrackListCommandAdapter) GetUsage() string {
	return "dw task-manager track list [--status <status>]"
}

func (c *TrackListCommandAdapter) GetHelp() string {
	return `Lists all tracks in the active roadmap with optional filtering.

Tracks are displayed sorted by rank (lower ranks first).

Flags:
  --status <status>      Filter by status (can be comma-separated)
                         Values: not-started, in-progress, complete, blocked, waiting
  --project <name>       Project name (optional, uses active project if not specified)

Examples:
  # List all tracks
  dw task-manager track list

  # List in-progress tracks
  dw task-manager track list --status in-progress

  # List multiple status values
  dw task-manager track list --status in-progress,blocked

Output:
  A table showing: ID, Title, Status, Rank, Dependencies count
  Tracks are ordered by rank (1=highest priority, 1000=lowest)`
}

func (c *TrackListCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--status":
			if i+1 < len(args) {
				c.status = args[i+1]
				i++
			}
		}
	}

	// Build filters
	filters := entities.TrackFilters{}
	if c.status != "" {
		filters.Status = strings.Split(strings.TrimSpace(c.status), ",")
		for i, s := range filters.Status {
			filters.Status[i] = strings.TrimSpace(s)
		}
	}

	// Get active roadmap ID
	roadmap, err := c.TrackService.GetActiveRoadmap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active roadmap: %w (create one with 'dw task-manager roadmap init')", err)
	}

	// List tracks via application service
	tracks, err := c.TrackService.ListTracks(ctx, roadmap.ID, filters)
	if err != nil {
		return fmt.Errorf("failed to list tracks: %w", err)
	}

	// Display tracks
	out := cmdCtx.GetStdout()
	if len(tracks) == 0 {
		fmt.Fprintf(out, "No tracks found.\n")
		return nil
	}

	// Print header
	fmt.Fprintf(out, "%-25s %-30s %-12s %-6s %s\n",
		"ID", "Title", "Status", "Rank", "Dependencies")
	fmt.Fprintf(out, "%s\n",
		strings.Repeat("-", 90))

	// Print each track
	for _, track := range tracks {
		depCount := len(track.Dependencies)
		depStr := fmt.Sprintf("%d", depCount)
		fmt.Fprintf(out, "%-25s %-30s %-12s %-6d %s\n",
			track.ID, truncateString(track.Title, 29), track.Status, track.Rank, depStr)
	}

	fmt.Fprintf(out, "\nTotal: %d track(s)\n", len(tracks))
	return nil
}

// Helper function to truncate strings for display
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ============================================================================
// TrackShowCommandAdapter - Shows detailed track information
// ============================================================================

type TrackShowCommandAdapter struct {
	TrackService     *application.TrackApplicationService
	DocumentService  *application.DocumentApplicationService

	// CLI flags
	project string
	trackID string
}

func (c *TrackShowCommandAdapter) GetName() string {
	return "track show"
}

func (c *TrackShowCommandAdapter) GetDescription() string {
	return "Show detailed information about a track"
}

func (c *TrackShowCommandAdapter) GetUsage() string {
	return "dw task-manager track show <track-id>"
}

func (c *TrackShowCommandAdapter) GetHelp() string {
	return `Shows detailed information about a track including dependencies and tasks.

Flags:
  --project <name>    Project name (optional)

Examples:
  # Show track details
  dw task-manager track show TM-track-1`
}

func (c *TrackShowCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse track ID
	if len(args) == 0 {
		return fmt.Errorf("track ID is required")
	}
	c.trackID = args[0]
	args = args[1:]

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		}
	}

	// Get track with tasks via application service
	track, err := c.TrackService.GetTrackWithTasks(ctx, c.trackID)
	if err != nil {
		return fmt.Errorf("failed to get track: %w", err)
	}

	// Display track details
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Track Details\n")
	fmt.Fprintf(out, "=============\n\n")
	fmt.Fprintf(out, "ID:          %s\n", track.ID)
	fmt.Fprintf(out, "Title:       %s\n", track.Title)
	fmt.Fprintf(out, "Status:      %s\n", track.Status)
	fmt.Fprintf(out, "Rank:        %d\n", track.Rank)
	if track.Description != "" {
		fmt.Fprintf(out, "Description: %s\n", track.Description)
	}

	// Show dependencies
	if len(track.Dependencies) > 0 {
		fmt.Fprintf(out, "\nDependencies:\n")
		for _, dep := range track.Dependencies {
			fmt.Fprintf(out, "  - %s\n", dep)
		}
	} else {
		fmt.Fprintf(out, "\nDependencies: None\n")
	}

	// Show attached documents
	fmt.Fprintf(out, "\nAttached Documents:\n")
	documents, err := c.DocumentService.ListDocuments(ctx, &c.trackID, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to list documents: %w", err)
	}

	if len(documents) == 0 {
		fmt.Fprintf(out, "  (none)\n")
	} else {
		// Print header
		fmt.Fprintf(out, "  %-20s %-30s %-12s %s\n",
			"ID", "Title", "Type", "Status")
		fmt.Fprintf(out, "  %s\n",
			strings.Repeat("-", 80))

		// Print each document
		for _, doc := range documents {
			fmt.Fprintf(out, "  %-20s %-30s %-12s %s\n",
				doc.ID, truncateString(doc.Title, 29), doc.Type, doc.Status)
		}
	}

	// Note: Task details would need to be fetched separately via TaskService
	// The TrackEntity doesn't embed task entities

	return nil
}

// ============================================================================
// TrackDeleteCommandAdapter - Deletes a track
// ============================================================================

type TrackDeleteCommandAdapter struct {
	TrackService *application.TrackApplicationService

	// CLI flags
	project string
	trackID string
	force   bool
}

func (c *TrackDeleteCommandAdapter) GetName() string {
	return "track delete"
}

func (c *TrackDeleteCommandAdapter) GetDescription() string {
	return "Delete a track"
}

func (c *TrackDeleteCommandAdapter) GetUsage() string {
	return "dw task-manager track delete <track-id> [--force]"
}

func (c *TrackDeleteCommandAdapter) GetHelp() string {
	return `Deletes a track from the roadmap.

Requires the --force flag for safety.

Flags:
  --force             Required to confirm deletion
  --project <name>    Project name (optional)

Examples:
  # Delete a track
  dw task-manager track delete TM-track-1 --force`
}

func (c *TrackDeleteCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse track ID
	if len(args) == 0 {
		return fmt.Errorf("track ID is required")
	}
	c.trackID = args[0]
	args = args[1:]

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--force":
			c.force = true
		}
	}

	// Validate --force flag
	if !c.force {
		return fmt.Errorf("--force flag is required to confirm deletion")
	}

	// Execute via application service
	if err := c.TrackService.DeleteTrack(ctx, c.trackID); err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Track deleted successfully\n")
	fmt.Fprintf(out, "  ID: %s\n", c.trackID)

	return nil
}

// ============================================================================
// TrackAddDependencyCommandAdapter - Adds a dependency between tracks
// ============================================================================

type TrackAddDependencyCommandAdapter struct {
	TrackService *application.TrackApplicationService

	// CLI flags
	project     string
	trackID     string
	dependsOnID string
}

func (c *TrackAddDependencyCommandAdapter) GetName() string {
	return "track add-dependency"
}

func (c *TrackAddDependencyCommandAdapter) GetDescription() string {
	return "Add a dependency between tracks"
}

func (c *TrackAddDependencyCommandAdapter) GetUsage() string {
	return "dw task-manager track add-dependency <track-id> <depends-on-id>"
}

func (c *TrackAddDependencyCommandAdapter) GetHelp() string {
	return `Adds a dependency relationship between two tracks.

This indicates that <track-id> depends on <depends-on-id> being completed first.
Circular dependencies are automatically detected and prevented.

Flags:
  --project <name>    Project name (optional)

Examples:
  # Track TM-track-2 depends on TM-track-1
  dw task-manager track add-dependency TM-track-2 TM-track-1`
}

func (c *TrackAddDependencyCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse track IDs
	if len(args) < 2 {
		return fmt.Errorf("both track-id and depends-on-id are required")
	}
	c.trackID = args[0]
	c.dependsOnID = args[1]
	args = args[2:]

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		}
	}

	// Execute via application service
	if err := c.TrackService.AddDependency(ctx, c.trackID, c.dependsOnID); err != nil {
		return fmt.Errorf("failed to add dependency: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Dependency added successfully\n")
	fmt.Fprintf(out, "  %s depends on %s\n", c.trackID, c.dependsOnID)

	return nil
}

// ============================================================================
// TrackRemoveDependencyCommandAdapter - Removes a dependency between tracks
// ============================================================================

type TrackRemoveDependencyCommandAdapter struct {
	TrackService *application.TrackApplicationService

	// CLI flags
	project     string
	trackID     string
	dependsOnID string
}

func (c *TrackRemoveDependencyCommandAdapter) GetName() string {
	return "track remove-dependency"
}

func (c *TrackRemoveDependencyCommandAdapter) GetDescription() string {
	return "Remove a dependency between tracks"
}

func (c *TrackRemoveDependencyCommandAdapter) GetUsage() string {
	return "dw task-manager track remove-dependency <track-id> <depends-on-id>"
}

func (c *TrackRemoveDependencyCommandAdapter) GetHelp() string {
	return `Removes a dependency relationship between two tracks.

Flags:
  --project <name>    Project name (optional)

Examples:
  # Remove dependency: TM-track-2 no longer depends on TM-track-1
  dw task-manager track remove-dependency TM-track-2 TM-track-1`
}

func (c *TrackRemoveDependencyCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse track IDs
	if len(args) < 2 {
		return fmt.Errorf("both track-id and depends-on-id are required")
	}
	c.trackID = args[0]
	c.dependsOnID = args[1]
	args = args[2:]

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		}
	}

	// Execute via application service
	if err := c.TrackService.RemoveDependency(ctx, c.trackID, c.dependsOnID); err != nil {
		return fmt.Errorf("failed to remove dependency: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Dependency removed successfully\n")
	fmt.Fprintf(out, "  %s no longer depends on %s\n", c.trackID, c.dependsOnID)

	return nil
}
