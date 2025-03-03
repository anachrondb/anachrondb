package db

import (
	"log"
)

type EventDispatcher struct {
	eventQueue  chan Event
	subscribers []func(Event)
	workerPool  *WorkerPool
}

func NewEventDispatcher(queueSize int, workerCount int) *EventDispatcher {
	dispatcher := &EventDispatcher{
		eventQueue: make(chan Event, queueSize),
		workerPool: NewWorkerPool(workerCount, queueSize),
	}
	go dispatcher.dispatch()
	return dispatcher
}

func (ed *EventDispatcher) Stop() {
	close(ed.eventQueue)
	ed.workerPool.Stop()
}

func (ed *EventDispatcher) Subscribe(handler func(Event)) {
	ed.subscribers = append(ed.subscribers, handler)
}

func (ed *EventDispatcher) Publish(event Event) {
	ed.eventQueue <- event
	log.Printf("[INFO] Event published %+v", event)
}

func (ed *EventDispatcher) dispatch() {
	for event := range ed.eventQueue {
		e := event // capture loop variable safely

		e.subscribers = ed.subscribers

		ed.workerPool.Submit(e)
	}
}
