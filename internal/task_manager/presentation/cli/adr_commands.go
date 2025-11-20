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
// NewADRCommands returns the ADR command group for Cobra
// ============================================================================

// NewADRCommands creates and returns the ADR command group with all subcommands.
func NewADRCommands(adrService *application.ADRApplicationService) *cobra.Command {
	adrCmd := &cobra.Command{
		Use:     "adr",
		Short:   "Manage Architecture Decision Records",
		Long:    "Commands for creating, updating, and managing ADRs within tracks",
		Aliases: []string{"a"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all ADR subcommands
	adrCmd.AddCommand(
		newADRCreateCommand(adrService),
		newADRListCommand(adrService),
		newADRShowCommand(adrService),
		newADRUpdateCommand(adrService),
		newADRSupersedeCommand(adrService),
		newADRDeprecateCommand(adrService),
		newADRCheckCommand(adrService),
	)

	return adrCmd
}

// ============================================================================
// adr create command
// ============================================================================

func newADRCreateCommand(adrService *application.ADRApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [track-id]",
		Short: "Create a new Architecture Decision Record",
		Long:  `Creates a new ADR within the specified track with required context, decision, and consequences fields.`,
		Example: `  # Create a simple ADR with positional track ID
  tm adr create TM-track-1 --title "Use PostgreSQL" --context "Need database" --decision "Use PostgreSQL" --consequences "Need to manage PostgreSQL deployment"

  # Create ADR with --track flag
  tm adr create --track TM-track-1 --title "Use PostgreSQL" --context "Need database" --decision "Use PostgreSQL" --consequences "Need to manage PostgreSQL deployment"

  # Create ADR with alternatives
  tm adr create --track TM-track-1 --title "API Framework" --context "Need API layer" --decision "Use gRPC" --consequences "High performance" --alternatives "Considered REST and GraphQL"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get track ID from positional arg or flag
			var trackID string
			if len(args) > 0 {
				trackID = args[0]
			} else {
				trackID, _ = cmd.Flags().GetString("track")
			}

			// Retrieve other flags
			title, _ := cmd.Flags().GetString("title")
			context, _ := cmd.Flags().GetString("context")
			decision, _ := cmd.Flags().GetString("decision")
			consequences, _ := cmd.Flags().GetString("consequences")
			alternatives, _ := cmd.Flags().GetString("alternatives")

			// Validate required fields
			if trackID == "" {
				return fmt.Errorf("track ID is required (as positional argument or --track flag)")
			}
			if title == "" {
				return fmt.Errorf("--title is required")
			}
			if context == "" {
				return fmt.Errorf("--context is required")
			}
			if decision == "" {
				return fmt.Errorf("--decision is required")
			}
			if consequences == "" {
				return fmt.Errorf("--consequences is required")
			}

			// Create ADR via application service
			input := dto.CreateADRDTO{
				TrackID:      trackID,
				Title:        title,
				Context:      context,
				Decision:     decision,
				Consequences: consequences,
				Alternatives: alternatives,
				Status:       "proposed",
			}

			adr, err := adrService.CreateADR(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to create ADR: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "ADR created successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:           %s\n", adr.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Track:        %s\n", adr.TrackID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:        %s\n", adr.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:       %s\n", adr.Status)

			return nil
		},
	}

	cmd.Flags().String("track", "", "Parent track ID (optional if provided as positional argument)")
	cmd.Flags().String("title", "", "ADR title (required)")
	cmd.Flags().String("context", "", "Problem context (required)")
	cmd.Flags().String("decision", "", "Decision made (required)")
	cmd.Flags().String("consequences", "", "Decision consequences (required)")
	cmd.Flags().String("alternatives", "", "Alternative approaches considered (optional)")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("context")
	cmd.MarkFlagRequired("decision")
	cmd.MarkFlagRequired("consequences")

	return cmd
}

// ============================================================================
// adr list command
// ============================================================================

func newADRListCommand(adrService *application.ADRApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all ADRs with optional filtering",
		Long:  `Lists all ADRs with optional filtering by track.`,
		Example: `  # List all ADRs
  tm adr list

  # List ADRs for a specific track
  tm adr list --track TM-track-1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Retrieve flags
			trackID, _ := cmd.Flags().GetString("track")

			// Build filter
			var trackIDPtr *string
			if trackID != "" {
				trackIDPtr = &trackID
			}

			// Execute via application service
			adrs, err := adrService.ListADRs(ctx, trackIDPtr)
			if err != nil {
				return fmt.Errorf("failed to list ADRs: %w", err)
			}

			// Format output
			if len(adrs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No ADRs found\n")
				return nil
			}

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-40s %-15s\n", "ID", "Track", "Title", "Status")
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-40s %-15s\n",
				strings.Repeat("-", 20),
				strings.Repeat("-", 20),
				strings.Repeat("-", 40),
				strings.Repeat("-", 15),
			)

			// Print ADRs
			for _, adr := range adrs {
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-40s %-15s\n",
					adr.ID,
					adr.TrackID,
					truncateString(adr.Title, 40),
					adr.Status,
				)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d ADR(s)\n", len(adrs))
			return nil
		},
	}

	cmd.Flags().String("track", "", "Filter by parent track ID (optional)")

	return cmd
}

