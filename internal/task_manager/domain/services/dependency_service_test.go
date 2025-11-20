package services_test

import (
	"context"
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/services"
)

func TestDependencyService_ValidateNoCycles(t *testing.T) {
	ctx := context.Background()
	service := services.NewDependencyService()

	// Test case 1: No cycle - simple chain A -> B -> C
	t.Run("no cycle - chain", func(t *testing.T) {
		getDeps := func(ctx context.Context, trackID string) ([]string, error) {
			switch trackID {
			case "A":
				return []string{"B"}, nil
			case "B":
				return []string{"C"}, nil
			case "C":
				return []string{}, nil
			default:
				return []string{}, nil
			}
		}

		err := service.ValidateNoCycles(ctx, "A", getDeps)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Test case 2: Simple 2-node cycle: A -> B -> A
	t.Run("cycle - 2 nodes", func(t *testing.T) {
		getDeps := func(ctx context.Context, trackID string) ([]string, error) {
			switch trackID {
			case "A":
				return []string{"B"}, nil
			case "B":
				return []string{"A"}, nil
			default:
				return []string{}, nil
			}
		}

		err := service.ValidateNoCycles(ctx, "A", getDeps)
		if err == nil {
			t.Error("expected cycle error, got nil")
		}
	})

	// Test case 3: 3-node cycle: A -> B -> C -> A
	t.Run("cycle - 3 nodes", func(t *testing.T) {
		getDeps := func(ctx context.Context, trackID string) ([]string, error) {
			switch trackID {
			case "A":
				return []string{"B"}, nil
			case "B":
				return []string{"C"}, nil
			case "C":
				return []string{"A"}, nil
			default:
				return []string{}, nil
			}
		}

		err := service.ValidateNoCycles(ctx, "A", getDeps)
		if err == nil {
			t.Error("expected cycle error, got nil")
		}
	})

	// Test case 4: Self-loop: A -> A
	t.Run("self-loop", func(t *testing.T) {
		getDeps := func(ctx context.Context, trackID string) ([]string, error) {
			if trackID == "A" {
				return []string{"A"}, nil
			}
			return []string{}, nil
		}

		err := service.ValidateNoCycles(ctx, "A", getDeps)
		if err == nil {
			t.Error("expected cycle error, got nil")
		}
	})

	// Test case 5: No dependencies
	t.Run("no dependencies", func(t *testing.T) {
		getDeps := func(ctx context.Context, trackID string) ([]string, error) {
			return []string{}, nil
		}

		err := service.ValidateNoCycles(ctx, "A", getDeps)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Test case 6: Complex DAG (no cycles)
	t.Run("complex dag", func(t *testing.T) {
		getDeps := func(ctx context.Context, trackID string) ([]string, error) {
			switch trackID {
			case "A":
				return []string{"B", "C"}, nil
			case "B":
				return []string{"D"}, nil
			case "C":
				return []string{"D"}, nil
			case "D":
				return []string{}, nil
			default:
				return []string{}, nil
			}
		}

		err := service.ValidateNoCycles(ctx, "A", getDeps)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
