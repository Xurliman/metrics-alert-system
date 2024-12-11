package constants

import "errors"

var (
	ErrLoadingEnv                = errors.New("error loading .env file, using default variables")
	ErrInvalidMetricType         = errors.New("invalid metrics type")
	ErrInvalidMetric             = errors.New("invalid metric, it's equal to nil")
	ErrInvalidCounterMetricValue = errors.New("invalid metrics value for counter type")
	ErrInvalidGaugeMetricValue   = errors.New("invalid metrics value for gauge type")
	ErrStatusNotOK               = errors.New("status not OK")
	ErrKeyMissing                = errors.New("key is missing")
)
