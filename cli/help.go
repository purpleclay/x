package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func helpFunc(theme Theme, width int) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, _ []string) {
		renderHelp(cmd.OutOrStdout(), cmd, theme, width)
	}
}

func usageFunc(theme Theme, width int) func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		renderHelp(cmd.OutOrStderr(), cmd, theme, width)
		return nil
	}
}

func renderHelp(w io.Writer, cmd *cobra.Command, theme Theme, width int) {
	if desc := cmd.Long; desc != "" {
		fmt.Fprintln(w, wrapText(dedent(desc), width))
		fmt.Fprintln(w)
	} else if desc := cmd.Short; desc != "" {
		fmt.Fprintln(w, wrapText(dedent(desc), width))
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, theme.Header.Render("USAGE"))
	fmt.Fprintf(w, "  %s\n", formatUsage(cmd, theme))

	if hasSubCommands(cmd) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("COMMANDS"))
		renderCommands(w, cmd, theme, width)
	}

	if cmd.Example != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("EXAMPLES"))
		renderExamples(w, dedent(cmd.Example), cmd, theme)
	}

	if cmd.HasAvailableLocalFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("FLAGS"))
		renderFlags(w, cmd.LocalFlags(), theme, width)
	}

	if cmd.HasAvailableInheritedFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("GLOBAL FLAGS"))
		renderFlags(w, cmd.InheritedFlags(), theme, width)
	}
}

func formatUsage(cmd *cobra.Command, theme Theme) string {
	var parts []string

	// Style the command path (e.g., "nsv next")
	parts = append(parts, theme.Command.Render(cmd.CommandPath()))

	if cmd.HasAvailableFlags() {
		parts = append(parts, theme.FlagType.Render("[FLAGS]"))
	}

	if args := extractArgs(cmd.Use); args != "" {
		parts = append(parts, theme.FlagType.Render(args))
	}

	if hasSubCommands(cmd) {
		parts = append(parts, theme.FlagType.Render("[COMMAND]"))
	}

	return strings.Join(parts, " ")
}

func extractArgs(use string) string {
	parts := strings.SplitN(use, " ", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func hasSubCommands(cmd *cobra.Command) bool {
	for _, sub := range cmd.Commands() {
		if !sub.Hidden {
			return true
		}
	}
	return false
}

func renderCommands(w io.Writer, cmd *cobra.Command, theme Theme, width int) {
	maxLen := 0
	for _, sub := range cmd.Commands() {
		if !sub.Hidden && len(sub.Name()) > maxLen {
			maxLen = len(sub.Name())
		}
	}

	indent := 2 + maxLen + 4

	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			continue
		}
		padding := strings.Repeat(" ", maxLen-len(sub.Name())+4)
		name := theme.Command.Render(sub.Name())

		descWidth := width - indent
		if descWidth <= 0 || width == 0 {
			descWidth = 0
		}
		wrapped := wrapText(sub.Short, descWidth)
		lines := strings.Split(wrapped, "\n")

		desc := theme.Description.Render(lines[0])
		fmt.Fprintf(w, "  %s%s%s\n", name, padding, desc)

		for _, line := range lines[1:] {
			fmt.Fprintf(w, "%s%s\n", strings.Repeat(" ", indent), theme.Description.Render(line))
		}
	}
}

func flagTypeName(t string) string {
	switch t {
	case "stringSlice", "stringArray":
		return "strings"
	case "intSlice":
		return "ints"
	case "float64":
		return "float"
	case "float64Slice":
		return "floats"
	case "boolSlice":
		return "bools"
	default:
		return t
	}
}

func renderFlags(w io.Writer, flags *pflag.FlagSet, theme Theme, width int) {
	const flagIndent = 10

	first := true
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}

		if !first {
			fmt.Fprintln(w)
		}
		first = false

		var flagStr string
		if f.Shorthand != "" {
			flagStr = fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
		} else {
			flagStr = fmt.Sprintf("    --%s", f.Name)
		}

		flagType := f.Value.Type()
		if flagType != "bool" {
			flagStr += " " + theme.FlagType.Render(fmt.Sprintf("<%s>", flagTypeName(flagType)))
		}

		fmt.Fprintf(w, "  %s\n", theme.Flag.Render(flagStr))

		descWidth := width - flagIndent
		if descWidth <= 0 || width == 0 {
			descWidth = 0
		}

		// Build description with inline default value
		desc := f.Usage
		hasDefault := f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]"
		defaultStr := ""
		if hasDefault {
			defaultStr = fmt.Sprintf("(default: %s)", f.DefValue)
		}

		wrapped := wrapText(desc, descWidth)
		lines := strings.Split(wrapped, "\n")

		for i, line := range lines {
			isLastLine := i == len(lines)-1
			if isLastLine && hasDefault {
				// Append default to last line
				line = line + " " + theme.FlagDefault.Render(defaultStr)
			}
			fmt.Fprintf(w, "          %s\n", theme.Description.Render(line))
		}
	})
}

func renderExamples(w io.Writer, s string, cmd *cobra.Command, theme Theme) {
	// Build set of known subcommand names
	subcommands := make(map[string]bool)
	root := cmd.Root()
	for _, c := range root.Commands() {
		if !c.Hidden {
			subcommands[c.Name()] = true
		}
	}

	for line := range strings.SplitSeq(s, "\n") {
		if line == "" {
			fmt.Fprintln(w)
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			fmt.Fprintf(w, "  %s\n", theme.Comment.Render(line))
		} else {
			styled := styleExampleLine(line, root.Name(), subcommands, theme)
			fmt.Fprintf(w, "  %s\n", styled)
		}
	}
}

func styleExampleLine(line, rootCmd string, subcommands map[string]bool, theme Theme) string {
	tokens := tokenizeExample(line)
	var result strings.Builder

	for i, token := range tokens {
		if token.isWhitespace {
			result.WriteString(token.value)
			continue
		}

		switch {
		case i == 0 && token.value == rootCmd:
			// Root command
			result.WriteString(theme.Command.Render(token.value))
		case subcommands[token.value]:
			// Known subcommand
			result.WriteString(theme.Command.Render(token.value))
		case strings.HasPrefix(token.value, "-"):
			if idx := strings.Index(token.value, "="); idx != -1 {
				flag := token.value[:idx+1]
				value := token.value[idx+1:]
				result.WriteString(theme.Flag.Render(flag))
				result.WriteString(value)
			} else {
				result.WriteString(theme.Flag.Render(token.value))
			}
		default:
			result.WriteString(token.value)
		}
	}

	return result.String()
}

type exampleToken struct {
	value        string
	isWhitespace bool
}

func tokenizeExample(line string) []exampleToken {
	var tokens []exampleToken
	var current strings.Builder
	inWhitespace := false

	for _, r := range line {
		isSpace := r == ' ' || r == '\t'
		if isSpace != inWhitespace {
			if current.Len() > 0 {
				tokens = append(tokens, exampleToken{value: current.String(), isWhitespace: inWhitespace})
				current.Reset()
			}
			inWhitespace = isSpace
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		tokens = append(tokens, exampleToken{value: current.String(), isWhitespace: inWhitespace})
	}

	return tokens
}
