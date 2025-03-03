package db

import (
	"sync"
	"time"
)

type Value struct {
	Data string
	UpdatedAt time.Time
}

type Store struct {
	mu sync.RWMutex
	data map[string]Value
	events []Event
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Value),
		events: []Event{},
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldValue := s.data[key].Data
	s.data[key] = Value{Data: value, UpdatedAt: time.Now()}

	s.events = append(s.events, Event{
		Timestamp: time.Now(),
		Type: EventSet,
		Key: key,
		OldValue: []byte(oldValue),
		NewValue: []byte(value),
	})
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.data[key]
	return val.Data, exists
}
