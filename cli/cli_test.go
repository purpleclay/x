package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestExecuteWithPositionalArgs(t *testing.T) {
	var capturedArgs []string
	subExecuted := false

	root := &cobra.Command{
		Use:   "myapp [PATHS...]",
		Short: "Example app",
		RunE: func(_ *cobra.Command, args []string) error {
			capturedArgs = args
			return nil
		},
	}

	sub := &cobra.Command{
		Use:   "sub",
		Short: "A subcommand",
		RunE: func(_ *cobra.Command, _ []string) error {
			subExecuted = true
			return nil
		},
	}
	root.AddCommand(sub)
	root.SetArgs([]string{"go.mod"})

	var buf bytes.Buffer
	err := Execute(root, WithStdout(&buf), WithStderr(&buf))

	require.NoError(t, err)
	require.Equal(t, []string{"go.mod"}, capturedArgs)
	require.False(t, subExecuted)

	capturedArgs = nil
	root.SetArgs([]string{"sub"})
	err = Execute(root, WithStdout(&buf), WithStderr(&buf))

	require.NoError(t, err)
	require.True(t, subExecuted)
	require.Nil(t, capturedArgs)
}

func TestExecuteWithArgsValidator(t *testing.T) {
	var capturedArgs []string

	cmd := &cobra.Command{
		Use:   "myapp [PATHS...]",
		Short: "Example app",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			capturedArgs = args
			return nil
		},
	}

	cmd.SetArgs([]string{"go.mod"})

	var buf bytes.Buffer
	err := Execute(cmd, WithStdout(&buf), WithStderr(&buf))

	require.NoError(t, err)
	require.Equal(t, []string{"go.mod"}, capturedArgs)

	cmd.SetArgs([]string{})
	err = Execute(cmd, WithStdout(&buf), WithStderr(&buf))
	require.Error(t, err)
}
