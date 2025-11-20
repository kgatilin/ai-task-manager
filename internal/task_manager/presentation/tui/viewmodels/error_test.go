package viewmodels_test

import (
	"testing"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/viewmodels"
)

func TestNewErrorViewModel(t *testing.T) {
	errorMsg := "Something went wrong"
	vm := viewmodels.NewErrorViewModel(errorMsg)

	if vm.ErrorMessage != errorMsg {
		t.Errorf("expected error message %q, got %q", errorMsg, vm.ErrorMessage)
	}

	if vm.Details != "" {
		t.Errorf("expected empty details, got %q", vm.Details)
	}

	if vm.CanGoBack {
		t.Error("expected CanGoBack to be false by default")
	}

	if vm.RetryAction != "" {
		t.Errorf("expected empty retry action, got %q", vm.RetryAction)
	}
}

func TestErrorViewModel_Fields(t *testing.T) {
	vm := viewmodels.NewErrorViewModel("Error")
	vm.Details = "More info"
	vm.CanGoBack = true
	vm.RetryAction = "Try again"

	if vm.ErrorMessage != "Error" {
		t.Errorf("expected error message %q, got %q", "Error", vm.ErrorMessage)
	}

	if vm.Details != "More info" {
		t.Errorf("expected details %q, got %q", "More info", vm.Details)
	}

	if !vm.CanGoBack {
		t.Error("expected CanGoBack to be true")
	}

	if vm.RetryAction != "Try again" {
		t.Errorf("expected retry action %q, got %q", "Try again", vm.RetryAction)
	}
}

func TestErrorViewModel_EmptyValues(t *testing.T) {
	vm := viewmodels.NewErrorViewModel("")
	vm.Details = ""
	vm.RetryAction = ""

	if vm.ErrorMessage != "" {
		t.Errorf("expected empty error message, got %q", vm.ErrorMessage)
	}

	if vm.Details != "" {
		t.Errorf("expected empty details, got %q", vm.Details)
	}

	if vm.RetryAction != "" {
		t.Errorf("expected empty retry action, got %q", vm.RetryAction)
	}
}
