package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func helpFunc() func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, _ []string) {
		renderHelp(cmd.OutOrStdout(), cmd)
	}
}

func usageFunc() func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		renderHelp(cmd.OutOrStderr(), cmd)
		return nil
	}
}

func renderHelp(w io.Writer, cmd *cobra.Command) {
	if desc := cmd.Long; desc != "" {
		fmt.Fprintln(w, dedent(desc))
		fmt.Fprintln(w)
	} else if desc := cmd.Short; desc != "" {
		fmt.Fprintln(w, dedent(desc))
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "USAGE")
	fmt.Fprintf(w, "  %s\n", formatUsage(cmd))

	if hasSubCommands(cmd) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "COMMANDS")
		renderCommands(w, cmd)
	}

	if cmd.Example != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "EXAMPLES")
		fmt.Fprintf(w, "%s\n", indentLines(dedent(cmd.Example), "  "))
	}

	if cmd.HasAvailableLocalFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "FLAGS")
		renderFlags(w, cmd.LocalFlags())
	}

	if cmd.HasAvailableInheritedFlags() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "GLOBAL FLAGS")
		renderFlags(w, cmd.InheritedFlags())
	}
}

func formatUsage(cmd *cobra.Command) string {
	var parts []string

	parts = append(parts, cmd.CommandPath())

	if cmd.HasAvailableFlags() {
		parts = append(parts, "[OPTIONS]")
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

func renderCommands(w io.Writer, cmd *cobra.Command) {
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
		fmt.Fprintf(w, "  %s%s%s\n", sub.Name(), padding, sub.Short)
	}
}

func renderFlags(w io.Writer, flags *pflag.FlagSet) {
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
			flagStr += fmt.Sprintf(" <%s>", strings.ToUpper(f.Name))
		}

		fmt.Fprintf(w, "  %s\n", flagStr)
		fmt.Fprintf(w, "          %s\n", f.Usage)

		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
			fmt.Fprintf(w, "\n          [default: %s]\n", f.DefValue)
		}
	})
}

func indentLines(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}
