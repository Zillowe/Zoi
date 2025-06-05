package ai

import (
	"fmt"
	"gct/src/config"
	"strings"
)

func NewProvider(cfg *config.Config) (AIProvider, error) {
	providerName := strings.ToLower(strings.ReplaceAll(cfg.Provider, " ", ""))

	switch providerName {
	case "googleaistudio", "google":
		return NewGoogleProvider(cfg.APIKey, cfg.Model)
	case "openai":
		return nil, fmt.Errorf("provider '%s' is not yet implemented", cfg.Provider)
	case "anthropic":
		return nil, fmt.Errorf("provider '%s' is not yet implemented", cfg.Provider)
	case "openrouter":
		return nil, fmt.Errorf("provider '%s' is not yet implemented", cfg.Provider)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.Provider)
	}
}
