package db

import "time"

type EventType string

const (
	EventSet EventType = "SET"
	EventDel EventType = "DEL"
)

type Event struct {
	Timestamp time.Time
	Type      EventType
	Key       string
	OldValue  string
	NewValue  string
}
