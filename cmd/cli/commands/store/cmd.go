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

	rootCmd.AddCommand(cmd)
}
