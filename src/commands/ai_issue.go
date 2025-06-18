package commands

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

const aiIssuePromptTemplate = `
You are a senior developer planning a solution for an issue. Based on the issue's title, body, and labels, propose a clear path forward.

Structure your response into three sections using Markdown headings:
### The Why
Restate the core problem or feature request from the issue. What is the user's goal?

### The How
Propose a high-level technical plan to address the issue. What files might need to be changed? What is the logical approach?

### The Solution
Describe what the completed work will look like from a user's perspective. What will they be able to do that they couldn't before?

--- ISSUE DATA START ---
Title: %s
Author: %s
Labels: %s
Body:
%s
--- ISSUE DATA END ---

Provide your proposed solution below:
`

func AIIssueCommand() {
	red := color.New(color.FgRed).SprintFunc()

	if len(os.Args) < 4 {
		fmt.Printf("%s Issue number is required.\n", red("Error:"))
		fmt.Println("Usage: gct ai issue <number>")
		return
	}
	issueNumber := os.Args[3]

	provider, err := NewGitHostingProvider()
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	details, err := provider.GetIssueDetails(issueNumber)
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	prompt := fmt.Sprintf(aiIssuePromptTemplate, details.Title, details.Author, strings.Join(details.Labels, ", "), details.Body)
	aiResponse, err := runAITask(prompt, false)
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	cleanMsg := strings.TrimSpace(aiResponse)
	viewerModel := NewAITextViewerModel(fmt.Sprintf("🤖 AI Proposed Solution for Issue #%s", issueNumber), cleanMsg)
	p := tea.NewProgram(viewerModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("%s Error displaying AI response: %v\n", red("Error:"), err)
	}
}
