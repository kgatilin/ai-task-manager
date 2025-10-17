package claude

import (
	"os"
	"path/filepath"
	"strings"
)

// DetectContext determines the current context from environment and working directory
func DetectContext() string {
	// Priority 1: Environment variable
	if ctx := os.Getenv("DW_CONTEXT"); ctx != "" {
		return ctx
	}

	// Priority 2: Parse from current working directory
	cwd, err := os.Getwd()
	if err == nil {
		if ctx := parseContextFromPath(cwd); ctx != "" {
			return ctx
		}
	}

	// Default to unknown
	return "unknown"
}

// parseContextFromPath extracts context information from a file path
func parseContextFromPath(path string) string {
	// Look for .darwinflow directory
	current := path
	for {
		dwPath := filepath.Join(current, ".darwinflow")
		if info, err := os.Stat(dwPath); err == nil && info.IsDir() {
			// Found .darwinflow directory, use the parent directory name as project
			projectName := filepath.Base(current)
			return "project/" + projectName
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached root
			break
		}
		current = parent
	}

	// Fallback: use the last directory component
	projectName := filepath.Base(path)
	if projectName != "" && projectName != "." && projectName != "/" {
		return "project/" + projectName
	}

	return ""
}

// NormalizeContent creates a normalized text representation of event data
func NormalizeContent(eventType, payload string) string {
	// For V1, just combine event type and payload
	// Future versions can implement smarter normalization
	return strings.TrimSpace(eventType + ": " + payload)
}
