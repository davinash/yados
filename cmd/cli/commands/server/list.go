package server

import (
	"os"

	"github.com/spf13/cobra"
)

//ListCommands cobra command for listing the servers
func ListCommands(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the servers in the cluster",
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
