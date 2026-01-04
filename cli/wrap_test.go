package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnfill(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "SingleLine",
			in:   "hello world",
			want: "hello world",
		},
		{
			name: "MultipleLines",
			in:   "hello\nworld",
			want: "hello world",
		},
		{
			name: "MultipleParagraphs",
			in:   "first paragraph\nwith two lines\n\nsecond paragraph",
			want: "first paragraph with two lines\n\nsecond paragraph",
		},
		{
			name: "CRLFLineEndings",
			in:   "hello\r\nworld",
			want: "hello world",
		},
		{
			name: "MixedLineEndings",
			in:   "hello\r\nworld\nfoo",
			want: "hello world foo",
		},
		{
			name: "TrailingNewline",
			in:   "hello\nworld\n",
			want: "hello world\n",
		},
		{
			name: "MultipleSpacesCollapsed",
			in:   "hello   world",
			want: "hello world",
		},
		{
			name: "LeadingAndTrailingSpaces",
			in:   "  hello world  ",
			want: "hello world",
		},
		{
			name: "EmptyString",
			in:   "",
			want: "",
		},
		{
			name: "OnlyNewlines",
			in:   "\n\n",
			want: "",
		},
		{
			name: "ParagraphWithTrailingNewline",
			in:   "first\n\nsecond\n",
			want: "first\n\nsecond\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unfill(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
