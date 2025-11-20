package services_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/errors"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
	"github.com/stretchr/testify/assert"
)

func TestNewValidationService(t *testing.T) {
	svc := services.NewValidationService()
	assert.NotNil(t, svc)
}

func TestValidationService_ValidateTrackID(t *testing.T) {
	svc := services.NewValidationService()

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		// New format: <CODE>-track-<number>
		{"new format - DW", "DW-track-1", false},
		{"new format - PROD", "PROD-track-123", false},
		{"new format - multi digit code", "MYAPP-track-1", false},
		{"new format - TM", "TM-track-5", false},
		{"new format - numeric code", "123-track-1", false},
		{"new format - mixed alphanumeric", "ABC123-track-999", false},

		// Old format: track-<slug>
		{"old format - simple", "track-auth", false},
		{"old format - hyphenated", "track-user-management", false},
		{"old format - numbers", "track-v2-api", false},
		{"old format - single char", "track-a", false},
		{"old format - multiple hyphens", "track-a-b-c", false},

		// Invalid formats
		{"invalid - no prefix", "nottrack-1", true},
		{"invalid - wrong separator new format", "DW_track_1", true},
		{"invalid - wrong separator old format", "track_test", true},
		{"invalid - space in old format", "track-auth system", true},
		{"invalid - space in new format", "DW-track-1 2", true},
		{"invalid - empty", "", true},
		{"invalid - just track", "track", true},
		{"invalid - uppercase in old format", "track-Auth", true},
		{"invalid - special chars in old format", "track-@#$", true},
		{"invalid - lowercase code in new format", "dw-track-1", true},
		{"invalid - missing number in new format", "DW-track-", true},
		{"invalid - missing code in new format", "-track-1", true},
		{"invalid - trailing hyphen in old format", "track-auth-", true},
		{"invalid - double hyphen in old format", "track-auth--system", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateTrackID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, errors.ErrInvalidArgument)
				assert.Contains(t, err.Error(), "track ID must follow convention")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationService_ValidateIterationNumber(t *testing.T) {
	svc := services.NewValidationService()

	tests := []struct {
		name    string
		number  int
		wantErr bool
	}{
		{"valid - 1", 1, false},
		{"valid - 10", 10, false},
		{"valid - 100", 100, false},
		{"valid - 999", 999, false},
		{"valid - large number", 999999, false},
		{"invalid - 0", 0, true},
		{"invalid - negative small", -1, true},
		{"invalid - negative large", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateIterationNumber(tt.number)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, errors.ErrInvalidArgument)
				assert.Contains(t, err.Error(), "iteration number must be positive")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationService_ValidateRank(t *testing.T) {
	svc := services.NewValidationService()

	tests := []struct {
		name    string
		rank    int
		wantErr bool
	}{
		{"valid - 1 (min)", 1, false},
		{"valid - 500 (mid)", 500, false},
		{"valid - 1000 (max)", 1000, false},
		{"valid - 100", 100, false},
		{"valid - 250", 250, false},
		{"valid - 750", 750, false},
		{"invalid - 0", 0, true},
		{"invalid - 1001 (just over)", 1001, true},
		{"invalid - negative", -1, true},
		{"invalid - large out of range", 5000, true},
		{"invalid - large negative", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateRank(tt.rank)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, errors.ErrInvalidArgument)
				assert.Contains(t, err.Error(), "rank")
				assert.Contains(t, err.Error(), "1 and 1000")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationService_ValidateNonEmpty(t *testing.T) {
	svc := services.NewValidationService()

	tests := []struct {
		name      string
		fieldName string
		value     string
		wantErr   bool
	}{
		{"valid - with value", "title", "Some Title", false},
		{"valid - with spaces", "description", "Text with spaces", false},
		{"valid - single char", "name", "A", false},
		{"valid - numbers", "id", "12345", false},
		{"valid - special chars", "data", "!@#$%", false},
		{"valid - whitespace only", "description", "   ", false},
		{"valid - newline", "content", "\n", false},
		{"valid - tab", "content", "\t", false},
		{"invalid - empty string", "title", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateNonEmpty(tt.fieldName, tt.value)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, errors.ErrInvalidArgument)
				assert.Contains(t, err.Error(), tt.fieldName)
				assert.Contains(t, err.Error(), "must be non-empty")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
