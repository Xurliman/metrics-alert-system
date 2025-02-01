package models

import "github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"

type Result struct {
	data    []byte
	request *requests.MetricsRequest
	url     string
	err     error
}

func (r Result) Bytes() []byte {
	return r.data
}

func (r Result) URL() string {
	return r.url
}

func (r Result) Error() error {
	return r.err
}

func (r Result) Request() *requests.MetricsRequest {
	return r.request
}