// ============================================================================
// adr show command
// ============================================================================

func newADRShowCommand(adrService *application.ADRApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <adr-id>",
		Short: "Show details of a specific ADR",
		Long:  `Displays detailed information about a specific ADR including context, decision, and consequences.`,
		Example: `  # Show ADR details
  tm adr show TM-adr-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			adrID := args[0]

			// Execute via application service
			adr, err := adrService.GetADR(ctx, adrID)
			if err != nil {
				return fmt.Errorf("failed to get ADR: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Architecture Decision Record\n")
			fmt.Fprintf(cmd.OutOrStdout(), "============================\n\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:           %s\n", adr.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Track:        %s\n", adr.TrackID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:        %s\n", adr.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:       %s\n", adr.Status)
			if adr.SupersededBy != nil && *adr.SupersededBy != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Superseded By: %s\n", *adr.SupersededBy)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nContext:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", adr.Context)

			fmt.Fprintf(cmd.OutOrStdout(), "\nDecision:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", adr.Decision)

			fmt.Fprintf(cmd.OutOrStdout(), "\nConsequences:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", adr.Consequences)

			if adr.Alternatives != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "\nAlternatives Considered:\n")
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", adr.Alternatives)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTimestamps:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Created: %s\n", adr.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
			fmt.Fprintf(cmd.OutOrStdout(), "  Updated: %s\n", adr.UpdatedAt.Format("2006-01-02 15:04:05 UTC"))

			return nil
		},
	}

	return cmd
}

// ============================================================================
// adr update command
// ============================================================================

func newADRUpdateCommand(adrService *application.ADRApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <adr-id>",
		Short: "Update an existing ADR",
		Long:  `Updates one or more fields of an existing ADR. At least one field must be specified.`,
		Example: `  # Update ADR title
  tm adr update TM-adr-1 --title "New Title"

  # Update ADR status
  tm adr update TM-adr-1 --status accepted

  # Update multiple fields
  tm adr update TM-adr-1 --title "Updated Title" --decision "New decision" --status accepted`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			adrID := args[0]

			// Check which flags were actually set
			titleSet := cmd.Flags().Changed("title")
			contextSet := cmd.Flags().Changed("context")
			decisionSet := cmd.Flags().Changed("decision")
			consSet := cmd.Flags().Changed("consequences")
			altSet := cmd.Flags().Changed("alternatives")
			statusSet := cmd.Flags().Changed("status")

			// Check that at least one field is being updated
			if !titleSet && !contextSet && !decisionSet && !consSet && !altSet && !statusSet {
				return fmt.Errorf("at least one field must be specified to update (--title, --context, --decision, --consequences, --alternatives, or --status)")
			}

			// Retrieve flags
			title, _ := cmd.Flags().GetString("title")
			context, _ := cmd.Flags().GetString("context")
			decision, _ := cmd.Flags().GetString("decision")
			consequences, _ := cmd.Flags().GetString("consequences")
			alternatives, _ := cmd.Flags().GetString("alternatives")
			status, _ := cmd.Flags().GetString("status")

			// Create DTO with only updated fields
			input := dto.UpdateADRDTO{
				ID: adrID,
			}

			if titleSet {
				input.Title = &title
			}
			if contextSet {
				input.Context = &context
			}
			if decisionSet {
				input.Decision = &decision
			}
			if consSet {
				input.Consequences = &consequences
			}
			if altSet {
				input.Alternatives = &alternatives
			}
			if statusSet {
				input.Status = &status
			}

			// Execute via application service
			adr, err := adrService.UpdateADR(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to update ADR: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "ADR updated successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:           %s\n", adr.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Track:        %s\n", adr.TrackID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:        %s\n", adr.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:       %s\n", adr.Status)

			return nil
		},
	}

	// Define flags without closure variables
	cmd.Flags().String("title", "", "New ADR title")
	cmd.Flags().String("context", "", "New problem context")
	cmd.Flags().String("decision", "", "New decision")
	cmd.Flags().String("consequences", "", "New consequences")
	cmd.Flags().String("alternatives", "", "New alternatives")
	cmd.Flags().String("status", "", "New status (proposed, accepted, superseded, deprecated, rejected)")

	return cmd
}

// ============================================================================
// adr supersede command
// ============================================================================

func newADRSupersedeCommand(adrService *application.ADRApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supersede <adr-id>",
		Short: "Mark an ADR as superseded by another",
		Long:  `Marks an ADR as superseded by a newer ADR when a new architectural decision replaces an old one.`,
		Example: `  # Mark TM-adr-1 as superseded by TM-adr-5
  tm adr supersede TM-adr-1 --superseded-by TM-adr-5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			adrID := args[0]

			// Retrieve flags
			supersededByID, _ := cmd.Flags().GetString("superseded-by")

			// Validate required flag
			if supersededByID == "" {
				return fmt.Errorf("--superseded-by is required")
			}

			// Execute via application service
			if err := adrService.SupersedeADR(ctx, adrID, supersededByID); err != nil {
				return fmt.Errorf("failed to supersede ADR: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "ADR superseded successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  %s is now superseded by %s\n", adrID, supersededByID)

			return nil
		},
	}

	cmd.Flags().String("superseded-by", "", "ID of the superseding ADR (required)")
	cmd.MarkFlagRequired("superseded-by")

	return cmd
}

