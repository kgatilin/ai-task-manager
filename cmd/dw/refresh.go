package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kgatilin/darwinflow-pub/internal/app"
	"github.com/kgatilin/darwinflow-pub/internal/infra"
)

// handleRefresh updates DarwinFlow to the latest version
// This includes:
// - Updating database schema (adding new columns, indexes, etc.)
// - Reinstalling/updating hooks
// - Updating configuration if needed
func handleRefresh(args []string) {
	dbPath := app.DefaultDBPath

	fmt.Println("Refreshing DarwinFlow to latest version...")
	fmt.Println()

	// Step 1: Update database schema
	fmt.Println("Updating database schema...")
	repository, err := infra.NewSQLiteEventRepository(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating repository: %v\n", err)
		os.Exit(1)
	}
	defer repository.Close()

	ctx := context.Background()
	if err := repository.Initialize(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating database schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Database schema updated:", dbPath)

	// Step 2: Update hooks
	fmt.Println()
	fmt.Println("Updating Claude Code hooks...")
	hookConfigManager, err := infra.NewHookConfigManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating hook config manager: %v\n", err)
		os.Exit(1)
	}

	if err := hookConfigManager.InstallDarwinFlowHooks(); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating hooks: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Hooks updated:", hookConfigManager.GetSettingsPath())

	// Step 3: Initialize config if needed
	fmt.Println()
	fmt.Println("Checking configuration...")
	logger := infra.NewDefaultLogger()
	configLoader := infra.NewConfigLoader(logger)

	// Try to load config
	config, err := configLoader.LoadConfig("")
	if err != nil || config == nil {
		// Config doesn't exist or is invalid, create default
		fmt.Println("Creating default configuration...")
		if err := configLoader.InitializeDefaultConfig(""); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not create default config: %v\n", err)
		} else {
			fmt.Println("✓ Configuration initialized:", infra.DefaultConfigFileName)
		}
	} else {
		fmt.Println("✓ Configuration is valid")
	}

	// Done
	fmt.Println()
	fmt.Println("DarwinFlow has been refreshed successfully!")
	fmt.Println()
	fmt.Println("Changes applied:")
	fmt.Println("  • Database schema updated with latest migrations")
	fmt.Println("  • Hooks updated to latest version")
	fmt.Println("  • Configuration verified")
	fmt.Println()
	fmt.Println("You may need to restart Claude Code for changes to take effect.")
}
