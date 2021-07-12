package server

import (
	"context"

	"github.com/davinash/yados/commands/utils"
	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// AddServerStopCmd command line options for stopping a server
func AddServerStopCmd(parentCmd *cobra.Command) {
	var serverAddress string
	var port int32
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a server",
		Args:  utils.ExactArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			conn, p, err := server.GetPeerConn(serverAddress, port)
			if err != nil {
				return err
			}
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
				}
			}(conn)
			_, err = p.StopServer(context.Background(), &pb.StopServerRequest{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&serverAddress, "server", "127.0.0.1", "IP address or host name")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	parentCmd.AddCommand(cmd)
}
