package sql

import "github.com/spf13/cobra"

//AddCommands commands related to server
func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "sql",
		Short: "manage and monitor store for sql store ( sqlite )",
	}

	CreateExecuteQueryCommand(cmd)
	CreateQueryCommand(cmd)

	rootCmd.AddCommand(cmd)
}
