package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

type ForgejoProvider struct{}

func NewForgejoProvider() (*ForgejoProvider, error) {
	if _, err := exec.LookPath("fj"); err != nil {
		return nil, fmt.Errorf("'fj' (Forgejo CLI) is not installed or not in your PATH")
	}
	return &ForgejoProvider{}, nil
}

func (p *ForgejoProvider) GetPRDetails(prNumber string) (*PRDetails, error) {
	cmd := exec.Command("fj", "pr", "view", prNumber)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get PR details from fj CLI: %w", err)
	}

	content := string(output)
	title := strings.SplitN(content, "\n", 2)[0]

	return &PRDetails{
		Title:  title,
		Body:   content,
		Author: "unknown",
		Diff:   "",
	}, nil
}

func (p *ForgejoProvider) GetIssueDetails(issueNumber string) (*IssueDetails, error) {
	cmd := exec.Command("fj", "issue", "view", issueNumber)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get issue details from fj CLI: %w", err)
	}

	content := string(output)
	title := strings.SplitN(content, "\n", 2)[0]

	return &IssueDetails{
		Title:  title,
		Body:   content,
		Author: "unknown",
		Labels: []string{},
	}, nil
}
