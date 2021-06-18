package commands

import (
	"github.com/davinash/yados/commands/server"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	server.AddCommands(rootCmd)
}
