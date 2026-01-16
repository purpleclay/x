package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
		renderGroupedFlags(w, cmd.LocalFlags(), "FLAGS", theme, width)
	}

	if cmd.HasAvailableInheritedFlags() && cmd.Annotations["hideInheritedFlags"] != "true" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render("GLOBAL FLAGS"))
		renderFlags(w, cmd.InheritedFlags(), theme, width)
	}
}

type flagGroup struct {
	name  string
	flags []*pflag.Flag
}

func collectFlagGroups(flags *pflag.FlagSet) (ungrouped []*pflag.Flag, groups []flagGroup) {
	groupOrder := make([]string, 0)
	groupFlags := make(map[string][]*pflag.Flag)
	seen := make(map[string]bool)

	flags.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}

		if ann, ok := f.Annotations[flagGroupAnnotation]; ok && len(ann) > 0 {
			group := ann[0]
			if !seen[group] {
				seen[group] = true
				groupOrder = append(groupOrder, group)
			}
			groupFlags[group] = append(groupFlags[group], f)
		} else {
			ungrouped = append(ungrouped, f)
		}
	})

	for _, name := range groupOrder {
		groups = append(groups, flagGroup{name: name, flags: groupFlags[name]})
	}

	return ungrouped, groups
}

func renderGroupedFlags(w io.Writer, flags *pflag.FlagSet, defaultHeader string, theme Theme, width int) {
	ungrouped, groups := collectFlagGroups(flags)

	if len(ungrouped) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render(defaultHeader))
		renderFlagList(w, ungrouped, theme, width)
	}

	for _, g := range groups {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.Header.Render(strings.ToUpper(g.name)))
		renderFlagList(w, g.flags, theme, width)
	}
}

func formatUsage(cmd *cobra.Command, theme Theme) string {
	var parts []string
	parts = append(parts, theme.Command.Render(cmd.CommandPath()))

	if cmd.HasAvailableFlags() && !cmd.DisableFlagsInUseLine {
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
	case "intSlice", "int32Slice", "int64Slice":
		return "ints"
	case "uintSlice":
		return "uints"
	case "float32", "float64":
		return "float"
	case "float32Slice", "float64Slice":
		return "floats"
	case "boolSlice":
		return "bools"
	case "durationSlice":
		return "durations"
	case "ipSlice":
		return "ips"
	case "ipNetSlice":
		return "cidrs"
	default:
		return t
	}
}

func formatDefaultValue(value, valueType string, style lipgloss.Style) string {
	switch valueType {
	case "string":
		return `"` + style.Render(value) + `"`
	case "stringSlice", "stringArray",
		"durationSlice",
		"ipSlice", "ipNetSlice":
		trimmed := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
		if trimmed == "" {
			return ""
		}
		parts := strings.Split(trimmed, ",")
		quoted := make([]string, len(parts))
		for i, p := range parts {
			quoted[i] = `"` + style.Render(p) + `"`
		}
		return strings.Join(quoted, ", ")
	case "intSlice", "int32Slice", "int64Slice",
		"uintSlice",
		"float32Slice", "float64Slice",
		"boolSlice":
		trimmed := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
		if trimmed == "" {
			return ""
		}
		parts := strings.Split(trimmed, ",")
		styled := make([]string, len(parts))
		for i, p := range parts {
			styled[i] = style.Render(p)
		}
		return strings.Join(styled, ", ")
	default:
		return style.Render(value)
	}
}

func renderFlags(w io.Writer, flags *pflag.FlagSet, theme Theme, width int) {
	var flagList []*pflag.Flag
	flags.VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			flagList = append(flagList, f)
		}
	})
	renderFlagList(w, flagList, theme, width)
}

func formatEnvVar(envVar string, theme Theme) string {
	val := os.Getenv(envVar)
	if val == "" {
		return "[env: " + theme.EnvVar.Render(envVar) + "]"
	}

	if len(val) > 20 {
		val = val[:20] + "..."
	}
	return "[env: " + theme.EnvVar.Render(envVar) + "=" + theme.EnvVarValue.Render(val) + "]"
}

func renderFlagList(w io.Writer, flags []*pflag.Flag, theme Theme, width int) {
	const flagIndent = 10

	for i, f := range flags {
		if i > 0 {
			fmt.Fprintln(w)
		}

		var flagStr string
		if f.Shorthand != "" {
			flagStr = fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
		} else {
			flagStr = fmt.Sprintf("    --%s", f.Name)
		}

		flagType := f.Value.Type()
		if flagType != "bool" {
			if helper, ok := f.Value.(EnumHelper); ok && helper.HasHelp() {
				flagType = helper.BaseType()
			}
			flagStr += " " + theme.FlagType.Render(fmt.Sprintf("<%s>", flagTypeName(flagType)))
		}

		if envVar := GetEnvVar(f); envVar != "" {
			flagStr += "  " + formatEnvVar(envVar, theme)
		}

		fmt.Fprintf(w, "  %s\n", theme.Flag.Render(flagStr))

		descWidth := width - flagIndent
		if descWidth <= 0 || width == 0 {
			descWidth = 0
		}

		desc := f.Usage
		hasDefault := f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]"

		wrapped := wrapText(desc, descWidth)
		lines := strings.Split(wrapped, "\n")

		for j, line := range lines {
			isLastLine := j == len(lines)-1
			if isLastLine && hasDefault {
				valueType := f.Value.Type()
				if helper, ok := f.Value.(EnumHelper); ok {
					valueType = helper.BaseType()
				}
				formatted := formatDefaultValue(f.DefValue, valueType, theme.FlagDefault)
				line = line + " (default: " + formatted + ")"
			}
			fmt.Fprintf(w, "          %s\n", theme.Description.Render(line))
		}

		if helper, ok := f.Value.(EnumHelper); ok && helper.HasHelp() {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "          %s\n", theme.Description.Render("Possible values:"))
			for _, entry := range helper.HelpEntries() {
				if entry.Help != "" {
					fmt.Fprintf(w, "          - %s: %s\n",
						theme.FlagType.Render(entry.Name),
						theme.Description.Render(entry.Help))
				} else {
					fmt.Fprintf(w, "          - %s\n", theme.FlagType.Render(entry.Name))
				}
			}
		}
	}
}

