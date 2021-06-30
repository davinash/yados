package server

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start and join new server in a cluster",
		Args:  utils.ExactArgs(0),
		RunE:  utils.ShowHelp(),
	}
	AddNodeStartCmd(cmd)
	rootCmd.AddCommand(cmd)
}
