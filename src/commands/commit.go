package commands

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

func CommitCommand() {
	tuiModel := NewCommitTUIModel()

	p := tea.NewProgram(tuiModel, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		return
	}

	commitModel, ok := finalModel.(CommitTUIModel)
	if !ok {
		fmt.Println("Could not cast final TUI model. This is an unexpected error.")
		return
	}

	if !commitModel.submitted || (commitModel.CommitType == "" && commitModel.Subject == "") {
		if commitModel.err != nil {
			fmt.Printf("%s TUI Error: %v\n", color.RedString("✗"), commitModel.err)
		} else {
			if !commitModel.quitting {
				fmt.Printf("%s Commit process cancelled.\n", color.YellowString("!"))
			}
		}
		return
	}

	commitSubjectLine := fmt.Sprintf("%s: %s", commitModel.CommitType, commitModel.Subject)
	commitBody := commitModel.Body
	trimmedBody := strings.TrimSpace(commitBody)

	var displayCommitMessage strings.Builder
	displayCommitMessage.WriteString(commitSubjectLine)
	if trimmedBody != "" {
		displayCommitMessage.WriteString("\n\n")
		displayCommitMessage.WriteString(commitBody)
	}
	finalCommitMsgForDisplay := displayCommitMessage.String()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("\n%s This is your generated commit message (as Git will store it):\n", cyan("ⓘ"))
	fmt.Printf("%s\n%s\n%s\n", yellow("--- Start of commit message ---"), green(finalCommitMsgForDisplay), yellow("--- End of commit message ---"))

	if !confirmPrompt("Are you happy with this commit message?") {
		fmt.Printf("%s Commit cancelled.\n", yellow("!"))
		return
	}

	fmt.Printf("\n%s Running git commit...\n", cyan("ℹ"))

	escapedSubject := strings.ReplaceAll(commitSubjectLine, "\"", "\\\"")

	var cmdToExecute string
	if trimmedBody != "" {
		escapedBody := strings.ReplaceAll(commitBody, "\"", "\\\"")
		cmdToExecute = fmt.Sprintf("git commit -m \"%s\" -m \"%s\"", escapedSubject, escapedBody)
	} else {
		cmdToExecute = fmt.Sprintf("git commit -m \"%s\"", escapedSubject)
	}

	gitErr := executeCommand(cmdToExecute)
	if gitErr != nil {
		fmt.Printf("%s Failed to commit: %v\n", red("✗"), gitErr)
		fmt.Printf("%s Command attempted: %s\n", red("↪"), cmdToExecute)
		return
	}

	fmt.Printf("\n%s Commit successful!\n", green("✓"))
}
