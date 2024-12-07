package workerpool

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"go.uber.org/zap"
	"sync"
)

type WorkerPool struct {
	handlers []func(metrics *models.Metrics) error
	jobs     chan *models.Metrics
	wg       sync.WaitGroup
	stopOnce sync.Once
}

func New(workerCount int, handlers []func(metrics *models.Metrics) error) *WorkerPool {
	pool := &WorkerPool{
		handlers: handlers,
		jobs:     make(chan *models.Metrics, workerCount),
	}
	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go pool.worker()
	}
	return pool
}
func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for job := range p.jobs {
		for _, handler := range p.handlers {
			log.Info("started handling a job")
			if err := handler(job); err != nil {
				log.Error("job failed", zap.Error(err))
			}
		}
	}
}

func (p *WorkerPool) AddJob(metric *models.Metrics) {
	select {
	case p.jobs <- metric:
	default:
		//
	}
}

func (p *WorkerPool) Stop() {
	p.stopOnce.Do(func() {
		close(p.jobs)
		p.wg.Wait()
	})
}
