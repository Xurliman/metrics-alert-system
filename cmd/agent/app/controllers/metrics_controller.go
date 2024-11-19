package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"net/http"
)

type MetricsController struct {
	client  http.Client
	service interfaces.MetricsService
	address string
	key     string
}

func NewMetricsController(client http.Client, service interfaces.MetricsService, address string, key string) interfaces.MetricsController {
	return &MetricsController{
		client:  client,
		service: service,
		address: address,
		key:     key,
	}
}

func (c *MetricsController) CollectMetrics() {
	c.service.CollectMetricValues()
}

func (c *MetricsController) SendMetrics() (errs error) {
	requestsToSend, err := c.service.GetRequestBodies()
	if err != nil {
		errs = errors.Join(errs, err)
	}

	url := fmt.Sprintf("http://%s/update/", c.address)
	for _, requestBody := range requestsToSend {
		request, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
		if err != nil {
			return errors.Join(errs, err)
		}

		request.Header.Set("Content-Type", "application/json")

		dst, err := c.hashRequest(requestBody)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			request.Header.Set("HashSHA256", dst)
		}

		response, err := c.client.Do(request)
		if err != nil {
			return errors.Join(errs, err)
		}

		err = response.Body.Close()
		if err != nil {
			return errors.Join(errs, err)
		}

		if response.StatusCode != http.StatusOK {
			return errors.Join(errs, constants.ErrStatusNotOK)
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

func (c *MetricsController) SendCompressedMetrics() (errs error) {
	requestsToSend, err := c.service.GetCompressedRequestBodies()
	if err != nil {
		errors.Join(errs, err)
	}
	url := fmt.Sprintf("http://%s/update/", c.address)

	for _, compressedRequest := range requestsToSend {
		request, err := http.NewRequest("POST", url, bytes.NewReader(compressedRequest))
		if err != nil {
			return errors.Join(errs, err)
		}

		dst, err := c.hashRequest(compressedRequest)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			request.Header.Set("HashSHA256", dst)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")

		response, err := c.client.Do(request)
		if err != nil {
			return errors.Join(errs, err)
		}

		err = response.Body.Close()
		if err != nil {
			return errors.Join(errs, err)
		}

		if response.StatusCode != http.StatusOK {
			return errors.Join(errs, constants.ErrStatusNotOK)
		}
	}
	return errs
}

func (c *MetricsController) SendCompressedMetricsWithParams() (err error) {
	urls, err := c.service.GetRequestURLs(c.address)
	if err != nil {
		return err
	}

	for _, url := range urls {
		request, err := http.NewRequest("POST", url, nil)
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

func (c *MetricsController) SendBatchMetrics() (err error) {
	requestToSend, err := c.service.GetCompressedRequestBody()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/updates/", c.address)
	request, err := http.NewRequest("POST", url, bytes.NewReader(requestToSend))
	if err != nil {
		return err
	}

	dst, err := c.hashRequest(requestToSend)
	if err != nil {
		return err
	} else {
		request.Header.Set("HashSHA256", dst)
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

	return nil
}

func (c *MetricsController) hashRequest(requestBody []byte) (string, error) {
	if c.key == "" {
		return "", constants.ErrKeyMissing
	}
	h := hmac.New(sha256.New, []byte(c.key))
	h.Write(requestBody)
	dst := h.Sum(nil)
	return hex.EncodeToString(dst), nil
}
