package server

import (
	"context"
	"fmt"
	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// AddServerListCmd Command for listing the members in the cluster
func AddServerListCmd(parentCmd *cobra.Command) {
	var serverAddress string
	var port int32
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all the members in cluster",

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
			listOfPeers, err := p.GetListOfPeers(context.Background(), &pb.ListOfPeersRequest{})
			if err != nil {
				return err
			}
			pj := protojson.MarshalOptions{
				Multiline: true,
				Indent:    "  ",
			}
			marshal, err := pj.Marshal(listOfPeers)
			if err != nil {
				return err
			}
			fmt.Println(string(marshal))
			return nil
		},
	}

	cmd.Flags().StringVar(&serverAddress, "server", "127.0.0.1", "IP address or host name")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	parentCmd.AddCommand(cmd)
}
