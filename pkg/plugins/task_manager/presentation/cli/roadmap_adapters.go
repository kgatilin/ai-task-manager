package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/dto"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ============================================================================
// RoadmapInitCommandAdapter - Adapts CLI to InitRoadmap use case
// ============================================================================

// RoadmapInitCommandAdapter adapts roadmap init CLI command to application use case
type RoadmapInitCommandAdapter struct {
	RoadmapService *application.RoadmapApplicationService

	// CLI flags (parsed from args)
	project         string
	vision          string
	successCriteria string
}

func (c *RoadmapInitCommandAdapter) GetName() string {
	return "roadmap init"
}

func (c *RoadmapInitCommandAdapter) GetDescription() string {
	return "Initialize a new roadmap"
}

func (c *RoadmapInitCommandAdapter) GetUsage() string {
	return "dw task-manager roadmap init --vision <vision> --success-criteria <criteria>"
}

func (c *RoadmapInitCommandAdapter) GetHelp() string {
	return `Creates a new roadmap with a vision statement and success criteria.

Only one active roadmap can exist at a time. If you need to replace the current
roadmap, delete it first using 'dw task-manager roadmap delete'.

Flags:
  --vision <vision>              The vision statement for the roadmap (required)
  --success-criteria <criteria>  Success criteria for the roadmap (required)
  --project <name>               Project name (optional, uses active project if not specified)

Examples:
  # Create a simple roadmap
  dw task-manager roadmap init \
    --vision "Build extensible framework" \
    --success-criteria "Support 10 plugins"

  # With multi-line vision
  dw task-manager roadmap init \
    --vision "Create unified productivity platform" \
    --success-criteria "100% test coverage, zero violations"

Notes:
  - Vision must be non-empty
  - Success criteria must be non-empty
  - Only one roadmap can be active at a time`
}

func (c *RoadmapInitCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--vision":
			if i+1 < len(args) {
				c.vision = args[i+1]
				i++
			}
		case "--success-criteria":
			if i+1 < len(args) {
				c.successCriteria = args[i+1]
				i++
			}
		}
	}

	// Validate required flags
	if c.vision == "" {
		return fmt.Errorf("--vision is required")
	}
	if c.successCriteria == "" {
		return fmt.Errorf("--success-criteria is required")
	}

	// Create DTO
	input := dto.CreateRoadmapDTO{
		Vision:          c.vision,
		SuccessCriteria: c.successCriteria,
	}

	// Execute via application service
	roadmap, err := c.RoadmapService.InitRoadmap(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create roadmap: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Roadmap created successfully\n")
	fmt.Fprintf(out, "ID:                %s\n", roadmap.ID)
	fmt.Fprintf(out, "Vision:            %s\n", roadmap.Vision)
	fmt.Fprintf(out, "Success Criteria:  %s\n", roadmap.SuccessCriteria)

	return nil
}

// ============================================================================
// RoadmapShowCommandAdapter - Adapts CLI to GetRoadmap use case
// ============================================================================

// RoadmapShowCommandAdapter adapts roadmap show CLI command to application use case
type RoadmapShowCommandAdapter struct {
	RoadmapService *application.RoadmapApplicationService

	// CLI flags (parsed from args)
	project string
}

func (c *RoadmapShowCommandAdapter) GetName() string {
	return "roadmap show"
}

func (c *RoadmapShowCommandAdapter) GetDescription() string {
	return "Display the current roadmap"
}

func (c *RoadmapShowCommandAdapter) GetUsage() string {
	return "dw task-manager roadmap show"
}

func (c *RoadmapShowCommandAdapter) GetHelp() string {
	return `Displays the details of the current active roadmap.

If no roadmap exists, you can create one using:
  dw task-manager roadmap init --vision <vision> --success-criteria <criteria>

Flags:
  --project <name>    Project name (optional, uses active project if not specified)

Examples:
  dw task-manager roadmap show

Output:
  ID:                roadmap-1234567890
  Vision:            Build extensible framework
  Success Criteria:  Support 10 plugins
  Created:           2025-10-31T10:00:00Z
  Updated:           2025-10-31T10:00:00Z`
}

