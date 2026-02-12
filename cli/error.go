package cli

import (
	"fmt"
	"io"
	"strings"
)

// ErrorHandler renders an error to the given writer using the provided theme.
type ErrorHandler func(w io.Writer, theme Theme, err error)

// DefaultErrorHandler prints a styled ERROR badge followed by the error message.
// For usage errors (unknown flags, missing arguments, etc.), it appends a hint
// to try --help.
func DefaultErrorHandler(w io.Writer, theme Theme, err error) {
	fmt.Fprintln(w, theme.ErrorHeader.SetString("ERROR").String())
	fmt.Fprintln(w)
	fmt.Fprintln(w, theme.ErrorDetails.Render(err.Error()))

	if isUsageError(err) {
		fmt.Fprintln(w)
		fmt.Fprintln(w, theme.ErrorDetails.Render("Try --help for usage."))
	}
}

// https://github.com/spf13/cobra/pull/2266
func isUsageError(err error) bool {
	s := err.Error()
	for _, prefix := range []string{
		"flag needs an argument:",
		"unknown flag:",
		"unknown shorthand flag:",
		"unknown command",
		"invalid argument",
		"accepts ",
		"required flag",
	} {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
