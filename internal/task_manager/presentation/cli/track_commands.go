package cli

import (
	"fmt"
	"strings"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/spf13/cobra"
)

// ============================================================================
// NewTrackCommands returns the track command group for Cobra
// ============================================================================

// NewTrackCommands creates and returns the track command group with all subcommands.
func NewTrackCommands(trackService *application.TrackApplicationService, docService *application.DocumentApplicationService) *cobra.Command {
	trackCmd := &cobra.Command{
		Use:     "track",
		Short:   "Manage tracks",
		Long:    "Commands for creating, updating, and managing tracks within roadmaps",
		Aliases: []string{"tr"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all track subcommands
	trackCmd.AddCommand(
		newTrackCreateCommand(trackService),
		newTrackListCommand(trackService),
		newTrackShowCommand(trackService, docService),
		newTrackUpdateCommand(trackService),
		newTrackDeleteCommand(trackService),
		newTrackAddDependencyCommand(trackService),
		newTrackRemoveDependencyCommand(trackService),
	)

	return trackCmd
}

// ============================================================================
// track create command
// ============================================================================

func newTrackCreateCommand(trackService *application.TrackApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new track",
		Long:  `Creates a new track within the active roadmap with optional description and rank.`,
		Example: `  # Create a simple track
  tm track create --title "Plugin System"

  # Create track with description and rank
  tm track create --title "Plugin System" --description "Implement extensible architecture" --rank 100`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get flags
			title, _ := cmd.Flags().GetString("title")
			description, _ := cmd.Flags().GetString("description")
			rank, _ := cmd.Flags().GetInt("rank")

			// Validate required flags
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			// Get active roadmap ID
			roadmap, err := trackService.GetActiveRoadmap(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active roadmap: %w (create one with 'tm roadmap init')", err)
			}

			// Create track via application service
			input := dto.CreateTrackDTO{
				RoadmapID:   roadmap.ID,
				Title:       title,
				Description: description,
				Status:      "not-started",
				Rank:        rank,
			}

			track, err := trackService.CreateTrack(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to create track: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Track created successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", track.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:       %s\n", track.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", track.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Rank:        %d\n", track.Rank)
			if track.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", track.Description)
			}

			return nil
		},
	}

	cmd.Flags().String("title", "", "Track title (required)")
	cmd.Flags().String("description", "", "Track description (optional)")
	cmd.Flags().Int("rank", 500, "Track rank (1-1000, default: 500)")

	cmd.MarkFlagRequired("title")

	return cmd
}

// ============================================================================
// track list command
// ============================================================================

func newTrackListCommand(trackService *application.TrackApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tracks with optional filtering",
		Long:  `Lists all tracks in the active roadmap with optional filtering by status.`,
		Example: `  # List all tracks
  tm track list

  # List in-progress tracks
  tm track list --status in-progress

  # List multiple status values
  tm track list --status in-progress,blocked`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get flags
			status, _ := cmd.Flags().GetString("status")

			// Build filters
			filters := entities.TrackFilters{}
			if status != "" {
				filters.Status = strings.Split(strings.TrimSpace(status), ",")
				for i, s := range filters.Status {
					filters.Status[i] = strings.TrimSpace(s)
				}
			}

			// Get active roadmap ID
			roadmap, err := trackService.GetActiveRoadmap(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active roadmap: %w (create one with 'tm roadmap init')", err)
			}

			// Execute via application service
			tracks, err := trackService.ListTracks(ctx, roadmap.ID, filters)
			if err != nil {
				return fmt.Errorf("failed to list tracks: %w", err)
			}

			// Format output
			if len(tracks) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No tracks found\n")
				return nil
			}

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-30s %-12s %-6s %s\n",
				"ID", "Title", "Status", "Rank", "Dependencies")
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
				strings.Repeat("-", 90))

			// Print tracks
			for _, track := range tracks {
				depCount := len(track.Dependencies)
				depStr := fmt.Sprintf("%d", depCount)
				fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-30s %-12s %-6d %s\n",
					track.ID, truncateString(track.Title, 29), track.Status, track.Rank, depStr)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d track(s)\n", len(tracks))
			return nil
		},
	}

	cmd.Flags().String("status", "", "Filter by status: not-started, in-progress, complete, blocked, waiting (optional)")

	return cmd
}

// ============================================================================
// track show command
// ============================================================================

