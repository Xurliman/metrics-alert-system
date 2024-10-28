package config

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"os"
	"strconv"
	"strings"
)

func GetEnvironmentValue(key string) (string, error) {
	if os.Getenv(key) == "" {
		return "", constants.ErrEnvValueMissing
	}
	return os.Getenv(key), nil
}

func GetPort() (string, error) {
	address, err := GetEnvironmentValue("ADDRESS")
	if err != nil {
		return "", err
	}

	options := strings.Split(address, ":")
	if len(options) < 2 {
		return "", constants.ErrWrongAddress
	}

	port, err := strconv.Atoi(options[1])
	if err != nil {
		return "", err
	}

	return ":" + strconv.Itoa(port), nil
}

func GetAppEnv() string {
	appEnv, err := GetEnvironmentValue("APP_ENV")
	if err != nil {
		return "development"
	}
	return appEnv
}
