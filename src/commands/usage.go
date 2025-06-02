package commands

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintUsage() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("%s - Universal Environment Setup Tool\n\n", cyan("GCT"))
	fmt.Printf("%s\n", yellow("Usage:"))
	fmt.Printf("  gct %s\n\n", green("<command> [arguments...]"))

	fmt.Printf("%s\n", cyan("Available Commands:"))
	fmt.Printf("  %s           Show GCT version information\n", green("version"))
	fmt.Printf("  %s             Display details and information about GCT\n", green("about"))
	fmt.Printf("  %s            Check for and apply updates to GCT itself\n", green("update"))
	fmt.Printf("  %s            Create a new git commit interactively\n", green("commit"))
	fmt.Printf("  %s              Show this help message\n\n", green("help"))

	fmt.Printf("%s\n", cyan("Flags:"))
	fmt.Printf("  %s, %s     Show GCT version information\n", green("-v"), green("--version"))

	fmt.Printf("\n%s Run %s for more details on a specific command.\n", yellow("Hint:"), green("'gct <command> --help'"))
}
