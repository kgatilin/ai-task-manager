package main

import (
	"fmt"
	"os"
)

func main() {
	// Bootstrap the application with all dependencies
	app, err := BootstrapApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize task manager: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	// Create and execute root Cobra command with app injected
	rootCmd := NewRootCmd(app)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
