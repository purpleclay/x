package cli

import (
	"math"
	"strings"
)

func dedent(s string) string {
	lines := strings.Split(s, "\n")

	minIndent := math.MaxInt
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent == math.MaxInt {
		return s
	}

	for i, line := range lines {
		if len(line) >= minIndent {
			lines[i] = line[minIndent:]
		}
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
