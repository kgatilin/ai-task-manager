package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/application/dto"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ============================================================================
// DocCreateCommandAdapter - Adapts CLI to CreateDocument use case
// ============================================================================

type DocCreateCommandAdapter struct {
	DocumentService *application.DocumentApplicationService

	// CLI flags
	project          string
	title            string
	docType          string
	status           string
	content          string
	fromFile         string
	track            *string
	iteration        *int
}

func (c *DocCreateCommandAdapter) GetName() string {
	return "doc create"
}

func (c *DocCreateCommandAdapter) GetDescription() string {
	return "Create a new document"
}

func (c *DocCreateCommandAdapter) GetUsage() string {
	return "dw task-manager doc create --title <title> --type <type> [--status <status>] [--content <content> | --from-file <path>] [--track <id> | --iteration <num>]"
}

func (c *DocCreateCommandAdapter) GetHelp() string {
	return `Creates a new document with markdown content.

Documents can be of various types (adr, plan, retrospective, other) and
can optionally be attached to a track or iteration. Content must be provided
either inline with --content or from a markdown file with --from-file.

Flags:
  --title <title>           Document title (required)
  --type <type>             Document type (required): adr, plan, retrospective, other
  --status <status>         Document status (optional, default: draft): draft, published, archived
  --content <content>       Document markdown content (required if --from-file not used)
  --from-file <path>        Read markdown content from file (required if --content not used)
  --track <id>              Attach to track ID (mutually exclusive with --iteration)
  --iteration <num>         Attach to iteration number (mutually exclusive with --track)
  --project <name>          Project name (optional, uses active project if not specified)

Notes:
  - Either --content or --from-file must be provided (not both)
  - Either --track or --iteration can be provided (not both)
  - Document can be created unattached if neither --track nor --iteration is specified

Examples:
  # Create document from file attached to track
  dw task-manager doc create \
    --title "API Design" \
    --type plan \
    --from-file ./docs/api-design.md \
    --track TM-track-1

  # Create inline document
  dw task-manager doc create \
    --title "Sprint 1 Retrospective" \
    --type retrospective \
    --status published \
    --content "## What went well\n- Good team communication"

  # Create attached to iteration
  dw task-manager doc create \
    --title "Architecture Decision" \
    --type adr \
    --from-file ./adr-001.md \
    --iteration 5`
}

func (c *DocCreateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse flags
	c.status = "draft" // default status
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
		case "--type":
			if i+1 < len(args) {
				c.docType = args[i+1]
				i++
			}
		case "--status":
			if i+1 < len(args) {
				c.status = args[i+1]
				i++
			}
		case "--content":
			if i+1 < len(args) {
				c.content = args[i+1]
				i++
			}
		case "--from-file":
			if i+1 < len(args) {
				c.fromFile = args[i+1]
				i++
			}
		case "--track":
			if i+1 < len(args) {
				trackID := args[i+1]
				c.track = &trackID
				i++
			}
		case "--iteration":
			if i+1 < len(args) {
				iterNum, err := strconv.Atoi(args[i+1])
				if err != nil {
					return fmt.Errorf("invalid iteration number: %v", err)
				}
				c.iteration = &iterNum
				i++
			}
		}
	}

	// Validate required flags
	if c.title == "" {
		return fmt.Errorf("--title is required")
	}
	if c.docType == "" {
		return fmt.Errorf("--type is required")
	}

	// Validate XOR: content vs from-file
	if c.content != "" && c.fromFile != "" {
		return fmt.Errorf("--content and --from-file are mutually exclusive (provide one, not both)")
	}
	if c.content == "" && c.fromFile == "" {
		return fmt.Errorf("either --content or --from-file is required")
	}

	// Read file if provided
	if c.fromFile != "" {
		data, err := os.ReadFile(c.fromFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", c.fromFile, err)
		}
		c.content = string(data)
	}

	// Create DTO
	input := dto.CreateDocumentDTO{
		Title:           c.title,
		Type:            c.docType,
		Status:          c.status,
		Content:         c.content,
		TrackID:         c.track,
		IterationNumber: c.iteration,
	}

	// Execute via application service
	docID, err := c.DocumentService.CreateDocument(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Document created: %s\n", docID)

	return nil
}

