package config

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfInterface interface {
	GetHost() (string, error)
	GetPollInterval() (time.Duration, error)
	GetReportInterval() (time.Duration, error)
	GetKey() (string, error)
}

type EnvConfig struct{}

func GetEnvironmentValue(key string) (string, error) {
	if os.Getenv(key) == "" {
		return "", constants.ErrEnvValueMissing
	}
	return os.Getenv(key), nil
}

func (c *EnvConfig) GetHost() (string, error) {
	address, err := GetEnvironmentValue("ADDRESS")
	if err != nil {
		return constants.DefaultServerAddress, err
	}

	options := strings.Split(address, ":")
	if len(options) < 2 {
		return "", constants.ErrWrongAddress
	}

	port, err := strconv.Atoi(options[1])
	if err != nil {
		return "", err
	}

	host := options[0]
	return host + ":" + strconv.Itoa(port), nil
}

func (c *EnvConfig) GetPollInterval() (time.Duration, error) {
	pollInterval, err := GetEnvironmentValue("POLL_INTERVAL")
	if err != nil {
		return time.Duration(2), err
	}

	pollIntervalInt, err := strconv.Atoi(pollInterval)
	if err != nil {
		return time.Duration(2), err
	}

	return time.Duration(pollIntervalInt) * time.Second, nil
}

func (c *EnvConfig) GetReportInterval() (time.Duration, error) {
	reportInterval, err := GetEnvironmentValue("REPORT_INTERVAL")
	if err != nil {
		return time.Duration(10), err
	}

	reportIntervalInt, err := strconv.Atoi(reportInterval)
	if err != nil {
		return time.Duration(10), err
	}

	return time.Duration(reportIntervalInt) * time.Second, nil
}

func (c *EnvConfig) GetKey() (string, error) {
	key, err := GetEnvironmentValue("KEY")
	if err != nil {
		return "", constants.ErrKeyMissing
	}
	return key, nil
}

func NewConfig() ConfInterface {
	return &EnvConfig{}
}
