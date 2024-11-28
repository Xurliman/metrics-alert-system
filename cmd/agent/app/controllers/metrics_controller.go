package controllers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
	"go.uber.org/zap"
	"log"
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
			errs = c.handleErrors(result.URL(), result.Bytes(), errs)
		}
	}
	return nil
}

func (c *MetricsController) SendCompressedMetrics(ctx context.Context) (errs error) {
	inputCh := c.service.Generator(ctx)
	reqCh := c.service.RequestConstructor(ctx, inputCh)
	bytesCh := c.service.ByteTransformer(ctx, reqCh)
	resultCh := c.service.RequestCompressor(ctx, bytesCh)

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
		if err != nil && !errors.Is(err, constants.ErrKeyMissing) {
			errs = errors.Join(errs, err)
		} else {
			request.Header.Set("HashSHA256", dst)
		}

		request.Header.Set("Content-Encoding", "gzip")
		if err = c.MakeRequest(request); err != nil {
			errs = c.handleErrors(url, result.Bytes(), errs)
		}
	}
	return errs
}

func (c *MetricsController) SendBatchMetrics(ctx context.Context) (err error) {
	inputCh := c.service.Generator(ctx)
	reqCh := c.service.RequestConstructor(ctx, inputCh)

	var reqs []*requests.MetricsRequest
	for result := range reqCh {
		if err = result.Error(); err != nil {
			return err
		}
		reqs = append(reqs, result.Request())
	}

	requestBytes, err := json.Marshal(reqs)
	if err != nil {
		return err
	}
	compressedRequest, err := utils.Compress(requestBytes)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/updates/", c.address)
	request, err := http.NewRequest("POST", url, bytes.NewReader(compressedRequest))
	if err != nil {
		return err
	}

	dst, err := c.hashRequest(compressedRequest)
	if err != nil && !errors.Is(err, constants.ErrKeyMissing) {
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
	reqCh := c.service.RequestConstructor(ctx, inputCh)
	resultCh := c.service.ByteTransformer(ctx, reqCh)

	for result := range resultCh {
		if err := result.Error(); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		log.Println("Before sending a request ", string(result.Bytes()))
		request, err := http.NewRequest("POST", url, bytes.NewReader(result.Bytes()))
		if err != nil {
			return errors.Join(errs, err)
		}

		dst, err := c.hashRequest(result.Bytes())
		if err != nil && !errors.Is(err, constants.ErrKeyMissing) {
			errs = errors.Join(errs, err)
		} else {
			request.Header.Set("HashSHA256", dst)
		}

		if err = c.MakeRequest(request); err != nil {
			errs = c.handleErrors(url, result.Bytes(), errs)
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
		return errors.Join(errs, err)
	}

	if response.StatusCode != http.StatusOK {
		return errors.Join(errs, constants.ErrStatusNotOK)
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

func (c *MetricsController) handleErrors(url string, body []byte, err error) (errs error) {
	errs = errors.Join(errs, err)
	if errs != nil {
		utils.Logger.Error("error while making request",
			zap.Error(errs),
			zap.String("url", url),
			zap.String("body", string(body)),
		)
	}
	return nil
}
