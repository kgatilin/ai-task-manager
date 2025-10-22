package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	main "github.com/kgatilin/darwinflow-pub/cmd/dw"
)

func TestInitializeApp_Success(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	configPath := "" // Use default config

	// Initialize app
	services, err := main.InitializeApp(dbPath, configPath, false)
	if err != nil {
		t.Fatalf("InitializeApp() failed: %v", err)
	}

	// Verify all services are initialized
	if services == nil {
		t.Fatal("Expected non-nil services")
	}

	if services.PluginRegistry == nil {
		t.Error("PluginRegistry is nil")
	}

	if services.CommandRegistry == nil {
		t.Error("CommandRegistry is nil")
	}

	if services.LogsService == nil {
		t.Error("LogsService is nil")
	}

	if services.AnalysisService == nil {
		t.Error("AnalysisService is nil")
	}

	if services.SetupService == nil {
		t.Error("SetupService is nil")
	}

	if services.ConfigLoader == nil {
		t.Error("ConfigLoader is nil")
	}

	if services.Logger == nil {
		t.Error("Logger is nil")
	}

	if services.EventRepo == nil {
		t.Error("EventRepo is nil")
	}

	if services.DBPath != dbPath {
		t.Errorf("DBPath = %q, want %q", services.DBPath, dbPath)
	}

	if services.WorkingDir == "" {
		t.Error("WorkingDir is empty")
	}
}

func TestInitializeApp_DebugMode(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	configPath := ""

	// Initialize app with debug mode
	services, err := main.InitializeApp(dbPath, configPath, true)
	if err != nil {
		t.Fatalf("InitializeApp() failed: %v", err)
	}

	// Verify services are initialized
	if services == nil {
		t.Fatal("Expected non-nil services")
	}

	if services.Logger == nil {
		t.Error("Logger is nil in debug mode")
	}
}

func TestInitializeApp_InvalidDBPath(t *testing.T) {
	// Try to create database in a non-existent parent directory that we cannot create
	// This simulates permission or filesystem errors
	dbPath := "/root/impossible/path/test.db" // Assuming tests don't run as root
	configPath := ""

	services, err := main.InitializeApp(dbPath, configPath, false)

	// On some systems this might succeed if running as root, so we check both cases
	if err == nil && services == nil {
		t.Error("Expected either error or valid services")
	}

	// If services are returned, they should be valid
	if err == nil && services != nil {
		if services.Logger == nil {
			t.Error("Logger should be initialized even if DB creation fails later")
		}
	}
}

func TestInitializeApp_CreatesDBDirectory(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create a nested path
	dbPath := filepath.Join(tmpDir, "nested", "path", "test.db")
	configPath := ""

	// Initialize app
	services, err := main.InitializeApp(dbPath, configPath, false)
	if err != nil {
		t.Fatalf("InitializeApp() failed: %v", err)
	}

	// Verify directory was created
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		t.Errorf("Database directory was not created: %s", dbDir)
	}

	// Verify services are valid
	if services == nil {
		t.Fatal("Expected non-nil services")
	}
}

func TestInitializeApp_NonFatalConfigError(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Use a non-existent config file (should fall back to defaults)
	configPath := filepath.Join(tmpDir, "nonexistent.yaml")

	// Initialize app - should succeed with default config
	services, err := main.InitializeApp(dbPath, configPath, false)
	if err != nil {
		t.Fatalf("InitializeApp() should handle missing config gracefully: %v", err)
	}

	// Verify services are initialized despite missing config
	if services == nil {
		t.Fatal("Expected non-nil services even with missing config")
	}

	if services.ConfigLoader == nil {
		t.Error("ConfigLoader should be initialized")
	}

	if services.Logger == nil {
		t.Error("Logger should be initialized")
	}
}

func TestInitializeApp_PluginsRegistered(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	configPath := ""

	// Initialize app
	services, err := main.InitializeApp(dbPath, configPath, false)
	if err != nil {
		t.Fatalf("InitializeApp() failed: %v", err)
	}

	// Verify built-in plugins are registered
	if services.PluginRegistry == nil {
		t.Fatal("PluginRegistry is nil")
	}

	pluginInfos := services.PluginRegistry.GetPluginInfos()
	if len(pluginInfos) == 0 {
		t.Error("Expected at least one built-in plugin to be registered")
	}

	// Verify claude-code plugin is registered
	foundClaudeCode := false
	for _, info := range pluginInfos {
		if info.Name == "claude-code" {
			foundClaudeCode = true
			break
		}
	}

	if !foundClaudeCode {
		t.Error("claude-code plugin should be registered by default")
	}
}

