package cli

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

func AddServerStatusCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Status of a server",
		Args:  utils.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	parentCmd.AddCommand(cmd)
}
