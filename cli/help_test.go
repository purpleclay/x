package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

// Log level enum for NSV commands
type LogLevel string

const (
	LogDebug LogLevel = "debug"
	LogInfo  LogLevel = "info"
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
)

// use a realistic example based on: https://github.com/purpleclay/nsv
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nsv",
		Short: "Manage semantic versioning without any config",
		Long: `
			NSV (Next Semantic Version) is a convention-based semantic versioning
			tool that leans on the power of conventional commits to make versioning
			your software a breeze.

			There is no need to manually maintain a version file or embed the
			version within your source code. NSV will do all of this for you.
		`,
	}

	logLevel := Enum(LogInfo, LogDebug, LogInfo, LogWarn, LogError)
	cmd.PersistentFlags().VarP(logLevel, "log-level", "l", "set the logging verbosity")
	cmd.PersistentFlags().Bool("no-color", false, "disable colored output")
	cmd.PersistentFlags().Bool("no-log", false, "disable all log output")

	return cmd
}

func newNextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next [PATH]...",
		Short: "Generate the next semantic version",
		Long: `
			Generate the next semantic version based on the conventional commit
			history of your repository.
		`,
		Example: `
			# Generate the next semantic version
			nsv next

			# Generate and output only the version number
			nsv next --show

			# Use a custom format
			nsv next --format "v{{.Version}}"
		`,
		Run: func(_ *cobra.Command, _ []string) {},
	}

	cmd.Flags().BoolP("show", "s", false, "show how the version was generated")
	cmd.Flags().StringP("format", "f", "", "provide a go template for changing the default version format")
	cmd.Flags().StringSlice("major-prefixes", nil, "a list of conventional commit prefixes that will trigger a major version increment")
	cmd.Flags().StringSlice("minor-prefixes", nil, "a list of conventional commit prefixes that will trigger a minor version increment")
	cmd.Flags().StringSlice("patch-prefixes", nil, "a list of conventional commit prefixes that will trigger a patch version increment")

	return cmd
}

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag [PATH]...",
		Short: "Tag the repository with the next semantic version based on the commit history",
		Run:   func(_ *cobra.Command, _ []string) {},
	}

	cmd.Flags().StringP("message", "m", "", "a custom message for the tag")

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build time version information",
		Run:   func(_ *cobra.Command, _ []string) {},
	}
}

