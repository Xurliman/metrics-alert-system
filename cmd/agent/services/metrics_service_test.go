package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/models"
	"math/rand"
	"net/http"
	"runtime"
	"testing"
)

func TestSendMetrics(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	type args struct {
		client  http.Client
		metrics models.Metrics
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "first",
			args: args{
				client: http.Client{},
				metrics: models.Metrics{
					Gauge: map[string]float64{
						"Alloc":         float64(memStats.Alloc),
						"BuckHashSys":   float64(memStats.BuckHashSys),
						"Frees":         float64(memStats.Frees),
						"GCCPUFraction": memStats.GCCPUFraction,
						"GCSys":         float64(memStats.GCSys),
						"HeapAlloc":     float64(memStats.HeapAlloc),
						"HeapIdle":      float64(memStats.HeapIdle),
						"HeapInuse":     float64(memStats.HeapInuse),
						"HeapObjects":   float64(memStats.HeapObjects),
						"HeapReleased":  float64(memStats.HeapReleased),
						"HeapSys":       float64(memStats.HeapSys),
						"LastGC":        float64(memStats.LastGC),
						"Lookups":       float64(memStats.Lookups),
						"MCacheInuse":   float64(memStats.MCacheInuse),
						"MCacheSys":     float64(memStats.MCacheSys),
						"MSpanInuse":    float64(memStats.MSpanInuse),
						"MSpanSys":      float64(memStats.MSpanSys),
						"Mallocs":       float64(memStats.Mallocs),
						"NextGC":        float64(memStats.NextGC),
						"NumForcedGC":   float64(memStats.NumForcedGC),
						"NumGC":         float64(memStats.NumGC),
						"OtherSys":      float64(memStats.OtherSys),
						"PauseTotalNs":  float64(memStats.PauseTotalNs),
						"StackInuse":    float64(memStats.StackInuse),
						"StackSys":      float64(memStats.StackSys),
						"Sys":           float64(memStats.Sys),
						"TotalAlloc":    float64(memStats.TotalAlloc),
						"RandomValue":   rand.Float64(),
					},
					Counter: map[string]int64{
						"PollCount": pollCount,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMetrics(tt.args.client, tt.args.metrics)
		})
	}
}