// ============================================================================
// DocUpdateCommandAdapter - Adapts CLI to UpdateDocument use case
// ============================================================================

type DocUpdateCommandAdapter struct {
	DocumentService *application.DocumentApplicationService

	// CLI flags
	project   string
	docID     string
	content   *string
	fromFile  string
	status    *string
	detach    bool
}

func (c *DocUpdateCommandAdapter) GetName() string {
	return "doc update"
}

func (c *DocUpdateCommandAdapter) GetDescription() string {
	return "Update an existing document"
}

func (c *DocUpdateCommandAdapter) GetUsage() string {
	return "dw task-manager doc update <doc-id> [--content <content> | --from-file <path>] [--status <status>] [--detach]"
}

func (c *DocUpdateCommandAdapter) GetHelp() string {
	return `Updates an existing document's content, status, or attachments.

At least one field must be specified to update.

Flags:
  --content <content>       New markdown content
  --from-file <path>        Read new markdown content from file
  --status <status>         New document status: draft, published, archived
  --detach                  Remove document's track/iteration attachment
  --project <name>          Project name (optional, uses active project if not specified)

Notes:
  - Either --content or --from-file can be provided (not both)
  - Use --detach to remove all attachments without specifying new ones
  - Status must be one of: draft, published, archived

Examples:
  # Update document content from file
  dw task-manager doc update TM-doc-1 --from-file ./updated-content.md

  # Change document status and detach
  dw task-manager doc update TM-doc-1 --status published --detach

  # Update inline content
  dw task-manager doc update TM-doc-1 --content "Updated content"`
}

func (c *DocUpdateCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse doc ID from first arg
	if len(args) == 0 {
		return fmt.Errorf("document ID is required")
	}
	c.docID = args[0]

	// Parse flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--content":
			if i+1 < len(args) {
				content := args[i+1]
				c.content = &content
				i++
			}
		case "--from-file":
			if i+1 < len(args) {
				c.fromFile = args[i+1]
				i++
			}
		case "--status":
			if i+1 < len(args) {
				status := args[i+1]
				c.status = &status
				i++
			}
		case "--detach":
			c.detach = true
		}
	}

	// Validate XOR: content vs from-file
	if c.content != nil && c.fromFile != "" {
		return fmt.Errorf("--content and --from-file are mutually exclusive")
	}

	// Read file if provided
	if c.fromFile != "" {
		data, err := os.ReadFile(c.fromFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", c.fromFile, err)
		}
		content := string(data)
		c.content = &content
	}

	// Create DTO with detach flag
	input := dto.UpdateDocumentDTO{
		ID:      c.docID,
		Content: c.content,
		Status:  c.status,
		Detach:  c.detach,
	}

	// Execute via application service
	err := c.DocumentService.UpdateDocument(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Document %s updated successfully\n", c.docID)

	return nil
}

// ============================================================================
// DocShowCommandAdapter - Adapts CLI to GetDocument use case
// ============================================================================

type DocShowCommandAdapter struct {
	DocumentService *application.DocumentApplicationService

	// CLI flags
	project string
	docID   string
}

func (c *DocShowCommandAdapter) GetName() string {
	return "doc show"
}

func (c *DocShowCommandAdapter) GetDescription() string {
	return "Show document details"
}

func (c *DocShowCommandAdapter) GetUsage() string {
	return "dw task-manager doc show <doc-id>"
}

func (c *DocShowCommandAdapter) GetHelp() string {
	return `Shows detailed information about a document, including its title, type, status,
attachments, and full markdown content.

Flags:
  --project <name>  Project name (optional, uses active project if not specified)

Examples:
  dw task-manager doc show TM-doc-1`
}

