package services

import (
	"fmt"
	"regexp"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ValidationService provides centralized domain validation rules
type ValidationService struct{}

// NewValidationService creates a new validation service
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// ValidateTrackID validates track ID format
// Accepts both old format (track-<slug>) and new format (<CODE>-track-<number>)
func (s *ValidationService) ValidateTrackID(id string) error {
	// New format: <CODE>-track-<number> (e.g., DW-track-1, PROD-track-123)
	newPattern := `^[A-Z0-9]+-track-[0-9]+$`
	newRegex := regexp.MustCompile(newPattern)
	if newRegex.MatchString(id) {
		return nil
	}

	// Old format: track-<slug> (for backward compatibility)
	oldPattern := `^track-[a-z0-9]+(-[a-z0-9]+)*$`
	oldRegex := regexp.MustCompile(oldPattern)
	if oldRegex.MatchString(id) {
		return nil
	}

	return fmt.Errorf("%w: track ID must follow convention: track-<slug> or <CODE>-track-<number>", pluginsdk.ErrInvalidArgument)
}

// ValidateIterationNumber validates iteration number is positive
func (s *ValidationService) ValidateIterationNumber(number int) error {
	if number <= 0 {
		return fmt.Errorf("%w: iteration number must be positive", pluginsdk.ErrInvalidArgument)
	}
	return nil
}

// ValidateRank validates rank is within valid range (1-1000)
func (s *ValidationService) ValidateRank(rank int) error {
	if rank < 1 || rank > 1000 {
		return fmt.Errorf("%w: rank must be between 1 and 1000", pluginsdk.ErrInvalidArgument)
	}
	return nil
}

// ValidateNonEmpty validates a string field is non-empty
func (s *ValidationService) ValidateNonEmpty(fieldName, value string) error {
	if value == "" {
		return fmt.Errorf("%w: %s must be non-empty", pluginsdk.ErrInvalidArgument, fieldName)
	}
	return nil
}
