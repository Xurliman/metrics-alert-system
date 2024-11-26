package controllers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"net/http"
	"time"
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
			errs = errors.Join(errs, err)
		}

		dst, err := c.hashRequest(requestBody)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			request.Header.Set("HashSHA256", dst)
		}

		errors.Join(errs, c.MakeRequest(request))
	}
	return nil
}

func (c *MetricsController) SendMetricsWithParams() (errs error) {
	urls, err := c.service.GetRequestURLs(c.address)
	if err != nil {
		return errors.Join(errs, err)
	}

	for _, url := range urls {
		request, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return errors.Join(errs, err)
		}

		request.Header.Set("Content-Type", "text/plain")

		if err = c.MakeRequest(request); err != nil {
			return errors.Join(errs, err)
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

		request.Header.Set("Content-Encoding", "gzip")
		errors.Join(errs, c.MakeRequest(request))
	}
	return errs
}

func (c *MetricsController) SendCompressedMetricsWithParams() (errs error) {
	urls, err := c.service.GetRequestURLs(c.address)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	for _, url := range urls {
		request, err := http.NewRequest("POST", url, nil)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		request.Header.Set("Content-Encoding", "gzip")

		if err = c.MakeRequest(request); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
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

	request.Header.Set("Content-Encoding", "gzip")
	return c.MakeRequest(request)
}

func (c *MetricsController) SendMetricsChan() (errs error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://%s/update/", c.address)

	inputCh := c.service.Generator(ctx)
	resultCh := c.service.Converter(ctx, inputCh)

	for result := range resultCh {
		if result.HasError() {
			errors.Join(errs, result.Error())
			continue
		}

		request, err := http.NewRequest("POST", url, bytes.NewReader(result.Bytes()))
		if err != nil {
			return errors.Join(errs, err)
		}

		dst, err := c.hashRequest(result.Bytes())
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			request.Header.Set("HashSHA256", dst)
		}

		errors.Join(errs, c.MakeRequest(request))
	}
	return nil
}

func (c *MetricsController) MakeRequest(request *http.Request) (errs error) {
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return errors.Join(errs, err)
	}

	err = response.Body.Close()
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if response.StatusCode != http.StatusOK {
		errs = errors.Join(errs, constants.ErrStatusNotOK)
	}
	return errs
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
