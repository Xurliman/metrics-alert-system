package utils

import (
	"encoding/json"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"os"
)

type ArchiveWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewArchiveWriter(filename string) (*ArchiveWriter, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(wd+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ArchiveWriter{file: file, encoder: json.NewEncoder(file)}, nil
}

func (aw *ArchiveWriter) Close() error {
	return aw.file.Close()
}

func (aw *ArchiveWriter) Archive(metrics map[string]*models.Metrics) error {
	return aw.encoder.Encode(metrics)
}

type ArchiveReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewArchiveReader(filename string) (*ArchiveReader, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(wd+filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &ArchiveReader{file: file, decoder: json.NewDecoder(file)}, nil
}

func (ar *ArchiveReader) Close() error {
	return ar.file.Close()
}

func (ar *ArchiveReader) LoadMetrics() (map[string]*models.Metrics, error) {
	metrics := make(map[string]*models.Metrics)
	err := ar.decoder.Decode(&metrics)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
