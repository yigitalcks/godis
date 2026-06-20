package interpreter

import (
	"context"
	"sync"
	"time"

	"godis/internal/logging"
)

type Expiryjob struct {
	key   string
	delay time.Duration
}

type ExpiryScheduler struct {
	jobs   chan Expiryjob
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	data *sync.Map
}

func (e *ExpiryScheduler) Start() {
	go func() {
		for {
			select {
			case job, ok := <-e.jobs:
				if !ok {
					return
				}

				e.wg.Add(1)
				go func(job Expiryjob) {
					defer e.wg.Done()
					time.Sleep(job.delay)

					e.data.Delete(job.key)
					logging.Println("godis: key expired -", job.key)
				}(job)
			case <-e.ctx.Done():
				return
			}
		}
	}()
}

func (e *ExpiryScheduler) AddJob(job Expiryjob) {
	e.jobs <- job
}

func (e *ExpiryScheduler) Shutdown() {
	e.cancel()
	close(e.jobs)
	e.wg.Wait()
}

func NewExpiryScheduler(jobs chan Expiryjob, data *sync.Map) *ExpiryScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &ExpiryScheduler{
		jobs:   jobs,
		ctx:    ctx,
		cancel: cancel,
		data:   data,
	}
}

func NewExpiryJob(key string, delay time.Duration) *Expiryjob {
	return &Expiryjob{
		key:   key,
		delay: delay,
	}
}
