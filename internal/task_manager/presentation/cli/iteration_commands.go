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
// NewIterationCommands returns the iteration command group for Cobra
// ============================================================================

// NewIterationCommands creates and returns the iteration command group with all subcommands.
func NewIterationCommands(iterationService *application.IterationApplicationService, docService *application.DocumentApplicationService, acService *application.ACApplicationService) *cobra.Command {
	iterCmd := &cobra.Command{
		Use:     "iteration",
		Short:   "Manage iterations",
		Long:    "Commands for creating, updating, and managing iterations",
		Aliases: []string{"iter"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all iteration subcommands
	iterCmd.AddCommand(
		newIterationCreateCommand(iterationService),
		newIterationListCommand(iterationService),
		newIterationShowCommand(iterationService, docService),
		newIterationCurrentCommand(iterationService, acService),
		newIterationStartCommand(iterationService),
		newIterationCompleteCommand(iterationService),
		newIterationAddTaskCommand(iterationService),
		newIterationRemoveTaskCommand(iterationService),
		newIterationDeleteCommand(iterationService),
		newIterationUpdateCommand(iterationService),
	)

	return iterCmd
}

// ============================================================================
// iteration create command
// ============================================================================

func newIterationCreateCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new iteration",
		Long:  `Creates a new iteration with auto-incremented number.`,
		Example: `  # Create a simple iteration
  tm iteration create --name "Sprint 1" --goal "Complete core features" --deliverable "MVP release"

  # Create with custom rank
  tm iteration create --name "Sprint 2" --goal "Bug fixes" --deliverable "Patch release" --rank 100`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name, _ := cmd.Flags().GetString("name")
			goal, _ := cmd.Flags().GetString("goal")
			deliverable, _ := cmd.Flags().GetString("deliverable")

			// Validate required flags
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if goal == "" {
				return fmt.Errorf("--goal is required")
			}
			if deliverable == "" {
				return fmt.Errorf("--deliverable is required")
			}

			// Create iteration via application service
			input := dto.CreateIterationDTO{
				Name:        name,
				Goal:        goal,
				Deliverable: deliverable,
			}

			iteration, err := iterationService.CreateIteration(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to create iteration: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Iteration created successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Number:      %d\n", iteration.Number)
			fmt.Fprintf(cmd.OutOrStdout(), "  Name:        %s\n", iteration.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "  Goal:        %s\n", iteration.Goal)
			fmt.Fprintf(cmd.OutOrStdout(), "  Deliverable: %s\n", iteration.Deliverable)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", iteration.Status)

			return nil
		},
	}

	cmd.Flags().String("name", "", "Iteration name (required)")
	cmd.Flags().String("goal", "", "Iteration goal (required)")
	cmd.Flags().String("deliverable", "", "Deliverable description (required)")
	cmd.Flags().Int("rank", 500, "Iteration rank (1-1000, default: 500)")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("goal")
	cmd.MarkFlagRequired("deliverable")

	return cmd
}

// ============================================================================
// iteration list command
// ============================================================================

func newIterationListCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all iterations",
		Long:  `Lists all iterations with their status, goal, and task count.`,
		Example: `  # List all iterations
  tm iteration list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Execute via application service
			iterations, err := iterationService.ListIterations(ctx)
			if err != nil {
				return fmt.Errorf("failed to list iterations: %w", err)
			}

			// Format output
			if len(iterations) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No iterations found\n")
				return nil
			}

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "%-5s %-25s %-20s %-12s\n",
				"#", "Name", "Goal", "Status")
			fmt.Fprintf(cmd.OutOrStdout(), "%-5s %-25s %-20s %-12s\n",
				strings.Repeat("-", 5),
				strings.Repeat("-", 25),
				strings.Repeat("-", 20),
				strings.Repeat("-", 12),
			)

			// Print iterations
			for _, iter := range iterations {
				name := truncateString(iter.Name, 25)
				goal := truncateString(iter.Goal, 20)

				fmt.Fprintf(cmd.OutOrStdout(), "%-5d %-25s %-20s %-12s\n",
					iter.Number,
					name,
					goal,
					iter.Status,
				)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d iteration(s)\n", len(iterations))
			return nil
		},
	}

	return cmd
}

// ============================================================================
// iteration show command
// ============================================================================

func newIterationShowCommand(iterationService *application.IterationApplicationService, docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <iteration-number>",
		Short: "Show details of a specific iteration",
		Long:  `Displays detailed information about a specific iteration including its status, goal, and associated tasks.`,
		Example: `  # Show iteration details
  tm iteration show 1

  # Show detailed information
  tm iteration show 1 --full`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var number int
			_, err := fmt.Sscanf(args[0], "%d", &number)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %w", err)
			}

			full, _ := cmd.Flags().GetBool("full")

			// Execute via application service
			iteration, err := iterationService.GetIteration(ctx, number)
			if err != nil {
				return fmt.Errorf("failed to get iteration: %w", err)
			}

			// Get tasks in iteration
			tasks, err := iterationService.GetIterationTasks(ctx, number)
			if err != nil {
				return fmt.Errorf("failed to get iteration tasks: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Iteration Details\n")
			fmt.Fprintf(cmd.OutOrStdout(), "=================\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Number:      %d\n", iteration.Number)
			fmt.Fprintf(cmd.OutOrStdout(), "  Name:        %s\n", iteration.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "  Goal:        %s\n", iteration.Goal)
			fmt.Fprintf(cmd.OutOrStdout(), "  Deliverable: %s\n", iteration.Deliverable)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", iteration.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Task Count:  %d\n", len(tasks))

			// Display tasks if any (always show, not just with --full)
			if len(tasks) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\n  Tasks:\n")
				for _, task := range tasks {
					fmt.Fprintf(cmd.OutOrStdout(), "    - %s (%s): %s\n", task.ID, task.Status, task.Title)
				}
			}

			if full {
				fmt.Fprintf(cmd.OutOrStdout(), "  Created:     %s\n", iteration.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
				fmt.Fprintf(cmd.OutOrStdout(), "  Updated:     %s\n", iteration.UpdatedAt.Format("2006-01-02 15:04:05 UTC"))

				if iteration.StartedAt != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  Started:     %s\n", iteration.StartedAt.Format("2006-01-02 15:04:05 UTC"))
				}
				if iteration.CompletedAt != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  Completed:   %s\n", iteration.CompletedAt.Format("2006-01-02 15:04:05 UTC"))
				}
			}

			// Show attached documents
			docs, err := docService.ListDocuments(ctx, nil, &number, nil)
			if err == nil && len(docs) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\nAttached Documents:\n")
				for _, doc := range docs {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s  %s  %s\n", doc.ID, doc.Title, doc.Type, doc.Status)
				}
			}

			return nil
		},
	}

	cmd.Flags().Bool("full", false, "Show full details including timestamps and tasks")

	return cmd
}

// ============================================================================
// iteration current command
// ============================================================================

func newIterationCurrentCommand(iterationService *application.IterationApplicationService, acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the current active iteration",
		Long:  `Displays the iteration with status 'current'. If no current iteration exists, shows the next planned iteration.`,
		Example: `  # Show current iteration
  tm iteration current

  # Show detailed information
  tm iteration current --full`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			full, _ := cmd.Flags().GetBool("full")

			// Execute via application service
			result, err := iterationService.GetCurrentIteration(ctx)
			if err != nil {
				return fmt.Errorf("failed to get current iteration: %w", err)
			}

			// Check if no iterations found at all
			if result.Iteration == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", result.FallbackMsg)
				return nil
			}

			// Cast iteration to entity (it's always *entities.IterationEntity when not nil)
			iteration, ok := result.Iteration.(*entities.IterationEntity)
			if !ok {
				return fmt.Errorf("unexpected iteration type")
			}

			// Display fallback message if applicable
			if result.IsFallback {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", result.FallbackMsg)
			}

			// Get tasks in iteration
			tasks, err := iterationService.GetIterationTasks(ctx, iteration.Number)
			if err != nil {
				return fmt.Errorf("failed to get iteration tasks: %w", err)
			}

			// Format output
			if result.IsFallback {
				fmt.Fprintf(cmd.OutOrStdout(), "Iteration: %d - %s\n", iteration.Number, iteration.Name)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Current Iteration: %d - %s\n", iteration.Number, iteration.Name)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", iteration.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Goal:        %s\n", iteration.Goal)
			fmt.Fprintf(cmd.OutOrStdout(), "  Deliverable: %s\n", iteration.Deliverable)
			fmt.Fprintf(cmd.OutOrStdout(), "  Task Count:  %d\n", len(tasks))

			if full {
				fmt.Fprintf(cmd.OutOrStdout(), "  Created:     %s\n", iteration.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
				fmt.Fprintf(cmd.OutOrStdout(), "  Updated:     %s\n", iteration.UpdatedAt.Format("2006-01-02 15:04:05 UTC"))

				if iteration.StartedAt != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  Started:     %s\n", iteration.StartedAt.Format("2006-01-02 15:04:05 UTC"))
				}
				if iteration.CompletedAt != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "  Completed:   %s\n", iteration.CompletedAt.Format("2006-01-02 15:04:05 UTC"))
				}

				// Display tasks if any
				if len(tasks) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "\n  Tasks:\n")
					for _, task := range tasks {
						fmt.Fprintf(cmd.OutOrStdout(), "    - %s (%s): %s\n", task.ID, task.Status, task.Title)

						// Fetch and display ACs for this task
						acs, err := acService.ListAC(ctx, task.ID)
						if err == nil && len(acs) > 0 {
							fmt.Fprintf(cmd.OutOrStdout(), "      Acceptance Criteria:\n")
							for _, ac := range acs {
								statusIcon := getStatusIndicator(ac.Status)
								fmt.Fprintf(cmd.OutOrStdout(), "        %s [%s] %s\n", statusIcon, ac.ID, ac.Description)
							}
						}
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().Bool("full", false, "Show full details including timestamps and tasks")

	return cmd
}

// ============================================================================
// iteration start command
// ============================================================================

func newIterationStartCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <iteration-number>",
		Short: "Start an iteration (set as current)",
		Long:  `Sets an iteration as the current active iteration. Only one iteration can be current at a time.`,
		Example: `  # Start iteration 1
  tm iteration start 1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var number int
			_, err := fmt.Sscanf(args[0], "%d", &number)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %w", err)
			}

			// Execute via application service
			if err := iterationService.StartIteration(ctx, number); err != nil {
				return fmt.Errorf("failed to start iteration: %w", err)
			}

			// Get updated iteration for output
			iteration, err := iterationService.GetIteration(ctx, number)
			if err != nil {
				return fmt.Errorf("failed to get iteration: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Iteration %d started successfully\n", iteration.Number)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status: %s\n", iteration.Status)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// iteration complete command
// ============================================================================

func newIterationCompleteCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete <iteration-number>",
		Short: "Mark an iteration as complete",
		Long:  `Marks an iteration as complete. Sets the completion timestamp.`,
		Example: `  # Complete iteration 1
  tm iteration complete 1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var number int
			_, err := fmt.Sscanf(args[0], "%d", &number)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %w", err)
			}

			// Execute via application service
			if err := iterationService.CompleteIteration(ctx, number); err != nil {
				return fmt.Errorf("failed to complete iteration: %w", err)
			}

			// Get updated iteration for output
			iteration, err := iterationService.GetIteration(ctx, number)
			if err != nil {
				return fmt.Errorf("failed to get iteration: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Iteration %d completed successfully\n", iteration.Number)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status: %s\n", iteration.Status)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// iteration add-task command
// ============================================================================

func newIterationAddTaskCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-task <iteration-number> <task-id> [<task-id>...]",
		Short: "Add task(s) to an iteration",
		Long:  `Adds one or more tasks to an iteration.`,
		Example: `  # Add single task
  tm iteration add-task 1 TM-task-1

  # Add multiple tasks
  tm iteration add-task 1 TM-task-1 TM-task-2 TM-task-3`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var number int
			_, err := fmt.Sscanf(args[0], "%d", &number)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %w", err)
			}

			taskIDs := args[1:]
			successCount := 0
			var lastErr error

			// Add each task
			for _, taskID := range taskIDs {
				if err := iterationService.AddTask(ctx, number, taskID); err != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "Failed to add task %s: %v\n", taskID, err)
					lastErr = err
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Added task %s to iteration %d\n", taskID, number)
					successCount++
				}
			}

			if successCount > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Successfully added %d task(s)\n", successCount)
			}

			if lastErr != nil && successCount == 0 {
				return lastErr
			}

			return nil
		},
	}

	return cmd
}

// ============================================================================
// iteration remove-task command
// ============================================================================

func newIterationRemoveTaskCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-task <iteration-number> <task-id> [<task-id>...]",
		Short: "Remove task(s) from an iteration",
		Long:  `Removes one or more tasks from an iteration. The task itself is not deleted, only the iteration association.`,
		Example: `  # Remove single task
  tm iteration remove-task 1 TM-task-1

  # Remove multiple tasks
  tm iteration remove-task 1 TM-task-1 TM-task-2`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var number int
			_, err := fmt.Sscanf(args[0], "%d", &number)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %w", err)
			}

			taskIDs := args[1:]
			successCount := 0
			var lastErr error

			// Remove each task
			for _, taskID := range taskIDs {
				if err := iterationService.RemoveTask(ctx, number, taskID); err != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "Failed to remove task %s: %v\n", taskID, err)
					lastErr = err
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Removed task %s from iteration %d\n", taskID, number)
					successCount++
				}
			}

			if successCount > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed %d task(s)\n", successCount)
			}

			if lastErr != nil && successCount == 0 {
				return lastErr
			}

			return nil
		},
	}

	return cmd
}

// ============================================================================
// iteration delete command
// ============================================================================

func newIterationDeleteCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <iteration-number>",
		Short: "Delete an iteration",
		Long:  `Deletes an iteration from the project and removes its task associations.`,
		Example: `  # Delete iteration 1
  tm iteration delete 1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var number int
			_, err := fmt.Sscanf(args[0], "%d", &number)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %w", err)
			}

			// Execute via application service
			if err := iterationService.DeleteIteration(ctx, number); err != nil {
				return fmt.Errorf("failed to delete iteration: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Iteration %d deleted successfully\n", number)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// iteration update command
// ============================================================================

func newIterationUpdateCommand(iterationService *application.IterationApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <iteration-number>",
		Short: "Update an existing iteration",
		Long:  `Updates one or more fields of an existing iteration. At least one field must be specified.`,
		Example: `  # Update iteration name
  tm iteration update 1 --name "Sprint 1 - Updated"

  # Update multiple fields
  tm iteration update 1 --name "Sprint 2" --goal "Refactoring" --rank 200`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var number int
			_, err := fmt.Sscanf(args[0], "%d", &number)
			if err != nil {
				return fmt.Errorf("invalid iteration number: %w", err)
			}

			nameSet := cmd.Flags().Changed("name")
			goalSet := cmd.Flags().Changed("goal")
			deliverableSet := cmd.Flags().Changed("deliverable")
			rankSet := cmd.Flags().Changed("rank")

			// Check that at least one field is being updated
			if !nameSet && !goalSet && !deliverableSet && !rankSet {
				return fmt.Errorf("at least one field must be specified to update (--name, --goal, --deliverable, or --rank)")
			}

			// Create DTO with only updated fields
			input := dto.UpdateIterationDTO{
				Number: number,
			}

			if nameSet {
				name, _ := cmd.Flags().GetString("name")
				input.Name = &name
			}
			if goalSet {
				goal, _ := cmd.Flags().GetString("goal")
				input.Goal = &goal
			}
			if deliverableSet {
				deliverable, _ := cmd.Flags().GetString("deliverable")
				input.Deliverable = &deliverable
			}
			// (rank field is currently ignored in the DTO, but we keep the flag for future compatibility)

			// Execute via application service
			iteration, err := iterationService.UpdateIteration(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to update iteration: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Iteration updated successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Number:      %d\n", iteration.Number)
			fmt.Fprintf(cmd.OutOrStdout(), "  Name:        %s\n", iteration.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "  Goal:        %s\n", iteration.Goal)
			fmt.Fprintf(cmd.OutOrStdout(), "  Deliverable: %s\n", iteration.Deliverable)

			return nil
		},
	}

	cmd.Flags().String("name", "", "New iteration name")
	cmd.Flags().String("goal", "", "New iteration goal")
	cmd.Flags().String("deliverable", "", "New deliverable description")
	cmd.Flags().Int("rank", 0, "New iteration rank (1-1000)")

	return cmd
}
