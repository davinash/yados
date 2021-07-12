package server

import (
	"github.com/davinash/yados/commands/utils"
	"github.com/spf13/cobra"
)

// AddServerStopCmd command line options for stopping a server
func AddServerStopCmd(parentCmd *cobra.Command) {
	var serverAddress string
	var port int
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "StopServerFn a server",
		Args:  utils.ExactArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			//	_, err := server.SendMessage(&server.MemberServer{
			//		Port:    port,
			//		Address: serverAddress,
			//	}, &server.Request{
			//		ID:        server.StopServer,
			//		Arguments: server.StopMember{},
			//	}, nil)
			//	if err != nil {
			//		return err
			//	}
			return nil
		},
	}
	cmd.Flags().StringVar(&serverAddress, "server", "127.0.0.1", "IP address or host name")
	cmd.Flags().IntVar(&port, "port", 9191, "Port to use for communication")
	_ = cmd.MarkFlagRequired("server")
	_ = cmd.MarkFlagRequired("port")

	parentCmd.AddCommand(cmd)
}
