package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

const (
	// DefaultDBPath is the default location for the event database
	DefaultDBPath = ".darwinflow/logs/events.db"
)

// SetupService orchestrates initialization of the DarwinFlow framework infrastructure.
// This is framework-level setup only (database, schema, etc.).
// Plugins handle their own initialization via their init commands.
type SetupService struct {
	repository domain.EventRepository
	logger     Logger
}

// NewSetupService creates a new setup service
func NewSetupService(
	repository domain.EventRepository,
	logger Logger,
) *SetupService {
	return &SetupService{
		repository: repository,
		logger:     logger,
	}
}

// Initialize sets up the framework-level infrastructure.
// This includes creating the database directory and initializing the repository schema.
// Plugin-specific initialization is handled by each plugin's init command.
func (s *SetupService) Initialize(ctx context.Context, dbPath string) error {
	// Create database directory
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Initialize repository (create schema, indexes, etc.)
	if err := s.repository.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	s.logger.Info("Framework infrastructure initialized")
	return nil
}
