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

// HookConfigManager defines the interface for managing Claude Code hook configuration
type HookConfigManager interface {
	InstallDarwinFlowHooks() error
	GetSettingsPath() string
}

// SetupService orchestrates initialization of the DarwinFlow logging infrastructure
type SetupService struct {
	repository        domain.EventRepository
	hookConfigManager HookConfigManager
}

// NewSetupService creates a new setup service
func NewSetupService(
	repository domain.EventRepository,
	hookConfigManager HookConfigManager,
) *SetupService {
	return &SetupService{
		repository:        repository,
		hookConfigManager: hookConfigManager,
	}
}

// Initialize sets up the complete logging infrastructure
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

	// Install hooks in Claude Code settings
	if err := s.hookConfigManager.InstallDarwinFlowHooks(); err != nil {
		return fmt.Errorf("failed to install hooks: %w", err)
	}

	return nil
}

// GetSettingsPath returns the path to the Claude Code settings file
func (s *SetupService) GetSettingsPath() string {
	return s.hookConfigManager.GetSettingsPath()
}
