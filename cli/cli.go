package cli

import (
	"github.com/spf13/cobra"
)

// Configure applies the custom help and usage templates to a cobra command.
// This should be called on the root command before Execute.
func Configure(cmd *cobra.Command) {
	cmd.SetHelpFunc(helpFunc())
	cmd.SetUsageFunc(usageFunc())

	// Disable the default completion command
	cmd.CompletionOptions.DisableDefaultCmd = true
}
