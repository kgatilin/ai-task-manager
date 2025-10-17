package main

import (
	"fmt"
	"os"

	"github.com/kgatilin/darwinflow-pub/pkg/claude"
)

func handleClaudeCommand(args []string) {
	if len(args) < 1 {
		printClaudeUsage()
		os.Exit(1)
	}

	subcommand := args[0]

	switch subcommand {
	case "init":
		handleInit(args[1:])
	case "log":
		handleLog(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown claude subcommand: %s\n\n", subcommand)
		printClaudeUsage()
		os.Exit(1)
	}
}

func printClaudeUsage() {
	fmt.Println("Usage: dw claude <subcommand>")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  init              Initialize Claude Code logging infrastructure")
	fmt.Println("  log <event-type>  Log a Claude Code event (reads JSON from stdin)")
	fmt.Println()
}

func handleInit(args []string) {
	dbPath := claude.DefaultDBPath

	fmt.Println("Initializing Claude Code logging for DarwinFlow...")
	fmt.Println()

	// Initialize logging infrastructure
	if err := claude.InitializeLogging(dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get settings path
	settingsManager, err := claude.NewSettingsManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Created logging database:", dbPath)
	fmt.Println("✓ Added hooks to Claude Code settings:", settingsManager.GetSettingsPath())
	fmt.Println()
	fmt.Println("DarwinFlow logging is now active for all Claude Code sessions.")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Restart Claude Code to activate the hooks")
	fmt.Println("  2. Events will be automatically logged to", dbPath)
}

func handleLog(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: event type required")
		fmt.Fprintln(os.Stderr, "Usage: dw claude log <event-type>")
		os.Exit(1)
	}

	eventTypeStr := args[0]

	// Get max param length from environment or use default
	maxParamLength := 30
	if envVal := os.Getenv("DW_MAX_PARAM_LENGTH"); envVal != "" {
		if parsed, err := fmt.Sscanf(envVal, "%d", &maxParamLength); err == nil && parsed == 1 {
			// Successfully parsed
		}
	}

	// Delegate to pkg layer
	_ = claude.LogFromStdin(eventTypeStr, maxParamLength)
}
