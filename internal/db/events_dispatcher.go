package db

import (
	"log"
	"sync"
)

type EventSubscriber func(Event)

type EventDispatcher struct {
	mu sync.RWMutex
	subscribers []EventSubscriber
	eventQueue chan Event
	quit chan struct{}
}

func NewEventDispatcher(buffersize int) *EventDispatcher {
	ed := &EventDispatcher{
		subscribers: make([]EventSubscriber, 0),
		eventQueue: make(chan Event, buffersize),
		quit: make(chan struct{}),
	}

	    go func() {
        for {
            select {
            case event := <-ed.eventQueue:
                ed.dispatch(event)
            case <-ed.quit:
                return
            }
        }
    }()
    
    return ed
}

func (ed *EventDispatcher) Publish(event Event) {
	select {
		case ed.eventQueue <- event:
			log.Printf("[INFO] Event published %+v", event)
		default:
			log.Printf("[WARN] Event queue is full, dropping event %+v", event)
	}
}

func (ed *EventDispatcher) Subscribe(Subscriber EventSubscriber) {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	ed.subscribers = append(ed.subscribers, Subscriber)
}

func (ed *EventDispatcher) dispatch(event Event) {
    ed.mu.RLock()
    defer ed.mu.RUnlock()
    for _, subscriber := range ed.subscribers {
        go subscriber(event)
    }
}

func (ed *EventDispatcher) Stop() {
    close(ed.quit)
}