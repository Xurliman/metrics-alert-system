package constants

import "errors"

var (
	ErrLoadingEnv        = errors.New("error loading .env file, using default variables")
	ErrInvalidMetricType = errors.New("invalid metrics type")
	ErrInvalidMetric     = errors.New("invalid metric, it's equal to nil")
	ErrStatusNotOK       = errors.New("status not OK")
	ErrKeyMissing        = errors.New("key is missing")
)
