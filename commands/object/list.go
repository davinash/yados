package object

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddListObjectCommand Cobra command to list objects in a store
func AddListObjectCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all objects in a store",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
