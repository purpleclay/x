package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkFlagRequires(t *testing.T) {
	var buf bytes.Buffer
	var check, workspace bool

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().BoolVarP(&check, "check", "c", false, "check for drift")
	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "use workspace")
	MarkFlagRequires(cmd.Flags().Lookup("workspace"), "check")

	cmd.SetArgs([]string{"--workspace", "--check"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)
}

func TestMarkFlagRequiresMissingRequired(t *testing.T) {
	var buf bytes.Buffer
	var check, workspace bool

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().BoolVarP(&check, "check", "c", false, "check for drift")
	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "use workspace")
	MarkFlagRequires(cmd.Flags().Lookup("workspace"), "check")

	cmd.SetArgs([]string{"--workspace"})

	err := Execute(cmd, WithStdout(&buf))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "flag --workspace requires --check")
}

func TestMarkFlagRequiresIndependentFlagWorks(t *testing.T) {
	var buf bytes.Buffer
	var check, workspace bool

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().BoolVarP(&check, "check", "c", false, "check for drift")
	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "use workspace")
	MarkFlagRequires(cmd.Flags().Lookup("workspace"), "check")

	cmd.SetArgs([]string{"--check"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)
}

func TestMarkFlagRequiresMultipleDependencies(t *testing.T) {
	var buf bytes.Buffer
	var format, output, verbose bool

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().BoolVarP(&format, "format", "f", false, "format output")
	cmd.Flags().BoolVarP(&output, "output", "o", false, "write to file")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	MarkFlagRequires(cmd.Flags().Lookup("format"), "output", "verbose")

	cmd.SetArgs([]string{"--format"})

	err := Execute(cmd, WithStdout(&buf))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "flag --format requires")
	assert.Contains(t, err.Error(), "--output")
	assert.Contains(t, err.Error(), "--verbose")
}

func TestMarkFlagRequiresMultipleDependenciesSatisfied(t *testing.T) {
	var buf bytes.Buffer
	var format, output, verbose bool

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().BoolVarP(&format, "format", "f", false, "format output")
	cmd.Flags().BoolVarP(&output, "output", "o", false, "write to file")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	MarkFlagRequires(cmd.Flags().Lookup("format"), "output", "verbose")

	cmd.SetArgs([]string{"--format", "--output", "--verbose"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)
}

func TestMarkFlagRequiresNilFlag(_ *testing.T) {
	MarkFlagRequires(nil, "check")
}

func TestMarkFlagRequiresPreservesExistingPreRunE(t *testing.T) {
	var buf bytes.Buffer
	var check, workspace bool
	preRunExecuted := false

	cmd := &cobra.Command{
		Use: "test",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			preRunExecuted = true
			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().BoolVarP(&check, "check", "c", false, "check for drift")
	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "use workspace")
	MarkFlagRequires(cmd.Flags().Lookup("workspace"), "check")

	cmd.SetArgs([]string{"--check"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)
	assert.True(t, preRunExecuted)
}

func TestMarkFlagRequiresPreservesExistingPreRun(t *testing.T) {
	var buf bytes.Buffer
	var check, workspace bool
	preRunExecuted := false

	cmd := &cobra.Command{
		Use: "test",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			preRunExecuted = true
		},
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().BoolVarP(&check, "check", "c", false, "check for drift")
	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "use workspace")
	MarkFlagRequires(cmd.Flags().Lookup("workspace"), "check")

	cmd.SetArgs([]string{"--check"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)
	assert.True(t, preRunExecuted)
}
