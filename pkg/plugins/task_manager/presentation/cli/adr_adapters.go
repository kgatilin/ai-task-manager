package cli

import (
	"context"
	"fmt"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/dto"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ============================================================================
// ADRCreateCommandAdapter - Adapts CLI to CreateADRCommand use case
// ============================================================================

type ADRCreateCommandAdapter struct {
	ADRService   *application.ADRApplicationService

	// CLI flags
	project       string
	trackID       string
	title         string
	context       string
	decision      string
	consequences  string
	alternatives  string
}

func (c *ADRCreateCommandAdapter) GetName() string {
	return "adr create"
}

func (c *ADRCreateCommandAdapter) GetDescription() string {
	return "Create an Architecture Decision Record for a track"
}

func (c *ADRCreateCommandAdapter) GetUsage() string {
	return "dw task-manager adr create <track-id> --title <title> --context <ctx> --decision <dec> --consequences <cons> [--alternatives <alt>]"
}

func (c *ADRCreateCommandAdapter) GetHelp() string {
	return `Creates an Architecture Decision Record (ADR) for a track.

Flags:
  --title <title>              ADR title (required)
  --context <context>          Problem context (required)
  --decision <decision>        Decision made (required)
  --consequences <cons>        Decision consequences (required)
  --alternatives <alt>         Alternative approaches considered (optional)
  --project <name>             Project name (optional)`
}

func (c *ADRCreateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
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
		case "--title":
			if i+1 < len(args) {
				c.title = args[i+1]
				i++
			}
		case "--context":
			if i+1 < len(args) {
				c.context = args[i+1]
				i++
			}
		case "--decision":
			if i+1 < len(args) {
				c.decision = args[i+1]
				i++
			}
		case "--consequences":
			if i+1 < len(args) {
				c.consequences = args[i+1]
				i++
			}
		case "--alternatives":
			if i+1 < len(args) {
				c.alternatives = args[i+1]
				i++
			}
		}
	}

	// Validate required flags
	if c.title == "" {
		return fmt.Errorf("--title is required")
	}
	if c.context == "" {
		return fmt.Errorf("--context is required")
	}
	if c.decision == "" {
		return fmt.Errorf("--decision is required")
	}
	if c.consequences == "" {
		return fmt.Errorf("--consequences is required")
	}


	// Create DTO
	input := dto.CreateADRDTO{
		TrackID:      c.trackID,
		Title:        c.title,
		Context:      c.context,
		Decision:     c.decision,
		Consequences: c.consequences,
		Alternatives: c.alternatives,
	}

	// Execute via application service
	adr, err := c.ADRService.CreateADR(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create ADR: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "ADR created successfully\n")
	fmt.Fprintf(out, "  ID:           %s\n", adr.ID)
	fmt.Fprintf(out, "  Track:        %s\n", adr.TrackID)
	fmt.Fprintf(out, "  Title:        %s\n", adr.Title)
	fmt.Fprintf(out, "  Status:       %s\n", adr.Status)

	return nil
}

// ============================================================================
// ADRUpdateCommandAdapter - Adapts CLI to UpdateADRCommand use case
// ============================================================================

type ADRUpdateCommandAdapter struct {
	ADRService   *application.ADRApplicationService

	// CLI flags
	project      string
	adrID        string
	title        *string
	context      *string
	decision     *string
	consequences *string
	alternatives *string
	status       *string
}

func (c *ADRUpdateCommandAdapter) GetName() string {
	return "adr update"
}

func (c *ADRUpdateCommandAdapter) GetDescription() string {
	return "Update an existing ADR"
}

func (c *ADRUpdateCommandAdapter) GetUsage() string {
	return "dw task-manager adr update <adr-id> [options]"
}

func (c *ADRUpdateCommandAdapter) GetHelp() string {
	return `Updates an existing ADR's fields.

Flags:
  --title <title>              New ADR title
  --context <context>          New problem context
  --decision <decision>        New decision
  --consequences <cons>        New consequences
  --alternatives <alt>         New alternatives
  --status <status>            New status (proposed, accepted, superseded, deprecated)
  --project <name>             Project name (optional)`
}

