package main_test

import (
	"strings"
	"testing"

	main "github.com/kgatilin/darwinflow-pub/cmd/dw"
)

func TestPrintLogsHelp(t *testing.T) {
	output := captureStdout(func() {
		main.PrintLogsHelp()
	})

	expectedStrings := []string{
		"DarwinFlow Logs",
		"DATABASE STRUCTURE",
		"Table: events",
		"Columns:",
		"event_type",
		"session_id",
		"COMMON EVENT TYPES",
		"QUERY EXAMPLES",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("PrintLogsHelp output should contain %q, but it doesn't", expected)
		}
	}
}

func TestPrintLogsHelp_ContainsSQLExamples(t *testing.T) {
	output := captureStdout(func() {
		main.PrintLogsHelp()
	})

	// Should have SQL query examples
	if !strings.Contains(output, "SELECT") {
		t.Error("PrintLogsHelp should include SQL query examples with SELECT")
	}
	if !strings.Contains(output, "COUNT") {
		t.Error("PrintLogsHelp should include examples using COUNT")
	}
}

func TestPrintLogsHelp_DescribesEventTypes(t *testing.T) {
	output := captureStdout(func() {
		main.PrintLogsHelp()
	})

	// Should describe common event types
	expectedTypes := []string{"tool.invoked", "chat.message.user", "file.read"}
	for _, eventType := range expectedTypes {
		if !strings.Contains(output, eventType) {
			t.Errorf("PrintLogsHelp should mention event type %q", eventType)
		}
	}
}
