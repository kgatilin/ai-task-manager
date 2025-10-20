package app_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kgatilin/darwinflow-pub/internal/app"
	"github.com/kgatilin/darwinflow-pub/internal/domain"
)

// MockTool implements domain.Tool for testing
type MockTool struct {
	name        string
	description string
	usage       string
	executeFunc func(ctx context.Context, args []string, projectCtx *domain.ProjectContext) error
}

func (m *MockTool) GetName() string {
	return m.name
}

func (m *MockTool) GetDescription() string {
	return m.description
}

func (m *MockTool) GetUsage() string {
	return m.usage
}

func (m *MockTool) Execute(ctx context.Context, args []string, projectCtx *domain.ProjectContext) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, args, projectCtx)
	}
	return nil
}

// MockToolProvider is a plugin that provides tools
type MockToolProvider struct {
	*MockPlugin
	tools []domain.Tool
}

func NewMockToolProvider(name string, entityTypes []domain.EntityTypeInfo, tools []domain.Tool) *MockToolProvider {
	return &MockToolProvider{
		MockPlugin: NewMockPlugin(name, entityTypes),
		tools:      tools,
	}
}

func (p *MockToolProvider) GetTools() []domain.Tool {
	return p.tools
}

func TestToolRegistry_NewToolRegistry(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)
	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	if toolRegistry == nil {
		t.Error("Expected non-nil ToolRegistry")
	}
}

func TestToolRegistry_GetTool(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)

	// Create a plugin with tools
	tools := []domain.Tool{
		&MockTool{name: "tool1", description: "First tool"},
		&MockTool{name: "tool2", description: "Second tool"},
	}

	plugin := NewMockToolProvider("test-plugin", []domain.EntityTypeInfo{
		{Type: "task", DisplayName: "Task", Capabilities: []string{"IExtensible"}},
	}, tools)

	err := pluginRegistry.RegisterPlugin(plugin)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	// Test getting existing tool
	tool, err := toolRegistry.GetTool("tool1")
	if err != nil {
		t.Fatalf("Failed to get tool: %v", err)
	}
	if tool.GetName() != "tool1" {
		t.Errorf("Expected tool name 'tool1', got %s", tool.GetName())
	}

	// Test getting second tool
	tool2, err := toolRegistry.GetTool("tool2")
	if err != nil {
		t.Fatalf("Failed to get tool2: %v", err)
	}
	if tool2.GetDescription() != "Second tool" {
		t.Error("Tool2 description mismatch")
	}
}

func TestToolRegistry_GetTool_NotFound(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)
	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	_, err := toolRegistry.GetTool("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent tool, got nil")
	}
}

func TestToolRegistry_GetAllTools(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)

	// Register first plugin with tools
	tools1 := []domain.Tool{
		&MockTool{name: "alpha", description: "Alpha tool"},
		&MockTool{name: "beta", description: "Beta tool"},
	}
	plugin1 := NewMockToolProvider("plugin1", []domain.EntityTypeInfo{
		{Type: "task", DisplayName: "Task", Capabilities: []string{"IExtensible"}},
	}, tools1)
	pluginRegistry.RegisterPlugin(plugin1)

	// Register second plugin with tools
	tools2 := []domain.Tool{
		&MockTool{name: "gamma", description: "Gamma tool"},
	}
	plugin2 := NewMockToolProvider("plugin2", []domain.EntityTypeInfo{
		{Type: "note", DisplayName: "Note", Capabilities: []string{"IExtensible"}},
	}, tools2)
	pluginRegistry.RegisterPlugin(plugin2)

	// Register a plugin without tools (doesn't implement IToolProvider)
	plugin3 := NewMockPlugin("plugin3", []domain.EntityTypeInfo{
		{Type: "other", DisplayName: "Other", Capabilities: []string{"IExtensible"}},
	})
	pluginRegistry.RegisterPlugin(plugin3)

	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	allTools := toolRegistry.GetAllTools()
	if len(allTools) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(allTools))
	}

	// Verify tools are sorted by name (alpha, beta, gamma)
	if allTools[0].GetName() != "alpha" {
		t.Errorf("Expected first tool to be 'alpha', got %s", allTools[0].GetName())
	}
	if allTools[1].GetName() != "beta" {
		t.Errorf("Expected second tool to be 'beta', got %s", allTools[1].GetName())
	}
	if allTools[2].GetName() != "gamma" {
		t.Errorf("Expected third tool to be 'gamma', got %s", allTools[2].GetName())
	}
}

func TestToolRegistry_GetAllTools_Empty(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)
	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	allTools := toolRegistry.GetAllTools()
	if len(allTools) != 0 {
		t.Errorf("Expected 0 tools, got %d", len(allTools))
	}
}

