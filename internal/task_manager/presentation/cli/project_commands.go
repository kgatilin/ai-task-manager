package cli

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/spf13/cobra"
)

// Project name validation regex: alphanumeric + hyphens/underscores only
var projectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// NewProjectCommands creates and returns the project command group with all subcommands.
func NewProjectCommands(projectService *application.ProjectApplicationService) *cobra.Command {
	projectCmd := &cobra.Command{
		Use:     "project",
		Short:   "Manage projects",
		Long:    "Commands for managing multiple isolated project databases",
		Aliases: []string{"proj", "p"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all project subcommands
	projectCmd.AddCommand(
		newProjectCreateCommand(projectService),
		newProjectListCommand(projectService),
		newProjectShowCommand(projectService),
		newProjectSwitchCommand(projectService),
		newProjectDeleteCommand(projectService),
	)

	return projectCmd
}

// ============================================================================
// project create command
// ============================================================================

func newProjectCreateCommand(provider *application.ProjectApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <project-name>",
		Short: "Create a new project",
		Long: `Creates a new project with its own isolated database.

Each project maintains separate roadmaps, tracks, tasks, and iterations.
Project names must be alphanumeric with hyphens or underscores only.

The project name is used as the directory name under .tm/projects/ (or custom directory).
Optionally provide a project code for human-readable IDs (e.g., DW, PROD, TEST).`,
		Example: `  # Create a project with default code
  tm project create myproject

  # Create a project with custom code
  tm project create my-product --code PROD`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			// Validate project name
			if !projectNameRegex.MatchString(projectName) {
				return fmt.Errorf("invalid project name: must be alphanumeric with hyphens or underscores only")
			}

			// Get project code flag
			projectCode, _ := cmd.Flags().GetString("code")

			// If empty, application service will generate default
			if projectCode != "" && !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(projectCode) {
				return fmt.Errorf("invalid project code: must be alphanumeric uppercase (e.g., DW, PROD, TEST)")
			}

			// Create project via provider
			generatedCode, err := provider.CreateProject(projectName, projectCode)
			if err != nil {
				return fmt.Errorf("failed to create project: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Project created successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  Name: %s\n", projectName)
			fmt.Fprintf(cmd.OutOrStdout(), "  Project code: %s\n", generatedCode)
			fmt.Fprintf(cmd.OutOrStdout(), "  Location: .tm/projects/%s/\n", projectName)
			fmt.Fprintf(cmd.OutOrStdout(), "\nTo switch to this project, run:\n  tm project switch %s\n", projectName)

			return nil
		},
	}

	cmd.Flags().String("code", "", "Project code for human-readable IDs (optional, default: uppercase first 2-3 letters)")

	return cmd
}

// ============================================================================
// project list command
// ============================================================================

func newProjectListCommand(provider *application.ProjectApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Long: `Lists all projects in the current workspace.

Shows project names and indicates which project is currently active.`,
		Example: `  tm project list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get list of projects
			projects, err := provider.ListProjects()
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}

			if len(projects) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No projects found\n")
				return nil
			}

			// Get active project
			activeProject, _ := provider.GetActiveProject()

			fmt.Fprintf(cmd.OutOrStdout(), "Projects:\n\n")
			for _, proj := range projects {
				marker := "  "
				if proj == activeProject {
					marker = "* "
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", marker, proj)
			}

			if activeProject != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "\nActive: %s\n", activeProject)
			}

			return nil
		},
	}

	return cmd
}

// ============================================================================
// project show command
// ============================================================================

func newProjectShowCommand(provider *application.ProjectApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [project-name]",
		Short: "Show project information",
		Long: `Displays detailed information about a project.

If no project name is provided, shows the currently active project.`,
		Example: `  # Show active project
  tm project show

  # Show specific project
  tm project show myproject`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := ""
			if len(args) > 0 {
				projectName = args[0]
			}

			// If no project specified, use active project
			if projectName == "" {
				active, err := provider.GetActiveProject()
				if err != nil || active == "" {
					return fmt.Errorf("no active project set and no project specified")
				}
				projectName = active
			}

			// Get project info
			info, err := provider.GetProjectInfo(projectName)
			if err != nil {
				return fmt.Errorf("failed to get project info: %w", err)
			}

			// Get active project for comparison
			activeProject, _ := provider.GetActiveProject()

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Active project: %s\n\n", projectName)

			// Sort keys for consistent output
			keys := make([]string, 0, len(info))
			for k := range info {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, key := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s: %s\n", key, info[key])
			}

			if projectName == activeProject {
				fmt.Fprintf(cmd.OutOrStdout(), "\n(This is the active project)\n")
			}

			return nil
		},
	}

	return cmd
}

// ============================================================================
// project switch command
// ============================================================================

func newProjectSwitchCommand(provider *application.ProjectApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch <project-name>",
		Short: "Switch to a different project",
		Long: `Sets the active project for the current workspace.

All commands will use the active project by default.
You can override the active project on any command with the --project flag.`,
		Example: `  tm project switch myproject`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			// Switch project via provider
			if err := provider.SwitchProject(projectName); err != nil {
				return fmt.Errorf("failed to switch project: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Switched to project: %s\n", projectName)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// project delete command
// ============================================================================

func newProjectDeleteCommand(provider *application.ProjectApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <project-name>",
		Short: "Delete a project",
		Long: `Deletes a project and its entire database.

Warning: This operation is irreversible. All data for the project will be lost.`,
		Example: `  # Delete a project (with confirmation prompt)
  tm project delete myproject

  # Force delete without confirmation
  tm project delete myproject --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			// Get force flag
			force, _ := cmd.Flags().GetBool("force")

			// Confirm deletion unless forced
			if !force {
				fmt.Fprintf(cmd.OutOrStdout(), "Are you sure you want to delete project '%s'? This cannot be undone.\n", projectName)
				fmt.Fprintf(cmd.OutOrStdout(), "Type the project name to confirm: ")

				var confirmation string
				_, err := fmt.Scanln(&confirmation)
				if err != nil || confirmation != projectName {
					fmt.Fprintf(cmd.OutOrStdout(), "Deletion cancelled\n")
					return nil
				}
			}

			// Delete project via provider
			if err := provider.DeleteProject(projectName); err != nil {
				return fmt.Errorf("failed to delete project: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Project deleted successfully: %s\n", projectName)

			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Force delete without confirmation prompt")

	return cmd
}
