package constants

import "errors"

var (
	ErrEnvValueMissing           = errors.New("environment variable is missing")
	ErrWrongAddress              = errors.New("need 2 values as host:port")
	ErrWrongPort                 = errors.New("wrong port value given")
	ErrLoadingEnv                = errors.New("error loading .env file, using default variables")
	ErrInvalidMetricType         = errors.New("invalid metrics type")
	ErrInvalidCounterMetricValue = errors.New("invalid metrics value for counter type")
	ErrInvalidGaugeMetricValue   = errors.New("invalid metrics value for gauge type")
	ErrEmptyMetricName           = errors.New("metrics name is empty")
)
