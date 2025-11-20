package cli

import (
	"fmt"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/spf13/cobra"
)

// ============================================================================
// NewRoadmapCommands returns the roadmap command group for Cobra
// ============================================================================

// NewRoadmapCommands creates and returns the roadmap command group with all subcommands.
func NewRoadmapCommands(roadmapService *application.RoadmapApplicationService) *cobra.Command {
	roadmapCmd := &cobra.Command{
		Use:     "roadmap",
		Short:   "Manage roadmaps",
		Long:    "Commands for creating, updating, and managing roadmaps",
		Aliases: []string{"r"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all roadmap subcommands
	roadmapCmd.AddCommand(
		newRoadmapInitCommand(roadmapService),
		newRoadmapShowCommand(roadmapService),
		newRoadmapUpdateCommand(roadmapService),
	)

	return roadmapCmd
}

// ============================================================================
// roadmap init command
// ============================================================================

func newRoadmapInitCommand(roadmapService *application.RoadmapApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new roadmap",
		Long:  `Creates a new roadmap with a vision statement and success criteria. Only one roadmap can exist at a time.`,
		Example: `  # Create a simple roadmap
  tm roadmap init \
    --vision "Build extensible framework" \
    --success-criteria "Support 10 plugins"

  # With multi-line vision
  tm roadmap init \
    --vision "Create unified productivity platform" \
    --success-criteria "100% test coverage, zero violations"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get flags
			vision, _ := cmd.Flags().GetString("vision")
			successCriteria, _ := cmd.Flags().GetString("success-criteria")

			// Validate required flags
			if vision == "" {
				return fmt.Errorf("--vision is required")
			}
			if successCriteria == "" {
				return fmt.Errorf("--success-criteria is required")
			}

			// Create roadmap via application service
			input := dto.CreateRoadmapDTO{
				Vision:          vision,
				SuccessCriteria: successCriteria,
			}

			roadmap, err := roadmapService.InitRoadmap(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to create roadmap: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Roadmap created successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:                %s\n", roadmap.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Vision:            %s\n", roadmap.Vision)
			fmt.Fprintf(cmd.OutOrStdout(), "  Success Criteria:  %s\n", roadmap.SuccessCriteria)

			return nil
		},
	}

	cmd.Flags().String("vision", "", "The vision statement for the roadmap (required)")
	cmd.Flags().String("success-criteria", "", "Success criteria for the roadmap (required)")

	cmd.MarkFlagRequired("vision")
	cmd.MarkFlagRequired("success-criteria")

	return cmd
}

// ============================================================================
// roadmap show command
// ============================================================================

func newRoadmapShowCommand(roadmapService *application.RoadmapApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display the current roadmap",
		Long:  `Displays the details of the current active roadmap.`,
		Example: `  # Show current roadmap
  tm roadmap show`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get active roadmap via application service
			roadmap, err := roadmapService.GetRoadmap(ctx)
			if err != nil {
				return fmt.Errorf("failed to get roadmap: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Roadmap:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:                %s\n", roadmap.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Vision:            %s\n", roadmap.Vision)
			fmt.Fprintf(cmd.OutOrStdout(), "  Success Criteria:  %s\n", roadmap.SuccessCriteria)
			fmt.Fprintf(cmd.OutOrStdout(), "  Created:           %s\n", roadmap.CreatedAt.Format(time.RFC3339))
			fmt.Fprintf(cmd.OutOrStdout(), "  Updated:           %s\n", roadmap.UpdatedAt.Format(time.RFC3339))

			return nil
		},
	}

	return cmd
}

// ============================================================================
// roadmap update command
// ============================================================================

func newRoadmapUpdateCommand(roadmapService *application.RoadmapApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the current roadmap",
		Long:  `Updates properties of the current active roadmap. At least one flag must be provided to update.`,
		Example: `  # Update vision
  tm roadmap update --vision "Create unified platform"

  # Update both
  tm roadmap update \
    --vision "New vision" \
    --success-criteria "New criteria"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Check which flags were set
			visionSet := cmd.Flags().Changed("vision")
			successCriteriaSet := cmd.Flags().Changed("success-criteria")

			// Check that at least one field is being updated
			if !visionSet && !successCriteriaSet {
				return fmt.Errorf("at least one field must be specified to update (--vision or --success-criteria)")
			}

			// Get flag values
			vision, _ := cmd.Flags().GetString("vision")
			successCriteria, _ := cmd.Flags().GetString("success-criteria")

			// Create DTO with only updated fields
			input := dto.UpdateRoadmapDTO{}

			if visionSet {
				input.Vision = &vision
			}
			if successCriteriaSet {
				input.SuccessCriteria = &successCriteria
			}

			// Execute via application service
			roadmap, err := roadmapService.UpdateRoadmap(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to update roadmap: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Roadmap updated successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:                %s\n", roadmap.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Vision:            %s\n", roadmap.Vision)
			fmt.Fprintf(cmd.OutOrStdout(), "  Success Criteria:  %s\n", roadmap.SuccessCriteria)
			fmt.Fprintf(cmd.OutOrStdout(), "  Updated:           %s\n", roadmap.UpdatedAt.Format(time.RFC3339))

			return nil
		},
	}

	cmd.Flags().String("vision", "", "New vision statement")
	cmd.Flags().String("success-criteria", "", "New success criteria")

	return cmd
}
