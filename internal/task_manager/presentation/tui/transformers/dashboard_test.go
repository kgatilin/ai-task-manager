package transformers_test

import (
	"testing"
	"time"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
	"github.com/kgatilin/ai-task-manager/internal/task_manager/presentation/tui/transformers"
)

func TestFilterActiveIterations(t *testing.T) {
	now := time.Now()

	iterations := []*entities.IterationEntity{
		mustCreateIteration(1, "Iteration 1", "Goal 1", "Deliverable 1", []string{}, "planned", 100, now, now),
		mustCreateIteration(2, "Iteration 2", "Goal 2", "Deliverable 2", []string{}, "current", 200, now, now),
		mustCreateIteration(3, "Iteration 3", "Goal 3", "Deliverable 3", []string{}, "complete", 300, now, now),
		mustCreateIteration(4, "Iteration 4", "Goal 4", "Deliverable 4", []string{}, "planned", 400, now, now),
	}

	active := transformers.FilterActiveIterations(iterations)

	if len(active) != 3 {
		t.Errorf("expected 3 active iterations, got %d", len(active))
	}

	// Check that complete iteration is excluded
	for _, iter := range active {
		if iter.Status == "complete" {
			t.Errorf("expected no complete iterations in active list, found one: %d", iter.Number)
		}
	}

	// Verify specific iterations
	expectedNumbers := map[int]bool{1: true, 2: true, 4: true}
	for _, iter := range active {
		if !expectedNumbers[iter.Number] {
			t.Errorf("unexpected iteration number in active list: %d", iter.Number)
		}
	}
}

func TestFilterActiveIterations_Empty(t *testing.T) {
	active := transformers.FilterActiveIterations([]*entities.IterationEntity{})

	if active == nil {
		t.Error("expected non-nil slice for empty input")
	}

	if len(active) != 0 {
		t.Errorf("expected empty slice, got %d items", len(active))
	}
}

func TestFilterActiveIterations_AllComplete(t *testing.T) {
	now := time.Now()
	iterations := []*entities.IterationEntity{
		mustCreateIteration(1, "Iteration 1", "Goal 1", "Deliverable 1", []string{}, "complete", 100, now, now),
		mustCreateIteration(2, "Iteration 2", "Goal 2", "Deliverable 2", []string{}, "complete", 200, now, now),
	}

	active := transformers.FilterActiveIterations(iterations)

	if active == nil {
		t.Error("expected non-nil slice")
	}

	if len(active) != 0 {
		t.Errorf("expected empty slice when all iterations complete, got %d items", len(active))
	}
}

func TestFilterActiveTracks(t *testing.T) {
	now := time.Now()

	tracks := []*entities.TrackEntity{
		mustCreateTrack("TM-track-1", "roadmap-1", "Track 1", "Description 1", "not-started", 100, []string{}, now, now),
		mustCreateTrack("TM-track-2", "roadmap-1", "Track 2", "Description 2", "in-progress", 200, []string{}, now, now),
		mustCreateTrack("TM-track-3", "roadmap-1", "Track 3", "Description 3", "complete", 300, []string{}, now, now),
		mustCreateTrack("TM-track-4", "roadmap-1", "Track 4", "Description 4", "blocked", 400, []string{}, now, now),
	}

	active := transformers.FilterActiveTracks(tracks)

	if len(active) != 3 {
		t.Errorf("expected 3 active tracks, got %d", len(active))
	}

	// Check that complete track is excluded
	for _, track := range active {
		if track.Status == "complete" {
			t.Errorf("expected no complete tracks in active list, found one: %s", track.ID)
		}
	}

	// Verify specific tracks
	expectedIDs := map[string]bool{"TM-track-1": true, "TM-track-2": true, "TM-track-4": true}
	for _, track := range active {
		if !expectedIDs[track.ID] {
			t.Errorf("unexpected track ID in active list: %s", track.ID)
		}
	}
}

