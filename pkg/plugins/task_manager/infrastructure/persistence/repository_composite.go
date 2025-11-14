package persistence

import (
	"context"
	"database/sql"

	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/repositories"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// Compile-time interface check
var _ domain.RoadmapRepository = (*SQLiteRepositoryComposite)(nil)

// SQLiteRepositoryComposite implements domain.RoadmapRepository by delegating to focused repositories.
// This provides backward compatibility during the migration from the old monolithic repository
// to the new focused repository architecture.
type SQLiteRepositoryComposite struct {
	Roadmap   repositories.RoadmapRepository
	Track     repositories.TrackRepository
	Task      repositories.TaskRepository
	Iteration repositories.IterationRepository
	ADR       repositories.ADRRepository
	AC        repositories.AcceptanceCriteriaRepository
	Aggregate repositories.AggregateRepository

	DB     *sql.DB
	logger pluginsdk.Logger
}

// NewSQLiteRepositoryComposite creates a composite repository with all focused repositories.
func NewSQLiteRepositoryComposite(db *sql.DB, logger pluginsdk.Logger) *SQLiteRepositoryComposite {
	return &SQLiteRepositoryComposite{
		Roadmap:   NewSQLiteRoadmapOnlyRepository(db, logger),
		Track:     NewSQLiteTrackRepository(db, logger),
		Task:      NewSQLiteTaskRepository(db, logger),
		Iteration: NewSQLiteIterationRepository(db, logger),
		ADR:       NewSQLiteADRRepository(db, logger),
		AC:        NewSQLiteAcceptanceCriteriaRepository(db, logger),
		Aggregate: NewSQLiteAggregateRepository(db, logger),
		DB:        db,
		logger:    logger,
	}
}

// GetDB returns the database connection (for migration command).
func (c *SQLiteRepositoryComposite) GetDB() *sql.DB {
	return c.DB
}

// ============================================================================
// Roadmap operations (4 methods) - delegate to Roadmap repository
// ============================================================================

// SaveRoadmap persists a new roadmap to storage.
func (c *SQLiteRepositoryComposite) SaveRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error {
	return c.Roadmap.SaveRoadmap(ctx, roadmap)
}

// GetRoadmap retrieves a roadmap by its ID.
func (c *SQLiteRepositoryComposite) GetRoadmap(ctx context.Context, id string) (*entities.RoadmapEntity, error) {
	return c.Roadmap.GetRoadmap(ctx, id)
}

// GetActiveRoadmap retrieves the most recently created roadmap.
func (c *SQLiteRepositoryComposite) GetActiveRoadmap(ctx context.Context) (*entities.RoadmapEntity, error) {
	return c.Roadmap.GetActiveRoadmap(ctx)
}

// UpdateRoadmap updates an existing roadmap.
func (c *SQLiteRepositoryComposite) UpdateRoadmap(ctx context.Context, roadmap *entities.RoadmapEntity) error {
	return c.Roadmap.UpdateRoadmap(ctx, roadmap)
}

// ============================================================================
// Track operations (10 methods) - delegate to Track repository
// ============================================================================

// SaveTrack persists a new track to storage.
func (c *SQLiteRepositoryComposite) SaveTrack(ctx context.Context, track *entities.TrackEntity) error {
	return c.Track.SaveTrack(ctx, track)
}

// GetTrack retrieves a track by its ID.
func (c *SQLiteRepositoryComposite) GetTrack(ctx context.Context, id string) (*entities.TrackEntity, error) {
	return c.Track.GetTrack(ctx, id)
}

// ListTracks returns all tracks for a roadmap, optionally filtered.
func (c *SQLiteRepositoryComposite) ListTracks(ctx context.Context, roadmapID string, filters entities.TrackFilters) ([]*entities.TrackEntity, error) {
	return c.Track.ListTracks(ctx, roadmapID, filters)
}

// UpdateTrack updates an existing track.
func (c *SQLiteRepositoryComposite) UpdateTrack(ctx context.Context, track *entities.TrackEntity) error {
	return c.Track.UpdateTrack(ctx, track)
}

// DeleteTrack removes a track from storage.
func (c *SQLiteRepositoryComposite) DeleteTrack(ctx context.Context, id string) error {
	return c.Track.DeleteTrack(ctx, id)
}

// AddTrackDependency adds a dependency from trackID to dependsOnID.
func (c *SQLiteRepositoryComposite) AddTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	return c.Track.AddTrackDependency(ctx, trackID, dependsOnID)
}

