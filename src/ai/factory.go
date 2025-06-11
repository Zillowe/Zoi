package ai

import (
	"fmt"
	"gct/src/config"
	"strings"
)

var SupportedProviders = []string{
	"Google AI Studio",
	"Google Vertex AI",
	"OpenRouter",
	"OpenAI",
	"Azure OpenAI",
	"OpenAI Compatible",
	"Anthropic",
	"DeepSeek",
	"Mistral",
	"Alibaba",
	"Hugging Face",
	"Amazon Bedrock",
	"xAI",
}

func NewProvider(cfg *config.Config) (AIProvider, error) {
	providerName := strings.ToLower(strings.ReplaceAll(cfg.Provider, " ", ""))

	switch providerName {
	case "googleaistudio", "google", "gemini":
		return NewGoogleProvider(cfg.APIKey, cfg.Model)

	case "googlevertexai", "vertexai", "vertex":
		return NewVertexAIProvider(cfg.APIKey, cfg.Model, cfg.GCPProjectID, cfg.GCPRegion)

	case "openrouter":
		return NewOpenRouterProvider(cfg.APIKey, cfg.Model)

	case "openai", "gpt":
		return NewOpenAIProvider(cfg.APIKey, cfg.Model)

	case "azureopenai", "azure":
		return NewAzureProvider(cfg.APIKey, cfg.AzureResourceName, cfg.Model)

	case "openaicompatible":
		return NewOpenAICompatibleProvider(cfg.APIKey, cfg.Model, cfg.Endpoint)

	case "anthropic", "claude":
		return NewAnthropicProvider(cfg.APIKey, cfg.Model)

	case "deepseek":
		endpoint := "https://api.deepseek.com/v1"
		if cfg.Endpoint != "" {
			endpoint = cfg.Endpoint
		}
		return NewOpenAICompatibleProvider(cfg.APIKey, cfg.Model, endpoint)

	case "mistral":
		return NewOpenAICompatibleProvider(cfg.APIKey, cfg.Model, "https://api.mistral.ai/v1")

	case "alibaba", "qwen":
		return NewOpenAICompatibleProvider(cfg.APIKey, cfg.Model, "https://dashscope.aliyuncs.com/api/v1")

	case "huggingface", "hf":
		return NewHuggingFaceProvider(cfg.APIKey, cfg.Model)

	case "amazonbedrock", "bedrock", "amazon":
		return NewBedrockProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.AWSRegion, cfg.Model)

	case "xai", "grok":
		endpoint := "https://api.x.ai/v1"
		if cfg.Endpoint != "" {
			endpoint = cfg.Endpoint
		}
		return NewOpenAICompatibleProvider(cfg.APIKey, cfg.Model, endpoint)

	default:
		return nil, fmt.Errorf(
			"unsupported AI provider: '%s'. Supported providers are: %s",
			cfg.Provider,
			strings.Join(SupportedProviders, ", "),
		)
	}
}
