package constants

import "errors"

var (
	ErrLoadingEnv          = errors.New("error loading .env file")
	ErrHostNotPassedAsFlag = errors.New("host is not passed as with a flag")
	ErrWrongAddress        = errors.New("address should be given as host:port")
	ErrWrongPort           = errors.New("wrong port value given")
	ErrEnvValueMissing     = errors.New("environment variable is missing")
)
