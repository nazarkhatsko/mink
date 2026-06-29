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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http: read response: %w", err)
	}

	headers := make(map[string]any, len(resp.Header))
	for k := range resp.Header {
		headers[k] = resp.Header.Get(k)
	}

	var parsedBody any
	if err := json.Unmarshal(respBody, &parsedBody); err != nil {
		parsedBody = string(respBody)
	}

	return driver.Output{
		"status":  resp.StatusCode,
		"body":    parsedBody,
		"headers": headers,
	}, nil
}
