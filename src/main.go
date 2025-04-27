package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	VerBranch = "Dev."
	VerStatus = "Pre-Alpha"
	VerNumber = "2.4.0"
	VerCommit = "dev"
)

var requiredTools = map[string]string{
	"git": ">=2.25",
}

func main() {
	if len(os.Args) < 2 {
		PrintUsage()
		return
	}

	switch os.Args[1] {
	case "create":
		if len(os.Args) < 4 {
			yellow := color.New(color.FgYellow).SprintFunc()

			fmt.Println(yellow("Usage: zoi create <app-template> <app-name>"))
			return
		}
		CreateCommand(os.Args[2], os.Args[3])
	case "make":
		if len(os.Args) < 3 {
			yellow := color.New(color.FgYellow).SprintFunc()

			fmt.Println(yellow("Usage: zoi make <config.yaml>"))
			return
		}
		MakeCommand(os.Args[2])
	case "set":
		if len(os.Args) < 4 {
			yellow := color.New(color.FgYellow).SprintFunc()

			fmt.Println(yellow("Usage: zoi set <key> <value>"))
			return
		}
		SetCommand(os.Args[2], os.Args[3])
	case "install":
		if len(os.Args) < 3 {
			yellow := color.New(color.FgYellow).SprintFunc()

			fmt.Println(yellow("Usage: zoi install <package>[@version]"))
			return
		}
		InstallCommand(os.Args[2])
	case "check":
		CheckCommand()
	// case "update":
	// 	UpdateCommand()
	case "info":
		InfoCommand()
	case "version":
		VersionCommand()
	case "about":
		AboutCommand()
	case "help":
		PrintUsage()
	case "--version", "-v":
		VersionCommand()
	default:
		NotFoundCommand()
	}
}
