package yados

import (
	"github.com/davinash/yados/commands/server"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	//server.AddNodeStartCmd(rootCmd)
	server.AddServerCmd(rootCmd)
}
