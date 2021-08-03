package server

import "github.com/spf13/cobra"

//StopCommands command to stop one or more servers
func StopCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop one or more server",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	rootCmd.AddCommand(cmd)
}
