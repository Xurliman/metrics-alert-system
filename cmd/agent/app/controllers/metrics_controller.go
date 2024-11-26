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

func (c *MetricsController) SendMetricsWithParams(ctx context.Context) (errs error) {
	inputCh := c.service.Generator(ctx)
	resultCh := c.service.URLConstructor(ctx, inputCh, c.address)

	for result := range resultCh {
		if err := result.Error(); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		request, err := http.NewRequest("POST", result.URL(), nil)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		request.Header.Add("Content-Type", "text/plain")
		if err = c.MakeRequest(request); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
	}
	return nil
}

func (c *MetricsController) SendCompressedMetrics(ctx context.Context) (errs error) {
	inputCh := c.service.Generator(ctx)
	resultCh := c.service.RequestCompressor(ctx, c.service.ByteTransformer(ctx, inputCh))

	url := fmt.Sprintf("http://%s/update/", c.address)

	for result := range resultCh {
		if err := result.Error(); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		request, err := http.NewRequest("POST", url, bytes.NewReader(result.Bytes()))
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		dst, err := c.hashRequest(result.Bytes())
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			request.Header.Set("HashSHA256", dst)
		}

		request.Header.Set("Content-Encoding", "gzip")
		if err = c.MakeRequest(request); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (c *MetricsController) SendCompressedMetricsWithParams(ctx context.Context) (errs error) {
	inputCh := c.service.Generator(ctx)
	resultCh := c.service.URLConstructor(ctx, inputCh, c.address)

	for result := range resultCh {
		if err := result.Error(); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		request, err := http.NewRequest("POST", result.URL(), nil)
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

func (c *MetricsController) SendMetrics(ctx context.Context) (errs error) {
	url := fmt.Sprintf("http://%s/update/", c.address)

	inputCh := c.service.Generator(ctx)
	resultCh := c.service.ByteTransformer(ctx, inputCh)

	for result := range resultCh {
		if err := result.Error(); err != nil {
			errs = errors.Join(errs, err)
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

		if err = c.MakeRequest(request); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return nil
}

func (c *MetricsController) MakeRequest(request *http.Request) (errs error) {
	request.Header.Add("Content-Type", "application/json")

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
