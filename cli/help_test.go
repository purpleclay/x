package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestRenderHelp(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *cobra.Command
		golden string
	}{
		{
			name: "basic command",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
					Long:  "A longer description of the test application.",
				}
			},
			golden: "basic_command.golden",
		},
		{
			name: "with flags",
			setup: func() *cobra.Command {
				var output string
				var verbose bool
				cmd := &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path")
				cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
				return cmd
			},
			golden: "with_flags.golden",
		},
		{
			name: "with subcommands",
			setup: func() *cobra.Command {
				root := &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
				}
				root.AddCommand(
					&cobra.Command{
						Use:   "build",
						Short: "Build the project",
						Run:   func(cmd *cobra.Command, args []string) {},
					},
					&cobra.Command{
						Use:   "init",
						Short: "Initialize a new project",
						Run:   func(cmd *cobra.Command, args []string) {},
					},
				)
				return root
			},
			golden: "with_subcommands.golden",
		},
		{
			name: "with examples",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
					Example: `# Build the project
myapp build --output ./dist

# Initialize a new project
myapp init --name myproject`,
					Run: func(cmd *cobra.Command, args []string) {},
				}
			},
			golden: "with_examples.golden",
		},
		{
			name: "with global flags",
			setup: func() *cobra.Command {
				var debug bool
				root := &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
				}
				root.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
				root.AddCommand(&cobra.Command{
					Use:   "sub",
					Short: "A subcommand",
					Run:   func(cmd *cobra.Command, args []string) {},
				})
				return root
			},
			golden: "with_global_flags.golden",
		},
		{
			name: "hidden commands not shown",
			setup: func() *cobra.Command {
				root := &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
				}
				root.AddCommand(
					&cobra.Command{
						Use:   "visible",
						Short: "A visible command",
						Run:   func(cmd *cobra.Command, args []string) {},
					},
					&cobra.Command{
						Use:    "hidden",
						Short:  "A hidden command",
						Hidden: true,
						Run:    func(cmd *cobra.Command, args []string) {},
					},
				)
				return root
			},
			golden: "hidden_commands.golden",
		},
		{
			name: "hidden flags not shown",
			setup: func() *cobra.Command {
				var visible, hidden string
				cmd := &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				cmd.Flags().StringVar(&visible, "visible", "", "A visible flag")
				cmd.Flags().StringVar(&hidden, "hidden", "", "A hidden flag")
				cmd.Flags().MarkHidden("hidden")
				return cmd
			},
			golden: "hidden_flags.golden",
		},
		{
			name: "with default value",
			setup: func() *cobra.Command {
				var port int
				cmd := &cobra.Command{
					Use:   "myapp",
					Short: "A test application",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				cmd.Flags().IntVarP(&port, "port", "p", 8080, "Server port")
				return cmd
			},
			golden: "with_default_value.golden",
		},
		{
			name: "with arguments",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use:   "myapp <source> [destination]",
					Short: "Copy files",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
			},
			golden: "with_arguments.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup()
			Configure(cmd)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs([]string{"--help"})
			err := cmd.Execute()
			assert.NilError(t, err)

			golden.Assert(t, buf.String(), tt.golden)
		})
	}
}

func TestRenderHelp_SubcommandGlobalFlags(t *testing.T) {
	var debug bool
	root := &cobra.Command{
		Use:   "myapp",
		Short: "A test application",
	}
	root.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
	root.AddCommand(&cobra.Command{
		Use:   "sub",
		Short: "A subcommand",
		Run:   func(cmd *cobra.Command, args []string) {},
	})
	Configure(root)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"sub", "--help"})
	err := root.Execute()
	assert.NilError(t, err)

	golden.Assert(t, buf.String(), "subcommand_global_flags.golden")
}

func TestExtractArgs(t *testing.T) {
	tests := []struct {
		use      string
		expected string
	}{
		{"myapp", ""},
		{"myapp <file>", "<file>"},
		{"myapp <source> [destination]", "<source> [destination]"},
		{"build [flags]", "[flags]"},
	}

	for _, tt := range tests {
		t.Run(tt.use, func(t *testing.T) {
			result := extractArgs(tt.use)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestHasSubCommands(t *testing.T) {
	t.Run("no subcommands", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		assert.Assert(t, !hasSubCommands(cmd))
	})

	t.Run("with visible subcommand", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(&cobra.Command{Use: "sub"})
		assert.Assert(t, hasSubCommands(cmd))
	})

	t.Run("with only hidden subcommand", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(&cobra.Command{Use: "sub", Hidden: true})
		assert.Assert(t, !hasSubCommands(cmd))
	})
}

func TestIndentLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prefix   string
		expected string
	}{
		{
			name:     "single line",
			input:    "hello",
			prefix:   "  ",
			expected: "  hello",
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3",
			prefix:   "  ",
			expected: "  line1\n  line2\n  line3",
		},
		{
			name:     "with empty lines",
			input:    "line1\n\nline3",
			prefix:   "  ",
			expected: "  line1\n\n  line3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := indentLines(tt.input, tt.prefix)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestConfigure(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "myapp",
		Short: "A test application",
	}
	Configure(cmd)

	assert.Assert(t, cmd.CompletionOptions.DisableDefaultCmd)
}