func renderExamples(w io.Writer, s string, cmd *cobra.Command, theme Theme) {
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

	expectCommand := true
	for _, token := range tokens {
		switch token.tokenType {
		case tokenWhitespace:
			result.WriteString(token.value)

		case tokenOperator:
			result.WriteString(theme.Operator.Render(token.value))
			// After a pipe or semicolon, the next word is a command
			if token.value == "|" || token.value == ";" || token.value == "&&" || token.value == "||" {
				expectCommand = true
			}

		case tokenString:
			result.WriteString(token.value)

		case tokenEnvAssign:
			// Style environment variable name and value separately
			if idx := strings.Index(token.value, "="); idx != -1 {
				varName := token.value[:idx+1]
				varValue := token.value[idx+1:]
				result.WriteString(theme.EnvVar.Render(varName))
				result.WriteString(theme.EnvVarValue.Render(varValue))
			} else {
				result.WriteString(theme.EnvVar.Render(token.value))
			}
			// After env var, we still expect a command

		case tokenLineContinuation:
			result.WriteString(theme.Operator.Render(token.value))

		case tokenWord:
			switch {
			case expectCommand && token.value == rootCmd:
				result.WriteString(theme.Command.Render(token.value))
				expectCommand = false
			case subcommands[token.value]:
				result.WriteString(theme.Command.Render(token.value))
				expectCommand = false
			case expectCommand:
				result.WriteString(theme.Command.Render(token.value))
				expectCommand = false
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
	}

	return result.String()
}

// tokenType represents the type of token in an example line.
type tokenType int

const (
	tokenWhitespace tokenType = iota
	tokenWord
	tokenOperator
	tokenString
	tokenEnvAssign
	tokenLineContinuation
)

type exampleToken struct {
	value     string
	tokenType tokenType
}

// shellOperators contains shell operators ordered by length (longest first for proper matching).
var shellOperators = []string{
	">>", "<<", "&&", "||",
	"|", ">", "<", ";", "&",
}

func tokenizeExample(line string) []exampleToken {
	var tokens []exampleToken
	runes := []rune(line)
	i := 0

	for i < len(runes) {
		r := runes[i]

		// Handle whitespace
		if r == ' ' || r == '\t' {
			start := i
			for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t') {
				i++
			}
			tokens = append(tokens, exampleToken{value: string(runes[start:i]), tokenType: tokenWhitespace})
			continue
		}

		// Handle quoted strings
		if r == '"' || r == '\'' {
			quote := r
			start := i
			i++
			for i < len(runes) && runes[i] != quote {
				if runes[i] == '\\' && i+1 < len(runes) {
					i += 2
				} else {
					i++
				}
			}
			if i < len(runes) {
				i++ // consume closing quote
			}
			tokens = append(tokens, exampleToken{value: string(runes[start:i]), tokenType: tokenString})
			continue
		}

		// Handle line continuation (backslash at end of line)
		if r == '\\' && i == len(runes)-1 {
			tokens = append(tokens, exampleToken{value: "\\", tokenType: tokenLineContinuation})
			i++
			continue
		}

		// Handle shell operators
		if op := matchOperator(runes, i); op != "" {
			tokens = append(tokens, exampleToken{value: op, tokenType: tokenOperator})
			i += len([]rune(op))
			continue
		}

		// Handle words (including flags, commands, env assignments)
		start := i
		for i < len(runes) {
			r := runes[i]
			if r == ' ' || r == '\t' || r == '"' || r == '\'' {
				break
			}
			if r == '\\' && i == len(runes)-1 {
				break
			}
			if matchOperator(runes, i) != "" {
				break
			}
			i++
		}

		word := string(runes[start:i])
		if word != "" {
			tt := tokenWord
			// Check for environment variable assignment (VAR=value at start of command)
			// Env vars can appear consecutively: VAR1=val1 VAR2=val2 cmd
			if allEnvOrWhitespace(tokens) {
				if idx := strings.Index(word, "="); idx > 0 && !strings.HasPrefix(word, "-") {
					tt = tokenEnvAssign
				}
			}
			tokens = append(tokens, exampleToken{value: word, tokenType: tt})
		}
	}

	return tokens
}

func matchOperator(runes []rune, pos int) string {
	remaining := string(runes[pos:])
	for _, op := range shellOperators {
		if strings.HasPrefix(remaining, op) {
			return op
		}
	}
	return ""
}

func allEnvOrWhitespace(tokens []exampleToken) bool {
	for _, t := range tokens {
		if t.tokenType != tokenWhitespace && t.tokenType != tokenEnvAssign {
			return false
		}
	}
	return true
}
