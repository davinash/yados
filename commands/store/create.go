package store

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddCreateStoreCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create store",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
