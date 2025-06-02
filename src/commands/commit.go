package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func CommitCommand() {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s Welcome to the GCT Commit Tool!\n", cyan("\n✨"))

	typeInput := promptForInput("Enter commit Type (e.g. ✨ Feat, Fix, Chore)")
	for strings.TrimSpace(typeInput) == "" {
		fmt.Printf("%s Type cannot be empty.\n", red("✗"))
		typeInput = promptForInput("Enter commit Type (e.g. ✨ Feat, Fix, Chore)")
	}

	subjectText := promptForInput("Enter Subject")
	for strings.TrimSpace(subjectText) == "" {
		fmt.Printf("%s Subject cannot be empty.\n", red("✗"))
		subjectText = promptForInput("Enter Subject")
	}

	commitSubjectLine := fmt.Sprintf("%s: %s", typeInput, subjectText)

	fmt.Printf("%s Enter Body (optional, press Enter twice on consecutive empty lines to finish):\n", cyan("?"))
	bodyLines := []string{}
	reader := bufio.NewReader(os.Stdin)
	consecutiveEmptyLineInputs := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")

		if line == "" {
			consecutiveEmptyLineInputs++
			if consecutiveEmptyLineInputs == 2 {
				break
			}
			bodyLines = append(bodyLines, line)
		} else {
			consecutiveEmptyLineInputs = 0
			bodyLines = append(bodyLines, line)
		}
	}

	if consecutiveEmptyLineInputs == 2 && len(bodyLines) > 0 && bodyLines[len(bodyLines)-1] == "" {
		bodyLines = bodyLines[:len(bodyLines)-1]
	}

	commitBody := strings.Join(bodyLines, "\n")
	trimmedBody := strings.TrimSpace(commitBody)

	var displayCommitMessage strings.Builder
	displayCommitMessage.WriteString(commitSubjectLine)
	if trimmedBody != "" {
		displayCommitMessage.WriteString("\n\n")
		displayCommitMessage.WriteString(commitBody)
	}
	finalCommitMsgForDisplay := displayCommitMessage.String()

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

	err := executeCommand(cmdToExecute)
	if err != nil {
		fmt.Printf("%s Failed to commit: %v\n", red("✗"), err)
		fmt.Printf("%s Command attempted: %s\n", red("↪"), cmdToExecute)
		return
	}

	fmt.Printf("\n%s Commit successful!\n", green("✓"))
}