// RemoveTrackDependency removes a dependency from trackID to dependsOnID.
func (c *SQLiteRepositoryComposite) RemoveTrackDependency(ctx context.Context, trackID, dependsOnID string) error {
	return c.Track.RemoveTrackDependency(ctx, trackID, dependsOnID)
}

// GetTrackDependencies returns the IDs of all tracks that trackID depends on.
func (c *SQLiteRepositoryComposite) GetTrackDependencies(ctx context.Context, trackID string) ([]string, error) {
	return c.Track.GetTrackDependencies(ctx, trackID)
}

// ValidateNoCycles checks if adding/updating the track would create a circular dependency.
func (c *SQLiteRepositoryComposite) ValidateNoCycles(ctx context.Context, trackID string) error {
	return c.Track.ValidateNoCycles(ctx, trackID)
}

// GetTrackWithTasks retrieves a track with all its tasks.
func (c *SQLiteRepositoryComposite) GetTrackWithTasks(ctx context.Context, trackID string) (*entities.TrackEntity, error) {
	return c.Track.GetTrackWithTasks(ctx, trackID)
}

// ============================================================================
// Task operations (7 methods) - delegate to Task repository
// ============================================================================

// SaveTask persists a new task to storage.
func (c *SQLiteRepositoryComposite) SaveTask(ctx context.Context, task *entities.TaskEntity) error {
	return c.Task.SaveTask(ctx, task)
}

// GetTask retrieves a task by its ID.
func (c *SQLiteRepositoryComposite) GetTask(ctx context.Context, id string) (*entities.TaskEntity, error) {
	return c.Task.GetTask(ctx, id)
}

// ListTasks returns all tasks matching the filters.
func (c *SQLiteRepositoryComposite) ListTasks(ctx context.Context, filters entities.TaskFilters) ([]*entities.TaskEntity, error) {
	return c.Task.ListTasks(ctx, filters)
}

// UpdateTask updates an existing task.
func (c *SQLiteRepositoryComposite) UpdateTask(ctx context.Context, task *entities.TaskEntity) error {
	return c.Task.UpdateTask(ctx, task)
}

// DeleteTask removes a task from storage.
func (c *SQLiteRepositoryComposite) DeleteTask(ctx context.Context, id string) error {
	return c.Task.DeleteTask(ctx, id)
}

// MoveTaskToTrack moves a task from its current track to a new track.
func (c *SQLiteRepositoryComposite) MoveTaskToTrack(ctx context.Context, taskID, newTrackID string) error {
	return c.Task.MoveTaskToTrack(ctx, taskID, newTrackID)
}

// GetBacklogTasks returns all tasks that are not in any iteration and not done.
func (c *SQLiteRepositoryComposite) GetBacklogTasks(ctx context.Context) ([]*entities.TaskEntity, error) {
	return c.Task.GetBacklogTasks(ctx)
}

// ============================================================================
// Iteration operations (13 methods) - delegate to Iteration repository
// ============================================================================

// SaveIteration persists a new iteration to storage.
func (c *SQLiteRepositoryComposite) SaveIteration(ctx context.Context, iteration *entities.IterationEntity) error {
	return c.Iteration.SaveIteration(ctx, iteration)
}

// GetIteration retrieves an iteration by its number.
func (c *SQLiteRepositoryComposite) GetIteration(ctx context.Context, number int) (*entities.IterationEntity, error) {
	return c.Iteration.GetIteration(ctx, number)
}

// GetCurrentIteration returns the iteration with status "current".
func (c *SQLiteRepositoryComposite) GetCurrentIteration(ctx context.Context) (*entities.IterationEntity, error) {
	return c.Iteration.GetCurrentIteration(ctx)
}

// ListIterations returns all iterations, ordered by number.
func (c *SQLiteRepositoryComposite) ListIterations(ctx context.Context) ([]*entities.IterationEntity, error) {
	return c.Iteration.ListIterations(ctx)
}

// UpdateIteration updates an existing iteration.
func (c *SQLiteRepositoryComposite) UpdateIteration(ctx context.Context, iteration *entities.IterationEntity) error {
	return c.Iteration.UpdateIteration(ctx, iteration)
}

// DeleteIteration removes an iteration from storage.
func (c *SQLiteRepositoryComposite) DeleteIteration(ctx context.Context, number int) error {
	return c.Iteration.DeleteIteration(ctx, number)
}

// AddTaskToIteration adds a task to an iteration.
func (c *SQLiteRepositoryComposite) AddTaskToIteration(ctx context.Context, iterationNum int, taskID string) error {
	return c.Iteration.AddTaskToIteration(ctx, iterationNum, taskID)
}

