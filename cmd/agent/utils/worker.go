package utils

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

type WorkerPool struct {
	rateLimit int
	handlers  []func(ctx context.Context) error
}

func NewWorkerPool(rateLimit int, handlers []func(ctx context.Context) error) *WorkerPool {
	return &WorkerPool{
		rateLimit: rateLimit,
		handlers:  handlers,
	}
}

func (wp *WorkerPool) Run() {
	var wg sync.WaitGroup
	for range wp.rateLimit {
		wg.Add(wp.rateLimit)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			for _, handler := range wp.handlers {
				err := handler(ctx)
				if err != nil {
					wp.handleError(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (wp *WorkerPool) handleError(err error) {
	if err != nil {
		Logger.Error("error while sending metrics", zap.Error(err))
	}
}
