// Package interfaces define set of methods to implement in server side logic
package interfaces

import "github.com/Xurliman/metrics-alert-system/cmd/server/app/models"

type ArchiveServiceInterface interface {
	Archive(metrics map[string]*models.Metrics) error
	Load() (map[string]*models.Metrics, error)
}
