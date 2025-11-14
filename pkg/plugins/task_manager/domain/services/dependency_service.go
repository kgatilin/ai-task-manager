package services

import (
	"context"
	"fmt"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// DependencyService handles circular dependency detection for tracks
type DependencyService struct{}

// NewDependencyService creates a new dependency service
func NewDependencyService() *DependencyService {
	return &DependencyService{}
}

// ValidateNoCycles checks if the given track has any circular dependencies
// Uses depth-first search algorithm to detect cycles
// Returns ErrInvalidArgument if a cycle is detected
func (s *DependencyService) ValidateNoCycles(
	ctx context.Context,
	trackID string,
	getDependencies func(context.Context, string) ([]string, error),
) error {
	visited := make(map[string]bool)
	return s.detectCycleDFS(ctx, trackID, visited, getDependencies)
}

// detectCycleDFS performs depth-first search to detect cycles
// The visited map tracks:
// - true = currently in the path (visiting)
// - false = fully processed (visited)
func (s *DependencyService) detectCycleDFS(
	ctx context.Context,
	trackID string,
	visited map[string]bool,
	getDependencies func(context.Context, string) ([]string, error),
) error {
	// If we're revisiting a node that's in the current path, we have a cycle
	if visited[trackID] {
		return fmt.Errorf("%w: circular dependency detected for track %s", pluginsdk.ErrInvalidArgument, trackID)
	}

	// Mark node as in the current path
	visited[trackID] = true

	// Get dependencies for this track
	deps, err := getDependencies(ctx, trackID)
	if err != nil {
		return fmt.Errorf("failed to get dependencies for track %s: %w", trackID, err)
	}

	// Recursively check all dependencies
	for _, depID := range deps {
		if err := s.detectCycleDFS(ctx, depID, visited, getDependencies); err != nil {
			return err
		}
	}

	// Mark node as fully processed (no longer in current path)
	visited[trackID] = false
	return nil
}
