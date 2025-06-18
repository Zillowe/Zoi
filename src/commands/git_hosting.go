package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

type PRDetails struct {
	Title  string
	Body   string
	Author string
	Diff   string
}

type IssueDetails struct {
	Title  string
	Body   string
	Author string
	Labels []string
}

type GitHostingProvider interface {
	GetPRDetails(prNumber string) (*PRDetails, error)
	GetIssueDetails(issueNumber string) (*IssueDetails, error)
}

func NewGitHostingProvider() (GitHostingProvider, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not determine git remote 'origin'. Are you in a git repository?")
	}
	remoteURL := strings.TrimSpace(string(output))

	if strings.Contains(remoteURL, "github.com") {
		return NewGitHubProvider()
	} else if strings.Contains(remoteURL, "codeberg.org") || strings.Contains(remoteURL, "gitea.com") || strings.Contains(remoteURL, "code.forgejo.org") {
		return NewForgejoProvider()
	}

	return nil, fmt.Errorf("unsupported git hosting platform for remote: %s", remoteURL)
}
