package entities

// TrackStatus represents valid status values for tracks
type TrackStatus string

const (
	TrackStatusNotStarted TrackStatus = "not-started"
	TrackStatusInProgress TrackStatus = "in-progress"
	TrackStatusComplete   TrackStatus = "complete"
	TrackStatusBlocked    TrackStatus = "blocked"
	TrackStatusWaiting    TrackStatus = "waiting"
)

// Valid status values for tracks
var validTrackStatuses = map[string]bool{
	string(TrackStatusNotStarted): true,
	string(TrackStatusInProgress): true,
	string(TrackStatusComplete):   true,
	string(TrackStatusBlocked):    true,
	string(TrackStatusWaiting):    true,
}

// IsValidTrackStatus validates a track status string
func IsValidTrackStatus(status string) bool {
	return validTrackStatuses[status]
}

// TaskStatus represents valid status values for tasks
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusDone       TaskStatus = "done"
)

// Valid status values for tasks
var validTaskStatuses = map[string]bool{
	string(TaskStatusTodo):       true,
	string(TaskStatusInProgress): true,
	string(TaskStatusDone):       true,
}

// IsValidTaskStatus validates a task status string
func IsValidTaskStatus(status string) bool {
	return validTaskStatuses[status]
}

// IterationStatus represents valid status values for iterations
type IterationStatus string

const (
	IterationStatusPlanned  IterationStatus = "planned"
	IterationStatusCurrent  IterationStatus = "current"
	IterationStatusComplete IterationStatus = "complete"
)

// Valid status values for iterations
var validIterationStatuses = map[string]bool{
	string(IterationStatusPlanned):  true,
	string(IterationStatusCurrent):  true,
	string(IterationStatusComplete): true,
}

// IsValidIterationStatus validates an iteration status string
func IsValidIterationStatus(status string) bool {
	return validIterationStatuses[status]
}

// ADRStatus represents the lifecycle status of an ADR
type ADRStatus string

const (
	ADRStatusProposed    ADRStatus = "proposed"
	ADRStatusAccepted    ADRStatus = "accepted"
	ADRStatusDeprecated  ADRStatus = "deprecated"
	ADRStatusSuperseded  ADRStatus = "superseded"
)

// Valid ADR statuses
var validADRStatuses = map[string]bool{
	string(ADRStatusProposed):   true,
	string(ADRStatusAccepted):   true,
	string(ADRStatusDeprecated): true,
	string(ADRStatusSuperseded): true,
}

// IsValidADRStatus validates an ADR status string
func IsValidADRStatus(status string) bool {
	return validADRStatuses[status]
}

// AcceptanceCriteriaStatus represents the current status of an acceptance criterion
type AcceptanceCriteriaStatus string

const (
	// ACStatusNotStarted - AC has not been verified yet
	ACStatusNotStarted AcceptanceCriteriaStatus = "not_started"
	// ACStatusAutomaticallyVerified - AC was verified by automated process
	ACStatusAutomaticallyVerified AcceptanceCriteriaStatus = "automatically_verified"
	// ACStatusPendingHumanReview - AC is awaiting human verification
	ACStatusPendingHumanReview AcceptanceCriteriaStatus = "pending_human_review"
	// ACStatusVerified - AC has been manually verified by human
	ACStatusVerified AcceptanceCriteriaStatus = "verified"
	// ACStatusFailed - AC did not meet verification requirements
	ACStatusFailed AcceptanceCriteriaStatus = "failed"
)

// AcceptanceCriteriaVerificationType indicates who should verify this AC
type AcceptanceCriteriaVerificationType string

const (
	// VerificationTypeManual - Requires manual human verification
	VerificationTypeManual AcceptanceCriteriaVerificationType = "manual"
	// VerificationTypeAutomated - Can be automatically verified by coding agent
	VerificationTypeAutomated AcceptanceCriteriaVerificationType = "automated"
)

// Filter types for queries

// TrackFilters represents filter criteria for track queries
type TrackFilters struct {
	Status   []string // Filter by status values (e.g., "not-started", "in-progress")
	Priority []string // Legacy - not used
}

// TaskFilters represents filter criteria for task queries
type TaskFilters struct {
	TrackID  string   // Filter by parent track ID
	Status   []string // Filter by status values (e.g., "todo", "in-progress", "done")
	Priority []string // Legacy - not used
}

// ACFilters represents filter criteria for acceptance criteria queries
type ACFilters struct {
	IterationNum *int   // Filter by iteration number
	TrackID      string // Filter by track ID (via tasks)
	TaskID       string // Filter by task ID
}
