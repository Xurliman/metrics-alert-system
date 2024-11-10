package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"go.uber.org/zap"
)

type ArchiveService struct {
	path             string
	lastSavedMetrics map[string]*models.Metrics
}

func (a ArchiveService) Archive(metrics map[string]*models.Metrics) error {
	toSave := make(map[string]*models.Metrics)
	for currentMetricKey, currentMetric := range metrics {
		if lastMetric, exists := a.lastSavedMetrics[currentMetricKey]; exists && currentMetric.Equals(lastMetric) {
			continue
		}
		toSave[currentMetricKey] = currentMetric
	}

	writer, err := utils.NewArchiveWriter(a.path)
	if err != nil {
		return err
	}
	defer func(archiveWriter *utils.ArchiveWriter) {
		err = archiveWriter.Close()
		if err != nil {
			utils.Logger.Error("error closing write archive", zap.Error(err))
		}
	}(writer)

	return writer.Archive(toSave)
}

func (a ArchiveService) Load() (map[string]*models.Metrics, error) {
	reader, err := utils.NewArchiveReader(a.path)
	if err != nil {
		return nil, err
	}
	defer func(archiveReader *utils.ArchiveReader) {
		err = archiveReader.Close()
		if err != nil {
			utils.Logger.Error("closing read archive error", zap.Error(err))
		}
	}(reader)

	loadedMetrics, err := reader.Load()
	if err != nil {
		return nil, err
	}

	return loadedMetrics, nil
}

func NewArchiveService(path string) interfaces.ArchiveServiceInterface {
	lastSavedMetrics := make(map[string]*models.Metrics)
	return &ArchiveService{
		path:             path,
		lastSavedMetrics: lastSavedMetrics,
	}
}
