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

	// EnvVar styles the environment variable name in the env binding hint
	// (e.g., GPG_PRIVATE_KEY in [env: GPG_PRIVATE_KEY]).
	EnvVar lipgloss.Style

	// EnvVarValue styles the environment variable value in the env binding hint
	// (e.g., my-key in [env: GPG_PRIVATE_KEY=my-key]).
	EnvVarValue lipgloss.Style

	// Flag styles flag names including short and long forms
	// (e.g., -v, --verbose).
	Flag lipgloss.Style

	// FlagDefault styles the default value indicator shown beneath flags
	// (e.g., [default: 8080]).
	FlagDefault lipgloss.Style

	// FlagType styles the type hint for flags that accept values
	// (e.g., <string> in --config <string>).
	FlagType lipgloss.Style

	// Header styles section headings such as USAGE, COMMANDS, FLAGS,
	// GLOBAL FLAGS, and EXAMPLES.
	Header lipgloss.Style

	// Operator styles shell operators in the EXAMPLES section
	// (e.g., |, >, >>, <, &&, ||, ;).
	Operator lipgloss.Style
}

// DefaultTheme returns a theme with no styling applied.
func DefaultTheme() Theme {
	return Theme{
		Command:     lipgloss.NewStyle(),
		Comment:     lipgloss.NewStyle(),
		Description: lipgloss.NewStyle(),
		EnvVar:      lipgloss.NewStyle(),
		EnvVarValue: lipgloss.NewStyle(),
		Flag:        lipgloss.NewStyle(),
		FlagDefault: lipgloss.NewStyle(),
		FlagType:    lipgloss.NewStyle(),
		Header:      lipgloss.NewStyle(),
		Operator:    lipgloss.NewStyle(),
	}
}
