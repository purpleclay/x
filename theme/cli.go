package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/x/cli"
)

// Adaptive colors for CLI help rendering.
// Each color adapts for readability on light and dark terminals.
var (
	// CommandText styles command and subcommand names.
	CommandText = lipgloss.AdaptiveColor{Light: string(Purple400), Dark: string(Purple50)}

	// CommentText styles comment lines in examples.
	CommentText = lipgloss.AdaptiveColor{Light: string(Green600), Dark: string(Green50)}

	// EnvVarText styles environment variable names in examples.
	EnvVarText = lipgloss.AdaptiveColor{Light: string(Blue400), Dark: string(Blue100)}

	// EnvVarValueText styles environment variable values in examples (dimmer than name).
	EnvVarValueText = lipgloss.AdaptiveColor{Light: string(Blue600), Dark: string(Blue300)}

	// FlagText styles flag names.
	FlagText = lipgloss.AdaptiveColor{Light: string(Orange500), Dark: string(Orange50)}

	// FlagMetaText styles flag metadata such as type hints and default values.
	FlagMetaText = lipgloss.AdaptiveColor{Light: string(Purple500), Dark: string(Purple100)}

	// OperatorText styles shell operators in examples.
	OperatorText = lipgloss.AdaptiveColor{Light: string(Red500), Dark: string(Red50)}
)

// PurpleClayCLI returns the official PurpleClay CLI theme. Colors adapt
// automatically for light and dark terminals, and can be integrated with
// the cli package.
//
//	import (
//	    "github.com/purpleclay/x/cli"
//	    "github.com/purpleclay/x/theme"
//	)
//
//	func main() {
//	    root := &cobra.Command{Use: "myapp"}
//	    cli.Execute(root, cli.WithTheme(theme.PurpleClayCLI()))
//	}
func PurpleClayCLI() cli.Theme {
	flagMeta := lipgloss.NewStyle().Foreground(FlagMetaText)

	return cli.Theme{
		Command:     Bold.Foreground(CommandText),
		Comment:     lipgloss.NewStyle().Foreground(CommentText),
		Description: lipgloss.NewStyle(),
		EnvVar:      Bold.Foreground(EnvVarText),
		EnvVarValue: lipgloss.NewStyle().Foreground(EnvVarValueText),
		Flag:        Bold.Foreground(FlagText),
		FlagDefault: flagMeta,
		FlagType:    flagMeta,
		Header:      H5.MarginBottom(1),
		Operator:    Bold.Foreground(OperatorText),
	}
}
