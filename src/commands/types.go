package commands

import (
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Branch string
	Status string
	Number string
}

var statusOrder = map[string]int{
	"Pre-Alpha":    0,
	"Alpha":        1,
	"Pre-Beta":     2,
	"Beta":         3,
	"Pre-Release":  4,
	"Early-Access": 5,
	"Demo":         6,
	"Release":      7,
}

func parseVersion(branch, status, number string) Version {
	return Version{
		Branch: branch,
		Status: status,
		Number: number,
	}
}

func (v Version) Compare(other Version) int {
	currentRank, okCurrent := statusOrder[v.Status]
	otherRank, okOther := statusOrder[other.Status]
	if !okCurrent {
		currentRank = -1
	}
	if !okOther {
		otherRank = -1
	}

	if otherRank > currentRank {
		return 1
	}
	if otherRank < currentRank {
		return -1
	}

	return compareSemver(v.Number, other.Number)
}

func compareSemver(a, b string) int {
	cleanRegex := regexp.MustCompile(`^(\d+\.\d+(?:\.\d+)?)?.*`)
	aClean := cleanRegex.ReplaceAllString(a, "$1")
	bClean := cleanRegex.ReplaceAllString(b, "$1")

	aParts := strings.Split(aClean, ".")
	bParts := strings.Split(bClean, ".")

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			numStr := regexp.MustCompile(`^(\d+).*`).ReplaceAllString(aParts[i], "$1")
			aNum, _ = strconv.Atoi(numStr)
		}
		if i < len(bParts) {
			numStr := regexp.MustCompile(`^(\d+).*`).ReplaceAllString(bParts[i], "$1")
			bNum, _ = strconv.Atoi(numStr)
		}

		if bNum > aNum {
			return 1
		}
		if bNum < aNum {
			return -1
		}
	}

	return 0
}

type VersionInfo struct {
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
