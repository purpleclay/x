package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func helpFunc(theme Theme) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, _ []string) {
		renderHelp(cmd.OutOrStdout(), cmd, theme)
	}
}

func usageFunc(theme Theme) func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		renderHelp(cmd.OutOrStderr(), cmd, theme)
		return nil
	}
}

func renderHelp(w io.Writer, cmd *cobra.Command, theme Theme) {
	if desc := cmd.Long; desc != "" {
		fmt.Fprintln(w, dedent(desc))
		fmt.Fprintln(w)
	} else if desc := cmd.Short; desc != "" {
		fmt.Fprintln(w, dedent(desc))
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, theme.Header.Render("USAGE"))
	fmt.Fprintf(w, "  %s\n", formatUsage(cmd))

	if hasSubCommands(cmd) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("COMMANDS"))
		renderCommands(w, cmd, theme)
	}

	if cmd.Example != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("EXAMPLES"))
		renderExamples(w, dedent(cmd.Example), theme)
	}

	if cmd.HasAvailableLocalFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("FLAGS"))
		renderFlags(w, cmd.LocalFlags(), theme)
	}

	if cmd.HasAvailableInheritedFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("GLOBAL FLAGS"))
		renderFlags(w, cmd.InheritedFlags(), theme)
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

func renderCommands(w io.Writer, cmd *cobra.Command, theme Theme) {
	maxLen := 0
	for _, sub := range cmd.Commands() {
		if !sub.Hidden && len(sub.Name()) > maxLen {
			maxLen = len(sub.Name())
		}
	}

	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			continue
		}
		padding := strings.Repeat(" ", maxLen-len(sub.Name())+4)
		name := theme.Command.Render(sub.Name())
		desc := theme.Description.Render(sub.Short)
		fmt.Fprintf(w, "  %s%s%s\n", name, padding, desc)
	}
}

func renderFlags(w io.Writer, flags *pflag.FlagSet, theme Theme) {
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

		if f.Value.Type() != "bool" {
			flagStr += " " + theme.FlagArg.Render(fmt.Sprintf("<%s>", strings.ToUpper(f.Name)))
		}

		fmt.Fprintf(w, "  %s\n", theme.Flag.Render(flagStr))
		fmt.Fprintf(w, "          %s\n", theme.Description.Render(f.Usage))

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