func TestFilterActiveTracks_Empty(t *testing.T) {
	active := transformers.FilterActiveTracks([]*entities.TrackEntity{})

	if active == nil {
		t.Error("expected non-nil slice for empty input")
	}

	if len(active) != 0 {
		t.Errorf("expected empty slice, got %d items", len(active))
	}
}

func TestFilterBacklogTasks(t *testing.T) {
	now := time.Now()

	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "in-progress", 200, "", now, now),
		mustCreateTask("TM-task-3", "TM-track-1", "Task 3", "Description 3", "done", 300, "", now, now),
		mustCreateTask("TM-task-4", "TM-track-1", "Task 4", "Description 4", "review", 400, "", now, now),
	}

	backlog := transformers.FilterBacklogTasks(tasks)

	if len(backlog) != 3 {
		t.Errorf("expected 3 backlog tasks, got %d", len(backlog))
	}

	// Check that done task is excluded
	for _, task := range backlog {
		if task.Status == "done" {
			t.Errorf("expected no done tasks in backlog, found one: %s", task.ID)
		}
	}

	// Verify specific tasks
	expectedIDs := map[string]bool{"TM-task-1": true, "TM-task-2": true, "TM-task-4": true}
	for _, task := range backlog {
		if !expectedIDs[task.ID] {
			t.Errorf("unexpected task ID in backlog: %s", task.ID)
		}
	}
}

func TestFilterBacklogTasks_Empty(t *testing.T) {
	backlog := transformers.FilterBacklogTasks([]*entities.TaskEntity{})

	if backlog == nil {
		t.Error("expected non-nil slice for empty input")
	}

	if len(backlog) != 0 {
		t.Errorf("expected empty slice, got %d items", len(backlog))
	}
}

func TestFilterBacklogTasks_ExcludesCancelled(t *testing.T) {
	now := time.Now()

	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "in-progress", 200, "", now, now),
		mustCreateTask("TM-task-3", "TM-track-1", "Task 3", "Description 3", "done", 300, "", now, now),
		mustCreateTask("TM-task-4", "TM-track-1", "Task 4", "Description 4", "cancelled", 400, "", now, now),
		mustCreateTask("TM-task-5", "TM-track-1", "Task 5", "Description 5", "review", 500, "", now, now),
	}

	backlog := transformers.FilterBacklogTasks(tasks)

	if len(backlog) != 3 {
		t.Errorf("expected 3 backlog tasks, got %d", len(backlog))
	}

	// Check that done task is excluded
	for _, task := range backlog {
		if task.Status == "done" {
			t.Errorf("expected no done tasks in backlog, found one: %s", task.ID)
		}
	}

	// Check that cancelled task is excluded
	for _, task := range backlog {
		if task.Status == "cancelled" {
			t.Errorf("expected no cancelled tasks in backlog, found one: %s", task.ID)
		}
	}

	// Verify specific tasks (todo, in-progress, review)
	expectedIDs := map[string]bool{"TM-task-1": true, "TM-task-2": true, "TM-task-5": true}
	for _, task := range backlog {
		if !expectedIDs[task.ID] {
			t.Errorf("unexpected task ID in backlog: %s", task.ID)
		}
	}
}

