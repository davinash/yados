package server

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddServerAddCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new server in the cluster",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
