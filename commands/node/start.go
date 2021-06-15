package node

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddNodeStartCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a node",
		Args:  utils.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
