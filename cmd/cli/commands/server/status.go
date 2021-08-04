package server

import (
	"os"

	"github.com/spf13/cobra"
)

//StatusCommands command for showing the status of one more servers
func StatusCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "status of the one or more servers",
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