func (c *ADRUpdateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse ADR ID
	if len(args) == 0 {
		return fmt.Errorf("ADR ID is required")
	}
	c.adrID = args[0]
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
		case "--context":
			if i+1 < len(args) {
				val := args[i+1]
				c.context = &val
				i++
			}
		case "--decision":
			if i+1 < len(args) {
				val := args[i+1]
				c.decision = &val
				i++
			}
		case "--consequences":
			if i+1 < len(args) {
				val := args[i+1]
				c.consequences = &val
				i++
			}
		case "--alternatives":
			if i+1 < len(args) {
				val := args[i+1]
				c.alternatives = &val
				i++
			}
		case "--status":
			if i+1 < len(args) {
				val := args[i+1]
				c.status = &val
				i++
			}
		}
	}

	// Validate at least one field
	if c.title == nil && c.context == nil && c.decision == nil && c.consequences == nil && c.alternatives == nil && c.status == nil {
		return fmt.Errorf("at least one field must be specified to update")
	}

	// Create DTO
	input := dto.UpdateADRDTO{
		ID:           c.adrID,
		Title:        c.title,
		Context:      c.context,
		Decision:     c.decision,
		Consequences: c.consequences,
		Alternatives: c.alternatives,
		Status:       c.status,
	}

	// Execute via application service
	adr, err := c.ADRService.UpdateADR(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update ADR: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "ADR updated successfully\n")
	fmt.Fprintf(out, "  ID:           %s\n", adr.ID)
	fmt.Fprintf(out, "  Track:        %s\n", adr.TrackID)
	fmt.Fprintf(out, "  Title:        %s\n", adr.Title)
	fmt.Fprintf(out, "  Status:       %s\n", adr.Status)

	return nil
}

// ============================================================================
// ADRListCommandAdapter - Lists ADRs
// ============================================================================

type ADRListCommandAdapter struct {
	ADRService   *application.ADRApplicationService

	// CLI flags
	project string
	trackID string
}

func (c *ADRListCommandAdapter) GetName() string {
	return "adr list"
}

func (c *ADRListCommandAdapter) GetDescription() string {
	return "List Architecture Decision Records"
}

func (c *ADRListCommandAdapter) GetUsage() string {
	return "dw task-manager adr list [--track <track-id>]"
}

func (c *ADRListCommandAdapter) GetHelp() string {
	return `Lists all ADRs or ADRs for a specific track.

Flags:
  --track <track-id>    Filter by track ID (optional)
  --project <name>      Project name (optional)

Examples:
  # List all ADRs
  dw task-manager adr list

  # List ADRs for a track
  dw task-manager adr list --track TM-track-1`
}

func (c *ADRListCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--track":
			if i+1 < len(args) {
				c.trackID = args[i+1]
				i++
			}
		}
	}

	// List ADRs via application service
	var trackIDPtr *string
	if c.trackID != "" {
		trackIDPtr = &c.trackID
	}
	adrs, err := c.ADRService.ListADRs(ctx, trackIDPtr)
	if err != nil {
		return fmt.Errorf("failed to list ADRs: %w", err)
	}

	// Display ADRs
	out := cmdCtx.GetStdout()
	if len(adrs) == 0 {
		fmt.Fprintf(out, "No ADRs found.\n")
		return nil
	}

	// Print header
	fmt.Fprintf(out, "%-20s %-25s %-40s %-12s\n", "ID", "Track", "Title", "Status")
	fmt.Fprintf(out, "%s\n", "------------------------------------------------------------------------------")

	// Print each ADR
	for _, adr := range adrs {
		title := adr.Title
		if len(title) > 39 {
			title = title[:36] + "..."
		}
		fmt.Fprintf(out, "%-20s %-25s %-40s %-12s\n", adr.ID, adr.TrackID, title, adr.Status)
	}

	fmt.Fprintf(out, "\nTotal: %d ADR(s)\n", len(adrs))
	return nil
}

// ============================================================================
// ADRShowCommandAdapter - Shows detailed ADR information
// ============================================================================

type ADRShowCommandAdapter struct {
	ADRService   *application.ADRApplicationService

	// CLI flags
	project string
	adrID   string
}

func (c *ADRShowCommandAdapter) GetName() string {
	return "adr show"
}

func (c *ADRShowCommandAdapter) GetDescription() string {
	return "Show detailed information about an ADR"
}

func (c *ADRShowCommandAdapter) GetUsage() string {
	return "dw task-manager adr show <adr-id>"
}

func (c *ADRShowCommandAdapter) GetHelp() string {
	return `Shows detailed information about an Architecture Decision Record.

Flags:
  --project <name>    Project name (optional)

Examples:
  # Show ADR details
  dw task-manager adr show TM-adr-1`
}

