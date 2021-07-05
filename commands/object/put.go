package object

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddPutObjectCommand Cobra comamnd to put new object in store
func AddPutObjectCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "put",
		Short: "put object",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
