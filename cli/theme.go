package cli

import "github.com/charmbracelet/lipgloss"

// Theme defines the styles used for rendering CLI help output.
// Each field controls the appearance of a specific element.
type Theme struct {
	// Command styles command and subcommand names in the COMMANDS section.
	Command lipgloss.Style

	// Comment styles comment lines within the EXAMPLES section.
	Comment lipgloss.Style

	// Description styles the help text that describes commands and flags.
	Description lipgloss.Style

	// Flag styles flag names including short and long forms
	// (e.g., -v, --verbose).
	Flag lipgloss.Style

	// FlagArg styles the argument placeholder for flags that accept values
	// (e.g., <FILE> in --config <FILE>).
	FlagArg lipgloss.Style

	// FlagDefault styles the default value indicator shown beneath flags
	// (e.g., [default: 8080]).
	FlagDefault lipgloss.Style

	// Header styles section headings such as USAGE, COMMANDS, FLAGS,
	// GLOBAL FLAGS, and EXAMPLES.
	Header lipgloss.Style
}

// DefaultTheme returns a theme with no styling applied.
func DefaultTheme() Theme {
	return Theme{
		Command:     lipgloss.NewStyle(),
		Comment:     lipgloss.NewStyle(),
		Description: lipgloss.NewStyle(),
		Flag:        lipgloss.NewStyle(),
		FlagArg:     lipgloss.NewStyle(),
		FlagDefault: lipgloss.NewStyle(),
		Header:      lipgloss.NewStyle(),
	}
}
