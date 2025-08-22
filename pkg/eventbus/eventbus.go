package eventbus

import (
	"context"
	"github.com/alishojaeiir/mahdaad/internal/events"
	"sync"
)

// EventBus is a generic in-memory pub/sub system for EDA.
type EventBus struct {
	subscribers map[string][]chan events.Event
	mu          sync.RWMutex
}

// NewEventBus creates a new EventBus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan events.Event),
	}
}

// Subscribe registers a channel to receive eventbus of a specific type.
func (eb *EventBus) Subscribe(eventType string, ch chan events.Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
}

// Publish sends an event to all subscribed channels asynchronously.
func (eb *EventBus) Publish(ctx context.Context, event events.Event) {
	eb.mu.RLock()
	subs := eb.subscribers[event.Type()]
	eb.mu.RUnlock()

	for _, ch := range subs {
		go func(ch chan events.Event, event events.Event) {
			select {
			case ch <- event:
			case <-ctx.Done():
				// Handle timeout or cancellation if needed
			}
		}(ch, event)
	}
}
