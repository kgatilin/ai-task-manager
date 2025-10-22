package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/kgatilin/darwinflow-pub/internal/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// EventDispatcher manages async event streams from multiple plugins.
// It provides a central buffered channel where plugins push events,
// and a background goroutine that consumes and persists them.
//
// This enables real-time event streaming from multiple sources without
// blocking plugin operations.
type EventDispatcher struct {
	eventChan    chan pluginsdk.Event
	eventRepo    domain.EventRepository
	logger       Logger
	emitters     []pluginsdk.IEventEmitter
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.RWMutex
	running      bool
	pluginCtx    pluginsdk.PluginContext
	eventCounter int64 // For metrics/debugging
}

const (
	// EventChannelBuffer is the size of the buffered event channel.
	// This allows plugins to emit bursts of events without blocking.
	EventChannelBuffer = 100
)

// NewEventDispatcher creates a new event dispatcher.
func NewEventDispatcher(
	eventRepo domain.EventRepository,
	logger Logger,
	pluginCtx pluginsdk.PluginContext,
) *EventDispatcher {
	return &EventDispatcher{
		eventChan: make(chan pluginsdk.Event, EventChannelBuffer),
		eventRepo: eventRepo,
		logger:    logger,
		emitters:  make([]pluginsdk.IEventEmitter, 0),
		pluginCtx: pluginCtx,
	}
}

// RegisterEmitter registers a plugin that implements IEventEmitter.
// This should be called during plugin registration, before Start().
func (d *EventDispatcher) RegisterEmitter(emitter pluginsdk.IEventEmitter) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.emitters = append(d.emitters, emitter)
	pluginInfo := emitter.GetInfo()
	d.logger.Debug("Registered event emitter: %s", pluginInfo.Name)
}

// Start begins background event processing.
// It starts the consumer goroutine and tells all registered emitters to begin streaming.
func (d *EventDispatcher) Start(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.running {
		return fmt.Errorf("event dispatcher already running")
	}

	// Create cancellable context
	d.ctx, d.cancel = context.WithCancel(ctx)

	// Start consumer goroutine
	d.wg.Add(1)
	go d.consumeEvents()

	// Start all registered emitters
	for _, emitter := range d.emitters {
		if err := emitter.StartEventStream(d.ctx, d.eventChan); err != nil {
			pluginInfo := emitter.GetInfo()
			d.logger.Warn("Failed to start event stream for %s: %v", pluginInfo.Name, err)
			// Continue with other emitters even if one fails
		} else {
			pluginInfo := emitter.GetInfo()
			d.logger.Info("Started event stream for plugin: %s", pluginInfo.Name)
		}
	}

	d.running = true
	d.logger.Info("EventDispatcher started with %d emitters", len(d.emitters))

	return nil
}

// Stop gracefully shuts down the dispatcher.
// It stops all emitters, drains the event channel, and waits for goroutines to finish.
func (d *EventDispatcher) Stop() error {
	d.mu.Lock()
	if !d.running {
		d.mu.Unlock()
		return fmt.Errorf("event dispatcher not running")
	}

	d.logger.Info("Stopping EventDispatcher...")

	// Stop all emitters first
	for _, emitter := range d.emitters {
		if err := emitter.StopEventStream(); err != nil {
			pluginInfo := emitter.GetInfo()
			d.logger.Warn("Error stopping event stream for %s: %v", pluginInfo.Name, err)
		}
	}

	// Cancel context to signal consumer to stop
	d.cancel()
	d.running = false
	d.mu.Unlock()

	// Wait for consumer goroutine to finish
	d.wg.Wait()

	d.logger.Info("EventDispatcher stopped (processed %d events)", d.eventCounter)

	return nil
}

// consumeEvents is the background goroutine that processes events from the channel.
// It runs until the context is cancelled and the channel is drained.
func (d *EventDispatcher) consumeEvents() {
	defer d.wg.Done()

	d.logger.Debug("Event consumer goroutine started")

	for {
		select {
		case event, ok := <-d.eventChan:
			if !ok {
				// Channel closed, exit
				d.logger.Debug("Event channel closed, consumer exiting")
				return
			}

			// Process the event
			if err := d.processEvent(event); err != nil {
				d.logger.Error("Failed to process event %s from %s: %v",
					event.Type, event.Source, err)
			} else {
				d.eventCounter++
			}

		case <-d.ctx.Done():
			// Context cancelled, drain remaining events then exit
			d.logger.Debug("Context cancelled, draining event channel...")
			d.drainChannel()
			return
		}
	}
}

// processEvent persists a single event to the repository.
func (d *EventDispatcher) processEvent(event pluginsdk.Event) error {
	// Use the plugin context adapter to convert SDK event to domain event
	// This reuses the existing conversion logic in pluginContextAdapter
	return d.pluginCtx.EmitEvent(context.Background(), event)
}

// drainChannel processes any remaining events in the channel before shutdown.
func (d *EventDispatcher) drainChannel() {
	drained := 0
	for {
		select {
		case event, ok := <-d.eventChan:
			if !ok {
				d.logger.Debug("Drained %d events during shutdown", drained)
				return
			}
			if err := d.processEvent(event); err != nil {
				d.logger.Error("Failed to process event during drain: %v", err)
			} else {
				drained++
				d.eventCounter++
			}
		default:
			// Channel empty
			d.logger.Debug("Drained %d events during shutdown", drained)
			return
		}
	}
}

// GetMetrics returns current dispatcher metrics (for monitoring/debugging).
func (d *EventDispatcher) GetMetrics() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]interface{}{
		"running":        d.running,
		"emitter_count":  len(d.emitters),
		"events_handled": d.eventCounter,
		"channel_len":    len(d.eventChan),
		"channel_cap":    cap(d.eventChan),
	}
}

// GetEventChannel returns a read-only channel for subscribing to events.
// This allows multiple consumers to listen to the event stream.
// The returned channel receives copies of all events processed by the dispatcher.
func (d *EventDispatcher) GetEventChannel() <-chan pluginsdk.Event {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.eventChan
}
