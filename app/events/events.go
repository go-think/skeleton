package events

import "sync"

type OrderCreatedEvent struct {
	OrderID string
	Amount  float64
}

type EventTracker struct {
	mu            sync.Mutex
	HandledEvents []string
}

var Tracker = &EventTracker{
	HandledEvents: make([]string, 0),
}

func (t *EventTracker) Record(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.HandledEvents = append(t.HandledEvents, name)
}

func (t *EventTracker) Has(name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, e := range t.HandledEvents {
		if e == name {
			return true
		}
	}
	return false
}

func (t *EventTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.HandledEvents = make([]string, 0)
}
