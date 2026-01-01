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

	// FlagText styles flag names.
	FlagText = lipgloss.AdaptiveColor{Light: string(Orange500), Dark: string(Orange50)}

	// FlagArgText styles flag argument placeholders.
	FlagArgText = lipgloss.AdaptiveColor{Light: string(Purple400), Dark: string(Purple50)}

	// FlagDefaultText styles default value indicators.
	FlagDefaultText = lipgloss.AdaptiveColor{Light: string(Purple500), Dark: string(Purple100)}

	// HeaderBackground styles section header backgrounds.
	HeaderBackground = lipgloss.AdaptiveColor{Light: string(Purple600), Dark: string(Purple400)}
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
	return cli.Theme{
		Command:     lipgloss.NewStyle().Foreground(CommandText).Bold(true),
		Comment:     lipgloss.NewStyle().Foreground(CommentText),
		Description: lipgloss.NewStyle(),
		Flag:        lipgloss.NewStyle().Foreground(FlagText).Bold(true),
		FlagArg:     lipgloss.NewStyle().Foreground(FlagArgText),
		FlagDefault: lipgloss.NewStyle().Foreground(FlagDefaultText),
		Header:      lipgloss.NewStyle().Background(HeaderBackground).Foreground(BrightWhite).Bold(true).Padding(0, 1).MarginBottom(1),
	}
}
