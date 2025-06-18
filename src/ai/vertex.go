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

type VertexAIProvider struct {
	client    *http.Client
	apiKey    string
	model     string
	projectID string
	region    string
	baseURL   string
}

type vertexAIRequest struct {
	Contents         []vertexAIContent `json:"contents"`
	GenerationConfig vertexAIGenConfig `json:"generation_config"`
}

type vertexAIContent struct {
	Role  string         `json:"role"`
	Parts []vertexAIPart `json:"parts"`
}
type vertexAIPart struct {
	Text string `json:"text"`
}

type vertexAIGenConfig struct {
	MaxOutputTokens int `json:"maxOutputTokens"`
}

type vertexAIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewVertexAIProvider(apiKey, modelName, projectID, region string) (*VertexAIProvider, error) {
	if apiKey == "" || modelName == "" || projectID == "" || region == "" {
		return nil, fmt.Errorf("API Key, Model, GCP Project ID, and GCP Region are all required for Vertex AI")
	}

	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:generateContent",
		region, projectID, region, modelName)

	return &VertexAIProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:    apiKey,
		model:     modelName,
		projectID: projectID,
		region:    region,
		baseURL:   url,
	}, nil
}

func (p *VertexAIProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := vertexAIRequest{
		Contents: []vertexAIContent{
			{Role: "user", Parts: []vertexAIPart{{Text: prompt}}},
		},
		GenerationConfig: vertexAIGenConfig{MaxOutputTokens: 8192},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vertex ai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request for vertex ai: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to vertex ai: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read vertex ai response body: %w", err)
	}

	var apiResp vertexAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse vertex ai json response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("vertex ai api error: %s", apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from vertex ai: %d", resp.StatusCode)
	}

	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("received an empty or invalid response from vertex ai")
	}

	return apiResp.Candidates[0].Content.Parts[0].Text, nil
}
