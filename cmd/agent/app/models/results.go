package models

import "github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"

type Result struct {
	data    []byte
	request *requests.MetricsRequest
	url     string
	err     error
}

func NewResult(data []byte, err error) Result {
	return Result{data: data, err: err}
}

func NewURLResult(url string, err error) Result {
	return Result{url: url, err: err}
}

func NewRequestResult(request *requests.MetricsRequest, err error) Result {
	return Result{request: request, err: err}
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
