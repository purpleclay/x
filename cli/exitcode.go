package cli

import (
	"sync"

	"github.com/spf13/cobra"
)

// ExitCode documents a single exit code returned by a command, rendered
// in a dedicated EXIT CODES section of the help output.
type ExitCode struct {
	Code int
	Desc string
}

var (
	exitCodeRegistryMu sync.RWMutex
	exitCodeRegistry   = map[*cobra.Command][]ExitCode{}
)

// ExitCodes documents the exit codes a command can return, rendered as a
// dedicated EXIT CODES section in the help output.
//
//	cli.ExitCodes(cmd,
//	    cli.ExitCode{Code: 0, Desc: "success"},
//	    cli.ExitCode{Code: 1, Desc: "drift detected"},
//	    cli.ExitCode{Code: 2, Desc: "error"},
//	)
func ExitCodes(cmd *cobra.Command, codes ...ExitCode) {
	cloned := append([]ExitCode(nil), codes...)

	exitCodeRegistryMu.Lock()
	exitCodeRegistry[cmd] = cloned
	exitCodeRegistryMu.Unlock()
}

func getExitCodes(cmd *cobra.Command) []ExitCode {
	exitCodeRegistryMu.RLock()
	defer exitCodeRegistryMu.RUnlock()
	return exitCodeRegistry[cmd]
}
