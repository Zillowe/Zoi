package commands

import (
	"context"
	"fmt"
	"gct/src/ai"
	"gct/src/config"

	"github.com/fatih/color"
)

func runAITask(prompt string) (string, error) {
	cyan := color.New(color.FgCyan).SprintFunc()

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	provider, err := ai.NewProvider(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	fmt.Println(cyan("[Thinking]..."))
	ctx := context.Background()
	generatedText, err := provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	return generatedText, nil
}
