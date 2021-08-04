package server

import (
	"github.com/spf13/cobra"
)

//AddCommands commands related to server
func AddCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "YADOS server manage and monitor commands",
	}
	ListCommands(cmd)
	StartCommands(cmd)
	StopCommands(cmd)
	StatusCommands(cmd)

	rootCmd.AddCommand(cmd)
}
