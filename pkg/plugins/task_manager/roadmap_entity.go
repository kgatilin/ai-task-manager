package task_manager

import (
	"fmt"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// RoadmapEntity represents a roadmap and implements SDK capability interfaces.
// It implements the IExtensible interface.
type RoadmapEntity struct {
	ID               string    `json:"id"`
	Vision           string    `json:"vision"`
	SuccessCriteria  string    `json:"success_criteria"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NewRoadmapEntity creates a new roadmap entity
func NewRoadmapEntity(id, vision, successCriteria string, createdAt, updatedAt time.Time) (*RoadmapEntity, error) {
	if vision == "" {
		return nil, fmt.Errorf("%w: vision must be non-empty", pluginsdk.ErrInvalidArgument)
	}
	if successCriteria == "" {
		return nil, fmt.Errorf("%w: success criteria must be non-empty", pluginsdk.ErrInvalidArgument)
	}

	return &RoadmapEntity{
		ID:              id,
		Vision:          vision,
		SuccessCriteria: successCriteria,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

// IExtensible implementation

// GetID returns the unique identifier for this entity
func (r *RoadmapEntity) GetID() string {
	return r.ID
}

// GetType returns the entity type
func (r *RoadmapEntity) GetType() string {
	return "roadmap"
}

// GetCapabilities returns list of capability names this entity supports
func (r *RoadmapEntity) GetCapabilities() []string {
	return []string{"IExtensible"}
}

// GetField retrieves a named field value
func (r *RoadmapEntity) GetField(name string) interface{} {
	fields := r.GetAllFields()
	return fields[name]
}

// GetAllFields returns all fields as a map
func (r *RoadmapEntity) GetAllFields() map[string]interface{} {
	return map[string]interface{}{
		"id":                 r.ID,
		"vision":             r.Vision,
		"success_criteria":   r.SuccessCriteria,
		"created_at":         r.CreatedAt,
		"updated_at":         r.UpdatedAt,
	}
}
