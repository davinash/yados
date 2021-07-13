package object

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddDeleteObjectCommand Delete objects from the store
func AddDeleteObjectCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete object",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
