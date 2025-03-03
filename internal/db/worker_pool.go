// internal/db/worker_pool.go
package db

import (
	"sync"
)

// WorkerPool handles concurrent event processing.
type WorkerPool struct {
	workers int
	jobs    chan Event
	wg      sync.WaitGroup
}

// NewWorkerPool creates a new pool with a given number of workers.
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		workers: workers,
		jobs:    make(chan Event, queueSize),
	}
	wp.start()
	return wp
}

// start launches worker goroutines.
func (wp *WorkerPool) start() {
	wp.wg.Add(wp.workers)
	for i := 0; i < wp.workers; i++ {
		go func() {
			defer wp.wg.Done()
			for event := range wp.jobs {
				event.Handle()
			}
		}()
	}
}

// Submit adds a job to the pool.
func (wp *WorkerPool) Submit(event Event) {
	wp.jobs <- event
}

// Stop gracefully shuts down the worker pool.
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}