package constants

import "errors"

var (
	ErrLoadingEnv                = errors.New("error loading .env file, using default variables")
	ErrHostNotPassedAsFlag       = errors.New("host is not passed as with a flag")
	ErrWrongAddress              = errors.New("address should be given as host:port")
	ErrWrongPort                 = errors.New("wrong port value given")
	ErrInvalidMetricType         = errors.New("invalid metrics type")
	ErrEnvValueMissing           = errors.New("environment variable is missing")
	ErrInvalidCounterMetricValue = errors.New("invalid metrics value for counter type")
	ErrInvalidGaugeMetricValue   = errors.New("invalid metrics value for gauge type")
	ErrStatusNotOK               = errors.New("status not OK")
	ErrInvalidRateLimit          = errors.New("invalid rate limit")
	ErrKeyMissing                = errors.New("key is missing")
)
