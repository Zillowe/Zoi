package commands

import (
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
	fmt.Printf("%s You can now run %s to generate commit messages.\n", cyan("›"), yellow("gct ai commit"))
}
