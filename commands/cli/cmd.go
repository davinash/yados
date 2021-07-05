package cli

import (
	"github.com/davinash/yados/commands/object"
	"github.com/davinash/yados/commands/server"
	"github.com/davinash/yados/commands/store"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	server.AddCliCommands(rootCmd)
	object.AddCommands(rootCmd)
	store.AddCommands(rootCmd)
}
