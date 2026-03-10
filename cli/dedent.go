package cli

import (
	"math"
	"strings"
)

func dedent(s string) string {
	s = expandTabs(s, 4)
	lines := strings.Split(s, "\n")

	minIndent := math.MaxInt
	for i, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))
		if i == 0 && indent == 0 {
			continue
		}

		if indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent == math.MaxInt {
		return s
	}

	for i, line := range lines {
		indent := len(line) - len(strings.TrimLeft(line, " "))
		strip := min(indent, minIndent)
		lines[i] = line[strip:]
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func expandTabs(s string, tabWidth int) string {
	if !strings.Contains(s, "\t") {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))
	col := 0
	for _, ch := range s {
		if ch == '\t' {
			spaces := tabWidth - (col % tabWidth)
			for range spaces {
				b.WriteByte(' ')
			}
			col += spaces
		} else {
			b.WriteRune(ch)
			if ch == '\n' {
				col = 0
			} else {
				col++
			}
		}
	}
	return b.String()
}
