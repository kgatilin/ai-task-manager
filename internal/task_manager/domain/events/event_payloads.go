package events

import (
	"github.com/kgatilin/ai-task-manager/internal/task_manager/domain/entities"
)

// Event payload types for domain events
// Payloads are the actual entity data emitted with events

// Roadmap event payloads
type (
	// RoadmapCreatedPayload contains the created roadmap entity
	RoadmapCreatedPayload = entities.RoadmapEntity

	// RoadmapUpdatedPayload contains the updated roadmap entity
	RoadmapUpdatedPayload = entities.RoadmapEntity
)

// Track event payloads
type (
	// TrackCreatedPayload contains the created track entity
	TrackCreatedPayload = entities.TrackEntity

	// TrackUpdatedPayload contains the updated track entity
	TrackUpdatedPayload = entities.TrackEntity

	// TrackStatusChangedPayload contains the track with changed status
	TrackStatusChangedPayload = entities.TrackEntity

	// TrackCompletedPayload contains the completed track entity
	TrackCompletedPayload = entities.TrackEntity

	// TrackBlockedPayload contains the blocked track entity
	TrackBlockedPayload = entities.TrackEntity
)

// Task event payloads
type (
	// TaskCreatedPayload contains the created task entity
	TaskCreatedPayload = entities.TaskEntity

	// TaskUpdatedPayload contains the updated task entity
	TaskUpdatedPayload = entities.TaskEntity

	// TaskStatusChangedPayload contains the task with changed status
	TaskStatusChangedPayload = entities.TaskEntity

	// TaskCompletedPayload contains the completed task entity
	TaskCompletedPayload = entities.TaskEntity
)

// Iteration event payloads
type (
	// IterationCreatedPayload contains the created iteration entity
	IterationCreatedPayload = entities.IterationEntity

	// IterationStartedPayload contains the started iteration entity
	IterationStartedPayload = entities.IterationEntity

	// IterationCompletedPayload contains the completed iteration entity
	IterationCompletedPayload = entities.IterationEntity

	// IterationUpdatedPayload contains the updated iteration entity
	IterationUpdatedPayload = entities.IterationEntity
)

// Acceptance Criteria event payloads
type (
	// ACCreatedPayload contains the created acceptance criteria entity
	ACCreatedPayload = entities.AcceptanceCriteriaEntity

	// ACUpdatedPayload contains the updated acceptance criteria entity
	ACUpdatedPayload = entities.AcceptanceCriteriaEntity

	// ACVerifiedPayload contains the verified acceptance criteria entity
	ACVerifiedPayload = entities.AcceptanceCriteriaEntity

	// ACAutomaticallyVerifiedPayload contains the automatically verified acceptance criteria entity
	ACAutomaticallyVerifiedPayload = entities.AcceptanceCriteriaEntity

	// ACPendingReviewPayload contains the acceptance criteria awaiting review
	ACPendingReviewPayload = entities.AcceptanceCriteriaEntity

	// ACFailedPayload contains the failed acceptance criteria entity
	ACFailedPayload = entities.AcceptanceCriteriaEntity

	// ACDeletedPayload contains the deleted acceptance criteria entity
	ACDeletedPayload = entities.AcceptanceCriteriaEntity
)

// ADR event payloads
type (
	// ADRCreatedPayload contains the created ADR entity
	ADRCreatedPayload = entities.ADREntity

	// ADRUpdatedPayload contains the updated ADR entity
	ADRUpdatedPayload = entities.ADREntity

	// ADRSupersededPayload contains the superseded ADR entity
	ADRSupersededPayload = entities.ADREntity

	// ADRDeprecatedPayload contains the deprecated ADR entity
	ADRDeprecatedPayload = entities.ADREntity
)
