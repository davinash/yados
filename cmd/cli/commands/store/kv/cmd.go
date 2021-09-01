package kv

import (
	"github.com/spf13/cobra"
)

//AddCommands commands related to server
func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "kv",
		Short: "manage and monitor store for in-memory key value store",
	}
	CreatePutCommand(cmd)
	CreateGetCommand(cmd)

	rootCmd.AddCommand(cmd)
}
