package commands

import (
	"bufio"
	"fmt"
	"gct/src/config"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

const configFileName = "gct.yaml"

func InitCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if _, err := os.Stat(configFileName); err == nil {
		fmt.Printf("%s Config file '%s' already exists.\n", yellow("Warning:"), configFileName)
		if !confirmPrompt("Do you want to overwrite it?") {
			fmt.Println("Initialization cancelled.")
			return
		}
	}

	tuiModel := NewInitTUIModel()
	p := tea.NewProgram(tuiModel, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("%s Error running setup: %v\n", red("✗"), err)
		return
	}

	initModel, ok := finalModel.(InitTUIModel)
	if !ok {
		fmt.Printf("%s Could not read setup results. This is an unexpected error.\n", red("✗"))
		return
	}

	if !initModel.submitted {
		fmt.Println("\nInitialization cancelled.")
		return
	}

	var guidePaths []string
	trimmedGuides := strings.TrimSpace(initModel.Guides)
	if trimmedGuides != "" {

		guidePaths = strings.Fields(trimmedGuides)
	} else {
		guidePaths = []string{}
	}

	newConfig := config.Config{
		Name:     initModel.Name,
		Provider: initModel.Provider,
		Model:    initModel.Model,
		APIKey:   initModel.APIKey,
		Guides:   guidePaths,
		Endpoint: initModel.Endpoint,
	}

	yamlData, err := yaml.Marshal(&newConfig)
	if err != nil {
		fmt.Printf("%s Failed to create YAML config: %v\n", red("✗"), err)
		return
	}

	err = os.WriteFile(configFileName, yamlData, 0644)
	if err != nil {
		fmt.Printf("%s Failed to write %s: %v\n", red("✗"), configFileName, err)
		return
	}

	fmt.Printf("\n%s Config file '%s' created successfully!\n", green("✓"), configFileName)

	err = addPathToGitignore(configFileName, "/.gct/")
	if err != nil {
		fmt.Printf("%s Could not automatically update .gitignore: %v\n", yellow("Warning:"), err)
		fmt.Printf("%s Please add '%s' to your .gitignore file manually to protect your API key.\n", yellow("Hint:"), configFileName)
	} else {
		fmt.Printf("%s '%s' was added to your .gitignore file.\n", green("✓"), configFileName)
	}

	fmt.Printf("%s You can now run %s to generate commit messages.\n", cyan("›"), yellow("gct ai commit"))
}

func addPathToGitignore(pathsToAdd ...string) error {
	const gitignoreFileName = ".gitignore"

	file, err := os.OpenFile(gitignoreFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create .gitignore: %w", err)
	}
	defer file.Close()

	existingEntries := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		existingEntries[strings.TrimSpace(scanner.Text())] = true
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}

	var newContent strings.Builder
	for _, path := range pathsToAdd {
		if !existingEntries[path] {
			var comment string
			if path == "gct.yaml" {
				comment = "\n\n# Added by GCT\n"
			} else if path == "/.gct/" {
				comment = ""
			}
			newContent.WriteString(fmt.Sprintf("%s%s\n", comment, path))
		}
	}

	if newContent.Len() > 0 {
		if _, err := file.WriteString(newContent.String()); err != nil {
			return fmt.Errorf("failed to write to .gitignore: %w", err)
		}
	}

	return nil
}
