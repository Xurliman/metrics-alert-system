package utils

import (
	"encoding/json"
	"os"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
)

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

func (ar *ArchiveReader) Load() (map[string]*models.Metrics, error) {
	metrics := make(map[string]*models.Metrics)
	err := ar.decoder.Decode(&metrics)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
