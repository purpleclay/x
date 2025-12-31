package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
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

	cmd.PersistentFlags().String("log-level", "info", "set the logging verbosity")
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
	cmd.Flags().StringSlice("major-prefixes", nil, "a list of conventional commit prefixes for a major version increment")
	cmd.Flags().StringSlice("minor-prefixes", nil, "a list of conventional commit prefixes for a minor version increment")
	cmd.Flags().StringSlice("patch-prefixes", nil, "a list of conventional commit prefixes for a patch version increment")

	return cmd
}

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag [PATH]...",
		Short: "Tag the repository with the next semantic version",
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
	root := newRootCmd()
	root.AddCommand(newNextCmd(), newTagCmd(), newVersionCmd())
	Configure(root)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--help"})
	err := root.Execute()
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help.golden")
}

func TestHelpWithExamples(t *testing.T) {
	root := newRootCmd()
	next := newNextCmd()
	root.AddCommand(next)
	Configure(root)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"next", "--help"})
	err := root.Execute()
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_examples.golden")
}

func TestHelpWithGlobalFlags(t *testing.T) {
	root := newRootCmd()
	tag := newTagCmd()
	root.AddCommand(tag)
	Configure(root)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"tag", "--help"})
	err := root.Execute()
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_global_flags.golden")
}

func TestHelpWithSubcommands(t *testing.T) {
	root := newRootCmd()
	root.AddCommand(newNextCmd(), newTagCmd(), newVersionCmd())
	Configure(root)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--help"})
	err := root.Execute()
	require.NoError(t, err)

	golden.Assert(t, buf.String(), "help_with_subcommands.golden")
}
