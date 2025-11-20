package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/application"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/application/dto"
	"github.com/spf13/cobra"
)

// ============================================================================
// NewDocCommands returns the doc command group for Cobra
// ============================================================================

// NewDocCommands creates and returns the doc command group with all subcommands.
func NewDocCommands(docService *application.DocumentApplicationService) *cobra.Command {
	docCmd := &cobra.Command{
		Use:     "doc",
		Short:   "Manage documents",
		Long: `Document Management Commands

Commands for creating, updating, and managing documents (ADRs, plans, retrospectives, etc.).

Examples:
  # Create document from file attached to track
  tm doc create --title "API Design" --type plan --from-file ./docs/api.md --track TM-track-1

  # Show document details
  tm doc show TM-doc-1

  # List all ADR documents
  tm doc list --type adr

  # Update document content
  tm doc update TM-doc-1 --from-file ./updated.md

  # Attach document to track or iteration
  tm doc attach TM-doc-1 --track TM-track-2
  tm doc attach TM-doc-1 --iteration 5

  # Detach document
  tm doc detach TM-doc-1

  # Delete document
  tm doc delete TM-doc-1`,
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add all document subcommands
	docCmd.AddCommand(
		newDocCreateCommand(docService),
		newDocListCommand(docService),
		newDocShowCommand(docService),
		newDocUpdateCommand(docService),
		newDocAttachCommand(docService),
		newDocDetachCommand(docService),
		newDocDeleteCommand(docService),
	)

	return docCmd
}

// ============================================================================
// doc create command
// ============================================================================

func newDocCreateCommand(docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new document",
		Long:  `Creates a new document with markdown content. Content can be provided inline or from a file.`,
		Example: `  # Create document from file attached to track
  tm doc create --title "API Design" --type plan --from-file ./docs/api.md --track TM-track-1

  # Create inline document
  tm doc create --title "Sprint Retrospective" --type retrospective --content "## What went well..."

  # Create attached to iteration
  tm doc create --title "Architecture Decision" --type adr --from-file ./adr.md --iteration 5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			title, _ := cmd.Flags().GetString("title")
			docType, _ := cmd.Flags().GetString("type")
			status, _ := cmd.Flags().GetString("status")
			content, _ := cmd.Flags().GetString("content")
			fromFile, _ := cmd.Flags().GetString("from-file")
			track, _ := cmd.Flags().GetString("track")
			iteration, _ := cmd.Flags().GetString("iteration")

			// Validate required flags
			if title == "" {
				return fmt.Errorf("--title is required")
			}
			if docType == "" {
				return fmt.Errorf("--type is required")
			}

			// Validate XOR: content vs from-file
			if content != "" && fromFile != "" {
				return fmt.Errorf("--content and --from-file are mutually exclusive (provide one, not both)")
			}
			if content == "" && fromFile == "" {
				return fmt.Errorf("either --content or --from-file is required")
			}

			// Read file if provided
			if fromFile != "" {
				data, err := os.ReadFile(fromFile)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", fromFile, err)
				}
				content = string(data)
			}

			// Parse track and iteration flags
			var trackID *string
			var iterationNum *int
			if track != "" {
				trackID = &track
			}
			if iteration != "" {
				num, err := strconv.Atoi(iteration)
				if err != nil {
					return fmt.Errorf("invalid iteration number: %v", err)
				}
				iterationNum = &num
			}

			// Create DTO
			input := dto.CreateDocumentDTO{
				Title:           title,
				Type:            docType,
				Status:          status,
				Content:         content,
				TrackID:         trackID,
				IterationNumber: iterationNum,
			}

			// Execute via application service
			docID, err := docService.CreateDocument(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Document created successfully\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:    %s\n", docID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title: %s\n", title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Type:  %s\n", docType)
			if track != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Track: %s\n", track)
			}
			if iteration != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Iteration: %s\n", iteration)
			}

			return nil
		},
	}

	cmd.Flags().String("title", "", "Document title (required)")
	cmd.Flags().String("type", "", "Document type (required): adr, plan, retrospective, other")
	cmd.Flags().String("status", "draft", "Document status (optional): draft, published, archived")
	cmd.Flags().String("content", "", "Document markdown content (XOR with --from-file)")
	cmd.Flags().String("from-file", "", "Read markdown content from file (XOR with --content)")
	cmd.Flags().String("track", "", "Attach to track ID (optional, XOR with --iteration)")
	cmd.Flags().String("iteration", "", "Attach to iteration number (optional, XOR with --track)")

	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("type")

	return cmd
}

// ============================================================================
// doc list command
// ============================================================================

func newDocListCommand(docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List documents with optional filters",
		Long:  `Lists all documents with optional filtering by track, iteration, or type.`,
		Example: `  # List all documents
  tm doc list

  # List documents attached to a track
  tm doc list --track TM-track-1

  # List all ADR documents
  tm doc list --type adr`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			track, _ := cmd.Flags().GetString("track")
			iteration, _ := cmd.Flags().GetString("iteration")
			docType, _ := cmd.Flags().GetString("type")

			// Convert flags to pointers for optional values
			var trackPtr *string
			var iterationPtr *int
			var docTypePtr *string

			if track != "" {
				trackPtr = &track
			}
			if iteration != "" {
				num, err := strconv.Atoi(iteration)
				if err != nil {
					return fmt.Errorf("invalid iteration number: %v", err)
				}
				iterationPtr = &num
			}
			if docType != "" {
				docTypePtr = &docType
			}

			// Execute via application service
			docs, err := docService.ListDocuments(ctx, trackPtr, iterationPtr, docTypePtr)
			if err != nil {
				return fmt.Errorf("failed to list documents: %w", err)
			}

			// Format output
			if len(docs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No documents found\n")
				return nil
			}

			// Print header
			fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-30s %-15s %-12s %-20s\n", "ID", "Title", "Type", "Status", "Attachment")
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s %s %s\n",
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

				fmt.Fprintf(cmd.OutOrStdout(), "%-15s %-30s %-15s %-12s %-20s\n",
					doc.ID,
					truncateString(doc.Title, 30),
					doc.Type,
					doc.Status,
					attachment,
				)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d document(s)\n", len(docs))
			return nil
		},
	}

	cmd.Flags().String("track", "", "Filter by attached track ID (optional)")
	cmd.Flags().String("iteration", "", "Filter by attached iteration number (optional)")
	cmd.Flags().String("type", "", "Filter by document type: adr, plan, retrospective, other (optional)")

	return cmd
}

// ============================================================================
// doc show command
// ============================================================================

func newDocShowCommand(docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <doc-id>",
		Short: "Show document details",
		Long:  `Displays detailed information about a document including its content, type, and attachments.`,
		Example: `  # Show document details
  tm doc show TM-doc-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			docID := args[0]

			// Execute via application service
			doc, err := docService.GetDocument(ctx, docID)
			if err != nil {
				return fmt.Errorf("failed to get document: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Document Details\n")
			fmt.Fprintf(cmd.OutOrStdout(), "================\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ID:     %s\n", doc.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Title:  %s\n", doc.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "  Type:   %s\n", doc.Type)
			fmt.Fprintf(cmd.OutOrStdout(), "  Status: %s\n", doc.Status)

			// Show attachment info
			if doc.TrackID != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  Track %s\n", *doc.TrackID)
			} else if doc.IterationNumber != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  Iteration %d\n", *doc.IterationNumber)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "  Attachment: None (unattached)\n")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "  Created: %s\n", doc.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
			fmt.Fprintf(cmd.OutOrStdout(), "  Updated: %s\n", doc.UpdatedAt.Format("2006-01-02 15:04:05 UTC"))

			// Show content
			fmt.Fprintf(cmd.OutOrStdout(), "\nContent:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "--------\n")
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", doc.Content)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// doc update command
// ============================================================================

func newDocUpdateCommand(docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <doc-id>",
		Short: "Update an existing document",
		Long:  `Updates one or more fields of an existing document. At least one field must be specified.`,
		Example: `  # Update document content from file
  tm doc update TM-doc-1 --from-file ./updated.md

  # Change document status
  tm doc update TM-doc-1 --status published

  # Update inline content
  tm doc update TM-doc-1 --content "Updated content"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			docID := args[0]

			content, _ := cmd.Flags().GetString("content")
			fromFile, _ := cmd.Flags().GetString("from-file")
			status, _ := cmd.Flags().GetString("status")
			detach, _ := cmd.Flags().GetBool("detach")

			contentSet := cmd.Flags().Changed("content")
			statusSet := cmd.Flags().Changed("status")

			// Check that at least one field is being updated
			if !contentSet && !statusSet && !detach && fromFile == "" {
				return fmt.Errorf("at least one field must be specified to update (--content, --from-file, --status, or --detach)")
			}

			// Validate XOR: content vs from-file
			if content != "" && fromFile != "" {
				return fmt.Errorf("--content and --from-file are mutually exclusive")
			}

			// Read file if provided
			var finalContent *string
			if fromFile != "" {
				data, err := os.ReadFile(fromFile)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", fromFile, err)
				}
				c := string(data)
				finalContent = &c
			} else if contentSet {
				finalContent = &content
			}

			// Convert status to pointer
			var finalStatus *string
			if statusSet {
				finalStatus = &status
			}

			// Create DTO with only updated fields
			input := dto.UpdateDocumentDTO{
				ID:      docID,
				Content: finalContent,
				Status:  finalStatus,
				Detach:  detach,
			}

			// Execute via application service
			err := docService.UpdateDocument(ctx, input)
			if err != nil {
				return fmt.Errorf("failed to update document: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Document %s updated successfully\n", docID)

			return nil
		},
	}

	cmd.Flags().String("content", "", "New document markdown content")
	cmd.Flags().String("from-file", "", "Read new markdown content from file")
	cmd.Flags().String("status", "", "New document status: draft, published, archived")
	cmd.Flags().Bool("detach", false, "Remove document's track/iteration attachment")

	return cmd
}

// ============================================================================
// doc attach command
// ============================================================================

func newDocAttachCommand(docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach <doc-id>",
		Short: "Attach a document to a track or iteration",
		Long:  `Attaches an unattached document to a track or iteration. You must provide exactly one of --track or --iteration.`,
		Example: `  # Attach to track
  tm doc attach TM-doc-1 --track TM-track-1

  # Attach to iteration
  tm doc attach TM-doc-1 --iteration 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			docID := args[0]

			track, _ := cmd.Flags().GetString("track")
			iteration, _ := cmd.Flags().GetString("iteration")

			// Validate XOR: track vs iteration (one required)
			if track == "" && iteration == "" {
				return fmt.Errorf("either --track or --iteration is required")
			}
			if track != "" && iteration != "" {
				return fmt.Errorf("--track and --iteration are mutually exclusive")
			}

			// Convert to pointers
			var trackPtr *string
			var iterationPtr *int
			if track != "" {
				trackPtr = &track
			}
			if iteration != "" {
				num, err := strconv.Atoi(iteration)
				if err != nil {
					return fmt.Errorf("invalid iteration number: %v", err)
				}
				iterationPtr = &num
			}

			// Execute via application service
			err := docService.AttachDocument(ctx, docID, trackPtr, iterationPtr)
			if err != nil {
				return fmt.Errorf("failed to attach document: %w", err)
			}

			// Format output
			if track != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Document %s attached to track %s successfully\n", docID, track)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Document %s attached to iteration %s successfully\n", docID, iteration)
			}

			return nil
		},
	}

	cmd.Flags().String("track", "", "Attach to track ID (mutually exclusive with --iteration)")
	cmd.Flags().String("iteration", "", "Attach to iteration number (mutually exclusive with --track)")

	return cmd
}

