package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const envVarAnnotation = "purpleclay_cli_env"

// BindEnv associates an environment variable with a flag. If the environment
// variable is set and the flag has not been explicitly provided, the
// environment variable value is used as the flag's value during command
// execution.
//
// Precedence (highest to lowest):
//  1. Explicit flag value on command line
//  2. Environment variable
//  3. Flag default value
//
// If flag is nil, BindEnv silently returns without effect (no-op).
//
//	cmd.Flags().StringVarP(&key, "key", "k", "", "GPG private key")
//	cli.BindEnv(cmd.Flags().Lookup("key"), "GPG_PRIVATE_KEY")
func BindEnv(flag *pflag.Flag, envVar string) {
	if flag == nil {
		return
	}

	if flag.Annotations == nil {
		flag.Annotations = make(map[string][]string)
	}
	flag.Annotations[envVarAnnotation] = []string{envVar}
}

// GetEnvVar returns the environment variable associated with a flag,
// or an empty string if no binding exists.
func GetEnvVar(flag *pflag.Flag) string {
	if flag == nil || flag.Annotations == nil {
		return ""
	}
	if envs, ok := flag.Annotations[envVarAnnotation]; ok && len(envs) > 0 {
		return envs[0]
	}
	return ""
}

func applyEnvBindings(cmd *cobra.Command) error {
	var applyErr error

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if applyErr != nil {
			return
		}
		if err := applyEnvToFlag(f); err != nil {
			applyErr = err
		}
	})

	if applyErr != nil {
		return applyErr
	}

	for _, sub := range cmd.Commands() {
		if err := applyEnvBindings(sub); err != nil {
			return err
		}
	}

	return nil
}

func applyEnvToFlag(flag *pflag.Flag) error {
	envVar := GetEnvVar(flag)
	if envVar == "" {
		return nil
	}

	val := os.Getenv(envVar)
	if val == "" || flag.Changed {
		return nil
	}

	if err := flag.Value.Set(val); err != nil {
		return fmt.Errorf("invalid value for --%s from environment variable %s: %w", flag.Name, envVar, err)
	}

	return nil
}
