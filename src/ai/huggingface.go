package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const huggingFaceAPIBaseURL = "https://api-inference.huggingface.co/models/"

type HuggingFaceProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type huggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

type huggingFaceResponse []struct {
	GeneratedText string `json:"generated_text"`
}

func NewHuggingFaceProvider(apiKey, modelName string) (*HuggingFaceProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Hugging Face API key is required")
	}
	return &HuggingFaceProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: huggingFaceAPIBaseURL,
	}, nil
}

func (p *HuggingFaceProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := huggingFaceRequest{Inputs: prompt}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal huggingface request: %w", err)
	}

	url := p.baseURL + p.model
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to huggingface: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status from huggingface: %s", string(respBody))
	}

	var apiResp huggingFaceResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse huggingface json response: %w", err)
	}

	if len(apiResp) == 0 || apiResp[0].GeneratedText == "" {
		return "", fmt.Errorf("received an empty response from huggingface")
	}

	return strings.TrimPrefix(apiResp[0].GeneratedText, prompt), nil
}
