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
// NewACCommands returns the AC command group for Cobra
// ============================================================================

// NewACCommands creates and returns the acceptance criteria command group with all subcommands.
func NewACCommands(acService *application.ACApplicationService, taskService *application.TaskApplicationService) *cobra.Command {
	acCmd := &cobra.Command{
		Use:     "ac",
		Short:   "Manage acceptance criteria",
		Long:    "Commands for creating, verifying, and managing acceptance criteria for tasks",
		Aliases: []string{"criterion", "criteria"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all AC subcommands
	acCmd.AddCommand(
		newACAddCommand(acService),
		newACListCommand(acService),
		newACListIterationCommand(acService),
		newACShowCommand(acService),
		newACUpdateCommand(acService),
		newACVerifyCommand(acService),
		newACFailCommand(acService),
		newACSkipCommand(acService),
		newACFailedCommand(acService),
		newACDeleteCommand(acService),
	)

	return acCmd
}

// ============================================================================
// ac add command
// ============================================================================

func newACAddCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <task-id>",
		Short: "Add an acceptance criterion to a task",
		Long:  `Adds an acceptance criterion to a task with description and optional testing instructions.`,
		Example: `  # Add simple AC
  tm ac add TM-task-1 --description "User can log in"

  # Add AC with testing instructions
  tm ac add TM-task-1 --description "User can log in" --testing-instructions "1. Click login\n2. Enter credentials\n3. Verify redirected"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			taskID := args[0]

			description, _ := cmd.Flags().GetString("description")
			testingInstructions, _ := cmd.Flags().GetString("testing-instructions")

			// Validate required flags
			if description == "" {
				return fmt.Errorf("--description is required")
			}

			// Create AC via application service
			input := dto.CreateACDTO{
				TaskID:              taskID,
				Description:         description,
				TestingInstructions: testingInstructions,
			}

			ac, err := acService.CreateAC(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to add acceptance criterion: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance criterion added successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", ac.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Task:        %s\n", ac.TaskID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", ac.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", ac.Status)
			if ac.TestingInstructions != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Testing:     %s\n", ac.TestingInstructions)
			}

			return nil
		},
	}

	cmd.Flags().String("description", "", "AC description (required)")
	cmd.Flags().String("testing-instructions", "", "Step-by-step testing instructions (optional)")

	cmd.MarkFlagRequired("description")

	return cmd
}

// ============================================================================
// ac list command
// ============================================================================

func newACListCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <task-id>",
		Short: "List acceptance criteria for a task",
		Long:  `Lists all acceptance criteria for a task with their verification status.`,
		Example: `  # List all ACs for a task
  tm ac list TM-task-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			taskID := args[0]

			// Get ACs for task via application service
			acs, err := acService.ListAC(ctx, taskID)
			if err != nil {
				return fmt.Errorf("failed to list ACs: %w", err)
			}

			if len(acs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No acceptance criteria found for task %s\n", taskID)
				return nil
			}

			// Count verified
			verifiedCount := 0
			for _, ac := range acs {
				if ac.Status == "verified" || ac.Status == "automatically-verified" {
					verifiedCount++
				}
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance Criteria for Task: %s\n", taskID)
			fmt.Fprintf(cmd.OutOrStdout(), "Summary: %d/%d verified\n\n", verifiedCount, len(acs))

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "%-10s %-20s %-50s %-15s\n", "Status", "ID", "Description", "Status")
			fmt.Fprintf(cmd.OutOrStdout(), "%-10s %-20s %-50s %-15s\n",
				strings.Repeat("-", 10),
				strings.Repeat("-", 20),
				strings.Repeat("-", 50),
				strings.Repeat("-", 15),
			)

			// Print ACs
			for _, ac := range acs {
				statusIcon := getStatusIndicator(ac.Status)
				fmt.Fprintf(cmd.OutOrStdout(), "%-10s %-20s %-50s %-15s\n",
					statusIcon,
					ac.ID,
					truncateString(ac.Description, 50),
					ac.Status,
				)
			}

			return nil
		},
	}

	return cmd
}

// ============================================================================
// ac list-iteration command
// ============================================================================

func newACListIterationCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-iteration <iteration-number>",
		Short: "List acceptance criteria for an iteration",
		Long:  `Lists all acceptance criteria for all tasks in an iteration, grouped by task with status indicators.`,
		Example: `  # List ACs for iteration 1
  tm ac list-iteration 1

  # List ACs with testing instructions
  tm ac list-iteration 1 --with-testing`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			withTesting, _ := cmd.Flags().GetBool("with-testing")

			// Parse iteration number
			var iteration int
			_, err := fmt.Sscanf(args[0], "%d", &iteration)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %s", args[0])
			}

			// Get ACs for iteration via application service
			acs, err := acService.ListACByIteration(ctx, iteration)
			if err != nil {
				return fmt.Errorf("failed to get ACs for iteration: %w", err)
			}

			if len(acs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Iteration %d has no acceptance criteria\n", iteration)
				return nil
			}

			// Count verification status
			var verifiedCount, pendingCount, failedCount, notStartedCount int
			for _, ac := range acs {
				switch ac.Status {
				case "verified", "automatically-verified":
					verifiedCount++
				case "pending-review":
					pendingCount++
				case "failed":
					failedCount++
				default:
					notStartedCount++
				}
			}

			// Group ACs by task
			acsByTask := make(map[string][]*entities.AcceptanceCriteriaEntity)
			for _, ac := range acs {
				acsByTask[ac.TaskID] = append(acsByTask[ac.TaskID], ac)
			}

			// Display results
			fmt.Fprintf(cmd.OutOrStdout(), "Iteration %d\n", iteration)
			fmt.Fprintf(cmd.OutOrStdout(), "\nAcceptance Criteria Summary:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ✓  Verified:        %d\n", verifiedCount)
			fmt.Fprintf(cmd.OutOrStdout(), "  ⏸  Pending Review:  %d\n", pendingCount)
			fmt.Fprintf(cmd.OutOrStdout(), "  ✗  Failed:          %d\n", failedCount)
			fmt.Fprintf(cmd.OutOrStdout(), "  ○  Not Started:     %d\n", notStartedCount)
			fmt.Fprintf(cmd.OutOrStdout(), "  Total:              %d\n", len(acs))

			fmt.Fprintf(cmd.OutOrStdout(), "\nAcceptance Criteria by Task:\n\n")

			for taskID, taskACs := range acsByTask {
				fmt.Fprintf(cmd.OutOrStdout(), "Task: %s\n", taskID)

				for _, ac := range taskACs {
					statusIcon := getStatusIndicator(ac.Status)
					fmt.Fprintf(cmd.OutOrStdout(), "  %s [%s] %s\n", statusIcon, ac.ID, ac.Description)

					// Display testing instructions if flag is set and AC has them
					if withTesting && ac.TestingInstructions != "" {
						fmt.Fprintf(cmd.OutOrStdout(), "    Testing Instructions:\n")
						// Indent each line of testing instructions with 6 spaces
						lines := strings.Split(ac.TestingInstructions, "\n")
						for _, line := range lines {
							if line != "" {
								fmt.Fprintf(cmd.OutOrStdout(), "      %s\n", line)
							}
						}
						fmt.Fprintf(cmd.OutOrStdout(), "\n")
					}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\n")
			}

			return nil
		},
	}

	cmd.Flags().Bool("with-testing", false, "Include testing instructions in output")

	return cmd
}

// ============================================================================
// ac show command
// ============================================================================

func newACShowCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <ac-id>",
		Short: "Show details of an acceptance criterion",
		Long:  `Displays detailed information about an acceptance criterion including description, testing instructions, and status.`,
		Example: `  # Show AC details
  tm ac show TM-ac-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			acID := args[0]

			// Get AC via application service
			ac, err := acService.GetAC(ctx, acID)
			if err != nil {
				return fmt.Errorf("failed to get AC: %w", err)
			}

			// Display AC details
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance Criterion Details\n")
			fmt.Fprintf(cmd.OutOrStdout(), "============================\n\n")

			fmt.Fprintf(cmd.OutOrStdout(), "ID:          %s\n", ac.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Task ID:     %s\n", ac.TaskID)
			fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", ac.Description)
			statusIcon := getStatusIndicator(ac.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "Status:      %s %s\n", statusIcon, ac.Status)

			// Show testing instructions if present
			if ac.TestingInstructions != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "\nTesting Instructions:\n")
				fmt.Fprintf(cmd.OutOrStdout(), "---------------------\n")
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", ac.TestingInstructions)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "\nTesting Instructions: (none)\n")
			}

			// Show failure notes if AC failed
			if ac.Status == "failed" && ac.Notes != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "\nFailure Feedback:\n")
				fmt.Fprintf(cmd.OutOrStdout(), "-----------------\n")
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", ac.Notes)
			}

			// Show timestamps
			fmt.Fprintf(cmd.OutOrStdout(), "\nTimestamps:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "-----------\n")
			fmt.Fprintf(cmd.OutOrStdout(), "Created: %s\n", ac.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Fprintf(cmd.OutOrStdout(), "Updated: %s\n", ac.UpdatedAt.Format("2006-01-02 15:04:05"))

			return nil
		},
	}

	return cmd
}

// ============================================================================
// ac update command
// ============================================================================

func newACUpdateCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <ac-id>",
		Short: "Update an acceptance criterion",
		Long:  `Updates an acceptance criterion's description or testing instructions. At least one field must be specified.`,
		Example: `  # Update description
  tm ac update TM-ac-1 --description "Updated requirement"

  # Update testing instructions
  tm ac update TM-ac-1 --testing-instructions "1. New step\n2. Another step"

  # Update both
  tm ac update TM-ac-1 --description "New desc" --testing-instructions "New steps"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			acID := args[0]

			descSet := cmd.Flags().Changed("description")
			testSet := cmd.Flags().Changed("testing-instructions")

			// Check that at least one field is being updated
			if !descSet && !testSet {
				return fmt.Errorf("at least one field must be specified to update (--description or --testing-instructions)")
			}

			// Create DTO with only updated fields
			input := dto.UpdateACDTO{
				ID: acID,
			}

			if descSet {
				description, _ := cmd.Flags().GetString("description")
				input.Description = &description
			}
			if testSet {
				testingInstructions, _ := cmd.Flags().GetString("testing-instructions")
				input.TestingInstructions = &testingInstructions
			}

			// Execute via application service
			ac, err := acService.UpdateAC(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to update AC: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance criterion updated successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", ac.ID)
			if descSet {
				fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", ac.Description)
			}
			if testSet {
				if ac.TestingInstructions != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  Testing Instructions: Updated\n")
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "  Testing Instructions: Cleared\n")
				}
			}

			return nil
		},
	}

	cmd.Flags().String("description", "", "New AC description")
	cmd.Flags().String("testing-instructions", "", "New testing instructions")

	return cmd
}

// ============================================================================
// ac verify command
// ============================================================================

func newACVerifyCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify <ac-id>",
		Short: "Mark an acceptance criterion as verified",
		Long:  `Marks an acceptance criterion as verified after manual testing.`,
		Example: `  # Verify an AC
  tm ac verify TM-ac-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			acID := args[0]

			// Create DTO with verification metadata
			input := dto.VerifyACDTO{
				ID:         acID,
				VerifiedBy: "user",
				VerifiedAt: "now",
			}

			// Execute via application service
			if err := acService.VerifyAC(ctx, input); err != nil {
				return fmt.Errorf("failed to verify acceptance criterion: %w", err)
			}

			// Get updated AC for output
			ac, err := acService.GetAC(ctx, acID)
			if err != nil {
				return fmt.Errorf("failed to get AC: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance criterion verified successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:     %s\n", ac.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status: %s\n", ac.Status)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// ac fail command
// ============================================================================

func newACFailCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fail <ac-id>",
		Short: "Mark an acceptance criterion as failed",
		Long:  `Marks an acceptance criterion as failed with feedback explaining why.`,
		Example: `  # Mark AC as failed
  tm ac fail TM-ac-1 --feedback "Login button not responding"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			acID := args[0]

			feedback, _ := cmd.Flags().GetString("feedback")

			// Validate required flag
			if feedback == "" {
				return fmt.Errorf("--feedback is required")
			}

			// Create DTO
			input := dto.FailACDTO{
				ID:       acID,
				Feedback: feedback,
			}

			// Execute via application service
			if err := acService.FailAC(ctx, input); err != nil {
				return fmt.Errorf("failed to mark acceptance criterion as failed: %w", err)
			}

			// Get updated AC for output
			ac, err := acService.GetAC(ctx, acID)
			if err != nil {
				return fmt.Errorf("failed to get AC: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance criterion marked as failed\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:       %s\n", ac.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:   %s\n", ac.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Feedback: %s\n", feedback)

			return nil
		},
	}

	cmd.Flags().String("feedback", "", "Failure feedback (required)")
	cmd.MarkFlagRequired("feedback")

	return cmd
}

// ============================================================================
// ac skip command
// ============================================================================

func newACSkipCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skip <ac-id>",
		Short: "Mark an acceptance criterion as skipped",
		Long:  `Marks an acceptance criterion as skipped with a reason explaining why it's not applicable.`,
		Example: `  # Skip an AC
  tm ac skip TM-ac-1 --reason "Not applicable for this task"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			acID := args[0]

			reason, _ := cmd.Flags().GetString("reason")

			// Validate required flag
			if reason == "" {
				return fmt.Errorf("--reason is required")
			}

			// Create DTO
			input := dto.SkipACDTO{
				ID:     acID,
				Reason: reason,
			}

			// Execute via application service
			if err := acService.SkipAC(ctx, input); err != nil {
				return fmt.Errorf("failed to skip acceptance criterion: %w", err)
			}

			// Get updated AC for output
			ac, err := acService.GetAC(ctx, acID)
			if err != nil {
				return fmt.Errorf("failed to get AC: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance criterion skipped successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:     %s\n", ac.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status: %s\n", ac.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Reason: %s\n", reason)

			return nil
		},
	}

	cmd.Flags().String("reason", "", "Reason for skipping (required)")
	cmd.MarkFlagRequired("reason")

	return cmd
}

// ============================================================================
// ac failed command
// ============================================================================

func newACFailedCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "failed",
		Short: "List failed acceptance criteria with optional filtering",
		Long:  `Lists all acceptance criteria with status "failed" with optional filtering by iteration, track, or task.`,
		Example: `  # List all failed ACs
  tm ac failed

  # List failed ACs in iteration 3
  tm ac failed --iteration 3

  # List failed ACs for a specific track
  tm ac failed --track TM-track-1

  # List failed ACs for a specific task
  tm ac failed --task TM-task-1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hasIteration := cmd.Flags().Changed("iteration")
			hasTrack := cmd.Flags().Changed("track")
			hasTask := cmd.Flags().Changed("task")

			iterationNum, _ := cmd.Flags().GetInt("iteration")
			trackID, _ := cmd.Flags().GetString("track")
			taskID, _ := cmd.Flags().GetString("task")

			// Build filters
			var iterNumPtr *int
			if hasIteration {
				iterNumPtr = &iterationNum
			}

			filters := entities.ACFilters{
				IterationNum: iterNumPtr,
				TrackID:      trackID,
				TaskID:       taskID,
			}

			// Get failed ACs via application service
			failedACs, err := acService.ListFailedAC(ctx, filters)
			if err != nil {
				return fmt.Errorf("failed to list failed ACs: %w", err)
			}

			if len(failedACs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No failed acceptance criteria found")
				if hasIteration {
					fmt.Fprintf(cmd.OutOrStdout(), " for iteration %d", iterationNum)
				}
				if hasTrack {
					fmt.Fprintf(cmd.OutOrStdout(), " for track %s", trackID)
				}
				if hasTask {
					fmt.Fprintf(cmd.OutOrStdout(), " for task %s", taskID)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\n")
				return nil
			}

			// Display header
			fmt.Fprintf(cmd.OutOrStdout(), "Failed Acceptance Criteria")
			if hasIteration {
				fmt.Fprintf(cmd.OutOrStdout(), " (Iteration %d)", iterationNum)
			}
			if hasTrack {
				fmt.Fprintf(cmd.OutOrStdout(), " (Track: %s)", trackID)
			}
			if hasTask {
				fmt.Fprintf(cmd.OutOrStdout(), " (Task: %s)", taskID)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n")
			fmt.Fprintf(cmd.OutOrStdout(), "Total: %d\n\n", len(failedACs))

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-50s\n", "AC ID", "Task ID", "Description")
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-50s\n",
				strings.Repeat("-", 20),
				strings.Repeat("-", 20),
				strings.Repeat("-", 50),
			)

			// Display each failed AC
			for _, ac := range failedACs {
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-20s %-50s\n",
					ac.ID,
					ac.TaskID,
					truncateString(ac.Description, 50),
				)
				if ac.Notes != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  Feedback: %s\n", ac.Notes)
				}
			}

			return nil
		},
	}

	cmd.Flags().Int("iteration", 0, "Filter by iteration number (optional)")
	cmd.Flags().String("track", "", "Filter by track ID (optional)")
	cmd.Flags().String("task", "", "Filter by task ID (optional)")

	return cmd
}

