package utils

import (
	"errors"
	"flag"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	host           string
	port           int
	reportInterval int
	pollInterval   int
}

func (o *Options) SetHost(host string) {
	o.host = host
}

func (o *Options) SetPort(port int) {
	o.port = port
}

func (o *Options) SetReportInterval(reportInterval int) {
	o.reportInterval = reportInterval
}

func (o *Options) SetPollInterval(pollInterval int) {
	o.pollInterval = pollInterval
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

func ParseFlags() *Options {
	options := &Options{
		host: "localhost",
		port: 8080,
	}
	flag.IntVar(&options.reportInterval, "r", 10, "set report interval")
	flag.IntVar(&options.pollInterval, "p", 2, "set poll interval")
	flag.Var(options, "a", "give server host:port (default: localhost:8080)")
	flag.Parse()
	return options
}

func (o *Options) GetAddr() string {
	return o.host + ":" + strconv.Itoa(o.port)
}

func (o *Options) GetPollInterval() time.Duration {
	return time.Duration(o.pollInterval) * time.Second
}

func (o *Options) GetReportInterval() time.Duration {
	return time.Duration(o.reportInterval) * time.Second
}
