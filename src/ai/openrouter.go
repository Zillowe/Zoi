package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	openRouterAPIBaseURL = "https://openrouter.ai/api/v1"
	gctAppName           = "https://codeberg.org/Zusty/GCT"
)

type OpenRouterProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type openRouterRequest struct {
	Model    string              `json:"model"`
	Messages []openRouterMessage `json:"messages"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewOpenRouterProvider(apiKey, modelName string) (*OpenRouterProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is required")
	}

	return &OpenRouterProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: openRouterAPIBaseURL,
	}, nil
}

func (p *OpenRouterProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := openRouterRequest{
		Model: p.model,
		Messages: []openRouterMessage{
			{Role: "user", Content: prompt},
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal openrouter request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("HTTP-Referer", gctAppName)
	req.Header.Set("X-Title", "GCT AI Commit")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to openrouter: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp openRouterResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse openrouter json response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("openrouter api error (%d): %s", resp.StatusCode, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from openrouter: %d", resp.StatusCode)
	}

	if len(apiResp.Choices) == 0 || apiResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("received an empty or invalid response from openrouter")
	}

	return apiResp.Choices[0].Message.Content, nil
}
