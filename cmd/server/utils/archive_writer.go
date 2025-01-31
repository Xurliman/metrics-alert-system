package utils

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
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
	err := aw.file.Truncate(0)
	if err != nil {
		return err
	}
	// Reset the file's offset to the beginning
	_, err = aw.file.Seek(0, 0)
	if err != nil {
		return err
	}
	aw.writer = bufio.NewWriter(aw.file)
	aw.encoder = json.NewEncoder(aw.writer)

	err = aw.encoder.Encode(metrics)
	if err != nil {
		return err
	}
	return aw.writer.Flush()
}