func (c *RoadmapShowCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
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

	// Get active roadmap via application service
	roadmap, err := c.RoadmapService.GetRoadmap(ctx)
	if err != nil {
		if err == pluginsdk.ErrNotFound {
			out := cmdCtx.GetStdout()
			fmt.Fprintf(out, "No roadmap found.\n")
			fmt.Fprintf(out, "Run 'dw task-manager roadmap init' to create one.\n")
			return nil
		}
		return fmt.Errorf("failed to get roadmap: %w", err)
	}

	// Display roadmap details
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Roadmap:\n")
	fmt.Fprintf(out, "  ID:                %s\n", roadmap.ID)
	fmt.Fprintf(out, "  Vision:            %s\n", roadmap.Vision)
	fmt.Fprintf(out, "  Success Criteria:  %s\n", roadmap.SuccessCriteria)
	fmt.Fprintf(out, "  Created:           %s\n", roadmap.CreatedAt.Format(time.RFC3339))
	fmt.Fprintf(out, "  Updated:           %s\n", roadmap.UpdatedAt.Format(time.RFC3339))

	return nil
}

// ============================================================================
// RoadmapUpdateCommandAdapter - Adapts CLI to UpdateRoadmap use case
// ============================================================================

// RoadmapUpdateCommandAdapter adapts roadmap update CLI command to application use case
type RoadmapUpdateCommandAdapter struct {
	RoadmapService *application.RoadmapApplicationService

	// CLI flags (parsed from args)
	project         string
	vision          *string
	successCriteria *string
}

func (c *RoadmapUpdateCommandAdapter) GetName() string {
	return "roadmap update"
}

func (c *RoadmapUpdateCommandAdapter) GetDescription() string {
	return "Update the current roadmap"
}

func (c *RoadmapUpdateCommandAdapter) GetUsage() string {
	return "dw task-manager roadmap update [--vision <vision>] [--success-criteria <criteria>]"
}

func (c *RoadmapUpdateCommandAdapter) GetHelp() string {
	return `Updates properties of the current active roadmap.

At least one flag must be provided to update.

Flags:
  --vision <vision>              New vision statement
  --success-criteria <criteria>  New success criteria
  --project <name>               Project name (optional, uses active project if not specified)

Examples:
  # Update vision
  dw task-manager roadmap update --vision "Create unified platform"

  # Update both
  dw task-manager roadmap update \
    --vision "New vision" \
    --success-criteria "New criteria"

Notes:
  - At least one flag is required
  - Run 'dw task-manager roadmap show' to see current values
  - Updated_at timestamp is automatically updated`
}

func (c *RoadmapUpdateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--vision":
			if i+1 < len(args) {
				v := args[i+1]
				c.vision = &v
				i++
			}
		case "--success-criteria":
			if i+1 < len(args) {
				sc := args[i+1]
				c.successCriteria = &sc
				i++
			}
		}
	}

	// Validate at least one field is provided
	if c.vision == nil && c.successCriteria == nil {
		return fmt.Errorf("at least one flag must be provided (--vision or --success-criteria)")
	}

	// Create DTO
	input := dto.UpdateRoadmapDTO{
		Vision:          c.vision,
		SuccessCriteria: c.successCriteria,
	}

	// Execute via application service
	roadmap, err := c.RoadmapService.UpdateRoadmap(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update roadmap: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Roadmap updated successfully\n")
	fmt.Fprintf(out, "ID:                %s\n", roadmap.ID)
	fmt.Fprintf(out, "Vision:            %s\n", roadmap.Vision)
	fmt.Fprintf(out, "Success Criteria:  %s\n", roadmap.SuccessCriteria)
	fmt.Fprintf(out, "Updated:           %s\n", roadmap.UpdatedAt.Format(time.RFC3339))

	return nil
}

// ============================================================================
// RoadmapFullCommandAdapter - Adapts CLI to GetFullOverview use case
// ============================================================================

// RoadmapFullCommandAdapter adapts roadmap full CLI command to application use case
type RoadmapFullCommandAdapter struct {
	RoadmapService *application.RoadmapApplicationService

	// CLI flags (parsed from args)
	project  string
	verbose  bool
	format   string
	sections []string
}

func (c *RoadmapFullCommandAdapter) GetName() string {
	return "roadmap full"
}

func (c *RoadmapFullCommandAdapter) GetDescription() string {
	return "Display complete roadmap overview in LLM-optimized format"
}

func (c *RoadmapFullCommandAdapter) GetUsage() string {
	return "dw task-manager roadmap full [--verbose] [--format json] [--sections <list>]"
}

