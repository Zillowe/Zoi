package ai

import (
	"fmt"
	"gct/src/config"
	"strings"
)

var SupportedProviders = []string{
	"Google AI Studio",
	"OpenRouter",
	"OpenAI",
	"Anthropic",
}

func NewProvider(cfg *config.Config) (AIProvider, error) {
	providerName := strings.ToLower(strings.ReplaceAll(cfg.Provider, " ", ""))

	switch providerName {
	case "googleaistudio", "google":
		return NewGoogleProvider(cfg.APIKey, cfg.Model)

	case "openrouter":
		return NewOpenRouterProvider(cfg.APIKey, cfg.Model)

	case "openai":
		return NewOpenAIProvider(cfg.APIKey, cfg.Model)

	case "anthropic":
		return NewAnthropicProvider(cfg.APIKey, cfg.Model)

	default:
		return nil, fmt.Errorf(
			"unsupported AI provider: '%s'. Supported providers are: %s",
			cfg.Provider,
			strings.Join(SupportedProviders, ", "),
		)
	}
}
