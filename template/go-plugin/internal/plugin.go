package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ItemPlugin implements a simple item management plugin.
// It demonstrates the DarwinFlow external plugin protocol using JSON-RPC over stdin/stdout.
type ItemPlugin struct {
	// workingDir is the current working directory provided during initialization
	workingDir string

	// items stores all items in memory (keyed by ID)
	items map[string]*Item

	// eventStreaming indicates whether event emission is active
	eventStreaming bool

	// eventChannel is used for direct event emission (CLI mode)
	eventChannel chan pluginsdk.RPCEvent
}

// NewItemPlugin creates a new ItemPlugin instance.
func NewItemPlugin() *ItemPlugin {
	return &ItemPlugin{
		items:        make(map[string]*Item),
		eventChannel: make(chan pluginsdk.RPCEvent, 100),
	}
}

// AddItem adds an item to the plugin's storage.
// This is a convenience method for initializing sample data.
func (p *ItemPlugin) AddItem(item *Item) {
	p.items[item.ID] = item
}

// Serve runs the JSON-RPC server loop.
// It reads newline-delimited JSON requests from stdin and writes responses to stdout.
// This method blocks until stdin is closed.
func (p *ItemPlugin) Serve() {
	// Create a buffered scanner for reading stdin
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024) // 64KB initial, 1MB max

	// Read and process requests line by line
	for scanner.Scan() {
		var req pluginsdk.RPCRequest
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			// Send parse error response
			p.sendError(req.ID, pluginsdk.RPCErrorParseError, "parse error: "+err.Error())
			continue
		}

		// Dispatch request to appropriate handler
		p.handleRequest(&req)
	}
}

// handleRequest routes an RPC request to the appropriate handler method.
func (p *ItemPlugin) handleRequest(req *pluginsdk.RPCRequest) {
	switch req.Method {
	case pluginsdk.RPCMethodInit:
		p.handleInit(req)
	case pluginsdk.RPCMethodGetInfo:
		p.handleGetInfo(req)
	case pluginsdk.RPCMethodGetCapabilities:
		p.handleGetCapabilities(req)
	case pluginsdk.RPCMethodGetEntityTypes:
		p.handleGetEntityTypes(req)
	case pluginsdk.RPCMethodQueryEntities:
		p.handleQueryEntities(req)
	case pluginsdk.RPCMethodGetEntity:
		p.handleGetEntity(req)
	case pluginsdk.RPCMethodUpdateEntity:
		p.handleUpdateEntity(req)
	case pluginsdk.RPCMethodStartEventStream:
		p.handleStartEventStream(req)
	case pluginsdk.RPCMethodStopEventStream:
		p.handleStopEventStream(req)
	default:
		p.sendError(req.ID, pluginsdk.RPCErrorMethodNotFound, "method not found: "+req.Method)
	}
}

// sendResult sends a successful RPC response.
// If result is nil, an empty result is sent.
func (p *ItemPlugin) sendResult(id interface{}, result interface{}) {
	var resultJSON json.RawMessage
	if result != nil {
		data, err := json.Marshal(result)
		if err != nil {
			p.sendError(id, pluginsdk.RPCErrorInternal, "failed to marshal result: "+err.Error())
			return
		}
		resultJSON = data
	}

	resp := pluginsdk.RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  resultJSON,
	}

	data, _ := json.Marshal(resp)
	fmt.Fprintf(os.Stdout, "%s\n", string(data))
}

// sendError sends an RPC error response.
// Use standard error codes from pluginsdk (e.g., RPCErrorInvalidParams).
func (p *ItemPlugin) sendError(id interface{}, code int, message string) {
	resp := pluginsdk.RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &pluginsdk.RPCError{
			Code:    code,
			Message: message,
		},
	}

	data, _ := json.Marshal(resp)
	fmt.Fprintf(os.Stdout, "%s\n", string(data))
}

// emitEvent sends an event to the main process.
// Events are only sent when event streaming is active.
// In RPC mode, events are written to stdout.
// In CLI mode (when eventChannel is buffered), events are sent to the channel only.
func (p *ItemPlugin) emitEvent(eventType string, payload map[string]interface{}) {
	if !p.eventStreaming {
		return
	}

	event := pluginsdk.RPCEvent{
		Event:     "event",
		Type:      eventType,
		Source:    "myplugin",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload:   payload,
	}

	// If event channel is available (CLI mode), send to channel instead of stdout
	if p.eventChannel != nil && cap(p.eventChannel) > 0 {
		select {
		case p.eventChannel <- event:
		default:
			// Channel full, skip event
		}
		return
	}

	// RPC mode: write to stdout
	data, _ := json.Marshal(event)
	fmt.Fprintf(os.Stdout, "%s\n", string(data))
}

