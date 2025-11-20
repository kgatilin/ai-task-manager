package domain

import (
	"context"

	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// RoadmapRepository is the legacy monolithic repository interface.
// It is preserved for backward compatibility with existing code.
// New code should use the focused interfaces in domain/repositories/.
//
// This interface aggregates all repository operations for all entities:
// - Roadmap management
// - Track CRUD and dependencies
// - Task CRUD and queries
// - Iteration management and lifecycle
// - ADR (Architecture Decision Records)
// - Acceptance Criteria
// - Aggregate queries
//
// Implementation: infrastructure/persistence/repository_composite.go
type RoadmapRepository interface {
	// Roadmap operations
	SaveRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error
	GetRoadmap(ctx context.Context, id string) (*entities.RoadmapEntity, error)
	GetActiveRoadmap(ctx context.Context) (*entities.RoadmapEntity, error)
	UpdateRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error

	// Track operations
	SaveTrack(ctx context.Context, track *entities.TrackEntity) error
	GetTrack(ctx context.Context, id string) (*entities.TrackEntity, error)
	ListTracks(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error)
	UpdateTrack(ctx context.Context, track *entities.TrackEntity) error
	DeleteTrack(ctx context.Context, id string) error

	// Track dependency operations
	AddTrackDependency(ctx context.Context, trackID, dependsOnID string) error
	RemoveTrackDependency(ctx context.Context, trackID, dependsOnID string) error
	GetTrackDependencies(ctx context.Context, trackID string) ([]string, error)
	ValidateNoCycles(ctx context.Context, trackID string) error
	GetTrackWithTasks(ctx context.Context, trackID string) (*entities.TrackEntity, error)

	// Task operations
	SaveTask(ctx context.Context, task *entities.TaskEntity) error
	GetTask(ctx context.Context, id string) (*entities.TaskEntity, error)
	ListTasks(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error)
	UpdateTask(ctx context.Context, task *entities.TaskEntity) error
	DeleteTask(ctx context.Context, id string) error
	MoveTaskToTrack(ctx context.Context, taskID, newTrackID string) error
	GetBacklogTasks(ctx context.Context) ([]*entities.TaskEntity, error)
	GetIterationsForTask(ctx context.Context, taskID string) ([]*entities.IterationEntity, error)

	// Iteration operations
	SaveIteration(ctx context.Context, iteration *entities.IterationEntity) error
	GetIteration(ctx context.Context, number int) (*entities.IterationEntity, error)
	GetCurrentIteration(ctx context.Context) (*entities.IterationEntity, error)
	ListIterations(ctx context.Context) ([]*entities.IterationEntity, error)
	UpdateIteration(ctx context.Context, iteration *entities.IterationEntity) error
	DeleteIteration(ctx context.Context, number int) error

	// Iteration-task relationship operations
	AddTaskToIteration(ctx context.Context, iterationNum int, taskID string) error
	RemoveTaskFromIteration(ctx context.Context, iterationNum int, taskID string) error
	GetIterationTasks(ctx context.Context, iterationNum int) ([]*entities.TaskEntity, error)
	GetIterationTasksWithWarnings(ctx context.Context, iterationNum int) ([]*entities.TaskEntity, []string, error)

	// Iteration lifecycle operations
	StartIteration(ctx context.Context, iterationNumber int) error
	CompleteIteration(ctx context.Context, iterationNumber int) error
	RevertIteration(ctx context.Context, iterationNumber int) error
	GetIterationByNumber(ctx context.Context, iterationNumber int) (*entities.IterationEntity, error)

	// ADR operations
	SaveADR(ctx context.Context, adr *entities.ADREntity) error
	GetADR(ctx context.Context, id string) (*entities.ADREntity, error)
	ListADRs(ctx context.Context, trackID *string) ([]*entities.ADREntity, error)
	UpdateADR(ctx context.Context, adr *entities.ADREntity) error
	SupersedeADR(ctx context.Context, adrID, supersededByID string) error
	DeprecateADR(ctx context.Context, adrID string) error
	GetADRsByTrack(ctx context.Context, trackID string) ([]*entities.ADREntity, error)

	// Acceptance Criteria operations
	SaveAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error
	GetAC(ctx context.Context, id string) (*entities.AcceptanceCriteriaEntity, error)
	ListAC(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error)
	UpdateAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error
	DeleteAC(ctx context.Context, id string) error
	ListACByTask(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error)
	ListACByTrack(ctx context.Context, trackID string) ([]*entities.AcceptanceCriteriaEntity, error)
	ListACByIteration(ctx context.Context, iterationNum int) ([]*entities.AcceptanceCriteriaEntity, error)
	ListFailedAC(ctx context.Context, filters entities.ACFilters) ([]*entities.AcceptanceCriteriaEntity, error)

	// Document operations
	SaveDocument(ctx context.Context, doc *entities.DocumentEntity) error
	FindDocumentByID(ctx context.Context, id string) (*entities.DocumentEntity, error)
	FindAllDocuments(ctx context.Context) ([]*entities.DocumentEntity, error)
	FindDocumentsByTrack(ctx context.Context, trackID string) ([]*entities.DocumentEntity, error)
	FindDocumentsByIteration(ctx context.Context, iterationNumber int) ([]*entities.DocumentEntity, error)
	FindDocumentsByType(ctx context.Context, docType entities.DocumentType) ([]*entities.DocumentEntity, error)
	UpdateDocument(ctx context.Context, doc *entities.DocumentEntity) error
	DeleteDocument(ctx context.Context, id string) error

	// Aggregate queries
	GetRoadmapWithTracks(ctx context.Context, roadmapID string) (*entities.RoadmapEntity, error)
	GetProjectMetadata(ctx context.Context, key string) (string, error)
	SetProjectMetadata(ctx context.Context, key, value string) error
	GetProjectCode(ctx context.Context) string
	GetNextSequenceNumber(ctx context.Context, entityType string) (int, error)
}

// RoadmapRepositoryFactory is a function that creates a RoadmapRepository instance.
// Used for dependency injection in commands.
type RoadmapRepositoryFactory func() RoadmapRepository
