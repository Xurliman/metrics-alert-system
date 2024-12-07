package config

import (
	"flag"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/caarlos0/env/v11"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Host            string `env:"ADDRESS" envDefault:"localhost:8080"`
	Port            int    `env:"PORT" envDefault:"8080"`
	StoreInterval   int    `env:"STORE_INTERVAL" envDefault:"5"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Restore         bool   `env:"RESTORE" envDefault:"false"`
	Key             string `env:"KEY" envDefault:""`
}

func Setup() (*Config, error) {
	cfg := &Config{
		Host:            "localhost",
		Port:            0,
		StoreInterval:   1,
		FileStoragePath: "",
		Restore:         true,
		Key:             "",
	}
	flag.Var(cfg, constants.AddressFlag, constants.AddressFlagDescription)
	flag.IntVar(&cfg.StoreInterval, constants.StoreIntervalFlag, 1, constants.StoreIntervalFlagDescription)
	flag.StringVar(&cfg.FileStoragePath, constants.FileStoragePathFlag, constants.DefaultFileStoragePath, constants.FileStoragePathFlagDescription)
	flag.BoolVar(&cfg.Restore, constants.RestoreFlag, constants.DefaultRestore, constants.RestoreFlagDescription)
	flag.StringVar(&cfg.DatabaseDSN, constants.DatabaseDSNFlag, "", constants.DatabaseDSNFlagDescription)
	flag.StringVar(&cfg.Key, constants.KeyFlag, "", constants.KeyFlagDescription)
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (o *Config) Set(flagValue string) (err error) {
	options := strings.Split(flagValue, ":")
	if len(options) != 2 {
		return constants.ErrWrongAddress
	}

	port, err := strconv.Atoi(options[1])
	if err != nil {
		return constants.ErrWrongPort
	}

	o.Host = options[0]
	o.Port = port
	return nil
}

func (o *Config) String() string {
	return o.Host + ":" + strconv.Itoa(o.Port)
}

func (o *Config) GetPort() string {
	port := ":" + strconv.Itoa(o.Port)
	if port == ":0" {
		return constants.DefaultPort
	}
	return port
}

func (o *Config) GetStoreInterval() time.Duration {
	return time.Duration(o.StoreInterval) * time.Second
}
