package cli

import (
	"fmt"
	"maps"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Shell represents a supported shell for completion.
type Shell string

const (
	// ShellBash represents the Bourne Again Shell.
	ShellBash Shell = "bash"
	// ShellZsh represents the Z Shell.
	ShellZsh Shell = "zsh"
	// ShellFish represents the Friendly Interactive Shell.
	ShellFish Shell = "fish"
	// ShellPowerShell represents PowerShell.
	ShellPowerShell Shell = "powershell"
	// ShellElvish represents the Elvish shell.
	ShellElvish Shell = "elvish"
	// ShellNushell represents Nushell.
	ShellNushell Shell = "nushell"
	// ShellIon represents the Ion shell.
	ShellIon Shell = "ion"
	// ShellOil represents the Oil shell.
	ShellOil Shell = "oil"
	// ShellTcsh represents the TENEX C Shell.
	ShellTcsh Shell = "tcsh"
	// ShellXonsh represents Xonsh.
	ShellXonsh Shell = "xonsh"
	// ShellCmd represents Windows Command Prompt.
	ShellCmd Shell = "cmd"
)

// shellInfo contains metadata for a shell.
type shellInfo struct {
	name        string
	description string
	example     string
}

var shellRegistry = map[Shell]shellInfo{
	ShellBash:       {"bash", "Bourne Again Shell", "source <(%s completion bash)"},
	ShellZsh:        {"zsh", "Z Shell", "source <(%s completion zsh)"},
	ShellFish:       {"fish", "Friendly Interactive Shell", "%s completion fish | source"},
	ShellPowerShell: {"powershell", "PowerShell", "%s completion powershell | Out-String | Invoke-Expression"},
	ShellElvish:     {"elvish", "Elvish", "eval (%s completion elvish | slurp)"},
	ShellNushell:    {"nushell", "Nushell", "%s completion nushell | save ~/.cache/nushell/completion.nu"},
	ShellIon:        {"ion", "Ion Shell", "eval $(%s completion ion)"},
	ShellOil:        {"oil", "Oil Shell", "source <(%s completion oil)"},
	ShellTcsh:       {"tcsh", "TENEX C Shell", "eval `%s completion tcsh`"},
	ShellXonsh:      {"xonsh", "Xonsh", "exec($(%s completion xonsh))"},
	ShellCmd:        {"cmd", "Windows Command Prompt", "%s completion cmd > completion.bat"},
}

// DefaultShells returns the default set of supported shells.
func DefaultShells() []Shell {
	return []Shell{ShellBash, ShellZsh, ShellFish}
}

// Completer defines a completion source that can be converted to a carapace Action.
type Completer interface {
	toAction() carapace.Action
}

// filesCompleter completes file paths.
type filesCompleter struct {
	extensions []string
}

func (c filesCompleter) toAction() carapace.Action {
	return carapace.ActionFiles(c.extensions...)
}

// Files returns a [Completer] for file paths, optionally filtered by extension.
//
//	cli.CompleteFlag("config", cli.Files(".yaml", ".yml", ".json"))
func Files(extensions ...string) Completer {
	return filesCompleter{extensions: extensions}
}

// directoriesCompleter completes directory paths.
type directoriesCompleter struct{}

func (c directoriesCompleter) toAction() carapace.Action {
	return carapace.ActionDirectories()
}

// Directories returns a [Completer] for directory paths.
//
//	cli.CompletePositional(0, cli.Directories())
func Directories() Completer {
	return directoriesCompleter{}
}

// valuesCompleter completes from a fixed list.
type valuesCompleter struct {
	values []string
}

func (c valuesCompleter) toAction() carapace.Action {
	return carapace.ActionValues(c.values...)
}

// Values returns a [Completer] for a fixed list of values.
//
//	cli.CompleteFlag("format", cli.Values("json", "text", "yaml"))
func Values(values ...string) Completer {
	return valuesCompleter{values: values}
}

// valuesDescribedCompleter completes from a list with descriptions.
type valuesDescribedCompleter struct {
	pairs []string
}

func (c valuesDescribedCompleter) toAction() carapace.Action {
	return carapace.ActionValuesDescribed(c.pairs...)
}

// ValuesDescribed returns a [Completer] for values with descriptions.
// Arguments are provided as pairs: value1, description1, value2, description2, ...
//
//	cli.CompleteFlag("format", cli.ValuesDescribed(
//	    "json", "JSON format",
//	    "yaml", "YAML format",
//	))
func ValuesDescribed(pairs ...string) Completer {
	return valuesDescribedCompleter{pairs: pairs}
}

// executablesCompleter completes executable names.
type executablesCompleter struct{}

func (c executablesCompleter) toAction() carapace.Action {
	return carapace.ActionExecutables()
}

// Executables returns a [Completer] for executable names from PATH.
//
//	cli.CompleteFlag("editor", cli.Executables())
func Executables() Completer {
	return executablesCompleter{}
}

// noneCompleter disables completion.
type noneCompleter struct{}

func (c noneCompleter) toAction() carapace.Action {
	return carapace.ActionValues()
}

// None returns a [Completer] that disables file completion fallback.
//
//	cli.CompleteFlag("password", cli.None())
func None() Completer {
	return noneCompleter{}
}

// actionFuncCompleter wraps a function returning a carapace.Action.
type actionFuncCompleter struct {
	fn func() carapace.Action
}

func (c actionFuncCompleter) toAction() carapace.Action {
	return c.fn()
}

// ActionFunc returns a [Completer] from a function that produces a [carapace.Action].
// This is the escape hatch for advanced carapace usage.
//
//	cli.CompleteFlag("branch", cli.ActionFunc(func() carapace.Action {
//	    return carapace.ActionExecCommand("git", "branch", "--format=%(refname:short)")(
//	        func(output []byte) carapace.Action {
//	            branches := strings.Split(strings.TrimSpace(string(output)), "\n")
//	            return carapace.ActionValues(branches...)
//	        },
//	    )
//	}))
func ActionFunc(fn func() carapace.Action) Completer {
	return actionFuncCompleter{fn: fn}
}

// CompletionOption configures shell completion behavior.
type CompletionOption func(*completionOptions)

type completionOptions struct {
	shells        []Shell
	flags         map[string]Completer
	positional    map[int]Completer
	positionalAny Completer
	subcommands   map[string]*completionOptions
}

func defaultCompletionOptions() *completionOptions {
	return &completionOptions{
		shells: DefaultShells(),
	}
}

// WithShells sets the supported shells, replacing the defaults.
// The help output and completion command will only show these shells.
//
//	cli.WithCompletionCommand(
//	    cli.WithShells(cli.ShellBash, cli.ShellZsh, cli.ShellPowerShell),
//	)
func WithShells(shells ...Shell) CompletionOption {
	return func(o *completionOptions) {
		o.shells = shells
	}
}

// WithExtraShells adds shells to the default set (bash, zsh, fish).
//
//	cli.WithCompletionCommand(
//	    cli.WithExtraShells(cli.ShellPowerShell, cli.ShellNushell),
//	)
func WithExtraShells(shells ...Shell) CompletionOption {
	return func(o *completionOptions) {
		o.shells = append(o.shells, shells...)
	}
}

// CompleteFlag defines completion for a flag.
//
//	cli.WithCompletionCommand(
//	    cli.CompleteFlag("config", cli.Files(".yaml", ".json")),
//	    cli.CompleteFlag("format", cli.Values("json", "text")),
//	)
func CompleteFlag(flag string, completer Completer) CompletionOption {
	return func(o *completionOptions) {
		if o.flags == nil {
			o.flags = make(map[string]Completer)
		}
		o.flags[flag] = completer
	}
}

// CompletePositional defines completion for a positional argument (0-indexed).
//
//	cli.WithCompletionCommand(
//	    cli.CompletePositional(0, cli.Directories()),
//	)
func CompletePositional(position int, completer Completer) CompletionOption {
	return func(o *completionOptions) {
		if o.positional == nil {
			o.positional = make(map[int]Completer)
		}
		o.positional[position] = completer
	}
}

// CompletePositionalAny defines completion for remaining positional arguments.
//
//	cli.WithCompletionCommand(
//	    cli.CompletePositionalAny(cli.Files()),
//	)
func CompletePositionalAny(completer Completer) CompletionOption {
	return func(o *completionOptions) {
		o.positionalAny = completer
	}
}

// CompleteSubcommand defines completions for a specific subcommand.
//
//	cli.WithCompletionCommand(
//	    cli.CompleteSubcommand("deploy",
//	        cli.CompleteFlag("environment", cli.Values("dev", "staging", "prod")),
//	        cli.CompletePositional(0, cli.Directories()),
//	    ),
//	)
func CompleteSubcommand(name string, opts ...CompletionOption) CompletionOption {
	return func(o *completionOptions) {
		if o.subcommands == nil {
			o.subcommands = make(map[string]*completionOptions)
		}
		sub := &completionOptions{}
		for _, opt := range opts {
			opt(sub)
		}
		o.subcommands[name] = sub
	}
}

func applyCompletions(cmd *cobra.Command, opts *completionOptions) {
	if opts == nil {
		return
	}
	inferredActions := inferFlagCompletions(cmd)

	if len(opts.flags) > 0 || len(inferredActions) > 0 {
		actions := make(carapace.ActionMap)
		maps.Copy(actions, inferredActions)
		for name, completer := range opts.flags {
			actions[name] = completer.toAction()
		}
		carapace.Gen(cmd).FlagCompletion(actions)
	}

	if len(opts.positional) > 0 {
		maxPos := 0
		for pos := range opts.positional {
			if pos > maxPos {
				maxPos = pos
			}
		}
		actions := make([]carapace.Action, maxPos+1)
		for i := 0; i <= maxPos; i++ {
			if completer, ok := opts.positional[i]; ok {
				actions[i] = completer.toAction()
			} else {
				actions[i] = carapace.ActionValues()
			}
		}
		carapace.Gen(cmd).PositionalCompletion(actions...)
	}

	if opts.positionalAny != nil {
		carapace.Gen(cmd).PositionalAnyCompletion(opts.positionalAny.toAction())
	}

	for _, sub := range cmd.Commands() {
		if subOpts, ok := opts.subcommands[sub.Name()]; ok {
			applyCompletions(sub, subOpts)
		}
	}
}

func inferFlagCompletions(cmd *cobra.Command) carapace.ActionMap {
	actions := make(carapace.ActionMap)

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if helper, ok := f.Value.(EnumHelper); ok {
			var values []string
			for _, entry := range helper.HelpEntries() {
				values = append(values, entry.Name)
			}
			actions[f.Name] = carapace.ActionValues(values...)
		}
	})

	return actions
}