// --- Public methods for CLI access ---

// SetWorkingDir sets the working directory (for CLI testing)
func (p *ItemPlugin) SetWorkingDir(dir string) {
	p.workingDir = dir
}

// StartEventStreaming enables event emission
func (p *ItemPlugin) StartEventStreaming() {
	p.eventStreaming = true
	// Emit initial event
	p.emitEvent("stream.started", map[string]interface{}{
		"item_count": len(p.items),
	})
}

// StopEventStreaming disables event emission
func (p *ItemPlugin) StopEventStreaming() {
	p.eventStreaming = false
}

// GetEventChannel returns the event channel for CLI monitoring
func (p *ItemPlugin) GetEventChannel() <-chan pluginsdk.RPCEvent {
	return p.eventChannel
}

// GetInfo returns plugin information (CLI-friendly)
func (p *ItemPlugin) GetInfo() pluginsdk.PluginInfo {
	return pluginsdk.PluginInfo{
		Name:        "myplugin",
		Version:     "1.0.0",
		Description: "Example external plugin template",
		IsCore:      false,
	}
}

// GetCapabilities returns plugin capabilities (CLI-friendly)
func (p *ItemPlugin) GetCapabilities() []string {
	return []string{
		"IEntityProvider",
		"IEntityUpdater",
		"IEventEmitter",
	}
}

// GetEntityTypes returns entity type metadata (CLI-friendly)
func (p *ItemPlugin) GetEntityTypes() []pluginsdk.EntityTypeInfo {
	return []pluginsdk.EntityTypeInfo{
		{
			Type:              "item",
			DisplayName:       "Item",
			DisplayNamePlural: "Items",
			Capabilities:      []string{},
			Icon:              "📦",
			Description:       "A generic item entity",
		},
	}
}

// QueryEntities queries items based on filters (CLI-friendly)
func (p *ItemPlugin) QueryEntities(query pluginsdk.EntityQuery) []map[string]interface{} {
	// Only handle queries for "item" entity type
	if query.EntityType != "item" {
		return []map[string]interface{}{}
	}

	// Convert all items to maps
	items := make([]map[string]interface{}, 0, len(p.items))
	for _, item := range p.items {
		items = append(items, item.ToMap())
	}

	// Apply pagination limit if specified
	if query.Limit > 0 && len(items) > query.Limit {
		items = items[:query.Limit]
	}

	return items
}

// GetEntity retrieves a specific item by ID (CLI-friendly)
func (p *ItemPlugin) GetEntity(entityID string) (map[string]interface{}, error) {
	item, ok := p.items[entityID]
	if !ok {
		return nil, fmt.Errorf("item not found: %s", entityID)
	}
	return item.ToMap(), nil
}

// UpdateEntity updates an item's fields (CLI-friendly)
func (p *ItemPlugin) UpdateEntity(entityID string, fields map[string]interface{}) (map[string]interface{}, error) {
	item, ok := p.items[entityID]
	if !ok {
		return nil, fmt.Errorf("item not found: %s", entityID)
	}

	// Update fields if provided
	if name, ok := fields["name"].(string); ok {
		item.Name = name
	}
	if description, ok := fields["description"].(string); ok {
		item.Description = description
	}
	if tags, ok := fields["tags"].([]string); ok {
		item.Tags = tags
	} else if tagsInterface, ok := fields["tags"].([]interface{}); ok {
		// Handle []interface{} from JSON unmarshaling
		item.Tags = make([]string, 0, len(tagsInterface))
		for _, tag := range tagsInterface {
			if tagStr, ok := tag.(string); ok {
				item.Tags = append(item.Tags, tagStr)
			}
		}
	}

	// Update timestamp
	item.UpdatedAt = time.Now()

	// Emit update event
	p.emitEvent("item.updated", map[string]interface{}{
		"item_id": item.ID,
		"name":    item.Name,
	})

	return item.ToMap(), nil
}
