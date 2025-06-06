package main

import (
	"fmt"
	"gct/src/commands"
	"os"

	"github.com/fatih/color"
)

var (
	VerBranch = "Prod."
	VerStatus = "Release"
	VerNumber = "2.1.0"
	VerCommit = "dev"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		commands.PrintUsage()
		return
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		commands.VersionCommand(VerBranch, VerStatus, VerNumber, VerCommit)
		return
	}

	if len(os.Args) > 2 {
		command := os.Args[1]
		subCommand := os.Args[2]

		if command == "ai" {
			switch subCommand {
			case "commit":
				commands.AICommitCommand()
				return
			case "diff":
				commands.AIDiffCommand()
				return
			}
		}

		if command == "commit" && subCommand == "edit" {
			commands.EditCommitCommand()
			return
		}
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "init":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct init (no arguments expected)"))
			return
		}
		commands.InitCommand()
	case "version":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct version (no arguments expected)"))
			return
		}
		commands.VersionCommand(VerBranch, VerStatus, VerNumber, VerCommit)
	case "about":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct about (no arguments expected)"))
			return
		}
		commands.AboutCommand()
	case "update":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct update (no arguments expected)"))
			return
		}
		commands.UpdateCommand(VerBranch, VerStatus, VerNumber)
	case "commit":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct commit (no arguments expected)"))
			return
		}
		commands.CommitCommand()
	case "ai":
		fmt.Printf("%s 'ai' command requires a subcommand.\n", color.RedString("Error:"))
		fmt.Println("Usage: gct ai [commit|diff]")
		return
	default:
		commands.NotFoundCommand()
	}
}
