package utils

import (
	"flag"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"strconv"
	"strings"
)

type Options struct {
	host string
	port int
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

func GetPort() (string, error) {
	options := &Options{
		host: "localhost",
		port: 8080,
	}
	flag.Var(options, "a", "give server host:port (default: localhost:8080)")
	flag.Parse()
	port := ":" + strconv.Itoa(options.port)
	if port == "" {
		return "", constants.ErrWrongPort
	}
	return port, nil
}
