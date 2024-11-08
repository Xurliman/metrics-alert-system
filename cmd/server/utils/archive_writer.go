package utils

import (
	"bufio"
	"encoding/json"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"os"
)

type ArchiveWriter struct {
	file    *os.File
	writer  *bufio.Writer
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

	bufferedWriter := bufio.NewWriter(file)
	return &ArchiveWriter{
		file:    file,
		writer:  bufferedWriter,
		encoder: json.NewEncoder(bufferedWriter),
	}, nil
}

func (aw *ArchiveWriter) Close() error {
	return aw.file.Close()
}

func (aw *ArchiveWriter) Archive(metrics map[string]*models.Metrics) error {
	err := aw.encoder.Encode(metrics)
	if err != nil {
		return err
	}
	return aw.writer.Flush()
}
