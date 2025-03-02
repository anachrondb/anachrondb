package db

import "time"

func (s *Store) Diff(key string, from, to time.Time) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var diffs []Event
	for _, event := range s.events {
		if event.Timestamp.After(from) && event.Timestamp.Before(to) {
			if event.Key == key {
				diffs = append(diffs, event)
			}
		}
	}
	return diffs
}