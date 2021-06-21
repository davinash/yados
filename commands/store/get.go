package store

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddGetStoreCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get store",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
