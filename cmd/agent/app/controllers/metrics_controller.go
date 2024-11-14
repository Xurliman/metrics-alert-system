package controllers

import (
	"bytes"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
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

func (c *MetricsController) SendMetrics() (err error) {
	requestsToSend, err := c.service.GetRequestBodies()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/update/", c.address)
	for _, requestBody := range requestsToSend {
		response, err := c.client.Post(url, "application/json", bytes.NewReader(requestBody))
		if err != nil {
			return err
		}

		err = response.Body.Close()
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return constants.ErrStatusNotOK
		}
	}
	return nil
}

func (c *MetricsController) SendMetricsWithParams() (err error) {
	urls, err := c.service.GetRequestURLs(c.address)
	if err != nil {
		return err
	}

	for _, url := range urls {
		response, err := c.client.Post(url, "text/plain", nil)
		if err != nil {
			return err
		}

		err = response.Body.Close()
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return constants.ErrStatusNotOK
		}
	}
	return nil
}

func (c *MetricsController) SendCompressedMetrics() (err error) {
	requestsToSend, err := c.service.GetCompressedRequestBodies()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/update/", c.address)

	for _, compressedRequest := range requestsToSend {
		request, err := http.NewRequest("POST", url, bytes.NewReader(compressedRequest))
		if err != nil {
			return err
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")

		response, err := c.client.Do(request)
		if err != nil {
			return err
		}

		err = response.Body.Close()
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return constants.ErrStatusNotOK
		}
	}
	return nil
}

func (c *MetricsController) SendCompressedMetricsWithParams() (err error) {
	urls, err := c.service.GetRequestURLs(c.address)
	if err != nil {
		return err
	}

	for _, url := range urls {
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		response, err := c.client.Do(req)
		if err != nil {
			return err
		}

		err = response.Body.Close()
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return constants.ErrStatusNotOK
		}
	}
	return nil
}

func (c *MetricsController) SendBatchMetrics() (err error) {
	requestToSend, err := c.service.GetCompressedRequestBody()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/updates/", c.address)
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestToSend))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return constants.ErrStatusNotOK
	}

	return nil
}
