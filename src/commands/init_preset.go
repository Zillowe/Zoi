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

func InitPresetCommand() {
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

	tuiModel := NewInitModelTUIModel()
	p := tea.NewProgram(tuiModel, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("%s Error running setup: %v\n", red("✗"), err)
		return
	}

	initModel, ok := finalModel.(InitModelTUIModel)
	if !ok {
		fmt.Printf("%s Could not read setup results. This is an unexpected error.\n", red("✗"))
		return
	}

	if !initModel.submitted {
		fmt.Println("\nInitialization cancelled.")
		return
	}

	commitGuidePaths := []string{}
	if trimmed := strings.TrimSpace(initModel.CommitGuides); trimmed != "" {
		commitGuidePaths = strings.Fields(trimmed)
	}
	changelogGuidePaths := []string{}
	if trimmed := strings.TrimSpace(initModel.ChangelogGuides); trimmed != "" {
		changelogGuidePaths = strings.Fields(trimmed)
	}

	newConfig := config.Config{
		Name:               initModel.Name,
		Provider:           initModel.Provider,
		Model:              initModel.Model,
		APIKey:             initModel.APIKey,
		Endpoint:           initModel.Endpoint,
		Commits:            config.GuidesConfig{Paths: commitGuidePaths},
		Changelogs:         config.GuidesConfig{Paths: changelogGuidePaths},
		GCPProjectID:       initModel.GCPProjectID,
		GCPRegion:          initModel.GCPRegion,
		AWSRegion:          initModel.AWSRegion,
		AWSAccessKeyID:     initModel.AWSAccessKeyID,
		AWSSecretAccessKey: initModel.AWSSecretAccessKey,
		AzureResourceName:  initModel.AzureResourceName,
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
	err = addPathToGitignore(configFileName)
	if err != nil {
		fmt.Printf("%s Could not automatically update .gitignore: %v\n", yellow("Warning:"), err)
		fmt.Printf("%s Please add '%s' to your .gitignore file manually.\n", yellow("Hint:"), configFileName)
	} else {
		fmt.Printf("%s '%s' was added to your .gitignore file.\n", green("✓"), configFileName)
	}
	fmt.Printf("%s You can now run %s.\n", cyan("›"), yellow("gct ai commit"))
}
