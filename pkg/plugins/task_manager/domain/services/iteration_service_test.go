package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/services"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/stretchr/testify/assert"
)

func TestNewIterationService(t *testing.T) {
	svc := services.NewIterationService()
	assert.NotNil(t, svc)
}

func TestIterationService_CanCompleteIteration(t *testing.T) {
	svc := services.NewIterationService()

	tests := []struct {
		name       string
		status     string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "current iteration - valid",
			status:  string(entities.IterationStatusCurrent),
			wantErr: false,
		},
		{
			name:       "planned iteration - invalid",
			status:     string(entities.IterationStatusPlanned),
			wantErr:    true,
			wantErrMsg: "current status",
		},
		{
			name:       "complete iteration - invalid",
			status:     string(entities.IterationStatusComplete),
			wantErr:    true,
			wantErrMsg: "current status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iter := createTestIteration(1, "Test", tt.status)

			err := svc.CanCompleteIteration(iter)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIterationService_CanStartIteration(t *testing.T) {
	svc := services.NewIterationService()
	ctx := context.Background()

	t.Run("valid - no current iteration", func(t *testing.T) {
		iter := createTestIteration(1, "Test", string(entities.IterationStatusPlanned))

		// Callback that returns ErrNotFound (no current iteration)
		getCurrentIter := func(ctx context.Context) (*entities.IterationEntity, error) {
			return nil, pluginsdk.ErrNotFound
		}

		err := svc.CanStartIteration(ctx, iter, getCurrentIter)
		assert.NoError(t, err)
	})

	t.Run("invalid - not planned status", func(t *testing.T) {
		iter := createTestIteration(1, "Test", string(entities.IterationStatusCurrent))

		getCurrentIter := func(ctx context.Context) (*entities.IterationEntity, error) {
			return nil, pluginsdk.ErrNotFound
		}

		err := svc.CanStartIteration(ctx, iter, getCurrentIter)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "planned status")
	})

	t.Run("invalid - another iteration is current", func(t *testing.T) {
		iter := createTestIteration(2, "Test", string(entities.IterationStatusPlanned))

		// Callback that returns current iteration #1
		getCurrentIter := func(ctx context.Context) (*entities.IterationEntity, error) {
			return createTestIteration(1, "Current", string(entities.IterationStatusCurrent)), nil
		}

		err := svc.CanStartIteration(ctx, iter, getCurrentIter)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already current")
	})

	t.Run("valid - same iteration already current (idempotent)", func(t *testing.T) {
		iter := createTestIteration(1, "Test", string(entities.IterationStatusPlanned))

		// Callback that returns the same iteration as current
		getCurrentIter := func(ctx context.Context) (*entities.IterationEntity, error) {
			return createTestIteration(1, "Test", string(entities.IterationStatusCurrent)), nil
		}

		err := svc.CanStartIteration(ctx, iter, getCurrentIter)
		assert.NoError(t, err)
	})

	t.Run("error - callback returns error", func(t *testing.T) {
		iter := createTestIteration(1, "Test", string(entities.IterationStatusPlanned))

		// Callback that returns an error
		getCurrentIter := func(ctx context.Context) (*entities.IterationEntity, error) {
			return nil, pluginsdk.ErrInternal
		}

		err := svc.CanStartIteration(ctx, iter, getCurrentIter)
		assert.Error(t, err)
	})
}

// Helper function to create test iterations
func createTestIteration(number int, name, status string) *entities.IterationEntity {
	now := time.Now()
	iter := &entities.IterationEntity{
		Number:    number,
		Name:      name,
		Goal:      "Test goal",
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return iter
}
