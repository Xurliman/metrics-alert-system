package constants

import "errors"

var (
	ErrWrongAddress              = errors.New("need 2 values as host:port")
	ErrWrongPort                 = errors.New("wrong port value given")
	ErrLoadingEnv                = errors.New("error loading .env file, using default variables")
	ErrInvalidMetricType         = errors.New("invalid metrics type")
	ErrInvalidCounterMetricValue = errors.New("invalid metrics value for counter type")
	ErrInvalidGaugeMetricValue   = errors.New("invalid metrics value for gauge type")
	ErrEmptyMetricName           = errors.New("metrics name is empty")
	ErrMetricExists              = errors.New("metric exists with other type")
	ErrDatabaseDSNEmpty          = errors.New("database DSN is empty, using file storage")
	ErrConnectingDatabase        = errors.New("error connecting to database")
	ErrMetricNotFound            = errors.New("metric not found")
	ErrLoadingMetricsFromArchive = errors.New("error loading metrics from archive")
	ErrInvalidHash               = errors.New("invalid hash")
)
