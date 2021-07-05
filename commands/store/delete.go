package store

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

//AddDeleteStoreCommand Cobra command for delete store command
func AddDeleteStoreCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete store",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
