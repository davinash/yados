package server

import "github.com/spf13/cobra"

//StartCommands command for starting a server
func StartCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "create and start a new server",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	rootCmd.AddCommand(cmd)
}
