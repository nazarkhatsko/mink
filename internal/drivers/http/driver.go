package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "http" }

func (d *Driver) Execute(ctx context.Context, options map[string]any) (driver.Output, error) {
	method, _ := options["method"].(string)
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	url, _ := options["url"].(string)
	if url == "" {
		return nil, driver.NewError(driver.ErrConfig, "http: url is required")
	}

	var bodyReader io.Reader
	if body, ok := options["body"]; ok && body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, driver.Wrap(driver.ErrConfig, err, "http: marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, driver.Wrap(driver.ErrConfig, err, "http: create request: %v", err)
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if headers, ok := options["headers"].(map[string]any); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, driver.Wrap(driver.ErrTransport, err, "http: do request: %v", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, driver.Wrap(driver.ErrTransport, err, "http: read response: %v", err)
	}

	respHeaders := make(map[string]any, len(resp.Header))
	for k := range resp.Header {
		respHeaders[k] = resp.Header.Get(k)
	}

	var respBody any
	if err := json.Unmarshal(respBodyBytes, &respBody); err != nil {
		respBody = string(respBodyBytes)
	}

	reqHeaders := make(map[string]any, len(req.Header))
	for k := range req.Header {
		reqHeaders[k] = req.Header.Get(k)
	}

	return driver.Output{
		"req": map[string]any{
			"method":  method,
			"url":     url,
			"headers": reqHeaders,
			"body":    options["body"],
		},
		"resp": map[string]any{
			"status":  resp.StatusCode,
			"headers": respHeaders,
			"body":    respBody,
		},
	}, nil
}
