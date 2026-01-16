package cli

import (
	"bytes"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

func TestWithCompletionCommand(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.AddCommand(newNextCmd(), newTagCmd())
	root.SetArgs([]string{"--help"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand())
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_completion.golden")
}

func TestWithCompletionCommandExtraShells(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "--help"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand(
		WithExtraShells(ShellPowerShell, ShellNushell),
	))
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bash, zsh, fish, powershell, nushell")
}

func TestWithCompletionCommandCustomShells(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "--help"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand(
		WithShells(ShellBash, ShellPowerShell),
	))
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bash, powershell")
	assert.NotContains(t, output, "zsh")
	assert.NotContains(t, output, "fish")
}

func TestCompletionSubcommandHelp(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "--help"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand())
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "completion_help.golden")
}

func TestCompletionSubcommandHelpExtraShells(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "--help"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand(
		WithExtraShells(ShellPowerShell, ShellNushell),
	))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "completion_help_extra_shells.golden")
}

func TestCompletionGeneratesBashScript(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "bash"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand())
	require.NoError(t, err)

	output := buf.String()
	// Check for bash completion script markers
	assert.True(t, len(output) > 0, "completion script should not be empty")
	assert.Contains(t, output, "bash", "should contain bash-specific content")
}

func TestCompletionGeneratesZshScript(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "zsh"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand())
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, len(output) > 0, "completion script should not be empty")
}

func TestCompletionGeneratesFishScript(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "fish"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand())
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, len(output) > 0, "completion script should not be empty")
}

func TestCompletionRejectsUnconfiguredShell(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "powershell"})

	err := Execute(root, WithStdout(&buf), WithStderr(&buf), WithCompletionCommand())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell")
}

func TestCompletionAcceptsConfiguredShell(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.SetArgs([]string{"completion", "powershell"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand(
		WithExtraShells(ShellPowerShell),
	))
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, len(output) > 0, "completion script should not be empty")
}

func TestDefaultShells(t *testing.T) {
	shells := DefaultShells()
	assert.Equal(t, []Shell{ShellBash, ShellZsh, ShellFish}, shells)
}

func TestCompleterFiles(t *testing.T) {
	completer := Files(".yaml", ".json")
	action := completer.toAction()
	assert.NotNil(t, action)
}

func TestCompleterDirectories(t *testing.T) {
	completer := Directories()
	action := completer.toAction()
	assert.NotNil(t, action)
}

func TestCompleterValues(t *testing.T) {
	completer := Values("one", "two", "three")
	action := completer.toAction()
	assert.NotNil(t, action)
}

func TestCompleterValuesDescribed(t *testing.T) {
	completer := ValuesDescribed("json", "JSON format", "yaml", "YAML format")
	action := completer.toAction()
	assert.NotNil(t, action)
}

func TestCompleterExecutables(t *testing.T) {
	completer := Executables()
	action := completer.toAction()
	assert.NotNil(t, action)
}

func TestCompleterNone(t *testing.T) {
	completer := None()
	action := completer.toAction()
	assert.NotNil(t, action)
}

func TestCompleterActionFunc(t *testing.T) {
	completer := ActionFunc(func() carapace.Action {
		return carapace.ActionValues("custom1", "custom2")
	})
	action := completer.toAction()
	assert.NotNil(t, action)
}

func TestCompleteFlag(t *testing.T) {
	var buf bytes.Buffer

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run:   func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().String("format", "", "output format")
	cmd.SetArgs([]string{"--help"})

	err := Execute(cmd, WithStdout(&buf), WithCompletionCommand(
		CompleteFlag("format", Values("json", "yaml", "text")),
	))
	require.NoError(t, err)
}

func TestCompletePositional(t *testing.T) {
	var buf bytes.Buffer

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run:   func(_ *cobra.Command, _ []string) {},
	}
	cmd.SetArgs([]string{"--help"})

	err := Execute(cmd, WithStdout(&buf), WithCompletionCommand(
		CompletePositional(0, Directories()),
	))
	require.NoError(t, err)
}

func TestCompleteSubcommand(t *testing.T) {
	var buf bytes.Buffer

	root := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}

	sub := &cobra.Command{
		Use:   "sub",
		Short: "Subcommand",
		Run:   func(_ *cobra.Command, _ []string) {},
	}
	sub.Flags().String("env", "", "environment")
	root.AddCommand(sub)
	root.SetArgs([]string{"--help"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand(
		CompleteSubcommand("sub",
			CompleteFlag("env", Values("dev", "staging", "prod")),
		),
	))
	require.NoError(t, err)
}

func TestInferFlagCompletionsForEnum(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	logLevel := Enum(LogInfo, LogDebug, LogInfo, LogWarn, LogError)
	cmd.Flags().Var(logLevel, "log-level", "set the logging level")

	actions := inferFlagCompletions(cmd)
	assert.Contains(t, actions, "log-level")
}

func TestCompletionNotShownOnSubcommandHelp(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.AddCommand(newNextCmd())
	root.SetArgs([]string{"next", "--help"})

	err := Execute(root, WithStdout(&buf), WithCompletionCommand())
	require.NoError(t, err)

	output := buf.String()
	assert.NotContains(t, output, "SHELL COMPLETION")
}
