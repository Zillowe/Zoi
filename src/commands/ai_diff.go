package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

const aiDiffPromptTemplate = `
You are an expert code reviewer. Your task is to provide a high-level, human-readable explanation of the following git diff.

Focus on the following aspects:
1.  **Overall Purpose:** What is the main goal of these changes?
2.  **Key Changes:** Describe the most important modifications. What was added, removed, or refactored?
3.  **Potential Impact:** Are there any potential risks, breaking changes, or important considerations for other developers?

Do not describe the changes line-by-line. Structure your response using Markdown for clarity (e.g. headings, bullet points).

--- GIT DIFF START ---
%s
--- GIT DIFF END ---

Provide your expert summary below:
`

func AIDiffCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	var diffCmd *exec.Cmd
	var description string

	args := os.Args[3:]

	switch {
	case len(args) == 0:
		description = "unstaged changes in the working directory"
		diffCmd = exec.Command("git", "diff")
	case len(args) == 1 && args[0] == "--staged":
		description = "staged changes"
		diffCmd = exec.Command("git", "diff", "--staged")
	case len(args) == 1:
		ref := args[0]
		description = fmt.Sprintf("changes between HEAD and '%s'", ref)
		diffCmd = exec.Command("git", "diff", ref)
	default:
		fmt.Printf("%s Invalid arguments for 'ai diff'.\n", red("Error:"))
		fmt.Println("Usage: gct ai diff [--staged | <commit|branch>]")
		return
	}

	fmt.Printf("%s Analyzing %s...\n", cyan("🔍"), description)
	diffOutput, err := diffCmd.Output()
	if err != nil {
		fmt.Printf("%s Failed to get git diff. Is the reference valid?\n", red("Error:"))
		return
	}

	if len(diffOutput) == 0 {
		fmt.Printf("%s No changes found to analyze for %s.\n", green("✓"), description)
		return
	}

	prompt := fmt.Sprintf(aiDiffPromptTemplate, string(diffOutput))
	aiResponse, err := runAITask(prompt, false)
	if err != nil {
		if err.Error() == "operation cancelled by user" {
			fmt.Println(color.YellowString("Diff analysis cancelled."))
		} else {
			fmt.Printf("%s %v\n", color.RedString("Error:"), err)
		}
		return
	}

	cleanMsg := strings.TrimSpace(aiResponse)

	viewerModel := NewAITextViewerModel("🤖 AI Explanation of Changes", cleanMsg)
	p := tea.NewProgram(viewerModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("%s Error displaying AI response: %v\n", color.RedString("Error:"), err)
	}
}
