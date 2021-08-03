package server

import "github.com/spf13/cobra"

//ListCommands cobra command for listing the servers
func ListCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the servers in the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	rootCmd.AddCommand(cmd)
}
