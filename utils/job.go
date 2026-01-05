package utils

import "sync"

// a single domain scan request.

type Job struct {
	Domain string
	Flags  *Flags
}

// output of the job
type JobResult struct {
	Job     Job
	Results []ScanResult
	Err     error
}

// JobHandler processes a single Job and returns a JobResult.
type JobHandler func(Job) JobResult

// fixed number of goroutines consuming Jobs.
type WorkerPool struct {
	jobs    chan Job
	results chan JobResult
	wg      sync.WaitGroup
}

// spins up workerCount workers that invoke handler for each job
// queueSize controls backpressure; values <=0 default to workerCount*2.
func NewWorkerPool(workerCount, queueSize int, handler JobHandler) *WorkerPool {
	if workerCount < 1 {
		workerCount = 1
	}
	if queueSize <= 0 {
		queueSize = workerCount * 2
	}

	pool := &WorkerPool{
		jobs:    make(chan Job, queueSize),
		results: make(chan JobResult, queueSize),
	}

	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer pool.wg.Done()
			for job := range pool.jobs {
				pool.results <- handler(job)
			}
		}()
	}

	return pool
}

// Submit enqueues a Job for processing.
func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

// Results exposes the channel used by workers to report back.
func (p *WorkerPool) Results() <-chan JobResult {
	return p.results
}

// Close stops intake, waits for workers, and then closes the results channel.
func (p *WorkerPool) Close() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
