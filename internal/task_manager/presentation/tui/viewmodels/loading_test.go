package viewmodels_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

func TestNewLoadingViewModel(t *testing.T) {
	message := "Loading data..."
	vm := viewmodels.NewLoadingViewModel(message)

	if vm.Message != message {
		t.Errorf("expected message %q, got %q", message, vm.Message)
	}

	if !vm.ShowSpinner {
		t.Error("expected ShowSpinner to be true by default")
	}
}

func TestLoadingViewModel_Fields(t *testing.T) {
	vm := viewmodels.NewLoadingViewModel("test")
	vm.Message = "Updated"
	vm.ShowSpinner = false

	if vm.Message != "Updated" {
		t.Errorf("expected message %q, got %q", "Updated", vm.Message)
	}

	if vm.ShowSpinner {
		t.Error("expected ShowSpinner to be false")
	}
}
