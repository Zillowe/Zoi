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

type GoogleProvider struct {
	client  *http.Client
	apiKey  string
	model   string
	baseURL string
}

type googleRESTRequest struct {
	Contents []googleContent `json:"contents"`
}

type googleContent struct {
	Parts []googlePart `json:"parts"`
}

type googlePart struct {
	Text string `json:"text"`
}

type googleRESTResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func NewGoogleProvider(apiKey, modelName string) (*GoogleProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Google AI Studio API key is required")
	}

	return &GoogleProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		model:   modelName,
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/",
	}, nil
}

func (p *GoogleProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := googleRESTRequest{
		Contents: []googleContent{
			{
				Parts: []googlePart{
					{Text: prompt},
				},
			},
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal google rest request: %w", err)
	}

	url := fmt.Sprintf("%s%s:generateContent?key=%s", p.baseURL, p.model, p.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to google ai: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp googleRESTResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse google json response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("google api error (%d - %s): %s", apiResp.Error.Code, apiResp.Error.Status, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from google ai: %d", resp.StatusCode)
	}

	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("received an empty or invalid response from google ai")
	}

	return apiResp.Candidates[0].Content.Parts[0].Text, nil
}
