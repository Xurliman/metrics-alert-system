package config

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"os"
	"strconv"
	"strings"
	"time"
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
		return constants.DefaultPort, err
	}

	options := strings.Split(address, ":")
	if len(options) < 2 {
		return constants.DefaultPort, constants.ErrWrongAddress
	}

	port, err := strconv.Atoi(options[1])
	if err != nil {
		return constants.DefaultPort, err
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

func GetStoreInterval() time.Duration {
	storeIntervalStr, err := GetEnvironmentValue("STORE_INTERVAL")
	if err != nil {
		return time.Second * constants.DefaultStoreInterval
	}

	storeInterval, err := strconv.Atoi(storeIntervalStr)
	if err != nil {
		return time.Second * constants.DefaultStoreInterval
	}

	if storeInterval == 0 {
		return time.Second
	}
	return time.Second * time.Duration(storeInterval)
}

func GetFileStoragePath() string {
	fileStoragePath, err := GetEnvironmentValue("FILE_STORAGE_PATH")
	if err != nil {
		return constants.DefaultFileStoragePath
	}
	return fileStoragePath
}

func GetShouldRestore() bool {
	shouldRestore, err := GetEnvironmentValue("RESTORE")
	if err != nil {
		return constants.DefaultRestore
	}
	return shouldRestore == "true"
}

func GetKey() string {
	key, err := GetEnvironmentValue("KEY")
	if err != nil {
		return ""
	}
	return key
}
