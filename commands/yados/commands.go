package yados

import (
	"github.com/davinash/yados/commands/server"
	"github.com/spf13/cobra"
)

// AddCommands Cobra command for yados cli
func AddCommands(rootCmd *cobra.Command) {
	server.AddServerCmd(rootCmd)
}
