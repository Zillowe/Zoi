package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GoogleProvider struct {
	client *genai.GenerativeModel
}

func NewGoogleProvider(apiKey, modelName string) (*GoogleProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Google AI Studio API key is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google AI client: %w", err)
	}

	model := client.GenerativeModel(modelName)
	return &GoogleProvider{client: model}, nil
}

func (g *GoogleProvider) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := g.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content from Google AI: %w", err)
	}

	var result strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					result.WriteString(string(txt))
				}
			}
		}
	}

	if result.Len() == 0 {
		return "", fmt.Errorf("received an empty response from the AI")
	}

	return result.String(), nil
}