func (c *ADRShowCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse ADR ID
	if len(args) == 0 {
		return fmt.Errorf("ADR ID is required")
	}
	c.adrID = args[0]
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

	// Get ADR via application service
	adr, err := c.ADRService.GetADR(ctx, c.adrID)
	if err != nil {
		return fmt.Errorf("failed to get ADR: %w", err)
	}

	// Display ADR details
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Architecture Decision Record\n")
	fmt.Fprintf(out, "============================\n\n")
	fmt.Fprintf(out, "ID:           %s\n", adr.ID)
	fmt.Fprintf(out, "Track:        %s\n", adr.TrackID)
	fmt.Fprintf(out, "Title:        %s\n", adr.Title)
	fmt.Fprintf(out, "Status:       %s\n", adr.Status)
	if adr.SupersededBy != nil && *adr.SupersededBy != "" {
		fmt.Fprintf(out, "Superseded By: %s\n", *adr.SupersededBy)
	}

	fmt.Fprintf(out, "\nContext:\n")
	fmt.Fprintf(out, "--------\n")
	fmt.Fprintf(out, "%s\n", adr.Context)

	fmt.Fprintf(out, "\nDecision:\n")
	fmt.Fprintf(out, "---------\n")
	fmt.Fprintf(out, "%s\n", adr.Decision)

	fmt.Fprintf(out, "\nConsequences:\n")
	fmt.Fprintf(out, "-------------\n")
	fmt.Fprintf(out, "%s\n", adr.Consequences)

	if adr.Alternatives != "" {
		fmt.Fprintf(out, "\nAlternatives Considered:\n")
		fmt.Fprintf(out, "------------------------\n")
		fmt.Fprintf(out, "%s\n", adr.Alternatives)
	}

	fmt.Fprintf(out, "\nTimestamps:\n")
	fmt.Fprintf(out, "-----------\n")
	fmt.Fprintf(out, "Created: %s\n", adr.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(out, "Updated: %s\n", adr.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// ============================================================================
// ADRSupersedeCommandAdapter - Supersedes an ADR with another
// ============================================================================

type ADRSupersedeCommandAdapter struct {
	ADRService    *application.ADRApplicationService

	// CLI flags
	project       string
	adrID         string
	supersededByID string
}

func (c *ADRSupersedeCommandAdapter) GetName() string {
	return "adr supersede"
}

func (c *ADRSupersedeCommandAdapter) GetDescription() string {
	return "Mark an ADR as superseded by another ADR"
}

func (c *ADRSupersedeCommandAdapter) GetUsage() string {
	return "dw task-manager adr supersede <adr-id> --superseded-by <new-adr-id>"
}

func (c *ADRSupersedeCommandAdapter) GetHelp() string {
	return `Marks an ADR as superseded by a newer ADR.

This is used when a new architectural decision replaces an old one.

Flags:
  --superseded-by <new-adr-id>    ID of the superseding ADR (required)
  --project <name>                Project name (optional)

Examples:
  # Mark TM-adr-1 as superseded by TM-adr-5
  dw task-manager adr supersede TM-adr-1 --superseded-by TM-adr-5`
}

func (c *ADRSupersedeCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse ADR ID
	if len(args) == 0 {
		return fmt.Errorf("ADR ID is required")
	}
	c.adrID = args[0]
	args = args[1:]

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--superseded-by":
			if i+1 < len(args) {
				c.supersededByID = args[i+1]
				i++
			}
		}
	}

	// Validate required flag
	if c.supersededByID == "" {
		return fmt.Errorf("--superseded-by is required")
	}

	// Execute via application service
	if err := c.ADRService.SupersedeADR(ctx, c.adrID, c.supersededByID); err != nil {
		return fmt.Errorf("failed to supersede ADR: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "ADR superseded successfully\n")
	fmt.Fprintf(out, "  %s is now superseded by %s\n", c.adrID, c.supersededByID)

	return nil
}

// ============================================================================
// ADRDeprecateCommandAdapter - Deprecates an ADR
// ============================================================================

type ADRDeprecateCommandAdapter struct {
	ADRService   *application.ADRApplicationService

	// CLI flags
	project string
	adrID   string
}

func (c *ADRDeprecateCommandAdapter) GetName() string {
	return "adr deprecate"
}

func (c *ADRDeprecateCommandAdapter) GetDescription() string {
	return "Mark an ADR as deprecated"
}

func (c *ADRDeprecateCommandAdapter) GetUsage() string {
	return "dw task-manager adr deprecate <adr-id>"
}

func (c *ADRDeprecateCommandAdapter) GetHelp() string {
	return `Marks an ADR as deprecated.

Use this when an ADR is no longer relevant but isn't directly superseded
by another ADR.

Flags:
  --project <name>    Project name (optional)

Examples:
  # Deprecate an ADR
  dw task-manager adr deprecate TM-adr-3`
}

func (c *ADRDeprecateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse ADR ID
	if len(args) == 0 {
		return fmt.Errorf("ADR ID is required")
	}
	c.adrID = args[0]
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

	// Execute via application service
	if err := c.ADRService.DeprecateADR(ctx, c.adrID); err != nil {
		return fmt.Errorf("failed to deprecate ADR: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "ADR deprecated successfully\n")
	fmt.Fprintf(out, "  ID: %s\n", c.adrID)

	return nil
}