func TestHelp(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.AddCommand(newNextCmd(), newTagCmd(), newVersionCmd())
	root.SetArgs([]string{"--help"})

	err := Execute(root, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help.golden")
}

func TestHelpWithExamples(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.AddCommand(newNextCmd())
	root.SetArgs([]string{"next", "--help"})

	err := Execute(root, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_examples.golden")
}

func TestHelpWithGlobalFlags(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.AddCommand(newTagCmd())
	root.SetArgs([]string{"tag", "--help"})

	err := Execute(root, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_global_flags.golden")
}

func TestHelpWithSubcommands(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.AddCommand(newNextCmd(), newTagCmd(), newVersionCmd())
	root.SetArgs([]string{"--help"})

	err := Execute(root, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_subcommands.golden")
}

func TestHelpWithNoWrapping(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	root.AddCommand(newNextCmd(), newTagCmd(), newVersionCmd())
	root.SetArgs([]string{"--help"})

	err := Execute(root, WithStdout(&buf), WithWidth(0))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_no_wrapping.golden")
}

func TestHelpWithFlagGroups(t *testing.T) {
	var buf bytes.Buffer

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application to the cloud",
		Run:   func(_ *cobra.Command, _ []string) {},
	}

	cmd.Flags().BoolP("verbose", "v", false, "enable verbose logging")
	cmd.Flags().String("config", "", "path to config file")
	cmd.Flags().String("format", "json", "output format")
	cmd.Flags().String("output", "", "output file path")
	cmd.Flags().String("token", "", "authentication token")
	cmd.Flags().String("url", "", "API endpoint URL")

	FlagGroup(cmd, "Output Options", "format", "output")
	FlagGroup(cmd, "Authentication", "token", "url")

	cmd.SetArgs([]string{"--help"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_flag_groups.golden")
}

func TestHelpWithEnum(t *testing.T) {
	var buf bytes.Buffer

	type TrustLevel int
	const (
		TrustUnknown TrustLevel = iota + 1
		TrustNever
		TrustMarginal
		TrustFull
		TrustUltimate
	)

	trust := Enum(TrustUnknown, TrustUnknown, TrustNever, TrustMarginal, TrustFull, TrustUltimate).
		WithHelp(
			"I don't know or won't say",
			"I do NOT trust",
			"I trust marginally",
			"I trust fully",
			"I trust ultimately",
		)

	cmd := &cobra.Command{
		Use:   "gpg-import",
		Short: "Import your GPG private key into the local keyring",
		Long: `
			Import your GPG private key into the local keyring of your CI
			environment. Supports automatic detection and deletion of the
			imported key after use.
		`,
		Run: func(_ *cobra.Command, _ []string) {},
	}

	cmd.Flags().VarP(trust, "trust-level", "t", "a level of trust to associate with the GPG private key")
	cmd.SetArgs([]string{"--help"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_enum.golden")
}

func TestHelpWithEnvVars(t *testing.T) {
	var buf bytes.Buffer

	root := newRootCmd()
	tag := newTagCmd()
	BindEnv(tag.Flags().Lookup("message"), "NSV_TAG_MESSAGE")
	root.AddCommand(tag)
	root.SetArgs([]string{"tag", "--help"})

	err := Execute(root, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_env_vars.golden")
}

func TestHelpWithEnvVarsSet(t *testing.T) {
	t.Setenv("NSV_TAG_MESSAGE", "chore: this is a release created by nsv")

	var buf bytes.Buffer

	root := newRootCmd()
	tag := newTagCmd()
	BindEnv(tag.Flags().Lookup("message"), "NSV_TAG_MESSAGE")
	root.AddCommand(tag)
	root.SetArgs([]string{"tag", "--help"})

	err := Execute(root, WithStdout(&buf))
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_env_vars_set.golden")
}

func TestTokenizeExample(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []exampleToken
	}{
		{
			name:  "WithSimpleCommand",
			input: "cmd arg",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "arg", tokenType: tokenWord},
			},
		},
		{
			name:  "WithPipeOperator",
			input: "cmd | grep foo",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "|", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "grep", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "foo", tokenType: tokenWord},
			},
		},
		{
			name:  "WithRedirectOperator",
			input: "cmd > out.txt",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: ">", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "out.txt", tokenType: tokenWord},
			},
		},
		{
			name:  "WithAppendRedirect",
			input: "cmd >> log.txt",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: ">>", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "log.txt", tokenType: tokenWord},
			},
		},
		{
			name:  "WithDoubleQuotedString",
			input: `cmd "hello world"`,
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: `"hello world"`, tokenType: tokenString},
			},
		},
		{
			name:  "WithSingleQuotedString",
			input: "cmd 'hello world'",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "'hello world'", tokenType: tokenString},
			},
		},
		{
			name:  "WithEnvironmentVariable",
			input: "VAR=value cmd",
			expected: []exampleToken{
				{value: "VAR=value", tokenType: tokenEnvAssign},
				{value: " ", tokenType: tokenWhitespace},
				{value: "cmd", tokenType: tokenWord},
			},
		},
		{
			name:  "WithMultipleEnvironmentVariables",
			input: "VAR1=val1 VAR2=val2 cmd",
			expected: []exampleToken{
				{value: "VAR1=val1", tokenType: tokenEnvAssign},
				{value: " ", tokenType: tokenWhitespace},
				{value: "VAR2=val2", tokenType: tokenEnvAssign},
				{value: " ", tokenType: tokenWhitespace},
				{value: "cmd", tokenType: tokenWord},
			},
		},
		{
			name:  "WithLineContinuation",
			input: `cmd --flag \`,
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "--flag", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: `\`, tokenType: tokenLineContinuation},
			},
		},
		{
			name:  "WithCompoundOperators",
			input: "cmd1 && cmd2 || cmd3",
			expected: []exampleToken{
				{value: "cmd1", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "&&", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "cmd2", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "||", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "cmd3", tokenType: tokenWord},
			},
		},
		{
			name:  "WithFlagEquals",
			input: "cmd --format=json",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "--format=json", tokenType: tokenWord},
			},
		},
		{
			name:  "WithComplexPipeline",
			input: `API_KEY=secret cmd --verbose | grep "pattern" > results.txt`,
			expected: []exampleToken{
				{value: "API_KEY=secret", tokenType: tokenEnvAssign},
				{value: " ", tokenType: tokenWhitespace},
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "--verbose", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "|", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "grep", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: `"pattern"`, tokenType: tokenString},
				{value: " ", tokenType: tokenWhitespace},
				{value: ">", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "results.txt", tokenType: tokenWord},
			},
		},
		{
			name:  "WithSemicolonSeparator",
			input: "cmd1; cmd2",
			expected: []exampleToken{
				{value: "cmd1", tokenType: tokenWord},
				{value: ";", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "cmd2", tokenType: tokenWord},
			},
		},
		{
			name:  "WithBackgroundOperator",
			input: "cmd &",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "&", tokenType: tokenOperator},
			},
		},
		{
			name:  "WithHeredocOperator",
			input: "cmd << EOF",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "<<", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "EOF", tokenType: tokenWord},
			},
		},
		{
			name:  "WithInputRedirect",
			input: "cmd < input.txt",
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: "<", tokenType: tokenOperator},
				{value: " ", tokenType: tokenWhitespace},
				{value: "input.txt", tokenType: tokenWord},
			},
		},
		{
			name:  "WithEscapedQuoteInString",
			input: `cmd "hello \"world\""`,
			expected: []exampleToken{
				{value: "cmd", tokenType: tokenWord},
				{value: " ", tokenType: tokenWhitespace},
				{value: `"hello \"world\""`, tokenType: tokenString},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenizeExample(tt.input)
			require.Equal(t, tt.expected, got)
		})
	}
}
