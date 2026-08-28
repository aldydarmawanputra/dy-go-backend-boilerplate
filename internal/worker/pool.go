package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type job struct {
	name string
	run  func(ctx context.Context) error
}

// Pool is a small in-process job queue: bounded workers, retry with backoff,
// and graceful drain. For durable/distributed jobs, swap for a Redis-backed
// queue (e.g. asynq); the Submit call sites stay the same.
type Pool struct {
	jobs    chan job
	wg      sync.WaitGroup
	workers int
	retries int
}

func New(workers, queueSize int) *Pool {
	if workers <= 0 {
		workers = 4
	}
	if queueSize <= 0 {
		queueSize = 128
	}
	return &Pool{jobs: make(chan job, queueSize), workers: workers, retries: 3}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.loop()
	}
}

func (p *Pool) loop() {
	defer p.wg.Done()
	for j := range p.jobs {
		p.process(j)
	}
}

func (p *Pool) process(j job) {
	for attempt := 1; attempt <= p.retries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := j.run(ctx)
		cancel()
		if err == nil {
			return
		}
		slog.Warn("job failed", "job", j.name, "attempt", attempt, "err", err)
		if attempt < p.retries {
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}
	}
	slog.Error("job exhausted retries", "job", j.name)
}

// Submit enqueues a job without blocking the caller; if the queue is full the
// job is dropped (and logged) rather than stalling the request path.
func (p *Pool) Submit(name string, run func(ctx context.Context) error) {
	select {
	case p.jobs <- job{name: name, run: run}:
	default:
		slog.Warn("job queue full, dropping job", "job", name)
	}
}

// Shutdown stops accepting new jobs and waits for in-flight work to drain, or
// until ctx is done.
func (p *Pool) Shutdown(ctx context.Context) {
	close(p.jobs)
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
	}
}
