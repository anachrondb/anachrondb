package db

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration, t *testing.T) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return
	case <-time.After(timeout):
		t.Fatalf("test timed out after %s", timeout)
	}
}

func TestEventDispatcher_SingleSubscriber(t *testing.T) {
	dispatcher := NewEventDispatcher(1000, 10)
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	dispatcher.Subscribe(func(e Event) {
		defer wg.Done()
		if e.Key != "test-key" {
			t.Errorf("unexpected event key: %s", e.Key)
		}
	})

	dispatcher.Publish(Event{
		Timestamp: time.Now(),
		Type:      EventSet,
		Key:       "test-key",
	})

	waitWithTimeout(&wg, 2*time.Second, t)
}

func TestEventDispatcher_MultipleSubscribers(t *testing.T) {
	dispatcher := NewEventDispatcher(1000, 10)
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	subscriberCount := 10
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
	})

	waitWithTimeout(&wg, 2*time.Second, t)
}

func TestEventDispatcher_ConcurrentPublish(t *testing.T) {
	dispatcher := NewEventDispatcher(2_000_000, 500)
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	subscriberCount := 100
	eventCount := 1_000_000
	wg.Add(subscriberCount * eventCount)

	for i := 0; i < subscriberCount; i++ {
		dispatcher.Subscribe(func(e Event) {
			defer wg.Done()
		})
	}

	for i := 0; i < eventCount; i++ {
		dispatcher.Publish(Event{
			Timestamp: time.Now(),
			Type:      EventSet,
			Key:       fmt.Sprintf("key-%d", i),
		})
	}

	waitWithTimeout(&wg, 30*time.Second, t)
}

func TestEventDispatcher_QueueOverflow(t *testing.T) {
	dispatcher := NewEventDispatcher(1, 2)
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	dispatcher.Subscribe(func(e Event) {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
	})

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

	waitWithTimeout(&wg, 1*time.Second, t)
}