func newCompletionCommand(opts *completionOptions, rootName string) *cobra.Command {
	validArgs := make([]string, len(opts.shells))
	descPairs := make([]string, 0, len(opts.shells)*2)

	var examples strings.Builder
	for i, shell := range opts.shells {
		info := shellRegistry[shell]
		validArgs[i] = info.name
		descPairs = append(descPairs, info.name, info.description)

		if i > 0 {
			examples.WriteString("\n\n")
		}
		fmt.Fprintf(&examples, "# %s\n", info.description)
		fmt.Fprintf(&examples, info.example, rootName)
	}

	cmd := &cobra.Command{
		Use:   "completion <shell>",
		Short: "Generate shell completion scripts for your shell",
		Long: fmt.Sprintf(`Generate shell completion scripts for your shell.

Supported shells: %s`, strings.Join(validArgs, ", ")),
		Example:               examples.String(),
		DisableFlagsInUseLine: true,
		ValidArgs:             validArgs,
		Args:                  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"hideInheritedFlags": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			valid := false
			for _, s := range opts.shells {
				if string(s) == shell {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("unsupported shell: %s", shell)
			}
			snippet, err := carapace.Gen(cmd.Root()).Snippet(shell)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), snippet)
			return nil
		},
	}

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionValuesDescribed(descPairs...),
	)

	return cmd
}
