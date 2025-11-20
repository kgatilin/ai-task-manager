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
// NewTaskCommands returns the task command group for Cobra
// ============================================================================

// NewTaskCommands creates and returns the task command group with all subcommands.
func NewTaskCommands(taskService *application.TaskApplicationService, acService *application.ACApplicationService) *cobra.Command {
	taskCmd := &cobra.Command{
		Use:     "task",
		Short:   "Manage tasks",
		Long:    "Commands for creating, updating, and managing tasks within tracks",
		Aliases: []string{"t"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all task subcommands
	taskCmd.AddCommand(
		newTaskCreateCommand(taskService),
		newTaskListCommand(taskService),
		newTaskShowCommand(taskService),
		newTaskUpdateCommand(taskService),
		newTaskDeleteCommand(taskService),
		newTaskMoveCommand(taskService),
		newTaskBacklogCommand(taskService),
		newTaskCheckReadyCommand(taskService, acService),
	)

	return taskCmd
}

// ============================================================================
// task create command
// ============================================================================

func newTaskCreateCommand(taskService *application.TaskApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new task in a track",
		Long:  `Creates a new task within the specified track with optional description, rank, and branch name.`,
		Example: `  # Create a simple task
  tm task create --track TM-track-1 --title "Implement login"

  # Create task with description and rank
  tm task create --track TM-track-1 --title "Add tests" --description "Unit tests for auth" --rank 300`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get flags
			trackID, _ := cmd.Flags().GetString("track")
			title, _ := cmd.Flags().GetString("title")
			description, _ := cmd.Flags().GetString("description")
			rank, _ := cmd.Flags().GetInt("rank")
			branch, _ := cmd.Flags().GetString("branch")

			// Validate required flags
			if trackID == "" {
				return fmt.Errorf("--track is required")
			}
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			// Create task via application service
			input := dto.CreateTaskDTO{
				TrackID:     trackID,
				Title:       title,
				Description: description,
				Status:      "todo",
				Rank:        rank,
			}
			_ = branch // Branch field not in DTO, reserved for future use

			task, err := taskService.CreateTask(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Task created successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", task.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Track:       %s\n", task.TrackID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:       %s\n", task.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", task.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Rank:        %d\n", task.Rank)
			if task.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", task.Description)
			}
			if task.Branch != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Branch:      %s\n", task.Branch)
			}

			return nil
		},
	}

	cmd.Flags().String("track", "", "Parent track ID (required)")
	cmd.Flags().String("title", "", "Task title (required)")
	cmd.Flags().String("description", "", "Task description (optional)")
	cmd.Flags().Int("rank", 500, "Task rank (1-1000, default: 500)")
	cmd.Flags().String("branch", "", "Git branch name (optional)")

	cmd.MarkFlagRequired("track")
	cmd.MarkFlagRequired("title")

	return cmd
}

// ============================================================================
// task list command
// ============================================================================

func newTaskListCommand(taskService *application.TaskApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks with optional filtering",
		Long:  `Lists all tasks with optional filtering by track or status.`,
		Example: `  # List all tasks
  tm task list

  # List tasks in a specific track
  tm task list --track TM-track-1

  # List tasks with specific status
  tm task list --status todo

  # Combine filters
  tm task list --track TM-track-1 --status in-progress`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get flags
			trackID, _ := cmd.Flags().GetString("track")
			status, _ := cmd.Flags().GetString("status")

			// Build filters
			filters := entities.TaskFilters{
				TrackID: trackID,
			}
			if status != "" {
				filters.Status = []string{status}
			}

			// Execute via application service
			tasks, err := taskService.ListTasks(ctx, filters)
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			// Format output
			if len(tasks) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No tasks found\n")
				return nil
			}

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-20s %-15s %-40s\n", "ID", "Track", "Status", "Title")
			fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-20s %-15s %-40s\n",
				strings.Repeat("-", 15),
				strings.Repeat("-", 20),
				strings.Repeat("-", 15),
				strings.Repeat("-", 40),
			)

			// Print tasks
			for _, task := range tasks {
				fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-20s %-15s %-40s\n",
					task.ID,
					task.TrackID,
					task.Status,
					truncateString(task.Title, 40),
				)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d task(s)\n", len(tasks))
			return nil
		},
	}

	cmd.Flags().String("track", "", "Filter by parent track ID (optional)")
	cmd.Flags().String("status", "", "Filter by status: todo, in-progress, review, done (optional)")

	return cmd
}

