package db

import (
	"time"
)

type EventType string

const (
	EventSet    EventType = "SET"
	EventDelete EventType = "DELETE"
)

type Event struct {
	Timestamp time.Time
	Type      EventType
	Key       string
	OldValue  []byte
	NewValue  []byte
	Version   int64

	subscribers []func(Event)
}

func (e Event) Handle() {
	for _, subscriber := range e.subscribers {
		subscriber(e)
	}
}
