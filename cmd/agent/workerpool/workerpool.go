// Package workerpool manages sending metrics concurrently
package workerpool

import (
	"context"
	"sync"

	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"go.uber.org/zap"
)

type job struct {
	ctx    context.Context
	metric *models.Metrics
}

type WorkerPool struct {
	handlers []func(ctx context.Context, metrics *models.Metrics) error
	jobs     chan job
	wg       sync.WaitGroup
	stopOnce sync.Once
}

func New(workerCount int, handlers []func(ctx context.Context, metrics *models.Metrics) error) *WorkerPool {
	pool := &WorkerPool{
		handlers: handlers,
		jobs:     make(chan job, workerCount),
	}
	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go pool.worker()
	}
	return pool
}
func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.jobs {
		select {
		case <-task.ctx.Done():
			log.Warn("job canceled due to context cancellation")
			continue
		default:
			for _, handler := range p.handlers {
				log.Info("started handling a job")
				if err := handler(task.ctx, task.metric); err != nil {
					log.Error("job failed", zap.Error(err))
				}
			}
		}
	}
}

func (p *WorkerPool) AddJob(ctx context.Context, metric *models.Metrics) {
	select {
	case p.jobs <- job{ctx: ctx, metric: metric}:
	default:
		log.Warn("job queue is full, job was dropped")
	}
}

func (p *WorkerPool) Stop() {
	p.stopOnce.Do(func() {
		close(p.jobs)
		p.wg.Wait()
	})
}
