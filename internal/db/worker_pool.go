package db

import (
	"sync"
)

type WorkerPool struct {
	workers int
	jobs    chan Event
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		workers: workers,
		jobs:    make(chan Event, queueSize),
	}
	wp.start()
	return wp
}

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

func (wp *WorkerPool) Submit(event Event) {
	wp.jobs <- event
}

func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}