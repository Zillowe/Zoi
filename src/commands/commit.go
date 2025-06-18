package commands

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

func parseCommitMessage(fullMessage string) (cType, cSubject, body string) {
	parts := strings.SplitN(fullMessage, "\n", 2)
	subjectLine := parts[0]
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}

	subjectParts := strings.SplitN(subjectLine, ": ", 2)
	if len(subjectParts) == 2 {
		cType = subjectParts[0]
		cSubject = subjectParts[1]
	} else {
		cType = "Refactor"
		cSubject = subjectLine
	}
	return
}

func executeGitCommit(subjectLine, body string, isAmend bool) error {
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	escapedSubject := strings.ReplaceAll(subjectLine, "\"", "\\\"")
	var cmdToExecute string

	baseCmd := "git commit"
	if isAmend {
		baseCmd = "git commit --amend"
	}

	if strings.TrimSpace(body) != "" {
		escapedBody := strings.ReplaceAll(body, "\"", "\\\"")
		cmdToExecute = fmt.Sprintf("%s -m \"%s\" -m \"%s\"", baseCmd, escapedSubject, escapedBody)
	} else {
		cmdToExecute = fmt.Sprintf("%s -m \"%s\"", baseCmd, escapedSubject)
	}

	fmt.Printf("\n%s Running git command...\n", cyan("ℹ"))
	err := executeCommand(cmdToExecute)
	if err != nil {
		fmt.Printf("%s Failed to commit: %v\n", red("✗"), err)
		fmt.Printf("%s Command attempted: %s\n", red("↪"), cmdToExecute)
		return err
	}
	return nil
}

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
		return
	}

	_, cSubject, body := parseCommitMessage(string(output))
	cType, cSubject, _ := parseCommitMessage(cSubject)

	fullMessage := strings.Replace(string(output), delimiter, "\n", 1)
	cType, cSubject, body = parseCommitMessage(fullMessage)

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

	commitModel, _ := finalModel.(CommitTUIModel)
	if !commitModel.submitted || (commitModel.CommitType == "" && commitModel.Subject == "") {
		if !commitModel.quitting {
			fmt.Printf("%s Process cancelled.\n", color.YellowString("!"))
		}
		return
	}

	commitSubjectLine := fmt.Sprintf("%s: %s", commitModel.CommitType, commitModel.Subject)
	err = executeGitCommit(commitSubjectLine, commitModel.Body, isAmend)
	if err != nil {
		return
	}

	successMsg := "Commit successful!"
	if isAmend {
		successMsg = "Amend successful!"
	}
	fmt.Printf("\n%s %s\n", color.GreenString("✓"), successMsg)
}
