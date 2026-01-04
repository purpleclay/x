package cli

import (
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}
	return wordwrap.String(unfill(s), width)
}

// Inspiration taken from https://github.com/mgeisler/textwrap/blob/master/src/refill.rs
func unfill(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")

	hasTrailingNewline := strings.HasSuffix(s, "\n")
	s = strings.TrimSuffix(s, "\n")

	paragraphs := strings.Split(s, "\n\n")
	for i, para := range paragraphs {
		para = strings.ReplaceAll(para, "\n", " ")

		for strings.Contains(para, "  ") {
			para = strings.ReplaceAll(para, "  ", " ")
		}

		paragraphs[i] = strings.TrimSpace(para)
	}

	s = strings.Join(paragraphs, "\n\n")
	if hasTrailingNewline && s != "" {
		s += "\n"
	}
	return s
}
