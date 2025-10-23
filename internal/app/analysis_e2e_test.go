package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/kgatilin/darwinflow-pub/internal/app"
	"github.com/kgatilin/darwinflow-pub/internal/domain"
	"github.com/kgatilin/darwinflow-pub/internal/infra"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/claude_code"
)


// TestAnalysisE2E_FreshDatabase tests the full analysis flow with a fresh database
func TestAnalysisE2E_FreshDatabase(t *testing.T) {
	// 1. Create in-memory SQLite database
	ctx := context.Background()
	repo, err := infra.NewSQLiteEventRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Initialize schema (this is what was missing in the bug!)
	if err := repo.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. Create mock LLM
	mockLLM := &MockLLM{
		Response:   "This session demonstrates effective tool usage patterns.",
		ModelValue: "mock-model",
	}

	// 3. Create logger
	logger := infra.NewDefaultLogger()

	// 4. Create config
	config := domain.DefaultConfig()

	// 5. Create services
	logsService := app.NewLogsService(repo, repo)
	analysisService := app.NewAnalysisService(repo, repo, logsService, mockLLM, logger, config)

	// Set session view factory
	analysisService.SetSessionViewFactory(func(sessionID string, events []pluginsdk.Event) pluginsdk.AnalysisView {
		return claude_code.NewSessionView(sessionID, events)
	})

	// 6. Create test events
	sessionID := "test-session-123"
	events := createTestEvents(sessionID)

	// Save events to database
	for _, event := range events {
		if err := repo.Save(ctx, event); err != nil {
			t.Fatalf("Failed to save event: %v", err)
		}
	}

	// 7. Analyze the session
	analysis, err := analysisService.AnalyzeSession(ctx, sessionID)
	if err != nil {
		t.Fatalf("AnalyzeSession failed: %v", err)
	}

	// 8. Verify analysis was saved
	if analysis == nil {
		t.Fatal("Analysis is nil")
	}
	if analysis.SessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, analysis.SessionID)
	}
	if analysis.AnalysisResult == "" {
		t.Error("Analysis result is empty")
	}
	// Model comes from config, not from LLM mock
	if analysis.ModelUsed == "" {
		t.Error("Model used is empty")
	}

	// 9. Query analysis back from database
	retrieved, err := analysisService.GetAnalysis(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetAnalysis failed: %v", err)
	}

	// 10. Verify retrieved analysis matches
	if retrieved == nil {
		t.Fatal("Retrieved analysis is nil")
	}
	if retrieved.SessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, retrieved.SessionID)
	}
	if retrieved.AnalysisResult != analysis.AnalysisResult {
		t.Errorf("Analysis result mismatch: expected %s, got %s", analysis.AnalysisResult, retrieved.AnalysisResult)
	}

	// 11. Verify LLM was called
	if mockLLM.QueryCalls == 0 {
		t.Error("LLM was not called")
	}

	// 12. Verify analysis was also saved to generic analyses table
	genericAnalyses, err := repo.FindAnalysisByViewID(ctx, sessionID)
	if err != nil {
		t.Fatalf("FindAnalysisByViewID failed: %v", err)
	}
	if len(genericAnalyses) == 0 {
		t.Error("No generic analyses found")
	}
}