func (c *RoadmapFullCommandAdapter) GetHelp() string {
	return `Displays the complete roadmap overview in LLM-optimized markdown format.

Shows:
  - Vision and success criteria
  - All tracks with their tasks (titles only by default)
  - All iterations with assigned tasks and progress
  - Backlog (tasks not in any iteration)

Flags:
  --verbose             Include task descriptions and additional details
  --format <format>     Output format (default: markdown)
                        Values: markdown, json
  --sections <list>     Only show specific sections (comma-separated)
                        Values: vision, tracks, iterations, backlog
  --project <name>      Project name (optional, uses active project if not specified)

Examples:
  # Basic overview
  dw task-manager roadmap full

  # Verbose with full details
  dw task-manager roadmap full --verbose

  # Only show tracks and iterations
  dw task-manager roadmap full --sections tracks,iterations

  # Output as JSON
  dw task-manager roadmap full --format json

Notes:
  - Uses status icons: ✅ (complete), ⏸️ (planned/in-progress), ○ (todo)
  - Optimized for LLM consumption
  - JSON format includes all metadata`
}

func (c *RoadmapFullCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	c.format = "markdown"
	c.sections = nil

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--verbose":
			c.verbose = true
		case "--format":
			if i+1 < len(args) {
				c.format = args[i+1]
				i++
			}
		case "--sections":
			if i+1 < len(args) {
				c.sections = splitSections(args[i+1])
				i++
			}
		}
	}

	// Build options
	options := dto.RoadmapOverviewOptions{
		Verbose:  c.verbose,
		Sections: c.sections,
	}

	// Get full overview via application service
	overview, err := c.RoadmapService.GetFullOverview(ctx, options)
	if err != nil {
		if err == pluginsdk.ErrNotFound {
			out := cmdCtx.GetStdout()
			fmt.Fprintf(out, "No roadmap found.\n")
			fmt.Fprintf(out, "Run 'dw task-manager roadmap init' to create one.\n")
			return nil
		}
		return fmt.Errorf("failed to get roadmap overview: %w", err)
	}

	// Format output based on requested format
	if c.format == "json" {
		return c.outputJSON(cmdCtx, overview)
	}

	return c.outputMarkdown(cmdCtx, overview, options)
}

