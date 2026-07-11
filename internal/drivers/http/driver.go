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

func (d *Driver) Methods() []string {
	return []string{"get", "post", "put", "patch", "delete", "request"}
}

func (d *Driver) Options(method string) []driver.Option {
	base := []driver.Option{
		{Name: "url", Required: true},
		{Name: "headers"},
		{Name: "body"},
	}
	if method == "request" {
		return append([]driver.Option{{Name: "method", Required: true}}, base...)
	}
	return base
}

func (d *Driver) Execute(ctx context.Context, method string, options map[string]any) (driver.Output, error) {
	switch method {
	case "get":
		return d.doRequest(ctx, http.MethodGet, options)
	case "post":
		return d.doRequest(ctx, http.MethodPost, options)
	case "put":
		return d.doRequest(ctx, http.MethodPut, options)
	case "patch":
		return d.doRequest(ctx, http.MethodPatch, options)
	case "delete":
		return d.doRequest(ctx, http.MethodDelete, options)
	case "request":
		httpMethod, _ := options["method"].(string)
		if httpMethod == "" {
			return nil, driver.NewError(driver.ErrConfig, "http: method is required for the request method")
		}
		return d.doRequest(ctx, strings.ToUpper(httpMethod), options)
	default:
		return nil, driver.NewError(driver.ErrConfig, "http: unknown method %q", method)
	}
}

func (d *Driver) doRequest(ctx context.Context, method string, options map[string]any) (driver.Output, error) {
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