// RemoveTaskFromIteration removes a task from an iteration.
func (c *SQLiteRepositoryComposite) RemoveTaskFromIteration(ctx context.Context, iterationNum int, taskID string) error {
	return c.Iteration.RemoveTaskFromIteration(ctx, iterationNum, taskID)
}

// GetIterationTasks returns all tasks in an iteration.
func (c *SQLiteRepositoryComposite) GetIterationTasks(ctx context.Context, iterationNum int) ([]*entities.TaskEntity, error) {
	return c.Iteration.GetIterationTasks(ctx, iterationNum)
}

// GetIterationTasksWithWarnings retrieves all tasks for an iteration,
// gracefully handling missing tasks by returning them separately.
func (c *SQLiteRepositoryComposite) GetIterationTasksWithWarnings(ctx context.Context, iterationNum int) ([]*entities.TaskEntity, []string, error) {
	return c.Iteration.GetIterationTasksWithWarnings(ctx, iterationNum)
}

// StartIteration marks an iteration as current and sets started_at timestamp.
func (c *SQLiteRepositoryComposite) StartIteration(ctx context.Context, iterationNum int) error {
	return c.Iteration.StartIteration(ctx, iterationNum)
}

// CompleteIteration marks an iteration as complete and sets completed_at timestamp.
func (c *SQLiteRepositoryComposite) CompleteIteration(ctx context.Context, iterationNum int) error {
	return c.Iteration.CompleteIteration(ctx, iterationNum)
}

// GetIterationByNumber is an alias for GetIteration for consistency with other repositories.
func (c *SQLiteRepositoryComposite) GetIterationByNumber(ctx context.Context, number int) (*entities.IterationEntity, error) {
	return c.Iteration.GetIterationByNumber(ctx, number)
}

// GetIterationsForTask returns all iterations that contain a specific task.
func (c *SQLiteRepositoryComposite) GetIterationsForTask(ctx context.Context, taskID string) ([]*entities.IterationEntity, error) {
	return c.Task.GetIterationsForTask(ctx, taskID)
}

// ============================================================================
// ADR operations (7 methods) - delegate to ADR repository
// ============================================================================

// SaveADR persists a new ADR to storage.
func (c *SQLiteRepositoryComposite) SaveADR(ctx context.Context, adr *entities.ADREntity) error {
	return c.ADR.SaveADR(ctx, adr)
}

// GetADR retrieves an ADR by its ID.
func (c *SQLiteRepositoryComposite) GetADR(ctx context.Context, id string) (*entities.ADREntity, error) {
	return c.ADR.GetADR(ctx, id)
}

// ListADRs returns all ADRs, optionally filtered by track.
func (c *SQLiteRepositoryComposite) ListADRs(ctx context.Context, trackID *string) ([]*entities.ADREntity, error) {
	return c.ADR.ListADRs(ctx, trackID)
}

// UpdateADR updates an existing ADR.
func (c *SQLiteRepositoryComposite) UpdateADR(ctx context.Context, adr *entities.ADREntity) error {
	return c.ADR.UpdateADR(ctx, adr)
}

// SupersedeADR marks an ADR as superseded by another ADR.
func (c *SQLiteRepositoryComposite) SupersedeADR(ctx context.Context, adrID, supersededByID string) error {
	return c.ADR.SupersedeADR(ctx, adrID, supersededByID)
}

// DeprecateADR marks an ADR as deprecated.
func (c *SQLiteRepositoryComposite) DeprecateADR(ctx context.Context, adrID string) error {
	return c.ADR.DeprecateADR(ctx, adrID)
}

// GetADRsByTrack returns all ADRs for a specific track.
func (c *SQLiteRepositoryComposite) GetADRsByTrack(ctx context.Context, trackID string) ([]*entities.ADREntity, error) {
	return c.ADR.GetADRsByTrack(ctx, trackID)
}

// ============================================================================
// Acceptance Criteria operations (8 methods) - delegate to AC repository
// ============================================================================

// SaveAC persists a new acceptance criterion to storage.
func (c *SQLiteRepositoryComposite) SaveAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error {
	return c.AC.SaveAC(ctx, ac)
}

// GetAC retrieves an acceptance criterion by its ID.
func (c *SQLiteRepositoryComposite) GetAC(ctx context.Context, id string) (*entities.AcceptanceCriteriaEntity, error) {
	return c.AC.GetAC(ctx, id)
}

// ListAC returns all acceptance criteria for a task.
func (c *SQLiteRepositoryComposite) ListAC(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
	return c.AC.ListAC(ctx, taskID)
}

