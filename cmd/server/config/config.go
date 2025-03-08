package config

import (
	"encoding/json"
	"flag"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Host            string `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	Port            int    `env:"PORT" envDefault:"8080" json:"port"`
	StoreInterval   int    `env:"STORE_INTERVAL" envDefault:"5" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`
	Restore         bool   `env:"RESTORE" envDefault:"false" json:"restore"`
	Key             string `env:"KEY" envDefault:"invalid_key" json:"key"`
	CryptoKey       string `env:"CRYPTO_KEY" envDefault:"" json:"crypto_key"`
	ConfigFile      string `env:"CONFIG" envDefault:""`
}

func Setup() (*Config, error) {
	var cfg Config
	flag.Var(&cfg, constants.AddressFlag, constants.AddressFlagDescription)
	flag.IntVar(&cfg.StoreInterval, constants.StoreIntervalFlag, 1, constants.StoreIntervalFlagDescription)
	flag.StringVar(&cfg.FileStoragePath, constants.FileStoragePathFlag, constants.DefaultFileStoragePath, constants.FileStoragePathFlagDescription)
	flag.BoolVar(&cfg.Restore, constants.RestoreFlag, constants.DefaultRestore, constants.RestoreFlagDescription)
	flag.StringVar(&cfg.DatabaseDSN, constants.DatabaseDSNFlag, "", constants.DatabaseDSNFlagDescription)
	flag.StringVar(&cfg.Key, constants.KeyFlag, "", constants.KeyFlagDescription)
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "set private key path")
	flag.StringVar(&cfg.ConfigFile, "c", "", "set config file path")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	if cfg.ConfigFile != "" {
		if err := cfg.parseConfigFile(); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
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
	log.Warn("Set updated config", zap.Any("config", o))
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

func (o *Config) parseConfigFile() error {
	file, err := os.Open(o.ConfigFile)
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

	*o = fileCfg
	return nil
}
