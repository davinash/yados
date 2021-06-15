package commands

import (
	"github.com/davinash/yados/commands/node"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	node.AddCommands(rootCmd)
}
