package commands

import (
	"context"
	"fmt"
	"gct/src/ai"
	"gct/src/config"
	"os"
	"strings"

	"github.com/fatih/color"
)

const tokenWarningThreshold = 4000
const charsPerToken = 4

func runAITask(prompt string, isSilent bool) (string, error) {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	if !isSilent {
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
		fmt.Printf("%s Using %s from %s\n", cyan("›"), magenta(cfg.Model), magenta(cfg.Provider))
		fmt.Println(cyan("[Thinking]..."))
	}

	provider, err := ai.NewProvider(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	ctx := context.Background()
	generatedText, err := provider.Generate(ctx, prompt)

	if !isSilent {
		fmt.Printf("\r%s\n", green("✓ Done!                     "))
	}

	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	return generatedText, nil
}

func readGuidelines(paths []string) (string, error) {
	yellow := color.New(color.FgYellow).SprintFunc()
	var guidelines strings.Builder

	for _, path := range paths {
		if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".txt") {
			fmt.Printf("%s Skipping unsupported guide file: %s\n", yellow("Warning:"), path)
			continue
		}
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("%s Could not read guide file %s: %v\n", yellow("Warning:"), path, err)
			continue
		}
		guidelines.Write(content)
		guidelines.WriteString("\n---\n")
	}
	return guidelines.String(), nil
}
