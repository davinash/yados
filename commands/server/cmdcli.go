package server

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddCliCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "server operation commands",
		Args:  utils.ExactArgs(0),
		RunE:  utils.ShowHelp(),
	}
	AddServerRemoveCmd(cmd)
	AddServerStatusCmd(cmd)
	AddServerStopCmd(cmd)
	rootCmd.AddCommand(cmd)
}
