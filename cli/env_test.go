package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindEnv(t *testing.T) {
	t.Setenv("TEST_KEY", "from-env")

	var buf bytes.Buffer
	var val string

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().StringVar(&val, "key", "default", "test flag")
	BindEnv(cmd.Flags().Lookup("key"), "TEST_KEY")

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)
	assert.Equal(t, "from-env", val)
}

func TestBindEnvExplicitFlagOverrides(t *testing.T) {
	t.Setenv("TEST_KEY", "from-env")

	var buf bytes.Buffer
	var val string

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().StringVar(&val, "key", "", "test flag")
	BindEnv(cmd.Flags().Lookup("key"), "TEST_KEY")
	cmd.SetArgs([]string{"--key=from-flag"})

	err := Execute(cmd, WithStdout(&buf))
	require.NoError(t, err)
	assert.Equal(t, "from-flag", val)
}

func TestBindEnvInvalidValueReturnsError(t *testing.T) {
	t.Setenv("TEST_PORT", "not-a-number")

	var buf bytes.Buffer

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	cmd.Flags().Int("port", 8080, "test flag")
	BindEnv(cmd.Flags().Lookup("port"), "TEST_PORT")

	err := Execute(cmd, WithStdout(&buf))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid value for --port from environment variable TEST_PORT")
}