func TestInitializeApp_WorkingDirectory(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	configPath := ""

	// Initialize app
	services, err := main.InitializeApp(dbPath, configPath, false)
	if err != nil {
		t.Fatalf("InitializeApp() failed: %v", err)
	}

	// Verify working directory is set
	if services.WorkingDir == "" {
		t.Error("WorkingDir should not be empty")
	}

	// Working directory should exist
	if _, err := os.Stat(services.WorkingDir); err != nil {
		t.Errorf("WorkingDir does not exist: %s", services.WorkingDir)
	}
}

func TestAppServices_Structure(t *testing.T) {
	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	configPath := ""

	// Initialize app
	services, err := main.InitializeApp(dbPath, configPath, false)
	if err != nil {
		t.Fatalf("InitializeApp() failed: %v", err)
	}

	// Test that AppServices struct has all expected fields populated
	tests := []struct {
		name  string
		check func() bool
	}{
		{"PluginRegistry", func() bool { return services.PluginRegistry != nil }},
		{"CommandRegistry", func() bool { return services.CommandRegistry != nil }},
		{"LogsService", func() bool { return services.LogsService != nil }},
		{"AnalysisService", func() bool { return services.AnalysisService != nil }},
		{"SetupService", func() bool { return services.SetupService != nil }},
		{"ConfigLoader", func() bool { return services.ConfigLoader != nil }},
		{"Logger", func() bool { return services.Logger != nil }},
		{"EventRepo", func() bool { return services.EventRepo != nil }},
		{"DBPath", func() bool { return services.DBPath == dbPath }},
		{"WorkingDir", func() bool { return services.WorkingDir != "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Errorf("AppServices.%s check failed", tt.name)
			}
		})
	}
}

// TestInitializeApp_LoadsExternalPlugins tests that external plugins
// are loaded from plugins.yaml during bootstrap.
func TestInitializeApp_LoadsExternalPlugins(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	darwinflowDir := filepath.Join(tmpDir, ".darwinflow")
	logsDir := filepath.Join(darwinflowDir, "logs")
	binDir := filepath.Join(darwinflowDir, "bin")

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("failed to create bin dir: %v", err)
	}

	// Build a simple test plugin
	pluginSrc := `package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Request struct {
	JSONRPC string          ` + "`json:\"jsonrpc\"`" + `
	ID      interface{}     ` + "`json:\"id\"`" + `
	Method  string          ` + "`json:\"method\"`" + `
	Params  json.RawMessage ` + "`json:\"params,omitempty\"`" + `
}

type Response struct {
	JSONRPC string          ` + "`json:\"jsonrpc\"`" + `
	ID      interface{}     ` + "`json:\"id\"`" + `
	Result  json.RawMessage ` + "`json:\"result,omitempty\"`" + `
}

type PluginInfo struct {
	Name         string   ` + "`json:\"name\"`" + `
	Version      string   ` + "`json:\"version\"`" + `
	Description  string   ` + "`json:\"description\"`" + `
	Capabilities []string ` + "`json:\"capabilities\"`" + `
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req Request
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			continue
		}

		var resp Response
		resp.JSONRPC = "2.0"
		resp.ID = req.ID

		switch req.Method {
		case "get_info":
			info := PluginInfo{
				Name:        "bootstrap-test-plugin",
				Version:     "1.0.0",
				Description: "Test plugin for bootstrap integration",
				Capabilities: []string{},
			}
			result, _ := json.Marshal(info)
			resp.Result = result
		case "get_capabilities":
			resp.Result = json.RawMessage("[]")
		default:
			resp.Result = json.RawMessage("{}")
		}

		data, _ := json.Marshal(resp)
		fmt.Fprintf(os.Stdout, "%s\n", string(data))
	}
}
`

	// Write and compile plugin
	pluginSrcPath := filepath.Join(tmpDir, "plugin.go")
	if err := os.WriteFile(pluginSrcPath, []byte(pluginSrc), 0644); err != nil {
		t.Fatalf("failed to write plugin source: %v", err)
	}

	pluginBinPath := filepath.Join(binDir, "bootstrap-test-plugin")
	cmd := exec.Command("go", "build", "-o", pluginBinPath, pluginSrcPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build test plugin: %v\nOutput: %s", err, output)
	}

	// Create plugins.yaml
	pluginsYAML := `plugins:
  bootstrap-test:
    command: bin/bootstrap-test-plugin
    enabled: true
`
	pluginsYAMLPath := filepath.Join(darwinflowDir, "plugins.yaml")
	if err := os.WriteFile(pluginsYAMLPath, []byte(pluginsYAML), 0644); err != nil {
		t.Fatalf("failed to write plugins.yaml: %v", err)
	}

	// Initialize app with the test database path
	dbPath := filepath.Join(logsDir, "events.db")
	services, err := main.InitializeApp(dbPath, tmpDir, false)
	if err != nil {
		t.Fatalf("InitializeApp failed: %v", err)
	}
	// Cast EventRepo to access Close method
	if closer, ok := services.EventRepo.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	// Verify that external plugin was loaded
	plugin, err := services.PluginRegistry.GetPlugin("bootstrap-test-plugin")
	if err != nil {
		t.Fatalf("GetPlugin failed: %v", err)
	}
	if plugin == nil {
		// List all registered plugins for debugging
		allPlugins := services.PluginRegistry.GetAllPlugins()
		t.Logf("Registered plugins: %d", len(allPlugins))
		for _, p := range allPlugins {
			info := p.GetInfo()
			t.Logf("  - %s", info.Name)
		}
		t.Fatal("external plugin 'bootstrap-test-plugin' was not loaded")
	}

	// Verify plugin info
	info := plugin.GetInfo()
	if info.Name != "bootstrap-test-plugin" {
		t.Errorf("expected plugin name 'bootstrap-test-plugin', got %q", info.Name)
	}
	if info.Description != "Test plugin for bootstrap integration" {
		t.Errorf("unexpected description: %q", info.Description)
	}
}

