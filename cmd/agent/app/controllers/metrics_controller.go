package controllers

import (
	"bytes"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"net/http"
)

type MetricsController struct {
	client  http.Client
	service interfaces.MetricsService
	address string
}

func NewMetricsController(client http.Client, service interfaces.MetricsService, address string) interfaces.MetricsController {
	return &MetricsController{
		client:  client,
		service: service,
		address: address,
	}
}

func (c *MetricsController) CollectMetrics() {
	c.service.CollectMetricValues()
}

func (c *MetricsController) SendMetrics() {
	requestsToSend, err := c.service.GetRequestBodies()
	if err != nil {
		return
	}

	url := fmt.Sprintf("http://%s/update/", c.address)
	for _, requestBody := range requestsToSend {
		response, err := c.client.Post(url, "application/json", bytes.NewReader(requestBody))
		if err != nil {
			return
		}

		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func (c *MetricsController) SendMetricsWithParams() {
	urls, err := c.service.GetRequestURLs(c.address)
	if err != nil {
		return
	}

	for _, url := range urls {
		response, err := c.client.Post(url, "text/plain", nil)
		if err != nil {
			return
		}
		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func (c *MetricsController) SendCompressedMetrics() {
	requestsToSend, err := c.service.GetCompressedRequestBodies()
	if err != nil {
		return
	}
	url := fmt.Sprintf("http://%s/update/", c.address)

	for _, compressedRequest := range requestsToSend {
		request, err := http.NewRequest("POST", url, bytes.NewReader(compressedRequest))
		if err != nil {
			return
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")

		response, err := c.client.Do(request)
		if err != nil {
			return
		}

		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func (c *MetricsController) SendCompressedMetricsWithParams() {
	urls, err := c.service.GetRequestURLs(c.address)
	if err != nil {
		return
	}

	for _, url := range urls {
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		response, err := c.client.Do(req)
		if err != nil {
			return
		}
		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func (c *MetricsController) SendManyMetrics() {
	requestToSend, err := c.service.GetCompressedRequestBody()
	if err != nil {
		return
	}
	url := fmt.Sprintf("http://%s/updates/", c.address)
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestToSend))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	response, err := c.client.Do(req)
	if err != nil {
		return
	}
	err = response.Body.Close()
	if err != nil {
		return
	}
}
