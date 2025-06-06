package commands

import (
	"fmt"
	"gct/src/config"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

const aiCommitPromptTemplate = `
You are an expert programmer creating a commit message.
Your task is to generate a concise, conventional commit message based on the provided guidelines and staged code changes.

Adhere strictly to the following guidelines:
--- GUIDELINES START ---
%s
--- GUIDELINES END ---

Here are the staged changes (git diff):
--- GIT DIFF START ---
%s
--- GIT DIFF END ---

Based on the guidelines and the diff, generate the complete commit message.
The message must have a subject line, a blank line, and then the body.
ONLY output the raw commit message itself, without any extra commentary, introductory text, or markdown formatting like backticks.
`

func AICommitCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Println(cyan("🔍 Loading configuration..."))
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	fmt.Println(cyan("📚 Reading commit guidelines..."))
	var guidelines strings.Builder
	for _, guidePath := range cfg.Guides {
		if !strings.HasSuffix(guidePath, ".md") && !strings.HasSuffix(guidePath, ".txt") {
			fmt.Printf("%s Skipping unsupported guide file: %s\n", yellow("Warning:"), guidePath)
			continue
		}
		content, err := os.ReadFile(guidePath)
		if err != nil {
			fmt.Printf("%s Could not read guide file %s: %v\n", yellow("Warning:"), guidePath, err)
			continue
		}
		guidelines.Write(content)
		guidelines.WriteString("\n")
	}

	fmt.Println(cyan("📝 Analyzing staged changes..."))
	diffCmd := exec.Command("git", "diff", "--staged")
	diffOutput, err := diffCmd.Output()
	if err != nil {
		fmt.Printf("%s Failed to get git diff: %v\n", red("Error:"), err)
		return
	}
	if len(diffOutput) == 0 {
		fmt.Println(yellow("No changes are staged. Nothing to commit."))
		return
	}

	prompt := fmt.Sprintf(aiCommitPromptTemplate, guidelines.String(), string(diffOutput))

	fmt.Println(cyan("[Thinking]..."))
	generatedMsg, err := runAITask(prompt)
	if err != nil {
		if err.Error() == "operation cancelled by user" {
			fmt.Println(color.YellowString("Commit cancelled."))
		} else {
			fmt.Printf("%s %v\n", color.RedString("Error:"), err)
		}
		return
	}

	cleanMsg := strings.TrimSpace(generatedMsg)
	cleanMsg = strings.Trim(cleanMsg, "`")

	fmt.Printf("\n%s AI Generated Commit Message:\n", cyan("🤖"))
	fmt.Printf("%s\n%s\n%s\n", yellow("--- Start ---"), green(cleanMsg), yellow("--- End ---"))

	if !confirmPrompt("Use this commit message?") {
		fmt.Println(yellow("Commit cancelled."))
		return
	}

	parts := strings.SplitN(cleanMsg, "\n", 2)
	subject := parts[0]
	body := ""
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}

	escapedSubject := strings.ReplaceAll(subject, "\"", "\\\"")
	var cmdToExecute string
	if body != "" {
		escapedBody := strings.ReplaceAll(body, "\"", "\\\"")
		cmdToExecute = fmt.Sprintf("git commit -m \"%s\" -m \"%s\"", escapedSubject, escapedBody)
	} else {
		cmdToExecute = fmt.Sprintf("git commit -m \"%s\"", escapedSubject)
	}

	gitErr := executeCommand(cmdToExecute)
	if gitErr != nil {
		fmt.Printf("%s Failed to commit: %v\n", red("Error:"), gitErr)
		return
	}

	fmt.Printf("\n%s AI Commit successful!\n", green("✓"))
}
