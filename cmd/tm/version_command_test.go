package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()

	if cmd == nil {
		t.Fatal("NewVersionCommand() returned nil")
	}

	if cmd.Use != "version" {
		t.Errorf("Use = %v, want version", cmd.Use)
	}

	// Test execution
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "tm - Task Manager CLI") {
		t.Errorf("output = %v, want to contain 'tm - Task Manager CLI'", output)
	}
}
