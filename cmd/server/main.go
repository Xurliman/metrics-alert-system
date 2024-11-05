package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"runtime"
	"time"
)

var lastSavedMetrics = make(map[string]*models.Metrics)

func main() {
	utils.Logger = utils.NewLogger(config.GetAppEnv())

	err := godotenv.Load(constants.EnvFilePath)
	if err != nil {
		utils.Logger.Error("", zap.Error(constants.ErrLoadingEnv))
	}

	flagOptions := utils.NewOptions()

	port, err := flagOptions.GetPort()
	if err != nil {
		port, _ = config.GetPort()
	}

	fileStoragePath, err := flagOptions.GetFileStoragePath()
	if err != nil {
		fileStoragePath = config.GetFileStoragePath()
	}

	storeInterval, err := flagOptions.GetStoreInterval()
	if err != nil {
		storeInterval = config.GetStoreInterval()
	}

	shouldRestore := flagOptions.GetShouldRestore() && config.GetShouldRestore()
	SetMetrics(shouldRestore, fileStoragePath)
	r := routes.SetupRoutes()
	utils.Logger.Error("Starting server on port %v",
		zap.String("port", port),
	)
	log.Printf("store interval %v", storeInterval)

	go func() {
		storeTicker := time.NewTicker(storeInterval)
		defer storeTicker.Stop()
		for {
			select {
			case <-storeTicker.C:
				err = Archive(services.MetricsCollection, fileStoragePath)
				if err != nil {
					utils.Logger.Error("archiving data went wrong",
						zap.Error(err),
					)
				}
			}
		}
	}()
	err = r.Run(port)
	if err != nil {
		return
	}
}

func Archive(metrics map[string]*models.Metrics, fileStoragePath string) (err error) {
	log.Println("archiving")
	toSave := make(map[string]*models.Metrics)
	for currentMetricKey, currentMetric := range metrics {
		if lastMetric, exists := lastSavedMetrics[currentMetricKey]; exists && currentMetric.Equals(lastMetric) {
			continue
		}
		toSave[currentMetricKey] = currentMetric
	}

	archiveWriter, err := utils.NewArchiveWriter(fileStoragePath)
	if err != nil {
		utils.Logger.Error("creating archive writer error", zap.Error(err))
		log.Fatal("error creating archive writer ", err)
	}
	defer func(archiveWriter *utils.ArchiveWriter) {
		err = archiveWriter.Close()
		if err != nil {
			utils.Logger.Error("closing write archive error", zap.Error(err))
		}
	}(archiveWriter)

	return archiveWriter.Archive(toSave)
}

func LoadDefaultMetricCollection() map[string]*models.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gaugeMetricExample := float64(memStats.Alloc)
	var counterMetricExample int64
	return map[string]*models.Metrics{
		"Alloc": {
			ID:    "Alloc",
			MType: constants.GaugeMetricType,
			Value: &gaugeMetricExample,
		},
		"PollCount": {
			ID:    "PollCount",
			MType: constants.CounterMetricType,
			Delta: &counterMetricExample,
		},
	}
}

func LoadMetricsFromFile(fileStoragePath string) (map[string]*models.Metrics, error) {
	archiveReader, err := utils.NewArchiveReader(fileStoragePath)
	if err != nil {
		return nil, err
	}
	defer func(archiveReader *utils.ArchiveReader) {
		err = archiveReader.Close()
		if err != nil {
			utils.Logger.Error("closing read archive error", zap.Error(err))
		}
	}(archiveReader)

	loadedMetrics, err := archiveReader.LoadMetrics()
	if err != nil {
		return nil, err
	}

	return loadedMetrics, nil
}

func SetMetrics(shouldRestore bool, fileStoragePath string) {
	if !shouldRestore {
		services.MetricsCollection = LoadDefaultMetricCollection()
		return
	}

	metrics, err := LoadMetricsFromFile(fileStoragePath)
	if err != nil {
		utils.Logger.Error("error loading metrics from file", zap.Error(err), zap.String("file", fileStoragePath))
		services.MetricsCollection = LoadDefaultMetricCollection()
		return
	}
	services.MetricsCollection = metrics
}
