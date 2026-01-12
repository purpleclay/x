package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

func newVersionTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "myapp",
		Short: "A test application",
		Run:   func(_ *cobra.Command, _ []string) {},
	}
}

func testVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   "1.2.3",
		GitCommit: "abc1234",
		GitBranch: "main",
		BuildDate: "2024-01-15T10:30:00Z",
		GoVersion: "go1.21.0",
	}
}

func TestVersionFlag(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"--version"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionFlag(testVersionInfo()),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "version.golden")
}

func TestVersionFlagShort(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"-V"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionFlag(testVersionInfo()),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "version.golden")
}

func TestVersionCommand(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"version"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionCommand(testVersionInfo()),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "version.golden")
}

func TestVersionCommandShort(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"version", "--short"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionCommand(testVersionInfo()),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "version_short.golden")
}

func TestVersionCommandJSON(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"version", "--json"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionCommand(testVersionInfo()),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "version_json.golden")
}

func TestVersionMinimal(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"--version"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionFlag(VersionInfo{Version: "0.1.0"}),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "version_minimal.golden")
}

func TestHelpWithVersionFlag(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"--help"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionFlag(testVersionInfo()),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_version_flag.golden")
}

func TestHelpWithVersionCommand(t *testing.T) {
	var buf bytes.Buffer

	cmd := newVersionTestCmd()
	cmd.SetArgs([]string{"--help"})

	err := Execute(cmd,
		WithStdout(&buf),
		WithVersionCommand(testVersionInfo()),
	)
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_version_command.golden")
}