func TestTransformToRoadmapListViewModel(t *testing.T) {
	now := time.Now()

	// Create test data
	roadmap := &entities.RoadmapEntity{
		ID:              "roadmap-1",
		Vision:          "Build the future",
		SuccessCriteria: "Ship by Q4",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	iterations := []*entities.IterationEntity{
		mustCreateIteration(1, "Iteration 1", "Goal 1", "Deliverable 1", []string{"TM-task-1", "TM-task-2"}, "planned", 100, now, now),
		mustCreateIteration(2, "Iteration 2", "Goal 2", "Deliverable 2", []string{"TM-task-3"}, "current", 200, now, now),
		mustCreateIteration(3, "Iteration 3", "Goal 3", "Deliverable 3", []string{}, "complete", 300, now, now),
	}

	tracks := []*entities.TrackEntity{
		mustCreateTrack("TM-track-1", "roadmap-1", "Track 1", "Description 1", "in-progress", 100, []string{}, now, now),
		mustCreateTrack("TM-track-2", "roadmap-1", "Track 2", "Description 2", "complete", 200, []string{}, now, now),
	}

	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "in-progress", 200, "", now, now),
		mustCreateTask("TM-task-3", "TM-track-1", "Task 3", "Description 3", "done", 300, "", now, now),
	}

	vm := transformers.TransformToRoadmapListViewModel(roadmap, iterations, tracks, tasks)

	// Verify roadmap vision and success criteria
	if vm.Vision != "Build the future" {
		t.Errorf("expected vision 'Build the future', got %q", vm.Vision)
	}

	if vm.SuccessCriteria != "Ship by Q4" {
		t.Errorf("expected success criteria 'Ship by Q4', got %q", vm.SuccessCriteria)
	}

	// Verify iterations (exclude complete)
	if len(vm.ActiveIterations) != 2 {
		t.Errorf("expected 2 active iterations, got %d", len(vm.ActiveIterations))
	}

	if vm.ActiveIterations[0].Number != 1 {
		t.Errorf("expected first iteration number 1, got %d", vm.ActiveIterations[0].Number)
	}

	if vm.ActiveIterations[0].TaskCount != 2 {
		t.Errorf("expected first iteration task count 2, got %d", vm.ActiveIterations[0].TaskCount)
	}

	// Verify tracks (exclude complete)
	if len(vm.ActiveTracks) != 1 {
		t.Errorf("expected 1 active track, got %d", len(vm.ActiveTracks))
	}

	if vm.ActiveTracks[0].ID != "TM-track-1" {
		t.Errorf("expected track ID 'TM-track-1', got %q", vm.ActiveTracks[0].ID)
	}

	// Track task count should exclude done/cancelled tasks (2 active out of 3 total)
	if vm.ActiveTracks[0].TaskCount != 2 {
		t.Errorf("expected track task count 2 (excluding done/cancelled), got %d", vm.ActiveTracks[0].TaskCount)
	}

	// Verify backlog tasks (exclude done)
	if len(vm.BacklogTasks) != 2 {
		t.Errorf("expected 2 backlog tasks, got %d", len(vm.BacklogTasks))
	}

	if vm.BacklogTasks[0].ID != "TM-task-1" {
		t.Errorf("expected first task ID 'TM-task-1', got %q", vm.BacklogTasks[0].ID)
	}
}

func TestTransformToRoadmapListViewModel_EmptyData(t *testing.T) {
	vm := transformers.TransformToRoadmapListViewModel(
		&entities.RoadmapEntity{
			ID:              "roadmap-1",
			Vision:          "",
			SuccessCriteria: "",
		},
		[]*entities.IterationEntity{},
		[]*entities.TrackEntity{},
		[]*entities.TaskEntity{},
	)

	if vm == nil {
		t.Fatal("expected non-nil view model")
	}

	if len(vm.ActiveIterations) != 0 {
		t.Errorf("expected 0 active iterations, got %d", len(vm.ActiveIterations))
	}

	if len(vm.ActiveTracks) != 0 {
		t.Errorf("expected 0 active tracks, got %d", len(vm.ActiveTracks))
	}

	if len(vm.BacklogTasks) != 0 {
		t.Errorf("expected 0 backlog tasks, got %d", len(vm.BacklogTasks))
	}
}

