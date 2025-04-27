package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func getOS() string {
	switch os := runtime.GOOS; os {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return os
	}
}

func getLinuxDistro() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}

	re := regexp.MustCompile(`ID=("?)(.+?)\1`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) > 2 {
		return matches[2]
	}
	return "unknown"
}

func getArch() string {
	switch arch := runtime.GOARCH; arch {
	case "amd64":
		return "x86_64"
	case "386":
		return "x86"
	case "arm64":
		return "ARM64"
	default:
		return arch
	}
}

func executeCommand(command string) error {
	cmd := exec.Command("powershell", "-c", command)
	if runtime.GOOS != "windows" {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateCommand(app string, appName string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	resp, err := http.Get(getAppsURL())
	if err != nil {
		fmt.Printf("%s Failed to fetch app definitions: %v\n", red("✗"), err)
		return
	}
	defer resp.Body.Close()

	var apps map[string]AppDefinition
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		fmt.Printf("%s Invalid app definitions: %v\n", red("✗"), err)
		return
	}

	appDef, exists := apps[app]
	if !exists {
		fmt.Printf("%s Application template '%s' not found\n", red("✗"), app)
		return
	}

	for depName, dep := range appDef.Dependencies {
		fmt.Printf("%s Checking %s...\n", cyan("ℹ"), depName)
		if err := executeCommand(dep.Check); err == nil {
			fmt.Printf("%s %s already installed\n", green("✓"), depName)
			continue
		}

		var installCmd string
		switch runtime.GOOS {
		case "linux":
			distro := getLinuxDistro()
			installCmd = dep.Install.Linux[distro]
			if installCmd == "" {
				installCmd = dep.Install.Linux["default"]
			}
		case "darwin":
			installCmd = dep.Install.Darwin
		case "windows":
			installCmd = dep.Install.Win32
		}

		if installCmd == "" {
			installCmd = dep.Install.Default
		}

		fmt.Printf("%s Installing %s...\n", cyan("ℹ"), depName)
		if err := executeCommand(installCmd); err != nil {
			fmt.Printf("%s Failed to install %s: %v\n", red("✗"), depName, err)
			return
		}

		fmt.Printf("%s Successfully installed %s\n", green("✓"), depName)
	}

	fmt.Printf("%s Creating application...\n", cyan("ℹ"))
	createCmd := strings.ReplaceAll(appDef.Create, "${appName}", appName)
	if err := executeCommand(createCmd); err != nil {
		fmt.Printf("%s Failed to create application: %v\n", red("✗"), err)
		return
	}
	fmt.Printf("%s Successfully created %s application\n", green("✓"), appName)
}

func MakeCommand(yamlFile string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("%s Failed to read YAML file: %v\n", red("✗"), err)
		return
	}

	var config YamlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("%s Invalid YAML format: %v\n", red("✗"), err)
		return
	}

	for _, pkg := range config.Packages {
		fmt.Printf("%s Checking %s...\n", cyan("ℹ"), pkg.Name)

		if err := executeCommand(pkg.Check); err == nil {
			fmt.Printf("%s %s already installed\n", green("✓"), pkg.Name)
			continue
		}

		var installCmd string
		switch runtime.GOOS {
		case "linux":
			installCmd = pkg.Install.Linux
		case "windows":
			installCmd = pkg.Install.Windows
		case "darwin":
			installCmd = pkg.Install.Darwin
		}

		if installCmd == "" {
			installCmd = pkg.Install.Default
		}

		if installCmd == "" {
			fmt.Printf("%s No install command for %s\n", red("✗"), pkg.Name)
			return
		}

		fmt.Printf("%s Installing %s...\n", cyan("ℹ"), pkg.Name)
		if err := executeCommand(installCmd); err != nil {
			fmt.Printf("%s Failed to install %s: %v\n", red("✗"), pkg.Name, err)
			return
		}
		fmt.Printf("%s Successfully installed %s\n", green("✓"), pkg.Name)
	}

	fmt.Printf("%s Creating application...\n", cyan("ℹ"))
	createCmd := strings.ReplaceAll(config.CreateCommand, "${appName}", config.AppName)
	if err := executeCommand(createCmd); err != nil {
		fmt.Printf("%s Failed to create application: %v\n", red("✗"), err)
		return
	}
	fmt.Printf("%s Successfully created '%s'\n", green("✓"), config.AppName)
}

func CheckCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(cyan("🏥 Running Zoi Environment Health Check\n"))

	fmt.Println(cyan("🔧 Essential Tools:"))
	for tool, constraint := range requiredTools {
		checkTool(tool, constraint)
	}

	fmt.Println("\n" + cyan("🌐 Network Connectivity:"))
	checkNetwork()
}

func checkTool(tool, constraint string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	path, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf("  %s %s: %s\n", red("✗"), tool, "Not installed")
		return
	}

	cmd := exec.Command(tool, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("  %s %s: %s\n", red("✗"), tool, "Version check failed")
		return
	}

	version := strings.Split(string(out), "\n")[0]
	fmt.Printf("  %s %s: %s (%s)\n",
		green("✓"),
		tool,
		strings.TrimSpace(version),
		cyan(path),
	)

}

func checkNetwork() {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	_, err := net.LookupHost("google.com")
	if err == nil {
		fmt.Printf("  %s DNS resolution working\n", green("✓"))
	} else {
		fmt.Printf("  %s DNS resolution failed\n", red("✗"))
	}

	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get("https://example.com")
	if err == nil && resp.StatusCode == 200 {
		fmt.Printf("  %s HTTPS connectivity working\n", green("✓"))
	} else {
		fmt.Printf("  %s HTTPS connection failed\n", red("✗"))
	}

	fmt.Print("  Testing latency... ")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "ping", "-n", "1", "8.8.8.8")
	} else {
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "8.8.8.8")
	}

	if err := cmd.Run(); err == nil {
		fmt.Printf("%s\n", green("✓ Network responsive"))
	} else {
		fmt.Printf("%s\n", red("✗ High latency/no connection"))
	}
}

func SetCommand(key, value string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	if key != "appsUrl" {
		fmt.Printf("%s Unknown config key: %s\n", red("✗"), key)
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("%s Failed to load config: %v\n", red("✗"), err)
		return
	}

	if value == "default" {
		cfg.AppsURL = defaultAppsURL
	} else {
		cfg.AppsURL = value
	}

	if err := saveConfig(cfg); err != nil {
		fmt.Printf("%s Failed to save config: %v\n", red("✗"), err)
		return
	}

	fmt.Printf("%s Updated apps URL to: %s\n",
		green("✓"),
		cyan(cfg.AppsURL),
	)
}

func getAppsURL() string {
	cfg, err := loadConfig()
	if err != nil {
		return defaultAppsURL
	}
	return cfg.AppsURL
}

func checkToolDirectly(pkg string) (string, error) {
	cmd := fmt.Sprintf("%s --version", pkg)
	output, err := runCommandOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("command failed")
	}
	return extractVersion(output), nil
}

