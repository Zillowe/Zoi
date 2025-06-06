package commands

import (
	"context"
	"fmt"
	"gct/src/ai"
	"gct/src/config"

	"github.com/fatih/color"
)

const tokenWarningThreshold = 4000

const charsPerToken = 4

func runAITask(prompt string) (string, error) {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	estimatedTokens := len(prompt) / charsPerToken

	fmt.Printf("%s  Estimated tokens: ~%d\n", cyan("ℹ"), estimatedTokens)

	if estimatedTokens > tokenWarningThreshold {
		warningMsg := fmt.Sprintf(
			"The input is large (~%d tokens). This may result in a long wait or higher costs.",
			estimatedTokens,
		)
		fmt.Printf("%s %s\n", yellow("Warning:"), warningMsg)
		if !confirmPrompt("Do you want to continue?") {
			return "", fmt.Errorf("operation cancelled by user")
		}
	}

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

	fmt.Printf("\r%s\n", green("✓ Done!        "))

	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	return generatedText, nil
}
