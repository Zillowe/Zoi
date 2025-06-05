package commands

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintUsage() {
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	faint := color.New(color.Faint).SprintFunc()

	fmt.Printf("%s - A smart, interactive Git tool powered by AI\n\n", bold("GCT"))
	fmt.Printf("%s\n", yellow("USAGE"))
	fmt.Printf("  gct %s\n\n", green("<command> [arguments...]"))

	fmt.Printf("%s\n", yellow("CORE COMMANDS"))
	fmt.Printf("  %s          Interactively create the 'gct.yaml' config file\n", green("%-15s", "init"))
	fmt.Printf("  %s           Show GCT version information\n", green("%-15s", "version"))
	fmt.Printf("  %s          Display details and information about GCT\n", green("%-15s", "about"))
	fmt.Printf("  %s         Check for and apply updates to GCT itself\n", green("%-15s", "update"))
	fmt.Printf("  %s              Show this help message\n\n", green("%-15s", "help"))

	fmt.Printf("%s\n", yellow("MANUAL GIT COMMANDS"))
	fmt.Printf("  %s            Create a new git commit using an interactive form\n", green("%-15s", "commit"))
	fmt.Printf("  %s       Edit the previous commit's message interactively\n\n", green("%-15s", "commit edit"))

	fmt.Printf("%s\n", yellow("AI GIT COMMANDS"))
	fmt.Printf("  %s     Generate a commit message using AI based on staged changes\n", green("%-15s", "ai commit"))
	fmt.Printf("  %s       Explain code changes using AI\n", green("%-1s", "ai diff [args]"))
	fmt.Printf("    %s %s\n", faint("└─"), "Explain unstaged changes")
	fmt.Printf("    %s %s\n", faint("  └─"), green("--staged"))
	fmt.Printf("    %s %s\n\n", faint("  └─"), green("<commit|branch>"))

	fmt.Printf("%s\n", yellow("GLOBAL FLAGS"))
	fmt.Printf("  %s, %s     Show GCT version information\n", green("-v"), green("--version"))
}
