package commands

import (
	"fmt"

	"github.com/fatih/color"
)

func AboutCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(cyan("About GCT:"))
	fmt.Printf("  GCT is a git commit tool.\n")
	fmt.Printf("\n")
	fmt.Printf("  GCT is a part of the Zillowe Development Suite (ZDS)\n")
	fmt.Printf("\n")
	fmt.Printf("  Created by Zillowe Foundation > Zusty\n")
	fmt.Printf("  Hosted on %s\n", yellow("https://Codeberg.org/Zusty/GCT"))
}
