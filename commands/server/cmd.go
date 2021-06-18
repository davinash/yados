package server

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage & Monitor Nodes in YADOS Cluster",
		Args:  utils.ExactArgs(0),
		RunE:  utils.ShowHelp(),
	}
	AddNodeStartCmd(cmd)
	AddNodeStopCmd(cmd)
	AddNodeStatusCmd(cmd)
	rootCmd.AddCommand(cmd)
}