func TestToolRegistry_ExecuteTool(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)

	executed := false
	var receivedArgs []string
	tools := []domain.Tool{
		&MockTool{
			name:        "test-tool",
			description: "Test tool",
			executeFunc: func(ctx context.Context, args []string, projectCtx *domain.ProjectContext) error {
				executed = true
				receivedArgs = args
				return nil
			},
		},
	}

	plugin := NewMockToolProvider("plugin", []domain.EntityTypeInfo{
		{Type: "task", DisplayName: "Task", Capabilities: []string{"IExtensible"}},
	}, tools)
	pluginRegistry.RegisterPlugin(plugin)

	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	ctx := context.Background()
	projectCtx := &domain.ProjectContext{
		CWD: "/test",
	}
	args := []string{"arg1", "arg2"}

	err := toolRegistry.ExecuteTool(ctx, "test-tool", args, projectCtx)
	if err != nil {
		t.Fatalf("ExecuteTool failed: %v", err)
	}

	if !executed {
		t.Error("Tool was not executed")
	}
	if len(receivedArgs) != 2 || receivedArgs[0] != "arg1" {
		t.Error("Tool did not receive correct args")
	}
}

func TestToolRegistry_ExecuteTool_Error(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)

	expectedErr := fmt.Errorf("execution error")
	tools := []domain.Tool{
		&MockTool{
			name: "failing-tool",
			executeFunc: func(ctx context.Context, args []string, projectCtx *domain.ProjectContext) error {
				return expectedErr
			},
		},
	}

	plugin := NewMockToolProvider("plugin", []domain.EntityTypeInfo{
		{Type: "task", DisplayName: "Task", Capabilities: []string{"IExtensible"}},
	}, tools)
	pluginRegistry.RegisterPlugin(plugin)

	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	ctx := context.Background()
	projectCtx := &domain.ProjectContext{}
	err := toolRegistry.ExecuteTool(ctx, "failing-tool", nil, projectCtx)

	if err == nil {
		t.Error("Expected error from ExecuteTool, got nil")
	}
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestToolRegistry_ExecuteTool_NotFound(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)
	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	ctx := context.Background()
	err := toolRegistry.ExecuteTool(ctx, "nonexistent", nil, nil)

	if err == nil {
		t.Error("Expected error for nonexistent tool, got nil")
	}
}

func TestToolRegistry_ExecuteToolWithContext(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)

	var receivedProjectCtx *domain.ProjectContext
	tools := []domain.Tool{
		&MockTool{
			name: "ctx-tool",
			executeFunc: func(ctx context.Context, args []string, projectCtx *domain.ProjectContext) error {
				receivedProjectCtx = projectCtx
				return nil
			},
		},
	}

	plugin := NewMockToolProvider("plugin", []domain.EntityTypeInfo{
		{Type: "task", DisplayName: "Task", Capabilities: []string{"IExtensible"}},
	}, tools)
	pluginRegistry.RegisterPlugin(plugin)

	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	ctx := context.Background()
	config := domain.DefaultConfig()
	cwd := "/test/dir"
	dbPath := "/test/db.db"

	err := toolRegistry.ExecuteToolWithContext(ctx, "ctx-tool", nil, nil, nil, config, cwd, dbPath)
	if err != nil {
		t.Fatalf("ExecuteToolWithContext failed: %v", err)
	}

	if receivedProjectCtx == nil {
		t.Fatal("ProjectContext was not passed to tool")
	}
	if receivedProjectCtx.CWD != cwd {
		t.Errorf("Expected CWD = %s, got %s", cwd, receivedProjectCtx.CWD)
	}
	if receivedProjectCtx.DBPath != dbPath {
		t.Errorf("Expected DBPath = %s, got %s", dbPath, receivedProjectCtx.DBPath)
	}
	if receivedProjectCtx.Config != config {
		t.Error("Config was not passed correctly")
	}
}

func TestToolRegistry_ListTools(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)

	tools := []domain.Tool{
		&MockTool{name: "tool1", description: "First tool"},
		&MockTool{name: "tool2", description: "Second tool"},
	}

	plugin := NewMockToolProvider("plugin", []domain.EntityTypeInfo{
		{Type: "task", DisplayName: "Task", Capabilities: []string{"IExtensible"}},
	}, tools)
	pluginRegistry.RegisterPlugin(plugin)

	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	output := toolRegistry.ListTools()
	if output == "" {
		t.Error("Expected non-empty output")
	}
	if !contains(output, "tool1") || !contains(output, "First tool") {
		t.Error("Output should contain tool1 and its description")
	}
	if !contains(output, "tool2") || !contains(output, "Second tool") {
		t.Error("Output should contain tool2 and its description")
	}
}

func TestToolRegistry_ListTools_Empty(t *testing.T) {
	logger := &app.NoOpLogger{}
	pluginRegistry := app.NewPluginRegistry(logger)
	toolRegistry := app.NewToolRegistry(pluginRegistry, logger)

	output := toolRegistry.ListTools()
	if !contains(output, "No tools available") {
		t.Error("Expected 'No tools available' message")
	}
}