// ============================================================================
// task show command
// ============================================================================

func newTaskShowCommand(taskService *application.TaskApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <task-id>",
		Short: "Show details of a specific task",
		Long:  `Displays detailed information about a specific task including its status, rank, and metadata.`,
		Example: `  # Show task details
  tm task show TM-task-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			taskID := args[0]

			// Execute via application service
			task, err := taskService.GetTask(ctx, taskID)
			if err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Task Details\n")
			fmt.Fprintf(cmd.OutOrStdout(), "============\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", task.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Track:       %s\n", task.TrackID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:       %s\n", task.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", task.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", task.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Rank:        %d\n", task.Rank)
			if task.Branch != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Branch:      %s\n", task.Branch)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  Created:     %s\n", task.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
			fmt.Fprintf(cmd.OutOrStdout(), "  Updated:     %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05 UTC"))

			return nil
		},
	}

	return cmd
}

// ============================================================================
// task update command
// ============================================================================

func newTaskUpdateCommand(taskService *application.TaskApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <task-id>",
		Short: "Update an existing task",
		Long:  `Updates one or more fields of an existing task. At least one field must be specified.`,
		Example: `  # Update task title
  tm task update TM-task-1 --title "New Title"

  # Update task status
  tm task update TM-task-1 --status in-progress

  # Update multiple fields
  tm task update TM-task-1 --title "New Title" --status done --rank 100`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			taskID := args[0]

			// Check which flags were set
			titleSet := cmd.Flags().Changed("title")
			descSet := cmd.Flags().Changed("description")
			statusSet := cmd.Flags().Changed("status")
			rankSet := cmd.Flags().Changed("rank")

			// Check that at least one field is being updated
			if !titleSet && !descSet && !statusSet && !rankSet {
				return fmt.Errorf("at least one field must be specified to update (--title, --description, --status, or --rank)")
			}

			// Get flag values
			title, _ := cmd.Flags().GetString("title")
			description, _ := cmd.Flags().GetString("description")
			status, _ := cmd.Flags().GetString("status")
			rank, _ := cmd.Flags().GetInt("rank")

			// Create DTO with only updated fields
			input := dto.UpdateTaskDTO{
				ID: taskID,
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
			task, err := taskService.UpdateTask(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to update task: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Task updated successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:          %s\n", task.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Track:       %s\n", task.TrackID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:       %s\n", task.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status:      %s\n", task.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "  Rank:        %d\n", task.Rank)
			if task.Description != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", task.Description)
			}
			if task.Branch != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Branch:      %s\n", task.Branch)
			}

			return nil
		},
	}

	cmd.Flags().String("title", "", "New task title")
	cmd.Flags().String("description", "", "New task description")
	cmd.Flags().String("status", "", "New task status (todo, in-progress, review, done)")
	cmd.Flags().Int("rank", 0, "New task rank (1-1000)")

	return cmd
}

// ============================================================================
// task delete command
// ============================================================================

func newTaskDeleteCommand(taskService *application.TaskApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <task-id>",
		Short: "Delete a task",
		Long:  `Deletes a task and removes it from any iterations it belongs to.`,
		Example: `  # Delete a task
  tm task delete TM-task-1 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			taskID := args[0]

			// Execute via application service
			if err := taskService.DeleteTask(ctx, taskID); err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Task %s deleted successfully\n", taskID)

			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Force deletion without confirmation")

	return cmd
}

// ============================================================================
// task move command
// ============================================================================

