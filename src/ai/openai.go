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
	openAIBaseURL = "https://api.openai.com/v1"
)

type OpenAIProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func NewOpenAIProvider(apiKey, modelName string) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	return &OpenAIProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: openAIBaseURL,
	}, nil
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := openAIRequest{
		Model: p.model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal openai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to openai: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp openAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse openai json response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("openai api error (type: %s): %s", apiResp.Error.Type, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from openai: %d", resp.StatusCode)
	}

	if len(apiResp.Choices) == 0 || apiResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("received an empty or invalid response from openai")
	}

	return apiResp.Choices[0].Message.Content, nil
}
