package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"os"
	"time"
)

type Config struct {
	ServerAddress  string `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2" json:"poll_interval"`
	Key            string `envDefault:"" json:"key"`
	RateLimit      int    `env:"RATE_LIMIT" envDefault:"0" json:"rate_limit"`
	CryptoKey      string `env:"CRYPTO_KEY" envDefault:"" json:"crypto_key"`
	ConfigFile     string `env:"CONFIG" envDefault:""`
}

func Setup() (*Config, error) {
	var conf Config
	flag.IntVar(&conf.RateLimit, "l", 0, "set rate limit")
	flag.IntVar(&conf.ReportInterval, "r", 10, "set report interval")
	flag.IntVar(&conf.PollInterval, "p", 2, "set poll interval")
	flag.StringVar(&conf.Key, "k", "", "set key to hash")
	flag.StringVar(&conf.ServerAddress, "a", constants.DefaultServerAddress, "give server host:port (default: localhost:8080)")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "", "set public key path")
	flag.StringVar(&conf.ConfigFile, "c", "", "set config file path")
	flag.Parse()

	if err := env.Parse(&conf); err != nil {
		return nil, err
	}

	if conf.ConfigFile != "" {
		if err := conf.parseConfigFile(); err != nil {
			return nil, err
		}
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

func (cfg *Config) parseConfigFile() error {
	file, err := os.Open(cfg.ConfigFile)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error("Failed to close config file", zap.Error(err))
		}
	}(file)

	decoder := json.NewDecoder(file)
	fileCfg := Config{}
	if err = decoder.Decode(&fileCfg); err != nil {
		return err
	}

	if fileCfg.ServerAddress != "" {
		cfg.ServerAddress = fileCfg.ServerAddress
	}
	if fileCfg.ReportInterval != 0 {
		cfg.ReportInterval = fileCfg.ReportInterval
	}
	if fileCfg.PollInterval != 0 {
		cfg.PollInterval = fileCfg.PollInterval
	}
	if fileCfg.RateLimit != 0 {
		cfg.RateLimit = fileCfg.RateLimit
	}
	if fileCfg.Key != "" {
		cfg.Key = fileCfg.Key
	}
	if fileCfg.CryptoKey != "" {
		cfg.CryptoKey = fileCfg.CryptoKey
	}

	return nil
}
