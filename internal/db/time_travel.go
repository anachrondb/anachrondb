package db

import "time"

func (s *Store) GetAt(key string, at time.Time) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var value string
	exists := false

	for _, event := range s.events {
		if event.Timestamp.After(at) {
			break
		}
		if event.Key == key {
			switch event.Type {
			case EventSet:
				value = string(event.NewValue)
				exists = true
			case EventDelete:
				value = ""
				exists = false
			}
		}
	}
	return value, exists
}