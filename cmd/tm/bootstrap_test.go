package main

import (
	"path/filepath"
	"testing"
)

func TestBootstrapApp(t *testing.T) {
	// Setup: Create temporary working directory
	tempDir := t.TempDir()

	// Set working directory via env var
	t.Setenv("TM_WORKING_DIR", tempDir)

	// Bootstrap app
	app, err := BootstrapApp()
	if err != nil {
		t.Fatalf("BootstrapApp() failed: %v", err)
	}
	defer app.Close()

	// Verify app is initialized
	if app.Logger == nil {
		t.Error("Logger not initialized")
	}
	if app.WorkingDir != tempDir {
		t.Errorf("WorkingDir = %v, want %v", app.WorkingDir, tempDir)
	}
	if app.ActiveProject != "default" {
		t.Errorf("ActiveProject = %v, want default", app.ActiveProject)
	}
	if app.RepositoryCommon == nil {
		t.Error("RepositoryCommon not initialized")
	}
	if app.TrackService == nil {
		t.Error("TrackService not initialized")
	}
	if app.TaskService == nil {
		t.Error("TaskService not initialized")
	}
	if app.IterationService == nil {
		t.Error("IterationService not initialized")
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath() returned empty string")
	}
	// Should end with .tm/config.yaml
	if !filepath.IsAbs(path) && path != "./.tm/config.yaml" {
		t.Errorf("GetConfigPath() = %v, expected absolute path or ./.tm/config.yaml", path)
	}
}

func TestAppClose(t *testing.T) {
	// Setup: Create temporary working directory
	tempDir := t.TempDir()

	t.Setenv("TM_WORKING_DIR", tempDir)

	// Bootstrap app
	app, err := BootstrapApp()
	if err != nil {
		t.Fatalf("BootstrapApp() failed: %v", err)
	}

	// Close app
	if err := app.Close(); err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Closing again should not fail
	if err := app.Close(); err != nil {
		t.Errorf("Close() second time failed: %v", err)
	}
}
