package app

import (
	"context"
	"fmt"
	"io"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

// RefreshCommandHandler handles the refresh command logic.
// This performs framework-level refresh (database schema, config).
// Plugin-specific refresh is handled by each plugin's refresh command.
type RefreshCommandHandler struct {
	repo         domain.EventRepository
	configLoader ConfigLoader
	logger       Logger
	out          io.Writer
}

// ConfigLoader interface for config loading
type ConfigLoader interface {
	LoadConfig(path string) (*domain.Config, error)
	InitializeDefaultConfig(path string) (string, error)
}

// NewRefreshCommandHandler creates a new refresh command handler
func NewRefreshCommandHandler(
	repo domain.EventRepository,
	configLoader ConfigLoader,
	logger Logger,
	out io.Writer,
) *RefreshCommandHandler {
	return &RefreshCommandHandler{
		repo:         repo,
		configLoader: configLoader,
		logger:       logger,
		out:          out,
	}
}

// Execute runs the framework-level refresh operation
func (h *RefreshCommandHandler) Execute(ctx context.Context, dbPath string) error {
	fmt.Fprintln(h.out, "Refreshing DarwinFlow framework...")
	fmt.Fprintln(h.out)

	// Step 1: Update database schema
	fmt.Fprintln(h.out, "Updating database schema...")
	if err := h.repo.Initialize(ctx); err != nil {
		return fmt.Errorf("error updating database schema: %w", err)
	}
	fmt.Fprintf(h.out, "✓ Database schema updated: %s\n", dbPath)

	// Step 2: Initialize config if needed
	fmt.Fprintln(h.out)
	fmt.Fprintln(h.out, "Checking configuration...")

	// Try to load config
	config, err := h.configLoader.LoadConfig("")
	if err != nil || config == nil {
		// Config doesn't exist or is invalid, create default
		fmt.Fprintln(h.out, "Creating default configuration...")
		configPath, err := h.configLoader.InitializeDefaultConfig("")
		if err != nil {
			fmt.Fprintf(h.out, "Warning: Could not create default config: %v\n", err)
		} else {
			fmt.Fprintf(h.out, "✓ Configuration initialized: %s\n", configPath)
		}
	} else {
		fmt.Fprintln(h.out, "✓ Configuration is valid")
	}

	// Done
	fmt.Fprintln(h.out)
	fmt.Fprintln(h.out, "✓ Framework refreshed successfully!")
	fmt.Fprintln(h.out)
	fmt.Fprintln(h.out, "Framework changes applied:")
	fmt.Fprintln(h.out, "  • Database schema updated with latest migrations")
	fmt.Fprintln(h.out, "  • Configuration verified")
	fmt.Fprintln(h.out)
	fmt.Fprintln(h.out, "Note: To refresh plugin-specific configuration, run the plugin's init command")
	fmt.Fprintln(h.out, "Example: dw <plugin-name> init")

	return nil
}
