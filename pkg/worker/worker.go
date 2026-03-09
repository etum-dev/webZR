package worker

import (
	"sync"
	"sync/atomic"
)

// WaitGroup wraps sync.WaitGroup with an atomic counter so callers can
// query the number of in-flight jobs without a separate mutex.
type WaitGroup struct {
	sync.WaitGroup
	count int64
}

func (wg *WaitGroup) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *WaitGroup) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
}

// GetCount returns the number of jobs currently in flight.
func (wg *WaitGroup) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}
