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

const azureAPIVersion = "2024-02-01"

type AzureProvider struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

type azureRequest struct {
	Model    string         `json:"model"`
	Messages []azureMessage `json:"messages"`
}
type azureMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type azureResponse struct {
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

func NewAzureProvider(apiKey, resourceName, deploymentName string) (*AzureProvider, error) {
	if apiKey == "" || resourceName == "" || deploymentName == "" {
		return nil, fmt.Errorf("API Key, Azure Resource Name, and Deployment Name (in 'model' field) are all required for Azure OpenAI")
	}

	url := fmt.Sprintf("https://%s.openai.azure.com/openai/deployments/%s/chat/completions?api-version=%s",
		resourceName, deploymentName, azureAPIVersion)

	return &AzureProvider{
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: url,
	}, nil
}

func (p *AzureProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := azureRequest{
		Model: "",
		Messages: []azureMessage{
			{Role: "user", Content: prompt},
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal azure request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request for azure: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to azure: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read azure response body: %w", err)
	}

	var apiResp azureResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse azure json response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if apiResp.Error != nil {
			return "", fmt.Errorf("azure api error (type: %s): %s", apiResp.Error.Type, apiResp.Error.Message)
		}
		return "", fmt.Errorf("received non-200 status from azure: %s", string(respBody))
	}

	if len(apiResp.Choices) == 0 || apiResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("received an empty or invalid response from azure")
	}

	return apiResp.Choices[0].Message.Content, nil
}