// TestAnalysisE2E_MultipleAnalyses tests multiple analyses can coexist
func TestAnalysisE2E_MultipleAnalyses(t *testing.T) {
	// 1. Create in-memory database
	ctx := context.Background()
	repo, err := infra.NewSQLiteEventRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// 2. Initialize schema
	if err := repo.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. Setup services
	mockLLM := &MockLLM{
		Response:   "Analysis for session",
		ModelValue: "mock-model",
	}
	logger := infra.NewDefaultLogger()
	config := domain.DefaultConfig()
	logsService := app.NewLogsService(repo, repo)
	analysisService := app.NewAnalysisService(repo, repo, logsService, mockLLM, logger, config)
	analysisService.SetSessionViewFactory(func(sessionID string, events []pluginsdk.Event) pluginsdk.AnalysisView {
		return claude_code.NewSessionView(sessionID, events)
	})

	// 4. Create multiple sessions and analyze them
	sessionIDs := []string{"session-1", "session-2", "session-3"}
	for _, sessionID := range sessionIDs {
		events := createTestEvents(sessionID)
		for _, event := range events {
			if err := repo.Save(ctx, event); err != nil {
				t.Fatalf("Failed to save event: %v", err)
			}
		}

		_, err = analysisService.AnalyzeSession(ctx, sessionID)
		if err != nil {
			t.Fatalf("AnalyzeSession failed for %s: %v", sessionID, err)
		}
	}

	// 5. Verify all analyses exist in analyses table
	for _, sessionID := range sessionIDs {
		analyses, err := repo.FindAnalysisByViewID(ctx, sessionID)
		if err != nil {
			t.Fatalf("FindAnalysisByViewID failed for %s: %v", sessionID, err)
		}
		if len(analyses) == 0 {
			t.Errorf("No analyses found for session %s", sessionID)
		}
	}

	// 6. Verify ListRecentAnalyses returns all
	recent, err := repo.ListRecentAnalyses(ctx, 10)
	if err != nil {
		t.Fatalf("ListRecentAnalyses failed: %v", err)
	}
	if len(recent) < len(sessionIDs) {
		t.Errorf("Expected at least %d analyses, got %d", len(sessionIDs), len(recent))
	}

	// 7. Verify session_analyses backward compatibility
	for _, sessionID := range sessionIDs {
		analysis, err := repo.GetAnalysisBySessionID(ctx, sessionID)
		if err != nil {
			t.Fatalf("GetAnalysisBySessionID failed for %s: %v", sessionID, err)
		}
		if analysis == nil {
			t.Errorf("No session analysis found for %s", sessionID)
		}
	}
}

// TestAnalysisE2E_SessionView tests the full SessionView flow
func TestAnalysisE2E_SessionView(t *testing.T) {
	// 1. Setup
	ctx := context.Background()
	repo, err := infra.NewSQLiteEventRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	if err := repo.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	mockLLM := &MockLLM{
		Response:   "Comprehensive tool usage analysis",
		ModelValue: "mock-model",
	}
	logger := infra.NewDefaultLogger()
	config := domain.DefaultConfig()

	logsService := app.NewLogsService(repo, repo)
	analysisService := app.NewAnalysisService(repo, repo, logsService, mockLLM, logger, config)
	analysisService.SetSessionViewFactory(func(sessionID string, events []pluginsdk.Event) pluginsdk.AnalysisView {
		return claude_code.NewSessionView(sessionID, events)
	})

	// 2. Create events
	sessionID := "view-test-session"
	events := createTestEvents(sessionID)
	for _, event := range events {
		if err := repo.Save(ctx, event); err != nil {
			t.Fatalf("Failed to save event: %v", err)
		}
	}

	// 3. Test AnalyzeView directly
	// Fetch events
	query := pluginsdk.EventQuery{
		Metadata: map[string]string{"session_id": sessionID},
	}
	domainEvents, err := repo.FindByQuery(ctx, query)
	if err != nil {
		t.Fatalf("Failed to query events: %v", err)
	}

	// Convert to plugin events
	pluginEvents := make([]pluginsdk.Event, len(domainEvents))
	for i, e := range domainEvents {
		pluginEvents[i] = pluginsdk.Event{
			Type:      e.Type,
			Source:    "claude_code",
			Timestamp: e.Timestamp,
			Payload:   make(map[string]interface{}),
			Metadata:  map[string]string{"session_id": e.SessionID},
			Version:   e.Version,
		}
	}

	// Create session view
	sessionView := claude_code.NewSessionView(sessionID, pluginEvents)

	// Analyze view
	analysis, err := analysisService.AnalyzeView(ctx, sessionView, "tool_analysis")
	if err != nil {
		t.Fatalf("AnalyzeView failed: %v", err)
	}

	// 4. Verify analysis
	if analysis == nil {
		t.Fatal("Analysis is nil")
	}
	if analysis.ViewID != sessionID {
		t.Errorf("Expected view ID %s, got %s", sessionID, analysis.ViewID)
	}
	if analysis.ViewType != "session" {
		t.Errorf("Expected view type 'session', got %s", analysis.ViewType)
	}
	if analysis.Result == "" {
		t.Error("Analysis result is empty")
	}

	// 5. Retrieve from database
	analyses, err := repo.FindAnalysisByViewID(ctx, sessionID)
	if err != nil {
		t.Fatalf("FindAnalysisByViewID failed: %v", err)
	}
	if len(analyses) == 0 {
		t.Fatal("No analyses found")
	}
	if analyses[0].ViewID != sessionID {
		t.Errorf("Expected view ID %s, got %s", sessionID, analyses[0].ViewID)
	}
}

