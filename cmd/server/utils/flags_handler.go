package utils

import (
	"flag"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	host            string
	port            int
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func NewOptions() *Options {
	options := &Options{
		host:            "localhost",
		port:            0,
		storeInterval:   -1,
		fileStoragePath: "",
		restore:         true,
	}
	flag.Var(options, constants.AddressFlag, constants.AddressFlagDescription)
	flag.IntVar(&options.storeInterval, constants.StoreIntervalFlag, -1, constants.StoreIntervalFlagDescription)
	flag.StringVar(&options.fileStoragePath, constants.FileStoragePathFlag, constants.FileStoragePath, constants.FileStoragePathFlagDescription)
	flag.BoolVar(&options.restore, constants.RestoreFlag, constants.Restore, constants.RestoreFlagDescription)
	flag.Parse()
	return options
}

func (o *Options) Set(flagValue string) (err error) {
	options := strings.Split(flagValue, ":")
	if len(options) != 2 {
		return constants.ErrWrongAddress
	}

	port, err := strconv.Atoi(options[1])
	if err != nil {
		return constants.ErrWrongPort
	}

	o.host = options[0]
	o.port = port
	return nil
}

func (o *Options) String() string {
	return o.host + ":" + strconv.Itoa(o.port)
}

func (o *Options) GetPort() (string, error) {
	port := ":" + strconv.Itoa(o.port)
	if port == ":0" {
		return "", constants.ErrWrongPort
	}
	return port, nil
}

func (o *Options) GetStoreInterval() (time.Duration, error) {
	storeInterval := o.storeInterval
	if storeInterval == -1 {
		return -1, constants.ErrWrongStoreInterval
	}

	if storeInterval == 0 {
		return time.Second, nil
	}

	return time.Second * time.Duration(storeInterval), nil
}

func (o *Options) GetFileStoragePath() (string, error) {
	if o.fileStoragePath == "" {
		return "", constants.ErrWrongFileStoragePath
	}
	return o.fileStoragePath, nil
}

func (o *Options) GetShouldRestore() bool {
	return o.restore
}
