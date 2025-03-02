package db

import "time"

type EventType string

const (
	EventSet EventType = "SET"
	EventDel EventType = "DEL"
	EventUpd EventType = "UPD"
	EventGet EventType = "GET"
	EventLst EventType = "LST"
)

type Event struct {
	Timestamp time.Time
	Type      EventType
	Key       string
	OldValue  string
	NewValue  string
}
