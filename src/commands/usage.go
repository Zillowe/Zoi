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
	fmt.Printf("  %-18s          Interactively create the 'gct.yaml' config file\n", green("init"))
	fmt.Printf("    %s %s\n", faint("└─"), "Interactively create the 'gct.yaml' config file using preset of models")
	fmt.Printf("    %s %s\n", faint("  └─"), green("model"))
	fmt.Printf("  %-18s          Show GCT version information\n", green("version"))
	fmt.Printf("  %-18s          Display details and information about GCT\n", green("about"))
	fmt.Printf("  %-18s          Check for and apply updates to GCT itself\n", green("update"))
	fmt.Printf("  %-18s          Show this help message\n\n", green("help"))

	fmt.Printf("%s\n", yellow("MANUAL GIT COMMANDS"))
	fmt.Printf("  %-18s          Create a new git commit using an interactive form\n", green("commit"))
	fmt.Printf("  %-18s        Edit the previous commit's message interactively\n\n", green("commit edit"))

	fmt.Printf("%s\n", yellow("AI GIT COMMANDS"))
	fmt.Printf("  %-18s          Generate and conversationally refine a commit message\n", green("ai commit"))
	fmt.Printf("  %-18s     Explain code changes using AI\n", green("ai diff [args]"))
	fmt.Printf("    %s %s\n", faint("└─"), "Explain unstaged changes")
	fmt.Printf("    %s %s\n", faint("  └─"), green("--staged"))
	fmt.Printf("    %s %s\n", faint("  └─"), green("<commit|branch>"))
	fmt.Printf("  %-18s      Generate a changelog entry from code changes\n", green("ai log [args]"))
	fmt.Printf("    %s %s\n", faint("└─"), "For unstaged changes")
	fmt.Printf("    %s %s\n", faint("  └─"), green("--staged"))
	fmt.Printf("    %s %s\n", faint("  └─"), green("<commit|branch>"))
	fmt.Printf("  %-18s     Summarize a pull request\n", green("ai pr <number>"))
	fmt.Printf("  %-18s  Propose a solution for an issue\n\n", green("ai issue <number>"))

	fmt.Printf("%s\n", yellow("GLOBAL FLAGS"))
	fmt.Printf("  %s, %s      Show GCT version information\n", green("-v"), green("--version"))
}
