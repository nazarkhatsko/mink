package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (d *Driver) Execute(ctx context.Context, options map[string]any) (driver.Output, error) {
	apiKey, _ := options["api_key"].(string)
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("claude: api_key is required")
	}

	model, _ := options["model"].(string)
	if model == "" {
		model = "claude-sonnet-4-6"
	}

	messages, ok := options["messages"]
	if !ok {
		return nil, fmt.Errorf("claude: messages is required")
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
		return nil, fmt.Errorf("claude: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("claude: create request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claude: do request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("claude: read response: %w", err)
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
		return nil, fmt.Errorf("claude: unmarshal response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("claude: api error (%d): %s", resp.StatusCode, parsed.Error.Message)
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
