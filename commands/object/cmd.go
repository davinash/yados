package object

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "object",
		Short: "Object operations",
		Args:  utils.ExactArgs(0),
		RunE:  utils.ShowHelp(),
	}
	AddDeleteObjectCommand(cmd)
	AddGetObjectCommand(cmd)
	AddPutObjectCommand(cmd)
	rootCmd.AddCommand(cmd)
}
