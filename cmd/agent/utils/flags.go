package utils

import (
	"flag"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"strconv"
	"strings"
	"time"
)

type EnvConfig struct {
	host           string
	port           int
	reportInterval int
	pollInterval   int
	appEnv         string
}

func (cfg *EnvConfig) SetHost(host string) {
	cfg.host = host
}

func (cfg *EnvConfig) SetPort(port int) {
	cfg.port = port
}

func (cfg *EnvConfig) SetReportInterval(reportInterval int) {
	cfg.reportInterval = reportInterval
}

func (cfg *EnvConfig) SetPollInterval(pollInterval int) {
	cfg.pollInterval = pollInterval
}

func (cfg *EnvConfig) GetHost() (string, error) {
	host := cfg.String()
	if host == ":0" {
		return "", constants.ErrHostNotPassedAsFlag
	}
	return host, nil
}

func (cfg *EnvConfig) GetPollInterval() (time.Duration, error) {
	return time.Duration(cfg.pollInterval) * time.Second, nil
}

func (cfg *EnvConfig) GetReportInterval() (time.Duration, error) {
	return time.Duration(cfg.reportInterval) * time.Second, nil
}

func NewOptions() config.ConfInterface {
	options := &EnvConfig{
		host: "",
		port: 0,
	}
	flag.IntVar(&options.reportInterval, "r", 10, "set report interval")
	flag.IntVar(&options.pollInterval, "p", 2, "set poll interval")
	flag.Var(options, "a", "give server host:port (default: localhost:8080)")
	flag.Parse()
	return options
}

func (cfg *EnvConfig) Set(flagValue string) (err error) {
	options := strings.Split(flagValue, ":")
	if len(options) != 2 {
		return constants.ErrWrongAddress
	}
	port, err := strconv.Atoi(options[1])
	if err != nil {
		return constants.ErrWrongPort
	}
	cfg.host = options[0]
	cfg.port = port
	return nil
}

func (cfg *EnvConfig) String() string {
	return cfg.host + ":" + strconv.Itoa(cfg.port)
}
