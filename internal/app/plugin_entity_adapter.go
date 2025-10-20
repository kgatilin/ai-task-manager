package app

import (
	"time"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// entityAdapter adapts pluginsdk.IExtensible to domain.IExtensible
// This allows plugins using the SDK to provide entities that work with internal app code
type entityAdapter struct {
	inner pluginsdk.IExtensible
}

// newEntityAdapter wraps an SDK entity to implement domain interfaces
func newEntityAdapter(sdkEntity pluginsdk.IExtensible) domain.IExtensible {
	return &entityAdapter{inner: sdkEntity}
}

func (e *entityAdapter) GetID() string {
	return e.inner.GetID()
}

func (e *entityAdapter) GetType() string {
	return e.inner.GetType()
}

func (e *entityAdapter) GetCapabilities() []string {
	return e.inner.GetCapabilities()
}

func (e *entityAdapter) GetField(name string) interface{} {
	return e.inner.GetField(name)
}

func (e *entityAdapter) GetAllFields() map[string]interface{} {
	return e.inner.GetAllFields()
}

// Check if entity also implements IHasContext
func (e *entityAdapter) GetContext() *domain.EntityContext {
	hasContext, ok := e.inner.(pluginsdk.IHasContext)
	if !ok {
		return nil
	}

	sdkCtx := hasContext.GetContext()
	if sdkCtx == nil {
		return nil
	}

	// Adapt SDK EntityContext to domain EntityContext
	return &domain.EntityContext{
		RelatedEntities: sdkCtx.RelatedEntities,
		LinkedFiles:     sdkCtx.LinkedFiles,
		RecentActivity:  adaptActivityRecords(sdkCtx.RecentActivity),
		Metadata:        sdkCtx.Metadata,
	}
}

// Check if entity also implements ITrackable
func (e *entityAdapter) GetStatus() string {
	trackable, ok := e.inner.(pluginsdk.ITrackable)
	if !ok {
		return ""
	}
	return trackable.GetStatus()
}

func (e *entityAdapter) GetProgress() float64 {
	trackable, ok := e.inner.(pluginsdk.ITrackable)
	if !ok {
		return 0
	}
	return trackable.GetProgress()
}

func (e *entityAdapter) IsBlocked() bool {
	trackable, ok := e.inner.(pluginsdk.ITrackable)
	if !ok {
		return false
	}
	return trackable.IsBlocked()
}

func (e *entityAdapter) GetBlockReason() string {
	trackable, ok := e.inner.(pluginsdk.ITrackable)
	if !ok {
		return ""
	}
	return trackable.GetBlockReason()
}

// Check if entity also implements ISchedulable
func (e *entityAdapter) GetStartDate() *time.Time {
	schedulable, ok := e.inner.(pluginsdk.ISchedulable)
	if !ok {
		return nil
	}
	return schedulable.GetStartDate()
}

func (e *entityAdapter) GetDueDate() *time.Time {
	schedulable, ok := e.inner.(pluginsdk.ISchedulable)
	if !ok {
		return nil
	}
	return schedulable.GetDueDate()
}

func (e *entityAdapter) IsOverdue() bool {
	schedulable, ok := e.inner.(pluginsdk.ISchedulable)
	if !ok {
		return false
	}
	return schedulable.IsOverdue()
}

// Check if entity also implements IRelatable
func (e *entityAdapter) GetRelated(entityType string) []string {
	relatable, ok := e.inner.(pluginsdk.IRelatable)
	if !ok {
		return nil
	}
	return relatable.GetRelated(entityType)
}

func (e *entityAdapter) GetAllRelations() map[string][]string {
	relatable, ok := e.inner.(pluginsdk.IRelatable)
	if !ok {
		return nil
	}
	return relatable.GetAllRelations()
}

// adaptActivityRecords converts SDK activity records to domain activity records
func adaptActivityRecords(sdkRecords []pluginsdk.ActivityRecord) []domain.ActivityRecord {
	if sdkRecords == nil {
		return nil
	}

	domainRecords := make([]domain.ActivityRecord, len(sdkRecords))
	for i, sdkRec := range sdkRecords {
		domainRecords[i] = domain.ActivityRecord{
			Timestamp: sdkRec.Timestamp,
			Action:    sdkRec.Action,
			Details:   sdkRec.Details,
		}
	}
	return domainRecords
}

// adaptEntities wraps SDK entities in domain adapters
func adaptEntities(sdkEntities []pluginsdk.IExtensible) []domain.IExtensible {
	if sdkEntities == nil {
		return nil
	}

	domainEntities := make([]domain.IExtensible, len(sdkEntities))
	for i, sdkEntity := range sdkEntities {
		domainEntities[i] = newEntityAdapter(sdkEntity)
	}
	return domainEntities
}

// adaptEntityQuery converts domain EntityQuery to SDK EntityQuery
func adaptEntityQuery(domainQuery domain.EntityQuery) pluginsdk.EntityQuery {
	return pluginsdk.EntityQuery{
		EntityType:   domainQuery.EntityType,
		Capabilities: domainQuery.Capabilities,
		Filters:      domainQuery.Filters,
		Limit:        domainQuery.Limit,
		Offset:       domainQuery.Offset,
		OrderBy:      domainQuery.OrderBy,
		OrderDesc:    domainQuery.OrderDesc,
	}
}

// adaptPluginInfo converts SDK PluginInfo to domain PluginInfo
func adaptPluginInfo(sdkInfo pluginsdk.PluginInfo) domain.PluginInfo {
	return domain.PluginInfo{
		Name:        sdkInfo.Name,
		Version:     sdkInfo.Version,
		Description: sdkInfo.Description,
		IsCore:      sdkInfo.IsCore,
	}
}

// adaptEntityTypeInfo converts SDK EntityTypeInfo to domain EntityTypeInfo
func adaptEntityTypeInfo(sdkInfo pluginsdk.EntityTypeInfo) domain.EntityTypeInfo {
	return domain.EntityTypeInfo{
		Type:              sdkInfo.Type,
		DisplayName:       sdkInfo.DisplayName,
		DisplayNamePlural: sdkInfo.DisplayNamePlural,
		Capabilities:      sdkInfo.Capabilities,
		Icon:              sdkInfo.Icon,
	}
}

// adaptEntityTypeInfos converts slice of SDK EntityTypeInfo to domain EntityTypeInfo
func adaptEntityTypeInfos(sdkInfos []pluginsdk.EntityTypeInfo) []domain.EntityTypeInfo {
	if sdkInfos == nil {
		return nil
	}

	domainInfos := make([]domain.EntityTypeInfo, len(sdkInfos))
	for i, sdkInfo := range sdkInfos {
		domainInfos[i] = adaptEntityTypeInfo(sdkInfo)
	}
	return domainInfos
}
