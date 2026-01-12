package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

// VersionInfo contains build-time version information for the CLI.
type VersionInfo struct {
	// Version is the semantic version of the application.
	Version string `json:"version,omitempty"`

	// GitCommit is the git commit hash of the build.
	GitCommit string `json:"git_commit,omitempty"`

	// GitBranch is the git branch of the build.
	GitBranch string `json:"git_branch,omitempty"`

	// BuildDate is the date/time the binary was built.
	BuildDate string `json:"build_date,omitempty"`

	// GoVersion is the Go version used to build the binary.
	GoVersion string `json:"go_version,omitempty"`

	// Platform is the OS/architecture the binary was built for.
	// This is automatically populated from runtime.GOOS/GOARCH if not set.
	Platform string `json:"platform,omitempty"`
}

func renderVersion(info *VersionInfo, theme Theme) string {
	var buf strings.Builder

	buf.WriteString(theme.FlagDefault.Render(info.Version))
	buf.WriteString("\n")

	type field struct {
		label string
		value string
	}

	fields := []field{
		{"Git Commit", info.GitCommit},
		{"Git Branch", info.GitBranch},
		{"Build Date", info.BuildDate},
		{"Go Version", info.GoVersion},
		{"Platform", info.Platform},
	}

	hasFields := false
	for _, f := range fields {
		if f.value != "" {
			hasFields = true
			break
		}
	}

	if !hasFields {
		return buf.String()
	}

	buf.WriteString("\n")
	buf.WriteString(theme.Header.Render("BUILD INFORMATION"))
	buf.WriteString("\n\n")

	for _, f := range fields {
		if f.value == "" {
			continue
		}
		// Pad label before styling to avoid ANSI codes affecting width calculation
		paddedLabel := fmt.Sprintf("%-14s", f.label)
		fmt.Fprintf(&buf, "%s%s\n", theme.Description.Render(paddedLabel), theme.FlagDefault.Render(f.value))
	}

	return buf.String()
}

func renderVersionShort(w io.Writer, info *VersionInfo) {
	fmt.Fprintln(w, info.Version)
}

func renderVersionJSON(w io.Writer, info *VersionInfo) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(info)
}

func newVersionCommand(info *VersionInfo, theme Theme) *cobra.Command {
	var (
		short   bool
		jsonOut bool
	)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print build time version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if jsonOut {
				return renderVersionJSON(cmd.OutOrStdout(), info)
			}
			if short {
				renderVersionShort(cmd.OutOrStdout(), info)
				return nil
			}
			fmt.Fprint(cmd.OutOrStdout(), renderVersion(info, theme))
			return nil
		},
	}

	cmd.Flags().BoolVar(&short, "short", false, "display only the version number")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "display version information as JSON")
	cmd.MarkFlagsMutuallyExclusive("short", "json")

	return cmd
}
