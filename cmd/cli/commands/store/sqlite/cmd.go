package sqlite

import "github.com/spf13/cobra"

//AddCommands commands related to server
func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "sqlite",
		Short: "manage and monitor store for sqlite store ( sqlite )",
	}

	CreateExecuteQueryCommand(cmd)
	CreateQueryCommand(cmd)

	rootCmd.AddCommand(cmd)
}
