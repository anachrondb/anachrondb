package db

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestEventDispatcher_SingleSubscriber(t *testing.T) {
	dispatcher := NewEventDispatcher(10)
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	dispatcher.Subscribe(func(e Event) {
		defer wg.Done()
		if e.Key != "test-key" {
			t.Errorf("expected key 'test-key', got '%s'", e.Key)
		}
	})

	dispatcher.Publish(Event{
		Timestamp: time.Now(),
		Type:      EventSet,
		Key:       "test-key",
		NewValue:  []byte("value"),
	})

	waitWithTimeout(&wg, 1*time.Second, t)
}

func TestEventDispatcher_MultipleSubscribers(t *testing.T) {
	dispatcher := NewEventDispatcher(10)
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	subscriberCount := 5
	wg.Add(subscriberCount)

	for i := 0; i < subscriberCount; i++ {
		dispatcher.Subscribe(func(e Event) {
			defer wg.Done()
		})
	}

	dispatcher.Publish(Event{
		Timestamp: time.Now(),
		Type:      EventSet,
		Key:       "multi-key",
		NewValue:  []byte("value"),
	})

	waitWithTimeout(&wg, 1*time.Second, t)
}

func TestEventDispatcher_QueueOverflow(t *testing.T) {
	dispatcher := NewEventDispatcher(1)
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	dispatcher.Subscribe(func(e Event) {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
	})

	// First event gets through
	dispatcher.Publish(Event{
		Timestamp: time.Now(),
		Type:      EventSet,
		Key:       "first",
	})

	dispatcher.Publish(Event{
		Timestamp: time.Now(),
		Type:      EventSet,
		Key:       "second",
	})

	waitWithTimeout(&wg, 2*time.Second, t)
}

func TestEventDispatcher_ConcurrentPublish(t *testing.T) {
	dispatcher := NewEventDispatcher(5000000) // Increase queue size
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	subscriberCount := 10000
	eventCount := 1000000
	totalEvents := subscriberCount * eventCount
	wg.Add(totalEvents)

	for i := 0; i < subscriberCount; i++ {
		dispatcher.Subscribe(func(e Event) {
			wg.Done()
		})
	}

	for i := 0; i < eventCount; i++ {
		dispatcher.Publish(Event{
			Timestamp: time.Now(),
			Type:      EventSet,
			Key:       fmt.Sprintf("key-%d", i),
		})
	}

	waitWithTimeout(&wg, 1*time.Minute, t)
}

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration, t *testing.T) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatal("test timed out")
	}
}
