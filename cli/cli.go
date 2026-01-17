package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	mango "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

// Option is a functional option for configuring the CLI.
type Option func(*options)

type options struct {
	ctx            context.Context
	completion     *completionOptions
	manpages       bool
	stdout         io.Writer
	stderr         io.Writer
	theme          Theme
	version        *VersionInfo
	versionCommand bool
	width          int
}

func defaultOptions() *options {
	return &options{
		ctx:      context.Background(),
		manpages: true,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		theme:    DefaultTheme(),
		width:    80,
	}
}

// WithStdout sets the standard output writer for the CLI.
//
//	var buf strings.Builder
//	cli.Execute(root, cli.WithStdout(&buf))
//	fmt.Print(buf.String())
func WithStdout(w io.Writer) Option {
	return func(o *options) {
		o.stdout = w
	}
}

// WithStderr sets the standard error writer for the CLI.
//
//	var buf strings.Builder
//	cli.Execute(root, cli.WithStderr(&buf))
//	fmt.Print(buf.String())
func WithStderr(w io.Writer) Option {
	return func(o *options) {
		o.stderr = w
	}
}

// WithContext sets the context for the CLI, enabling cancellation
// and passing request-scoped values.
//
//	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
//	defer cancel()
//
//	cli.Execute(root, cli.WithContext(ctx))
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithTheme sets the theme for styling the CLI help output.
//
//	theme := cli.DefaultTheme()
//	theme.Header = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("141"))
//
//	cli.Execute(root, cli.WithTheme(theme))
func WithTheme(t Theme) Option {
	return func(o *options) {
		o.theme = t
	}
}

// WithWidth sets the maximum width for word wrapping CLI help output.
// Text will wrap at word boundaries to fit within the specified width.
// The default width is 80 characters. Set to 0 to disable wrapping.
//
//	cli.Execute(root, cli.WithWidth(100))
func WithWidth(w int) Option {
	return func(o *options) {
		o.width = w
	}
}

// WithoutManpage disables the hidden "man" command that generates a manpage.
// By default, a hidden "man" command is available that outputs a roff-formatted
// manpage which can be installed by piping to a file in your manpath.
func WithoutManpage() Option {
	return func(o *options) {
		o.manpages = false
	}
}

// WithVersionFlag adds a --version / -V flag to the root command that displays
// build information. The output is styled according to the configured theme.
//
//	cli.Execute(root,
//	    cli.WithVersionFlag(cli.VersionInfo{
//	        Version:   "0.5.0",
//	        GitCommit: "abc1234",
//	        GitBranch: "main",
//	        BuildDate: "2024-01-15T10:30:00Z",
//	        GoVersion: runtime.Version(),
//	    }),
//	)
func WithVersionFlag(info VersionInfo) Option {
	return func(o *options) {
		o.version = &info
		o.versionCommand = false
	}
}

// WithVersionCommand adds a "version" subcommand to the root command that
// displays build information. The subcommand supports additional flags for
// different output formats.
//
// Flags:
//
//   - --short: Display only the version number
//
//   - --json: Display version information as JSON
//
//     cli.Execute(root,
//     cli.WithVersionCommand(cli.VersionInfo{
//     Version:   "0.5.0",
//     GitCommit: "abc1234",
//     GitBranch: "main",
//     BuildDate: "2024-01-15T10:30:00Z",
//     GoVersion: runtime.Version(),
//     }),
//     )
func WithVersionCommand(info VersionInfo) Option {
	return func(o *options) {
		o.version = &info
		o.versionCommand = true
	}
}

// WithCompletionCommand adds a "completion" subcommand that generates shell
// completion scripts. By default, it supports bash, zsh, and fish shells.
//
// Basic usage:
//
//	cli.Execute(root, cli.WithCompletionCommand())
//
// With custom completions:
//
//	cli.Execute(root,
//	    cli.WithCompletionCommand(
//	        cli.CompleteFlag("config", cli.Files(".yaml", ".json")),
//	        cli.CompleteFlag("format", cli.Values("json", "text")),
//	    ),
//	)
//
// With additional shells:
//
//	cli.Execute(root,
//	    cli.WithCompletionCommand(
//	        cli.WithExtraShells(cli.ShellPowerShell, cli.ShellNushell),
//	    ),
//	)
func WithCompletionCommand(opts ...CompletionOption) Option {
	return func(o *options) {
		o.completion = defaultCompletionOptions()
		for _, opt := range opts {
			opt(o.completion)
		}
	}
}

const flagGroupAnnotation = "purpleclay_cli_group"

// FlagGroup assigns flags to a named group for organized help output.
// Grouped flags are rendered under their group header instead of the
// default FLAGS section. Groups appear in the order they are defined.
//
//	cli.FlagGroup(cmd, "Authentication", "token", "api-key")
//	cli.FlagGroup(cmd, "Output Options", "format", "output")
func FlagGroup(cmd *cobra.Command, group string, flags ...string) {
	for _, name := range flags {
		if f := cmd.Flags().Lookup(name); f != nil {
			if f.Annotations == nil {
				f.Annotations = make(map[string][]string)
			}
			f.Annotations[flagGroupAnnotation] = []string{group}
		}
	}
}

// Execute runs the provided cobra command with custom help rendering
// and sensible defaults. Options can be provided to customise behavior.
//
// Basic usage:
//
//	root := &cobra.Command{
//	    Use:   "myapp",
//	    Short: "A brief description",
//	    Long: `
//	        A longer description that spans multiple lines.
//	        Indentation is automatically removed.
//	    `,
//	}
//
//	if err := cli.Execute(root); err != nil {
//	    os.Exit(1)
//	}
//
// Functional options allow customisation:
//
//	cli.Execute(root,
//	    cli.WithContext(ctx),
//	    cli.WithStdout(os.Stdout),
//	    cli.WithStderr(os.Stderr),
//	)
func Execute(cmd *cobra.Command, opts ...Option) error {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	cmd.SetOut(o.stdout)
	cmd.SetErr(o.stderr)
	cmd.SetHelpFunc(helpFunc(o.theme, o.width))
	cmd.SetUsageFunc(usageFunc(o.theme, o.width))
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.TraverseChildren = true

	if o.manpages {
		cmd.AddCommand(&cobra.Command{
			Use:                   "man",
			Short:                 "Generate a manpage for the CLI",
			SilenceUsage:          true,
			DisableFlagsInUseLine: true,
			Hidden:                true,
			Args:                  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				page, err := mango.NewManPage(1, cmd)
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(o.stdout, page.Build(roff.NewDocument()))
				return err
			},
		})
	}

	if o.version != nil {
		if o.versionCommand {
			cmd.AddCommand(newVersionCommand(o.version, o.theme))
		} else {
			cmd.Version = renderVersion(o.version, o.theme)
			cmd.SetVersionTemplate("{{.Version}}")
			cmd.Flags().BoolP("version", "V", false, "print build time version information")
		}
	}

	if o.completion != nil {
		cmd.AddCommand(newCompletionCommand(o.completion, cmd.Name()))
		applyCompletions(cmd, o.completion)
	}

	if err := applyEnvBindings(cmd); err != nil {
		return err
	}

	return cmd.ExecuteContext(o.ctx)
}
