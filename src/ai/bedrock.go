package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockProvider struct {
	client *bedrockruntime.Client
	model  string
}

type bedrockClaudeRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	StopSequences     []string `json:"stop_sequences"`
}

type bedrockClaudeResponse struct {
	Completion string `json:"completion"`
}

func NewBedrockProvider(accessKey, secretKey, region, modelID string) (*BedrockProvider, error) {
	if accessKey == "" || secretKey == "" || region == "" {
		return nil, fmt.Errorf("AWS Access Key, Secret Key, and Region are required for Bedrock")
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	return &BedrockProvider{
		client: bedrockruntime.NewFromConfig(cfg),
		model:  modelID,
	}, nil
}

func (p *BedrockProvider) Generate(ctx context.Context, prompt string) (string, error) {
	claudePrompt := fmt.Sprintf("\n\nHuman: %s\n\nAssistant:", prompt)

	payload := bedrockClaudeRequest{
		Prompt:            claudePrompt,
		MaxTokensToSample: 8192,
		StopSequences:     []string{"\n\nHuman:"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal bedrock request: %w", err)
	}

	output, err := p.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(p.model),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke bedrock model: %w", err)
	}

	var resp bedrockClaudeResponse
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return "", fmt.Errorf("failed to unmarshal bedrock response: %w", err)
	}

	return resp.Completion, nil
}
