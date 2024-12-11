package config

import (
	"flag"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `json:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func Setup() (*Config, error) {
	var conf Config
	flag.IntVar(&conf.RateLimit, "l", 0, "set rate limit")
	flag.IntVar(&conf.ReportInterval, "r", 10, "set report interval")
	flag.IntVar(&conf.PollInterval, "p", 2, "set poll interval")
	flag.StringVar(&conf.Key, "k", "", "set key to hash")
	flag.StringVar(&conf.ServerAddress, "a", constants.DefaultServerAddress, "give server host:port (default: localhost:8080)")
	flag.Parse()

	if err := env.Parse(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (cfg *Config) GetHost() string {
	return fmt.Sprintf("http://%s", cfg.ServerAddress)
}

func (cfg *Config) GetPollInterval() time.Duration {
	return time.Duration(cfg.PollInterval) * time.Second
}

func (cfg *Config) GetReportInterval() time.Duration {
	return time.Duration(cfg.ReportInterval) * time.Second
}