func runCommandOutput(command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func extractVersion(output string) string {
	re := regexp.MustCompile(`(\d+\.\d+\.\d+|\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func versionMatches(want, actual string) bool {
	wantParts := strings.Split(want, ".")
	actualParts := strings.Split(actual, ".")

	for i := 0; i < len(wantParts); i++ {
		if i >= len(actualParts) {
			return false
		}
		if wantParts[i] != actualParts[i] {
			return false
		}
	}
	return true
}

func InstallCommand(pkg string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	parts := strings.Split(pkg, "@")
	pkgName := parts[0]
	pkgPath, err := exec.LookPath(pkgName)
	wantVersion := ""
	if len(parts) > 1 {
		wantVersion = parts[1]
	}

	directVersion, directErr := checkToolDirectly(pkgName)
	if directErr == nil {
		handleExistingInstallation(pkgName, directVersion, wantVersion, pkgPath)
		return
	}

	pm, err := detectPackageManager()
	if err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		return
	}

	installed, pkgVersion := isInstalled(pm, pkgName)
	if installed {
		handleExistingInstallation(pkgName, pkgVersion, wantVersion, pm.Name)
		return
	}

	installPkg := pkgName
	if wantVersion != "" {
		installPkg = fmt.Sprintf("%s=%s", pkgName, wantVersion)
	}
	installCmd := fmt.Sprintf(pm.InstallCmd, installPkg)

	fmt.Printf("%s About to install: %s\n", cyan("ℹ"), cyan(installCmd))
	if !confirmPrompt("Proceed with installation?") {
		return
	}

	fmt.Printf("%s Installing %s...\n", cyan("ℹ"), cyan(pkgName))
	if err := executeCommand(installCmd); err != nil {
		fmt.Printf("%s Installation failed: %v\n", red("✗"), err)
		return
	}

	fmt.Printf("%s Successfully installed %s\n", green("✓"), cyan(pkgName))
}

func handleExistingInstallation(pkg, currentVersion, wantVersion, source string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	versionInfo := ""
	if currentVersion != "" {
		versionInfo = fmt.Sprintf("@%s", currentVersion)
	}

	if wantVersion == "" || versionMatches(wantVersion, currentVersion) {
		fmt.Printf("%s %s%s is already installed (%s)\n",
			yellow("!"),
			cyan(pkg),
			green(versionInfo),
			cyan(source),
		)
		return
	}

	fmt.Printf("%s %s%s is already installed (%s)\n",
		yellow("!"),
		cyan(pkg),
		yellow(versionInfo),
		cyan(source),
	)
	if !confirmPrompt(fmt.Sprintf("Install %s@%s anyway?", pkg, wantVersion)) {
		return
	}
}

func detectPackageManager() (PackageManager, error) {
	os := runtime.GOOS
	distro := getLinuxDistro()

	switch {
	case os == "windows":
		return packageManagers["scoop"], nil

	case os == "darwin":
		return packageManagers["brew"], nil

	case os == "linux":
		switch distro {
		case "arch", "cachyos", "manjaro":
			return packageManagers["pacman"], nil
		case "fedora", "rhel", "centos":
			if commandExists("dnf") {
				return packageManagers["dnf"], nil
			}
			return packageManagers["yum"], nil
		case "alpine":
			return packageManagers["apk"], nil
		case "debian", "ubuntu", "kali", "raspbian":
			return packageManagers["apt"], nil
		default:
			return PackageManager{}, fmt.Errorf("unsupported Linux distribution: %s", distro)
		}

	default:
		return PackageManager{}, fmt.Errorf("unsupported operating system: %s", os)
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func isInstalled(pm PackageManager, pkg string) (bool, string) {
	checkCmd := fmt.Sprintf(pm.CheckCommand, pkg)
	out, err := runCommandOutput(checkCmd)
	if err != nil {
		return false, ""
	}

	if pm.Name == "scoop" {
		if strings.Contains(string(out), "Version:") {
			matches := pm.VersionRegex.FindStringSubmatch(string(out))
			if len(matches) > 1 {
				return true, strings.TrimSpace(matches[1])
			}
			return true, "unknown"
		}
		return false, ""
	}

	matches := pm.VersionRegex.FindStringSubmatch(string(out))
	if len(matches) > 1 {
		return true, strings.TrimSpace(matches[1])
	}
	return true, "unknown"
}

func confirmPrompt(prompt string) bool {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%s %s [y/N] ", yellow("?"), prompt)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

// TODO: Fix update command and install.ps1 script
func UpdateCommand() {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	buildChannel := "production"
	if strings.Contains(VerBranch, "Dev") {
		buildChannel = "development"
	}

	resp, err := http.Get("https://codeberg.org/Zusty/Zoi/raw/branch/app/version.json")
	if err != nil {
		fmt.Printf("%s Failed to check updates: %v\n", red("✗"), err)
		return
	}
	defer resp.Body.Close()

	var versionInfo struct {
		Latest struct {
			Production struct {
				Version string `json:"version"`
				Status  string `json:"status"`
			} `json:"production"`
			Development struct {
				Version string `json:"version"`
				Status  string `json:"status"`
			} `json:"development"`
		} `json:"latest"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
		fmt.Printf("%s Invalid version data: %v\n", red("✗"), err)
		return
	}

	current := parseVersion(
		VerBranch,
		VerStatus,
		VerNumber,
	)

	var latest Version
	if buildChannel == "production" {
		latest = parseVersion(
			"production",
			versionInfo.Latest.Production.Status,
			versionInfo.Latest.Production.Version,
		)
	} else {
		latest = parseVersion(
			"development",
			versionInfo.Latest.Development.Status,
			versionInfo.Latest.Development.Version,
		)
	}

	statusChange := ""
	if current.Status != latest.Status {
		statusChange = fmt.Sprintf("%s → %s ",
			yellow(current.Status),
			green(latest.Status),
		)
	}

	if current == latest {
		fmt.Printf("%s You're already on the latest %s version (%s)\n",
			green("✓"),
			cyan(buildChannel),
			green(current),
		)
		return
	}

	fmt.Printf("%s New %s version available: %s%s%s → %s\n",
		yellow("!"),
		cyan(current.Branch),
		statusChange,
		cyan(current.Number),
		yellow("→"),
		green(latest.Number),
	)

	if !confirmPrompt("Would you like to update?") {
		return
	}

	fmt.Printf("%s Starting update process...\n", cyan("ℹ"))
	if runtime.GOOS == "windows" {
		err = executeCommand("irm zusty.codeberg.page/Zoi/@app/install.ps1|iex")
	} else {
		err = executeCommand("curl -fsSL https://zusty.codeberg.page/Zoi/@app/install.sh | bash")
	}
	if err != nil {
		fmt.Printf("%s Update failed: %v\n", red("✗"), err)
		return
	}

	fmt.Printf("%s Successfully updated to %s\n", green("✓"), green(latest))
}

func InfoCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Println(cyan("Detected environment:"))
	fmt.Printf("- OS:    %s\n", green(getOS()))
	if getLinuxDistro() != "unknown" {
		fmt.Printf("- Linux Distro:  %s\n", green(getLinuxDistro()))
	}
	fmt.Printf("- Arch:  %s\n", green(getArch()))
}

func VersionCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s %s\n",
		cyan("Zoi"),
		green(VerBranch, " ", VerStatus, " ", VerNumber),
	)
	if getLinuxDistro() == "unknown" {
		fmt.Printf("Runtime: %s/%s\n",
			green(runtime.GOOS),
			green(runtime.GOARCH),
		)
	} else {
		fmt.Printf("Runtime: %s/%s/%s\n",
			green(runtime.GOOS),
			green(getLinuxDistro()),
			green(runtime.GOARCH),
		)
	}
	fmt.Printf("Commit: %s\n",
		yellow(VerCommit),
	)
}

func AboutCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(cyan("About Zoi:"))
	fmt.Printf("A universal environment setup tool for developers")
	fmt.Printf("\nCreated by Zillowe Foundation > Zusty")
	fmt.Printf("\nHosted on Codeberg.org/Zusty/Zoi")
}

func NotFoundCommand() {
	red := color.New(color.FgRed).SprintFunc()

	fmt.Println(red("Command not found"))
	fmt.Println(red("Run 'zoi help' for usage."))
}

func PrintUsage() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s\n", yellow("Usage: zoi <command> --<flag>"))
	fmt.Printf("%s\n", cyan("Commands:"))
	fmt.Println("  zoi info                             - Show system information")
	fmt.Println("  zoi version                          - Show version details")
	fmt.Println("  zoi help                             - Show usage")
	fmt.Println("  zoi about                            - Show about information")
	fmt.Println("  zoi check                            - Verify system health and requirements")
	// fmt.Println("  zoi update                           - Check for and install updates")
	fmt.Println("  zoi create <app> <name>              - Create new application")
	fmt.Println("  zoi make <file.yaml>                 - Create new application from a local file")
	fmt.Println("  zoi install <package>[@version]      - Install system packages")
	fmt.Println("  zoi set <key> <value>                - Update configuration")
	fmt.Println("                Available keys: appsUrl")
	fmt.Printf("\n")
	fmt.Printf("%s\n", cyan("Flags:"))
	fmt.Println("  zoi --version -v                     - Show version details")
}
