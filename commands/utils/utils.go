package utils

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ExactArgs returns an error if there is not the exact number of args
func ExactArgs(number int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == number {
			return nil
		}
		return errors.Errorf(
			"%q requires exactly %d %s.\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
			cmd.CommandPath(),
			number,
			pluralize("argument", number),
			cmd.CommandPath(),
			cmd.UseLine(),
			cmd.Short,
		)
	}
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}

//ShowHelp shows the help
func ShowHelp() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}
