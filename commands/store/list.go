package store

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddListStoreCommand Cobra command to list the stores
func AddListStoreCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all stores",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