// ============================================================================
// doc detach command
// ============================================================================

func newDocDetachCommand(docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach <doc-id>",
		Short: "Detach a document from its attachment",
		Long:  `Detaches a document from its track or iteration attachment, making it unattached.`,
		Example: `  # Detach a document
  tm doc detach TM-doc-1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			docID := args[0]

			// Execute via application service
			err := docService.DetachDocument(ctx, docID)
			if err != nil {
				return fmt.Errorf("failed to detach document: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Document %s detached successfully\n", docID)

			return nil
		},
	}

	return cmd
}

// ============================================================================
// doc delete command
// ============================================================================

func newDocDeleteCommand(docService *application.DocumentApplicationService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <doc-id>",
		Short: "Delete a document",
		Long:  `Deletes a document. By default, prompts for confirmation unless --force is used.`,
		Example: `  # Delete with confirmation prompt
  tm doc delete TM-doc-1

  # Delete without confirmation
  tm doc delete TM-doc-1 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			docID := args[0]

			force, _ := cmd.Flags().GetBool("force")

			// Prompt for confirmation unless --force
			if !force {
				fmt.Fprintf(cmd.OutOrStdout(), "Are you sure you want to delete document %s? (yes/no): ", docID)

				// Read user input
				var response string
				_, err := fmt.Scanf("%s", &response)
				if err != nil {
					return fmt.Errorf("failed to read confirmation: %w", err)
				}

				if strings.ToLower(response) != "yes" && strings.ToLower(response) != "y" {
					fmt.Fprintf(cmd.OutOrStdout(), "Deletion cancelled\n")
					return nil
				}
			}

			// Execute via application service
			err := docService.DeleteDocument(ctx, docID)
			if err != nil {
				return fmt.Errorf("failed to delete document: %w", err)
			}

			// Format output
			fmt.Fprintf(cmd.OutOrStdout(), "Document %s deleted successfully\n", docID)

			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")

	return cmd
}
