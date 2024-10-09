package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnvironmentValue(key string) (string, error) {
	if os.Getenv(key) == "" {
		return "", errors.New("environment variable is missing")
	}
	return os.Getenv(key), nil
}

func GetHost() (string, error) {
	address, err := GetEnvironmentValue("ADDRESS")
	if err != nil {
		return "", err
	}
	options := strings.Split(address, ":")
	if len(options) < 2 {
		return "", errors.New("need 2 values as host:port")
	}
	port, err := strconv.Atoi(options[1])
	if err != nil {
		return "", err
	}
	host := options[0]
	return host + ":" + strconv.Itoa(port), nil
}

func GetPollInterval() (time.Duration, error) {
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

func GetReportInterval() (time.Duration, error) {
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
