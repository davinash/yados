package server

import "github.com/spf13/cobra"

//StatusCommands command for showing the status of one more servers
func StatusCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "status of the one or more servers",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	rootCmd.AddCommand(cmd)
}
