package store

import (
	"github.com/davinash/yados/cmd/cli/commands/store/kv"
	"github.com/davinash/yados/cmd/cli/commands/store/sqlite"
	"github.com/spf13/cobra"
)

//AddCommands commands related to server
func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "manage and monitor store",
	}
	CreateCommand(cmd)
	CreateListCommand(cmd)
	CreateDeleteCommand(cmd)

	kv.AddCommands(cmd)
	sqlite.AddCommands(cmd)

	rootCmd.AddCommand(cmd)
}