// UpdateAC updates an existing acceptance criterion.
func (c *SQLiteRepositoryComposite) UpdateAC(ctx context.Context, ac *entities.AcceptanceCriteriaEntity) error {
	return c.AC.UpdateAC(ctx, ac)
}

// DeleteAC removes an acceptance criterion from storage.
func (c *SQLiteRepositoryComposite) DeleteAC(ctx context.Context, id string) error {
	return c.AC.DeleteAC(ctx, id)
}

// ListACByTask is an alias for ListAC for consistency with other repositories.
func (c *SQLiteRepositoryComposite) ListACByTask(ctx context.Context, taskID string) ([]*entities.AcceptanceCriteriaEntity, error) {
	return c.AC.ListACByTask(ctx, taskID)
}

// ListACByTrack returns all acceptance criteria for all tasks in a track.
// NOTE: This is a cross-entity query not yet in focused repositories, implemented directly.
func (c *SQLiteRepositoryComposite) ListACByTrack(ctx context.Context, trackID string) ([]*entities.AcceptanceCriteriaEntity, error) {
	rows, err := c.DB.QueryContext(
		ctx,
		`SELECT ac.id, ac.task_id, ac.description, ac.verification_type, ac.status, ac.notes, ac.testing_instructions, ac.created_at, ac.updated_at
		 FROM acceptance_criteria ac
		 JOIN tasks t ON ac.task_id = t.id
		 WHERE t.track_id = ?
		 ORDER BY ac.created_at ASC`,
		trackID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acs []*entities.AcceptanceCriteriaEntity
	for rows.Next() {
		ac := &entities.AcceptanceCriteriaEntity{}
		err := rows.Scan(
			&ac.ID, &ac.TaskID, &ac.Description, &ac.VerificationType,
			&ac.Status, &ac.Notes, &ac.TestingInstructions,
			&ac.CreatedAt, &ac.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		acs = append(acs, ac)
	}

	return acs, rows.Err()
}

// ListACByIteration returns all acceptance criteria for all tasks in an iteration.
func (c *SQLiteRepositoryComposite) ListACByIteration(ctx context.Context, iterationNum int) ([]*entities.AcceptanceCriteriaEntity, error) {
	return c.AC.ListACByIteration(ctx, iterationNum)
}

// ListFailedAC returns all acceptance criteria with status "failed".
func (c *SQLiteRepositoryComposite) ListFailedAC(ctx context.Context, filters entities.ACFilters) ([]*entities.AcceptanceCriteriaEntity, error) {
	return c.AC.ListFailedAC(ctx, filters)
}

// ============================================================================
// Aggregate queries (2 methods) - delegate to Aggregate repository
// ============================================================================

// GetRoadmapWithTracks retrieves a roadmap with all its tracks.
func (c *SQLiteRepositoryComposite) GetRoadmapWithTracks(ctx context.Context, roadmapID string) (*entities.RoadmapEntity, error) {
	return c.Aggregate.GetRoadmapWithTracks(ctx, roadmapID)
}

// ============================================================================
// Project metadata operations (3 methods) - delegate to Aggregate repository
// ============================================================================

// GetProjectMetadata retrieves a metadata value by key.
func (c *SQLiteRepositoryComposite) GetProjectMetadata(ctx context.Context, key string) (string, error) {
	return c.Aggregate.GetProjectMetadata(ctx, key)
}

// SetProjectMetadata sets a metadata value by key.
func (c *SQLiteRepositoryComposite) SetProjectMetadata(ctx context.Context, key, value string) error {
	return c.Aggregate.SetProjectMetadata(ctx, key, value)
}

// GetProjectCode retrieves the project code (e.g., "DW" for darwinflow).
func (c *SQLiteRepositoryComposite) GetProjectCode(ctx context.Context) string {
	return c.Aggregate.GetProjectCode(ctx)
}

// GetNextSequenceNumber retrieves the next sequence number for an entity type.
func (c *SQLiteRepositoryComposite) GetNextSequenceNumber(ctx context.Context, entityType string) (int, error) {
	return c.Aggregate.GetNextSequenceNumber(ctx, entityType)
}

// ============================================================================
// Backward Compatibility
// ============================================================================

// NewSQLiteRoadmapRepository creates a new composite repository.
// Preserved for backward compatibility with existing code and tests.
// Returns the composite repository that implements the full RoadmapRepository interface.
func NewSQLiteRoadmapRepository(db *sql.DB, logger pluginsdk.Logger) domain.RoadmapRepository {
	return NewSQLiteRepositoryComposite(db, logger)
}