func (c *RoadmapFullCommandAdapter) outputJSON(cmdCtx pluginsdk.CommandContext, overview *dto.RoadmapOverviewDTO) error {
	out := cmdCtx.GetStdout()

	// Convert to JSON
	data, err := json.MarshalIndent(overview, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Fprintf(out, "%s\n", string(data))
	return nil
}

func (c *RoadmapFullCommandAdapter) outputMarkdown(cmdCtx pluginsdk.CommandContext, overview *dto.RoadmapOverviewDTO, options dto.RoadmapOverviewOptions) error {
	out := cmdCtx.GetStdout()

	// Extract entities from overview
	roadmap := overview.Roadmap.(*entities.RoadmapEntity)
	tracks := c.convertToTracks(overview.Tracks)
	tasks := c.convertToTasks(overview.Tasks)
	iterations := c.convertToIterations(overview.Iterations)

	// Vision section
	if options.ShouldShowSection("vision") {
		fmt.Fprintf(out, "# Roadmap: %s\n\n", roadmap.ID)
		fmt.Fprintf(out, "## Vision\n%s\n\n", roadmap.Vision)
		fmt.Fprintf(out, "## Success Criteria\n%s\n\n", roadmap.SuccessCriteria)
	}

	// Tracks section
	if options.ShouldShowSection("tracks") {
		fmt.Fprintf(out, "## Tracks\n\n")
		for _, track := range tracks {
			statusIcon := GetStatusIcon(track.Status)
			fmt.Fprintf(out, "### %s %s\n", statusIcon, track.Title)
			fmt.Fprintf(out, "**ID**: %s | **Status**: %s | **Rank**: %d\n", track.ID, track.Status, track.Rank)

			if c.verbose && track.Description != "" {
				fmt.Fprintf(out, "**Description**: %s\n", track.Description)
			}

			// Get tasks for this track
			trackTasks := filterTasksByTrack(tasks, track.ID)
			if len(trackTasks) > 0 {
				fmt.Fprintf(out, "**Progress**: %d/%d tasks complete\n", countCompleteTasks(trackTasks), len(trackTasks))
				fmt.Fprintf(out, "**Tasks**:\n")
				for _, task := range trackTasks {
					taskIcon := GetStatusIcon(task.Status)
					fmt.Fprintf(out, "- %s %s", taskIcon, task.Title)
					if c.verbose && task.Description != "" {
						fmt.Fprintf(out, " - %s", task.Description)
					}
					fmt.Fprintf(out, "\n")
				}
			}
			fmt.Fprintf(out, "\n")
		}
	}

	// Iterations section
	if options.ShouldShowSection("iterations") {
		fmt.Fprintf(out, "## Iterations\n\n")
		for _, iter := range iterations {
			statusIcon := GetStatusIcon(iter.Status)
			fmt.Fprintf(out, "### %s Iteration %d: %s\n", statusIcon, iter.Number, iter.Name)
			fmt.Fprintf(out, "**Status**: %s | **Goal**: %s\n", iter.Status, iter.Goal)

			if c.verbose && iter.Deliverable != "" {
				fmt.Fprintf(out, "**Deliverable**: %s\n", iter.Deliverable)
			}

			// Get tasks for this iteration
			iterTasks := filterTasksByIDs(tasks, iter.TaskIDs)
			if len(iterTasks) > 0 {
				fmt.Fprintf(out, "**Progress**: %d/%d tasks complete (%.0f%%)\n",
					countCompleteTasks(iterTasks), len(iterTasks),
					float64(countCompleteTasks(iterTasks))/float64(len(iterTasks))*100)
				fmt.Fprintf(out, "**Tasks**:\n")
				for _, task := range iterTasks {
					taskIcon := GetStatusIcon(task.Status)
					fmt.Fprintf(out, "- %s %s", taskIcon, task.Title)
					if c.verbose && task.Description != "" {
						fmt.Fprintf(out, " - %s", task.Description)
					}
					fmt.Fprintf(out, "\n")
				}
			}
			fmt.Fprintf(out, "\n")
		}
	}

	// Backlog section
	if options.ShouldShowSection("backlog") {
		backlogTasks := getBacklogTasks(tasks, iterations)
		if len(backlogTasks) > 0 {
			fmt.Fprintf(out, "## Backlog\n\n")
			fmt.Fprintf(out, "%d tasks not assigned to any iteration:\n\n", len(backlogTasks))
			for _, task := range backlogTasks {
				taskIcon := GetStatusIcon(task.Status)
				fmt.Fprintf(out, "- %s %s", taskIcon, task.Title)
				if c.verbose && task.Description != "" {
					fmt.Fprintf(out, " - %s", task.Description)
				}
				fmt.Fprintf(out, "\n")
			}
		}
	}

	return nil
}

// Helper methods for type conversion

func (c *RoadmapFullCommandAdapter) convertToTracks(tracks []interface{}) []*entities.TrackEntity {
	result := make([]*entities.TrackEntity, len(tracks))
	for i, t := range tracks {
		result[i] = t.(*entities.TrackEntity)
	}
	return result
}

func (c *RoadmapFullCommandAdapter) convertToTasks(tasks []interface{}) []*entities.TaskEntity {
	result := make([]*entities.TaskEntity, len(tasks))
	for i, t := range tasks {
		result[i] = t.(*entities.TaskEntity)
	}
	return result
}

func (c *RoadmapFullCommandAdapter) convertToIterations(iterations []interface{}) []*entities.IterationEntity {
	result := make([]*entities.IterationEntity, len(iterations))
	for i, it := range iterations {
		result[i] = it.(*entities.IterationEntity)
	}
	return result
}

// ============================================================================
// Helper Functions
// ============================================================================

func filterTasksByTrack(tasks []*entities.TaskEntity, trackID string) []*entities.TaskEntity {
	var result []*entities.TaskEntity
	for _, task := range tasks {
		if task.TrackID == trackID {
			result = append(result, task)
		}
	}
	return result
}

func filterTasksByIDs(tasks []*entities.TaskEntity, taskIDs []string) []*entities.TaskEntity {
	var result []*entities.TaskEntity
	taskMap := make(map[string]*entities.TaskEntity)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}
	for _, id := range taskIDs {
		if task, ok := taskMap[id]; ok {
			result = append(result, task)
		}
	}
	return result
}

func countCompleteTasks(tasks []*entities.TaskEntity) int {
	count := 0
	for _, task := range tasks {
		if task.Status == "done" {
			count++
		}
	}
	return count
}

func getBacklogTasks(tasks []*entities.TaskEntity, iterations []*entities.IterationEntity) []*entities.TaskEntity {
	// Build set of all task IDs in iterations
	inIteration := make(map[string]bool)
	for _, iter := range iterations {
		for _, taskID := range iter.TaskIDs {
			inIteration[taskID] = true
		}
	}

	// Return tasks not in any iteration
	var backlog []*entities.TaskEntity
	for _, task := range tasks {
		if !inIteration[task.ID] {
			backlog = append(backlog, task)
		}
	}
	return backlog
}

func splitSections(sectionStr string) []string {
	if sectionStr == "" {
		return nil
	}
	parts := strings.Split(sectionStr, ",")
	var sections []string
	for _, s := range parts {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			sections = append(sections, trimmed)
		}
	}
	return sections
}