// ============================================================================
// ac delete command
// ============================================================================

func newACDeleteCommand(acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <ac-id>",
		Short: "Delete an acceptance criterion",
		Long:  `Deletes an acceptance criterion. Requires the --force flag for safety.`,
		Example: `  # Delete an AC
  tm ac delete TM-ac-1 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			acID := args[0]

			force, _ := cmd.Flags().GetBool("force")

			// Validate --force flag
			if !force {
				return fmt.Errorf("--force flag is required to confirm deletion")
			}

			// Execute via application service
			if err := acService.DeleteAC(ctx, acID); err != nil {
				return fmt.Errorf("failed to delete AC: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Acceptance criterion deleted\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID: %s\n", acID)

			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Required flag to confirm deletion")
	cmd.MarkFlagRequired("force")

	return cmd
}

// ============================================================================
// Helper functions
// ============================================================================

func getStatusIndicator(status entities.AcceptanceCriteriaStatus) string {
	switch status {
	case entities.ACStatusVerified, entities.ACStatusAutomaticallyVerified:
		return "✓"
	case entities.ACStatusPendingHumanReview:
		return "⏸"
	case entities.ACStatusFailed:
		return "✗"
	case entities.ACStatusSkipped:
		return "⊘"
	default:
		return "○"
	}
}
