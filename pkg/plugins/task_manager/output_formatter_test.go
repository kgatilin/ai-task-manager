package task_manager_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager"
)

func TestParseOutputFormat(t *testing.T) {
	tests := []struct{
		name     string
		input    string
		expected task_manager.OutputFormat
		wantErr  bool
	}{
		{"Empty string defaults to table", "", task_manager.FormatTable, false},
		{"Table format", "table", task_manager.FormatTable, false},
		{"LLM format", "llm", task_manager.FormatLLM, false},
		{"JSON format", "json", task_manager.FormatJSON, false},
		{"Case insensitive", "TABLE", task_manager.FormatTable, false},
		{"Invalid format", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := task_manager.ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseOutputFormat() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOutputFormatter_OutputJSON(t *testing.T) {
	var buf bytes.Buffer
	formatter := task_manager.NewOutputFormatter(&buf, task_manager.FormatJSON)

	testData := map[string]interface{}{
		"name": "test",
		"value": 42,
	}

	err := formatter.OutputJSON(testData)
	if err != nil {
		t.Fatalf("OutputJSON() error = %v", err)
	}

	// Verify output is valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name=test, got %v", result["name"])
	}
}

func TestContextualHints(t *testing.T) {
	var buf bytes.Buffer
	hints := task_manager.NewContextualHints()

	// Empty hints should output nothing
	hints.Output(&buf)
	if buf.String() != "" {
		t.Errorf("Empty hints should output nothing, got: %s", buf.String())
	}

	// Add hints and verify output
	buf.Reset()
	hints.Add("dw task-manager iteration list")
	hints.Add("dw task-manager task show <id>")
	hints.Output(&buf)

	output := buf.String()
	if !strings.Contains(output, "Next steps:") {
		t.Errorf("Expected 'Next steps:', got: %s", output)
	}
	if !strings.Contains(output, "dw task-manager iteration list") {
		t.Errorf("Expected first hint in output, got: %s", output)
	}
	if !strings.Contains(output, "dw task-manager task show <id>") {
		t.Errorf("Expected second hint in output, got: %s", output)
	}
}
