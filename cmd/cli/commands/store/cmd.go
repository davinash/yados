package store

import (
	"github.com/davinash/yados/cmd/cli/commands/store/sql"
	"github.com/spf13/cobra"
)

//NewStoreCommands commands related to server
func NewStoreCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "manage and monitor store",
	}
	CreateCommand(cmd)
	CreateListCommand(cmd)
	CreateDeleteCommand(cmd)

	sql.AddCommands(cmd)

	return cmd
}
