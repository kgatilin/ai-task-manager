package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgatilin/darwinflow-pub/internal/hooks"
)

// SettingsManager handles reading and writing Claude Code settings
type SettingsManager struct {
	settingsPath string
}

// NewSettingsManager creates a new settings manager
func NewSettingsManager() (*SettingsManager, error) {
	// Find Claude Code settings file
	// Priority: .claude/settings.local.json > .claude/settings.json > ~/.claude/settings.json
	settingsPath, err := findSettingsFile()
	if err != nil {
		return nil, err
	}

	return &SettingsManager{
		settingsPath: settingsPath,
	}, nil
}

// findSettingsFile locates the Claude Code settings file
func findSettingsFile() (string, error) {
	// Check local settings first
	localSettings := []string{
		".claude/settings.local.json",
		".claude/settings.json",
	}

	for _, path := range localSettings {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Check home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	globalSettings := filepath.Join(homeDir, ".claude", "settings.json")
	if _, err := os.Stat(globalSettings); err == nil {
		return globalSettings, nil
	}

	// Default to local settings.json if none exist
	return ".claude/settings.json", nil
}

// Settings represents the Claude Code settings structure
type Settings struct {
	Hooks map[string][]hooks.HookMatcher `json:"hooks,omitempty"`
	// Other settings fields can be added here as needed
	Other map[string]interface{} `json:"-"` // For preserving unknown fields
}

// ReadSettings reads the current settings file
func (sm *SettingsManager) ReadSettings() (*Settings, error) {
	// Check if file exists
	data, err := os.ReadFile(sm.settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return empty settings
			return &Settings{
				Hooks: make(map[string][]hooks.HookMatcher),
				Other: make(map[string]interface{}),
			}, nil
		}
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	// Parse JSON, preserving unknown fields
	var rawSettings map[string]interface{}
	if err := json.Unmarshal(data, &rawSettings); err != nil {
		return nil, fmt.Errorf("failed to parse settings JSON: %w", err)
	}

	settings := &Settings{
		Other: make(map[string]interface{}),
	}

	// Extract hooks
	if hooksData, ok := rawSettings["hooks"]; ok {
		hooksJSON, _ := json.Marshal(hooksData)
		if err := json.Unmarshal(hooksJSON, &settings.Hooks); err != nil {
			return nil, fmt.Errorf("failed to parse hooks: %w", err)
		}
		delete(rawSettings, "hooks")
	} else {
		settings.Hooks = make(map[string][]hooks.HookMatcher)
	}

	// Store other fields
	settings.Other = rawSettings

	return settings, nil
}

// WriteSettings writes settings to the file
func (sm *SettingsManager) WriteSettings(settings *Settings) error {
	// Ensure directory exists
	dir := filepath.Dir(sm.settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	// Merge hooks with other settings
	output := make(map[string]interface{})
	for k, v := range settings.Other {
		output[k] = v
	}
	if len(settings.Hooks) > 0 {
		output["hooks"] = settings.Hooks
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write to file
	if err := os.WriteFile(sm.settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// AddDarwinFlowHooks adds DarwinFlow logging hooks to settings
func (sm *SettingsManager) AddDarwinFlowHooks() error {
	// Read existing settings
	settings, err := sm.ReadSettings()
	if err != nil {
		return err
	}

	// Create backup if file exists
	if _, err := os.Stat(sm.settingsPath); err == nil {
		backupPath := sm.settingsPath + ".backup"
		if err := copyFile(sm.settingsPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Merge with default DarwinFlow hooks
	defaultConfig := hooks.DefaultConfig()
	merged := hooks.MergeConfig(
		hooks.HookConfig{Hooks: settings.Hooks},
		defaultConfig,
	)
	settings.Hooks = merged.Hooks

	// Write updated settings
	if err := sm.WriteSettings(settings); err != nil {
		return err
	}

	return nil
}

// GetSettingsPath returns the path to the settings file
func (sm *SettingsManager) GetSettingsPath() string {
	return sm.settingsPath
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
