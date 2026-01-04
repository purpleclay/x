package cli

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Option is a functional option for configuring the CLI.
type Option func(*options)

type options struct {
	ctx    context.Context
	stdout io.Writer
	stderr io.Writer
	theme  Theme
	width  int
}

func defaultOptions() *options {
	return &options{
		ctx:    context.Background(),
		stdout: os.Stdout,
		stderr: os.Stderr,
		theme:  DefaultTheme(),
		width:  80,
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

	return cmd.ExecuteContext(o.ctx)
}