// TestAnalysisE2E_GenericAnalysis tests generic Analysis storage and retrieval
func TestAnalysisE2E_GenericAnalysis(t *testing.T) {
	// 1. Setup
	ctx := context.Background()
	repo, err := infra.NewSQLiteEventRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	if err := repo.Initialize(ctx); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. Create generic analysis
	analysis := domain.NewAnalysis(
		"custom-view-123",
		"custom_type",
		"Custom analysis result",
		"claude-sonnet-4",
		"custom_prompt",
	)
	analysis.Metadata = map[string]interface{}{
		"custom_field": "custom_value",
		"count":        42,
	}

	// 3. Save to database
	if err := repo.SaveGenericAnalysis(ctx, analysis); err != nil {
		t.Fatalf("SaveGenericAnalysis failed: %v", err)
	}

	// 4. Retrieve by view ID
	retrieved, err := repo.FindAnalysisByViewID(ctx, "custom-view-123")
	if err != nil {
		t.Fatalf("FindAnalysisByViewID failed: %v", err)
	}
	if len(retrieved) == 0 {
		t.Fatal("No analyses found")
	}

	// 5. Verify metadata is preserved
	if retrieved[0].Metadata == nil {
		t.Fatal("Metadata is nil")
	}
	if val, ok := retrieved[0].Metadata["custom_field"].(string); !ok || val != "custom_value" {
		t.Errorf("Expected custom_field='custom_value', got %v", retrieved[0].Metadata["custom_field"])
	}

	// 6. Retrieve by view type
	byType, err := repo.FindAnalysisByViewType(ctx, "custom_type")
	if err != nil {
		t.Fatalf("FindAnalysisByViewType failed: %v", err)
	}
	if len(byType) == 0 {
		t.Fatal("No analyses found by type")
	}

	// 7. Retrieve by ID
	byID, err := repo.FindAnalysisById(ctx, analysis.ID)
	if err != nil {
		t.Fatalf("FindAnalysisById failed: %v", err)
	}
	if byID == nil {
		t.Fatal("Analysis not found by ID")
	}
	if byID.ViewID != "custom-view-123" {
		t.Errorf("Expected view ID 'custom-view-123', got %s", byID.ViewID)
	}

	// 8. List recent analyses
	recent, err := repo.ListRecentAnalyses(ctx, 10)
	if err != nil {
		t.Fatalf("ListRecentAnalyses failed: %v", err)
	}
	if len(recent) == 0 {
		t.Fatal("No recent analyses found")
	}
}

// Helper: Create test events with unique IDs per session
func createTestEvents(sessionID string) []*domain.Event {
	now := time.Now()
	return []*domain.Event{
		{
			ID:        sessionID + "-event-1",
			Timestamp: now.Add(-5 * time.Minute),
			Type:      "chat_started",
			SessionID: sessionID,
			Payload:   []byte(`{"message":"User started chat"}`),
			Content:   "User started chat session",
			Version:   "1.0",
		},
		{
			ID:        sessionID + "-event-2",
			Timestamp: now.Add(-4 * time.Minute),
			Type:      "tool_invoked",
			SessionID: sessionID,
			Payload:   []byte(`{"tool":"read","file":"test.go"}`),
			Content:   "Read file test.go",
			Version:   "1.0",
		},
		{
			ID:        sessionID + "-event-3",
			Timestamp: now.Add(-3 * time.Minute),
			Type:      "tool_invoked",
			SessionID: sessionID,
			Payload:   []byte(`{"tool":"edit","file":"test.go"}`),
			Content:   "Edit file test.go",
			Version:   "1.0",
		},
		{
			ID:        sessionID + "-event-4",
			Timestamp: now.Add(-2 * time.Minute),
			Type:      "chat_ended",
			SessionID: sessionID,
			Payload:   []byte(`{"message":"Session completed"}`),
			Content:   "Session completed successfully",
			Version:   "1.0",
		},
	}
}
