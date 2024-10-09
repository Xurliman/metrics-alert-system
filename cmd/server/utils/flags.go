package utils

import (
	"errors"
	"flag"
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
		return errors.New("address should be given as host:port")
	}
	port, err := strconv.Atoi(options[1])
	if err != nil {
		return errors.New("wrong port value given")
	}
	o.host = options[0]
	o.port = port
	return nil
}

func (o *Options) String() string {
	return o.host + ":" + strconv.Itoa(o.port)
}

func GetPort() string {
	options := &Options{
		host: "localhost",
		port: 8080,
	}
	flag.Var(options, "a", "give server host:port (default: localhost:8080)")
	flag.Parse()
	return ":" + strconv.Itoa(options.port)
}
