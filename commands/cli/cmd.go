package cli

import (
	"github.com/davinash/yados/commands/object"
	"github.com/davinash/yados/commands/store"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command) {
	AddServerRemoveCmd(rootCmd)
	AddServerStatusCmd(rootCmd)
	AddServerStopCmd(rootCmd)
	object.AddCommands(rootCmd)
	store.AddCommands(rootCmd)
}