func (c *DocShowCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse doc ID from first arg
	if len(args) == 0 {
		return fmt.Errorf("document ID is required")
	}
	c.docID = args[0]

	// Parse optional flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		}
	}

	// Execute via application service
	doc, err := c.DocumentService.GetDocument(ctx, c.docID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Document: %s\n", doc.ID)
	fmt.Fprintf(out, "  Title:       %s\n", doc.Title)
	fmt.Fprintf(out, "  Type:        %s\n", doc.Type)
	fmt.Fprintf(out, "  Status:      %s\n", doc.Status)
	fmt.Fprintf(out, "  Created:     %s\n", doc.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(out, "  Updated:     %s\n", doc.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Show attachment info
	if doc.TrackID != nil {
		fmt.Fprintf(out, "  Attachment:  Track %s\n", *doc.TrackID)
	} else if doc.IterationNumber != nil {
		fmt.Fprintf(out, "  Attachment:  Iteration %d\n", *doc.IterationNumber)
	} else {
		fmt.Fprintf(out, "  Attachment:  None (unattached)\n")
	}

	// Show content
	fmt.Fprintf(out, "\nContent:\n")
	fmt.Fprintf(out, "%s\n", doc.Content)

	return nil
}

// ============================================================================
// DocListCommandAdapter - Adapts CLI to ListDocuments use case
// ============================================================================

type DocListCommandAdapter struct {
	DocumentService *application.DocumentApplicationService

	// CLI flags
	project   string
	track     *string
	iteration *int
	docType   *string
}

func (c *DocListCommandAdapter) GetName() string {
	return "doc list"
}

func (c *DocListCommandAdapter) GetDescription() string {
	return "List documents with optional filters"
}

func (c *DocListCommandAdapter) GetUsage() string {
	return "dw task-manager doc list [--track <id>] [--iteration <num>] [--type <type>]"
}

func (c *DocListCommandAdapter) GetHelp() string {
	return `Lists all documents with optional filtering by track, iteration, or type.

If no filters are provided, lists all documents.

Flags:
  --track <id>      Filter by attached track ID
  --iteration <num> Filter by attached iteration number
  --type <type>     Filter by document type: adr, plan, retrospective, other
  --project <name>  Project name (optional, uses active project if not specified)

Notes:
  - Only one of --track, --iteration, --type should be provided
  - Documents can be filtered or listed without filters for all documents

Examples:
  # List all documents
  dw task-manager doc list

  # List documents attached to a track
  dw task-manager doc list --track TM-track-1

  # List all ARD documents
  dw task-manager doc list --type adr`
}

func (c *DocListCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
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
				trackID := args[i+1]
				c.track = &trackID
				i++
			}
		case "--iteration":
			if i+1 < len(args) {
				iterNum, err := strconv.Atoi(args[i+1])
				if err != nil {
					return fmt.Errorf("invalid iteration number: %v", err)
				}
				c.iteration = &iterNum
				i++
			}
		case "--type":
			if i+1 < len(args) {
				docType := args[i+1]
				c.docType = &docType
				i++
			}
		}
	}

	// Execute via application service
	docs, err := c.DocumentService.ListDocuments(ctx, c.track, c.iteration, c.docType)
	if err != nil {
		return fmt.Errorf("failed to list documents: %w", err)
	}

	// Format output as table
	out := cmdCtx.GetStdout()
	if len(docs) == 0 {
		fmt.Fprintf(out, "No documents found\n")
		return nil
	}

	// Print header
	fmt.Fprintf(out, "%-15s %-30s %-15s %-12s %-20s\n", "ID", "Title", "Type", "Status", "Attachment")
	fmt.Fprintf(out, "%s %s %s %s %s\n",
		strings.Repeat("-", 15),
		strings.Repeat("-", 30),
		strings.Repeat("-", 15),
		strings.Repeat("-", 12),
		strings.Repeat("-", 20),
	)

	// Print documents
	for _, doc := range docs {
		attachment := "Unattached"
		if doc.TrackID != nil {
			attachment = fmt.Sprintf("Track %s", *doc.TrackID)
		} else if doc.IterationNumber != nil {
			attachment = fmt.Sprintf("Iteration %d", *doc.IterationNumber)
		}

		// Truncate title if too long
		title := doc.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		fmt.Fprintf(out, "%-15s %-30s %-15s %-12s %-20s\n",
			doc.ID, title, doc.Type, doc.Status, attachment)
	}

	return nil
}

