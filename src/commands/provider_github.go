package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type GitHubProvider struct{}

func NewGitHubProvider() (*GitHubProvider, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil, fmt.Errorf("'gh' (GitHub CLI) is not installed or not in your PATH")
	}
	return &GitHubProvider{}, nil
}

func (p *GitHubProvider) GetPRDetails(prNumber string) (*PRDetails, error) {
	cmd := exec.Command("gh", "pr", "view", prNumber, "--json", "title,body,author,diff")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get PR details from gh CLI: %w. Is the PR number correct?", err)
	}

	var ghPR struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		Author struct {
			Login string `json:"login"`
		} `json:"author"`
		Diff string `json:"diff"`
	}

	if err := json.Unmarshal(output, &ghPR); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from gh CLI: %w", err)
	}

	return &PRDetails{
		Title:  ghPR.Title,
		Body:   ghPR.Body,
		Author: ghPR.Author.Login,
		Diff:   ghPR.Diff,
	}, nil
}

func (p *GitHubProvider) GetIssueDetails(issueNumber string) (*IssueDetails, error) {
	cmd := exec.Command("gh", "issue", "view", issueNumber, "--json", "title,body,author,labels")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get issue details from gh CLI: %w", err)
	}

	var ghIssue struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		Author struct {
			Login string `json:"login"`
		} `json:"author"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
	}

	if err := json.Unmarshal(output, &ghIssue); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from gh CLI: %w", err)
	}

	var labelNames []string
	for _, l := range ghIssue.Labels {
		labelNames = append(labelNames, l.Name)
	}

	return &IssueDetails{
		Title:  ghIssue.Title,
		Body:   ghIssue.Body,
		Author: ghIssue.Author.Login,
		Labels: labelNames,
	}, nil
}
