package node

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddNodeStopCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a node",
		Args:  utils.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