// ============================================================================
// DocAttachCommandAdapter - Adapts CLI to AttachDocument use case
// ============================================================================

type DocAttachCommandAdapter struct {
	DocumentService *application.DocumentApplicationService

	// CLI flags
	project   string
	docID     string
	track     *string
	iteration *int
}

func (c *DocAttachCommandAdapter) GetName() string {
	return "doc attach"
}

func (c *DocAttachCommandAdapter) GetDescription() string {
	return "Attach a document to a track or iteration"
}

func (c *DocAttachCommandAdapter) GetUsage() string {
	return "dw task-manager doc attach <doc-id> (--track <id> | --iteration <num>)"
}

func (c *DocAttachCommandAdapter) GetHelp() string {
	return `Attaches an unattached document to a track or iteration.

You must provide exactly one of --track or --iteration.

Flags:
  --track <id>      Attach to track ID (mutually exclusive with --iteration)
  --iteration <num> Attach to iteration number (mutually exclusive with --track)
  --project <name>  Project name (optional, uses active project if not specified)

Examples:
  # Attach to track
  dw task-manager doc attach TM-doc-1 --track TM-track-1

  # Attach to iteration
  dw task-manager doc attach TM-doc-1 --iteration 5`
}

func (c *DocAttachCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse doc ID from first arg
	if len(args) == 0 {
		return fmt.Errorf("document ID is required")
	}
	c.docID = args[0]

	// Parse flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		case "--track":
			if i+1 < len(args) {
				trackID := args[i+1]
				c.track = &trackID
				i++
			}
		case "--iteration":
			if i+1 < len(args) {
				iterNum, err := strconv.Atoi(args[i+1])
				if err != nil {
					return fmt.Errorf("invalid iteration number: %v", err)
				}
				c.iteration = &iterNum
				i++
			}
		}
	}

	// Validate XOR: track vs iteration (one required)
	if (c.track == nil || *c.track == "") && (c.iteration == nil || *c.iteration < 1) {
		return fmt.Errorf("either --track or --iteration is required")
	}
	if c.track != nil && *c.track != "" && c.iteration != nil && *c.iteration > 0 {
		return fmt.Errorf("--track and --iteration are mutually exclusive")
	}

	// Execute via application service
	err := c.DocumentService.AttachDocument(ctx, c.docID, c.track, c.iteration)
	if err != nil {
		return fmt.Errorf("failed to attach document: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	if c.track != nil && *c.track != "" {
		fmt.Fprintf(out, "Document %s attached to track %s\n", c.docID, *c.track)
	} else if c.iteration != nil && *c.iteration > 0 {
		fmt.Fprintf(out, "Document %s attached to iteration %d\n", c.docID, *c.iteration)
	}

	return nil
}

// ============================================================================
// DocDetachCommandAdapter - Adapts CLI to DetachDocument use case
// ============================================================================

type DocDetachCommandAdapter struct {
	DocumentService *application.DocumentApplicationService

	// CLI flags
	project string
	docID   string
}

func (c *DocDetachCommandAdapter) GetName() string {
	return "doc detach"
}

func (c *DocDetachCommandAdapter) GetDescription() string {
	return "Detach a document from its attachment"
}

func (c *DocDetachCommandAdapter) GetUsage() string {
	return "dw task-manager doc detach <doc-id>"
}

func (c *DocDetachCommandAdapter) GetHelp() string {
	return `Detaches a document from its track or iteration attachment, making it unattached.

Flags:
  --project <name>  Project name (optional, uses active project if not specified)

Examples:
  dw task-manager doc detach TM-doc-1`
}

func (c *DocDetachCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse doc ID from first arg
	if len(args) == 0 {
		return fmt.Errorf("document ID is required")
	}
	c.docID = args[0]

	// Parse optional flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		}
	}

	// Execute via application service
	err := c.DocumentService.DetachDocument(ctx, c.docID)
	if err != nil {
		return fmt.Errorf("failed to detach document: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Document %s detached\n", c.docID)

	return nil
}

