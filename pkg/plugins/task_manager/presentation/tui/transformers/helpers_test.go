package transformers_test

import (
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
)

// mustCreateIteration is a test helper that creates an IterationEntity, panicking on error
func mustCreateIteration(number int, name, goal, deliverable string, taskIDs []string, status string, rank float64, createdAt, updatedAt time.Time) *entities.IterationEntity {
	iter, err := entities.NewIterationEntity(number, name, goal, deliverable, taskIDs, status, rank, time.Time{}, time.Time{}, createdAt, updatedAt)
	if err != nil {
		panic(err)
	}
	return iter
}

// mustCreateTrack is a test helper that creates a TrackEntity, panicking on error
func mustCreateTrack(id, roadmapID, title, description, status string, rank int, dependencies []string, createdAt, updatedAt time.Time) *entities.TrackEntity {
	track, err := entities.NewTrackEntity(id, roadmapID, title, description, status, rank, dependencies, createdAt, updatedAt)
	if err != nil {
		panic(err)
	}
	return track
}

// mustCreateTask is a test helper that creates a TaskEntity, panicking on error
func mustCreateTask(id, trackID, title, description, status string, rank int, branch string, createdAt, updatedAt time.Time) *entities.TaskEntity {
	task, err := entities.NewTaskEntity(id, trackID, title, description, status, rank, branch, createdAt, updatedAt)
	if err != nil {
		panic(err)
	}
	return task
}
