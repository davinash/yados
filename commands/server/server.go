package server

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddServerCmd Parent command server operation
func AddServerCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "server commands",
		Args:  utils.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	AddServerStartCmd(cmd)
	parentCmd.AddCommand(cmd)
}