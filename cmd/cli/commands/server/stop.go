package server

import (
	"os"

	"github.com/spf13/cobra"
)

//StopCommands command to stop one or more servers
func StopCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop one or more server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				err := cmd.Help()
				if err != nil {
					return err
				}
				os.Exit(0)
			}
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
