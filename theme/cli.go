package theme

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/purpleclay/x/cli"
)

// Adaptive colors for CLI help rendering.
// Each color adapts for readability on light and dark terminals.
var (
	// CommandText styles command and subcommand names.
	CommandText = compat.AdaptiveColor{Light: Purple400, Dark: Purple50}

	// CommentText styles comment lines in examples.
	CommentText = compat.AdaptiveColor{Light: Green600, Dark: Green50}

	// EnvVarText styles environment variable names in examples.
	EnvVarText = compat.AdaptiveColor{Light: Blue400, Dark: Blue100}

	// EnvVarValueText styles environment variable values in examples (dimmer than name).
	EnvVarValueText = compat.AdaptiveColor{Light: Blue600, Dark: Blue300}

	// FlagText styles flag names.
	FlagText = compat.AdaptiveColor{Light: Orange500, Dark: Orange50}

	// FlagMetaText styles flag metadata such as type hints and default values.
	FlagMetaText = compat.AdaptiveColor{Light: Purple500, Dark: Purple100}

	// OperatorText styles shell operators in examples.
	OperatorText = compat.AdaptiveColor{Light: Red500, Dark: Red50}
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
		Header:      H5,
		Operator:    Bold.Foreground(OperatorText),
	}
}