func TestTransformToRoadmapListViewModel_TaskCountCalculation(t *testing.T) {
	now := time.Now()

	roadmap := &entities.RoadmapEntity{
		ID:              "roadmap-1",
		Vision:          "Test vision",
		SuccessCriteria: "Test criteria",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	tracks := []*entities.TrackEntity{
		mustCreateTrack("TM-track-1", "roadmap-1", "Track 1", "Description 1", "in-progress", 100, []string{}, now, now),
		mustCreateTrack("TM-track-2", "roadmap-1", "Track 2", "Description 2", "not-started", 200, []string{}, now, now),
	}

	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2", "Description 2", "todo", 200, "", now, now),
		mustCreateTask("TM-task-3", "TM-track-1", "Task 3", "Description 3", "todo", 300, "", now, now),
		mustCreateTask("TM-task-4", "TM-track-2", "Task 4", "Description 4", "todo", 400, "", now, now),
	}

	vm := transformers.TransformToRoadmapListViewModel(roadmap, []*entities.IterationEntity{}, tracks, tasks)

	// Track 1 should have 3 tasks (all active)
	if vm.ActiveTracks[0].TaskCount != 3 {
		t.Errorf("expected track 1 task count 3, got %d", vm.ActiveTracks[0].TaskCount)
	}

	// Track 2 should have 1 task (all active)
	if vm.ActiveTracks[1].TaskCount != 1 {
		t.Errorf("expected track 2 task count 1, got %d", vm.ActiveTracks[1].TaskCount)
	}
}

// TestTransformToRoadmapListViewModel_TaskCountExcludesDoneAndCancelled verifies that done/cancelled tasks are excluded from track counts
func TestTransformToRoadmapListViewModel_TaskCountExcludesDoneAndCancelled(t *testing.T) {
	now := time.Now()

	roadmap := &entities.RoadmapEntity{
		ID:              "roadmap-1",
		Vision:          "Test vision",
		SuccessCriteria: "Test criteria",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	tracks := []*entities.TrackEntity{
		mustCreateTrack("TM-track-1", "roadmap-1", "Track 1", "Description 1", "in-progress", 100, []string{}, now, now),
	}

	tasks := []*entities.TaskEntity{
		mustCreateTask("TM-task-1", "TM-track-1", "Task 1 - todo", "Description 1", "todo", 100, "", now, now),
		mustCreateTask("TM-task-2", "TM-track-1", "Task 2 - in-progress", "Description 2", "in-progress", 200, "", now, now),
		mustCreateTask("TM-task-3", "TM-track-1", "Task 3 - done", "Description 3", "done", 300, "", now, now),
		mustCreateTask("TM-task-4", "TM-track-1", "Task 4 - cancelled", "Description 4", "cancelled", 400, "", now, now),
		mustCreateTask("TM-task-5", "TM-track-1", "Task 5 - todo", "Description 5", "todo", 500, "", now, now),
	}

	vm := transformers.TransformToRoadmapListViewModel(roadmap, []*entities.IterationEntity{}, tracks, tasks)

	// Track should count only active tasks (3: todo, in-progress, todo)
	// Done and cancelled tasks should be excluded
	if vm.ActiveTracks[0].TaskCount != 3 {
		t.Errorf("expected track task count 3 (excluding done/cancelled), got %d", vm.ActiveTracks[0].TaskCount)
	}

	// Backlog should have 3 active tasks (exclude done and cancelled)
	if len(vm.BacklogTasks) != 3 {
		t.Errorf("expected 3 backlog tasks (excluding done/cancelled), got %d", len(vm.BacklogTasks))
	}
}

// TestTransformToRoadmapListViewModel_RoadmapNil verifies that nil roadmap is handled gracefully
func TestTransformToRoadmapListViewModel_RoadmapNil(t *testing.T) {
	now := time.Now()

	iterations := []*entities.IterationEntity{
		mustCreateIteration(1, "Iteration 1", "Goal 1", "Deliverable 1", []string{}, "planned", 100, now, now),
	}

	vm := transformers.TransformToRoadmapListViewModel(nil, iterations, []*entities.TrackEntity{}, []*entities.TaskEntity{})

	if vm == nil {
		t.Fatal("expected non-nil view model")
	}

	if vm.Vision != "" {
		t.Errorf("expected empty vision for nil roadmap, got %q", vm.Vision)
	}

	if vm.SuccessCriteria != "" {
		t.Errorf("expected empty success criteria for nil roadmap, got %q", vm.SuccessCriteria)
	}

	if len(vm.ActiveIterations) != 1 {
		t.Errorf("expected 1 iteration, got %d", len(vm.ActiveIterations))
	}
}
