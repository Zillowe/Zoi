package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

func getLinuxDistro() string {
	files := []string{"/etc/os-release", "/etc/lsb-release"}
	idRegex := regexp.MustCompile(`(?m)^ID=(?:["']?)(.+?)(?:["']?)$`)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err == nil {
			matches := idRegex.FindStringSubmatch(string(data))
			if len(matches) > 1 {
				return strings.ToLower(matches[1])
			}
		}
	}
	if _, err := os.Stat("/etc/arch-release"); err == nil {
		return "arch"
	}
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "fedora-rhel-centos"
	}
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}
	if _, err := os.Stat("/etc/alpine-release"); err == nil {
		return "alpine"
	}
	return "unknown"
}

func executeCommand(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("pwsh", "-Command", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func confirmPrompt(prompt string) bool {
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s %s [y/N] ", yellow("?"), cyan(prompt))
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func NotFoundCommand() {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	attemptedCmd := ""
	if len(os.Args) > 1 {
		attemptedCmd = os.Args[1]
	}

	fmt.Printf("%s Unknown command: '%s'\n", red("Error:"), yellow(attemptedCmd))
	fmt.Printf("%s Run 'gct help' for a list of available commands.\n", yellow("Hint:"))
}

func promptForAction(prompt string) rune {
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s %s ", yellow("?"), cyan(prompt))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if len(input) > 0 {
		switch input[0] {
		case 'c':
			return 'c'
		case 'e':
			return 'e'
		case 'q':
			return 'q'
		}
	}
	return '\n'
}

func promptForInput(prompt string) string {
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s %s: ", cyan("?"), prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
