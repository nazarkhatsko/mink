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

func (d *Driver) Describe() driver.Doc {
	return driver.Doc{
		Description: "Execute HTTP requests",
		Config: []driver.FieldDoc{
			{Name: "timeout", Type: "int", Required: false, Description: "Request timeout in milliseconds (default: no timeout)"},
		},
		Options: []driver.FieldDoc{
			{Name: "method", Type: "string", Required: false, Description: "HTTP method: GET, POST, PUT, PATCH, DELETE (default: GET)"},
			{Name: "url", Type: "string", Required: true, Description: "Full URL"},
			{Name: "headers", Type: "object", Required: false, Description: "Request headers"},
			{Name: "body", Type: "any", Required: false, Description: "Request body, serialized as JSON"},
		},
		Output: []driver.FieldDoc{
			{Name: "req", Type: "object", Description: "Outgoing request: method, url, headers, body"},
			{Name: "resp", Type: "object", Description: "Incoming response: status, headers, body"},
		},
	}
}

func (d *Driver) Execute(ctx context.Context, options map[string]any) (driver.Output, error) {
	method, _ := options["method"].(string)
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	url, _ := options["url"].(string)
	if url == "" {
		return nil, fmt.Errorf("http: url is required")
	}

	var bodyReader io.Reader
	if body, ok := options["body"]; ok && body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("http: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("http: create request: %w", err)
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
		return nil, fmt.Errorf("http: do request: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http: read response: %w", err)
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