// TestInitializeApp_SkipsDisabledPlugins tests that disabled plugins
// in plugins.yaml are not loaded.
func TestInitializeApp_SkipsDisabledPlugins(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	darwinflowDir := filepath.Join(tmpDir, ".darwinflow")
	logsDir := filepath.Join(darwinflowDir, "logs")

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}

	// Create plugins.yaml with disabled plugin
	pluginsYAML := `plugins:
  disabled-test:
    command: /usr/bin/true
    enabled: false
`
	pluginsYAMLPath := filepath.Join(darwinflowDir, "plugins.yaml")
	if err := os.WriteFile(pluginsYAMLPath, []byte(pluginsYAML), 0644); err != nil {
		t.Fatalf("failed to write plugins.yaml: %v", err)
	}

	// Initialize app
	dbPath := filepath.Join(logsDir, "events.db")
	services, err := main.InitializeApp(dbPath, tmpDir, false)
	if err != nil {
		t.Fatalf("InitializeApp failed: %v", err)
	}
	if closer, ok := services.EventRepo.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	// Verify that disabled plugin was not loaded
	plugin, err := services.PluginRegistry.GetPlugin("disabled-test")
	if err == nil && plugin != nil {
		t.Error("disabled plugin should not have been loaded")
	}

	// Should still have core plugins
	allPlugins := services.PluginRegistry.GetAllPlugins()
	if len(allPlugins) < 2 {
		t.Errorf("expected at least 2 core plugins, got %d", len(allPlugins))
	}
}

// TestInitializeApp_HandlesInvalidPlugins tests that invalid plugins
// in plugins.yaml don't crash the bootstrap.
func TestInitializeApp_HandlesInvalidPlugins(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	darwinflowDir := filepath.Join(tmpDir, ".darwinflow")
	logsDir := filepath.Join(darwinflowDir, "logs")

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}

	// Create plugins.yaml with invalid plugin path
	pluginsYAML := `plugins:
  invalid-test:
    command: /nonexistent/plugin
    enabled: true
`
	pluginsYAMLPath := filepath.Join(darwinflowDir, "plugins.yaml")
	if err := os.WriteFile(pluginsYAMLPath, []byte(pluginsYAML), 0644); err != nil {
		t.Fatalf("failed to write plugins.yaml: %v", err)
	}

	// Initialize app - should not crash
	dbPath := filepath.Join(logsDir, "events.db")
	services, err := main.InitializeApp(dbPath, tmpDir, false)
	if err != nil {
		t.Fatalf("InitializeApp failed: %v", err)
	}
	if closer, ok := services.EventRepo.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	// Verify that invalid plugin was not loaded
	plugin, err := services.PluginRegistry.GetPlugin("invalid-test")
	if err == nil && plugin != nil {
		t.Error("invalid plugin should not have been loaded")
	}

	// Should still have core plugins
	allPlugins := services.PluginRegistry.GetAllPlugins()
	if len(allPlugins) < 2 {
		t.Errorf("expected at least 2 core plugins, got %d", len(allPlugins))
	}
}

// TestInitializeApp_NoPluginsYAML tests that bootstrap works fine
// when plugins.yaml doesn't exist.
func TestInitializeApp_NoPluginsYAML(t *testing.T) {
	// Create temporary directory structure WITHOUT plugins.yaml
	tmpDir := t.TempDir()
	darwinflowDir := filepath.Join(tmpDir, ".darwinflow")
	logsDir := filepath.Join(darwinflowDir, "logs")

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}

	// Initialize app - should not crash
	dbPath := filepath.Join(logsDir, "events.db")
	services, err := main.InitializeApp(dbPath, tmpDir, false)
	if err != nil {
		t.Fatalf("InitializeApp failed: %v", err)
	}
	if closer, ok := services.EventRepo.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	// Should still have core plugins
	allPlugins := services.PluginRegistry.GetAllPlugins()
	if len(allPlugins) < 2 {
		t.Errorf("expected at least 2 core plugins, got %d", len(allPlugins))
	}
}