// ============================================================================
// adr deprecate command
// ============================================================================

func newADRDeprecateCommand(adrService *application.ADRApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deprecate <adr-id>",
		Short: "Mark an ADR as deprecated",
		Long:  `Marks an ADR as deprecated when it is no longer relevant but isn't directly superseded by another ADR.`,
		Example: `  # Deprecate an ADR
  tm adr deprecate TM-adr-3`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			adrID := args[0]

			// Execute via application service
			if err := adrService.DeprecateADR(ctx, adrID); err != nil {
				return fmt.Errorf("failed to deprecate ADR: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "ADR deprecated successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID: %s\n", adrID)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// adr check command
// ============================================================================

func newADRCheckCommand(adrService *application.ADRApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "List all ADRs with their status",
		Long:  `Lists all ADRs with their current status, showing which are proposed, accepted, superseded, deprecated, or rejected.`,
		Example: `  # Check all ADRs
  tm adr check`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Execute via application service
			adrs, err := adrService.ListADRs(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to list ADRs: %w", err)
			}

			// Format output
			if len(adrs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No ADRs found\n")
				return nil
			}

			// Group by status
			statusMap := make(map[string][]*entities.ADREntity)
			for _, adr := range adrs {
				statusMap[adr.Status] = append(statusMap[adr.Status], adr)
			}

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "Architecture Decision Records Status\n")
			fmt.Fprintf(cmd.OutOrStdout(), "====================================\n\n")

			// Statuses in order
			statuses := []string{"proposed", "accepted", "superseded", "deprecated", "rejected"}

			for _, status := range statuses {
				items := statusMap[status]
				if len(items) == 0 {
					continue
				}

				fmt.Fprintf(cmd.OutOrStdout(), "%s (%d)\n", strings.ToUpper(status), len(items))
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-40s\n", "ID", "Track", "Title")
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-40s\n",
					strings.Repeat("-", 20),
					strings.Repeat("-", 20),
					strings.Repeat("-", 40),
				)

				for _, adr := range items {
					fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-40s\n",
						adr.ID,
						adr.TrackID,
						truncateString(adr.Title, 40),
					)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\n")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Total: %d ADR(s)\n", len(adrs))
			return nil
		},
	}

	return cmd
}
