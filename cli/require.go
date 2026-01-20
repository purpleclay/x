package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const flagRequiresAnnotation = "purpleclay_cli_requires"

// MarkFlagRequires specifies that if flag is set, the named required flags
// must also be set. This is a one-way dependency - the required flags can
// be used independently.
//
// If flag is nil, MarkFlagRequires silently returns without effect (no-op).
//
//	cmd.Flags().BoolVarP(&check, "check", "c", false, "check for drift")
//	cmd.Flags().BoolVarP(&workspace, "workspace", "w", false, "use workspace")
//
//	cli.MarkFlagRequires(cmd.Flags().Lookup("workspace"), "check")
//
// Multiple dependencies can be specified:
//
//	cli.MarkFlagRequires(cmd.Flags().Lookup("format"), "output", "verbose")
//
// During command execution, if --workspace is provided without --check,
// an error is returned: "flag --workspace requires --check"
func MarkFlagRequires(flag *pflag.Flag, flagNames ...string) {
	if flag == nil {
		return
	}

	if flag.Annotations == nil {
		flag.Annotations = make(map[string][]string)
	}
	flag.Annotations[flagRequiresAnnotation] = append(
		flag.Annotations[flagRequiresAnnotation], flagNames...)
}

// GetFlagRequires returns the flags that the given flag requires,
// or nil if no requirements exist.
func GetFlagRequires(flag *pflag.Flag) []string {
	if flag == nil || flag.Annotations == nil {
		return nil
	}
	if requires, ok := flag.Annotations[flagRequiresAnnotation]; ok {
		return requires
	}
	return nil
}

func addFlagRequirementsValidation(cmd *cobra.Command) {
	existingPreRunE := cmd.PersistentPreRunE
	existingPreRun := cmd.PersistentPreRun

	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if err := validateFlagRequirements(c); err != nil {
			return err
		}

		if existingPreRunE != nil {
			return existingPreRunE(c, args)
		}
		if existingPreRun != nil {
			existingPreRun(c, args)
		}
		return nil
	}
	cmd.PersistentPreRun = nil

	for _, sub := range cmd.Commands() {
		addFlagRequirementsValidation(sub)
	}
}

func validateFlagRequirements(cmd *cobra.Command) error {
	var validateErr error

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if validateErr != nil {
			return
		}
		if err := validateFlagRequires(cmd.Flags(), f); err != nil {
			validateErr = err
		}
	})

	return validateErr
}

func validateFlagRequires(flags *pflag.FlagSet, flag *pflag.Flag) error {
	requires := GetFlagRequires(flag)
	if len(requires) == 0 || !flag.Changed {
		return nil
	}

	var missing []string
	for _, req := range requires {
		reqFlag := flags.Lookup(req)
		if reqFlag == nil || !reqFlag.Changed {
			missing = append(missing, "--"+req)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("flag --%s requires %s", flag.Name, strings.Join(missing, ", "))
	}

	return nil
}