// ============================================================================
// DocDeleteCommandAdapter - Adapts CLI to DeleteDocument use case
// ============================================================================

type DocDeleteCommandAdapter struct {
	DocumentService *application.DocumentApplicationService

	// CLI flags
	project string
	docID   string
	force   bool
}

func (c *DocDeleteCommandAdapter) GetName() string {
	return "doc delete"
}

func (c *DocDeleteCommandAdapter) GetDescription() string {
	return "Delete a document"
}

func (c *DocDeleteCommandAdapter) GetUsage() string {
	return "dw task-manager doc delete <doc-id> [--force]"
}

func (c *DocDeleteCommandAdapter) GetHelp() string {
	return `Deletes a document. By default, prompts for confirmation unless --force is used.

Flags:
  --force          Skip confirmation prompt
  --project <name> Project name (optional, uses active project if not specified)

Examples:
  # Delete with confirmation prompt
  dw task-manager doc delete TM-doc-1

  # Delete without confirmation
  dw task-manager doc delete TM-doc-1 --force`
}

func (c *DocDeleteCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	// Parse doc ID from first arg
	if len(args) == 0 {
		return fmt.Errorf("document ID is required")
	}
	c.docID = args[0]

	// Parse optional flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--force":
			c.force = true
		case "--project":
			if i+1 < len(args) {
				c.project = args[i+1]
				i++
			}
		}
	}

	// Prompt for confirmation unless --force
	if !c.force {
		out := cmdCtx.GetStdout()
		fmt.Fprintf(out, "Are you sure you want to delete document %s? (yes/no): ", c.docID)

		// Read user input
		var response string
		_, err := fmt.Scanf("%s", &response)
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		if strings.ToLower(response) != "yes" && strings.ToLower(response) != "y" {
			fmt.Fprintf(out, "Deletion cancelled\n")
			return nil
		}
	}

	// Execute via application service
	err := c.DocumentService.DeleteDocument(ctx, c.docID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Format output
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "Document %s deleted\n", c.docID)

	return nil
}

// ============================================================================
// DocHelpCommandAdapter - Shows help for all document commands
// ============================================================================

type DocHelpCommandAdapter struct {
	// No service needed for help command
}

func (c *DocHelpCommandAdapter) GetName() string {
	return "doc"
}

func (c *DocHelpCommandAdapter) GetDescription() string {
	return "Document management commands"
}

func (c *DocHelpCommandAdapter) GetUsage() string {
	return "dw task-manager doc <command> [options]"
}

func (c *DocHelpCommandAdapter) GetHelp() string {
	return `Document Management Commands

Documents are markdown content that can be attached to tracks or iterations,
such as architecture decisions (ADRs), planning documents, retrospectives, etc.

Available Commands:
  doc create       Create a new document
  doc update       Update an existing document
  doc show         Show document details
  doc list         List documents with optional filters
  doc attach       Attach a document to a track or iteration
  doc detach       Detach a document from its attachment
  doc delete       Delete a document
  doc --help       Show this help message

Quick Examples:

Create a document from a file:
  dw task-manager doc create \
    --title "API Design" \
    --type plan \
    --from-file ./docs/api-design.md \
    --track TM-track-1

Create an inline document:
  dw task-manager doc create \
    --title "Sprint Retrospective" \
    --type retrospective \
    --content "# What went well..."

View a document:
  dw task-manager doc show TM-doc-1

List all documents:
  dw task-manager doc list

Filter documents by type:
  dw task-manager doc list --type adr

Update document content:
  dw task-manager doc update TM-doc-1 --from-file ./updated.md

Attach to different location:
  dw task-manager doc attach TM-doc-1 --iteration 5

Delete a document:
  dw task-manager doc delete TM-doc-1 --force

For detailed help on each command, use:
  dw task-manager <command> --help

Examples:
  dw task-manager doc create --help
  dw task-manager doc list --help`
}

func (c *DocHelpCommandAdapter) Execute(ctx context.Context, cmdCtx pluginsdk.CommandContext, args []string) error {
	out := cmdCtx.GetStdout()
	fmt.Fprintf(out, "%s\n", c.GetHelp())
	return nil
}
