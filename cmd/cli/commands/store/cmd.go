package store

import "github.com/spf13/cobra"

//AddCommands commands related to server
func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "manage and monitor store",
	}
	CreateCommand(cmd)
	CreateListCommand(cmd)
	CreatePutCommand(cmd)
	CreateGetCommand(cmd)
	CreateDeleteCommand(cmd)

	rootCmd.AddCommand(cmd)
}
