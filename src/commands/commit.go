package commands

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

func CommitCommand() {
	tuiModel := NewCommitTUIModel("", "", "")
	runCommitTUI(tuiModel, false)
}

func EditCommitCommand() {
	const delimiter = "|||---GIT-LOG-SPLITTER---|||"
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%s"+delimiter+"%b")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s Could not get last commit: %v\n", color.RedString("✗"), err)
		fmt.Println("Is this a git repository with at least one commit?")
		return
	}

	fullMessage := string(output)
	parts := strings.SplitN(fullMessage, delimiter, 2)
	if len(parts) != 2 {
		fmt.Printf("%s Failed to parse last commit message.\n", color.RedString("✗"))
		return
	}
	subjectLine := parts[0]
	body := parts[1]

	var cType, cSubject string
	subjectParts := strings.SplitN(subjectLine, ": ", 2)
	if len(subjectParts) == 2 {
		cType = subjectParts[0]
		cSubject = subjectParts[1]
	} else {
		cType = "Refactor"
		cSubject = subjectLine
	}

	fmt.Println(color.CyanString("✍️  Editing last commit message..."))

	tuiModel := NewCommitTUIModel(cType, cSubject, body)
	runCommitTUI(tuiModel, true)
}

func runCommitTUI(tuiModel CommitTUIModel, isAmend bool) {
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
		} else if !commitModel.quitting {
			fmt.Printf("%s Process cancelled.\n", color.YellowString("!"))
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

	fmt.Printf("\n%s This is your commit message:\n", cyan("ⓘ"))
	fmt.Printf("%s\n%s\n%s\n", yellow("--- Start ---"), green(finalCommitMsgForDisplay), yellow("--- End ---"))

	confirmMsg := "Are you happy with this commit message?"
	if isAmend {
		confirmMsg = "Are you happy with this amended message?"
	}
	if !confirmPrompt(confirmMsg) {
		fmt.Printf("%s Operation cancelled.\n", yellow("!"))
		return
	}

	escapedSubject := strings.ReplaceAll(commitSubjectLine, "\"", "\\\"")
	var cmdToExecute string

	baseCmd := "git commit"
	if isAmend {
		baseCmd = "git commit --amend"
	}

	if trimmedBody != "" {
		escapedBody := strings.ReplaceAll(commitBody, "\"", "\\\"")
		cmdToExecute = fmt.Sprintf("%s -m \"%s\" -m \"%s\"", baseCmd, escapedSubject, escapedBody)
	} else {
		cmdToExecute = fmt.Sprintf("%s -m \"%s\"", baseCmd, escapedSubject)
	}

	fmt.Printf("\n%s Running git command...\n", cyan("ℹ"))
	gitErr := executeCommand(cmdToExecute)
	if gitErr != nil {
		fmt.Printf("%s Failed to commit: %v\n", red("✗"), gitErr)
		fmt.Printf("%s Command attempted: %s\n", red("↪"), cmdToExecute)
		return
	}

	successMsg := "Commit successful!"
	if isAmend {
		successMsg = "Amend successful!"
	}
	fmt.Printf("\n%s %s\n", green("✓"), successMsg)
}
