package commands

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

const aiPRPromptTemplate = `
You are a senior software engineer summarizing a pull request for a team member.
Based on the PR's title, body, and code diff, provide a clear and concise explanation.

Structure your response into three sections using Markdown headings:
### The Why
Explain the problem or feature request this PR addresses. What was the goal?

### The How
Describe the technical approach taken in this PR. How does the code solve the problem?

### The Solution
Detail the outcome for the user or the system. What new capabilities are enabled or what bugs are fixed?

--- PR DATA START ---
Title: %s
Author: %s
Body:
%s

Diff:
%s
--- PR DATA END ---

Provide your expert summary below:
`

func AIPRCommand() {
	red := color.New(color.FgRed).SprintFunc()

	if len(os.Args) < 4 {
		fmt.Printf("%s PR number is required.\n", red("Error:"))
		fmt.Println("Usage: gct ai pr <number>")
		return
	}
	prNumber := os.Args[3]

	provider, err := NewGitHostingProvider()
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	details, err := provider.GetPRDetails(prNumber)
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	prompt := fmt.Sprintf(aiPRPromptTemplate, details.Title, details.Author, details.Body, details.Diff)
	aiResponse, err := runAITask(prompt, false)
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	cleanMsg := strings.TrimSpace(aiResponse)
	viewerModel := NewAITextViewerModel(fmt.Sprintf("🤖 AI Summary of PR #%s", prNumber), cleanMsg)
	p := tea.NewProgram(viewerModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("%s Error displaying AI response: %v\n", red("Error:"), err)
	}
}
