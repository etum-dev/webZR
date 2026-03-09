package scanjob

import (
	config "github.com/etum-dev/WebZR/internal/flaginput"
	"github.com/etum-dev/WebZR/pkg/utils"
	"github.com/etum-dev/WebZR/pkg/worker"
)

// Job is a single domain scan request.
type Job struct {
	Domain string
	Flags  *config.Flags
}

// JobResult is the output of a processed Job.
type JobResult struct {
	Job     Job
	Results []utils.ScanResult
	Err     error
}

// JobHandler processes a single Job and returns a JobResult.
type JobHandler func(Job) JobResult

// Handler manages a pool of scan workers using a pool-of-channels pattern
// https://jsschools.com/golang/advanced-go-channel-patterns-for-building-robust-d/
type Handler struct {
	JobQueue chan Job
	// pool for workers to tell when available
	Pool      chan chan Job
	WaitGroup worker.WaitGroup
	threads   int
	handler   JobHandler
}

// scanWorker is a single worker process owned by a Handler.
type scanWorker struct {
	pool       chan chan Job
	jobChannel chan Job
	handler    JobHandler
}

// NewHandler creates a Handler that runs threads worker goroutines, each
// invoking h to process jobs.
func NewHandler(threads int, h JobHandler) Handler {
	if threads < 1 {
		threads = 1
	}
	return Handler{
		JobQueue: make(chan Job),
		Pool:     make(chan chan Job, threads),
		threads:  threads,
		handler:  h,
	}
}

// Run spawns all worker goroutines and starts the dispatch loop.
// Results are forwarded to listener. Run is non-blocking; use Wait to
// block until all enqueued jobs are finished.
func (h *Handler) Run(listener chan<- JobResult) {
	result := make(chan JobResult)

	for i := 0; i < h.threads; i++ {
		w := newScanWorker(h.Pool, h.handler)
		go w.spawn(result)
	}

	go func() {
		for {
			select {
			case job := <-h.JobQueue:
				go func(job Job) {
					jobChan := <-h.Pool
					jobChan <- job
				}(job)
			case r := <-result:
				listener <- r
				h.WaitGroup.Done()
			}
		}
	}()
}

// AddJob enqueues a Job for processing.
func (h *Handler) AddJob(job Job) {
	h.WaitGroup.Add(1)
	h.JobQueue <- job
}

// Wait blocks until all enqueued jobs have completed.
func (h *Handler) Wait() {
	h.WaitGroup.Wait()
}

// GetJobCount returns the number of jobs currently in flight.
func (h *Handler) GetJobCount() int {
	return h.WaitGroup.GetCount()
}

func newScanWorker(pool chan chan Job, h JobHandler) scanWorker {
	return scanWorker{
		pool:       pool,
		jobChannel: make(chan Job),
		handler:    h,
	}
}

// spawn registers the worker into the pool and processes jobs indefinitely.
func (w scanWorker) spawn(result chan<- JobResult) {
	for {
		w.pool <- w.jobChannel
		job := <-w.jobChannel
		result <- w.handler(job)
	}
}
