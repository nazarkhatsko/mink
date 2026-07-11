package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

const apiURL = "https://api.anthropic.com/v1/messages"

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "claude" }

func (d *Driver) Methods() []string { return []string{"message"} }

func (d *Driver) Options(method string) []driver.Option {
	return []driver.Option{
		{Name: "messages", Required: true},
		{Name: "system"},
		{Name: "max_tokens"},
		{Name: "temperature"},
	}
}

func (d *Driver) Execute(ctx context.Context, method string, options map[string]any) (driver.Output, error) {
	if method != "message" {
		return nil, driver.NewError(driver.ErrConfig, "claude: unknown method %q", method)
	}

	apiKey, _ := options["api_key"].(string)
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		return nil, driver.NewError(driver.ErrConfig, "claude: api_key is required")
	}

	model, _ := options["model"].(string)
	if model == "" {
		model = "claude-sonnet-4-6"
	}

	messages, ok := options["messages"]
	if !ok {
		return nil, driver.NewError(driver.ErrConfig, "claude: messages is required")
	}

	maxTokens := 1024
	switch v := options["max_tokens"].(type) {
	case int:
		maxTokens = v
	case float64:
		maxTokens = int(v)
	}

	reqBody := map[string]any{
		"model":      model,
		"messages":   messages,
		"max_tokens": maxTokens,
	}
	if system, ok := options["system"].(string); ok && system != "" {
		reqBody["system"] = system
	}
	if temperature, ok := options["temperature"]; ok {
		reqBody["temperature"] = temperature
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, driver.Wrap(driver.ErrInternal, err, "claude: marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, driver.Wrap(driver.ErrInternal, err, "claude: create request: %v", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, driver.Wrap(driver.ErrTransport, err, "claude: do request: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, driver.Wrap(driver.ErrTransport, err, "claude: read response: %v", err)
	}

	var parsed struct {
		ID         string `json:"id"`
		Model      string `json:"model"`
		Role       string `json:"role"`
		StopReason string `json:"stop_reason"`
		Content    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, driver.Wrap(driver.ErrInternal, err, "claude: unmarshal response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, driver.NewError(driver.ErrTransport, "claude: api error (%d): %s", resp.StatusCode, parsed.Error.Message)
	}

	var text strings.Builder
	for _, block := range parsed.Content {
		if block.Type == "text" {
			text.WriteString(block.Text)
		}
	}

	return driver.Output{
		"in": map[string]any{
			"model":      model,
			"messages":   messages,
			"system":     reqBody["system"],
			"max_tokens": maxTokens,
		},
		"out": map[string]any{
			"id":          parsed.ID,
			"model":       parsed.Model,
			"role":        parsed.Role,
			"content":     text.String(),
			"stop_reason": parsed.StopReason,
			"usage": map[string]any{
				"input_tokens":  parsed.Usage.InputTokens,
				"output_tokens": parsed.Usage.OutputTokens,
			},
		},
	}, nil
}