func newTaskMoveCommand(taskService *application.TaskApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "move <task-id>",
		Short: "Move a task to a different track",
		Long:  `Moves a task from its current track to a different track.`,
		Example: `  # Move task to a different track
  tm task move TM-task-1 --track TM-track-2`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			taskID := args[0]

			// Get flag
			newTrackID, _ := cmd.Flags().GetString("track")

			// Validate required flag
			if newTrackID == "" {
				return fmt.Errorf("--track is required")
			}

			// Execute via application service
			if err := taskService.MoveTask(ctx, taskID, newTrackID); err != nil {
				return fmt.Errorf("failed to move task: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Task %s moved to track %s successfully\n", taskID, newTrackID)

			return nil
		},
	}

	cmd.Flags().String("track", "", "Target track ID (required)")
	cmd.MarkFlagRequired("track")

	return cmd
}

// ============================================================================
// task backlog command
// ============================================================================

func newTaskBacklogCommand(taskService *application.TaskApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backlog",
		Short: "List all tasks in backlog",
		Long:  `Lists all tasks with status "todo" (backlog items).`,
		Example: `  # Show all backlog tasks
  tm task backlog`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Execute via application service
			tasks, err := taskService.GetBacklogTasks(ctx)
			if err != nil {
				return fmt.Errorf("failed to get backlog tasks: %w", err)
			}

			// Format output
			if len(tasks) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No backlog tasks found\n")
				return nil
			}

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "Backlog Tasks\n")
			fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-20s %-40s\n", "ID", "Track", "Title")
			fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-20s %-40s\n",
				strings.Repeat("-", 15),
				strings.Repeat("-", 20),
				strings.Repeat("-", 40),
			)

			// Print tasks
			for _, task := range tasks {
				fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-20s %-40s\n",
					task.ID,
					task.TrackID,
					truncateString(task.Title, 40),
				)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d backlog task(s)\n", len(tasks))
			return nil
		},
	}

	return cmd
}

// ============================================================================
// task check-ready command
// ============================================================================

func newTaskCheckReadyCommand(taskService *application.TaskApplicationService, acService *application.ACApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-ready <task-id>",
		Short: "Check if all acceptance criteria for a task are verified",
		Long:  `Checks if all acceptance criteria for a task are verified and ready for completion.`,
		Example: `  # Check if task is ready
  tm task check-ready TM-task-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			taskID := args[0]

			// Get task to verify it exists
			task, err := taskService.GetTask(ctx, taskID)
			if err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			// Get all ACs for the task
			acs, err := acService.ListAC(ctx, taskID)
			if err != nil {
				return fmt.Errorf("failed to get acceptance criteria: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Task: %s\n", task.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "Task ID: %s\n", task.ID)

			if len(acs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\nNo acceptance criteria defined\n")
				fmt.Fprintf(cmd.OutOrStdout(), "Status: READY (no criteria to verify)\n")
				return nil
			}

			// Check verification status
			allVerified := true
			verifiedCount := 0

			fmt.Fprintf(cmd.OutOrStdout(), "\nAcceptance Criteria:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-50s %-15s\n", "AC ID", "Description", "Status")
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-50s %-15s\n",
				strings.Repeat("-", 20),
				strings.Repeat("-", 50),
				strings.Repeat("-", 15),
			)

			for _, ac := range acs {
				status := ac.Status
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-50s %-15s\n",
					ac.ID,
					truncateString(ac.Description, 50),
					status,
				)

				if status == entities.ACStatusVerified {
					verifiedCount++
				} else {
					allVerified = false
				}
			}

			// Summary
			fmt.Fprintf(cmd.OutOrStdout(), "\nSummary: %d/%d criteria verified\n", verifiedCount, len(acs))
			if allVerified {
				fmt.Fprintf(cmd.OutOrStdout(), "Status: READY (all criteria verified)\n")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Status: NOT READY (some criteria pending)\n")
			}

			return nil
		},
	}

	return cmd
}
