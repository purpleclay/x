package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "WithBasicIndentation",
			input: `
                line one
                line two
            `,
			expected: "line one\nline two",
		},
		{
			name: "WithMixedIndentation",
			input: `
                line one
                    line two
                line three
            `,
			expected: "line one\n    line two\nline three",
		},
		{
			name: "WithEmptyLines",
			input: `
                line one

                line two
            `,
			expected: "line one\n\nline two",
		},
		{
			name:     "WithTabs",
			input:    "\t\tline one\n\t\tline two",
			expected: "line one\nline two",
		},
		{
			name:     "WithNoIndentation",
			input:    "line one\nline two",
			expected: "line one\nline two",
		},
		{
			name:     "WithLeadingSpaces",
			input:    "    line one",
			expected: "line one",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dedent(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
