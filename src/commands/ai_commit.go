package commands

import (
	"fmt"
	"gct/src/config"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

const aiEditCommitPromptTemplate = `
You are a helpful assistant revising a git commit message based on user feedback.
The user has provided an original commit message and an instruction on how to change it.

Here is the original commit message:
--- ORIGINAL MESSAGE START ---
%s
--- ORIGINAL MESSAGE END ---

Here is the user's instruction for the change:
--- USER INSTRUCTION START ---
%s
--- USER INSTRUCTION END ---

Your task is to generate the full, revised commit message that incorporates the user's instruction.
Maintain the original conventional commit format (e.g. "Type: Subject").
ONLY output the raw, complete, revised commit message. Do not add any extra commentary.
`

const aiCommitPromptTemplateWithContext = `
You are an expert programmer creating a commit message.
Your task is to generate a concise, conventional commit message based on the provided guidelines, staged code changes, and any additional context.

Adhere strictly to the following guidelines:
--- GUIDELINES START ---
%s
--- GUIDELINES END ---

Here is the additional context provided by the user. Incorporate this information into the commit message body where appropriate (e.g. for co-authorship, issue numbers, or specific explanations):
--- ADDITIONAL CONTEXT START ---
%s
--- ADDITIONAL CONTEXT END ---

Here are the staged changes (git diff):
--- GIT DIFF START ---
%s
--- GIT DIFF END ---

Based on all the information above, generate the complete commit message.
The message must have a subject line, a blank line, and then the body.
ONLY output the raw commit message itself, without any extra commentary, introductory text, or markdown formatting like backticks.
`

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

func AICommitCommand(additionalContext string) {
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
	guidelines, _ := readGuidelines(cfg.Commits.Paths)

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

	var prompt string
	if additionalContext != "" {
		fmt.Println(cyan("✍️ Applying additional user context..."))
		prompt = fmt.Sprintf(aiCommitPromptTemplateWithContext, guidelines, additionalContext, string(diffOutput))
	} else {
		prompt = fmt.Sprintf(aiCommitPromptTemplate, guidelines, string(diffOutput))
	}

	initialGeneratedMsg, err := runAITask(prompt, false)
	if err != nil {
		if err.Error() == "operation cancelled by user" {
			fmt.Println(color.YellowString("Commit cancelled."))
		} else {
			fmt.Printf("%s %v\n", color.RedString("Error:"), err)
		}
		return
	}

	currentMessage := strings.TrimSpace(initialGeneratedMsg)
	currentMessage = strings.Trim(currentMessage, "`")

	for {
		fmt.Printf("\n%s AI Generated Commit Message:\n", cyan("🤖"))
		fmt.Printf("%s\n%s\n%s\n", yellow("--- Start ---"), green(currentMessage), yellow("--- End ---"))

		action := promptForAction("Press [c] to chat/change, [e] to edit, [Enter] to commit, [q] to quit:")

		switch action {
		case 'c':
			changeRequest := promptForInput("What change would you like to make?")
			if changeRequest == "" {
				fmt.Println(yellow("No change request provided. Please try again."))
				continue
			}

			revisionPrompt := fmt.Sprintf(aiEditCommitPromptTemplate, currentMessage, changeRequest)

			revisedMsg, err := runAITask(revisionPrompt, true)
			if err != nil {
				fmt.Printf("%s %v\n", red("Error:"), err)
				continue
			}

			currentMessage = strings.TrimSpace(revisedMsg)
			currentMessage = strings.Trim(currentMessage, "`")
			continue

		case 'e':
			fmt.Println(cyan("\n✍️ Opening editor..."))
			cType, cSubject, body := parseCommitMessage(currentMessage)
			tuiModel := NewCommitTUIModel(cType, cSubject, body)
			runCommitTUI(tuiModel, false)
			return

		case '\n':
			subjectLine, body, _ := strings.Cut(currentMessage, "\n\n")
			err := executeGitCommit(subjectLine, body, false)
			if err != nil {
				return
			}
			fmt.Printf("\n%s AI Commit successful!\n", green("✓"))
			return

		case 'q':
			fmt.Println(yellow("\nCommit cancelled."))
			return
		}
	}
}
