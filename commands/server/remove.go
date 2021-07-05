package server

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddServerRemoveCmd Command to remove the member from the cluster
func AddServerRemoveCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a server from the cluster",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
