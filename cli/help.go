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
	fmt.Fprintf(w, "  %s\n", formatUsage(cmd))

	if hasSubCommands(cmd) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("COMMANDS"))
		renderCommands(w, cmd, theme, width)
	}

	if cmd.Example != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("EXAMPLES"))
		renderExamples(w, dedent(cmd.Example), theme)
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

func formatUsage(cmd *cobra.Command) string {
	var parts []string

	parts = append(parts, cmd.CommandPath())

	if cmd.HasAvailableFlags() {
		parts = append(parts, "[FLAGS]")
	}

	if args := extractArgs(cmd.Use); args != "" {
		parts = append(parts, args)
	}

	if hasSubCommands(cmd) {
		parts = append(parts, "[COMMAND]")
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
		wrapped := wrapText(f.Usage, descWidth)
		for line := range strings.SplitSeq(wrapped, "\n") {
			fmt.Fprintf(w, "          %s\n", theme.Description.Render(line))
		}

		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
			fmt.Fprintf(w, "\n          %s\n", theme.FlagDefault.Render(fmt.Sprintf("[default: %s]", f.DefValue)))
		}
	})
}

func renderExamples(w io.Writer, s string, theme Theme) {
	for line := range strings.SplitSeq(s, "\n") {
		if line == "" {
			fmt.Fprintln(w)
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			fmt.Fprintf(w, "  %s\n", theme.Comment.Render(line))
		} else {
			fmt.Fprintf(w, "  %s\n", line)
		}
	}
}
