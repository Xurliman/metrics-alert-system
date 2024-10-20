package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func GetEnvironmentValue(key string) (string, error) {
	if os.Getenv(key) == "" {
		return "", errors.New("environment variable is missing")
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
		return "", errors.New("need 2 values as host:port")
	}
	port, err := strconv.Atoi(options[1])
	if err != nil {
		return "", err
	}
	return ":" + strconv.Itoa(port), nil
}