func newTrackShowCommand(trackService *application.TrackApplicationService, docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <track-id>",
		Short: "Show details of a specific track",
		Long:  `Displays detailed information about a specific track including its status, rank, dependencies, and attached documents.`,
		Example: `  # Show track details
  tm track show TM-track-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			trackID := args[0]

			// Execute via application service
			track, err := trackService.GetTrackWithTasks(ctx, trackID)
			if err != nil {
				return fmt.Errorf("failed to get track: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Track Details\n")
			fmt.Fprintf(cmd.OutOrStdout(), "=============\n\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", track.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:       %s\n", track.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", track.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Rank:        %d\n", track.Rank)
			if track.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", track.Description)
			}

			// Show dependencies
			if len(track.Dependencies) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDependencies:\n")
				for _, dep := range track.Dependencies {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", dep)
				}
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "\nDependencies: None\n")
			}

			// Show attached documents
			docs, err := docService.ListDocuments(ctx, &trackID, nil, nil)
			if err == nil && len(docs) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\nAttached Documents:\n")
				for _, doc := range docs {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s  %s  %s\n", doc.ID, doc.Title, doc.Type, doc.Status)
				}
			}

			return nil
		},
	}

	return cmd
}

// ============================================================================
// track update command
// ============================================================================

func newTrackUpdateCommand(trackService *application.TrackApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <track-id>",
		Short: "Update an existing track",
		Long:  `Updates one or more fields of an existing track. At least one field must be specified.`,
		Example: `  # Update track title
  tm track update TM-track-1 --title "New Title"

  # Update track status
  tm track update TM-track-1 --status in-progress

  # Update multiple fields
  tm track update TM-track-1 --title "New Title" --status complete --rank 100`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			trackID := args[0]

			// Get flags
			title, _ := cmd.Flags().GetString("title")
			description, _ := cmd.Flags().GetString("description")
			status, _ := cmd.Flags().GetString("status")
			rank, _ := cmd.Flags().GetInt("rank")

			// Check which flags were actually set
			titleSet := cmd.Flags().Changed("title")
			descSet := cmd.Flags().Changed("description")
			statusSet := cmd.Flags().Changed("status")
			rankSet := cmd.Flags().Changed("rank")

			// Check that at least one field is being updated
			if !titleSet && !descSet && !statusSet && !rankSet {
				return fmt.Errorf("at least one field must be specified to update (--title, --description, --status, or --rank)")
			}

			// Create DTO with only updated fields
			input := dto.UpdateTrackDTO{
				ID: trackID,
			}

			if titleSet {
				input.Title = &title
			}
			if descSet {
				input.Description = &description
			}
			if statusSet {
				input.Status = &status
			}
			if rankSet {
				input.Rank = &rank
			}

			// Execute via application service
			track, err := trackService.UpdateTrack(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to update track: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Track updated successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", track.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:       %s\n", track.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", track.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Rank:        %d\n", track.Rank)
			if track.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", track.Description)
			}

			return nil
		},
	}

	// Use explicit flag registration
	cmd.Flags().String("title", "", "New track title")
	cmd.Flags().String("description", "", "New track description")
	cmd.Flags().String("status", "", "New track status (not-started, in-progress, complete, blocked, waiting)")
	cmd.Flags().Int("rank", 0, "New track rank (1-1000)")

	return cmd
}

// ============================================================================
// track delete command
// ============================================================================

func newTrackDeleteCommand(trackService *application.TrackApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <track-id>",
		Short: "Delete a track",
		Long:  `Deletes a track and removes it from the roadmap. Requires the --force flag for safety.`,
		Example: `  # Delete a track
  tm track delete TM-track-1 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			trackID := args[0]

			// Get flags
			force, _ := cmd.Flags().GetBool("force")

			// Validate --force flag
			if !force {
				return fmt.Errorf("--force flag is required to confirm deletion")
			}

			// Execute via application service
			if err := trackService.DeleteTrack(ctx, trackID); err != nil {
				return fmt.Errorf("failed to delete track: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Track %s deleted successfully\n", trackID)

			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Required to confirm deletion")

	return cmd
}

// ============================================================================
// track add-dependency command
// ============================================================================

func newTrackAddDependencyCommand(trackService *application.TrackApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-dependency <track-id> <depends-on-id>",
		Short: "Add a dependency between tracks",
		Long:  `Adds a dependency relationship between two tracks. This indicates that <track-id> depends on <depends-on-id> being completed first. Circular dependencies are automatically detected and prevented.`,
		Example: `  # Track TM-track-2 depends on TM-track-1
  tm track add-dependency TM-track-2 TM-track-1`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			trackID := args[0]
			dependsOnID := args[1]

			// Execute via application service
			if err := trackService.AddDependency(ctx, trackID, dependsOnID); err != nil {
				return fmt.Errorf("failed to add dependency: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Dependency added successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  %s depends on %s\n", trackID, dependsOnID)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// track remove-dependency command
// ============================================================================

func newTrackRemoveDependencyCommand(trackService *application.TrackApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-dependency <track-id> <depends-on-id>",
		Short: "Remove a dependency between tracks",
		Long:  `Removes a dependency relationship between two tracks.`,
		Example: `  # Remove dependency: TM-track-2 no longer depends on TM-track-1
  tm track remove-dependency TM-track-2 TM-track-1`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			trackID := args[0]
			dependsOnID := args[1]

			// Execute via application service
			if err := trackService.RemoveDependency(ctx, trackID, dependsOnID); err != nil {
				return fmt.Errorf("failed to remove dependency: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Dependency removed successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  %s no longer depends on %s\n", trackID, dependsOnID)

			return nil
		},
	}

	return cmd
}
