package object

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddGetObjectCommand Gets a object from a store
func AddGetObjectCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get object",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
