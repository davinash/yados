package server

import (
	"fmt"
	"github.com/davinash/yados/commands/utils"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

func AddNodeStartCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "start NAME PORT",
		Short: "Start a server",
		Args:  utils.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(args)
			server, err := server.CreateNewServer(args[0], args[1])
			if err != nil {
				return err
			}
			return server.Start()
		},
	}
	parentCmd.AddCommand(cmd)
}
