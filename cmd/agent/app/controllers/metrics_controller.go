package controllers

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/workerpool"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"net/http"
	"sync"
	"time"
)

type MetricsController struct {
	client     http.Client
	service    interfaces.MetricsService
	cfg        *config.Config
	workerPool *workerpool.WorkerPool
}

func NewMetricsController(service interfaces.MetricsService, cfg *config.Config) interfaces.MetricsController {
	handlers := []func(metrics *models.Metrics) error{
		service.SendMetric,
		service.SendCompressedMetric,
		service.SendMetricWithParams,
	}
	workerPool := workerpool.New(cfg.RateLimit, handlers)
	return &MetricsController{
		client:     http.Client{},
		service:    service,
		cfg:        cfg,
		workerPool: workerPool,
	}
}

func (c *MetricsController) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		c.Poll(ctx)
		wg.Done()
	}()

	go func() {
		c.Report(ctx)
		wg.Done()
	}()

	wg.Wait()
	c.workerPool.Stop()
}

func (c *MetricsController) Poll(ctx context.Context) {
	pollTicker := time.NewTicker(c.cfg.GetPollInterval())
	defer pollTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			err := c.service.CollectMetricValues()
			if err != nil {
				log.Error(err.Error())
			}
		case <-ctx.Done():
			log.Info("stopping polling metrics")
			return
		}
	}
}

func (c *MetricsController) Report(ctx context.Context) {
	reportTicker := time.NewTicker(c.cfg.GetReportInterval())
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			for _, m := range c.service.GetAll() {
				c.workerPool.AddJob(m)
			}
		case <-ctx.Done():
			log.Info("stopping reporting metrics")
			return
		}
	}
}
