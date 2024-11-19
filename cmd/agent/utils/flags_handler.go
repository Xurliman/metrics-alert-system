package utils

import (
	"flag"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"strconv"
	"strings"
	"time"
)

type FlagConfig struct {
	host           string
	port           int
	reportInterval int
	pollInterval   int
	key            string
}

func (cfg *FlagConfig) SetHost(host string) {
	cfg.host = host
}

func (cfg *FlagConfig) SetPort(port int) {
	cfg.port = port
}

func (cfg *FlagConfig) SetReportInterval(reportInterval int) {
	cfg.reportInterval = reportInterval
}

func (cfg *FlagConfig) SetPollInterval(pollInterval int) {
	cfg.pollInterval = pollInterval
}

func (cfg *FlagConfig) GetHost() (string, error) {
	host := cfg.String()
	if host == ":0" {
		return "", constants.ErrHostNotPassedAsFlag
	}
	return host, nil
}

func (cfg *FlagConfig) GetPollInterval() (time.Duration, error) {
	return time.Duration(cfg.pollInterval) * time.Second, nil
}

func (cfg *FlagConfig) GetReportInterval() (time.Duration, error) {
	return time.Duration(cfg.reportInterval) * time.Second, nil
}

func (cfg *FlagConfig) GetKey() (string, error) {
	if cfg.key == "" {
		return "", constants.ErrKeyMissing
	}
	return cfg.key, nil
}

func NewOptions() config.ConfInterface {
	options := &FlagConfig{
		host: "",
		port: 0,
	}
	flag.IntVar(&options.reportInterval, "r", 10, "set report interval")
	flag.IntVar(&options.pollInterval, "p", 2, "set poll interval")
	flag.StringVar(&options.key, "k", "", "set key to hash")
	flag.Var(options, "a", "give server host:port (default: localhost:8080)")
	flag.Parse()
	return options
}

func (cfg *FlagConfig) Set(flagValue string) (err error) {
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

func (cfg *FlagConfig) String() string {
	return cfg.host + ":" + strconv.Itoa(cfg.port)
}
