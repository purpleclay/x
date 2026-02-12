package cli

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

func TestDefaultErrorHandler(t *testing.T) {
	var buf bytes.Buffer
	theme := DefaultTheme()

	DefaultErrorHandler(&buf, theme, errors.New("something went wrong"))

	golden.Assert(t, buf.String(), "error.golden")
}

func TestDefaultErrorHandlerUsageError(t *testing.T) {
	var buf bytes.Buffer
	theme := DefaultTheme()

	DefaultErrorHandler(&buf, theme, errors.New("unknown flag: --verbose"))

	golden.Assert(t, buf.String(), "error_usage.golden")
}

func TestExecuteRendersStyledError(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "myapp",
		Short: "A test app",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("something went wrong")
		},
	}

	var stdout, stderr bytes.Buffer
	err := Execute(cmd, WithStdout(&stdout), WithStderr(&stderr))

	require.Error(t, err)
	golden.Assert(t, stderr.String(), "error.golden")
}

func TestExecuteWithCustomErrorHandler(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "myapp",
		Short: "A test app",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("custom error")
		},
	}

	var stdout, stderr bytes.Buffer
	handler := func(w io.Writer, _ Theme, err error) {
		_, _ = io.WriteString(w, "CUSTOM: "+err.Error()+"\n")
	}

	err := Execute(cmd, WithStdout(&stdout), WithStderr(&stderr), WithErrorHandler(handler))

	require.Error(t, err)
	require.Equal(t, "CUSTOM: custom error\n", stderr.String())
}

func TestIsUsageError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "UnknownFlag",
			err:      errors.New("unknown flag: --verbose"),
			expected: true,
		},
		{
			name:     "UnknownShorthandFlag",
			err:      errors.New("unknown shorthand flag: 'v' in -v"),
			expected: true,
		},
		{
			name:     "FlagNeedsArgument",
			err:      errors.New("flag needs an argument: --config"),
			expected: true,
		},
		{
			name:     "UnknownCommand",
			err:      errors.New("unknown command \"foo\" for \"myapp\""),
			expected: true,
		},
		{
			name:     "InvalidArgument",
			err:      errors.New("invalid argument \"foo\" for \"--count\""),
			expected: true,
		},
		{
			name:     "AcceptsArgs",
			err:      errors.New("accepts between 1 and 3 arg(s), received 0"),
			expected: true,
		},
		{
			name:     "RequiredFlag",
			err:      errors.New("required flag \"config\" not set"),
			expected: true,
		},
		{
			name:     "GenericError",
			err:      errors.New("something went wrong"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, isUsageError(tt.err))
		})
	}
}
